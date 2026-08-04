package simulation

import (
	"fmt"
	"log"
	"slices"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateCollectingEvents struct {
	SimulationEngine *SimulationEngine
}

func (s *StateCollectingEvents) StartState() {
	log.Printf("State: %v", s.GetStateName())
}

func (s *StateCollectingEvents) GetStateName() string {
	return utils.CollectingEventsStr
}

func (s *StateCollectingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	senderId := msg.SenderId

	// Check if orchestrator is waiting a response from client
	if !slices.Contains(s.SimulationEngine.Orchestrator.GetClientsRequest(), senderId) {
		return fmt.Errorf("No event expected from client: %v", senderId)
	}

	s.SimulationEngine.Storage.ClientsResponseSave(senderId, msg)
	s.SimulationEngine.Orchestrator.AppendClientsResponse(senderId)

	// if we dont have any client pending change to the next state
	if s.SimulationEngine.Orchestrator.RemoveFromRequest(senderId) {
		s.SimulationEngine.NextState()
	}
	return nil
}

func (s *StateCollectingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}

func (s *StateCollectingEvents) GetNextState() (utils.State, error) {
	nextState := NewStateDispatchingEvents(s.SimulationEngine)
	return nextState, nil
}

func (s *StateCollectingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return msg.GetMessageType() == pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE
}

func NewStateCollectingEvents(s *SimulationEngine) *StateCollectingEvents {
	return &StateCollectingEvents{SimulationEngine: s}
}
