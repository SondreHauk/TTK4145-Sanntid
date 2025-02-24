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
	"source/backup"
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

	PrimaryTXChan := make(chan string, 10)
	PrimaryRXChan := make(chan string, 10)
	//WorldviewTXChan := make(chan primary.Worldview, 10)

	BecomePrimary := make(chan bool)

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
	//Worldview := primary.Worldview{}

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
	go bcast.Transmitter(PORT_PRIMARY, PrimaryTXChan)
	go bcast.Receiver(PORT_PRIMARY, PrimaryRXChan)

	go backup.Run(PrimaryRXChan, BecomePrimary)
	go primary.Run(PeerUpdateChan, ElevatorRXChan, BecomePrimary, PrimaryTXChan/*, WorldviewTXChan, &Worldview*/)
	
	// Blocking select
	select {}
}
	
	//go assigner.TimeToIdle(elev)
	
/* 	el:=make([]Elevator,2)					// HALLUP HALLDWN  CAB
	el[0]=Elevator{Floor:2,Requests:[4][3]bool{{false, true, false}, //FLOOR 4
											   {false, false, true}, //FLOOR 3
											   {true, true, false}, //FLOOR 2
										       {false, false, true}, //FLOOR 1
											},
					PrevDirection:UP}
	el[1]=Elevator{Floor:3,Requests:[4][3]bool{{false, false, false},
											   {false, true, false},
											   {true, true, false},
										       {false, false, false},
											}}
	fmt.Println(assigner.ChooseElevator(el,Order{0,1}))
	fmt.Println("Time until pickup for el[0]: ",fsm.TimeUntilPickup(el[0],Order{0,1}))
	fmt.Println("Time until pickup for el[1]: ",fsm.TimeUntilPickup(el[1],Order{0,1})) */

	/* for {
		if elev.Floor!=-1 {
			fmt.Println("Time to idle: ",assigner.TimeToIdle(elev))
		}
		time.Sleep(time.Second)
		select {
		case msg := <-MsgChan:
			// Process and print received message
			fmt.Printf("Message received: ID = %d, Heartbeat = %s\n", msg.ID, msg.Heartbeat)
	} */
	//Primary backup protocol
	/*go backup(listens to bcast from primary) */

	/* //Blocking select
	select {
		/* 
		case primary dead
			if next in queue:
				go primary
		}*/
