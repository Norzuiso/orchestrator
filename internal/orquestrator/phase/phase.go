package Phase

import (
	"log"
	"slices"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type Phase struct {
	name         string
	allowMsgType []*pb.MessageType
}

func (p *Phase) String() string {
	return ""
}

func (p *Phase) GetName() string {
	return p.name
}

func (p *Phase) GetAllowMsgTypeStr() []string {
	allowTypesStr := []string{}
	for _, msgType := range p.allowMsgType {
		log.Println(msgType)

	}

	return allowTypesStr
}

func (p *Phase) GetAllowMsgType() []*pb.MessageType {
	return p.allowMsgType
}

func (p *Phase) IsMsgTypeAllowIt(msgType *pb.MessageType) bool {
	return slices.Contains(p.allowMsgType, msgType)
}

func NewPhase() *Phase {
	return &Phase{}
}

func (p *Phase) WaitingConnectionPhase() {
	p.name = "WaitingConnection"
	p.allowMsgType = []*pb.MessageType{pb.MessageType_MESSAGE_TYPE_OPEN_STREAM.Enum()}
}

func (p *Phase) RequestEventPhase() {
	p.name = "RequestEvent"
	p.allowMsgType = []*pb.MessageType{pb.MessageType_MESSAGE_TYPE_RECOLECT_EVENT.Enum()}
}

func (p *Phase) RecolectEventPhase() {
	p.name = "RecolectEvent"
	p.allowMsgType = []*pb.MessageType{pb.MessageType_MESSAGE_TYPE_SEND_EVENT.Enum()}
}

func (p *Phase) ApplyEventPhase() {
	p.name = "ApplyEvent"
	p.allowMsgType = []*pb.MessageType{pb.MessageType_MESSAGE_TYPE_APPLY_EVENT.Enum()}
}

func (p *Phase) RequestClientStatusPhase() {
	p.name = "RequestClientStatus"
	p.allowMsgType = []*pb.MessageType{pb.MessageType_MESSAGE_TYPE_CLIENT_STATUS.Enum()}
}

func (p *Phase) PausePhase() {
	p.name = "Pause"
	p.allowMsgType = []*pb.MessageType{}
}
