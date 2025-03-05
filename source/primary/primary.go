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
	cmd string //{"read","write one","write all"}
	id string
	elev Elevator
	elevMap map[string]Elevator
	readCh chan map[string]Elevator
}

func Run(
	peerUpdateChan <-chan peers.PeerUpdate,
	elevStateChan <-chan Elevator,
	becomePrimaryChan <-chan Worldview,
	worldviewChan chan <- Worldview,
	requestFromElevChan <- chan Order,
	orderToElevChan chan <- Order,
	hallLightsChan chan <- [][]bool,
	id string){

	MapActionChan := make(chan ElevMapAction, 10)
	ReadMapChan := make(chan map[string]Elevator, 2)
	//updateLights := new(bool)
	
	var worldview Worldview
	worldview.ElevatorSnapshot = make(map[string]Elevator)
	
	//Init hallLights matrix
	hallLights := make([][]bool, NUM_FLOORS)
	for i := range(hallLights){
		hallLights[i] = make([]bool, NUM_BUTTONS - 1)
	}

	//Handling reads and writes from/to ElevatorMap
	go ElevatorMap(MapActionChan)

	select{
	case worldview := <-becomePrimaryChan:
		fmt.Println("Taking over as Primary")
		MapActionChan <- ElevMapAction{cmd:"write all", elevMap: worldview.ElevatorSnapshot}
		
		//drain(elevStateChan) //FIX FLUSHING OF CHANNELS (?)
		HeartbeatTimer := time.NewTicker(T_HEARTBEAT)

		for{
			select{
			case worldview.PeerInfo = <-peerUpdateChan:
				//If elev lost: Reassign lost orders
				printPeers(worldview.PeerInfo)
				lost:=worldview.PeerInfo.Lost
				if len(lost)!=0{
					ReassignHallOrders(worldview, MapActionChan, orderToElevChan)
				}

			case elevUpdate := <-elevStateChan:
				//Request write
				MapActionChan <- ElevMapAction{cmd: "write one",id: elevUpdate.Id, elev: elevUpdate}
				//has a race condition but works fine
				updateHallLights(worldview, hallLights, MapActionChan, hallLightsChan)
			case request := <-requestFromElevChan:
				
				//Request read
				MapActionChan <- ElevMapAction{cmd: "read", readCh: ReadMapChan}
				select{ case worldview.ElevatorSnapshot = <-ReadMapChan: }
				
				AssignedId := assigner.ChooseElevator(worldview.ElevatorSnapshot,
													worldview.PeerInfo.Peers,
													request)
				orderToElevChan <- Order{Id: AssignedId, 
										Floor: request.Floor,
										Button: request.Button}

			case <-HeartbeatTimer.C:
				MapActionChan <- ElevMapAction{cmd: "read", readCh: ReadMapChan}
				select{ 
					case worldview.ElevatorSnapshot = <-ReadMapChan: 
				}
				worldviewChan <- worldview

			case <-becomePrimaryChan: //Needs logic //does it?
				fmt.Println("Another Primary taking over...")
				break
			}
		}
	}
}

func ReassignHallOrders(wv Worldview, MapAccessChan chan ElevMapAction, orderToElevChan chan<- Order){
	readChan := make(chan map[string]Elevator, 1)
	defer close(readChan)
	//Request read
	MapAccessChan <- ElevMapAction{cmd: "read", readCh: readChan}

	select{
	case elevMap := <-readChan:
		// Update with latest snapshot
		wv = Worldview{wv.PrimaryId,wv.PeerInfo,elevMap}
		for _,lostId := range(wv.PeerInfo.Lost){	
			orderMatrix := wv.ElevatorSnapshot[lostId].Orders
			for floor, floorOrders := range(orderMatrix){
				for btn, isOrder := range(floorOrders){
					if isOrder && btn!=CAB{
						lostOrder:= Order{
										Id: lostId,
										Floor: floor,
										Button: btn,
									}
						lostOrder.Id = assigner.ChooseElevator(wv.ElevatorSnapshot, wv.PeerInfo.Peers, lostOrder)
						orderToElevChan <- lostOrder
					}
				}
			}
		}
	}
}

func ElevatorMap(mapActionChan <- chan ElevMapAction){
	elevators:=make(map[string]Elevator)
	for{
		select{
		case newAction:= <- mapActionChan:
			switch newAction.cmd{
			case "read":
				deepCopy := make(map[string]Elevator, len(elevators))
				for key, value := range elevators{
					deepCopy[key] = value
				}
				newAction.readCh <- deepCopy
			case "write one":
				elevators[newAction.id]=newAction.elev
			case "write all":
				elevators = newAction.elevMap
			}
		}
	}
}

/* MAYBE implement function that owns hallLight state to avoid "trivial" race condition. Would be similar to ElevatorMap
	NOT 1st priority.  */

func updateHallLights(wv Worldview, hallLights [][]bool, MapActionChan chan<- ElevMapAction, hallLightsChan chan<-[][]bool){
	readChan := make(chan map[string]Elevator, 1)
	defer close(readChan)
	//Request read
	MapActionChan <- ElevMapAction{cmd: "read", readCh: readChan}
	shouldUpdate:=false
	prevHallLights := make([][]bool, NUM_FLOORS)
	for floor := range hallLights {
		prevHallLights[floor] = make([]bool, NUM_BUTTONS-1)
		copy(prevHallLights[floor], hallLights[floor]) // Copy row data
		for btn := range(NUM_BUTTONS-1){
			hallLights[floor][btn]=false
		}
	}
	select{
	case elevMap:= <-readChan:	
		wv=Worldview{wv.PrimaryId, wv.PeerInfo, elevMap}
		for _, id := range(wv.PeerInfo.Peers){
			orderMatrix := wv.ElevatorSnapshot[id].Orders
			for floor, floorOrders := range(orderMatrix){
				for btn, isOrder := range(floorOrders){
					if isOrder && btn!= int(elevio.BT_Cab){
						hallLights[floor][btn] = true
					}
				}
			}
		}
	}
	for floor := range(NUM_FLOORS){
		for btn := range(NUM_BUTTONS-1){
			if prevHallLights[floor][btn] != hallLights[floor][btn]{
			shouldUpdate = true
			}
		}
	}
	if shouldUpdate{
		hallLightsChan<-hallLights
	}
}

func printPeers(p peers.PeerUpdate){
	fmt.Printf("Peer update:\n")
	fmt.Printf("  Peers:    %q\n", p.Peers)
	fmt.Printf("  New:      %q\n", p.New)
	fmt.Printf("  Lost:     %q\n", p.Lost)
} 

/* func drain(ch <- chan Elevator){
	for len(ch) > 0{
		<- ch
	}
}

func printElevator(e Elevator){
	fmt.Println("Elevator State Updated")
	fmt.Printf("ID: %s\n", e.Id)
	fmt.Printf("Floor: %d\n", e.Floor)
}
*/