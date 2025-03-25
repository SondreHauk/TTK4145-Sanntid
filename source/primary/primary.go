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
	// var lostCabOrders []Order

	// Owns and handles access to maps
	go sync.FleetAccessManager(fleetActionChan)
	go sync.UnacceptedOrdersManager(orderActionChan)
	go sync.HallLightsManager(lightsActionChan)
	go obstructionHandler(elevUpdateObsChan, worldviewObsChan, fleetActionChan, orderActionChan)
	
	for{
		select {
		// Draining of channels prior to primary activation
		case <-worldviewRXChan:
		case <-elevStateChan:
			// Drain channel
		case <-requestsRXChan:
			// Drain channel
		case wv := <-becomePrimaryChan:
			fmt.Println("Taking over as Primary")
			worldview = wv
			sync.FullFleetWrite(worldview.FleetSnapshot, fleetActionChan)
			sync.WriteHallLights(lightsActionChan, wv.HallLightsSnapshot)
			heartbeatTimer := time.NewTicker(T_HEARTBEAT)
			// defer heartbeatTimer.Stop()

			primaryLoop:
			for {
				select {
				case worldview.PeerInfo = <-peerUpdateChan:
					printPeers(worldview.PeerInfo)
					lost := worldview.PeerInfo.Lost
					if len(lost) != 0 {
						reassignHallOrders(worldview, fleetActionChan, 
							orderActionChan, Reassignment{Cause: Disconnected})
						rememberLostCabOrders(lost, orderActionChan, worldview)
					}

				case elevUpdate := <-elevStateChan:
					sync.SingleFleetWrite(elevUpdate.Id, elevUpdate, fleetActionChan)
					unacceptedOrders := sync.GetUnacceptedOrders(orderActionChan)[elevUpdate.Id]
					checkforAcceptedOrders(orderActionChan, elevUpdate, unacceptedOrders)
					updateHallLights(worldview, hallLights, fleetActionChan, lightsActionChan)
					elevUpdateObsChan <- elevUpdate
          
				case requests := <-requestsRXChan:
					worldview.FleetSnapshot = sync.FleetRead(fleetActionChan)
					assigner.AssignRequests(requests, worldview, orderActionChan)

				case <-heartbeatTimer.C:
					worldview.FleetSnapshot = sync.FleetRead(fleetActionChan)
					worldview.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
					worldview.HallLightsSnapshot = sync.ReadHallLights(lightsActionChan)
					worldviewTXChan <- worldview
					worldviewObsChan <- worldview

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
