package elevator

import (
	"localElevator/elevio"
	//"localElevator/requests"
	"time"
)

// This module should contain the elevator struct and some actions that the elevator can perform
// What should the elevator struct contain (floor, direction, state, id, etc.)?
// What actions should the elevator be able to perform (open door, close door, move, etc.)?

const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
	OBSTRUCTED = 3
)

const(
	NUM_FLOORS = 4
	NUM_BUTTONS = 3
)

type Elevator struct {
	Floor     int
	Direction int
	State     int
	Requests  [4][3]bool
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