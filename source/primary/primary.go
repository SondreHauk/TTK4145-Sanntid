package primary

import (
	"fmt"
	. "source/config"
	"source/localElevator/elevio"
	"source/primary/assigner"
	"source/primary/sync"
	"time"
)

func Run(
	peerUpdateChan <-chan PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimaryChan <-chan Worldview,
	worldviewTXChan chan<- Worldview,
	worldviewRXChan <-chan Worldview,
	requestFromElevChan <-chan Order,
	orderToElevChan chan<- Order,
	/*hallLightsChan chan<- [][]bool,*/
	myId string) {

	// local channels
	fleetActionChan := make(chan FleetAccess, 10)
	orderActionChan := make(chan OrderAccess, 10)
	elevUpdateObsChan := make(chan Elevator, NUM_ELEVATORS)
	worldviewObsChan := make(chan Worldview, 10)

	// local variables
	var worldview Worldview
	worldview.FleetSnapshot = make(map[string]Elevator)
	worldview.UnacceptedOrdersSnapshot = make(map[string][]Order)

	//Init hallLights matrix
	hallLights := make([][]bool, NUM_FLOORS)
	for i := range hallLights {
		hallLights[i] = make([]bool, NUM_BUTTONS-1)
	}

	//Owns and handles access to maps
	go sync.FleetAccessManager(fleetActionChan)
	go sync.UnacceptedOrdersManager(orderActionChan)
	go obstructionHandler(elevUpdateObsChan, worldviewObsChan, fleetActionChan, orderToElevChan)

	select {
	case wv := <-becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		worldview = wv
		//drain(elevStateChan) //FIX FLUSHING OF CHANNELS(?)
		sync.FullFleetWrite(worldview.FleetSnapshot,fleetActionChan)
		heartbeatTimer := time.NewTicker(T_HEARTBEAT)
		defer heartbeatTimer.Stop()
    	
		//primaryLoop:
		for {
			select {
			case worldview.PeerInfo = <-peerUpdateChan:
				//If elev lost: Reassign lost orders
				printPeers(worldview.PeerInfo)
				lost := worldview.PeerInfo.Lost
				if len(lost) != 0 {
					ReassignHallOrders(worldview, fleetActionChan, orderToElevChan, Reassignment{Cause: Disconnected})
				}

			case elevUpdate := <-elevStateChan:
				sync.SingleFleetWrite(elevUpdate.Id,elevUpdate,fleetActionChan)
				// if accepted order in elevupdate matches unaccpted order in worldview,
				// remove order from unaccepted orders and update light matrix corresondingly.
				//has a race condition but works fine
				updateHallLights(worldview, hallLights, fleetActionChan, /*hallLightsChan*/)
				//Obstruction handler gets updated states
				elevUpdateObsChan <- elevUpdate

			case request := <-requestFromElevChan:
				worldview.FleetSnapshot=sync.FleetRead(fleetActionChan)
				AssignedId := assigner.ChooseElevator(worldview.FleetSnapshot, worldview.PeerInfo.Peers, request)
				orderToElevChan <- OrderConstructor(AssignedId, request.Floor, request.Button)
				fmt.Printf("Assigned elevator %s to order\n", AssignedId)

			case <-heartbeatTimer.C:
				worldview.FleetSnapshot=sync.FleetRead(fleetActionChan)
				worldviewTXChan <- worldview
				worldviewObsChan <- worldview

			/* case receivedWV := <-worldviewRXChan:
				receivedId := receivedWV.PrimaryId
				fmt.Print(receivedId)
				if receivedId < myId {
					fmt.Printf("Primary: %s, taking over\n", receivedId)
					break primaryLoop */
			 //defere break om mulig?
			}
		}
	}
}


/* MAYBE implement function that owns hallLight state to avoid "trivial" race condition. Would be similar to fleetAccessManager
NOT 1st priority.  */

func updateHallLights(wv Worldview, hallLights [][]bool, mapActionChan chan FleetAccess, /*hallLightsChan chan<- [][]bool*/) {
	shouldUpdate := false
	prevHallLights := make([][]bool, NUM_FLOORS)
	for floor := range hallLights {
		prevHallLights[floor] = make([]bool, NUM_BUTTONS-1)
		copy(prevHallLights[floor], hallLights[floor]) // Copy row data
		for btn := range NUM_BUTTONS - 1 {
			hallLights[floor][btn] = false
		}
	}
	wv = WorldviewConstructor(wv.PrimaryId, wv.PeerInfo, sync.FleetRead(mapActionChan))
	for _, id := range wv.PeerInfo.Peers {
		orderMatrix := wv.FleetSnapshot[id].Orders
		for floor, floorOrders := range orderMatrix {
			for btn, isOrder := range floorOrders {
				if isOrder && btn != int(elevio.BT_Cab) {
					hallLights[floor][btn] = true
				}
			}
		}
	}
	
  for floor := range NUM_FLOORS {
		for btn := range NUM_BUTTONS - 1 {
			if prevHallLights[floor][btn] != hallLights[floor][btn] {
				shouldUpdate = true

			}
		}
	}
	if shouldUpdate {
		hallLightsChan <- hallLights
	}
}

func ReassignHallOrders(wv Worldview, MapActionChan chan FleetAccess, orderToElevChan chan<- Order, reassign Reassignment){
	wv = WorldviewConstructor(wv.PrimaryId, wv.PeerInfo, sync.FleetRead(MapActionChan))
	switch reassign.Cause{
	case Disconnected:
		for _, lostId := range wv.PeerInfo.Lost {
		orderMatrix := wv.FleetSnapshot[lostId].Orders
		for floor, floorOrders := range orderMatrix {
			for btn, isOrder := range floorOrders {
				if isOrder && btn != int(elevio.BT_Cab) {
					lostOrder := Order{
						Id:     lostId,
						Floor:  floor,
						Button: btn,
					}
					lostOrder.Id = assigner.ChooseElevator(wv.FleetSnapshot, wv.PeerInfo.Peers, lostOrder)
					orderToElevChan <- lostOrder
				}
			}
		}
	}
	case Obstructed:
		orderMatrix := wv.FleetSnapshot[reassign.ObsId].Orders
		for floor, floorOrders := range(orderMatrix){
			for btn, isOrder := range(floorOrders){
			if isOrder && btn != int(elevio.BT_Cab){
				lostOrder:=Order{
					Id: reassign.ObsId,
					Floor: floor,
					Button: btn,
					}
				lostOrder.Id = assigner.ChooseElevator(wv.FleetSnapshot, wv.PeerInfo.Peers, lostOrder)
				orderToElevChan <- lostOrder
			}
			}
		}
	}
}

func printPeers(p PeerUpdate) {
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}

func obstructionHandler(
	elevUpdateObsChan chan Elevator,
	worldviewObsChan chan Worldview, 
	mapActionChan chan FleetAccess,
	orderToElevChan chan<- Order,
	){
	obstructedElevators := make([]string, NUM_ELEVATORS)
	obstructionTimers := make(map[string]*time.Timer)
	var worldview Worldview
	var elevUpdate Elevator
	for{
		select{
		case worldview = <-worldviewObsChan:
		case elevUpdate = <-elevUpdateObsChan:
			if elevUpdate.Obstructed {
				obstructedElevators = append(obstructedElevators, elevUpdate.Id)
				//If no timer, start one
				_, timerExists := obstructionTimers[elevUpdate.Id]
				if !timerExists{
					timer := time.AfterFunc(T_REASSIGN_PRIMARY, func() {
					reassignmentDetails := Reassignment{Cause: Obstructed, ObsId: obstructedElevators[len(obstructedElevators)-1]}
					ReassignHallOrders(worldview, mapActionChan, orderToElevChan, reassignmentDetails)})
					obstructionTimers[elevUpdate.Id] = timer
				}
			} else {
				//if ID in obstructedElevatorIds, pop id and stop timer
				// If the elevator is no longer obstructed, check if its ID is in the list of obstructed elevators
				for i, id := range obstructedElevators {
					if id == elevUpdate.Id {
						// If found, remove it from the slice
						obstructedElevators = append(obstructedElevators[:i], obstructedElevators[i+1:]...)
						//obstructedElevators = slices.Delete(obstructedElevators,i,i+1)
						// Stop the timer if it's active
						if timer, exists := obstructionTimers[elevUpdate.Id]; exists {
							timer.Stop()
							delete(obstructionTimers, elevUpdate.Id)
						}
						break
					}
				}
			}
		}
	}
} 
