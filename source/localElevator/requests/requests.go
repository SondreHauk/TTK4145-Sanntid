package requests

import (
	"fmt"
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

//PRIMARY MUST CLEAR HALL LIGHTS - NOT ELEVATOR!
func ClearOrder(elev *Elevator, floor int) {
	// Clear only the hall button in the right direction
	switch elev.Direction {
		case UP: // Clear hall up
			elev.Orders[floor][elevio.BT_HallUp] = false
			//elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
			if !OrdersAbove(*elev) {
				elev.Orders[floor][elevio.BT_HallDown] = false
				//elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			}
		case DOWN: // Clear hall down
			elev.Orders[floor][elevio.BT_HallDown] = false
			//elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
			if !OrdersBelow(*elev) {
				elev.Orders[floor][elevio.BT_HallUp] = false
				//elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
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

func MakeRequest(btnEvent <-chan elevio.ButtonEvent, 
	requestToPrimary chan <-Order, 
	orderChan chan <- Order,
	id string) {
	for{
		select {
			case btn := <-btnEvent:
				request := Order{Id: id, Floor: btn.Floor, Button: int(btn.Button)}
				if btn.Button == elevio.BT_Cab{
					orderChan <- request // Assign directly to elev
					elevio.SetButtonLamp(elevio.ButtonType(btn.Button), btn.Floor, true)
					//Remember: Lights on = Order MUST be taken
				} else {
					requestToPrimary<- request
				}
		}
		time.Sleep(T_SLEEP) //Necessary?
	}
}

//Make modular with for loop up to NUM_ELEV
func PrintRequests(elev Elevator){
	fmt.Printf("Floor 4: %t %t %t\n",elev.Orders[3][0],elev.Orders[3][1],elev.Orders[3][2])
	fmt.Printf("Floor 3: %t %t %t\n",elev.Orders[2][0],elev.Orders[2][1],elev.Orders[2][2])
	fmt.Printf("Floor 2: %t %t %t\n",elev.Orders[1][0],elev.Orders[1][1],elev.Orders[1][2])
	fmt.Printf("Floor 1: %t %t %t\n\n",elev.Orders[0][0],elev.Orders[0][1],elev.Orders[0][2])
}

func PrintState(elev Elevator){
	switch elev.State{
		case IDLE: fmt.Printf("State: IDLE\n")
		case MOVING: fmt.Printf("State: MOVING\n")
		case DOOR_OPEN: fmt.Printf("State: DOOR_OPEN\n")
	}
}