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
	myId string) {

	// Local channels
	fleetActionChan := make(chan FleetAccess, 10)
	orderActionChan := make(chan OrderAccess, 10)
	lightsActionChan := make(chan LightsAccess, 10)

	elevUpdateObsChan := make(chan Elevator, NUM_ELEVATORS)
	worldviewObsChan := make(chan Worldview, 10)

	// Local variables
	var worldview Worldview
	worldview.FleetSnapshot = make(map[string]Elevator)
	worldview.UnacceptedOrdersSnapshot = make(map[string][]Order)
	hallLights := HallLights{}

	// Owns and handles access to maps
	go sync.FleetAccessManager(fleetActionChan)
	go sync.UnacceptedOrdersManager(orderActionChan)
	go sync.HallLightsManager(lightsActionChan)
	go obstructionHandler(elevUpdateObsChan, worldviewObsChan, fleetActionChan, orderActionChan)
	for{
		select {
		case <-elevStateChan:
			// Drain channel
		case <-requestFromElevChan:
			// Drain channel
		case wv := <-becomePrimaryChan:
			fmt.Println("Taking over as Primary")
			worldview = wv
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

				case request := <-requestFromElevChan:
					// fmt.Printf("Request received from: %s\n ", request.Id)
					worldview.FleetSnapshot = sync.FleetRead(fleetActionChan)
					AssignedId := assigner.ChooseElevator(worldview.FleetSnapshot, worldview.PeerInfo.Peers, request)
					sync.AddUnacceptedOrder(orderActionChan, OrderConstructor(AssignedId, request.Floor, request.Button))
					// APPEND TO UNACCEPTED ORDERS IN WORLDVIEW
					/*orderToElevChan <- OrderConstructor(AssignedId, request.Floor, request.Button)*/
					// fmt.Printf("Elevator %s assigned\n", AssignedId)

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
}
