package primary

import (
	"fmt"
	. "source/config"
	"source/localElevator/elevio"
	"source/primary/assigner"
	"source/primary/sync"
	"time"
)

func checkforAcceptedOrders(
	orderActionChan chan OrderAccess,
	elevUpdate Elevator,
	unacceptedOrders []Order,
) {
	for floor, buttons := range elevUpdate.Orders {
		for btn, elevAcceptsOrder := range buttons {
			if elevAcceptsOrder {
				for _, unacceptedOrder := range unacceptedOrders {
					if unacceptedOrder.Floor == floor && unacceptedOrder.Button == btn {
						acceptedOrder := OrderConstructor(
							elevUpdate.Id,
							floor,
							btn,
						)
						sync.AcceptOrder(orderActionChan, acceptedOrder)
						break
					}
				}
			}
		}
	}
}

func updateHallLights(
	wv Worldview,
	lights HallMatrix,
	mapActionChan chan ElevatorsAccess,
	lightsActionChan chan LightsAccess,
) {
	wv.FleetSnapshot = sync.ElevatorsRead(mapActionChan)
	for _, id := range wv.PeerInfo.Peers {
		orderMatrix := wv.FleetSnapshot[id].Orders
		for floor, floorOrders := range orderMatrix {
			for btn, isOrder := range floorOrders {
				if isOrder && btn != int(elevio.BT_Cab) {
					lights[floor][btn] = true
				}
			}
		}
	}
	sync.WriteHallLights(lightsActionChan, lights)
}

func reassignHallOrders(
	wv Worldview,
	MapActionChan chan ElevatorsAccess,
	ordersActionChan chan OrderAccess,
	reassign Reassignment,
) {
	wv.FleetSnapshot = sync.ElevatorsRead(MapActionChan)
	switch reassign.Cause {
	case Disconnected:
		for _, lostId := range wv.PeerInfo.Lost {
			orderMatrix := wv.FleetSnapshot[lostId].Orders
			for floor, floorOrders := range orderMatrix {
				for btn, isOrder := range floorOrders {
					if isOrder && btn != int(elevio.BT_Cab) {
						lostOrder := OrderConstructor(
							lostId,
							floor,
							btn,
						)
						newId := assigner.ChooseElevator(
							wv.FleetSnapshot,
							wv.PeerInfo.Peers,
							lostOrder,
						)
						lostOrder.Id = newId
						sync.AddUnacceptedOrder(ordersActionChan, lostOrder)
					}
				}
			}
		}
	case Obstructed:
		orderMatrix := wv.FleetSnapshot[reassign.ObsId].Orders
		for floor, floorOrders := range orderMatrix {
			for btn, isOrder := range floorOrders {
				if isOrder && btn != int(elevio.BT_Cab) {
					lostOrder := OrderConstructor(
						reassign.ObsId,
						floor,
						btn,
					)
					newId := assigner.ChooseElevator(
						wv.FleetSnapshot,
						wv.PeerInfo.Peers,
						lostOrder,
					)
					lostOrder.Id = newId
					sync.AddUnacceptedOrder(ordersActionChan, lostOrder)
				}
			}
		}
	}
}

func rememberLostCabOrders(
	lostElevators []string,
	orderActionChan chan OrderAccess,
	wv Worldview,
	MapActionChan chan ElevatorsAccess,
) {
	//HAR LAGT TIL SYNCING AV SNAPSHOT. TRUR DET BLIR RIKTIG?
	wv.FleetSnapshot = sync.ElevatorsRead(MapActionChan)
	for _, id := range lostElevators {
		for floor, orders := range wv.FleetSnapshot[id].Orders {
			for btn, active := range orders {
				if active && btn == int(elevio.BT_Cab) {
					sync.AddUnacceptedOrder(orderActionChan,
						OrderConstructor(id, floor, btn))
				}
			}
		}
	}
}

func obstructionHandler(
	elevUpdateObsChan chan Elevator,
	wvObsChan chan Worldview,
	mapActionChan chan ElevatorsAccess,
	ordersActionChan chan OrderAccess,
) {
	// HER ER DENNE SATT TIL NUM_ELEVATORS
	obstructedElevs := make([]string, NUM_ELEVATORS)
	obstructionTimers := make(map[string]*time.Timer)
	var wv Worldview
	var elevUpdate Elevator
	for {
		select {
		case wv = <-wvObsChan:
		case elevUpdate = <-elevUpdateObsChan:
			if elevUpdate.Obstructed {
				obstructedElevs = append(obstructedElevs, elevUpdate.Id)
				_, timerExists := obstructionTimers[elevUpdate.Id]
				if !timerExists {
					timer := time.AfterFunc(T_REASSIGN_PRIMARY,
						func() {
							reassignHallOrders(
								wv,
								mapActionChan,
								ordersActionChan,
								Reassignment{
									Cause: Obstructed,
									ObsId: obstructedElevs[len(obstructedElevs)-1],
								},
							)
						},
					)
					obstructionTimers[elevUpdate.Id] = timer
				}
			} else {
				for i, id := range obstructedElevs {
					if id == elevUpdate.Id {
						obstructedElevs = append(obstructedElevs[:i], obstructedElevs[i+1:]...)
						timer, timerExists := obstructionTimers[elevUpdate.Id]
						if timerExists {
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

/* func PrintWorldview(wv Worldview) {
	fmt.Println("--- Worldview Snapshot ---")
	fmt.Println("PrimaryId:", wv.PrimaryId)
	fmt.Println("Peers:", wv.PeerInfo.Peers)
	fmt.Println("New Peer:", wv.PeerInfo.New)
	fmt.Println("Lost Peers:", wv.PeerInfo.Lost)
	fmt.Println("Fleet Snapshot:")
	 for id, elev := range wv.FleetSnapshot {
	 	fmt.Printf("  Elevator ID: %s\n", id)
	 	fmt.Printf("    Floor: %d, Direction: %d, PrevDirection: %d, State: %d\n",
	 		elev.Floor, elev.Direction, elev.PrevDirection, elev.State)
	 	fmt.Printf("    Obstructed: %t\n", elev.Obstructed)
	 	fmt.Println("    Orders:")
	 	for i := 0; i < NUM_FLOORS; i++ {
	 		fmt.Printf("      Floor %d: %v\n", i, elev.Orders[i])
	 	}
	 }
	 fmt.Println("Unaccepted Orders Snapshot:")
	for id, orders := range wv.UnacceptedOrdersSnapshot {
		fmt.Printf("  Orders for Elevator %s:\n", id)
		for _, order := range orders {
			fmt.Printf("    Floor: %d, Button: %d\n", order.Floor, order.Button)
		}
	}
	 fmt.Println("Hall Lights Snapshot:")
	 for i := 0; i < NUM_FLOORS; i++ {
	 	fmt.Printf("  Floor %d: %v\n", i, wv.HallLightsSnapshot[i])
	 }
	 fmt.Println("-------------------------")
} */
