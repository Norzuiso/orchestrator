package orchestrator

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type Orchestrator struct {
	Seed  int64
	Epoch float32

	ActiveClients             map[int64]*pb.Client
	ClientToClientConnections map[int64][]int64
	ClientEventsResponse      map[int64]*pb.Message
	ClientEventsRequest       []int64
}

func NewOrquestrator() *Orchestrator {
	o := &Orchestrator{}
	o.Epoch = 0
	return o
}

func (o *Orchestrator) StartSimualtion() {
	o.Epoch = 0
}
func (o *Orchestrator) PauseSimualtion() {
	o.Epoch = 0
}

func (o *Orchestrator) NextEpoch() {
	o.Epoch = +1
}
