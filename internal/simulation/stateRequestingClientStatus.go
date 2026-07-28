package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateRequestingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (w *StateRequestingClientStatus) GetStateName() string { return "" }
func (w *StateRequestingClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateRequestingClientStatus) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateRequestingClientStatus) NextState() (utils.State, error)       { return nil, nil }
func (w *StateRequestingClientStatus) IsMsgTypeAllowIt(msg *pb.Message) bool { return false }

func NewStateRequestingClientStatus(s *SimulationEngine) *StateRequestingClientStatus {
	return &StateRequestingClientStatus{SimulationEngine: s}
}
