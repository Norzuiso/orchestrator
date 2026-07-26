package utils

const (
	WaitingConnectionsStr     = "WaitingConnections"     // 1
	RequestingEventsStr       = "RequestingEvents"       // 2
	AwaitingEventResponsesStr = "AwaitingEventResponses" // 3
	CollectingEventsStr       = "CollectingEvents"       // 4
	DispatchingEventsStr      = "DispatchingEvents"      // 5
	RequestingClientStatusStr = "RequestingClientStatus" // 6
	AwaitingClientStatusStr   = "AwaitingClientStatus"   // 7
	StoringClientStatusStr    = "StoringClientStatus"    // 8
	FinishingStr              = "Finishing"              // 9
	PausedStr                 = "Paused"                 // 0 This state can be used to pause the simulation
)
