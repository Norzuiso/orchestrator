package models

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type Connection struct {
	Stream pb.Broadcast_ClientToClientMessageServer
	Outbox chan *pb.Message
}
