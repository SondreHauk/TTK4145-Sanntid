package requests

import "localElevator/elevator"

type DirnBehaviourPair struct {
    Dirn     Direction
    Behavior Behavior
}

func RequestsAbove(e Elevator) bool {
    for f := e.floor + 1; f < N_FLOORS; f++ {
        for btn = 0; btn < N_BUTTONS; btn++{
            if(e.Requests[f][btn]){
                return true;
            }
        }
    }
    return false;
}

func RequestsBelow(e Elevator) bool {
    for f = 0; f < e.floor; f++{
        for btn = 0; btn < N_BUTTONS; btn++{
            if(e.Requests[f][btn]){
                return true;
            }
        }
    }
    return false;
}

func RequestsHere(e Elevator) bool{
    for btn = 0; btn < N_BUTTONS; btn++{
        if(e.Requests[e.floor][btn]){
            return true;
        }
    }
    return false;
}


func requests_chooseDirection(e Elevator) DirnBehaviourPair {
    switch(e.Direction){
	case UP:
		if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{D_Down, EB_DoorOpen}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		}
	case DOWN:
		if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		} else if requestsHere(e) {
			return DirnBehaviourPair{D_Up, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		}
	case STOP:
		if requestsHere(e) {
			return DirnBehaviourPair{D_Stop, EB_DoorOpen}
		} else if requestsAbove(e) {
			return DirnBehaviourPair{D_Up, EB_Moving}
		} else if requestsBelow(e) {
			return DirnBehaviourPair{D_Down, EB_Moving}
		}
    default:
        return DirnBehaviourPair{D_Stop, EB_Idle}
    }
}

func requests_shouldStop(e Elevator) bool {
    switch(e.Direction){
    case D_Down:
        return
            e.Requests[e.Floor][B_HallDown] ||
            e.Requests[e.Floor][B_Cab]      ||
            !requests_below(e);
    case D_Up:
        return
            e.Requests[e.Floor][B_HallUp]   ||
            e.Requests[e.Floor][B_Cab]      ||
            !requests_above(e);
    case D_Stop:
    default:
        return 1;
    }

    switch(e.Directio) {
	case D_Down:
		return e.requests[e.Floor][B_HallDown] || e.requests[e.Floor][B_Cab] || !requestsBelow(e)
	case D_Up:
		return e.requests[e.Floor][B_HallUp] || e.requests[e.Floor][B_Cab] || !requestsAbove(e)
	case D_Stop:
		return true
	}
	return false
}
}

int requests_shouldClearImmediately(Elevator e, int btn_floor, Button btn_type){
    switch(e.config.clearRequestVariant){
    case CV_All:
        return e.floor == btn_floor;
    case CV_InDirn:
        return 
            e.floor == btn_floor && 
            (
                (e.dirn == D_Up   && btn_type == B_HallUp)    ||
                (e.dirn == D_Down && btn_type == B_HallDown)  ||
                e.dirn == D_Stop ||
                btn_type == B_Cab
            );  
    default:
        return 0;
    }
}

Elevator requests_clearAtCurrentFloor(Elevator e){
        
    switch(e.config.clearRequestVariant){
    case CV_All:
        for(Button btn = 0; btn < N_BUTTONS; btn++){
            e.requests[e.floor][btn] = 0;
        }
        break;
        
    case CV_InDirn:
        e.requests[e.floor][B_Cab] = 0;
        switch(e.dirn){
        case D_Up:
            if(!requests_above(e) && !e.requests[e.floor][B_HallUp]){
                e.requests[e.floor][B_HallDown] = 0;
            }
            e.requests[e.floor][B_HallUp] = 0;
            break;
            
        case D_Down:
            if(!requests_below(e) && !e.requests[e.floor][B_HallDown]){
                e.requests[e.floor][B_HallUp] = 0;
            }
            e.requests[e.floor][B_HallDown] = 0;
            break;
            
        case D_Stop:
        default:
            e.requests[e.floor][B_HallUp] = 0;
            e.requests[e.floor][B_HallDown] = 0;
            break;
        }
        break;
        
    default:
        break;
    }
    
    return e;
}