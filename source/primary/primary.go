package primary

import (
	//"source/network/bcast"
	. "source/localElevator/config"
	"time"
)

func MsgTX(port int, msg chan Message, id int){
	//go bcast.Transmitter(port, msg) // Start broadcasting in a separate goroutine
	for {
		msg <- Message{ID: id, Heartbeat: "Alive"}
		time.Sleep(T_HEARTBEAT)
	}
}