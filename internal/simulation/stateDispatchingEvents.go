package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateDispatchingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateDispatchingEvents) GetStateName() string                                   { return "" }
func (w *StateDispatchingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateDispatchingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateDispatchingEvents) NextState() (utils.State, error)                        { return nil, nil }
func (w *StateDispatchingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool                  { return false }

func NewStateDispatchingEvents(s *SimulationEngine) *StateDispatchingEvents {
	return &StateDispatchingEvents{SimulationEngine: s}
}
