package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	//"strconv"
	. "source/localElevator/config"
	"source/localElevator/elevio"
	"source/localElevator/fsm"
	"source/localElevator/inits"
	"source/localElevator/requests"
	//"source/backup"
	//"source/primary"
	//"source/network/bcast"
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
	inits.ElevatorInit(&elev, id)

	//Goroutines
	go requests.Update(ButtonChan, NewOrderChan)
	go elevio.PollButtons(ButtonChan)
	go elevio.PollFloorSensor(AtFloorChan)
	go elevio.PollObstructionSwitch(ObstructionChan)
	go elevio.PollStopButton(StopChan)
	go fsm.Run(&elev, /* ElevatorChan, */ AtFloorChan, NewOrderChan, ObstructionChan)
	go kill(StopChan)

	//UDP bcast testing
	// TXchan := make(chan Message, 100)
	// RXchan := make(chan Message, 100)
	// i, _ := strconv.Atoi(id)
	// go backup.MsgBcastRX(20020, RXchan)
	// go primary.MsgBcastTX(TXchan, int(i))
	// go bcast.Transmitter(20020, TXchan)

	//peers testing
	TransmitEnable := make(chan bool)
	PeerUpdateChan := make(chan peers.PeerUpdate)

	go peers.Transmitter(20030, id, TransmitEnable)
	go peers.Receiver(20030, PeerUpdateChan)

	for {
		select {
		case p := <-PeerUpdateChan:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)
	}
	}
	select {}
}