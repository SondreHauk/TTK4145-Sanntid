package backup

import (
	"fmt"
	. "source/config"
	"time"
)

func Run(
	wvRXChan <-chan Worldview,
	wvToElevChan chan<- Worldview,
	enablePrimaryChan chan<- Worldview,
	wvToPrimaryChan chan Worldview,
	myId string,
) {

	fmt.Println("Enter Backup mode - listening for primary")

	var latestWv Worldview
	latestWv.PrimaryId = myId
	latestWv.FleetSnapshot = make(map[string]Elevator)
	latestWv.UnacceptedOrdersSnapshot = make(map[string][]Order)

	select {
	case latestWv = <-wvRXChan:
	case <-time.After(T_PRIMARY_TIMEOUT):
		enablePrimaryChan <- latestWv
	}

	for {
		select {
		case latestWv = <-wvRXChan:
			wvToElevChan <- latestWv
			wvToPrimaryChan <- latestWv

		case <-time.After(T_PRIMARY_TIMEOUT):
			if shouldTakeOver(latestWv, myId) {
				latestWv.PrimaryId = myId
				enablePrimaryChan <- latestWv
				fmt.Println("Primary timeout - start takeover")
			} else {
				latestWv.PeerInfo.Peers = latestWv.PeerInfo.Peers[1:]
			}
		}
	}
}

func shouldTakeOver(backupWv Worldview, id string) bool {
	peerIds := backupWv.PeerInfo.Peers
	if len(peerIds) == 0 {
		return true
	} else {
		return peerIds[0] == id
	}
}
