package primary

import (
	//"source/network/bcast"
	. "source/localElevator/config"
	"time"
	"source/network/peers"
	"fmt"
)

func MsgBcastTX(msg chan Message, id int){
	//go bcast.Transmitter(port, msg) // Start broadcasting in a separate goroutine
	for {
		msg <- Message{ID: id, Heartbeat: "Alive"}
		time.Sleep(T_HEARTBEAT)
	}
}

func printPeers(p peers.PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}

func Run(peerUpdateChan <-chan peers.PeerUpdate){
	var activePeers peers.PeerUpdate
	for {
		select{
		case activePeers = <-peerUpdateChan:
			printPeers(activePeers)
		//case neworder = <- newOrderReceived:
			//assignelev
			//bcast order to right elev
			//wait for ack
		}
	}
}