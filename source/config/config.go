package config

import (
	"time"
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
	T_HEARTBEAT = time.Millisecond*100 //Must be much faster than .5 s
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
	PORT_WORLDVIEW  = 20040
)

type OrderMatrix [NUM_FLOORS][NUM_BUTTONS]bool
type HallMatrix [NUM_FLOORS][NUM_BUTTONS-1]bool

type ElevatorState int

type Elevator struct {
	Id            string
	Floor         int
	Direction     int
	PrevDirection int
	State         ElevatorState
	Orders        OrderMatrix
	Requests  	  OrderMatrix
 	Obstructed 	  bool
}

type Requests struct {
	Id 		 string
	Requests OrderMatrix
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

//----------------PRIMARY/BACKUP--------------------

type Worldview struct {
	PrimaryId     string
	PeerInfo      PeerUpdate
	FleetSnapshot map[string]Elevator
	UnacceptedOrdersSnapshot map[string][]Order
	HallLightsSnapshot HallMatrix
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
	NewHallLights HallMatrix
	ReadChan      chan HallMatrix
}

type Reassignment struct {
	Cause int
	ObsId string
}

//--------------------------------------------