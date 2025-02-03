package lights

import (
	"localElevator/elevator"
	"localElevator/elevio"
)

func LightsInit(elev elevator.Elevator){
	for floor:=0;floor<elevator.NUM_FLOORS;floor++{
		for btn:=0;btn<elevator.NUM_BUTTONS;btn++{
			elevio.SetButtonLamp(elevio.ButtonType(btn),floor,false)
		}
	}
}

func LightsHallRequests(elev elevator.Elevator){
	for fl := range elev.Requests{
		for btn := range fl{
			elevio.SetButtonLamp(elevio.ButtonType(btn),fl,elev.Requests[fl][btn])
		}
	}
}