package simulation

import (
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StatePaused struct {
	SimulationEngine *SimulationEngine
}

func (s *StatePaused) StartState() {
	log.Printf("State: %v", s.GetStateName())
}

func (s *StatePaused) GetStateName() string                                   { return "" }
func (s *StatePaused) ReadMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (s *StatePaused) SendMsg(msg *pb.Message, conn *models.Connection) error { return nil }
func (s *StatePaused) GetNextState() (utils.State, error)                     { return nil, nil }
func (s *StatePaused) IsMsgTypeAllowIt(msg *pb.Message) bool                  { return false }

func NewStatePaused(s *SimulationEngine) *StatePaused {
	return &StatePaused{SimulationEngine: s}
}
