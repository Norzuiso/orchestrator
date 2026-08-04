package orchestrator

import (
	"log"
	"math"
	"slices"
	"strconv"
	"sync"
)

type Orchestrator struct {
	Seed        int64
	Epoch       int64
	MaxOfEpochs int64
	SeedEpoch   int64

	mu              sync.Mutex
	clientsRequest  []int64
	clientsResponse []int64
}

func (o *Orchestrator) ResetClientsResponse() {
	o.clientsResponse = make([]int64, 0)
}

func (o *Orchestrator) GetClientsResponse() []int64 {
	o.mu.Lock()
	defer o.mu.Unlock()

	cp := make([]int64, len(o.clientsResponse))
	copy(cp, o.clientsResponse)
	return cp
}

func (o *Orchestrator) AppendClientsResponse(id int64) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.clientsResponse = append(o.clientsResponse, id)
}

func (o *Orchestrator) RemoveFromResponse(fromId int64) bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.clientsResponse = slices.DeleteFunc(o.clientsResponse, func(id int64) bool {
		return id == fromId
	})
	return len(o.clientsResponse) == 0
}

func (o *Orchestrator) ResetClientsRequest() {
	o.clientsRequest = make([]int64, 0)
}

func (o *Orchestrator) GetClientsRequest() []int64 {
	o.mu.Lock()
	defer o.mu.Unlock()

	cp := make([]int64, len(o.clientsRequest))
	copy(cp, o.clientsRequest)
	return cp
}

func (o *Orchestrator) AppendClientsRequest(id int64) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.clientsRequest = append(o.clientsRequest, id)
}

func (o *Orchestrator) RemoveFromRequest(fromId int64) bool {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.clientsRequest = slices.DeleteFunc(o.clientsRequest, func(id int64) bool {
		return id == fromId
	})
	return len(o.clientsRequest) == 0
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

	o.ResetClientsRequest()
	o.ResetClientsResponse()

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
