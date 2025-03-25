# Primary Module

The **primary** module is responsible for managing the system's central decision-making process, ensuring coordination between elevators, handling peer updates, and assigning orders.

## Overview
This module takes on the role of the **primary controller**, handling:
- Peer updates and network changes
- Elevator state updates
- Assign orders to available elevators
- Sending system state to other components

## Key Components

### `Worldview` Struct
The `Worldview` struct maintains the system's state, including:
- `PrimaryId`: The ID of the current primary node.
- `PeerInfo`: Information about which peers are active on the network.
- `Elevators`: A map of elevator states, indexed by their ID.

### `Run` Function
The main function that operates the primary controller:
- Listens for new peer updates and maintains an up-to-date list of active elevators.
- Handles elevator state updates.
- Receives requests and assigns orders using the `assigner.ChooseElevator` function.
- Periodically sends the system state (`worldview`) to backups.
- Implements failover handling, allowing another node to take over as primary when needed.

## How It Works
1. The module starts in a **standby** mode, listening for a signal on `becomePrimaryChan`.
2. Once it receives the signal, it takes over as primary and starts processing system updates.
3. It continuously:
   - Updates peer and elevator states
   - Assigns new orders to elevators
   - Sends system state updates at regular intervals
4. If another node with lower ID becomes the primary, it steps down.
