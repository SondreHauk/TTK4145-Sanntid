package elevator

import (
	"localElevator/elevio"
	//"localElevator/requests"
	"time"
)

// This module should contain the elevator struct and some actions that the elevator can perform
// What should the elevator struct contain (floor, direction, state, id, etc.)?
// What actions should the elevator be able to perform (open door, close door, move, etc.)?

type State int

const (
	IDLE State = iota
	MOVING
	DOOR_OPEN
	OBSTRUCTED
)

const(
	NUM_FLOORS = 4    //n
	NUM_BUTTONS = 3	  //m
)

type Elevator struct {
	Floor     int
	Direction Direction
	Behavior  ElevatorBehaviour
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
}

//Drives down to the nearest floor and updates floor indicator
func (elev *Elevator)ElevatorInit(){
	for elevio.GetFloor() == -1{
		time.Sleep(time.Millisecond*20)
		elevio.SetMotorDirection(elevio.MD_Down)
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	elev.Floor = elevio.GetFloor()
	elevio.SetFloorIndicator(elev.Floor)
}


//Moves to floor fl and updates floor indicators along the way.
func (elev *Elevator)MoveFloor(fl int){
	elev.Floor = elevio.GetFloor()
	if elev.Floor == -1{
		elev.ElevatorInit()
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

func (elev *Elevator)OpenDoor(){
	elevio.SetDoorOpenLamp(true)
}

func (elev *Elevator)CloseDoor(){
	elevio.SetDoorOpenLamp(false)
}