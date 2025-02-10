package main

import (
	"fmt"
	"os"
	"os/signal"
	. "source/localElevator/config"
	"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/lights"
	"source/localElevator/requests"
	"time"
)
var T_SLEEP = 20*time.Millisecond
func kill() {
	ch := make(chan os.Signal,1)
	signal.Notify(ch, os.Interrupt)
	//Blocks until an interrupt is received on ch
	select {
	case <-ch:
		elevio.SetMotorDirection(elevio.MD_Stop)
		os.Exit(1) //Terminates with error (Keyboard-interrupt etc)
	}
}

func main() {
	fmt.Println("")
	//Channels
	FsmChans := FsmChansType{
		ElevatorChan: make(chan Elevator),
        AtFloorChan:  make(chan int),
		NewOrderChan: make(chan Order),
	}

	ButtonChan := make(chan elevio.ButtonEvent)
	//Initializations
	elevio.Init("localhost:15657", NUM_FLOORS)
	elev := Elevator{}
	lights.LightsInit(elev)
	elevator.ElevatorInit(elev)
	
	//Goroutines
	
	go elevio.PollButtons(ButtonChan)
	go requests.Update(ButtonChan,FsmChans.NewOrderChan)
	go elevio.PollFloorSensor(FsmChans.AtFloorChan)
	go fsm.Run(elev, FsmChans)
	go kill()

	//Blocking. Deadlock if no goroutines are running.
	for{
		select{
			/* case a:= <-FsmChans.AtFloorChan:
				fmt.Println(a) */
			case a:= <-FsmChans.ElevatorChan:
				fmt.Println(a.Floor)
		
		}
		time.Sleep(T_SLEEP)
	}
	select{}
}
