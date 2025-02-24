package backup

import (
	. "source/localElevator/config"
	"fmt"
	"time"
)


func Run(fromprimary <-chan string, becomePrimary chan <- bool){
	fmt.Println("Enter Backup mode - listening for primary")
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