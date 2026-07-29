package simulation

import (
	"log"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateAwaitingEventResponses struct {
	SimulationEngine *SimulationEngine
}

func (s *StateAwaitingEventResponses) StartState() {
	log.Printf("State: %v", s.GetStateName())
}

func (w *StateAwaitingEventResponses) GetStateName() string { return "" }
func (w *StateAwaitingEventResponses) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateAwaitingEventResponses) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}
func (w *StateAwaitingEventResponses) GetNextState() (utils.State, error)    { return nil, nil }
func (w *StateAwaitingEventResponses) IsMsgTypeAllowIt(msg *pb.Message) bool { return false }

func NewStateAwaitingEventResponses(s *SimulationEngine) *StateAwaitingEventResponses {
	return &StateAwaitingEventResponses{SimulationEngine: s}
}
