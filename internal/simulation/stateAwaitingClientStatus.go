package simulation

import (
	"fmt"
	"slices"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateAwaitingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (s *StateAwaitingClientStatus) StartState() {
	// log.Printf("State: %v", s.GetStateName())
}

func (s *StateAwaitingClientStatus) GetStateName() string {
	return utils.AwaitingClientStatusStr
}

func (s *StateAwaitingClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {

	_, err := s.SimulationEngine.Storage.ClientsResponseGet(msg.SenderId)
	if err != nil {
		return err
	}

	// Check if orchestrator is waiting a response from client
	if !slices.Contains(s.SimulationEngine.Orchestrator.GetClientsRequest(), msg.SenderId) {
		return fmt.Errorf("No event expected from client: %v", msg.SenderId)
	}

	s.SimulationEngine.Storage.ClientsResponseSave(msg.SenderId, msg)

	// if we dont have any pending msg to read from a client change to the next state
	if s.SimulationEngine.Orchestrator.RemoveFromRequest(msg.SenderId) {
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
