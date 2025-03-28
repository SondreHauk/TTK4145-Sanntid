# Backup module
All elevators have an active backup module. The backup receives worldviews from the primary which is then channeled to the elevator and broadcasted back to the primary. The worldview is broadcasted back to the primary so that if there are more than one active primary online, the primaray will receive worldviews from another primary Id, and one of them will back down.

In the event of a primary disconnection, a `timeout` will trigger. The active peer list functions like a queue determining which elevator will take over as the new primary. After the timeout trigger happens, each running backup will check the first ID in the queue. If it matches its own, it takes over as primary. If not, it prematurely removes the first peer in the queue, exits the loop and waits for updates.

If the removed peer *is* alive and has taken over as primary, the latest worldview will be updated before the next timeout, and the modified worldview is void.

If not, the process will repeat and iterate until a primary is detected or itself becomes primary. This algorithm ensures slaves will compete for primary status.