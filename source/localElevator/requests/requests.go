package requests

import (
	"source/localElevator/elevator"
	"source/localElevator/elevio"
)

//import ("localElevator/elevio")
//This module should handle incoming requests and distribute them to the elevators
//c=[1,2,3], c[:]=c

func ClearCurrentFloor(elev elevator.Elevator){

}

func ShouldStop(elev elevator.Elevator){

}

func UpdateQueue(elev elevator.Elevator){
	input:=elevio.GetButton()
}