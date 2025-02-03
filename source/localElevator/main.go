package main

import (
	. "fmt"
	"localElevator/elevator"
	"localElevator/elevio"
	"localElevator/lights"
	"time"
)

const(
    SLEEP time.Duration = time.Millisecond*20
)
func Sleep(){
    time.Sleep(SLEEP)
}

func main(){
    Println()

    elevio.Init("localhost:15657", elevator.NUM_FLOORS)
    

    elev := elevator.Elevator{}    
    lights.LightsInit(elev)
    elev.ElevatorInit()
   
    //elev.MoveFloor(2)
}
