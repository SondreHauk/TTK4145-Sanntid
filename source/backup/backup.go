package backup

import (
	"fmt"
	. "source/config"
	"time"
)

func Run(
	worldviewRXChan <-chan Worldview,
	worldViewToElevChan chan <- Worldview,
	becomePrimaryChan chan<- Worldview,
	worldviewToPrimaryChan chan Worldview,
	myId string) {

	fmt.Println("Enter Backup mode - listening for primary")
	
	var latestWV Worldview
	latestWV.PrimaryId = myId
	latestWV.FleetSnapshot = make(map[string]Elevator)
	latestWV.UnacceptedOrdersSnapshot = make(map[string][]Order)

	select{ //INIT
	case latestWV = <- worldviewRXChan:

	case <-time.After(T_PRIMARY_TIMEOUT):
		becomePrimaryChan <- latestWV
	}

	for {
		select {
		case latestWV = <-worldviewRXChan:
			worldViewToElevChan <- latestWV
			worldviewToPrimaryChan <- latestWV
		
		case <-time.After(T_PRIMARY_TIMEOUT):
			if shouldTakeOver(latestWV, myId){
				latestWV.PrimaryId = myId
				becomePrimaryChan <- latestWV
				fmt.Println("Primary timeout - start takeover")
			} else {
				latestWV.PeerInfo.Peers = latestWV.PeerInfo.Peers[1:]
			}
		}
	}
}

func shouldTakeOver(backupWorldview Worldview, id string) bool {
	peerIds := backupWorldview.PeerInfo.Peers
	if len(peerIds) == 0 {
		return true
	} else {
		return peerIds[0] == id
	}
}