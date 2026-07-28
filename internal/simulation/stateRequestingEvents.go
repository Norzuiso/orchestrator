package simulation

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateRequestingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateRequestingEvents) GetStateName() string {
	return utils.RequestingEventsStr
}

func (w *StateRequestingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	conn.Outbox <- utils.BuildNoAllowMsgErrorMsg(msg)
	return nil
}

func (w *StateRequestingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error {
	clientStreams := w.SimulationEngine.GrpcServer.ClientStreams
	requestMsg := &pb.Message{
		SenderId:    0,
		Epoch:       w.SimulationEngine.Orchestrator.Epoch,
		MessageType: pb.MessageType_MESSAGE_TYPE_REQUEST_EVENT,
		Content:     "",
	}

	for id, client := range clientStreams {
		w.SimulationEngine.Orchestrator.ClientEventsRequest = append(w.SimulationEngine.Orchestrator.ClientEventsRequest, id)
		client.Outbox <- requestMsg
	}
	w.SimulationEngine.NextState()
	return nil
}

func (w *StateRequestingEvents) NextState() (utils.State, error) {
	return NewStateCollectingEvents(w.SimulationEngine), nil
}

func (w *StateRequestingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return false // This state does not allow any type of msg
}

func NewStateRequestingEvents(s *SimulationEngine) *StateRequestingEvents {
	return &StateRequestingEvents{SimulationEngine: s}
}
