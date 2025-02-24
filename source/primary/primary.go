package primary

import (
	//"source/network/bcast"
	"time"
	"fmt"
	. "source/localElevator/config"
	"source/network/peers"
)

/*
Seems to work all right! A bit strange behavior when it enters the second for-loop: 
at one instant it print outs many "Elevator State Updated", but then seem to operate normal again,
printing at regular time intervals.
TODO: Replace string in primary active chan with a worldview struct to be sent to backups
*/

type Worldview struct{
	ActivePeers peers.PeerUpdate
	Elevators map[string]Elevator
}

func Run(
	peerUpdateChan <-chan peers.PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimary <-chan bool,
	primaryActiveChan chan <- string,
	worldviewChan chan <- Worldview,
	worldview *Worldview){
	
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
					worldviewChan <- *worldview
				}
			}
		}
	}
}

func printPeers(p peers.PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}