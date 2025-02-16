package timer
import "time"
// Timer logic.
func DoorTimerReset(){
	a:=time.NewTimer(3*time.Second)
	select {
	case <-a.C:
		return
	}
}