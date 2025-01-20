package elevator

import (
	"localElevator/elevio"
	"localElevator/requests"
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

type elevator_unit struct {
	floor     int
	direction int
	state     int
	requests  [4][3]bool
}

func init_elevator(elev elevator_unit){
	for elevtio.GetFloor() == -1{
		time.Sleep(elevio._pollRate)
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elev.floor = elevio.GetFloor()
	
}