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
	clientEventsRequest := s.SimulationEngine.Orchestrator.ClientEventsRequest
	clientEventsResponse := s.SimulationEngine.Orchestrator.ClientEventsResponse

	// Check if orchestrator is waiting a response from client
	if !slices.Contains(clientEventsRequest, msg.SenderId) {
		return fmt.Errorf("No event expected from client: %v", msg.SenderId)
	}

	// check if client has a response
	if _, ok := clientEventsResponse[msg.SenderId]; ok {
		return fmt.Errorf("Client %v, already has an event register", msg.SenderId)
	}

	// remove clientid from the list of requested clients
	clientEventsRequest = slices.DeleteFunc(clientEventsRequest, func(id int64) bool {
		return id == msg.SenderId
	})

	// store msg from client into the events response
	clientEventsResponse[msg.SenderId] = msg

	s.SimulationEngine.Orchestrator.ClientEventsRequest = clientEventsRequest
	s.SimulationEngine.Orchestrator.ClientEventsResponse = clientEventsResponse

	// if we dont have any client pending change to the next state
	if len(s.SimulationEngine.Orchestrator.ClientEventsRequest) == 0 {
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
