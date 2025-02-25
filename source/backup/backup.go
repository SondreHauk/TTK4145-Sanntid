package backup

import (
	"fmt"
	. "source/localElevator/config"
	"source/primary"
	"time"
)

func Run(
	worldview <-chan primary.Worldview, 
	becomePrimary chan <- bool){

	fmt.Println("Enter Backup mode - listening for primary")

	var latestWorldview primary.Worldview

	for {
		select {

		case latestWorldview = <- worldview:
			// fmt.Println("Worldview received")
			// fmt.Printf("Active Peers: %v\n", latestWorldview.ActivePeers)
			// fmt.Printf("Elevators: %v\n", latestWorldview.Elevators)
		
		case <-time.After(T_TIMEOUT):
			//fmt.Println("Timout waiting for Primary")
			becomePrimary <- true
		}
	}
}