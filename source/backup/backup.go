package backup

import (
	. "source/localElevator/config"
	"fmt"
	"time"
)


func Run(fromprimary <-chan string, becomePrimary chan <- bool){
	//Listen for string from primary. If no message received within 2 seconds: Send bool = 1 on chan to activate primary func.
	for {
		select {
		case msg := <-fromprimary:
			fmt.Println("Received message from primary:", msg)
		
		case <-time.After(T_TIMEOUT):
			fmt.Println("Timout waiting for Primary")
			becomePrimary <- true
			return
		}
	}
}


// func MsgBcastRX(port int, msg chan Message) {
// 	go bcast.Receiver(port, msg)
// 	for {
// 		msg_rx := <- msg
// 		fmt.Printf("Message received: ID = %d, Heartbeat = %s\n", msg_rx.ID, msg_rx.Heartbeat)
// 	}
// }

// select{
// case msg_received := <- msg:
// 	fmt.Printf("Message received: ID = %d, Heartbeat = %s\n", msg_received.ID, msg_received.Heartbeat)
// case <- time.After(2 * time.Second):
// 	fmt.Println("No msg received in 2 sec. Still listening")
// }
// Listen for heartbeats from primary
// Update state and acknowlegde to primary
// If heartbeat time out, send primary dead event.