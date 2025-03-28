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
	elevRXChan <-chan Elevator,
	enablePrimaryChan <-chan Worldview,
	wvTXChan chan<- Worldview,
	wvRXChan <-chan Worldview,
	requestsRXChan <-chan Requests,
	myId string,
) {
	// Syncronization channels
	fleetAccessChan := make(chan FleetAccess, 10)
	orderActionChan := make(chan OrderAccess, 10)
	lightsActionChan := make(chan LightsAccess, 10)
	elevUpdateObsChan := make(chan Elevator, NUM_ELEVATORS)
	wvObsChan := make(chan Worldview, 10)

	// Local variables
	var wv Worldview
	var latestPeerUpdate PeerUpdate

	// Owns and handles access to shared data structures
	go sync.FleetAccessManager(fleetAccessChan)
	go sync.UnacceptedOrdersManager(orderActionChan)
	go sync.HallLightsManager(lightsActionChan)

	go obstructionHandler(elevUpdateObsChan, wvObsChan, fleetAccessChan, orderActionChan)

	for {
		select {
		// Draining of channels prior to primary activation
		case <-wvRXChan:
		case <-elevRXChan:
		case <-requestsRXChan:

		// Primary activation
		case wv = <-enablePrimaryChan:
			fmt.Println("Taking over as Primary")
			wv.PeerInfo = latestPeerUpdate
			sync.FullFleetWrite(wv.FleetSnapshot, fleetAccessChan)
			sync.WriteHallLights(lightsActionChan, wv.HallLightsSnapshot)
			heartbeatTimer := time.NewTicker(T_HEARTBEAT)
			defer heartbeatTimer.Stop()

		primaryLoop:
			for {
				select {
				// Drain in case of enable -> disable -> enable
				case <-enablePrimaryChan:

				case wv.PeerInfo = <-peerUpdateChan:
					lost := wv.PeerInfo.Lost
					if len(lost) != 0 {
						reassignHallOrders(
							wv,
							fleetAccessChan,
							orderActionChan,
							Reassignment{Cause: Disconnected},
						)
						rememberLostCabOrders(
							lost,
							orderActionChan,
							wv,
							fleetAccessChan,
						)
					}

				case elevUpdate := <-elevRXChan:
					sync.SingleElevFleetWrite(
						elevUpdate.Id,
						elevUpdate,
						fleetAccessChan,
					)
					unacceptedOrders := sync.GetUnacceptedOrders(orderActionChan)[elevUpdate.Id]
					checkforAcceptedOrders(
						orderActionChan,
						elevUpdate,
						unacceptedOrders,
					)
					updateHallLights(
						wv,
						fleetAccessChan,
						lightsActionChan,
					)
					elevUpdateObsChan <- elevUpdate

				case requests := <-requestsRXChan:
					wv.FleetSnapshot = sync.FleetRead(fleetAccessChan)
					wv.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
					assigner.AssignRequests(requests, wv, orderActionChan)

				case <-heartbeatTimer.C:
					wv.FleetSnapshot = sync.FleetRead(fleetAccessChan)
					wv.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
					wv.HallLightsSnapshot = sync.ReadHallLights(lightsActionChan)
					wvTXChan <- wv
					wvObsChan <- wv

				case receivedWV := <-wvRXChan:
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
