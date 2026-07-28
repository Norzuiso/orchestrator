package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateStoringClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (w *StateStoringClientStatus) GetStateName() string { return "" }
func (w *StateStoringClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateStoringClientStatus) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateStoringClientStatus) NextState() (utils.State, error)       { return nil, nil }
func (w *StateStoringClientStatus) IsMsgTypeAllowIt(msg *pb.Message) bool { return false }

func NewStateStoringClientStatus(s *SimulationEngine) *StateStoringClientStatus {
	return &StateStoringClientStatus{SimulationEngine: s}
}
