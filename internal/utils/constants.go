package utils

const (
	WaitingConnectionsStr     = "WaitingConnections"     // 1 <- Wait for new connections before to start simulation
	RequestingEventsStr       = "RequestingEvents"       // 2 <- Request events to connected clients -
	AwaitingEventResponsesStr = "AwaitingEventResponses" // 3 <-
	CollectingEventsStr       = "CollectingEvents"       // 4 <-
	DispatchingEventsStr      = "DispatchingEvents"      // 5 <- Deliver events to affected clients
	RequestingClientStatusStr = "RequestingClientStatus" // 6 <- Request Client status
	AwaitingClientStatusStr   = "AwaitingClientStatus"   // 7 <-
	StoringClientStatusStr    = "StoringClientStatus"    // 8 <- Store client status
	FinishingStr              = "Finishing"              // 9
	PausedStr                 = "Paused"                 // 0 This state can be used to pause the simulation

	// DB <-- Bbolt
	// ==============================
	BucketClients        = "clients"
	BucketClientToClient = "client-to-client"
	BucketClientResponse = "client-response"
)
