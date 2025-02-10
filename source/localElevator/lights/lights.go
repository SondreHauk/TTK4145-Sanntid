package lights

import (
	. "source/localElevator/config"
	//source/localElevator/elevator"
	"source/localElevator/elevio"
	"time"
	"fmt"
)

func LightsInit(elev Elevator){
	for fl:=0; fl<NUM_FLOORS; fl++{
		for btn:=0; btn<NUM_BUTTONS; btn++{
			elevio.SetButtonLamp(elevio.ButtonType(btn),fl,false)
		}
	}
}

//Updates cab and floor lights wrt current request matrix.
//Will be called often.
/* func Update(){
	for fl := range NUM_FLOORS{
		for btn := range NUM_BUTTONS{
			elevio.SetButtonLamp(elevio.ButtonType(btn),fl,elev.Requests[fl][btn])
		}
	}
} */
/* (ch chan elevio.ButtonEvent){
	for{
		select {
			case btn := <-ch:
				elevio.SetButtonLamp(btn.Button, btn.Floor, true)
		}
		time.Sleep(20*time.Millisecond)
	}
 */

//Turns on the door light and resets the timer
//Lights go off after 3 seconds
func OpenDoor(timer *time.Timer){
	elevio.SetDoorOpenLamp(true)
	timer.Stop()
	timer.Reset(3*time.Second)
	fmt.Println("Door open")
}

func OpenDoor2(){
	elevio.SetDoorOpenLamp(true)
	time.Sleep(3*time.Second)
	elevio.SetDoorOpenLamp(false)
}