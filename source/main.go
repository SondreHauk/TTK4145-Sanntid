package main

import (
	"flag"
	"source/backup"
	. "source/config"
	"source/localElevator/fsm"
	misc "source/miscellaneous"
	"source/network/bcast"
	"source/network/peers"
	"source/primary"
)

func main() {
	// Command line flags
	var port string
	var id string
	flag.StringVar(&port, "port", DEFAULT_ELEVIO_PORT, "Elevator port number")
	flag.StringVar(&id, "id", "1", "Elevator id")
	flag.IntVar(&NUM_FLOORS, "floors", DEFAULT_NUM_FLOORS, "Number of floors if simulating")
	flag.Parse()

	// Transmission/Receival channels
	enableTXChan := make(chan bool)
	elevTXChan := make(chan Elevator, 10)
	elevRXChan := make(chan Elevator, 10)
	wvTXChan := make(chan Worldview, 10)
	wvRXChan := make(chan Worldview, 10)
	requestsTXChan := make(chan Requests, 10)
	requestsRXChan := make(chan Requests, 10)

	// Local elevator channels
	stopChan := make(chan bool, 1)
	wvToElevChan := make(chan Worldview, 10)

	// Fault tolerance channels
	enablePrimaryChan := make(chan Worldview, 1)
	peerUpdateChan := make(chan PeerUpdate)
	wvToPrimaryChan := make(chan Worldview, 10)

	// Transmission/Receival goroutines
	go bcast.Transmitter(PORT_BCAST, elevTXChan, requestsTXChan, wvTXChan)
	go bcast.Receiver(PORT_BCAST, elevRXChan, requestsRXChan, wvRXChan)
	go peers.Transmitter(PORT_PEERS, id, enableTXChan)
	go peers.Receiver(PORT_PEERS, peerUpdateChan)

	// Local elevator protocol
	go fsm.Run(elevTXChan, requestsTXChan, wvToElevChan, stopChan, id, port)

	// Fault tolerance protocol
	go backup.Run(wvRXChan, wvToElevChan, enablePrimaryChan,
		wvToPrimaryChan, id)
	go primary.Run(peerUpdateChan, elevRXChan, enablePrimaryChan,
		wvTXChan, wvToPrimaryChan, requestsRXChan, id)

	// Terminates execution
	go misc.Kill(stopChan)
	select {}
}
