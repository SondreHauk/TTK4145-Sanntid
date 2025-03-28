package fsm

import (
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/requests"
	"source/localinit"
	"time"
)

func Run(
	elevTXChan chan<- Elevator,
	requestsTXChan chan Requests,
	worldviewToElevatorChan <-chan Worldview,
	stopChan chan bool,
	myId string,
	port string,
) {
	// Local variables
	var wv Worldview
	var NewOrder Order
	hallLights := HallMatrixConstructor()

	// Initializations
	elevio.Init("localhost:"+port, NUM_FLOORS)
	elev := localinit.ElevatorInit(myId)
	localinit.LightsInit()
	heartbeatTimer := time.NewTimer(T_HEARTBEAT)
	doorTimer := time.NewTimer(T_DOOR_OPEN)
	doorTimer.Stop()
	obstructionTimer := time.NewTimer(T_REASSIGN_LOCAL)
	obstructionTimer.Stop()
	motorstopTimer := time.NewTimer(T_MOTOR_STOP)
	motorstopTimer.Stop()

	// Local channels
	atFloorChan := make(chan int, 1)
	buttonChan := make(chan elevio.ButtonEvent, 10)
	obstructionChan := make(chan bool, 1)
	acceptedRequestsChan := make(chan OrderMatrix, 10)
	hallLightsChan := make(chan HallMatrix, 10)
	orderChan := make(chan Order, 10)

	// Local goroutines
	go elevio.PollButtons(buttonChan)
	go elevio.PollFloorSensor(atFloorChan)
	go elevio.PollObstructionSwitch(obstructionChan)
	go elevio.PollStopButton(stopChan)
	go requests.SendRequest(
		buttonChan,
		requestsTXChan,
		acceptedRequestsChan,
		orderChan,
		myId,
	)

	for {
		select {
		case wv = <-worldviewToElevatorChan:
			checkForNewOrders(
				wv,
				myId,
				orderChan,
				acceptedRequestsChan,
				elev.Orders,
			)
			checkForNewLights(
				wv,
				hallLights,
				hallLightsChan,
			)

		case NewOrder = <-orderChan:
			elev.Orders[NewOrder.Floor][NewOrder.Button] = true
			if NewOrder.Button == int(elevio.BT_Cab) {
				elevio.SetButtonLamp(
					elevio.BT_Cab,
					NewOrder.Floor,
					true,
				)
			}
			switch elev.State {
			case IDLE:
				elev.Direction = chooseDirection(elev)
				elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
				resetTimer(motorstopTimer, T_MOTOR_STOP)
				if elev.Direction == STOP {
					motorstopTimer.Stop()
					elevio.SetDoorOpenLamp(true)
					resetTimer(doorTimer, T_DOOR_OPEN)
					ackOrder(elev, elevTXChan) //Acknowledge order before clearing
					elev.Orders[elev.Floor][NewOrder.Button] = false
					if NewOrder.Button == int(elevio.BT_Cab) {
						elevio.SetButtonLamp(
							elevio.BT_Cab,
							NewOrder.Floor,
							false,
						)
					}
					elev.State = DOOR_OPEN
				} else {
					elev.State = MOVING
				}

			case MOVING: //NOOP

			case DOOR_OPEN:
				if elev.Floor == NewOrder.Floor {
					ackOrder(elev, elevTXChan) //Acknowledge order before clearing
					elev.Orders[elev.Floor][NewOrder.Button] = false
					elevio.SetButtonLamp(
						elevio.ButtonType(NewOrder.Button),
						elev.Floor,
						false,
					)
					if !elev.Obstructed {
						resetTimer(doorTimer, T_DOOR_OPEN)
					}
				}
			}
			elevTXChan <- elev

		case hallLights = <-hallLightsChan:
			setHallLights(hallLights)

		case elev.Floor = <-atFloorChan:
			if elev.MotorStop {
				elev.MotorStop = false
			}
			resetTimer(motorstopTimer, T_MOTOR_STOP)
			elevio.SetFloorIndicator(elev.Floor)
			if shouldStop(elev) {
				motorstopTimer.Stop()
				elevio.SetMotorDirection(elevio.MD_Stop)
				requests.ClearOrder(elev, elev.Floor)
				elev.Direction = STOP
				elevio.SetDoorOpenLamp(true)
				resetTimer(doorTimer, T_DOOR_OPEN)
				elev.State = DOOR_OPEN
			}
			elevTXChan <- elev

		case <-doorTimer.C:
			if !elev.Obstructed {
				elevio.SetDoorOpenLamp(false)
				elev.Direction = chooseDirection(elev)
				if elev.Direction == STOP {
					elev.State = IDLE
				} else {
					elevio.SetMotorDirection(elevio.MotorDirection(elev.Direction))
					resetTimer(motorstopTimer, T_MOTOR_STOP)
					elev.State = MOVING
				}
			} else {
				resetTimer(doorTimer, T_DOOR_OPEN)
			}
			elevTXChan <- elev

		case elev.Obstructed = <-obstructionChan:
			if !elev.Obstructed {
				resetTimer(doorTimer, T_DOOR_OPEN)
			} else {
				resetTimer(obstructionTimer, T_REASSIGN_LOCAL)
			}
			elevTXChan <- elev

		case <-obstructionTimer.C:
			for floor, floorOrders := range elev.Orders {
				for btn, orderActive := range floorOrders {
					if orderActive && btn != int(elevio.BT_Cab) {
						elev.Orders[floor][btn] = false
					}
				}
			}

		case <-heartbeatTimer.C:
			elevTXChan <- elev
			resetTimer(heartbeatTimer, T_HEARTBEAT)

		case <-motorstopTimer.C:
			elev.MotorStop = true
			ackOrder(elev, elevTXChan)
			for floor, floorOrders := range elev.Orders {
				for btn, orderActive := range floorOrders {
					if orderActive && btn != int(elevio.BT_Cab) {
						elev.Orders[floor][btn] = false
					}
				}
			}
		}
	}
}
