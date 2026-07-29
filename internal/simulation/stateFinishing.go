package simulation

import (
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateFinishing struct {
	SimulationEngine *SimulationEngine
}

func (s *StateFinishing) StartState() {
	log.Printf("State: %v", s.GetStateName())
}

func (s *StateFinishing) GetStateName() string {
	return ""
}

func (s *StateFinishing) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}

func (s *StateFinishing) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}

func (s *StateFinishing) GetNextState() (utils.State, error) {
	return nil, nil
}

func (s *StateFinishing) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return false
}

func NewStateFinishing(s *SimulationEngine) *StateFinishing {
	return &StateFinishing{SimulationEngine: s}
}
