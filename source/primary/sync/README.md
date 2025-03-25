## Syncronization module
As the name implies, this module ensures thread-safe access to shared variables in the primary. It does this by utilizing channels and blocking select statements.

In order for a routine to read or write from one of the shared variables, it must first send a r/w request on the corresponding access channel. The AccessManagers receives the request and utilizes blocking selects such that only one request is handled at a time. This ensures that only one routine has access to the variable at a time. 

Put simply; the channels work as requsest queues, and the AccessManager processes one request at a time.