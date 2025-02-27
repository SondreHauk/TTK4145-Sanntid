package backup

import (
	"fmt"
	. "source/localElevator/config"
	"source/primary"
	"time"
)

func Run(
	worldview <-chan primary.Worldview, 
	becomePrimary chan <- bool,
	id string){

	fmt.Println("Enter Backup mode - listening for primary")

	var latestWV primary.Worldview
	//Peers[0] doesnt exist before the first primary does
	select{
		case latestWV = <- worldview:
		case <-time.After(T_TIMEOUT):
			becomePrimary <- true
	}
	
	for {
		select {

		case latestWorldview = <- worldview:
			// fmt.Println("Worldview received")
			// fmt.Printf("Active Peers: %v\n", latestWorldview.ActivePeers)
			// fmt.Printf("Elevators: %v\n", latestWorldview.Elevators)
		
		case <-time.After(T_TIMEOUT):

			if shouldTakeOver(latestWV, id){
				becomePrimary <- true
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