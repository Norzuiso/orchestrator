package utils

import (
	"fmt"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

func BuildErrorMsg(msg *pb.Message, err error) *pb.Message {
	msg.Content = fmt.Sprintf("Error: %v", err)
	msg.MessageType = pb.MessageType_MESSAGE_TYPE_ERROR
	msg.SenderId = 0
	return msg
}
func BuildPhaseErrorMsg(msg *pb.Message, err error) *pb.Message {
	msg.Content = fmt.Sprintf("Error: %v", err)
	msg.MessageType = pb.MessageType_MESSAGE_TYPE_ERROR_PHASE
	msg.SenderId = 0
	return msg
}
func BuildNoAllowMsgErrorMsg(msg *pb.Message) *pb.Message {
	return BuildErrorMsg(msg, fmt.Errorf("Message not allow it. Orchestrator is not reciving any message"))
}
