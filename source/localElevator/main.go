package main

import (
	"flag"
	"os"
	"os/signal"
	. "source/localElevator/config"
	"source/localElevator/elevator"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/lights"
	"source/localElevator/requests"
	"fmt"
)

func kill() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	//Blocks until an interrupt is received on ch
	select {
	case <-ch:
		fmt.Println("Interrupt received")
		elevio.SetMotorDirection(elevio.MD_Stop)
		os.Exit(1) //Terminates with error (Keyboard-interrupt etc)
	}
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
	
	//Initializations
	elevio.Init("localhost:"+ port, NUM_FLOORS)
	elev := Elevator{}
	lights.LightsInit(elev)
	elevator.ElevatorInit(elev)

	//Goroutines
	go elevio.PollButtons(ButtonChan)
	go requests.Update(ButtonChan, NewOrderChan)
	go elevio.PollFloorSensor(AtFloorChan)
	go fsm.Run(&elev, /* ElevatorChan, */ AtFloorChan, NewOrderChan)
	go kill()

	//Blocking select
	select {}
}
