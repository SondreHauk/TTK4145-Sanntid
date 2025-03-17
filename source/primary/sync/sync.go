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

func SingleFleetWrite(id string, elev Elevator, mapActionChan chan FleetAccess){
	mapActionChan<-FleetAccess{Cmd:"write one", Id:id, Elev:elev}
}

func FullFleetWrite(elevMap map[string]Elevator, mapActionChan chan FleetAccess){
	mapActionChan<-FleetAccess{Cmd:"write all", ElevMap: elevMap}
}

func FleetRead(mapActionChan chan FleetAccess) map[string]Elevator{
	readChan := make(chan map[string]Elevator, 1)
	defer close(readChan)
	mapActionChan<-FleetAccess{Cmd:"read", ReadChan:readChan}
	return <- readChan
}

func UnacceptedOrdersManager(ordersActionChan <- chan OrderAccess) {
	orders := make(map[string][]Order) //The true map of unaccpeted orders
	for {
		select {
		case action := <- ordersActionChan:
			switch action.Cmd {
			case "read":
				deepCopy := make(map[string][]Order, len(orders))
				for key, value := range orders {
					deepCopy[key] = append([]Order{}, value...)
				}
				action.ReadChan <- deepCopy

			case "read all":
				action.ReadAllChan <- orders

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
					} else {
						delete(orders, action.Id)
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

func GetUnacceptedOrder(ordersActionChan chan<- OrderAccess, id string) []Order {
	readChan := make(chan map[string][]Order)
	defer close(readChan)

	ordersActionChan <- OrderAccess{
		Cmd:      "read",
		Id:		  id,
		ReadChan: readChan,
	}
	result := <-readChan
	return result[id]
}

func RemoveUnacceptedOrder(ordersActionChan chan<- OrderAccess, order Order) {
	ordersActionChan <- OrderAccess{
		Cmd:    "delete",
		Id:     order.Id,
		Orders: []Order{order},
	}
}

func GetAllUnacceptedOrders(orderActionChan chan<- OrderAccess) map[string][]Order {
    readAllChan := make(chan map[string][]Order)
	defer close(readAllChan)
    orderActionChan <- OrderAccess{Cmd: "read all", ReadAllChan: readAllChan}
    return <-readAllChan
}

func HallLightsManager(lightsActionChan <-chan LightsAccess) {
	hallLights := HallLights{}
	// for i := range hallLights {
    // 	hallLights[i] = make([]bool, NUM_BUTTONS-1) // Ensure all rows exist
	// }
	for {
		select {
		case action := <- lightsActionChan:
			switch action.Cmd {
			case "read":
				action.ReadChan <- hallLights
			case "write":
				hallLights = action.NewHallLights
			}
		}
	}
}

func ReadHallLights(lightsActionChan chan LightsAccess) HallLights {
	readChan := make(chan HallLights)
	defer close(readChan)

	lightsActionChan <- LightsAccess{
		Cmd: 	  "read",
		ReadChan: readChan,
	}
	return <- readChan
}

func WriteHallLights(lightsActionChan chan LightsAccess, newHallLights HallLights){
	lightsActionChan <- LightsAccess{
		Cmd: "write",
		NewHallLights: newHallLights,
	}
}