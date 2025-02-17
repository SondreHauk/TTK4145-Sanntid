package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	. "source/localElevator/config"
	"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/lights"
	"source/localElevator/requests"
)

func kill(StopButtonCh<-chan bool){
	KeyboardInterruptCh := make(chan os.Signal, 1)
	signal.Notify(KeyboardInterruptCh, os.Interrupt)
	//Blocks until an interrupt is received on ch
	select {
	case <-KeyboardInterruptCh:
	case <-StopButtonCh:
	}
	fmt.Println("Interrupt received")
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
	lights.LightsInit()
	elevator.ElevatorInit(&elev)

	//Goroutines
	go requests.Update(ButtonChan, NewOrderChan)
	go elevio.PollButtons(ButtonChan)
	go elevio.PollFloorSensor(AtFloorChan)
	go elevio.PollObstructionSwitch(ObstructionChan)
	go elevio.PollStopButton(StopChan)
	go fsm.Run(&elev, /* ElevatorChan, */ AtFloorChan, NewOrderChan, ObstructionChan)
	go kill(StopChan)

	//Blocking select
	select {}
}
