package requests

import (
	. "source/localElevator/config"
	"source/localElevator/elevio"
	"time"
)

//Clears Lights and Request Matrix
func ClearFloor(elev Elevator, floor int) {
	for btn := 0; btn < NUM_BUTTONS; btn++ {
		elev.Requests[floor][btn] = false
		elevio.SetButtonLamp(elevio.ButtonType(btn), floor, false)
	}
}

func ClearAll(elev Elevator) {
	for fl := 0; fl < NUM_FLOORS; fl++ {
		ClearFloor(elev, fl)
	}
}

func ShouldStop(elev Elevator) {

}

func Update(BtnCh chan elevio.ButtonEvent, FsmCh FsmChansType) {
	for{
		select {
			case btn := <-BtnCh:
				FsmCh.NewOrderChan<-Order{btn.Floor, int(btn.Button), false}
				elevio.SetButtonLamp(elevio.ButtonType(btn.Button), btn.Floor, true) //THIS SIGNIFIES ORDER IS ACCEPTED. CHANGE
		}
		time.Sleep(20*time.Millisecond)
	}
}