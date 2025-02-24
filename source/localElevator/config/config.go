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
	T_TIMEOUT = time.Second*2
)

const(
	UP = 1
	DOWN = -1
	STOP = 0
)

const(
	PORT_PEERS = 20020
	PORT_BCAST_ELEV = 20030
	PORT_PRIMARY = 20040
)

type Elevator struct {
	Floor     int
	Direction int
	State  ElevatorState
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
	ID string
}

type Order struct {
	Floor int
	Button int
}

type Message struct {
	// ID spesific paramters
	ID int
	// Parameters for all
	Heartbeat string
}