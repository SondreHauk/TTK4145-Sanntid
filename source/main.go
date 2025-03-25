package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"source/backup"
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/inits"
	"source/localElevator/requests"
	"source/network/bcast"
	"source/network/peers"
	"source/primary"
	"time"
)

func kill(StopButtonCh <-chan bool) {
	KeyboardInterruptCh := make(chan os.Signal, 1)
	signal.Notify(KeyboardInterruptCh, os.Interrupt)

	select {
	case <-KeyboardInterruptCh:
		fmt.Println("Keyboard interrupt")
	case <-StopButtonCh:
		for i := 0; i < 5; i++ {
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
	flag.StringVar(&id, "id", "", "Elevator port")
	flag.Parse()

	//Channels
	elevatorTXChan := make(chan Elevator, 10)
	elevatorRXChan := make(chan Elevator, 10)

	transmitEnableChan := make(chan bool)
	peerUpdateChan := make(chan PeerUpdate)

	worldviewTXChan := make(chan Worldview, 10)
	worldviewRXChan := make(chan Worldview, 10)
	becomePrimaryChan := make(chan Worldview, 1)

	worldviewToPrimaryChan := make(chan Worldview, 10)
	worldviewToElevatorChan := make(chan Worldview, 10)

	atFloorChan := make(chan int, 1)
	buttonChan := make(chan elevio.ButtonEvent, 10)
	obstructionChan := make(chan bool, 1)
	stopChan := make(chan bool, 1)

	requestsTXChan := make(chan Requests, 10)
	requestsRXChan := make(chan Requests, 10)
	accReqChan := make(chan OrderMatrix, 10)
	orderChan := make(chan Order, 10)

	//Initializations
	elevio.Init("localhost:" + port, NUM_FLOORS)
	elev := Elevator{}
	inits.LightsInit()
	inits.ElevatorInit(&elev, id)

	// Goroutines Local elevator
	go requests.SendRequest(buttonChan, requestsTXChan, accReqChan, id)
	go elevio.PollButtons(buttonChan)
	go elevio.PollFloorSensor(atFloorChan)
	go elevio.PollObstructionSwitch(obstructionChan)
	go elevio.PollStopButton(stopChan)
	go fsm.Run(&elev, elevatorTXChan, atFloorChan, orderChan,
		accReqChan, obstructionChan, worldviewToElevatorChan, id)

	// Goroutines communication
	go bcast.Transmitter(PORT_BCAST, elevatorTXChan, requestsTXChan)
	go bcast.Receiver(PORT_BCAST, elevatorRXChan, requestsRXChan)
	go peers.Transmitter(PORT_PEERS, id, transmitEnableChan)
	go peers.Receiver(PORT_PEERS, peerUpdateChan)

	go bcast.Transmitter(PORT_WORLDVIEW, worldviewTXChan)
	go bcast.Receiver(PORT_WORLDVIEW, worldviewRXChan)

	// Fault tolerance protocol
	go backup.Run(worldviewRXChan, worldviewToElevatorChan, becomePrimaryChan, worldviewToPrimaryChan, id)
	go primary.Run(peerUpdateChan, elevatorRXChan, becomePrimaryChan, 
		worldviewTXChan, worldviewToPrimaryChan, requestRXChan, id)

	// Kills terminal if interrupted
	go kill(stopChan)
	select {}
}
