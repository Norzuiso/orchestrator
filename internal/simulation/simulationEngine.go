package simulation

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/storage"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"

	"github.com/Norzuiso/orchestrator/internal/orchestrator"
	"github.com/Norzuiso/orchestrator/internal/servers"
	"github.com/Norzuiso/orchestrator/internal/utils"
	"google.golang.org/grpc"
)

type SimulationEngine struct {
	ClientService *servers.ClientToClientService
	Orchestrator  *orchestrator.Orchestrator
	Storage       *storage.Storage

	StateAwaitingClientStatus   utils.State
	StateAwaitingEventResponses utils.State
	StateCollectingEvents       utils.State
	StateDispatchingEvents      utils.State
	StateEnd                    utils.State
	StatePaused                 utils.State
	StateRequestingClientStatus utils.State
	StateRequestingEvents       utils.State
	StateStoringClientStatus    utils.State
	StateWaitingConnections     utils.State

	CurrentState utils.State

	mu           sync.Mutex
	StateChanged chan struct{}
	LogChannel   chan string
}

func (s *SimulationEngine) EndEpoch() {
	s.ClientService.LogStorage.LogMessage(fmt.Sprintf("EndEpoch: %v", s.Orchestrator.SeedEpoch), &pb.Message{Epoch: s.Orchestrator.Epoch}, fmt.Sprintf("EndEpoch: %v", s.Orchestrator.SeedEpoch))
	if !s.Orchestrator.StepsMode {
		s.StartEpoch()
	}
}

func (s *SimulationEngine) StartEpoch() {
	if s.Orchestrator.Epoch == s.Orchestrator.MaxOfEpochs {
		s.CurrentState = s.StateEnd
		s.ClientService.LogStorage.LogMessage("End of simulation", &pb.Message{Epoch: s.Orchestrator.Epoch}, "End of simulation")
		return
	}

	s.Orchestrator.NextEpoch()
	s.LogChannel <- fmt.Sprintf("StartEpoch: %v", s.Orchestrator.SeedEpoch)
	s.CurrentState = s.StateRequestingEvents
	s.CurrentState.StartState()
}

func (s *SimulationEngine) GetCurrentState() utils.State {
	return s.CurrentState
}

func (s *SimulationEngine) NextState() {
	s.CurrentState, _ = s.CurrentState.GetNextState()
	s.CurrentState.StartState()
	s.mu.Lock()
	close(s.StateChanged)
	s.StateChanged = make(chan struct{})
	s.mu.Unlock()
}

func (s *SimulationEngine) EndSimulation() {
	s.CurrentState = s.StateEnd
	s.CurrentState.StartState()

}

func (s *SimulationEngine) StateChangedChan() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.StateChanged
}

func (s *SimulationEngine) WriteLogs(str string) {
	s.LogChannel <- str
}

func (s *SimulationEngine) GrpcConnect() {
	// GRPC server logging
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// opts := []logging.Option{
	// 	logging.WithLogOnEvents(
	// 		logging.StartCall,
	// 		logging.FinishCall,
	// 		logging.PayloadReceived,
	// 		logging.PayloadSent,
	// 	),
	// }
	grpcServer := grpc.NewServer(
	// grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
	//	l.Log(ctx, slog.Level(lvl), msg, fields...)
	//}), opts...)),
	// grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
	//		l.Log(ctx, slog.Level(lvl), msg, fields...)
	//, opts...)),
	// grpc.StatsHandler(&register.StatsHandler{}),
	)

	pb.RegisterBroadcastServer(grpcServer, s.ClientService)

	// run server
	listener, err := net.Listen("tcp", ":8080") // TODO - Change port to get it from config
	if err != nil {
		panic(fmt.Errorf("Error creating server: %v", err))
	}

	fmt.Println("Server started at port :8080")

	if err := grpcServer.Serve(listener); err != nil {
		panic(fmt.Errorf("Error creating the server: %v", err))
	}
}

func NewSimulationEngine(seed int64) *SimulationEngine {

	s := &SimulationEngine{}
	s.Storage = &storage.Storage{}
	s.Storage.OpenDb("my.db")
	s.LogChannel = make(chan string, 10)
	s.StateChanged = make(chan struct{})

	// Orchestrator
	orch := orchestrator.NewOrquestrator(seed)
	clientService := &servers.ClientToClientService{
		ClientStreams:   make(map[int64]*models.Connection),
		Orchestrator:    orch,
		StateProvider:   s,
		StorageProvider: s.Storage,
		LogStorage:      &storage.LogStorage{},
		LogsProvider:    s,
	}
	s.ClientService = clientService
	s.Orchestrator = orch

	s.ClientService.LogStorage.OpenDb("log-One-Simulation.db")

	s.StateWaitingConnections = NewStateWaitingConnections(s)
	s.StateAwaitingClientStatus = NewStateAwaitingClientStatus(s)
	s.StateAwaitingEventResponses = NewStateAwaitingEventResponses(s)
	s.StateCollectingEvents = NewStateCollectingEvents(s)
	s.StateDispatchingEvents = NewStateDispatchingEvents(s)
	s.StateEnd = NewStateEnd(s)
	s.StatePaused = NewStatePaused(s)
	s.StateRequestingClientStatus = NewStateRequestingClientStatus(s)
	s.StateRequestingEvents = NewStateRequestingEvents(s)
	s.StateStoringClientStatus = NewStateStoringClientStatus(s)

	s.CurrentState = s.StateWaitingConnections

	return s
}

func (s *SimulationEngine) errorHandler() {
	log.Fatal("Error runing program")
}

func StartSimulationEnine() {

	se := NewSimulationEngine(1004) // TODO - Change the seed value to get it from config
	wg := sync.WaitGroup{}

	defer se.errorHandler()
	wg.Add(2)

	go se.GrpcConnect()
	log.Println("Listening for RPC on 127.0.0.1:8080")

	go se.HttpConnect()
	log.Println("Listening for HTTP on 127.0.0.1:8090")

	wg.Wait()
	se.Storage.CloseDb()
	se.ClientService.LogStorage.CloseDb()

}
