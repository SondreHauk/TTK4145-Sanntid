package requests

import (
	. "source/config"
	"source/localElevator/elevio"
	"time"
)

func OrdersAbove(elev Elevator) bool {
	for fl := elev.Floor + 1; fl < NUM_FLOORS; fl++ {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Orders[fl][btn] {
				return true
			}
		}
	}
	return false
}

func OrdersBelow(elev Elevator) bool {
	for fl := elev.Floor - 1; fl >= 0; fl-- {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Orders[fl][btn] {
				return true
			}
		}
	}
	return false
}

func ClearOrder(elev *Elevator, floor int) {
	switch elev.Direction {
		case UP:
			elev.Orders[floor][elevio.BT_HallUp] = false
			if !OrdersAbove(*elev) {
				elev.Orders[floor][elevio.BT_HallDown] = false
			}
		case DOWN:
			elev.Orders[floor][elevio.BT_HallDown] = false
			if !OrdersBelow(*elev) {
				elev.Orders[floor][elevio.BT_HallUp] = false
			}
	}
	elev.Orders[floor][elevio.BT_Cab] = false
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}	

func ClearAll(elev *Elevator) {
	for fl := 0; fl < NUM_FLOORS; fl++ {
		ClearOrder(elev, fl)
	}
}

func MakeRequest(
	reqEventChan <-chan elevio.ButtonEvent, 
	requestToPrimaryChan chan <- HallMatrix, 
	orderChan chan <- Order,
	accReqChan <- chan HallMatrix,
	id string) {

	requests := HallMatrix{}
	heartBeat := time.NewTicker(T_HEARTBEAT)
	defer heartBeat.Stop()

	for{
		select {
		case accReq := <- accReqChan:
			for floor, orders := range accReq{
				for btn := range orders{
					if accReq[floor][btn] {
						requests[floor][btn] = false // RACE CONDITION
					}
				}
			}

		case req := <- reqEventChan:

			requests[req.Floor][req.Button] = true // RACE CONDITION
			request := Order{Id: id, Floor: req.Floor, Button: int(req.Button)}
			
			if req.Button == elevio.BT_Cab{
				orderChan <- request // Assign directly to elev
				elevio.SetButtonLamp(elevio.ButtonType(req.Button), req.Floor, true)
			} else {
				requestToPrimaryChan <- requests
			}

		case <- heartBeat.C:
			if checkForActiveRequests(requests) {
				requestToPrimaryChan <- requests // RACE CONDITION
			}
		}
		time.Sleep(T_SLEEP)
	}
}

func checkForActiveRequests(requests HallMatrix) bool {
	for _, req := range requests{
		for _, active := range req {
			if active {
				return true
			}
		}
	}
	return false
}