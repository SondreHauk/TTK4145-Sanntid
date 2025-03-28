# Primary Module  

The `primary` manages the coordination of all online elevators. Only one primary is active at a time. 

The primary takes all the decisions for the system. It does this based on the received elevator states and requests, which is broadcasted from the elevators. It then processes this information, makes a decision, and broadcasts this in a worldview - which among other things contains hall light status and assigned orders for the elevators.

In the case of an elevator disconnect, an obstructed elevator or a motorstop, the primary `reassigns and remebers` the elevators hall orders and cab orders respectively. This reassign-and-rememeber operation is vital to ensure that no calls are lost.