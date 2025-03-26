package misc

import(
	"os/signal"
	"os"
	"source/localElevator/elevio"
	"fmt"
	. "source/config"
	"time"
)

func kill(StopButtonCh <-chan bool) {
	KeyboardInterruptCh := make(chan os.Signal, 1)
	signal.Notify(KeyboardInterruptCh, os.Interrupt)
	select {
	case <-KeyboardInterruptCh:
		fmt.Println("Keyboard interrupt")
	case <-StopButtonCh:
		for i := 0; i < 5; i++ {
			elevio.SetStopLamp(true)
			time.Sleep(time.Millisecond*50)
			elevio.SetStopLamp(false)
			time.Sleep(time.Millisecond*50)
		}
	}
	elevio.SetMotorDirection(elevio.MD_Stop)
	os.Exit(1)
}
