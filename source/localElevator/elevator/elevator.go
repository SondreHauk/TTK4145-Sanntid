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
	NUM_FLOORS = 4    //n
	NUM_BUTTONS = 3	  //m
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