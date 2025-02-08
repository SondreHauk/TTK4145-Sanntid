package elevator

import (
	"source/localElevator/elevio"
	//"localElevator/requests"
	"time"
)

// This module should contain the elevator struct and some actions that the elevator can perform
// What should the elevator struct contain (floor, direction, state, id, etc.)?
// What actions should the elevator be able to perform (open door, close door, move, etc.)?

/* type Elevator struct {
	Floor     int
	Direction elevio.MotorDirection
	State  ElevatorState
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
} */

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


//Moves to floor fl without checking queue. Mostly for testing
func MoveFloor(elev Elevator, fl int){
	elev.Floor = elevio.GetFloor()
	if elev.Floor == -1{
		ElevatorInit(elev)
	}
	
	if elev.Floor < fl{
		elevio.SetMotorDirection(elevio.MD_Up)
	}else if elev.Floor > fl{
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	
	for elev.Floor != fl{
		if elevio.GetFloor() != -1{
			elevio.SetFloorIndicator(elev.Floor)
		}
		time.Sleep(time.Millisecond*20)
		elev.Floor = elevio.GetFloor()
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetFloorIndicator(elev.Floor)
}