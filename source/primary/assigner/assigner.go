package assigner

import (
	. "source/config"
	"source/localElevator/fsm"
	"time"
	"source/localElevator/elevio"
	"source/primary/sync"
)

func AssignRequests(requests Requests, wv Worldview, orderActionChan chan OrderAccess){
	for floor, request := range requests.Requests {
		for req, active := range request {
			if active {
				order := OrderConstructor(requests.Id, floor, req)
				if order.Button == int(elevio.BT_Cab) {
					sync.AddUnacceptedOrder(orderActionChan, order)
				} else {
					AssignedId := ChooseElevator(wv.FleetSnapshot, wv.PeerInfo.Peers, order)
					sync.AddUnacceptedOrder(orderActionChan, OrderConstructor(AssignedId, order.Floor, order.Button))
				}
			}
		}
	}
}

func ChooseElevator(elevators map[string]Elevator, activeIds []string, NewOrder Order)string{
	
	bestTime := time.Hour //inf
	var bestId string
	
	for _,Id := range(activeIds){
		if !elevators[Id].Obstructed{
			pickupTime := fsm.TimeUntilPickup(elevators[Id],NewOrder)
			if pickupTime < bestTime{
				bestId = Id
				bestTime = pickupTime
			}
		}
	}
	return bestId
}

//Creates a copy of the elevator and simulates executing remaining orders
//NOT USED
// func TimeToIdle(elev Elevator) time.Duration {
// 	duration := time.Duration(0)
// 	// Determines initial state
// 	switch elev.State {
// 	case IDLE:
// 		elev.Direction = fsm.ChooseDirection(elev)
// 		if elev.Direction == STOP {
// 			return duration
// 		}
// 	case MOVING:
// 		duration += T_TRAVEL / 2
// 		elev.Floor += int(elev.Direction)
// 	case DOOR_OPEN:
// 		duration -= T_DOOR_OPEN / 2
// 	}
	
// 	//Simulates remaining orders
// 	for {
// 		if fsm.ShouldStop(elev) {
// 			requests.ClearOrder(&elev, elev.Floor) //Changes do not propagate back to main
// 			duration += T_DOOR_OPEN
// 			elev.Direction = fsm.ChooseDirection(elev)
// 			if elev.Direction == STOP {
// 				return duration
// 			}
// 		}
// 		elev.Floor += int(elev.Direction)
// 		duration += T_TRAVEL
// 	}
// }