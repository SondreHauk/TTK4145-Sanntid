package primary

import (
	//"source/network/bcast"
	"time"
	"fmt"
	. "source/localElevator/config"
	"source/network/peers"
)

type Worldview struct{
	PrimaryId string
	ActivePeers peers.PeerUpdate
	Elevators map[string]Elevator
}

func Run(
	peerUpdateChan <-chan peers.PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimary <-chan bool,
	worldviewChan chan <- Worldview,
	id string){

	var worldview Worldview
	worldview.Elevators = make(map[string]Elevator)
	worldview.PrimaryId = id

	for{
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
					worldview.Elevators[elevUpdate.Id] = elevUpdate
					printElevator(elevUpdate)

				case <- Neworder:

				case <-HeartbeatTimer.C:
					worldviewChan <- worldview

				case <-becomePrimary: // Should be deleted at some point
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
func printElevator(e Elevator){
	fmt.Println("Elevator State Updated")
	fmt.Printf("ID: %s\n", e.Id)
	fmt.Printf("Floor: %d\n", e.Floor)
}

func printPeers(p peers.PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
}