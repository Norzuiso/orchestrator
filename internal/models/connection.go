package models

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

type Connection struct {
	Stream pb.Broadcast_ClientToClientMessageServer
	Outbox chan OutboxItem
	ErrCh  chan error
}

type OutboxItem struct {
	Msg  *pb.Message
	Done chan error
}

type ClientToClientConnection struct {
	Client     int64
	Attributes map[string]*anypb.Any
}
