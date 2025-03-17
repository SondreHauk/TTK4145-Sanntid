package primary

import( 
	."source/config"
	"source/localElevator/elevio"
	"source/primary/sync"
	"source/primary/assigner"
	"time"
	"fmt"
)

func checkforAcceptedOrders(orderActionChan chan OrderAccess, elevUpdate Elevator, unacceptedOrders []Order){
	for floor, buttons := range elevUpdate.Orders {
		for btn, orderAccepted := range buttons {
			if orderAccepted {
				for _, unaccOrder := range unacceptedOrders {
					if unaccOrder.Floor == floor && unaccOrder.Button == btn {
						sync.RemoveUnacceptedOrder(orderActionChan, 
							Order{Id: elevUpdate.Id, Floor: floor, Button: btn})
						break
					}
				} 
			}
		}
	}
}

func updateHallLights(wv Worldview, hallLights HallLights, mapActionChan chan FleetAccess, lightsActionChan chan LightsAccess) {
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
	sync.WriteHallLights(lightsActionChan, hallLights)
}

func ReassignHallOrders(wv Worldview, MapActionChan chan FleetAccess, ordersActionChan chan OrderAccess, reassign Reassignment){
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
					// APPEND TO UNACCEPTED ORDERS IN WORLDVIEW
					sync.AddUnacceptedOrder(ordersActionChan, lostOrder)
					/*orderToElevChan <- lostOrder*/
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
				// APPEND TO UNACCEPTED ORDERS IN WORLDVIEW
				sync.AddUnacceptedOrder(ordersActionChan, lostOrder)
				/*orderToElevChan <- lostOrder*/
			}
			}
		}
	}
}

func obstructionHandler(
	elevUpdateObsChan chan Elevator,
	worldviewObsChan chan Worldview, 
	mapActionChan chan FleetAccess,
	ordersActionChan chan OrderAccess,
	/*orderToElevChan chan<- Order,*/
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
					ReassignHallOrders(worldview, mapActionChan,ordersActionChan, reassignmentDetails)})
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

func printPeers(p PeerUpdate) {
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}