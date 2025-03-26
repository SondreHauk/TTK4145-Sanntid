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
	worldviewRXChan <-chan Worldview,
	requestsRXChan <-chan Requests,
	myId string) {

	// Syncronization channels
	elevatorsActionChan := make(chan ElevatorsAccess, 10)
	orderActionChan := make(chan OrderAccess, 10)
	lightsActionChan := make(chan LightsAccess, 10)

	elevUpdateObsChan := make(chan Elevator, NUM_ELEVATORS)
	worldviewObsChan := make(chan Worldview, 10)

	// Shared variables
	var worldview Worldview
	var latestPeerUpdate PeerUpdate
	worldview.FleetSnapshot = make(map[string]Elevator)
	worldview.UnacceptedOrdersSnapshot = make(map[string][]Order)
	hallLights := HallMatrixConstructor()

	// Owns and handles access shared variables
	go sync.ElevatorsAccessManager(elevatorsActionChan)
	go sync.UnacceptedOrdersManager(orderActionChan)
	go sync.HallLightsManager(lightsActionChan)

	go obstructionHandler(elevUpdateObsChan, worldviewObsChan, elevatorsActionChan, orderActionChan)

	for {
		select {
		// Draining of channels prior to primary activation
		case <-worldviewRXChan:
		case <-elevStateChan:
		case <-requestsRXChan:
		// case latestPeerUpdate = <-peerUpdateChan:
		// Primary activation
		case wv := <-becomePrimaryChan:
			fmt.Println("Taking over as Primary")
			worldview = wv
			worldview.PeerInfo = latestPeerUpdate
			printPeers(worldview.PeerInfo)
			sync.AllElevatorsWrite(worldview.FleetSnapshot, elevatorsActionChan)
			sync.WriteHallLights(lightsActionChan, wv.HallLightsSnapshot)
			heartbeatTimer := time.NewTicker(T_HEARTBEAT)
			defer heartbeatTimer.Stop()

		primaryLoop:
			for {
				select {
				case <-becomePrimaryChan: //drain

				case worldview.PeerInfo = <-peerUpdateChan:
					printPeers(worldview.PeerInfo)
					lost := worldview.PeerInfo.Lost
					if len(lost) != 0 {
						fmt.Println("Reassign and remember")
						reassignHallOrders(worldview, elevatorsActionChan,
							orderActionChan, Reassignment{Cause: Disconnected})
						rememberLostCabOrders(lost, orderActionChan, worldview)
					}

				case elevUpdate := <-elevStateChan:
					sync.SingleElevatorWrite(elevUpdate.Id, elevUpdate, elevatorsActionChan)
					unacceptedOrders := sync.GetUnacceptedOrders(orderActionChan)[elevUpdate.Id]
					checkforAcceptedOrders(orderActionChan, elevUpdate, unacceptedOrders)
					updateHallLights(worldview, hallLights, elevatorsActionChan, lightsActionChan)
					elevUpdateObsChan <- elevUpdate

				case requests := <-requestsRXChan:
					worldview.FleetSnapshot = sync.ElevatorsRead(elevatorsActionChan)
					worldview.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
					assigner.AssignRequests(requests, worldview, orderActionChan)

				case <-heartbeatTimer.C:
					worldview.FleetSnapshot = sync.ElevatorsRead(elevatorsActionChan)
					worldview.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
					worldview.HallLightsSnapshot = sync.ReadHallLights(lightsActionChan)
					worldviewTXChan <- worldview
					worldviewObsChan <- worldview
					// PrintWorldview(worldview)

				case receivedWV := <-worldviewRXChan:
					if receivedWV.PrimaryId < myId {
						fmt.Printf("Primary: %s, taking over\n", receivedWV.PrimaryId)
						fmt.Println("Enter Backup mode - listening for primary")
						break primaryLoop
					}
				}
			}
		}
	}
}
