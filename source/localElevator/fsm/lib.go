package fsm

import (
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/requests"
	"time"
)


func ShouldStop(elev Elevator) bool {
	switch elev.Direction {
	case UP:
		if elev.Floor==NUM_FLOORS-1{
			return true
		}else{
			return elev.Orders[elev.Floor][elevio.BT_HallUp] || 
			elev.Orders[elev.Floor][elevio.BT_Cab] || 
			!requests.OrdersAbove(elev)
		}
	case DOWN:
		if elev.Floor==0{
			return true
		}else{
			return elev.Orders[elev.Floor][elevio.BT_HallDown] || 
			elev.Orders[elev.Floor][elevio.BT_Cab] || 
			!requests.OrdersBelow(elev)
		}
	case STOP:
		return true
	}
	return false
}

func ChooseDirection(elev Elevator) int {
	// In case of orders above and below; choose last moving direction
	if elev.PrevDirection == UP{
		if requests.OrdersAbove(elev) {
			return UP
		} else if requests.OrdersBelow(elev) {
			return DOWN
		}
	} else {
		if requests.OrdersBelow(elev) {
			return DOWN
		} else if requests.OrdersAbove(elev) {
			return UP
		}
	}
	return STOP

}

//Simulates elevator execution and returns approx time until pickup at NewOrder.Floor
func TimeUntilPickup(elev Elevator, NewOrder Order) time.Duration{
	duration := time.Duration(0)
	elev.Orders[NewOrder.Floor][NewOrder.Button]=true
	// Determines initial state
	switch elev.State {
	case IDLE:
		elev.Direction = ChooseDirection(elev)
		if elev.Direction == STOP && elev.Floor == NewOrder.Floor{
			return duration
		}
	case MOVING:
		duration += T_TRAVEL / 2
		elev.Floor += int(elev.Direction)
	case DOOR_OPEN:
		duration -= T_DOOR_OPEN / 2
	}

	for {
		if ShouldStop(elev) {
			if elev.Floor == NewOrder.Floor{
				return duration
			}else{
				for btn:=0; btn<NUM_BUTTONS; btn++{
					elev.Orders[elev.Floor][btn]=false
				}
				duration += T_DOOR_OPEN
				elev.Direction = ChooseDirection(elev)
			}
		}
		elev.Floor += int(elev.Direction)
		duration += T_TRAVEL
	}
}

func checkForNewOrders(
	wv Worldview,
	 myId string, 
	 orderChan chan <- Order, 
	 requestAcceptedChan chan <- HallMatrix,
	 acceptedorders [NUM_FLOORS][NUM_BUTTONS]bool) {
	
	// send acc hall orders to request module 
	accHallOrders := HallMatrix{}
	for _, accOrders := range wv.UnacceptedOrdersSnapshot{
			for _, ord := range accOrders{
				accHallOrders[ord.Floor][ord.Button] = true
			}
		}
	requestAcceptedChan <- accHallOrders

	// send assigned order to elevator
	orders, exists := wv.UnacceptedOrdersSnapshot[myId]
	if exists {
		for _, order := range orders{
			if !acceptedorders[order.Floor][order.Button] {
			orderChan <- order
			}
		}
	}
}

func checkForNewLights(wv Worldview, currenthallLights HallMatrix, hallLightsChan chan HallMatrix) {
	// if any update in hall lights. Send new lights on HallLightsChan
	for i := range currenthallLights {
		for j := range currenthallLights[i] {
			// Indexing empty hallightssnapshot error
			if currenthallLights[i][j] != wv.HallLightsSnapshot[i][j] {
				hallLightsChan <- wv.HallLightsSnapshot
				return
			}
		}
	}
}

// //Make modular with for loop up to NUM_ELEV
// func PrintRequests(elev Elevator){
// 	fmt.Printf("Floor 4: %t %t %t\n",elev.Orders[3][0],elev.Orders[3][1],elev.Orders[3][2])
// 	fmt.Printf("Floor 3: %t %t %t\n",elev.Orders[2][0],elev.Orders[2][1],elev.Orders[2][2])
// 	fmt.Printf("Floor 2: %t %t %t\n",elev.Orders[1][0],elev.Orders[1][1],elev.Orders[1][2])
// 	fmt.Printf("Floor 1: %t %t %t\n\n",elev.Orders[0][0],elev.Orders[0][1],elev.Orders[0][2])
// }

// func PrintState(elev Elevator){
// 	switch elev.State{
// 		case IDLE: fmt.Printf("State: IDLE\n")
// 		case MOVING: fmt.Printf("State: MOVING\n")
// 		case DOOR_OPEN: fmt.Printf("State: DOOR_OPEN\n")
// 	}
// }