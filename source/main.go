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

// func worldviewRouter(worldviewRXChan <-chan Worldview,
// 	/*worldviewToPrimaryChan chan<- Worldview,*/
// 	worldviewToBackupChan chan<- Worldview,
// 	worldviewToElevatorChan chan<- Worldview) {
// 	for wv := range worldviewRXChan {
// 		worldviewToBackupChan <- wv
// 		/*worldviewToPrimaryChan <- wv*/
// 		worldviewToElevatorChan <- wv
// 	}
// }

// func worldviewRouter(worldviewRXChan <-chan Worldview,
// 	/*worldviewToPrimaryChan chan<- Worldview,*/
// 	worldviewToBackupChan chan<- Worldview,
// 	worldviewToElevatorChan chan<- Worldview) {

// 	for wv := range worldviewRXChan {
// 		select {
// 		case worldviewToBackupChan <- wv:
// 			fmt.Println("Sent to backup")
// 		default:
// 			fmt.Println("Warning: Dropped worldviewToBackup due to full channel")
// 		}

// 		select {
// 		case worldviewToElevatorChan <- wv:
// 			// fmt.Println("Sent to elev")
// 		default:
// 			fmt.Println("Warning: Dropped worldviewToElevator due to full channel")
// 		}
// 	}
// }

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

	// worldviewToPrimaryChan := make(chan Worldview, 10)
	// worldviewToBackupChan := make(chan Worldview, 10)
	worldviewToElevatorChan := make(chan Worldview, 10)

	atFloorChan := make(chan int, 1)
	buttonChan := make(chan elevio.ButtonEvent, 10)
	obstructionChan := make(chan bool, 1)
	stopChan := make(chan bool, 1)

	requestToPrimaryChan := make(chan Order, 10)
	requestFromElevChan := make(chan Order, 10)
	/*orderToElevChan := make(chan Order, 10)*/
	orderChan := make(chan Order, 10)

	//Initializations
	elevio.Init("localhost:"+port, NUM_FLOORS)
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
		orderChan, /*hallLightsRXChan,*/ obstructionChan, 
		worldviewToElevatorChan, id)

	// Goroutines communication (TODO: reduce to two ports)
	go bcast.Transmitter(PORT_BCAST, elevatorTXChan, requestToPrimaryChan, worldviewTXChan)
	go bcast.Receiver(PORT_BCAST, elevatorRXChan, requestFromElevChan, worldviewRXChan)
	go peers.Transmitter(PORT_PEERS, id, transmitEnableChan)
	go peers.Receiver(PORT_PEERS, peerUpdateChan)

	// go worldviewRouter(worldviewRXChan, /*worldviewToPrimaryChan,*/ worldviewToBackupChan, worldviewToElevatorChan)

	//TODO: DRAIN CHANNELS GOING TO PRIMARY
	
	// Fault tolerance protocol
	go backup.Run(worldviewRXChan, worldviewToElevatorChan, becomePrimaryChan, id)
	go primary.Run(peerUpdateChan, elevatorRXChan,
		becomePrimaryChan, worldviewTXChan, /*worldviewToPrimaryChan,*/
		requestFromElevChan, /*orderToElevChan,*/
		/*hallLightsTXChan,*/ id)

	// Kills terminal if interrupted
	go kill(stopChan)
	select {}
}
