package simulation

import (
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateRequestingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (s *StateRequestingClientStatus) StartState() {
	log.Printf("\nState: %s", s.GetStateName)
}

func (s *StateRequestingClientStatus) GetStateName() string { return "" }
func (s *StateRequestingClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (s *StateRequestingClientStatus) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (s *StateRequestingClientStatus) GetNextState() (utils.State, error)    { return nil, nil }
func (s *StateRequestingClientStatus) IsMsgTypeAllowIt(msg *pb.Message) bool { return false }

func NewStateRequestingClientStatus(s *SimulationEngine) *StateRequestingClientStatus {
	return &StateRequestingClientStatus{SimulationEngine: s}
}
