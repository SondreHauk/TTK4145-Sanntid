package primary

import (
	"source/network/bcast"
	. "source/localElevator/config"
	"time"
)

func MsgTX(port int, msg chan Message){
	go func(){
		for{
			msg <- Message{ID: 1, Heartbeat: "Alive"}
			bcast.TransmitterModified(port, msg)
			time.Sleep(T_HEARTBEAT)
		}
	}()
}