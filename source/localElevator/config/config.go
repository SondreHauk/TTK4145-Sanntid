package config

import (
	"time"
)

type ElevatorState int

const (
    IDLE ElevatorState = iota
	MOVING
	DOOR_OPEN
	EMERGENCY_AT_FLOOR
	EMERGENCY_IN_SHAFT
)

const(
	NUM_FLOORS = 4
	NUM_BUTTONS = 3
)

const (
	T_HEARTBEAT = time.Millisecond*500
	T_SLEEP = time.Millisecond*20
	T_DOOR_OPEN = time.Second*3	
)

const(
	UP = 1
	DOWN = -1
	STOP = 0
)

type Elevator struct {
	Floor     int
	Direction int
	State  ElevatorState
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
}

type Order struct {
	Floor int
	Button int
}