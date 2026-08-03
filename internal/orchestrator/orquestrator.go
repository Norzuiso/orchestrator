package orchestrator

import (
	"log"
	"math"
	"slices"
	"strconv"
	"sync"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type Orchestrator struct {
	Seed        int64
	Epoch       int64
	MaxOfEpochs int64

	ActiveClients            map[int64]*pb.Client
	ClientToClientConnection map[int64]*pb.ClientConnectionList
	ClientsResponse          map[int64]*pb.Message

	mu sync.Mutex

	ClientsRequest []int64
	SeedEpoch      int64
}

func (o *Orchestrator) RemoveFromRequest(fromId int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.ClientsRequest = slices.DeleteFunc(o.ClientsRequest, func(id int64) bool {
		return id == fromId
	})

}

func (o *Orchestrator) GetEpoch() int64 {
	if o.Epoch != 0 {
		return o.Epoch
	}
	return 8080
}

func (o *Orchestrator) GetClientSeed(clientId int64) int64 {
	seedEpochStr := strconv.FormatInt(o.GetSeedEpoch(), 10)
	strResult := seedEpochStr + strconv.FormatInt(clientId, 10)
	seedEpochClient, err := strconv.ParseInt(strResult, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return seedEpochClient
}

func (o *Orchestrator) GetSeedEpoch() int64 {
	return o.SeedEpoch
}

func NewOrquestrator(seed int64) *Orchestrator {
	o := &Orchestrator{}
	o.ActiveClients = make(map[int64]*pb.Client)
	o.ClientsResponse = make(map[int64]*pb.Message)
	o.ClientsRequest = make([]int64, 0)
	o.ClientToClientConnection = make(map[int64]*pb.ClientConnectionList)
	o.Epoch = 0
	o.Seed = seed
	o.MaxOfEpochs = int64(math.Inf(1))
	o.generateSeedEpoch()
	return o
}

func (o *Orchestrator) StartSimualtion() {
	o.Epoch = 0
}
func (o *Orchestrator) PauseSimualtion() {
	o.Epoch = 0
}

func (o *Orchestrator) NextEpoch() {
	o.Epoch += 1
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
