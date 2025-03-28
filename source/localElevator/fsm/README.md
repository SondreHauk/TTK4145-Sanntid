# FSM Module  

The `fsm` (Finite State Machine) module controls the local elevator behavior, handling movement, door operations, and order management.  

## Key Responsibilities  
- **Elevator State Management**: Handles transitions between `IDLE`, `MOVING`, and `DOOR_OPEN` states.  
- **Order Handling**: Processes new orders, updates internal order lists, and acknowledges completed requests.  
- **Timer & Obstruction Handling**: Manages door timers, obstruction detection, and motor stop conditions.  
- **Communication**: Sends elevator status updates and processes incoming requests.