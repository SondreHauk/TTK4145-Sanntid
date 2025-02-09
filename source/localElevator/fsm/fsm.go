package fsm

// This module should contain the finite state machine for the local elevator

import (
	. "source/localElevator/config"
	//"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/localElevator/requests"
	"time"
)

type FsmChansType struct{
	ElevatorChan chan Elevator
	AtFloorChan chan int
	NewOrderChan chan elevio.ButtonEvent
}

func OrdersAbove(elev Elevator) bool{
	for fl:=elev.Floor+1; fl<NUM_FLOORS; fl++{
		for btn:=0; btn<NUM_BUTTONS; btn++{
			if elev.Requests[fl][btn]{
				return true
			}
		}
	}
	return false
}

func OrdersBelow(elev Elevator) bool{
	for fl:=elev.Floor-1; fl>=0; fl--{
		for btn:=0; btn<NUM_BUTTONS; btn++{
			if elev.Requests[fl][btn]{
				return true
			}
		}
	}
	return false
}

func ShouldStop(elev Elevator) bool{
	switch elev.Direction{
		case UP:
			return elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_Cab] || !OrdersAbove(elev)
		case DOWN:
			return elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] || !OrdersBelow(elev)
		case STOP:
			return true
		default:
			return false
	}
}

func ChooseDirection(elev Elevator) int{
	if OrdersAbove(elev){
		return UP
	}else if OrdersBelow(elev){
		return DOWN
	}else{
		return STOP
	}
}

func Run(elev Elevator, chans FsmChansType){
	chans.ElevatorChan <- elev
	DoorTimer:=time.NewTimer(3*time.Second)
	DoorTimer.Stop()

	for{
		select {
		case NewOrder := <- chans.NewOrderChan:
			elev.Requests[NewOrder.Floor][NewOrder.Button] = true
			
			chans.ElevatorChan <- elev //Update elevator state
		case elev.Floor = <- chans.AtFloorChan:
			if ShouldStop(elev){
				elevio.SetMotorDirection(elevio.MD_Stop)
				
				//DOORS
				elev.State = DOOR_OPEN
				elevio.SetDoorOpenLamp(true)
				DoorTimer.Reset(3*time.Second)
				
				requests.ClearCurrentFloor(elev)
				
			}
			chans.ElevatorChan <- elev //Update elevator state

		case <- DoorTimer.C:
			elevio.SetDoorOpenLamp(false)
			dir:=ChooseDirection(elev)
			if dir == STOP{
				elev.State = IDLE
			}else{
				elev.State = MOVING
				elevio.SetMotorDirection(elevio.MotorDirection(dir))
			}
			chans.ElevatorChan <- elev //Update elevator state
		}
	}
}