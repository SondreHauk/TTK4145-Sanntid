package localinit

import (
	. "source/config"
	"source/localElevator/elevio"
)

func ElevatorInit(id string) Elevator {
	currentFloor := elevio.GetFloor()
	if currentFloor == -1 {
		ch := make(chan int)
		go elevio.PollFloorSensor(ch)
		elevio.SetMotorDirection(elevio.MD_Down)
		currentFloor = <-ch //Blocks
		elevio.SetFloorIndicator(currentFloor)
		elevio.SetMotorDirection(elevio.MD_Stop)
	}
	elev := Elevator{
		Id:            id,
		Floor:         currentFloor,
		Direction:     int(elevio.MD_Stop),
		PrevDirection: DOWN,
		State:         IDLE,
		Orders:        OrderMatrixConstructor(),
		Requests:      OrderMatrixConstructor(),
		Obstructed:    false,
	}
	return elev
}

func LightsInit() {
	for fl := 0; fl < NUM_FLOORS; fl++ {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			elevio.SetButtonLamp(elevio.ButtonType(btn), fl, false)
		}
	}
}
