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
	flag.Parse()
	
	//Channels
	elevatorTXChan := make(chan Elevator, 10)
	elevatorRXChan := make(chan Elevator) 

	transmitEnableChan := make(chan bool)
	peerUpdateChan := make(chan peers.PeerUpdate)

	worldviewTXChan := make(chan primary.Worldview, 10)
	worldviewRXChan := make(chan primary.Worldview, 10)
	becomePrimaryChan := make(chan primary.Worldview, 1)

	hallLightsTXChan := make(chan HallLights, 10)
	hallLightsRXChan := make(chan HallLights, 10)

	atFloorChan := make(chan int, 1)
	buttonChan := make(chan elevio.ButtonEvent, 10)
	obstructionChan := make(chan bool, 1)
	stopChan := make(chan bool, 1)

	requestToPrimaryChan := make(chan Order, 10)
	requestFromElevChan := make(chan Order, 10)
	orderToElevChan := make(chan Order, 10)
	orderChan := make(chan Order, 10)
	
	//Initializations
	elevio.Init("localhost:" + port, NUM_FLOORS)
	elev := Elevator{}
	inits.LightsInit()
	inits.ElevatorInit(&elev, id)

	// Goroutines Local elevator
	go requests.MakeRequest(buttonChan, requestToPrimaryChan, orderChan, id)
	go elevio.PollButtons(buttonChan)
	go elevio.PollFloorSensor(atFloorChan)
	go elevio.PollObstructionSwitch(obstructionChan)
	go elevio.PollStopButton(stopChan)
	go fsm.Run(&elev, elevatorTXChan, atFloorChan, 
				orderChan, hallLightsRXChan, obstructionChan, id)

	// Goroutines communication
	go bcast.Transmitter(PORT_ELEVSTATE, elevatorTXChan)
	go bcast.Receiver(PORT_ELEVSTATE, elevatorRXChan)
	go peers.Transmitter(PORT_PEERS, id, transmitEnableChan)
	go peers.Receiver(PORT_PEERS, peerUpdateChan)
	go bcast.Transmitter(PORT_WORLDVIEW, worldviewTXChan)
	go bcast.Receiver(PORT_WORLDVIEW, worldviewRXChan)
  
	// Elevator --- Request ---> Primary --- Order ---> Elevator
	go bcast.Transmitter(PORT_REQUEST, requestToPrimaryChan)
	go bcast.Receiver(PORT_REQUEST, requestFromElevChan)
	go bcast.Transmitter(PORT_ORDER, orderToElevChan)
	go bcast.Receiver(PORT_ORDER, orderChan)
	go bcast.Transmitter(PORT_HALLLIGHTS, hallLightsTXChan)
	go bcast.Receiver(PORT_HALLLIGHTS, hallLightsRXChan)

	// Fault tolerance protocol
	go backup.Run(worldviewRXChan, becomePrimaryChan, id)
	go primary.Run(peerUpdateChan, elevatorRXChan, 
					becomePrimaryChan, worldviewTXChan,
					requestFromElevChan, orderToElevChan, 
					hallLightsTXChan, id)
	
	// Kills terminal if interrupted
	go kill(stopChan)
	select {}
}