package sync

import ."source/config"

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
				action.ReadCh <- deepCopy
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