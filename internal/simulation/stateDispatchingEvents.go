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
	responses := s.SimulationEngine.Orchestrator.GetClientsResponse()

	for _, senderId := range responses {
		conns, err := s.SimulationEngine.Storage.ClientToClientGet(senderId)
		if err != nil {
			return err
		}
		msgRes, err := s.SimulationEngine.Storage.ClientsResponseGet(senderId)
		if err != nil {
			return err
		}
		msgRes.Epoch = s.SimulationEngine.Orchestrator.Epoch
		msgRes.MessageType = pb.MessageType_MESSAGE_TYPE_EVENT_DISPATCH

		for _, connection := range conns.Connections {

			msgRes.Attributes = connection.Attributes

			done := make(chan error, 1)
			s.SimulationEngine.ClientService.ClientStreams[connection.ToId].Outbox <- models.OutboxItem{Msg: msgRes, Done: done}
			<-done
		}
	}
	s.SimulationEngine.Orchestrator.ResetClientsResponse()
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
