package assigner

import (
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/primary/sync"
	"time"
)

func AssignRequests(
	requests Requests,
	wv Worldview,
	orderActionChan chan OrderAccess,
) {
	unaccOrders := wv.UnacceptedOrdersSnapshot
	for floor, request := range requests.Requests {
		for req, active := range request {
			if active {
				order := OrderConstructor(
					requests.Id,
					floor,
					req,
				)
				orders := unaccOrders[order.Id]
				if !containsOrder(orders, order) {
					if order.Button == int(elevio.BT_Cab) {
						sync.AddUnacceptedOrder(orderActionChan, order)
					} else {
						AssignedId := ChooseElevator(
							wv.FleetSnapshot,
							wv.PeerInfo.Peers,
							order,
						)
						unacceptedOrder := OrderConstructor(
							AssignedId,
							order.Floor,
							order.Button,
						)
						sync.AddUnacceptedOrder(
							orderActionChan,
							unacceptedOrder,
						)
					}
				}
			}
		}
	}
}

func containsOrder(orders []Order, order Order) bool {
	for _, ord := range orders {
		if ord == order {
			return true // Found the value
		}
	}
	return false // Value not found
}

func ChooseElevator(
	elevators map[string]Elevator,
	activeIds []string,
	NewOrder Order,
) string {
	bestTime := time.Hour //inf
	var bestId string
	for _, Id := range activeIds {
		if !elevators[Id].Obstructed {
			pickupTime := fsm.TimeUntilPickup(elevators[Id], NewOrder)
			if pickupTime < bestTime {
				bestId = Id
				bestTime = pickupTime
			}
		}
	}
	return bestId
}
