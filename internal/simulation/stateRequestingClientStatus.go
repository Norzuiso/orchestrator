package simulation

import (
	"fmt"
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateRequestingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (s *StateRequestingClientStatus) StartState() {
	log.Printf("State: %v", s.GetStateName())
}

func (s *StateRequestingClientStatus) GetStateName() string {
	return utils.RequestingClientStatusStr
}

func (s *StateRequestingClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return fmt.Errorf("Message not allosit. Orchestrator is not reciving any message")
}

func (s *StateRequestingClientStatus) SendMsg(msg *pb.Message, conn *models.Connection) error {
	clientStreams := s.SimulationEngine.ClientService.ClientStreams
	for id, client := range clientStreams {
		s.SimulationEngine.Orchestrator.ClientsRequest = append(s.SimulationEngine.Orchestrator.ClientsRequest, id)
		done := make(chan error, 1)
		client.Outbox <- models.OutboxItem{Msg: &pb.Message{
			SenderId:    0,
			Epoch:       s.SimulationEngine.Orchestrator.Epoch,
			MessageType: pb.MessageType_MESSAGE_TYPE_REQUEST_CLIENT_STATUS,
			Content:     "",
			Seed:        s.SimulationEngine.Orchestrator.GetClientSeed(id),
		}, Done: done}
		<-done
	}
	s.SimulationEngine.NextState()
	return nil
}
func (s *StateRequestingClientStatus) GetNextState() (utils.State, error) {
	return NewStateStoringClientStatus(s.SimulationEngine), nil
}

func (s *StateRequestingClientStatus) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return false // This state does not allow any type of msg
}

func NewStateRequestingClientStatus(s *SimulationEngine) *StateRequestingClientStatus {
	return &StateRequestingClientStatus{SimulationEngine: s}
}
