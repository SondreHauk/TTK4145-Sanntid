package fsm

import (
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/requests"
	"time"
)

func shouldStop(elev Elevator) bool {
	switch elev.Direction {
	case UP:
		if elev.Floor == NUM_FLOORS-1 {
			return true
		} else {
			return elev.Orders[elev.Floor][elevio.BT_HallUp] ||
				elev.Orders[elev.Floor][elevio.BT_Cab] ||
				!requests.OrdersAbove(elev)
		}
	case DOWN:
		if elev.Floor == 0 {
			return true
		} else {
			return elev.Orders[elev.Floor][elevio.BT_HallDown] ||
				elev.Orders[elev.Floor][elevio.BT_Cab] ||
				!requests.OrdersBelow(elev)
		}
	case STOP:
		return true
	}
	return false
}

func chooseDirection(elev Elevator) int {
	// In case of orders above and below; choose last moving direction
	if elev.PrevDirection == UP {
		if requests.OrdersAbove(elev) {
			return UP
		} else if requests.OrdersBelow(elev) {
			return DOWN
		}
	} else {
		if requests.OrdersBelow(elev) {
			return DOWN
		} else if requests.OrdersAbove(elev) {
			return UP
		}
	}
	return STOP

}

// Is only used by primary/assigner but uses fsm helper functions
func TimeUntilPickup(elev Elevator, NewOrder Order) time.Duration {
	duration := time.Duration(0)
	elev.Orders[NewOrder.Floor][NewOrder.Button] = true
	// Determines initial state
	switch elev.State {
	case IDLE:
		elev.Direction = chooseDirection(elev)
		if elev.Direction == STOP && elev.Floor == NewOrder.Floor {
			return duration
		}
	case MOVING:
		duration += T_TRAVEL / 2
		elev.Floor += int(elev.Direction)
	case DOOR_OPEN:
		duration -= T_DOOR_OPEN / 2
	}

	for {
		if shouldStop(elev) {
			if elev.Floor == NewOrder.Floor {
				return duration
			} else {
				for btn := 0; btn < NUM_BUTTONS; btn++ {
					elev.Orders[elev.Floor][btn] = false
				}
				duration += T_DOOR_OPEN
				elev.Direction = chooseDirection(elev)
			}
		}
		elev.Floor += int(elev.Direction)
		duration += T_TRAVEL
	}
}

func checkForNewOrders(
	wv Worldview,
	myId string,
	orderChan chan<- Order,
	acceptedRequestsChan chan<- OrderMatrix,
	acceptedOrders OrderMatrix) {

	// send all assigned orders to request module
	accOrdersMatrix := OrderMatrixConstructor()
	for _, accOrders := range wv.UnacceptedOrdersSnapshot {
		for _, ord := range accOrders {
			accOrdersMatrix[ord.Floor][ord.Button] = true
		}
	}
	acceptedRequestsChan <- accOrdersMatrix

	// send ID assigned order to elevator
	orders, exists := wv.UnacceptedOrdersSnapshot[myId]
	if exists {
		for _, order := range orders {
			if !acceptedOrders[order.Floor][order.Button] {
				orderChan <- order
			}
		}
	}
}

func checkForNewLights(wv Worldview, lights HallMatrix, lightsChan chan HallMatrix) {
	for floor, buttons := range lights {
		for btn := range buttons {
			if lights[floor][btn] != wv.HallLightsSnapshot[floor][btn] {
				lightsChan <- wv.HallLightsSnapshot
				return
			}
		}
	}
}

func setHallLights(lights HallMatrix) {
	for floor, btns := range lights {
		for btn, status := range btns {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, status)
		}
	}
}

func resetTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(duration)
}

// Send multiple times to avoid hall light blinking, which can happen if there is *severe* packetloss.
// Message will never be truly lost, as primary will just reassign order.
func ackOrder(elev Elevator, elevChan chan<- Elevator) {
	for range 10 {
		elevChan <- elev
	}
	time.Sleep(T_SLEEP)
}
