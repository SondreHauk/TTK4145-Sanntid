package requests

import (
	. "source/config"
	"source/localElevator/elevio"
	"time"
	"fmt"
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
		if OrdersAbove(*elev){
			elev.Orders[floor][elevio.BT_HallUp] = false
		} else if OrdersBelow(*elev){
			elev.Orders[floor][elevio.BT_HallDown] = false
			fmt.Println("Switching direction")
		} else {
			elev.Orders[floor][elevio.BT_HallUp] = false
			elev.Orders[floor][elevio.BT_HallDown] = false
			fmt.Println("No more orders")
		}
	case DOWN:
		if OrdersBelow(*elev){
			elev.Orders[floor][elevio.BT_HallDown] = false
		} else if OrdersAbove(*elev) {
			elev.Orders[floor][elevio.BT_HallUp] = false
			fmt.Println("Switching direction")
		} else {
			elev.Orders[floor][elevio.BT_HallDown] = false
			elev.Orders[floor][elevio.BT_HallUp] = false
			fmt.Println("No more orders")
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

	hallRequests := HallMatrix{}
	heartBeat := time.NewTicker(T_HEARTBEAT)
	defer heartBeat.Stop()

	for{
		select {
		case accReq := <- accReqChan:
			for floor, orders := range accReq{
				for btn := range orders{
					if accReq[floor][btn] {
						hallRequests[floor][btn] = false
					}
				}
			}

		case req := <- reqEventChan:
			if req.Button == elevio.BT_Cab{
				orderChan <- OrderConstructor(id, req.Floor, int(req.Button)) // Assign directly to elev
				elevio.SetButtonLamp(elevio.ButtonType(req.Button), req.Floor, true)
			} else {
				hallRequests[req.Floor][req.Button] = true
			 	requestToPrimaryChan <- hallRequests
			}

		case <- heartBeat.C:
			if checkForActiveRequests(hallRequests) {
				requestToPrimaryChan <- hallRequests
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