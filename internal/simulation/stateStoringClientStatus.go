package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateStoringClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (s *StateStoringClientStatus) StartState() {
	// log.Printf("State: %v", s.GetStateName())
}

func (s *StateStoringClientStatus) GetStateName() string {
	return utils.StoringClientStatusStr
}

func (s *StateStoringClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {

	return nil
}
func (s *StateStoringClientStatus) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (s *StateStoringClientStatus) GetNextState() (utils.State, error)    { return nil, nil }
func (s *StateStoringClientStatus) IsMsgTypeAllowIt(msg *pb.Message) bool { return false }

func NewStateStoringClientStatus(s *SimulationEngine) *StateStoringClientStatus {
	return &StateStoringClientStatus{SimulationEngine: s}
}
