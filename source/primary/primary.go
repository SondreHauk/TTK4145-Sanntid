package primary

import (
	"fmt"
	. "source/config"
	"source/primary/sync"
	"time"
)

func Run(
	peerUpdateChan <-chan PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimaryChan <-chan Worldview,
	worldviewTXChan chan<- Worldview,
	/*worldviewRXChan <-chan Worldview,*/
	requestFromElevChan <-chan HallMatrix,
	myId string) {

	// Syncronization channels
	fleetActionChan := make(chan FleetAccess, 10)
	orderActionChan := make(chan OrderAccess, 10)
	lightsActionChan := make(chan LightsAccess, 10)

	elevUpdateObsChan := make(chan Elevator, NUM_ELEVATORS)
	worldviewObsChan := make(chan Worldview, 10)

	// Local variables
	var worldview Worldview
	worldview.FleetSnapshot = make(map[string]Elevator)
	worldview.UnacceptedOrdersSnapshot = make(map[string][]Order)
	hallLights := HallMatrix{}

	// Owns and handles access to maps
	go sync.FleetAccessManager(fleetActionChan)
	go sync.UnacceptedOrdersManager(orderActionChan)
	go sync.HallLightsManager(lightsActionChan)
	go obstructionHandler(elevUpdateObsChan, worldviewObsChan, fleetActionChan, orderActionChan)

	select {
	case wv := <-becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		worldview = wv
		// TODO: FIX FLUSHING/ROUTING OF CHANNELS
		sync.FullFleetWrite(worldview.FleetSnapshot, fleetActionChan)
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
					ReassignHallOrders(worldview, fleetActionChan, 
						orderActionChan, Reassignment{Cause: Disconnected})
				}

			case elevUpdate := <-elevStateChan:
				sync.SingleFleetWrite(elevUpdate.Id, elevUpdate, fleetActionChan)
				unacceptedOrders := sync.GetUnacceptedOrders(orderActionChan)[elevUpdate.Id]
				checkforAcceptedOrders(orderActionChan, elevUpdate, unacceptedOrders)
				updateHallLights(worldview, hallLights, fleetActionChan, lightsActionChan)
				elevUpdateObsChan <- elevUpdate

			case requests := <-requestFromElevChan:
				// fmt.Printf("Request received from: %s\n ", request.Id)
				worldview.FleetSnapshot = sync.FleetRead(fleetActionChan) // Should this be done each time and not once?
				// TODO: extract requests into order format and run the following for each order:

				assignRequests(requests, worldview, orderActionChan)
				// for floor, request := range requests {
				// 	for btn, active := range request {
				// 		if active {
				// 			order := OrderConstructor("arbitrary", floor, btn) // Id or not?
				// 			AssignedId := assigner.ChooseElevator(worldview.FleetSnapshot, worldview.PeerInfo.Peers, order)
				// 			sync.AddUnacceptedOrder(orderActionChan, OrderConstructor(AssignedId, order.Floor, order.Button))
				// 		}
				// 	}
				// }

			case <-heartbeatTimer.C:
				worldview.FleetSnapshot = sync.FleetRead(fleetActionChan)
				worldview.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
				worldview.HallLightsSnapshot = sync.ReadHallLights(lightsActionChan)
				// PrintWorldView(worldview)
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
