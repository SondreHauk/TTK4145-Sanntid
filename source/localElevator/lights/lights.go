package lights

import (
	. "source/localElevator/config"
	//source/localElevator/elevator"
	"source/localElevator/elevio"
)

func LightsInit(elev Elevator){
	for fl:=0; fl<NUM_FLOORS; fl++{
		for btn:=0; btn<NUM_BUTTONS; btn++{
			elevio.SetButtonLamp(elevio.ButtonType(btn),fl,false)
		}
	}
}

//Updates cab and floor lights wrt current request matrix.
//Will be called often.
func Update(elev Elevator){
	for fl := range elev.Requests{
		for btn := range fl{
			elevio.SetButtonLamp(elevio.ButtonType(btn),fl,elev.Requests[fl][btn])
		}
	}
}

/* func StopAtFloor(elev Elevator){
	elevio.SetFloorIndicator(elev.Floor)
	elevio.SetDoorOpenLamp(true)
	
}

func DrivePastFloor(elev Elevator){
	elevio.SetFloorIndicator(elev.Floor)
} */