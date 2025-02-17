// This module will assign a hall call to an elevator
// based on the current state of the elevator and the hall call
// This will be done in a cost function.

// Alternative 1: Assigning only the new request.

package assigner

import (
	"source/localElevator/elevator"
	."source/localElevator/config"
)

type Button int
type Direction int

const (
	D_Stop Direction  = 0
	N_FLOORS 	      = 3 
	N_BUTTONS         = 4 
	TRAVEL_TIME       = 5
	DOOR_OPEN_TIME    = 3
)

type Behaviour int

const (
	EB_Idle Behaviour = iota
	EB_Moving
	EB_DoorOpen
)

type ClearedRequestFunc func(Button, int)

func requestsClearAtCurrentFloor(e Elevator, onClearedRequest ClearedRequestFunc) Elevator {
	for btn := Button(0); btn < N_BUTTONS; btn++ {
		if e.requests[e.floor][btn] {
			e.requests[e.floor][btn] = false
			if onClearedRequest != nil {
				onClearedRequest(btn, e.floor)
			}
		}
	}
	return e
}

func requestsChooseDirection(e Elevator) Direction {
	DirnBehaviourPair requestsChooseDirection(e Elevator){
		switch(e.dirn){
		case D_Up:
			return  requests_above(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
					requests_here(e)  ? (DirnBehaviourPair){D_Down, EB_DoorOpen} :
					requests_below(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
										(DirnBehaviourPair){D_Stop, EB_Idle}     ;
		case D_Down:
			return  requests_below(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
					requests_here(e)  ? (DirnBehaviourPair){D_Up,   EB_DoorOpen} :
					requests_above(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
										(DirnBehaviourPair){D_Stop, EB_Idle}     ;
		case D_Stop: // there should only be one request in the Stop case. Checking up or down first is arbitrary.
			return  requests_here(e)  ? (DirnBehaviourPair){D_Stop, EB_DoorOpen} :
					requests_above(e) ? (DirnBehaviourPair){D_Up,   EB_Moving}   :
					requests_below(e) ? (DirnBehaviourPair){D_Down, EB_Moving}   :
										(DirnBehaviourPair){D_Stop, EB_Idle}     ;
		default:
			return (DirnBehaviourPair){D_Stop, EB_Idle};
		}
	}
	
	return D_Stop
}

func requestsShouldStop(e Elevator) bool {
	// Placeholder logic for stopping decision
	int requests_shouldStop(Elevator e){
		switch(e.dirn){
		case D_Down:
			return
				e.requests[e.floor][B_HallDown] ||
				e.requests[e.floor][B_Cab]      ||
				!requests_below(e);
		case D_Up:
			return
				e.requests[e.floor][B_HallUp]   ||
				e.requests[e.floor][B_Cab]      ||
				!requests_above(e);
		case D_Stop:
		default:
			return 1;
		}
	}
	return false
}

func timeToIdle(e Elevator) int {
	duration := 0
	// Enters the switch case once to determine the initial state of the elevator
	switch e.behaviour {
	case EB_Idle:
		e.dirn = requestsChooseDirection(e)
		if e.dirn == D_Stop {
			return duration
		}
	case EB_Moving:
		duration += TRAVEL_TIME / 2
		e.floor += int(e.dirn)
	case EB_DoorOpen:
		duration -= DOOR_OPEN_TIME / 2
	}
	// Adds exectution time for the elevator to reach IDLE (clear all pending requests).
	for {
		if requestsShouldStop(e) {
			e = requestsClearAtCurrentFloor(e, nil)
			duration += DOOR_OPEN_TIME
			e.dirn = requestsChooseDirection(e)
			if e.dirn == D_Stop {
				return duration
			}
		}
		e.floor += int(e.dirn)
		duration += TRAVEL_TIME
	}
}
