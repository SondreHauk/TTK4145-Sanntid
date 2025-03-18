package config

import (
	"time"
	"fmt"
)

const (
	IDLE ElevatorState = iota
	MOVING
	DOOR_OPEN
)

const (
	NUM_FLOORS    = 4
	NUM_BUTTONS   = 3
	NUM_ELEVATORS = 1 // TODO: User input?
	NUM_HALL_BTNS = 2
)

const (
	T_HEARTBEAT = time.Millisecond*500 //Must be much faster than .5 s
	T_SLEEP = time.Millisecond*20
	T_DOOR_OPEN = time.Second*3
	T_REASSIGN_PRIMARY = time.Second*3
	T_REASSIGN_LOCAL = time.Second*4
	T_TRAVEL = time.Second*2 	//Approximate time to travel from floor i to floor i+-1
	T_PRIMARY_TIMEOUT = time.Millisecond*2000
	T_BLINK = time.Millisecond*100
)

const (
	UP   = 1
	DOWN = -1
	STOP = 0
)

const(
	Obstructed = iota
	Disconnected
)

const (
	PORT_BCAST      = 20019
	PORT_PEERS      = 20020
	// PORT_ELEVSTATE  = 20030
	PORT_WV         = 20040
	// PORT_REQUEST    = 20050
	// PORT_ORDER      = 20060
	// PORT_HALLLIGHTS = 20070
)

type ElevatorState int

type Elevator struct {
	Id            string
	Floor         int
	Direction     int
	PrevDirection int
	State         ElevatorState
	Orders        [NUM_FLOORS][NUM_BUTTONS]bool
 	Obstructed 	  bool
}

type Order struct {
	Id     string
	Floor  int
	Button int
}

func OrderConstructor(Id string, Floor int, Button int) Order {
	return Order{Id: Id, Floor: Floor, Button: Button}
}
type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

type HallLights [NUM_FLOORS][NUM_BUTTONS-1]bool

//----------------PRIMARY/BACKUP--------------------

type Worldview struct {
	PrimaryId     string
	PeerInfo      PeerUpdate
	FleetSnapshot map[string]Elevator
	UnacceptedOrdersSnapshot map[string][]Order
	HallLightsSnapshot HallLights
}

func WorldviewConstructor(PrimaryId string, PeerInfo PeerUpdate, 
	FleetSnapshot map[string]Elevator,/*, UnacceptedOrdersSnapshot map[string][]Order,
	HallLightSnapshot [][]bool*/) Worldview {
	return Worldview{PrimaryId: PrimaryId, PeerInfo: PeerInfo, FleetSnapshot: FleetSnapshot,
	/*UnacceptedOrdersSnapshot: UnacceptedOrdersSnapshot, HallLightsSnapshot: HallLightSnapshot*/}
}

type FleetAccess struct {
	Cmd     string
	Id      string
	Elev    Elevator
	ElevMap map[string]Elevator
	ReadChan  chan map[string]Elevator
}

type OrderAccess struct {
	Cmd		         string
	Id 				 string
	Orders			 []Order
	UnacceptedOrders map[string][]Order
	ReadChan 		 chan map[string][]Order
	ReadAllChan 	 chan map[string][]Order
}

type LightsAccess struct {
	Cmd 	 	  string
	NewHallLights HallLights
	ReadChan      chan HallLights
}

type Reassignment struct {
	Cause int
	ObsId string
}

//--------------------------------------------

func PrintWorldView(wv Worldview) {
	fmt.Println("--- Worldview Snapshot ---")
	fmt.Println("PrimaryId:", wv.PrimaryId)
	fmt.Println("Peers:", wv.PeerInfo.Peers)
	fmt.Println("New Peer:", wv.PeerInfo.New)
	fmt.Println("Lost Peers:", wv.PeerInfo.Lost)
	fmt.Println("Fleet Snapshot:")
	for id, elev := range wv.FleetSnapshot {
		fmt.Printf("  Elevator ID: %s\n", id)
		fmt.Printf("    Floor: %d, Direction: %d, PrevDirection: %d, State: %d\n", 
			elev.Floor, elev.Direction, elev.PrevDirection, elev.State)
		fmt.Printf("    Obstructed: %t\n", elev.Obstructed)
		fmt.Println("    Orders:")
		for i := 0; i < NUM_FLOORS; i++ {
			fmt.Printf("      Floor %d: %v\n", i, elev.Orders[i])
		}
	}
	fmt.Println("Unaccepted Orders Snapshot:")
	for id, orders := range wv.UnacceptedOrdersSnapshot {
		fmt.Printf("  Orders for Elevator %s:\n", id)
		for _, order := range orders {
			fmt.Printf("    Floor: %d, Button: %d\n", order.Floor, order.Button)
		}
	}
	fmt.Println("Hall Lights Snapshot:")
	for i := 0; i < NUM_FLOORS; i++ {
		fmt.Printf("  Floor %d: %v\n", i, wv.HallLightsSnapshot[i])
	}
	fmt.Println("-------------------------")
}