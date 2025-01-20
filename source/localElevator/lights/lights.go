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
/*
func lights(elev elevator_unit){
	elevio.SetFloorIndicator(elev.floor)
	for req := range elev.requests{
		elevio.SetButtonLamp(elevio.BT_Cab, elev.floor,)
	}
}
	*/