// This module will assign a hall call to an elevator
// based on the current state of the elevator and the hall call
// This will be done in a cost function.

// Alternative 1: Assigning only the new request.

package assigner

import (
	. "source/localElevator/config"
	"source/localElevator/fsm"
	"source/localElevator/requests"
	"time"
)

//Creates a copy of the elevator and simulates executing remaining orders
func TimeToIdle(elev Elevator) time.Duration {
	duration := time.Duration(0)
	// Determines initial state
	switch elev.State {
	case IDLE:
		elev.Direction = fsm.ChooseDirection(elev)
		if elev.Direction == STOP {
			return duration
		}
	case MOVING:
		duration += T_TRAVEL / 2
		elev.Floor += int(elev.Direction)
	case DOOR_OPEN:
		duration -= T_DOOR_OPEN / 2
	}
	
	//Simulates remaining orders
	for {
		if fsm.ShouldStop(elev) {
			requests.ClearFloor(&elev, elev.Floor) //Changes do not propagate back to main
			duration += T_DOOR_OPEN
			elev.Direction = fsm.ChooseDirection(elev)
			if elev.Direction == STOP {
				return duration
			}
		}
		elev.Floor += int(elev.Direction)
		duration += T_TRAVEL
	}
}

func AssignElevator(elevs chan *Elevator){
	for{
		select {
		case e<-elevs:
			if 
			
		}
	}
}