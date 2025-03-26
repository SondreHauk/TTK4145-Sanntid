package fsm

import (
	"os/exec"
	"os"
	"runtime"
	"log"
	"strings"
	. "source/config"
	"source/localElevator/elevio"
	"source/localElevator/requests"
	"time"
	"fmt"
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
// WHY IN FSM MODULE?
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
	accReqChan chan <- OrderMatrix,
	acceptedorders OrderMatrix) {
	
	// send all assigned orders to request module 
	accOrdersMatrix := OrderMatrix{}
	for _, accOrders := range wv.UnacceptedOrdersSnapshot{
			for _, ord := range accOrders{
				accOrdersMatrix[ord.Floor][ord.Button] = true
			}
		}
	accReqChan <- accOrdersMatrix

	// send ID assigned order to elevator
	orders, exists := wv.UnacceptedOrdersSnapshot[myId]
	if exists {
		for _, order := range orders{
			if !acceptedorders[order.Floor][order.Button] {
			orderChan <- order
			}
		}
	}
}

func checkForNewLights(wv Worldview, lights HallMatrix, lightsChan chan HallMatrix) {
	for floor, buttons := range lights {
		for btn := range buttons {
			if lights[floor][btn] != wv.HallLightsSnapshot[floor][btn] {
				lightsChan <- wv.HallLightsSnapshot
				return
			}
		}
	}
}

func setHallLights(lights HallMatrix){
	for floor, btns := range lights {
		for btn, status := range btns {
			elevio.SetButtonLamp(elevio.ButtonType(btn), floor, status)
		}
	}
}

func spawnProcess() error{
	path, err := os.Executable()
	if err != nil {
		return err
	}
	args := strings.Join(os.Args[1:], " ")
	commandLine := path
	if args != "" {
		commandLine += " " + args
	}

	var cmd *exec.Cmd
	switch runtime.GOOS{
	case "linux":
		cmd = exec.Command("gnome-terminal", "--", "bash", "-c", path+" "+args+"; exec bash")
	case "darwin":
		cmd = exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "%s"`, commandLine))
	case "windows":
		cmd = exec.Command("cmd","/C","start","",commandLine)
	default:
		return fmt.Errorf("unsupported platform: %s. Valid platforms are Linux, Windows or MacOSX", runtime.GOOS)
	}
	
	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}

func restartUponMotorStop(){
	if err:=spawnProcess(); err != nil {
		log.Printf("Failed to restart process: %v", err)
	}
	log.Println("Motor stop detected. Restart")
	os.Exit(1)
}

func resetTimer(timer *time.Timer, duration time.Duration) {
	if !timer.Stop() {
		select {
		case <-timer.C:
		default:
		}
	}
	timer.Reset(duration)
}

// Send accepted order to Primary ten times to avoid loop
func ackOrder(elev *Elevator, elevChan chan <-Elevator){
	for range 10 {
		elevChan <- *elev
	}
	time.Sleep(T_SLEEP)
}