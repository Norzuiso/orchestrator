package simulation

import (
	"fmt"
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateDispatchingEvents struct {
	SimulationEngine *SimulationEngine
}

func (s *StateDispatchingEvents) StartState() {
	log.Printf("State: %v", s.GetStateName())
	_ = s.SendMsg(nil, nil)
}

func (s *StateDispatchingEvents) GetStateName() string {
	return utils.DispatchingEventsStr
}

func (s *StateDispatchingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return fmt.Errorf("Message not allow it. Orchestrator is not reciving any message")
}

func (s *StateDispatchingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error {
	responses := s.SimulationEngine.Orchestrator.ClientsResponse
	for id, msgRes := range responses {
		connections := s.SimulationEngine.Orchestrator.ClientToClientConnections[id]
		for _, id := range connections {
			msgRes.SenderId = 0
			msgRes.MessageType = pb.MessageType_MESSAGE_TYPE_EVENT_DISPATCH
			msgRes.Epoch = s.SimulationEngine.Orchestrator.Epoch
			s.SimulationEngine.ClientService.ClientStreams[id].Outbox <- msgRes
		}
	}
	s.SimulationEngine.Orchestrator.ClientsResponse = make(map[int64]*pb.Message)
	s.SimulationEngine.NextState()
	return nil
}
func (s *StateDispatchingEvents) GetNextState() (utils.State, error) {
	return NewStateRequestingEvents(s.SimulationEngine), nil
}
func (s *StateDispatchingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return false // This state does not allow any type of msg
}

func NewStateDispatchingEvents(s *SimulationEngine) *StateDispatchingEvents {
	return &StateDispatchingEvents{SimulationEngine: s}
}
