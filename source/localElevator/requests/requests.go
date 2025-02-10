package requests

import (
	"fmt"
	. "source/localElevator/config"
	//"source/localElevator/elevator"
	"source/localElevator/elevio"
	"time"
)

//Clears Lights and Request Matrix
func ClearFloor(elev *Elevator, floor int) {
	/* elev.Requests[floor]= [NUM_BUTTONS]bool{} */
	for btn := 0; btn < NUM_BUTTONS; btn++ {
		elev.Requests[floor][btn] = false
		elevio.SetButtonLamp(elevio.ButtonType(btn), floor, false)
	}
	//fmt.Printf("elev.Requests[%d] = %t\n",floor, elev.Requests[floor])
	fmt.Printf("Floor %d cleared\n", floor+1)
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
				Transmitter<-Order{btn.Floor, int(btn.Button)}//, false}
				elevio.SetButtonLamp(elevio.ButtonType(btn.Button), btn.Floor, true) //THIS SIGNIFIES ORDER IS ACCEPTED. CHANGE
		}
		time.Sleep(20*time.Millisecond)
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