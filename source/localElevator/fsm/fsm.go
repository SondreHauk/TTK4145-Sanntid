package fsm

// This module should contain the finite state machine for the local elevator

import (
	. "source/localElevator/config"
	//"source/localElevator/elevator"
	//"source/localElevator/elevator"
	"fmt"
	"math/rand"
	"source/localElevator/elevio"
	"source/localElevator/lights"
	"source/localElevator/requests"
	"time"
)


func OrdersAbove(elev Elevator) bool {
	for fl := elev.Floor + 1; fl < NUM_FLOORS; fl++ {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Requests[fl][btn] {
				return true
			}
		}
	}
	return false
}

func OrdersBelow(elev Elevator) bool {
	for fl := elev.Floor - 1; fl >= 0; fl-- {
		for btn := 0; btn < NUM_BUTTONS; btn++ {
			if elev.Requests[fl][btn] {
				return true
			}
		}
	}
	return false
}

func ShouldStop(elev Elevator) bool {
	switch elev.Direction {
	case UP:
		if elev.Floor==NUM_FLOORS-1{
			return true
		}else{
			return elev.Requests[elev.Floor][elevio.BT_HallUp] || elev.Requests[elev.Floor][elevio.BT_Cab] || !OrdersAbove(elev)
		}
	case DOWN:
		if elev.Floor==0{
			return true
		}else{
			return elev.Requests[elev.Floor][elevio.BT_HallDown] || elev.Requests[elev.Floor][elevio.BT_Cab] || !OrdersBelow(elev)
		}
	case STOP:
		return true
	default:
		fmt.Println("DEFAULT ERROR STOP")
		return false
	}
}

func ChooseDirection(elev Elevator) int {
	// In case of orders above and below; choose direction at random
	// Not very smaart, but it works
	rand.Seed(time.Now().UnixNano())
	rand := rand.Intn(10)
	if rand % 2 == 0{
		if OrdersAbove(elev) {
			return UP
		} else if OrdersBelow(elev) {
			return DOWN
		}
	} else {
		if OrdersBelow(elev) {
			return DOWN
		} else if OrdersAbove(elev) {
			return UP
		}
	}
	return STOP
}


//TODO: Fix spontainous elevator spasm bug when buttons are pressed at the same time
func Run(elev Elevator, ElevCh chan Elevator, AtFloorCh chan int, NewOrderCh chan Order){
	for{
		select{
			case NewOrder := <-NewOrderCh:
				elev.Requests[NewOrder.Floor][NewOrder.Button] = true
				fmt.Println("New order received")
			
			case elev.Floor = <-AtFloorCh:
				elevio.SetFloorIndicator(elev.Floor)
				if ShouldStop(elev){
					elevio.SetMotorDirection(elevio.MD_Stop)
					lights.OpenDoor()
					requests.ClearFloor(&elev, elev.Floor)
					elev.State = DOOR_OPEN
				}

			default:
				switch elev.State{
					case IDLE:
						elev.Direction = ChooseDirection(elev)
						elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))

					case MOVING:
						// NOOP

					case DOOR_OPEN:
						elev.Direction = ChooseDirection(elev)
						if elev.Direction == STOP {
							elev.State = IDLE
						} else {
							elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
							elev.State = MOVING
						}
			}
		time.Sleep(20 * time.Millisecond)
		}
	}
}


// This is the bug. Direction and requests should only be updated when elevator arrives at a floor.
// if elev.Floor == NewOrder.Floor{
// 	requests.ClearFloor(&elev, elev.Floor)
// 	lights.OpenDoor2()
// 	elev.State = DOOR_OPEN
// } else {
// 	elev.Direction = ChooseDirection(elev)
// 	elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
// 	elev.State = MOVING
// }

// func Run(elev Elevator, ElevCh chan Elevator, AtFloorCh chan int, NewOrderCh chan Order) {
// 	ElevCh <- elev //Update elevator state
// 	DoorTimer := time.NewTimer(3 * time.Second)
// 	DoorTimer.Stop()
// 	fmt.Println("")
// 	for {
// 		select {
// 		case NewOrder := <-NewOrderCh:
// 			/* if NewOrder.Done {
// 				requests.ClearFloor(elev, NewOrder.Floor)
// 			} else { */
// 				elev.Requests[NewOrder.Floor][NewOrder.Button] = true
// 			/* } */
			
// 			switch elev.State {
// 				case IDLE:
// 					elev.Direction = ChooseDirection(elev)
// 					elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
// 					if elev.Direction == STOP {
// 						lights.OpenDoor(DoorTimer)
// 						elev.State = DOOR_OPEN
// 					} else {
// 						elev.State = MOVING
// 					}
// 				case MOVING: //NOOP
// 				case DOOR_OPEN:
// 					if elev.Floor == NewOrder.Floor {
// 						requests.ClearFloor(&elev, elev.Floor)
// 						lights.OpenDoor(DoorTimer)
// 					}
// 			}

// 			ElevCh <- elev //Update elevator state

// 		case elev.Floor = <-AtFloorCh:
// 			elevio.SetFloorIndicator(elev.Floor)
// 			if ShouldStop(elev) {
// 				elevio.SetMotorDirection(elevio.MD_Stop)
// 				requests.ClearFloor(&elev, elev.Floor)
				
// 				lights.OpenDoor(DoorTimer)
// 				elev.State = DOOR_OPEN
// 			}
// 			ElevCh <- elev //Update elevator state

// 		case <-DoorTimer.C:

	// 	elevio.SetDoorOpenLamp(false)
			// elev.Direction = ChooseDirection(elev)
			// if elev.Direction == STOP {
			// 	elev.State = IDLE
			// } else {
			// 	elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
			// 	elev.State = MOVING
			// }
			// ElevCh <- elev //Update elevator state
// 		}
// 		time.Sleep(20*time.Millisecond)
// 	}
// }