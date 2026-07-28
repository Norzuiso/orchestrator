package simulation

import (
	"fmt"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StateWaitingConnections struct {
	SimulationEngine *SimulationEngine
}

func (s *StateWaitingConnections) GetStateName() string {
	return utils.WaitingConnectionsStr
}

func (s *StateWaitingConnections) ReadMsg(msg *pb.Message, conn *models.Connection) error {
	conn, ok := s.SimulationEngine.GrpcServer.ClientStreams[msg.SenderId]
	if ok {
		return s.SendMsg(utils.BuildErrorMsg(msg, fmt.Errorf("Client is already connected")), conn)
	}
	return nil
}

func (s *StateWaitingConnections) SendMsg(msg *pb.Message, conn *models.Connection) error {
	conn.Outbox <- msg
	return nil
}

func (s *StateWaitingConnections) NextState() (utils.State, error) {
	return NewStateRequestingEvents(s.SimulationEngine), nil
}

func (s *StateWaitingConnections) IsMsgTypeAllowIt(msg *pb.Message) bool {
	return msg.GetMessageType() == pb.MessageType_MESSAGE_TYPE_OPEN_STREAM
}

func NewStateWaitingConnections(s *SimulationEngine) *StateWaitingConnections {
	return &StateWaitingConnections{SimulationEngine: s}
}
