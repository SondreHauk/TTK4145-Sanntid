package lights

import (
	"source/localElevator/elevator"
	"source/localElevator/elevio"
)

func LightsInit(elev elevator.Elevator){
	for floor:=0;floor<elevator.NUM_FLOORS;floor++{
		for btn:=0;btn<elevator.NUM_BUTTONS;btn++{
			elevio.SetButtonLamp(elevio.ButtonType(btn),floor,false)
		}
	}
}

//Updates cab and floor lights wrt current request matrix.
//Will be called often.
func Update(elev elevator.Elevator){
	for fl := range elev.Requests{
		for btn := range fl{
			elevio.SetButtonLamp(elevio.ButtonType(btn),fl,elev.Requests[fl][btn])
		}
	}
}