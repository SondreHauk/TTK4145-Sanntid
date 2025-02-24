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
	"source/primary"
	"source/network/bcast"
	"source/network/peers"
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
	// Initialize from command line with: go run main.go -port=15657 -id=1
	var port string
	var id string
	flag.StringVar(&port, "port", "", "Elevator port number")
	flag.StringVar(&id, "id","", "Elevator port")
	//TODO: If not valid ID, kill.
	flag.Parse()

	//Channels
	ElevatorTXChan := make(chan Elevator, 10)
	ElevatorRXChan := make(chan Elevator)

	TransmitEnable := make(chan bool)
	PeerUpdateChan := make(chan peers.PeerUpdate)

	AtFloorChan := make(chan int, 1)
	NewOrderChan := make(chan Order, 10)
	ButtonChan := make(chan elevio.ButtonEvent, 10)
	ObstructionChan := make(chan bool, 1)
	StopChan := make(chan bool, 1)
	
	//Initializations
	elevio.Init("localhost:"+ port, NUM_FLOORS)
	elev := Elevator{}
	inits.LightsInit()
	inits.ElevatorInit(&elev, id)

	// Goroutines Local elevator
	go requests.Update(ButtonChan, NewOrderChan)
	go elevio.PollButtons(ButtonChan)
	go elevio.PollFloorSensor(AtFloorChan)
	go elevio.PollObstructionSwitch(ObstructionChan)
	go elevio.PollStopButton(StopChan)
	go fsm.Run(&elev, ElevatorTXChan, AtFloorChan, NewOrderChan, ObstructionChan)
	go kill(StopChan)

	// Goroutines communication
	go bcast.Transmitter(PORT_BCAST_ELEV, ElevatorTXChan)
	go bcast.Receiver(PORT_BCAST_ELEV, ElevatorRXChan)
	go peers.Transmitter(PORT_PEERS, id, TransmitEnable)
	go peers.Receiver(PORT_PEERS, PeerUpdateChan)

	go primary.Run(PeerUpdateChan,ElevatorRXChan)
	
	// Blocking select
	select {}
}