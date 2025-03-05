package backup

import (
	"fmt"
	. "source/config"
	"time"
)

func Run(
	worldViewChan <-chan Worldview,
	becomePrimaryChan chan<- Worldview,
	id string) {

	fmt.Println("Enter Backup mode - listening for primary")
	//Init an empty worldview
	var latestWV Worldview
	latestWV.PrimaryId = id
	latestWV.FleetSnapshot = make(map[string]Elevator)
	//Peers[0] doesnt exist before the first primary does

	select {
	case latestWV = <-worldViewChan:
	case <-time.After(T_TIMEOUT):
		becomePrimaryChan <- latestWV
	}

	for {
		select {
		case latestWV = <-worldViewChan:
		case <-time.After(T_TIMEOUT):
			if shouldTakeOver(latestWV, id) {
				becomePrimaryChan <- latestWV
				fmt.Println("Primary timeout - Taking over")
			} else {
				latestWV.PeerInfo.Peers = latestWV.PeerInfo.Peers[1:]
			}
		}
	}
}

func shouldTakeOver(backupWorldview Worldview, id string) bool {
	peerIds := backupWorldview.PeerInfo.Peers
	return peerIds[0] == id
}
