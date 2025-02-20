package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	. "source/localElevator/config"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/inits"
	"source/localElevator/requests"
)

func kill(StopButtonCh<-chan bool){
	KeyboardInterruptCh := make(chan os.Signal, 1)
	signal.Notify(KeyboardInterruptCh, os.Interrupt)
	
	select {
	case <-KeyboardInterruptCh:		
		fmt.Println("Keyboard interrupt")
	case <-StopButtonCh:
		fmt.Println("Stop button pressed")	
	}

	elevio.SetMotorDirection(elevio.MD_Stop)
	os.Exit(1) 
}

func main() {
	var port string
	flag.StringVar(&port, "port", "", "Elevator port number")
	flag.Parse()

	//Channels
	/* ElevatorChan := make(chan *Elevator, 10) */
	AtFloorChan := make(chan int, 1)
	NewOrderChan := make(chan Order, 10)
	ButtonChan := make(chan elevio.ButtonEvent, 10)
	ObstructionChan := make(chan bool, 1)
	StopChan := make(chan bool, 1)
	
	//Initializations
	elevio.Init("localhost:"+ port, NUM_FLOORS)
	elev := Elevator{}
	inits.LightsInit()
	inits.ElevatorInit(&elev)

	//Goroutines
	go requests.Update(ButtonChan, NewOrderChan)
	go elevio.PollButtons(ButtonChan)
	go elevio.PollFloorSensor(AtFloorChan)
	go elevio.PollObstructionSwitch(ObstructionChan)
	go elevio.PollStopButton(StopChan)
	go fsm.Run(&elev, /* ElevatorChan, */ AtFloorChan, NewOrderChan, ObstructionChan)
	go kill(StopChan)

	//Primary backup protocol
	/*go backup(listens to bcast from primary) */

	//Blocking select
	select {
		/* 
		case primary dead
			if next in queue:
				go primary
		*/
		
	}
	//
}
