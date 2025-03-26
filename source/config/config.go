package config

import (
	"time"
)

const (
	IDLE ElevatorState = iota
	MOVING
	DOOR_OPEN
)

var NUM_FLOORS int

const (
	NUM_BUTTONS        = 3
	NUM_ELEVATORS      = 3
	NUM_HALL_BTNS      = 2
	DEFAULT_NUM_FLOORS = 4
)

const (
	T_HEARTBEAT        = time.Millisecond * 100
	T_SLEEP            = time.Millisecond * 20
	T_DOOR_OPEN        = time.Second * 3
	T_REASSIGN_PRIMARY = time.Second * 3 // Time before primary clears reassigned orders
	T_REASSIGN_LOCAL   = time.Second * 4 // Time before elev clears reassigned orders
	T_TRAVEL           = time.Second * 2 // Approx time spent travelling between adjacent floors
	T_MOTOR_STOP       = 2 * T_TRAVEL    // Threshold for trigerring motor stop protocol while travelling
	T_PRIMARY_TIMEOUT  = time.Millisecond * 2000
)

const (
	UP   = 1
	DOWN = -1
	STOP = 0
)

const (
	Obstructed = iota
	Disconnected
)

const (
	PORT_BCAST          = 20019
	PORT_PEERS          = 20020
	PORT_WORLDVIEW      = 20040
	DEFAULT_ELEVIO_PORT = "15657"
)

// Technically Dynamic. Cannot pre-allocate due to compile-time limitations of Golang
type OrderMatrix [][]bool //Always [NUM_FLOORS][NUM_BUTTONS]
type HallMatrix [][]bool  //Always [NUM_FLOORS][NUM_BUTTONS-1]

func OrderMatrixConstructor() OrderMatrix {
	output := make(OrderMatrix, NUM_FLOORS)
	for i := range output {
		output[i] = make([]bool, NUM_BUTTONS)
	}
	return output
}
func HallMatrixConstructor() HallMatrix {
	output := make(HallMatrix, NUM_FLOORS)
	for i := range output {
		output[i] = make([]bool, NUM_BUTTONS-1)
	}
	return output
}

type ElevatorState int

type Elevator struct {
	Id            string
	Floor         int
	Direction     int
	PrevDirection int
	State         ElevatorState
	Orders        OrderMatrix
	Requests      OrderMatrix
	Obstructed    bool
}

type Requests struct {
	Id       string
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

type Worldview struct {
	PrimaryId                string
	PeerInfo                 PeerUpdate
	FleetSnapshot            map[string]Elevator
	UnacceptedOrdersSnapshot map[string][]Order
	HallLightsSnapshot       HallMatrix
}

func WorldviewConstructor(PrimaryId string, PeerInfo PeerUpdate, FleetSnapshot map[string]Elevator) Worldview {
	return Worldview{
		PrimaryId:     PrimaryId,
		PeerInfo:      PeerInfo,
		FleetSnapshot: FleetSnapshot,
	}
}

type ElevatorsAccess struct {
	Cmd      string
	Id       string
	Elev     Elevator
	ElevMap  map[string]Elevator
	ReadChan chan map[string]Elevator
}

type OrderAccess struct {
	Cmd              string
	Id               string
	Orders           []Order
	UnacceptedOrders map[string][]Order
	ReadChan         chan map[string][]Order
	ReadAllChan      chan map[string][]Order
}

type LightsAccess struct {
	Cmd           string
	NewHallLights HallMatrix
	ReadChan      chan HallMatrix
}

type Reassignment struct {
	Cause int
	ObsId string
}
