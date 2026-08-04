package simulation

import (
	"fmt"
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateRequestingEvents struct {
	SimulationEngine *SimulationEngine
}

func (s *StateRequestingEvents) StartState() {

	if s.SimulationEngine.Orchestrator.Epoch == s.SimulationEngine.Orchestrator.MaxOfEpochs {
		s.SimulationEngine.CurrentState = s.SimulationEngine.StateEnd
		s.SimulationEngine.CurrentState.StartState()
		return
	}

	s.SimulationEngine.Orchestrator.NextEpoch()
	log.Printf("State: %v", s.GetStateName())
	_ = s.SendMsg(nil, nil)
}

func (s *StateRequestingEvents) GetStateName() string {
	return utils.RequestingEventsStr
}

func (s *StateRequestingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return fmt.Errorf("Message not allosit. Orchestrator is not reciving any message")
}

func (s *StateRequestingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error {
	clientStreams := s.SimulationEngine.ClientService.ClientStreams
	for id, client := range clientStreams {
		s.SimulationEngine.Orchestrator.AppendClientsRequest(id)
		done := make(chan error, 1)
		client.Outbox <- models.OutboxItem{Msg: &pb.Message{
			SenderId:    0,
			Epoch:       s.SimulationEngine.Orchestrator.GetEpoch(),
			MessageType: pb.MessageType_MESSAGE_TYPE_REQUEST_EVENT,
			Content:     "",
			Seed:        s.SimulationEngine.Orchestrator.GetClientSeed(id),
		}, Done: done}
		<-done

	}
	s.SimulationEngine.NextState()
	return nil
}

func (s *StateRequestingEvents) GetNextState() (utils.State, error) {
	return NewStateCollectingEvents(s.SimulationEngine), nil
}

func (s *StateRequestingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return false // This state does not allow any type of msg
}

func NewStateRequestingEvents(s *SimulationEngine) *StateRequestingEvents {
	return &StateRequestingEvents{SimulationEngine: s}
}
