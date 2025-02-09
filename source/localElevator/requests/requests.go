package requests

import (
	. "source/localElevator/config"
	"source/localElevator/elevio"
)

//import ("localElevator/elevio")
//This module should handle incoming requests and distribute them to the elevators
//c=[1,2,3], c[:]=c

func ClearCurrentFloor(elev Elevator){
	for btn:=0; btn<NUM_BUTTONS; btn++{
		elev.Requests[elev.Floor][btn]=false
	}
}

func ShouldStop(elev Elevator){

}

/* func UpdateQueue(elev Elevator){
	input:=elevio.GetButton()
} */