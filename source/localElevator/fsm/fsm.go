package fsm

import (
	// "fmt"
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/requests"
	"time"
)

func Run(
	elev *Elevator, 
	elevChan chan <-Elevator, 
	atFloorChan <-chan int, 
	orderChan chan Order,
	accReqChan chan <- OrderMatrix,
	obstructionChan <-chan bool,
	worldviewToElevatorChan <-chan Worldview,
	myId string) {

	// Define local variables
	var wv Worldview
	hallLights := HallMatrix{}
	hallLightsChan := make(chan HallMatrix, 10)

	// Set timers
	heartbeatTimer := time.NewTimer(T_HEARTBEAT)
	doorTimer := time.NewTimer(T_DOOR_OPEN)
	doorTimer.Stop()
	obstructionTimer := time.NewTimer(T_REASSIGN_LOCAL)
	obstructionTimer.Stop()
	/* motorstopTimer := time.NewTimer(T_REASSIGN_LOCAL)
	motorstopTimer.Stop()
	go motorstopPoll() */

	for {
		select {
		case wv = <- worldviewToElevatorChan:
			// fmt.Println("Worldview received by elevator")
			checkForNewOrders(wv, myId, orderChan, accReqChan, elev.Orders)
			checkForNewLights(wv, hallLights, hallLightsChan)
		case NewOrder := <-orderChan:
			elev.Orders[NewOrder.Floor][NewOrder.Button] = true
			if NewOrder.Button == int(elevio.BT_Cab){
				elevio.SetButtonLamp(elevio.BT_Cab, NewOrder.Floor, true)
			}
			switch elev.State {
			case IDLE:
				elev.Direction = ChooseDirection(*elev)
				elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
				if elev.Direction == STOP {
					elevio.SetDoorOpenLamp(true)
					doorTimer.Reset(T_DOOR_OPEN)
					elevChan <- *elev //AVOID LOOP
					time.Sleep(T_SLEEP)
					elev.Orders[elev.Floor][NewOrder.Button] = false
					if(NewOrder.Button == int(elevio.BT_Cab)){
						elevio.SetButtonLamp(elevio.BT_Cab, NewOrder.Floor, false)
					}
					elev.State = DOOR_OPEN
				} else {
					elev.State = MOVING
				}
			case MOVING: //NOOP
			case DOOR_OPEN:
				if elev.Floor == NewOrder.Floor {
					elevChan <- *elev //AVOID LOOP BY ACKNOWLEDGING ORDER BEFORE CLEARING
					time.Sleep(T_SLEEP)
					elev.Orders[elev.Floor][NewOrder.Button] = false
					elevio.SetButtonLamp(elevio.ButtonType(NewOrder.Button), elev.Floor, false)
					if !elev.Obstructed{
						doorTimer.Reset(T_DOOR_OPEN)
					}
				}
			}
			elevChan <- *elev
		
		case hallLights = <- hallLightsChan:
			setHallLights(hallLights)


		case elev.Floor = <-atFloorChan:
			elevio.SetFloorIndicator(elev.Floor)
			if ShouldStop(*elev) {
				elevio.SetMotorDirection(elevio.MD_Stop)
				requests.ClearOrder(elev, elev.Floor)
				elev.Direction = STOP
				elevio.SetDoorOpenLamp(true)
				doorTimer.Reset(T_DOOR_OPEN)
				elev.State = DOOR_OPEN
			}
			elevChan <- *elev

		case <-doorTimer.C:
			elevio.SetDoorOpenLamp(false)
			elev.Direction = ChooseDirection(*elev)
			if elev.Direction == STOP {
				elev.State = IDLE
			} else {
				elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
				elev.State = MOVING
			}
			elevChan <- *elev
		
		case ObsEvent:= <-obstructionChan:
			if elev.State==DOOR_OPEN{
				switch ObsEvent{
					case true:
						elev.Obstructed = true
						doorTimer.Stop()
						obstructionTimer.Reset(T_REASSIGN_LOCAL)
					case false:
						elev.Obstructed = false
						doorTimer.Reset(T_DOOR_OPEN)
				}
			}
			elevChan <- *elev
		
		case <- obstructionTimer.C:
			//Delete active hall orders
			for floor, floorOrders := range(elev.Orders){
				for btn, orderActive := range(floorOrders){
					if orderActive && btn != int(elevio.BT_Cab) {
						elev.Orders[floor][btn] = false
					}
				}
			}
			obstructionTimer.Stop()

		case <-heartbeatTimer.C:
			elevChan <- *elev
			heartbeatTimer.Reset(T_HEARTBEAT)
		}

		time.Sleep(T_SLEEP)
	}
}