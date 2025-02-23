package backup

import (
	"source/network/bcast"
	. "source/localElevator/config"
)


func MsgRX(port int, msg chan Message) {
			bcast.Receiver(port, msg)
}
// Listen for heartbeats from primary
// Update state and acknowlegde to primary
// If heartbeat time out, send primary dead event.