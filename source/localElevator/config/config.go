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

const SleepTime time.Duration = time.Millisecond * 20

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
	//Done bool
}

type FsmChansType struct {
	ElevatorChan chan Elevator
	AtFloorChan  chan int
	NewOrderChan chan Order
}