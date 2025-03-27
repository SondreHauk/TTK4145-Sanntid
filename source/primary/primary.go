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
	elevsAccessChan := make(chan ElevatorsAccess, 10)
	orderActionChan := make(chan OrderAccess, 10)
	lightsActionChan := make(chan LightsAccess, 10)
	elevUpdateObsChan := make(chan Elevator, NUM_ELEVATORS)
	wvObsChan := make(chan Worldview, 10)

	// Shared variables
	var wv Worldview
	var latestPeerUpdate PeerUpdate
/* 	hallLights := HallMatrixConstructor() */

	// Owns and handles access shared variables
	go sync.ElevatorsAccessManager(elevsAccessChan)
	go sync.UnacceptedOrdersManager(orderActionChan)
	go sync.HallLightsManager(lightsActionChan)

	go obstructionHandler(elevUpdateObsChan, wvObsChan, elevsAccessChan, orderActionChan)

	for {
		select {
		// Draining of channels prior to primary activation
		case <-wvRXChan:
		case <-elevRXChan:
		case <-requestsRXChan:
		case latestPeerUpdate = <-peerUpdateChan:

		// Primary activation
		case wv = <-enablePrimaryChan:
			fmt.Println("Taking over as Primary")
			wv.PeerInfo = latestPeerUpdate
			printPeers(wv.PeerInfo)
			sync.AllElevatorsWrite(wv.FleetSnapshot, elevsAccessChan)
			sync.WriteHallLights(lightsActionChan, wv.HallLightsSnapshot)
			heartbeatTimer := time.NewTicker(T_HEARTBEAT)
			defer heartbeatTimer.Stop()

		primaryLoop:
			for {
				select {
				case <-enablePrimaryChan: //drain

				case wv.PeerInfo = <-peerUpdateChan:
					printPeers(wv.PeerInfo)
					lost := wv.PeerInfo.Lost
					if len(lost) != 0 {
						fmt.Println("Reassign and remember")
						reassignHallOrders(
							wv,
							elevsAccessChan,
							orderActionChan,
							Reassignment{Cause: Disconnected},
						)
						rememberLostCabOrders(
							lost,
							orderActionChan,
							wv,
							elevsAccessChan,
						)
					}

				case elevUpdate := <-elevRXChan:
					sync.SingleElevatorWrite(
						elevUpdate.Id,
						elevUpdate,
						elevsAccessChan,
					)
					unacceptedOrders := sync.GetUnacceptedOrders(orderActionChan)[elevUpdate.Id]
					checkforAcceptedOrders(
						orderActionChan,
						elevUpdate,
						unacceptedOrders,
					)
					updateHallLights(
						wv,
						/* hallLights, */
						elevsAccessChan,
						lightsActionChan,
					)
					elevUpdateObsChan <- elevUpdate

				case requests := <-requestsRXChan:
					wv.FleetSnapshot = sync.ElevatorsRead(elevsAccessChan)
					wv.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
					assigner.AssignRequests(requests, wv, orderActionChan)

				case <-heartbeatTimer.C:
					wv.FleetSnapshot = sync.ElevatorsRead(elevsAccessChan)
					wv.UnacceptedOrdersSnapshot = sync.GetUnacceptedOrders(orderActionChan)
					wv.HallLightsSnapshot = sync.ReadHallLights(lightsActionChan)
					wvTXChan <- wv
					wvObsChan <- wv
					// PrintWorldview(wv)

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
