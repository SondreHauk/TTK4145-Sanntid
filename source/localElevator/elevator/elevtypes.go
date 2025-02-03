package elevator

type ElevatorBehaviour int
type Direction int
type btn int

const (
    EB_Idle ElevatorBehaviour = iota
    EB_Moving
    EB_DoorOpen
)

const(
	D_Stop Direction = iota
	D_Down
	D_Up
)

const (
	IDLE = 0
	MOVING = 1
	DOOR_OPEN = 2
	OBSTRUCTED = 3
)

const(
	NUM_FLOORS = 4
	NUM_BUTTONS = 3
)