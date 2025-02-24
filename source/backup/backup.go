package backup

import (
	"source/network/bcast"
	. "source/localElevator/config"
	"fmt"
)


func MsgBcastRX(port int, msg chan Message) {
	go bcast.Receiver(port, msg)
	for {
		msg_rx := <- msg
		fmt.Printf("Message received: ID = %d, Heartbeat = %s\n", msg_rx.ID, msg_rx.Heartbeat)
	}
}

// select{
// case msg_received := <- msg:
// 	fmt.Printf("Message received: ID = %d, Heartbeat = %s\n", msg_received.ID, msg_received.Heartbeat)
// case <- time.After(2 * time.Second):
// 	fmt.Println("No msg received in 2 sec. Still listening")
// }
// Listen for heartbeats from primary
// Update state and acknowlegde to primary
// If heartbeat time out, send primary dead event.