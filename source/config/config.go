package config

import (
	"time"
)

type ElevatorState int

const (
    IDLE ElevatorState = iota
	MOVING
	DOOR_OPEN
	// EMERGENCY_AT_FLOOR
	// EMERGENCY_IN_SHAFT
	// OBSTRUCTED
)

const(
	NUM_FLOORS = 4
	NUM_BUTTONS = 3
	NUM_ELEVATORS = 1 // FOR NOW
)

const (
	T_HEARTBEAT = time.Millisecond*50 //Must be much faster than .5 s
	T_SLEEP = time.Millisecond*20
	T_DOOR_OPEN = time.Second*3
	T_OBSTRUCTED_PRIMARY = time.Second*3
	T_OBSTRUCTED_LOCAL = time.Second*4
	T_TRAVEL = time.Second*2 	//Approximate time to travel from floor i to floor i+-1
	T_PRIMARY_TIMEOUT = time.Millisecond*500
	T_BLINK = time.Millisecond*100
)

const(
	UP = 1
	DOWN = -1
	STOP = 0
)

//Is possible to use only one port for bcast.
const(
	PORT_PEERS = 20020
	PORT_ELEVSTATE = 20030
	PORT_WORLDVIEW = 20040
	PORT_REQUEST = 20050
	PORT_ORDER = 20060
	PORT_HALLLIGHTS = 20070
)

const(
	Obstructed = iota
	ConnectionLost
)

type Elevator struct {
	Id string
	Floor     int
	Direction int
	PrevDirection int
	State  ElevatorState
	Orders  [NUM_FLOORS][NUM_BUTTONS]bool
	Obstructed bool
}

type HallLights [][]bool 

type Order struct {
	Id string
	Floor int
	Button int
}