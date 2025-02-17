package requests

import (
	"fmt"
	. "source/localElevator/config"
	"source/localElevator/elevio"
	"time"
)

func OrdersAbove(elev Elevator) bool {
	for fl := elev.Floor + 1; fl < NUM_FLOORS; fl++ {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Requests[fl][btn] {
				return true
			}
		}
	}
	return false
}

func OrdersBelow(elev Elevator) bool {
	for fl := elev.Floor - 1; fl >= 0; fl-- {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Requests[fl][btn] {
				return true
			}
		}
	}
	return false
}

func ClearFloor(elev *Elevator, floor int) {
	// Clear only the hall button in the right direction
	switch elev.Direction {
		case UP: // Clear hall up
			elev.Requests[floor][elevio.BT_HallUp] = false
			elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			if !OrdersAbove(*elev) {
				elev.Requests[floor][elevio.BT_HallDown] = false
				elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			}
		case DOWN: // Clear hall down
			elev.Requests[floor][elevio.BT_HallDown] = false
			elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			if !OrdersBelow(*elev) {
				elev.Requests[floor][elevio.BT_HallUp] = false
				elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			}
	}
	elev.Requests[floor][elevio.BT_Cab] = false
	elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
}	

func ClearAll(elev *Elevator) {
	for fl := 0; fl < NUM_FLOORS; fl++ {
		ClearFloor(elev, fl)
	}
}

func Update(Receiver chan elevio.ButtonEvent, Transmitter chan Order) {
	for{
		select {
			case btn := <-Receiver:
				Transmitter<-Order{Floor: btn.Floor, Button: int(btn.Button)}
				elevio.SetButtonLamp(elevio.ButtonType(btn.Button), btn.Floor, true) //THIS SIGNIFIES ORDER IS ACCEPTED. CHANGE
		}
		time.Sleep(T_SLEEP)
	}
}

func PrintRequests(elev Elevator){
	fmt.Printf("Floor 4: %t %t %t\n",elev.Requests[3][0],elev.Requests[3][1],elev.Requests[3][2])
	fmt.Printf("Floor 3: %t %t %t\n",elev.Requests[2][0],elev.Requests[2][1],elev.Requests[2][2])
	fmt.Printf("Floor 2: %t %t %t\n",elev.Requests[1][0],elev.Requests[1][1],elev.Requests[1][2])
	fmt.Printf("Floor 1: %t %t %t\n\n",elev.Requests[0][0],elev.Requests[0][1],elev.Requests[0][2])
}

func PrintState(elev Elevator){
	switch elev.State{
		case IDLE: fmt.Printf("State: IDLE\n")
		case MOVING: fmt.Printf("State: MOVING\n")
		case DOOR_OPEN: fmt.Printf("State: DOOR_OPEN\n")
	}
}