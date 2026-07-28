package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StatePaused struct {
	SimulationEngine *SimulationEngine
}

func (w *StatePaused) GetStateName() string                                   { return "" }
func (w *StatePaused) ReadMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StatePaused) SendMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (w *StatePaused) NextState() (utils.State, error)                        { return nil, nil }
func (w *StatePaused) IsMsgTypeAllowIt(msg *pb.Message) bool                  { return false }

func NewStatePaused(s *SimulationEngine) *StatePaused {
	return &StatePaused{SimulationEngine: s}
}
