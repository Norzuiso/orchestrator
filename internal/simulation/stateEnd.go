package simulation

import (
	"fmt"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateEnd struct {
	SimulationEngine *SimulationEngine
}

func (s *StateEnd) StartState() {
	// log.Printf("State: %v", s.GetStateName())
	_ = s.SendMsg(nil, nil)
}

func (s *StateEnd) GetStateName() string {
	return utils.FinishingStr
}

func (s *StateEnd) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}

func (s *StateEnd) SendMsg(msg *pb.Message, conn *models.Connection) error {
	clientStreams := s.SimulationEngine.ClientService.ClientStreams
	for _, client := range clientStreams {
		client.ErrCh <- fmt.Errorf("End of simulation")
	}
	s.SimulationEngine.Orchestrator.ResetClientsRequest()
	return nil
}

func (s *StateEnd) GetNextState() (utils.State, error) {
	return nil, fmt.Errorf("No more states. This is the end of the simulation")
}

func (s *StateEnd) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return false
}

func NewStateEnd(s *SimulationEngine) *StateEnd {
	return &StateEnd{SimulationEngine: s}
}
