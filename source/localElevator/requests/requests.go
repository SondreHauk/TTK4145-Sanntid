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

func ClearOrder(elev Elevator, floor int) {
	switch elev.Direction {
	case UP:
		if OrdersAbove(elev) {
			elev.Orders[floor][elevio.BT_HallUp] = false
		} else if OrdersBelow(elev) {
			elev.Orders[floor][elevio.BT_HallDown] = false
		} else {
			elev.Orders[floor][elevio.BT_HallUp] = false
			elev.Orders[floor][elevio.BT_HallDown] = false
		}
	case DOWN:
		if OrdersBelow(elev) {
			elev.Orders[floor][elevio.BT_HallDown] = false
		} else if OrdersAbove(elev) {
			elev.Orders[floor][elevio.BT_HallUp] = false
		} else {
			elev.Orders[floor][elevio.BT_HallDown] = false
			elev.Orders[floor][elevio.BT_HallUp] = false
		}
	}
	elev.Orders[floor][elevio.BT_Cab] = false
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}

func SendRequest(
	reqEventChan <-chan elevio.ButtonEvent,
	requestsTXChan chan<- Requests,
	accReqChan <-chan OrderMatrix,
	id string
) {
	requests := OrderMatrixConstructor()
	heartBeat := time.NewTicker(T_HEARTBEAT)
	defer heartBeat.Stop()
	
  for {
		select {
		case accReq := <-accReqChan:
			for floor, orders := range accReq {
				for btn := range orders {
					if accReq[floor][btn] {
						requests[floor][btn] = false
					}
				}
			}

		case req := <-reqEventChan:
			requests[req.Floor][req.Button] = true
			requestsTXChan <- Requests{Id: id, Requests: requests}
			
			// if req.Button == elevio.BT_Cab {
			//	orderChan <- OrderConstructor(id, req.Floor, int(req.Button)) // Lets the elev operate when disconnected 
			//}

		case <-heartBeat.C:
			if checkForActiveRequests(requests) {
				requestsTXChan <- Requests{Id: id, Requests: requests}
			}
		}
		time.Sleep(T_SLEEP)
	}
}

func checkForActiveRequests(requests OrderMatrix) bool {
	for _, req := range requests {
		for _, active := range req {
			if active {
				return true
			}
		}
	}
	return false
}
