package main

import (
	"github.com/Norzuiso/orchestrator/internal/simulation"
)

func main() {
	simulation.StartSimulationEnine()
}

/*
TODO
==================================================
1. Http function to start simulation
1.1. Check data need it to start simulation
2. Http function to finish simulation
2.1. Check values need to return

3. Http function to send content msg to all client
4. Http function to send content msg to specific client
4.1. Http function to send content msg to specific list of clients

5. Http function to change orchestrator phase. Why? The flow with state pattern works fine.
5.1. Change phase will go to the next phase. But, should we add a function to change into an specific function? Naaa
5.2. Possibility to create phases to allow custom behavior. No need it


10. Orchestrator function to do what each phase its supose to be doing
10.0. Implement state pattern to orchestrator
10.1. WaitConnections
10.2. RequestEvents
10.3. RecolectEvent
10.4. ApplyEvent
10.5. RequestClientStatus
10.6. StoringClientStatus
10.7. Pause

13. Http function to next epoch
14. Http function to next phase
==================================================

6. Http function to get System status
6.1. What data is need it to return?
6.2. is this going to call every client or just check the latest record? Latest record sounds good

============= DEPENDENCY ON TODO #11 =============
7. Http function to get all active clients
8. Http function to get all client to client connection
9. Http function to get client to client connection by ClientId
==================================================

11. Implement an persistent way to store the connections and clients
11.1. Store client id, client address, name, description, seed
11.2. Store client status that needs to be any since each client has individual business logic
11.2.1. Organized by client, epoch. Also, the status has to be storing epoc seed.

12. EMPTY

15. Create function to bump seed by epoch

16. Check english gramar
16.1. Check responses
16.2. Check error msgs
16.3. Check packages, files, functions, var names

17. Add msg response to open stream message type

==================================================
Errors/BUGS:
==================================================
1. One open stream can open multiple streams with different clients ids.

==================================================
NOTES:
==================================================
- The seed value could break the idea of be an agnostic simulation engine, since it's used to generate values on client side in the students example

*/
