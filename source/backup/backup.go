package backup

import (
	. "source/localElevator/config"
	"fmt"
	"time"
)

/*
The backup must keep a copy of the primary. It must hold elevators, activepeers, orders etc.
The primary must therefore send this over to backup as a timed heartbeat and event driven.
Maybe we can create one big struct containing everything, that can be sent from prim to back.
*/

func Run(fromprimary <-chan string, becomePrimary chan <- bool){
	fmt.Println("Enter Backup mode - listening for primary")
	for {
		select {
		case msg := <-fromprimary:
			fmt.Println("Received message from primary:", msg)
		
		case <-time.After(T_TIMEOUT):
			fmt.Println("Timout waiting for Primary")
			becomePrimary <- true
			return
		}
	}
}