package backup

import (
	"fmt"
	. "source/localElevator/config"
	"source/primary"
	"time"
)

/*
The backup must keep a copy of the primary. It must hold elevators, activepeers, orders etc.
The primary must therefore send this over to backup as a timed heartbeat and event driven.
Maybe we can create one big struct containing everything, that can be sent from prim to back.
*/

func Run(
	/*fromprimary <-chan string, */
	worldview <-chan primary.Worldview, 
	becomePrimary chan <- bool){

	fmt.Println("Enter Backup mode - listening for primary")

	var latestWorldview primary.Worldview

	for {
		select {
		// case msg := <-fromprimary:
		// 	fmt.Println("Received message from primary:", msg)

		case latestWorldview = <- worldview:
			fmt.Println("Worldview received")
			fmt.Printf("Active Peers: %v\n", latestWorldview.ActivePeers)
			fmt.Printf("Elevators: %v\n", latestWorldview.Elevators)
		
		case <-time.After(T_TIMEOUT):
			fmt.Println("Timout waiting for Primary")
			becomePrimary <- true
			//return // Is this necessary, or can it just continue as backup?
				   // While also running the primary protocol?
		}
	}
}