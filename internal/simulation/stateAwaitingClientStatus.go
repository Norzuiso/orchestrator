package simulation

import (
	"fmt"
	"log"
	"slices"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateAwaitingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (s *StateAwaitingClientStatus) StartState() {
	log.Printf("State: %v", s.GetStateName())
}

func (s *StateAwaitingClientStatus) GetStateName() string {
	return utils.AwaitingClientStatusStr
}

func (s *StateAwaitingClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {

	clientEventsRequest := s.SimulationEngine.Orchestrator.ClientsRequest
	clientEventsResponse := s.SimulationEngine.Orchestrator.ClientsResponse

	// Check if orchestrator is waiting a response from client
	if !slices.Contains(clientEventsRequest, msg.SenderId) {
		return fmt.Errorf("No event expected from client: %v", msg.SenderId)
	}

	// check if client has a response
	if _, ok := clientEventsResponse[msg.SenderId]; ok {
		return fmt.Errorf("Client %v, already has an status register", msg.SenderId)
	}

	// remove clientid from the list of requested clients
	clientEventsRequest = slices.DeleteFunc(clientEventsRequest, func(id int64) bool {
		return id == msg.SenderId
	})

	// store msg from client into the events response
	clientEventsResponse[msg.SenderId] = msg

	s.SimulationEngine.Orchestrator.ClientsRequest = clientEventsRequest
	s.SimulationEngine.Orchestrator.ClientsResponse = clientEventsResponse

	// if we dont have any client pending change to the next state
	if len(s.SimulationEngine.Orchestrator.ClientsRequest) == 0 {
		s.SimulationEngine.NextState()
	}
	return nil
}
func (s *StateAwaitingClientStatus) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (s *StateAwaitingClientStatus) GetNextState() (utils.State, error)    { return nil, nil }
func (s *StateAwaitingClientStatus) IsMsgTypeAllowIt(msg *pb.Message) bool { return false }

func NewStateAwaitingClientStatus(s *SimulationEngine) *StateAwaitingClientStatus {
	return &StateAwaitingClientStatus{SimulationEngine: s}
}
