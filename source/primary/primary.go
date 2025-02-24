package primary

import (
	//"source/network/bcast"
	"fmt"
	"source/localElevator/config"
	. "source/localElevator/config"
	"source/network/peers"
)

// func MsgBcastTX(msg chan Message, id int){
// 	//go bcast.Transmitter(port, msg) // Start broadcasting in a separate goroutine
// 	for {
// 		msg <- Message{ID: id, Heartbeat: "Alive"}
// 		time.Sleep(T_HEARTBEAT)
// 	}
// }

// var Elevators = make(map[string]*Elevator)
// var ActivePeers peers.PeerUpdate

func printPeers(p peers.PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}

func Run(
	peerUpdateChan <-chan peers.PeerUpdate,
	elevStateChan <-chan config.Elevator){
	
	var activePeers peers.PeerUpdate
	var elevators = make(map[string]Elevator)
	
	for {
		select{
		case activePeers = <-peerUpdateChan:
			printPeers(activePeers)
			
		case elevUpdate := <-elevStateChan:
			elevators[elevUpdate.ID] = elevUpdate

			fmt.Println("Elevator State Updated")
			fmt.Printf("ID: %s\n", elevUpdate.ID)
			fmt.Printf("Floor: %d\n", elevUpdate.Floor)
			//fmt.Printf("Direction: %s\n", directionToString(elevUpdate.Direction))
			//fmt.Printf("State: %s\n", stateToString(elevUpdate.State))
			//fmt.Println("Requests:")
			//printRequests(elevUpdate.Requests)

		}		
	}
		//case neworder = <- newOrderReceived:
			//assignelev
			//bcast order to right elev
			//wait for ack
}