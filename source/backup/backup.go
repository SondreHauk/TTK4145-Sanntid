package backup

import (
	"fmt"
	. "source/localElevator/config"
	"source/primary"
	"time"
)

func Run(
	worldViewChan <-chan primary.Worldview, 
	becomePrimaryChan chan <- bool,
	id string){

	fmt.Println("Enter Backup mode - listening for primary")

	var latestWV primary.Worldview
	//Peers[0] doesnt exist before the first primary does
	select{
		case latestWV = <- worldViewChan:
		case <-time.After(T_TIMEOUT):
			becomePrimaryChan <- true
	}

	for {
		select {
		case latestWV = <-worldViewChan:
			// fmt.Println("Worldview received")
			// fmt.Printf("Active Peers: %v\n", latestWorldview.PeerInfo)
			// fmt.Printf("Elevators: %v\n", latestWorldview.Elevators)
		
		case <-time.After(T_TIMEOUT):
			if shouldTakeOver(latestWV, id){
				becomePrimaryChan <- true
			}else{
				latestWV.PeerInfo.Peers = latestWV.PeerInfo.Peers[1:]
			}
		}
	}
}

func shouldTakeOver(backupWorldview primary.Worldview, id string)bool{
	peerIds:=backupWorldview.PeerInfo.Peers
	return peerIds[0]==id
}