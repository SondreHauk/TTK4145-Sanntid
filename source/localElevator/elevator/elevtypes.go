package elevator

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