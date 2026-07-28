package simulation

import (
	"fmt"
	"slices"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateCollectingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateCollectingEvents) GetStateName() string {
	return utils.CollectingEventsStr
}

func (w *StateCollectingEvents) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	clientEventsRequest := w.SimulationEngine.Orchestrator.ClientEventsRequest
	clientEventsResponse := w.SimulationEngine.Orchestrator.ClientEventsResponse

	// Check if orchestrator is waiting a response from client
	if !slices.Contains(clientEventsRequest, msg.SenderId) {
		return fmt.Errorf("No event expected from client: %v", msg.SenderId)
	}

	// check if client has a response
	if _, ok := clientEventsResponse[msg.SenderId]; ok {
		return fmt.Errorf("Client %v, already has an event register", msg.SenderId)
	}

	// remove clientid from the list of requested clients
	clientEventsRequest = slices.DeleteFunc(clientEventsRequest, func(id int64) bool {
		return id == msg.SenderId
	})

	// store msg from client into the events response
	clientEventsResponse[msg.SenderId] = msg

	w.SimulationEngine.Orchestrator.ClientEventsRequest = clientEventsRequest
	w.SimulationEngine.Orchestrator.ClientEventsResponse = clientEventsResponse

	// if we dont have any client pending change to the next state
	if len(w.SimulationEngine.Orchestrator.ClientEventsRequest) == 0 {
		w.SimulationEngine.NextState()
	}

	return nil
}

func (w *StateCollectingEvents) SendMsg(msg *pb.Message, conn *models.Connection) error {
	return nil
}

func (w *StateCollectingEvents) NextState() (utils.State, error) {
	return NewStateDispatchingEvents(w.SimulationEngine), nil
}

func (w *StateCollectingEvents) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return msg.GetMessageType() == pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE
}

func NewStateCollectingEvents(s *SimulationEngine) *StateCollectingEvents {
	return &StateCollectingEvents{SimulationEngine: s}
}
