# TTK4145-Sanntid: The Elevators are Evolving!

## Main approach - Primary Backup
The problem: controll N elevators working in parallel across M floors.

The approach: Primary Backup system. All elevators on the network are backups, while *only one is primary*.The primary handles and distributes all hall requests and lights. The backups receives worldviews from the primary and stores the latest worldview received. If the primary disconnects, one of the backups will take over as primary. 

## Specifications
- No calls are lost.
- The lights and buttons function as expected.
- The doors function as expected.
- Each individual elevator behaves sensibly and efficiently.
- Multiple elevators are more efficient than a single one

## Initialization of elevator
An elevator can be initialized from the command line with: `go run main.go -port="..." -id="..."`.  
Each elevator must be assigned an unique id at initialization.

# The Button Light Contract
## Requests and Orders
In general, when a button is pressed in any elevator, a corresponding `request` is created. This request is then handled and an `order` is made. Each order is marked with an `id` and the order is accepted only by the elevator with the corresponding `id`.

If the request is of type `cab`, it is assigned directly as an order to the elevator:
`elevator -- btnevent --> makeRequest -- cab order --> elevator`

If the request is of type `hall`, it is sent to the primary, who creates an order and assigns it to the most suitable elevator on the network:
`elevator -- btnevent --> makeRequest -- hall request --> primary -- hall order --> elevator`

## When is the light set?
The handling of the `cab lights` is done locally on the elevator. If an elevator recevies a cab order it updates its `order matrix` and sets the corresponding cab light. Likewise, if it completes a cab order, it updates its order matrix and turns off the cab light.

The `hall lights` is handled by the primary. The primary knows that an `order is accepted` when the assigned elevator returns an `elevator state` with the corresponding order set active, i.e with an updated order matrix. With this in mind, the primary uses the order matrices from the elevators to update the hall lights. It does this in a `hall light matrix`, which is essentially a union of all the order matrices. The hall light matrix is then broadcasted to the elevators, who updates their corresponding hall lights.

`primary -- order --> elevator -- order matrix --> primary -- hall light matrix --> elevator`

# Improvements
## Allow for elevators to remove active orders in case of an obstruction.
As per now, it is not possible to remove an assigned order from an elevator. In the case of an elevator disconnect, the order can be reassigned to the other active elevator, but the disconnected elevator will still have the order as active. This causes redundancy and is necessary in the case of a disconnect. 

However, in the case of an obstruction, this redundancy is unecessary and should be improved. When an elevator is obstructed, the orders are redestributed to all non-obstructed elevators after some time, but since the elevator is still online, it is unecessary to take the orders twise. To fix this problem however, demands some deep structural changes in the program.

## One port for all message types
As per now, each message type is sendt through an unique port. Fix this with a message handler. Assign all messages an ID, send them through the same port, and let the message handler assign the message to the right channel.

When the primary assigns an order to an elevator, it starts a countdown timer. If the primary does not receive a correct stateUpdate from the assgined elevator within the deadline: declear the elevator for dead/broken and reassign all hall orders!

# Miscellaneous
## Variable naming convention
As a general rule, camel-case convention is used. Initial capitalization is handled as such: lowercase for localVariables and helperFunctions, uppercase for ExportedVariables and FunctionsUsedElsewhere. Global constants are fully capitalized with underscores between words: PORT_WORLDVIEW, NUM_FLOORS, etc.

