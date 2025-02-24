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
	NUM_ELEVATORS = 1 // FOR NOW
)

const (
	T_HEARTBEAT = time.Millisecond*500
	T_SLEEP = time.Millisecond*20
	T_DOOR_OPEN = time.Second*3
	T_TRAVEL = time.Second*2 	//Approximate time to travel from floor i to floor i+-1
)

const(
	UP = 1
	DOWN = -1
	STOP = 0
)

type Elevator struct {
	Id 	int
	Floor     int
	Direction int
	State  ElevatorState
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
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

func DeepCopyElev(elev Elevator){
	
}