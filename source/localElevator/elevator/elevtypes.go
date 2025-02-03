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
	D_STOP Direction = iota
	D_DOWN
	D_UP
)

const(
	B_HALL_DOWN Button = iota
	B_HALL_UP
	B_CAB
)

const(
	NUM_FLOORS = 4
	NUM_BUTTONS = 3
)
