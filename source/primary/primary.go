package primary

import (
	"fmt"
	. "source/config"
	"source/localElevator/elevio"
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
	becomePrimaryChan <-chan Worldview,
	worldviewChan chan <- Worldview,
	requestFromElevChan <- chan Order,
	orderToElevChan chan <- Order,
	hallLightsChan chan <- HallLights,
	id string){

	// Local variables
	updateLights := new(bool)
	var worldview Worldview
	worldview.Elevators = make(map[string]Elevator)
	obstructedElevators := make([]string, NUM_ELEVATORS)

	// Init hall lights matrix
	hallLights := make([][] bool, NUM_FLOORS)
	for i := range(hallLights){
		hallLights[i] = make([]bool, NUM_BUTTONS - 1)
	}

	select{
	case worldview := <-becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		//drain(elevStateChan) //FIX FLUSHING OF CHANNELS
		heartbeatTimer := time.NewTicker(T_HEARTBEAT)
		obstructionTimers := make(map[string]*time.Timer)

		for{
			select{
			case worldview.PeerInfo = <-peerUpdateChan:
				//If elev lost: Reassign lost orders
				printPeers(worldview.PeerInfo)
				lost:=worldview.PeerInfo.Lost
				if len(lost)!=0{
					ReassignHallOrders(worldview, orderToElevChan, ConnectionLost, "")
				}

			case elevUpdate := <-elevStateChan:
				worldview.Elevators[elevUpdate.Id] = elevUpdate
				//Not working properly?
				updateHallLights(worldview, hallLights, updateLights)
				if (*updateLights){
					hallLightsChan <- hallLights}

				// ------ OBSTRUCTION ------- //
				//If elevator is obstructed for 3 sec, reassign hall orders
				if elevUpdate.Obstructed {
					obstructedElevators = append(obstructedElevators, elevUpdate.Id)
					//start timer
					if _, exists := obstructionTimers[elevUpdate.Id]; !exists{
						timer := time.AfterFunc(T_OBSTRUCTED_PRIMARY, func() {
							obstructedId := obstructedElevators[len(obstructedElevators)-1]
							ReassignHallOrders(worldview, orderToElevChan, Obstructed, obstructedId)
						})
						obstructionTimers[elevUpdate.Id] = timer
					}
				} else {
					//if ID in obstructedElevatorIds, pop id and stop timer
					// If the elevator is no longer obstructed, check if its ID is in the list of obstructed elevators
					for i, id := range obstructedElevators {
						if id == elevUpdate.Id {
							// If found, remove it from the slice
							obstructedElevators = append(obstructedElevators[:i], obstructedElevators[i+1:]...)
							//obstructedElevators = slices.Delete(obstructedElevators,i,i+1)
							// Stop the timer if it's active
							if timer, exists := obstructionTimers[elevUpdate.Id]; exists {
								timer.Stop()
								delete(obstructionTimers, elevUpdate.Id)
							}
							break
						}
					}
				}

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

			case <-heartbeatTimer.C:
				worldviewChan <- worldview

			case <-becomePrimaryChan: //Needs logic
				fmt.Println("Another Primary taking over...")
				break
			}
		}
	}
}

// Updates Hall Lights, returns true if new light matrix is different from old 
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
				if isOrder && btn!= int(elevio.BT_Cab){
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

func ReassignHallOrders(wv Worldview, orderToElevChan chan<- Order, situation int, id string){
switch situation {
case Obstructed:
	orderMatrix := wv.Elevators[id].Orders
		for floor, floorOrders := range(orderMatrix){
			for btn, isOrder := range(floorOrders){
				if isOrder && btn != int(elevio.BT_Cab){
					lostOrder:=Order{
								Id: id,
								Floor: floor,
								Button: btn,
							}
					lostOrder.Id = assigner.ChooseElevator(wv.Elevators,wv.PeerInfo.Peers,lostOrder)
					orderToElevChan <- lostOrder
				}
			}
		}	
case ConnectionLost:
	for _,lostId := range(wv.PeerInfo.Lost){
		orderMatrix := wv.Elevators[lostId].Orders
		for floor, floorOrders := range(orderMatrix){
			for btn, isOrder := range(floorOrders){
				if isOrder && btn != int(elevio.BT_Cab){
					lostOrder:=Order{
								Id: lostId,
								Floor: floor,
								Button: btn,
							}
					lostOrder.Id = assigner.ChooseElevator(wv.Elevators,wv.PeerInfo.Peers,lostOrder)
					orderToElevChan <- lostOrder
				}
			}
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