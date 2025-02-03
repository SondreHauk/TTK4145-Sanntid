package elevator

type ElevatorBehaviour int
type Direction int
type Button int

const (
    	EB_IDLE ElevatorBehaviour = iota
    	EB_MOVING
    	EB_DOOR_OPEN
	EB_OBSTRUCTED
)

const(
	NUM_FLOORS = 4
	NUM_BUTTONS = 3
)
