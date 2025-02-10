package fsm

// This module should contain the finite state machine for the local elevator

import (
	. "source/localElevator/config"
	//"source/localElevator/elevator"
	"fmt"
	"source/localElevator/elevio"
	"source/localElevator/lights"
	"source/localElevator/requests"
	"time"
)


func OrdersAbove(elev Elevator) bool {
	for fl := elev.Floor + 1; fl < NUM_FLOORS; fl++ {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Requests[fl][btn] {
				return true
			}
		}
	}
	return false
}

func OrdersBelow(elev Elevator) bool {
	for fl := elev.Floor - 1; fl >= 0; fl-- {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Requests[fl][btn] {
				return true
			}
		}
	}
	return false
}

func ShouldStop(elev Elevator) bool {
	switch elev.Direction {
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

func ChooseDirection(elev Elevator) int {
	if OrdersAbove(elev) {
		return UP
	} else if OrdersBelow(elev) {
		return DOWN
	} else {
		return STOP
	}
}

func Run(elev Elevator, chans FsmChansType) {
	chans.ElevatorChan <- elev //Update elevator state
	DoorTimer := time.NewTimer(3 * time.Second)
	DoorTimer.Stop()

	for {
		select {
		case NewOrder := <-chans.NewOrderChan:
			if NewOrder.Done {
				requests.ClearFloor(elev, NewOrder.Floor)
			} else {
				elev.requests[NewOrder.Floor][NewOrder.Button] = true
			}
			switch elev.State {
			case IDLE:
				elev.Direction = ChooseDirection(elev)
				elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
				if elev.Direction == STOP {
					lights.OpenDoor(DoorTimer)
					elev.State = DOOR_OPEN
				} else {
					elev.State = MOVING
				}
			case MOVING: //NOOP
			case DOOR_OPEN:
				if elev.Floor == NewOrder.Floor {
					requests.ClearFloor(elev, elev.Floor)
					NewOrder.Done = true
				}
				/* 			case EMERGENCY_AT_FLOOR:
				   			case EMERGENCY_IN_SHAFT: Bør kanskje legges til funksjonalitet her?*/
			}

			fmt.Println("NewOrder UPDATE")
			chans.ElevatorChan <- elev //Update elevator state

		case elev.Floor = <-chans.AtFloorChan:
			if ShouldStop(elev) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				elev.State = DOOR_OPEN
				lights.OpenDoor(DoorTimer)
				requests.ClearFloor(elev, elev.Floor)
				fmt.Println("ShouldStop")
			}
			fmt.Println("AtFloor UPDATE")
			chans.ElevatorChan <- elev //Update elevator state

		case <-DoorTimer.C:
			elevio.SetDoorOpenLamp(false)
			dir := ChooseDirection(elev)
			if dir == STOP {
				elev.State = IDLE
			} else {
				elev.State = MOVING
				elevio.SetMotorDirection(elevio.MotorDirection(dir))
			}
			fmt.Println("DoorTimer UPDATE")
			chans.ElevatorChan <- elev //Update elevator state
		}
		time.Sleep(20*time.Millisecond)
	}
}
