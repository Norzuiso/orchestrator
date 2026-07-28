package utils

import (
	"github.com/Norzuiso/orchestrator/internal/models"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type State interface {
	GetStateName() string
	ReadMsg(msg *pb.Message, conn *models.Connection) error
	SendMsg(msg *pb.Message, conn *models.Connection) error
	NextState() (State, error)
	IsMsgTypeAllowIt(msg *pb.Message) bool
}
