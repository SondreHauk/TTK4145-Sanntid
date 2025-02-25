package primary

import (
	//"source/network/bcast"
	"time"
	"fmt"
	. "source/localElevator/config"
	"source/network/peers"
)

/*
TODO: 
- Replace string in primary active chan with a worldview struct to be sent to backups
- Fix problem with takeover printing whole history of elev state updates. Maybe peerUpdateChan is full when takeover happens.
*/

// var activePeers peers.PeerUpdate
// var elevators = make(map[string]Elevator)

type Worldview struct{
	ActivePeers peers.PeerUpdate
	Elevators map[string]Elevator
}

func Run(
	peerUpdateChan <-chan peers.PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimary <-chan bool,
	/*primaryActiveChan chan <- string,*/
	worldviewChan chan <- Worldview,
	/*worldview *Worldview*/){

	var worldview Worldview
	worldview.Elevators = make(map[string]Elevator)

	for {

		select{
		case <- becomePrimary:
			fmt.Println("Taking over as Primary")
			//drainElevatorStateUpdates(elevStateChan, &worldview.Elevators)
			HeartbeatTimer := time.NewTicker(T_HEARTBEAT)

			for{
				select{
				case worldview.ActivePeers = <-peerUpdateChan:
					printPeers(worldview.ActivePeers)
					
				case elevUpdate := <-elevStateChan:
					worldview.Elevators[elevUpdate.ID] = elevUpdate

					fmt.Println("Elevator State Updated")
					fmt.Printf("ID: %s\n", elevUpdate.ID)
					fmt.Printf("Floor: %d\n", elevUpdate.Floor)

				case <-HeartbeatTimer.C:
					//primaryActiveChan <- "Hello from Primary"
					worldviewChan <- worldview

				case <-becomePrimary:
					fmt.Println("Another Primary taking over...")
					break
				}
			}
		}
	}
}

// **Helper Function: Drain all pending updates before normal operation**
// func drainElevatorStateUpdates(elevStateChan <-chan Elevator, elevators *map[string]Elevator) {
// 	fmt.Println(" Draining old elevator state updates before taking over...")

// func drainElevatorStateUpdates(elevStateChan <-chan Elevator, elevators *map[string]Elevator) {
// 	for {
// 		select {
// 		case elevUpdate := <-elevStateChan:
// 			(*elevators)[elevUpdate.ID] = elevUpdate
// 		default:
// 			return // Exit when no more messages are available
// 		}
// 	}
// }


func printPeers(p peers.PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}