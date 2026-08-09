package register

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type MsgHandler struct {
	*pb.Message
	stream pb.Broadcast_ClientToClientMessageServer
}

func (m *MsgHandler) BuildDefaultMsg() {
	m.Message.MessageType = pb.MessageType_MESSAGE_TYPE_DEFAULT
}

func (m *MsgHandler) BuildErrorMsg(err error) {
	m.MessageType = pb.MessageType_MESSAGE_TYPE_ERROR
	m.Content = err.Error()
}

func (m *MsgHandler) SendError(err error) error {
	m.BuildErrorMsg(err)
	return m.SendMsg(m.Message)
}

func (m *MsgHandler) SendMsg(msg *pb.Message) error {
	return m.stream.Send(msg)
}

func (m *MsgHandler) SendMsgContent(content string) error {
	m.Message.Content = content
	return m.stream.Send(m.Message)
}

func (m *MsgHandler) SendPhaseError(content string) {
	m.MessageType = pb.MessageType_MESSAGE_TYPE_ERROR_PHASE
	m.Content = content
}

func NewMessageHandler(msg *pb.Message, stream pb.Broadcast_ClientToClientMessageServer) *MsgHandler {
	return &MsgHandler{msg, stream}
}
