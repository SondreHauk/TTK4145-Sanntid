package sync

import (
	. "source/config"
)

// "fleet" is supposed to signify that the primary acts like a maritime general with subordinate ships(elevators)
func FleetAccessManager(fleetAccessChan <-chan FleetAccess) {
	fleet := make(map[string]Elevator)
	for {
		select {
		case newAction := <-fleetAccessChan:
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
				fleet = newAction.FullFleet
			}
		}
	}
}

func SingleElevFleetWrite(
	id string,
	elev Elevator,
	fleetAccessChan chan FleetAccess,
) {
	fleetAccessChan <- FleetAccess{Cmd: "write one", Id: id, Elev: elev}
}

func FullFleetWrite(
	fullFleet map[string]Elevator,
	fleetAccessChan chan FleetAccess,
) {
	fleetAccessChan <- FleetAccess{Cmd: "write all", FullFleet: fullFleet}
}

func FleetRead(fleetAccessChan chan FleetAccess) map[string]Elevator {
	readChan := make(chan map[string]Elevator, 1)
	defer close(readChan)
	fleetAccessChan <- FleetAccess{Cmd: "read", ReadChan: readChan}
	return <-readChan
}

func UnacceptedOrdersManager(ordersActionChan <-chan OrderAccess) {
	orders := make(map[string][]Order)
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
						if !(o.Floor == action.Orders[0].Floor &&
							o.Button == action.Orders[0].Button) {
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
	if order.Id != "" {
		ordersActionChan <- OrderAccess{
			Cmd:    "write",
			Id:     order.Id,
			Orders: []Order{order},
		}
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

func AcceptOrder(ordersActionChan chan<- OrderAccess, order Order) {
	ordersActionChan <- OrderAccess{
		Cmd:    "delete",
		Id:     order.Id,
		Orders: []Order{order},
	}
}

func HallLightsManager(lightsActionChan <-chan LightsAccess) {
	hallLights := HallMatrixConstructor()
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

func WriteHallLights(
	lightsActionChan chan LightsAccess,
	newHallLights HallMatrix,
) {
	lightsActionChan <- LightsAccess{
		Cmd:           "write",
		NewHallLights: newHallLights,
	}
}
