package main

import (
	//"fmt"
	"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/localElevator/lights"
    "source/localElevator/fsm"
	//"source/localElevator/requests"
	"time"
)

const SleepTime time.Duration = time.Millisecond*20
const (
    NUM_FLOORS=elevator.NUM_FLOORS
    NUM_BUTTONS=elevator.NUM_BUTTONS
)

func main(){
    elevio.Init("localhost:15657", elevator.NUM_FLOORS)
    
    elev := elevator.Elevator{}    
    lights.LightsInit(elev)
    elev.ElevatorInit()
    ButtonChan := make(chan elevio.ButtonEvent)
    
    for{
        // One part for Requests
        //requests.Update()
       /*  PrevRequests:= make([][]int,NUM_FLOORS)
        for f:=0; f<NUM_FLOORS;f++{
            for b:=0; b<NUM_BUTTONS; b++{
                inp:=elevio.GetButton()
            }
        } */
        go elevio.PollButtons(ButtonChan)
        select{
        case ButtonPress:=<-ButtonChan:
            fsm.OnButtonPress(ButtonPress)
        }

        // One part for FloorSensor
        time.Sleep(SleepTime)
    }
}