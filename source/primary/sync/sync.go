package sync

import (
	. "source/config"
)

func FleetAccessManager(mapActionChan <-chan FleetAccess) {
	fleet := make(map[string]Elevator) // Real fleet map. All others are snapshots of this
	for {
		select {
		case newAction := <-mapActionChan:
			switch newAction.Cmd {
			case "read":
				deepCopy := make(map[string]Elevator, len(fleet))
				for key, value := range fleet {
					deepCopy[key] = value
				}
				newAction.ReadChan <- deepCopy
			case "write one":
				fleet[newAction.Id] = newAction.Elev
			case "write all":
				fleet = newAction.ElevMap
			}
		}
	}
}

func SingleFleetWrite(id string, elev Elevator, mapActionChan chan FleetAccess) {
	mapActionChan <- FleetAccess{Cmd: "write one", Id: id, Elev: elev}
}

func FullFleetWrite(elevMap map[string]Elevator, mapActionChan chan FleetAccess) {
	mapActionChan <- FleetAccess{Cmd: "write all", ElevMap: elevMap}
}

func FleetRead(mapActionChan chan FleetAccess) map[string]Elevator {
	readChan := make(chan map[string]Elevator, 1)
	defer close(readChan)
	mapActionChan <- FleetAccess{Cmd: "read", ReadChan: readChan}
	return <-readChan
}

func UnacceptedOrdersManager(ordersActionChan <-chan OrderAccess) {
	orders := make(map[string][]Order) //The true map of unaccpeted orders
	for {
		select {
		case action := <-ordersActionChan:
			switch action.Cmd {
			case "read":
				deepCopy := make(map[string][]Order, len(orders))
				for key, value := range orders {
					deepCopy[key] = append([]Order{}, value...)
				}
				action.ReadChan <- deepCopy

			case "write":
				orders[action.Id] = append(orders[action.Id], action.Orders...)

			case "delete":
				if existingOrders, exists := orders[action.Id]; exists {
					newOrders := []Order{}
					for _, o := range existingOrders {
						// Keep only orders that don't match the given order
						if !(o.Floor == action.Orders[0].Floor && o.Button == action.Orders[0].Button) {
							newOrders = append(newOrders, o)
						}
					}
					// If no orders remain, delete the key from the map
					if len(newOrders) > 0 {
						orders[action.Id] = newOrders
					// } else {
					// 	delete(orders, action.Id)
					}
				}
			}
		}
	}
}

func AddUnacceptedOrder(ordersActionChan chan<- OrderAccess, order Order) {
	ordersActionChan <- OrderAccess{
		Cmd:    "write",
		Id:     order.Id,
		Orders: []Order{order}, // Send a single order as a slice
	}
}

func GetUnacceptedOrders(ordersActionChan chan<- OrderAccess) map[string][]Order {
	readChan := make(chan map[string][]Order)
	defer close(readChan)

	ordersActionChan <- OrderAccess{
		Cmd:      "read",
		ReadChan: readChan,
	}
	return <-readChan
}

func RemoveUnacceptedOrder(ordersActionChan chan<- OrderAccess, order Order) {
	ordersActionChan <- OrderAccess{
		Cmd:    "delete",
		Id:     order.Id,
		Orders: []Order{order},
	}
}

func HallLightsManager(lightsActionChan <-chan LightsAccess) {
	hallLights := HallMatrix{}
	for {
		select {
		case action := <-lightsActionChan:
			switch action.Cmd {
			case "read":
				action.ReadChan <- hallLights
			case "write":
				hallLights = action.NewHallLights
			}
		}
	}
}

func ReadHallLights(lightsActionChan chan LightsAccess) HallMatrix {
	readChan := make(chan HallMatrix)
	defer close(readChan)

	lightsActionChan <- LightsAccess{
		Cmd:      "read",
		ReadChan: readChan,
	}
	return <-readChan
}

func WriteHallLights(lightsActionChan chan LightsAccess, newHallLights HallMatrix) {
	lightsActionChan <- LightsAccess{
		Cmd:           "write",
		NewHallLights: newHallLights,
	}
}