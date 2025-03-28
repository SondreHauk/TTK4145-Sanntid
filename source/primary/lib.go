package primary

import (
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
	mapActionChan chan FleetAccess,
	lightsActionChan chan LightsAccess,
) {
	lights := HallMatrixConstructor()
	wv.FleetSnapshot = sync.FleetRead(mapActionChan)
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

// We are fully aware that this function can be shortened down quite a bit.
// There simply was not time left to prioritize this.
func reassignHallOrders(
	wv Worldview,
	MapActionChan chan FleetAccess,
	ordersActionChan chan OrderAccess,
	reassign Reassignment,
) {
	wv.FleetSnapshot = sync.FleetRead(MapActionChan)
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
		orderMatrix := wv.FleetSnapshot[reassign.Id].Orders
		for floor, floorOrders := range orderMatrix {
			for btn, isOrder := range floorOrders {
				if isOrder && btn != int(elevio.BT_Cab) {
					lostOrder := OrderConstructor(
						reassign.Id,
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
	case MotorStop:
		orderMatrix := wv.FleetSnapshot[reassign.Id].Orders
		for floor, floorOrders := range orderMatrix {
			for btn, isOrder := range floorOrders {
				if isOrder && btn != int(elevio.BT_Cab) {
					lostOrder := OrderConstructor(
						reassign.Id,
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
	MapActionChan chan FleetAccess,
) {
	wv.FleetSnapshot = sync.FleetRead(MapActionChan)
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
	mapActionChan chan FleetAccess,
	ordersActionChan chan OrderAccess,
) {
	// This slice length can be generalized like NUM_FLOORS
	obstructedElevs := make([]string, NUM_ELEVATORS)
	obstructionTimers := make(map[string]*time.Timer)
	var wv Worldview
	var elevUpdate Elevator
	for {
		select {
		case wv = <-wvObsChan:
		case elevUpdate = <-elevUpdateObsChan:
			if elevUpdate.MotorStop {
				reassignHallOrders(
					wv,
					mapActionChan,
					ordersActionChan,
					Reassignment{
						Cause: MotorStop,
						Id:    elevUpdate.Id,
					},
				)
			}
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
									Id:    obstructedElevs[len(obstructedElevs)-1],
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
