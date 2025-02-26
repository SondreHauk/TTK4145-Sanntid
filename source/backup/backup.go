package backup

import (
	"fmt"
	. "source/localElevator/config"
	"source/primary"
	"time"
)

func Run(
	worldViewChan <-chan primary.Worldview, 
	becomePrimaryChan chan <- bool){

	fmt.Println("Enter Backup mode - listening for primary")

	//var latestWorldview primary.Worldview

	for {
		select {
		case /*latestWorldview =*/ <- worldViewChan:
			// fmt.Println("Worldview received")
			// fmt.Printf("Active Peers: %v\n", latestWorldview.PeerInfo)
			// fmt.Printf("Elevators: %v\n", latestWorldview.Elevators)
		
		case <-time.After(T_TIMEOUT):
			//fmt.Println("Timout waiting for Primary")
			becomePrimaryChan <- true
			//Send latestWorldView to new primary
		}
	}
}