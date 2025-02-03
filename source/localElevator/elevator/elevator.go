package elevator

import (
	. "localElevator/elevio"
	//"localElevator/requests"
	"time"
)

// This module should contain the elevator struct and some actions that the elevator can perform
// What should the elevator struct contain (floor, direction, state, id, etc.)?
// What actions should the elevator be able to perform (open door, close door, move, etc.)?

type Elevator struct {
	Floor     int
	Direction MotorDirection
	State  ElevatorState
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
	//IsConnected bool  ??
}

//Drives down to the nearest floor and updates floor indicator
func (elev *Elevator)ElevatorInit(){
	for GetFloor() == -1{
		time.Sleep(time.Millisecond*20)
		SetMotorDirection(MD_Down)
	}
	SetMotorDirection(MD_Stop)
	elev.Floor = GetFloor()
	SetFloorIndicator(elev.Floor)
}


//Moves to floor fl and updates floor indicators along the way.
func (elev *Elevator)MoveFloor(fl int){
	elev.Floor = GetFloor()
	if elev.Floor == -1{
		elev.ElevatorInit()
	}
	
	if elev.Floor < fl{
		SetMotorDirection(MD_Up)
	}else if elev.Floor > fl{
		SetMotorDirection(MD_Down)
	}
	
	for elev.Floor != fl{
		if GetFloor() != -1{
			SetFloorIndicator(elev.Floor)
		}
		time.Sleep(time.Millisecond*20)
		elev.Floor = GetFloor()
	}
	SetMotorDirection(MD_Stop)
	SetFloorIndicator(elev.Floor)
}

func (elev *Elevator)OpenDoor(){
	SetDoorOpenLamp(true)
}

func (elev *Elevator)CloseDoor(){
	SetDoorOpenLamp(false)
}
