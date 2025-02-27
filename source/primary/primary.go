package primary

import (
	"fmt"
	. "source/config"
	"source/network/peers"
	"source/primary/assigner"
	"time"
)

type Worldview struct{
	PrimaryId string
	PeerInfo peers.PeerUpdate
	Elevators map[string]Elevator
}

func Run(
	peerUpdateChan <-chan peers.PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimaryChan <-chan bool,
	worldviewChan chan <- Worldview,
	requestFromElevChan <- chan Order,
	orderToElevChan chan <- Order,
	/*hallLightschan chan <- Halllights*/ 
	id string){

	var worldview Worldview
	worldview.Elevators = make(map[string]Elevator)
	worldview.PrimaryId = id

	select{
	case <- becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		drain(elevStateChan) //FIX FLUSHING OF CHANNELS
		HeartbeatTimer := time.NewTicker(T_HEARTBEAT)

		for{
			select{
			case worldview.PeerInfo = <-peerUpdateChan:
				//If elev lost: Reassign lost orders
				printPeers(worldview.PeerInfo)

				
			case elevUpdate := <-elevStateChan:
				worldview.Elevators[elevUpdate.Id] = elevUpdate
				//printElevator(elevUpdate)
				//Check if the elevUpdate includes the order sent from primary
				//If elevUpdate.Id == Order.Id && "elevUpdate.OrderMatrix == Order" 
				//Set hall-lights!
				//hallLightsChan <- hallLights

			case request := <- requestFromElevChan:
				fmt.Printf("Request received from id: %s \n", request.Id)
				AssignedId := assigner.ChooseElevator(worldview.Elevators,
													worldview.PeerInfo.Peers,
													request)
				orderToElevChan <- Order{Id: AssignedId, 
											Floor: request.Floor,
											Button: request.Button}
				fmt.Printf("Order sent to id: %s \n", AssignedId)
				//Start a timer. If no elevUpdate is received from the assigned 
				//elev within timeout, decelar it dead and reassign orders!


			case <-HeartbeatTimer.C:
				worldviewChan <- worldview

			case <-becomePrimaryChan: //Needs logic
				fmt.Println("Another Primary taking over...")
				break
			}
		}
	}
}

func setHallLights(elevators map[string]Elevator){
	// for _,id := range(worldview.Elevators)
	// orders = worldview.Elevators
	// Iterate through the order matrix of each Elevator,
	// Make a light-in-hall-matrix, send it to elevs
	// fsm.Run() sets lights accordingly.
}

// **Helper Function: Drain all pending updates before normal operation**
// func drainElevatorStateUpdates(elevStateChan <-chan Elevator, elevators *map[string]Elevator) {
// 	fmt.Println(" Draining old elevator state updates before taking over...")

func drain(ch <- chan Elevator){
	for len(ch) > 0{
		<- ch
	}
}

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