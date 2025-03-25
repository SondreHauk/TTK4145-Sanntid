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