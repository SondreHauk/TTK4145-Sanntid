package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
	"source/backup"
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/inits"
	"source/localElevator/requests"
	"source/network/bcast"
	"source/network/peers"
	"source/primary"
)

func kill(StopButtonCh<-chan bool){
	KeyboardInterruptCh := make(chan os.Signal, 1)
	signal.Notify(KeyboardInterruptCh, os.Interrupt)
	
	select {
	case <-KeyboardInterruptCh:		
		fmt.Println("Keyboard interrupt")
	case <-StopButtonCh:
		for i:=0;i<5;i++{
			elevio.SetStopLamp(true)
			time.Sleep(T_BLINK)
			elevio.SetStopLamp(false)
			time.Sleep(T_BLINK)	
		}		
	}

	elevio.SetMotorDirection(elevio.MD_Stop)
	os.Exit(1) 
}

func main() {
	
	var port string
	var id string
	flag.StringVar(&port, "port", "", "Elevator port number")
	flag.StringVar(&id, "id","", "Elevator port")
	//TODO: If not valid ID, kill.
	flag.Parse()

	//Channels
	ElevatorTXChan := make(chan Elevator, 10)
	ElevatorRXChan := make(chan Elevator) 

	TransmitEnableChan := make(chan bool)
	PeerUpdateChan := make(chan peers.PeerUpdate)

	WorldviewTXChan := make(chan primary.Worldview, 10)
	WorldviewRXChan := make(chan primary.Worldview, 10)
	BecomePrimaryChan := make(chan bool)

	AtFloorChan := make(chan int, 1)
	ButtonChan := make(chan elevio.ButtonEvent, 10)
	ObstructionChan := make(chan bool, 1)
	StopChan := make(chan bool, 1)

	RequestToPrimaryChan := make(chan Order, 10)
	RequestFromElevChan := make(chan Order, 10)
	OrderToElevChan := make(chan Order, 10)
	OrderChan := make(chan Order, 10)
	
	//Initializations
	elevio.Init("localhost:"+ port, NUM_FLOORS)
	elev := Elevator{}
	inits.LightsInit()
	inits.ElevatorInit(&elev, id)

	// Goroutines Local elevator
	go requests.MakeRequest(ButtonChan, RequestToPrimaryChan, 
							OrderChan, id)
	go elevio.PollButtons(ButtonChan)
	go elevio.PollFloorSensor(AtFloorChan)
	go elevio.PollObstructionSwitch(ObstructionChan)
	go elevio.PollStopButton(StopChan)
	go kill(StopChan)
	go fsm.Run(&elev, ElevatorTXChan, AtFloorChan, 
				OrderChan, ObstructionChan, id)

	// Goroutines communication
	go bcast.Transmitter(PORT_ELEVSTATE, ElevatorTXChan)
	go bcast.Receiver(PORT_ELEVSTATE, ElevatorRXChan)
	go peers.Transmitter(PORT_PEERS, id, TransmitEnableChan)
	go peers.Receiver(PORT_PEERS, PeerUpdateChan)
	go bcast.Transmitter(PORT_WORLDVIEW, WorldviewTXChan)
	go bcast.Receiver(PORT_WORLDVIEW, WorldviewRXChan)
  
	// Elevator --- Request ---> Primary --- Order ---> Elevator
	go bcast.Transmitter(PORT_REQUEST, RequestToPrimaryChan)
	go bcast.Receiver(PORT_REQUEST, RequestFromElevChan)
	go bcast.Transmitter(PORT_ORDER, OrderToElevChan)
	go bcast.Receiver(PORT_ORDER, OrderChan)

	go backup.Run(WorldviewRXChan, BecomePrimaryChan, id)

	go primary.Run(PeerUpdateChan, ElevatorRXChan, 
					BecomePrimaryChan, WorldviewTXChan,
					RequestFromElevChan, OrderToElevChan, id)
	
	// Blocking select
	select {}
}