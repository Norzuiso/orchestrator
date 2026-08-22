package sse

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type LogsBroadcaster = Broadcaster[string]
type OpenStreamBroadcaster = Broadcaster[int64]
type MsgBroadcaster = Broadcaster[*pb.Message]
