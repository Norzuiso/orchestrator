package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateRequestingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateRequestingEvents) GetStateName() string                                   { return "" }
func (w *StateRequestingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateRequestingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateRequestingEvents) NextState() (utils.State, error)                        { return nil, nil }
func (w *StateRequestingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool                  { return false }

func NewStateRequestingEvents(s *SimulationEngine) *StateRequestingEvents {
	return &StateRequestingEvents{SimulationEngine: s}
}
