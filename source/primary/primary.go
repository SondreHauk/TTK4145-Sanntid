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
	hallLightsChan chan <- HallLights,
	id string){

	var worldview Worldview
	worldview.Elevators = make(map[string]Elevator)
	worldview.PrimaryId = id

	updateLights := new(bool)
	
	//Init hall lights matrix
	hallLights := make([][] bool, NUM_FLOORS)
	for i := range(hallLights){
		hallLights[i] = make([]bool, NUM_BUTTONS - 1)
	}

	select{
	case <- becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		drain(elevStateChan) //FIX FLUSHING OF CHANNELS
		HeartbeatTimer := time.NewTicker(T_HEARTBEAT)

		for{
			select{
			case worldview.PeerInfo = <-peerUpdateChan:
				//If elev lost: Reassign lost orders
				//printPeers(worldview.PeerInfo)

			case elevUpdate := <-elevStateChan:
				worldview.Elevators[elevUpdate.Id] = elevUpdate
				//Not working properly
				updateHallLights(worldview, hallLights, updateLights)
				if (*updateLights){
					hallLightsChan <- hallLights}

			case request := <- requestFromElevChan:
				//fmt.Printf("Request received from id: %s \n", request.Id)
				AssignedId := assigner.ChooseElevator(worldview.Elevators,
													worldview.PeerInfo.Peers,
													request)
				orderToElevChan <- Order{Id: AssignedId, 
											Floor: request.Floor,
											Button: request.Button}
				//fmt.Printf("Order sent to id: %s \n", AssignedId)
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
//NOT WORKING PROPERLY
func updateHallLights(wv Worldview, 
					hallLights [][]bool,
					updateHallLights *bool){

	*updateHallLights = false // Reset flag

	// Create a deep copy of hallLights (to properly compare changes)
	prevHallLights := make([][]bool, NUM_FLOORS)
	for i := range hallLights {
		prevHallLights[i] = make([]bool, NUM_BUTTONS-1)
		copy(prevHallLights[i], hallLights[i]) // Copy row data
	}

	// Reset hallLights matrix (assume no lights first, then set needed ones)
	for floor := range hallLights {
		for btn := range hallLights[floor] {
			hallLights[floor][btn] = false
		}
	}

	// Update hallLights based on the order matrix from all peers
	for _, id := range(wv.PeerInfo.Peers){
		orderMatrix := wv.Elevators[id].Orders
		for floor, floorOrders := range(orderMatrix){
			for btn, isOrder := range(floorOrders){
				if isOrder && btn!=2{
					hallLights[floor][btn] = hallLights[floor][btn] || isOrder
				}
			}
		}
	}
	// Compare hallLights with prevHallLights
	for floor := 0; floor < NUM_FLOORS; floor++ {
        for btn := 0; btn < NUM_BUTTONS-1; btn++ {
            if hallLights[floor][btn] != prevHallLights[floor][btn] {
                *updateHallLights = true
            }
        }
    }
}

func drain(ch <- chan Elevator){
	for len(ch) > 0{
		<- ch
	}
}

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