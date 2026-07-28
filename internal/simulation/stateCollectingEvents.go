package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateCollectingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateCollectingEvents) GetStateName() string                                   { return "" }
func (w *StateCollectingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateCollectingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StateCollectingEvents) NextState() (utils.State, error)                        { return nil, nil }
func (w *StateCollectingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool                  { return false }

func NewStateCollectingEvents(s *SimulationEngine) *StateCollectingEvents {
	return &StateCollectingEvents{SimulationEngine: s}
}
