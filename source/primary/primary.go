package primary

import (
	"fmt"
	. "source/config"
	"source/primary/assigner"
	"source/primary/sync"
	"time"
)

func Run(
	peerUpdateChan <-chan PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimaryChan <-chan Worldview,
	worldviewTXChan chan<- Worldview,
	/*worldviewRXChan <-chan Worldview,*/
	requestFromElevChan <-chan Order,
	/*orderToElevChan chan<- Order,*/
	/*hallLightsChan chan<- [][]bool,*/
	myId string) {

	// local channels
	fleetActionChan := make(chan FleetAccess, 10)
	orderActionChan := make(chan OrderAccess, 10)
	lightsActionChan := make(chan LightsAccess, 10)

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
	go sync.HallLightsManager(lightsActionChan)
	go obstructionHandler(elevUpdateObsChan, worldviewObsChan, fleetActionChan, orderActionChan)

	select {
	case wv := <-becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		worldview = wv
		// TODO: FIX FLUSHING/ROUTING OF CHANNELS
		sync.FullFleetWrite(worldview.FleetSnapshot,fleetActionChan)
		sync.WriteHallLights(lightsActionChan, wv.HallLightsSnapshot)
		heartbeatTimer := time.NewTicker(T_HEARTBEAT)
		defer heartbeatTimer.Stop()
    	
		//primaryLoop:
		for {
			select {
			case worldview.PeerInfo = <-peerUpdateChan:
				printPeers(worldview.PeerInfo)
				lost := worldview.PeerInfo.Lost
				if len(lost) != 0 {
					ReassignHallOrders(worldview, fleetActionChan, orderActionChan, Reassignment{Cause: Disconnected})
				}

			case elevUpdate := <-elevStateChan:
				sync.SingleFleetWrite(elevUpdate.Id,elevUpdate,fleetActionChan)
				unacceptedOrders := sync.GetUnacceptedOrder(orderActionChan, elevUpdate.Id)
				checkforAcceptedOrders(orderActionChan, elevUpdate, unacceptedOrders)
				updateHallLights(worldview, hallLights, fleetActionChan, lightsActionChan)
				elevUpdateObsChan <- elevUpdate

			case request := <-requestFromElevChan:
				fmt.Printf("Request received from: %s\n ", request.Id)
				worldview.FleetSnapshot=sync.FleetRead(fleetActionChan)
				AssignedId := assigner.ChooseElevator(worldview.FleetSnapshot, worldview.PeerInfo.Peers, request)
				sync.AddUnacceptedOrder(orderActionChan, OrderConstructor(AssignedId, request.Floor, request.Button))
				// APPEND TO UNACCEPTED ORDERS IN WORLDVIEW
				/*orderToElevChan <- OrderConstructor(AssignedId, request.Floor, request.Button)*/
				fmt.Printf("Elevator %s assigned\n", AssignedId)

			case <-heartbeatTimer.C:
				worldview.FleetSnapshot = sync.FleetRead(fleetActionChan)
				worldview.UnacceptedOrdersSnapshot = sync.GetAllUnacceptedOrders(orderActionChan)
				worldview.HallLightsSnapshot = sync.ReadHallLights(lightsActionChan)

				// matrix := worldview.HallLightsSnapshot
				// // Print the slice matrix
				// fmt.Println("4x2 Boolean Matrix:")
				// for _, row := range matrix {
				// 	for _, val := range row {
				// 		if val {
				// 			fmt.Printf("1 ")
				// 		} else {
				// 			fmt.Printf("0 ")
				// 		}
				// 	}
				// 	fmt.Println() // New line after each row
				// }

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