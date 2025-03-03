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
	ElevatorSnapshot map[string]Elevator // Used for transmission, not valid for concurrent access
}

type ElevMapAction struct{
	cmd string
	id string
	elev Elevator
	initMap map[string]Elevator
	respCh chan any
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

	ElevMapChan := make(chan ElevMapAction, 20)
	
	updateLights := new(bool)
	var worldview Worldview
	worldview.ElevatorSnapshot = make(map[string]Elevator)
	
	//Init hall lights matrix
	hallLights := make([][] bool, NUM_FLOORS)
	for i := range(hallLights){
		hallLights[i] = make([]bool, NUM_BUTTONS - 1)
	}

	// waitloop:
	// for{
	// 	select{
	// 	case <- peerUpdateChan:
	// 	case <- elevStateChan:
	// 	case <- requestFromElevChan:
	// 	case worldview = <-becomePrimaryChan:
	// 		break waitloop
	// }
	// //Primary mode
	// fmt.Println("Taking over as Primary")
	// HeartbeatTimer := time.NewTicker(T_HEARTBEAT)

	select{
	case worldview := <-becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		ElevMapChan <- ElevMapAction{cmd:"initialize", initMap: worldview.ElevatorSnapshot}
		
		//drain(elevStateChan) //FIX FLUSHING OF CHANNELS
		HeartbeatTimer := time.NewTicker(T_HEARTBEAT)

		for{
			select{
			case worldview.PeerInfo = <-peerUpdateChan:
				//If elev lost: Reassign lost orders
				printPeers(worldview.PeerInfo)
				lost:=worldview.PeerInfo.Lost
				if len(lost)!=0{
					ReassignHallOrders(worldview, ElevMapChan, orderToElevChan)
				}

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
					updateHallLights *bool,
					MapAccessChan chan ElevMapAction){

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

	readChan := make(chan any, 1)
	defer close(readChan)
	//Request read
	MapAccessChan <- ElevMapAction{cmd: "read", respCh: readChan}
	
	// Update hallLights based on the order matrix from all peers
	select{
	case elevMap:= <-readChan:
		wv=Worldview{wv.PrimaryId, wv.PeerInfo, elevMap.(map[string]Elevator)}
		for _, id := range(wv.PeerInfo.Peers){
			orderMatrix := wv.ElevatorSnapshot[id].Orders
			for floor, floorOrders := range(orderMatrix){
				for btn, isOrder := range(floorOrders){
					if isOrder && btn!= int(elevio.BT_Cab){
						hallLights[floor][btn] = hallLights[floor][btn] || isOrder
					}
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

func ReassignHallOrders(wv Worldview, MapAccessChan chan ElevMapAction, orderToElevChan chan<- Order){
	readChan := make(chan any, 1)
	defer close(readChan)
	//Request read
	MapAccessChan <- ElevMapAction{cmd: "read", /* id: lostId, */ respCh: readChan}

	select{
		case elevMap := <-readChan:
			wv = Worldview{wv.PrimaryId,wv.PeerInfo,elevMap.(map[string]Elevator)} //Updated worldview
			for _,lostId := range(wv.PeerInfo.Lost){	
				orderMatrix := elevMap.(map[string]Elevator)[lostId].Orders
				
				for floor, floorOrders := range(orderMatrix){
					for btn, isOrder := range(floorOrders){
						
						if isOrder && btn!=CAB{
							lostOrder:=Order{
										Id: lostId,
										Floor: floor,
										Button: btn,
									}
							lostOrder.Id = assigner.ChooseElevator(wv.ElevatorSnapshot,wv.PeerInfo.Peers,lostOrder)
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
func elevatorMap(mapActionChan <- chan ElevMapAction){
	elevators:=make(map[string]Elevator)
	for{
		select{
		case newAction:= <- mapActionChan:
			switch newAction.cmd{
			case "read":
				newAction.respCh <- elevators[newAction.id]			
			case "write":
				elevators[newAction.id] = newAction.elev
			case "migrate":
				newAction.respCh <- elevators
			case "initialize":
				elevators = newAction.initMap
			}
		}
	}
}