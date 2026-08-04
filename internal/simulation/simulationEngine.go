package simulation

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/storage"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"

	"github.com/Norzuiso/orchestrator/internal/orchestrator"
	"github.com/Norzuiso/orchestrator/internal/servers"
	"github.com/Norzuiso/orchestrator/internal/utils"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
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
}

func (s *SimulationEngine) EndEpoch() {
	s.ClientService.LogStorage.LogMessage(fmt.Sprintf("EndEpoch: %v", s.Orchestrator.SeedEpoch), &pb.Message{Epoch: s.Orchestrator.Epoch}, fmt.Sprintf("EndEpoch: %v", s.Orchestrator.SeedEpoch))
	if !s.Orchestrator.StepsMode {
		s.StartEpoch()
	}
}

func (s *SimulationEngine) StartEpoch() {
	s.Orchestrator.NextEpoch()
	s.ClientService.LogStorage.LogMessage(fmt.Sprintf("StartEpoch: %v", s.Orchestrator.SeedEpoch), &pb.Message{Epoch: s.Orchestrator.Epoch}, fmt.Sprintf("StartEpoch: %v", s.Orchestrator.SeedEpoch))
	s.CurrentState = s.StateRequestingEvents
	s.CurrentState.StartState()
}

func (s *SimulationEngine) GetCurrentState() utils.State {
	return s.CurrentState
}

func (s *SimulationEngine) NextState() {
	s.CurrentState, _ = s.CurrentState.GetNextState()
	s.CurrentState.StartState()
}

func NewSimulationEngine() *SimulationEngine {

	s := &SimulationEngine{}
	s.Storage = &storage.Storage{}
	s.Storage.OpenDb()

	// Orchestrator
	orch := orchestrator.NewOrquestrator(1004) // TODO - Change the seed value to get it from config
	clientService := &servers.ClientToClientService{
		ClientStreams:   make(map[int64]*models.Connection),
		Orchestrator:    orch,
		StateProvider:   s,
		StorageProvider: s.Storage,
		LogStorage:      &storage.LogStorage{},
	}
	s.ClientService = clientService
	s.Orchestrator = orch

	s.ClientService.LogStorage.OpenDb("log.db")

	stateWaitingConnections := NewStateWaitingConnections(s)

	s.CurrentState = stateWaitingConnections
	s.StateWaitingConnections = stateWaitingConnections

	s.StateAwaitingClientStatus = NewStateAwaitingClientStatus(s)
	s.StateAwaitingEventResponses = NewStateAwaitingEventResponses(s)
	s.StateCollectingEvents = NewStateCollectingEvents(s)
	s.StateDispatchingEvents = NewStateDispatchingEvents(s)
	s.StateEnd = NewStateEnd(s)
	s.StatePaused = NewStatePaused(s)
	s.StateRequestingClientStatus = NewStateRequestingClientStatus(s)
	s.StateRequestingEvents = NewStateRequestingEvents(s)
	s.StateStoringClientStatus = NewStateStoringClientStatus(s)

	return s
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
	// grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(interceptorLogger(logger), opts...)),
	// grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(interceptorLogger(logger), opts...)),
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

func (s *SimulationEngine) HttpConnect() {
	http.HandleFunc("POST /msg/all", s.SendMsgToAllClients)           // TODO
	http.HandleFunc("POST /msg/client", s.SendMsgToClient)            // TODO
	http.HandleFunc("POST /msg/clients/list", s.SendMsgToClientsList) // TODO
	http.HandleFunc("POST /msg/client-to-client", s.RegisterClientToClientconnection)

	http.HandleFunc("POST /simulation/start", s.StartSimulation)
	http.HandleFunc("GET /simulation/pause", s.StopSimulation) // TODO
	http.HandleFunc("GET /simulation/state/waiting-connections", s.WaitingConnections)
	http.HandleFunc("GET /simulation/end", s.EndSimulation)

	http.HandleFunc("GET /simulation/next-phase", s.NextStateHttp) // TODO
	http.HandleFunc("GET /simulation/next-epoch", s.NextEpoch)     // TODO

	err := http.ListenAndServe(":8090", nil) // TODO - Port Get it from config

	if err != nil {
		panic(err)
	}
}

func (s *SimulationEngine) errorHandler() {
	log.Fatal("Error runing program")
}

func StartSimulationEnine() {

	se := NewSimulationEngine()
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

// adaptador para slog (la lib es agnóstica del logger)
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
