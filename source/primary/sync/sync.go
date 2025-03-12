package sync

import ."source/config"

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
	select{
	case output := <-readChan:
		return output
	}
}

func UnacceptedOrdersManager(ordersActionChan <- chan OrderAccess) {
	orders := make(map[string][]Order)
	for {
		select {
		case action := <- ordersActionChan:
			switch action.Cmd {
			case "read":
				deepCopy := make(map[string][]Order, len(orders))
				for key, value := range orders {
					deepCopy[key] = value
				}
				action.ReadChan <- deepCopy
			case "write one":
				orders[action.Id] = action.Orders
			case "write all":
				orders = action.UnacceptedOrders
			}
		}
	}
}

func HallLightsManager() {
	
}