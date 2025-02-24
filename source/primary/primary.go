package primary

import (
	//"source/network/bcast"
	"time"
	"fmt"
	. "source/localElevator/config"
	"source/network/peers"
)

func printPeers(p peers.PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}

func Run(
	peerUpdateChan <-chan peers.PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimary <-chan bool,
	primaryActiveChan chan <- string){
	
	var activePeers peers.PeerUpdate
	var elevators = make(map[string]Elevator)

	for {
		select{
		case <- becomePrimary:
			fmt.Println("Taking over as Primary")
			HeartbeatTimer := time.NewTicker(T_HEARTBEAT)
			for{
				select{
				case activePeers = <-peerUpdateChan:
					printPeers(activePeers)
					
				case elevUpdate := <-elevStateChan:
					elevators[elevUpdate.ID] = elevUpdate

					fmt.Println("Elevator State Updated")
					fmt.Printf("ID: %s\n", elevUpdate.ID)
					fmt.Printf("Floor: %d\n", elevUpdate.Floor)

				case <- HeartbeatTimer.C:
					primaryActiveChan <- "Hello from Primary"
				}
			}
		}
	}
}