package backup

import (
	"fmt"
	. "source/config"
	"time"
)

func Run(
	worldViewToBackupChan <-chan Worldview,
	worldViewToElevChan chan <- Worldview,
	becomePrimaryChan chan<- Worldview,
	myId string) {

	fmt.Println("Enter Backup mode - listening for primary")
	//Init an empty worldview
	var latestWV Worldview
	hallLights := HallLights{}
	latestWV.PrimaryId = myId
	latestWV.FleetSnapshot = make(map[string]Elevator)
	latestWV.UnacceptedOrdersSnapshot = make(map[string][]Order)
	latestWV.HallLightsSnapshot = hallLights
	//Peers[0] doesnt exist before the first primary does

	// go func() {
	// 	for wv := range worldViewToBackupChan {
	// 		fmt.Printf("Backup: Received worldview update from primary %s\n", wv.PrimaryId)
	// 		latestWV = wv
	// 	}
	// }()

	select{
	case latestWV = <- worldViewToBackupChan:
		fmt.Printf ("Wordview prio received by primary: %s\n", latestWV.PrimaryId)
	case <-time.After(T_PRIMARY_TIMEOUT):
		becomePrimaryChan <- latestWV
	}

	for {
		select {
		case latestWV = <-worldViewToBackupChan:
			worldViewToElevChan <- latestWV
			fmt.Printf("Worldview post received by primary: %s\n", latestWV.PrimaryId)
		
		case <-time.After(T_PRIMARY_TIMEOUT):
			if shouldTakeOver(latestWV, myId){
				latestWV.PrimaryId = myId
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
	if len(peerIds) == 0 {
		return true
	} else {
		return peerIds[0] == id
	}
}
