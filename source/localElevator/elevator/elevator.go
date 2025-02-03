package elevator

import (
	"localElevator/elevio"
	//"localElevator/requests"
	"time"
)

// This module should contain the elevator struct and some actions that the elevator can perform
// What should the elevator struct contain (floor, direction, state, id, etc.)?
// What actions should the elevator be able to perform (open door, close door, move, etc.)?

type Elevator struct {
	Floor     int
	Direction Direction
	Behavior  ElevatorBehaviour
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
}

//Drives down to the nearest floor and updates floor indicator
func ElevatorInit(elev Elevator){
	for elevio.GetFloor() == -1{
		time.Sleep(time.Millisecond*20)
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elev.Floor = elevio.GetFloor()
	elevio.SetFloorIndicator(elev.Floor)
}
