package main

import (
	"fmt"
	. "source/localElevator/config"
	"source/localElevator/elevio"
	//"source/localElevator/fsm"
	"time"
)

func main(){
	fmt.Println("")
	elev:=Elevator{
		Floor: 2,
		Direction: UP,
		State: MOVING,
		Requests: [NUM_FLOORS][NUM_BUTTONS]bool{},
	}

	elevio.Init("localhost:15657", NUM_FLOORS)
	//OrderChan:=make(chan Order)
	AtFloorChan:=make(chan int,5)
	
	//OrderChan<-Order{3, 1}
	go elevio.PollFloorSensor(AtFloorChan)
	for{
		select{
			case elev.Floor = <-AtFloorChan:
				fmt.Printf("Floor: %d\n",elev.Floor)
		}
		time.Sleep(20*time.Millisecond)
	}	
}