package models

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"google.golang.org/protobuf/types/known/anypb"
)

type MsgToAllClients struct {
	Msg MsgRequest `json:"msg"`
}

type MsgToClient struct {
	ClientId int64      `json:"client_id"`
	Msg      MsgRequest `json:"msg"`
}

type MsgToClientsList struct {
	ClientIds []int64    `json:"clients"`
	Msg       MsgRequest `json:"msg"`
}

type MsgRequest struct {
	Content    string                `json:"content"`
	MsgType    pb.MessageType        `json:"msg_type"`
	Attributes map[string]*anypb.Any `json:"attributes"`
}

func (m *MsgRequest) MsgRequestToMessage(epoch float32, senderId int64) *pb.Message {
	return &pb.Message{
		SenderId:    senderId,
		Epoch:       epoch,
		MessageType: m.MsgType,
		Content:     m.Content,
		Attributes:  m.Attributes,
	}
}
