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
				if order.Button == int(elevio.BT_Cab){
					if !containsOrder(orderActionChan, order){
						sync.AddUnacceptedOrder(orderActionChan, order)
					}
				} else {
					assignedId := ChooseElevator(wv.FleetSnapshot, wv.PeerInfo.Peers, order)
					order = OrderConstructor(assignedId, order.Floor, order.Button)
					if !containsOrder(orderActionChan, order) {
						sync.AddUnacceptedOrder(orderActionChan, order)
					}
				}
			}
		}
	}
}

func containsOrder(orderActionChan chan OrderAccess, order Order) bool {
	orders := sync.GetUnacceptedOrders(orderActionChan)[order.Id]
    for _, ord := range orders {
        if ord == order {
            return true // Found the order
        }
    }
    return false // Order not found
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