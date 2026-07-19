package register

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type MsgHandler struct {
	*pb.Message
	stream pb.Broadcast_ClientToClientMessageServer
}

func (m *MsgHandler) BuildDefaultMsg() {
	m.Message.MessageType = pb.MessageType_MESSAGE_TYPE_DEFAULT.Enum()
}

func (m *MsgHandler) BuildErrorMsg(err error) {
	m.MessageType = pb.MessageType_MESSAGE_TYPE_ERROR.Enum()
	m.Content = err.Error()
}

func (m *MsgHandler) SendError(err error) {
	m.BuildErrorMsg(err)
	m.SendMsg()
}
func (m *MsgHandler) SendMsgContent(content string) {
	m.Message.Content = content
	m.SenderId = 0
	m.stream.Send(m.Message)
}

func (m *MsgHandler) SendMsg() {
	m.stream.Send(m.Message)
}

func (m *MsgHandler) SendPhaseError(content string) {
	m.SenderId = 0
	m.MessageType = pb.MessageType_MESSAGE_TYPE_ERROR_PHASE.Enum()
	m.Content = content
}

func NewMessageHandler(msg *pb.Message, stream pb.Broadcast_ClientToClientMessageServer) *MsgHandler {
	return &MsgHandler{msg, stream}
}
