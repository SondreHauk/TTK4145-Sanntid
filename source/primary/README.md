# Primary Module  

The `primary` package manages the coordination and synchronization of the elevators. It ensures proper handling of peer updates, elevator status changes, and request assignments.  

## Key Responsibilities  
- **Primary Activation & Synchronization**: Handles primary election and maintains consistency in the elevator network.  
- **Data Access Management**: Controls access to shared data structures for elevators, orders, and lights.  
- **Order Assignment & Reassignment**: Ensures requests are assigned efficiently and handles order redistribution in case of failures.  
- **Heartbeat & Failover Handling**: Periodically updates the system state and detects failures to prevent split-brain scenarios.  

## Main Components  
- `sync` package: Manages shared data access and state updates.  
- `assigner` package: Assigns requests to elevators based on availability.  
- `Run` function: Core loop that processes updates and maintains the primary role.  