package requests

import (
    . "source/localElevator/elevator"
    . "source/localElevator/elevio"
)

type DirnBehaviourPair struct {
    Dirn     MotorDirection
    Behavior ElevatorBehaviour
}

func RequestsAbove(e Elevator) bool {
    for f := e.Floor + 1; f < NUM_FLOORS; f++ {
        for btn := 0; btn < NUM_BUTTONS; btn++{
            if(e.Requests[f][btn]){
                return true;
            }
        }
    }
    return false;
}

func RequestsBelow(e Elevator) bool {
    for f := 0; f < e.Floor; f++{
        for btn := 0; btn < NUM_BUTTONS; btn++{
            if(e.Requests[f][btn]){
                return true;
            }
        }
    }
    return false;
}

func RequestsHere(e Elevator) bool{
    for btn := 0; btn < NUM_BUTTONS; btn++{
        if(e.Requests[e.Floor][btn]){
            return true;
        }
    }
    return false;
}

func RequestsChooseDirection(e Elevator) DirnBehaviourPair {
    switch(e.Direction){
	case MD_Up:
		if RequestsAbove(e) {
			return DirnBehaviourPair{MD_Up, EB_MOVING}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{MD_Down, EB_DOOR_OPEN}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{MD_Down, EB_MOVING}
		}
	case MD_Down:
		if RequestsBelow(e) {
			return DirnBehaviourPair{MD_Down, EB_MOVING}
		} else if RequestsHere(e) {
			return DirnBehaviourPair{MD_Up, EB_DOOR_OPEN}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{MD_Up, EB_MOVING}
		}
	case MD_Stop:
		if RequestsHere(e) {
			return DirnBehaviourPair{MD_Stop, EB_DOOR_OPEN}
		} else if RequestsAbove(e) {
			return DirnBehaviourPair{MD_Up, EB_MOVING}
		} else if RequestsBelow(e) {
			return DirnBehaviourPair{MD_Down, EB_MOVING}
		}
    default:
        return DirnBehaviourPair{MD_Stop, EB_IDLE}
    }
}

func RequestsShouldStop(e Elevator) bool {
    switch(e.Direction) {
	case MD_Down:
		return e.Requests[e.Floor][BT_HallDown] || e.Requests[e.Floor][BT_Cab] || !RequestsBelow(e)
	case MD_Up:
		return e.Requests[e.Floor][BT_HallUp] || e.Requests[e.Floor][BT_Cab] || !RequestsAbove(e)
	case MD_Stop:
		return true
	}
	return true
}

func RequestsShouldClearImmediately(e Elevator, btnFloor int, btnType Button) bool {
    // Assumes all people enter the elevator even though the elevator is moving in the opposite direction
	return e.Floor == btnFloor
}

func Requests_ClearAtCurrentFloor(e Elevator) Elevator {
        for btn := 0; btn < NUM_BUTTONS; btn++ {
            e.Requests[e.Floor][btn] = 0;
        }
        return e;
    }