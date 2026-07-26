package Phase

import (
	"log"
	"slices"

	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

var nextPhases map[string]*Phase = make(map[string]*Phase)

type Phase struct {
	name         string
	allowMsgType []pb.MessageType
}

func (p *Phase) IsWaitingConnection() bool {
	return p.name == utils.WaitingConnectionsStr
}

func (p *Phase) String() string {
	return ""
}

func (p *Phase) GetName() string {
	return p.name
}

func (p *Phase) GetNextPhase() *Phase {
	return nextPhases[p.name]
}

func (p *Phase) NextPhase() {
	p = nextPhases[p.name]
}

func (p *Phase) GetAllowMsgTypeStr() []string {
	allowTypesStr := []string{}
	for _, msgType := range p.allowMsgType {
		log.Println(msgType)

	}

	return allowTypesStr
}

func (p *Phase) GetAllowMsgType() []pb.MessageType {
	return p.allowMsgType
}

func (p *Phase) IsMsgTypeAllowIt(msgType pb.MessageType) bool {
	return slices.Contains(p.allowMsgType, msgType)
}

func NewPhase() *Phase {
	p := &Phase{}

	// nextPhases[utils.WaitingConnectionStr] = p.RequestEventPhase()
	// nextPhases[utils.AwaitingEventRequestStr] = p.RecolectEventPhase()
	// nextPhases[utils.CollectEventsStr] = p.ApplyEventPhase()
	// nextPhases[utils.ApplyEventStr] = p.RequestClientStatusPhase()
	// nextPhases[utils.RequestClientStatusStr] = p.StoringClientStatusPhase()
	// nextPhases[utils.StoringClientStatusStr] = p.RequestEventPhase()
	// nextPhases[utils.PauseStr] = p.RequestClientStatusPhase()

	return p
}

func (p *Phase) WaitingConnectionPhase() *Phase {
	return &Phase{
		utils.WaitingConnectionsStr,
		[]pb.MessageType{pb.MessageType_MESSAGE_TYPE_OPEN_STREAM},
	}
}

func (p *Phase) RequestEventPhase() *Phase {

	return &Phase{
		utils.AwaitingEventResponsesStr,
		[]pb.MessageType{pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE},
	}
}

func (p *Phase) RecolectEventPhase() *Phase {

	return &Phase{
		utils.AwaitingEventResponsesStr,
		[]pb.MessageType{pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE},
	}
}

func (p *Phase) ApplyEventPhase() *Phase {

	return &Phase{
		utils.AwaitingEventResponsesStr,
		[]pb.MessageType{pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE},
	}
}

func (p *Phase) RequestClientStatusPhase() *Phase {

	return &Phase{
		utils.PausedStr,
		[]pb.MessageType{pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE},
	}
}
func (p *Phase) StoringClientStatusPhase() *Phase {

	return &Phase{
		utils.StoringClientStatusStr,
		[]pb.MessageType{pb.MessageType_MESSAGE_TYPE_EVENT_RESPONSE},
	}
}

// We should have something as waiting responses... please, evaluate this options

func (p *Phase) PausePhase() *Phase {

	return &Phase{
		utils.PausedStr,
		[]pb.MessageType{},
	}
}
