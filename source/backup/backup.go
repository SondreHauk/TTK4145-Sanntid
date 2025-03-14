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
	latestWV.UnacceptedOrdersSnapshot = make(map[string][]Order)

	hallLights := make([][]bool, NUM_ELEVATORS)
	for i := range hallLights {
		hallLights[i] = make([]bool, NUM_BUTTONS-1)
		for j := range hallLights[i] {
			hallLights[i][j] = false 
		}
	}

	latestWV.HallLightsSnapshot = hallLights
	//Peers[0] doesnt exist before the first primary does

	select{
		case latestWV = <- worldViewChan:
		case <-time.After(T_PRIMARY_TIMEOUT):
			becomePrimaryChan <- latestWV
	}

	for {
		select {
		case latestWV = <-worldViewChan:
			// fmt.Println("Worldview received by backup")
			//fmt.Printf("Active Peers: %v\n", latestWV.PeerInfo)
			//fmt.Printf("Lights: %v\n", latestWV.HallLightsSnapshot)
			//fmt.Printf("Unaccepted Orders: %v\n", latestWV.UnacceptedOrdersSnapshot)
		
		case <-time.After(T_PRIMARY_TIMEOUT):
			if shouldTakeOver(latestWV, id){

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
