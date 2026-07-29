package orchestrator

import (
	"log"
	"strconv"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type Orchestrator struct {
	Seed  int64
	Epoch int64

	ActiveClients             map[int64]*pb.Client
	ClientToClientConnections map[int64][]int64
	ClientEventsResponse      map[int64]*pb.Message
	ClientEventsRequest       []int64
	SeedEpoch                 int64
}

func (o *Orchestrator) GetSeedEpochToClient(clientId int64) int64 {
	seedEpochStr := strconv.FormatInt(o.SeedEpoch, 10)
	clientIdStr := strconv.FormatInt(clientId, 10)
	seedEpochClientStr := seedEpochStr + clientIdStr
	seedEpochClient, err := strconv.ParseInt(seedEpochClientStr, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return int64(seedEpochClient)
}

func (o *Orchestrator) GetSeedEpoch() int64 {
	return o.SeedEpoch
}

func NewOrquestrator(seed int64) *Orchestrator {
	o := &Orchestrator{}
	o.ActiveClients = make(map[int64]*pb.Client)
	o.ClientEventsResponse = make(map[int64]*pb.Message)
	o.ClientEventsRequest = make([]int64, 0)
	o.ClientToClientConnections = make(map[int64][]int64)
	o.Epoch = 0
	return o
}

func (o *Orchestrator) StartSimualtion() {
	o.Epoch = 0
}
func (o *Orchestrator) PauseSimualtion() {
	o.Epoch = 0
}

func (o *Orchestrator) NextEpoch() {
	o.Epoch = +1
	o.generateSeedEpoch()

}
func (o *Orchestrator) generateSeedEpoch() {
	seedStr := strconv.FormatInt(o.Seed, 10)
	epochStr := strconv.FormatInt(o.Epoch, 10)
	seedEpochStr := epochStr + seedStr
	seedEpoch, err := strconv.ParseInt(seedEpochStr, 10, 64)
	if err != nil {
		log.Fatal(err)
		return
	}
	o.SeedEpoch = seedEpoch
}
