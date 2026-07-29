package simulation

import (
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateAwaitingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (s *StateAwaitingClientStatus) StartState() {
	log.Printf("State: %v", s.GetStateName())
}

func (w *StateAwaitingClientStatus) GetStateName() string { return "" }
func (w *StateAwaitingClientStatus) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateAwaitingClientStatus) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateAwaitingClientStatus) GetNextState() (utils.State, error)    { return nil, nil }
func (w *StateAwaitingClientStatus) IsMsgTypeAllowIt(msg *pb.Message) bool { return false }

func NewStateAwaitingClientStatus(s *SimulationEngine) *StateAwaitingClientStatus {
	return &StateAwaitingClientStatus{SimulationEngine: s}
}
