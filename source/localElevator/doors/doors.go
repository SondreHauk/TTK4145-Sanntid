package doors

import "source/localElevator/elevio"

//Shell function that does not alter state
func Open(){
	elevio.SetDoorOpenLamp(true)
}

//Shell function that does not alter state
func Close(){
	elevio.SetDoorOpenLamp(false)
}
