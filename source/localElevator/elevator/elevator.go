package elevator

import (
	"localElevator/elevio"
	requests "localElevator/requestst"
)

// This module should contain the elevator struct and some actions that the elevator can perform
// What should the elevator struct contain (floor, direction, state, id, etc.)?
// What actions should the elevator be able to perform (open door, close door, move, etc.)?

type elevator_unit struct {
	floor     int
	direction int
	state     int
	requests  [4][3]bool
}

func init_elevator(){

}

func lights(elev elevator_unit){
	elevio.SetFloorIndicator(elev.floor)
	for req := range elev.requests{
		elevio.SetButtonLamp(elevio.BT_Cab,elev.floor,)
}

func door_open(){

}

func door_close(){
	
}
