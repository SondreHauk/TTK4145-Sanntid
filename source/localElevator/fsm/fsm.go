package fsm

// This module should contain the finite state machine for the local elevator

import (
	"elevator"
	"elevio"
	"requests"
	"fmt"
	"time"
)



func FSM(){
	for{
		switch {
			case elevator.Behaviour == elevator.EB_IDLE:
			case elevator.Behaviour == elevator.EB_MOVING:
			case elevator.Behaviour == elevator.EB_DOOR_OPEN:
			case elevator.Behaviour == elevator.EB_OBSTRUCTED:			
		}
	}
}