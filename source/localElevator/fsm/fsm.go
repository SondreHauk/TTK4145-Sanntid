package fsm

// This module should contain the finite state machine for the local elevator

import (
	"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/primary/requests"
	. "source/localElevator/config"
)

type FsmChansType struct{
	Elevator chan Elevator
	AtFloor chan int
}

func AtFloor(elev Elevator){
	switch currentState{
		case elevator.MOVING:
			if requests.ShouldStop(elev){
				elevio.SetMotorDirection(elevio.MD_Stop)
				doors.Open(elev)
				elev.State=elevator.DOOR_OPEN
				requests.ClearCurrentFloor(elev)
			}
			break
		default:
			break
		}
	}
}

func OnButtonPress(ButtonPress ButtonEvent){
	switch ButtonPress:
	case 
}

func Run(chans FsmChans){

}