package primary

import (
	"source/network/bcast"
	. "source/localElevator/config"
	"time"
)

func MsgTX(port int, msg chan Message){
	for{
		msg <- Message{ID: 1, Heartbeat: "Alive"}
		bcast.Transmitter(port, msg)
		time.Sleep(T_HEARTBEAT)
	}
}