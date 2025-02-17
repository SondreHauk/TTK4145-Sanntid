package fsm

// This module should contain the finite state machine for the local elevator

import (
	. "source/localElevator/config"
	//"source/localElevator/elevator"
	//"source/localElevator/elevator"
	"fmt"
	"math/rand"
	"source/localElevator/elevio"
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
		if elev.Floor==NUM_FLOORS-1{
			return true
		}else{
			return elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_Cab] || !OrdersAbove(elev)
		}
	case DOWN:
		if elev.Floor==0{
			return true
		}else{
			return elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] || !OrdersBelow(elev)
		}
	case STOP:
		return true
	default:
		fmt.Println("DEFAULT ERROR STOP")
		return false
	}
}

func ChooseDirection(elev Elevator) int {
	// In case of orders above and below; choose direction at random
	// Not very smaart, but it works
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	if r % 2 == 0{
		if OrdersAbove(elev) {
			return UP
		} else if OrdersBelow(elev) {
			return DOWN
		}
	} else {
		if OrdersBelow(elev) {
			return DOWN
		} else if OrdersAbove(elev) {
			return UP
		}
	}
	return STOP
}

func Run(elev *Elevator, /* ElevCh chan *Elevator, */ AtFloorCh chan int, NewOrderCh chan Order) {
	//ElevCh <- elev //Send updated elevator state to master
	DoorTimer := time.NewTimer(T_DOOR_OPEN)
	DoorTimer.Stop()
	for {
		select {
		case NewOrder := <-NewOrderCh:
			elev.Requests[NewOrder.Floor][NewOrder.Button] = true
			switch elev.State {
				case IDLE:
					elev.Direction = ChooseDirection(*elev)
					elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
					if elev.Direction == STOP {
						elevio.SetDoorOpenLamp(true)
						DoorTimer.Reset(T_DOOR_OPEN)
						elev.State = DOOR_OPEN
					} else {
						elev.State = MOVING
					}
				case MOVING: //NOOP
				case DOOR_OPEN:
					if elev.Floor == NewOrder.Floor {
						requests.ClearFloor(elev, elev.Floor)
						DoorTimer.Reset(T_DOOR_OPEN)
					}
			}
			//ElevCh <- elev //Send updated elevator state to master

		case elev.Floor = <-AtFloorCh:
			elevio.SetFloorIndicator(elev.Floor)
			if ShouldStop(*elev) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				requests.ClearFloor(elev, elev.Floor)
				elevio.SetDoorOpenLamp(true)
				DoorTimer.Reset(T_DOOR_OPEN)
				elev.State = DOOR_OPEN
			}
			//ElevCh <- elev //Send updated elevator state to master

		case <-DoorTimer.C:

			elevio.SetDoorOpenLamp(false)
			elev.Direction = ChooseDirection(*elev)
			if elev.Direction == STOP {
				elev.State = IDLE
			} else {
				elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
				elev.State = MOVING
			}
			//ElevCh <- elev //Send updated elevator state to master
		}
		time.Sleep(T_SLEEP)
	}
}