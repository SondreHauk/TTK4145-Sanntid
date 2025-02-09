package requests

import (
	. "source/localElevator/config"
)

//import ("localElevator/elevio")
//This module should handle incoming requests and distribute them to the elevators
//c=[1,2,3], c[:]=c

func ClearFloor(elev Elevator, floor int){
	for btn:=0; btn<NUM_BUTTONS; btn++{
		elev.Requests[floor][btn]=false
	}
}

func ClearAll(elev Elevator){
	for fl:=0; fl<NUM_FLOORS; fl++{
		ClearFloor(elev,fl)
	}
}

func ShouldStop(elev Elevator){

}

func Update(elev Elevator, order Order){
	elev.Requests[order.Floor][order.Button]=true
}