package backup

import (
	"fmt"
	. "source/config"
	"source/primary"
	"time"
)

func Run(
	worldViewChan <- chan primary.Worldview, 
	becomePrimaryChan chan <- primary.Worldview,
	id string){

	fmt.Println("Enter Backup mode - listening for primary")
	//Init an empty worldview
	var latestWV primary.Worldview
	latestWV.PrimaryId = id
	latestWV.ElevatorSnapshot = make(map[string]Elevator)	
	//Peers[0] doesnt exist before the first primary does
	select{
		case latestWV = <- worldViewChan:
		case <-time.After(T_TIMEOUT):
			becomePrimaryChan <- latestWV
	}

	for {
		select {
		case latestWV = <-worldViewChan:
		
		case <-time.After(T_TIMEOUT):
			if shouldTakeOver(latestWV, id){
				becomePrimaryChan <- latestWV
				fmt.Println("Primary timeout - Taking over")
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