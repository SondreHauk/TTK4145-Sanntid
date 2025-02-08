package main

import (
	"fmt"
	"os"
	"os/signal"
	"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/lights"

	//"source/localElevator/requests"
	"time"
)

const SleepTime time.Duration = time.Millisecond*20
const (
    NUM_FLOORS=elevator.NUM_FLOORS
    NUM_BUTTONS=elevator.NUM_BUTTONS
)

func kill(){
    ch := make(chan os.Signal)
    signal.Notify(ch, os.Interrupt)
    
    //Blocks until an interrupt is received on ch
    select{
    case <-ch:
        elevio.SetMotorDirection(elevio.MD_Stop)
        os.Exit(1) //Terminates with error (Keyboard-interrupt etc)
    }
}

func main(){
    //Initializations
    elevio.Init("localhost:15657", elevator.NUM_FLOORS)
    elev := elevator.Elevator{}    
    lights.LightsInit(elev)
    elev.ElevatorInit()
    
    //Channels
    ButtonChan := make(chan elevio.ButtonEvent)
    ElevChan := make(chan elevator.Elevator)
    
    //Goroutines    
    go elevio.PollButtons(ButtonChan)
    go kill()

    //Blocking. Deadlock if no goroutines are running.
    select{}
}