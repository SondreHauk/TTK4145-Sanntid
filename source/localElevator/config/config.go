package config
import "time"

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

const(
	BTN_HALL_UP = iota
	BTN_HALL_DOWN
	BTN_CAB
)

type Elevator struct {
	Floor     int
	Direction int
	State  ElevatorState
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
}