package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateFinishing struct {
	SimulationEngine *SimulationEngine
}

func (w *StateFinishing) GetStateName() string                                   { return "" }
func (w *StateFinishing) ReadMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateFinishing) SendMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateFinishing) NextState() (utils.State, error)                        { return nil, nil }
func (w *StateFinishing) IsMsgTypeAllowIt(msg *pb.Message) bool                  { return false }

func NewStateFinishing(s *SimulationEngine) *StateFinishing {
	return &StateFinishing{SimulationEngine: s}
}
