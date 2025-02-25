package fsm

// This module should contain the finite state machine for the local elevator

import (
	. "source/localElevator/config"
	"source/localElevator/elevio"
	"source/localElevator/requests"
	"time"
)

func ShouldStop(elev Elevator) bool {
	switch elev.Direction {
	case UP:
		if elev.Floor==NUM_FLOORS-1{
			return true
		}else{
			return elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_Cab] || !requests.OrdersAbove(elev)
		}
	case DOWN:
		if elev.Floor==0{
			return true
		}else{
			return elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] || !requests.OrdersBelow(elev)
		}
	case STOP:
		return true
	}
	return false
}

func ChooseDirection(elev Elevator) int {
	// In case of orders above and below; choose last moving direction
	if elev.PrevDirection == UP{
		if requests.OrdersAbove(elev) {
			return UP
		} else if requests.OrdersBelow(elev) {
			return DOWN
		}
	} else {
		if requests.OrdersBelow(elev) {
			return DOWN
		} else if requests.OrdersAbove(elev) {
			return UP
		}
	}
	return STOP

}

//Simulates elevator execution and returns approx time until pickup at NewOrder.Floor
func TimeUntilPickup(elev Elevator, NewOrder Order) time.Duration{
	duration := time.Duration(0)
	elev.Requests[NewOrder.Floor][NewOrder.Button]=true
	// Determines initial state
	switch elev.State {
	case IDLE:
		elev.Direction = ChooseDirection(elev)
		if elev.Direction == STOP && elev.Floor == NewOrder.Floor{
			return duration
		}
	case MOVING:
		duration += T_TRAVEL / 2
		elev.Floor += int(elev.Direction)
	case DOOR_OPEN:
		duration -= T_DOOR_OPEN / 2
	}

	for {
		if ShouldStop(elev) {
			if elev.Floor == NewOrder.Floor{
				return duration
			}else{
				for btn:=0; btn<NUM_BUTTONS; btn++{
					elev.Requests[elev.Floor][btn]=false
				}
				duration += T_DOOR_OPEN
				elev.Direction = ChooseDirection(elev)
			}
		}
		elev.Floor += int(elev.Direction)
		duration += T_TRAVEL
	}
}

func Run(
	elev *Elevator, 
	ElevCh chan <-Elevator, 
	AtFloorCh <-chan int, 
	OrderFromPrimaryChan <-chan Order, 
	ObsCh <-chan bool,
	myId string) {

	ElevCh <- *elev
	HeartbeatTimer := time.NewTimer(T_HEARTBEAT)
	DoorTimer := time.NewTimer(T_DOOR_OPEN)
	DoorTimer.Stop()
	Obstructed := false
	
	for {
		select {
		case NewOrder := <-OrderFromPrimaryChan:
			if NewOrder.Id == myId{
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
						elev.Requests[elev.Floor][NewOrder.Button] = false
						elevio.SetButtonLamp(elevio.ButtonType(NewOrder.Button), elev.Floor, false)
						if !Obstructed{
							DoorTimer.Reset(T_DOOR_OPEN)
						}
					}
				}
				ElevCh <- *elev
			}

		case elev.Floor = <-AtFloorCh:
			elevio.SetFloorIndicator(elev.Floor)
			if ShouldStop(*elev) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				requests.ClearFloor(elev, elev.Floor)
				elev.Direction = STOP
				elevio.SetDoorOpenLamp(true)
				DoorTimer.Reset(T_DOOR_OPEN)
				elev.State = DOOR_OPEN
			}
			ElevCh <- *elev

		case <-DoorTimer.C:

			elevio.SetDoorOpenLamp(false)
			elev.Direction = ChooseDirection(*elev)
			if elev.Direction == STOP {
				elev.State = IDLE
			} else {
				elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
				elev.State = MOVING
			}
			ElevCh <- *elev
		
		case ObsEvent:= <-ObsCh:
			if elev.State==DOOR_OPEN{
				switch ObsEvent{
					case true:
						Obstructed = true
						DoorTimer.Stop()
					case false:
						Obstructed = false
						DoorTimer.Reset(T_DOOR_OPEN)
				}
			}
		case <-HeartbeatTimer.C:
			ElevCh <- *elev
			HeartbeatTimer.Reset(T_HEARTBEAT)
		}

		time.Sleep(T_SLEEP)
	}
}