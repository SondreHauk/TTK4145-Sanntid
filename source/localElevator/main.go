package main

import (
	"os"
	"os/signal"
	"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/lights"
	. "source/localElevator/config"
)

func kill() {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	//Blocks until an interrupt is received on ch
	select {
	case <-ch:
		elevio.SetMotorDirection(elevio.MD_Stop)
		os.Exit(1) //Terminates with error (Keyboard-interrupt etc)
	}
}

func main() {
	//Channels
	FsmChans := fsm.FsmChansType{
		ElevatorChan: make(chan Elevator),
        AtFloorChan:  make(chan int),
	}

	ButtonChan := make(chan elevio.ButtonEvent)

	//Initializations
	elevio.Init("localhost:15657", NUM_FLOORS)
	elev := Elevator{}
	lights.LightsInit(elev)
	elevator.ElevatorInit(elev)

	//Goroutines
	go elevio.PollButtons(ButtonChan)
	go elevio.PollFloorSensor(FsmChans.AtFloorChan)
	go fsm.Run(elev, FsmChans)
	go lights.Update(elev)
	go kill()

	//Blocking. Deadlock if no goroutines are running.
	select {}
}
