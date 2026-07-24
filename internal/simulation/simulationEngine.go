package simulation

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"sync"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"

	"github.com/Norzuiso/orchestrator/internal/orchestrator"
	"github.com/Norzuiso/orchestrator/internal/servers"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

type SimulationEngine struct {
	GrpcServer   *servers.ClientToClientService
	Orchestrator *orchestrator.Orchestrator
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
	// Creation of grpc server
	s.GrpcServer = servers.NewClientToClientService(1001, s.Orchestrator) // TODO - Change the seed value to get it from config
	pb.RegisterBroadcastServer(grpcServer, s.GrpcServer)

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
	http.HandleFunc("POST /msg/all", s.SendMsgToAllClients)
	http.HandleFunc("POST /msg/client", s.SendMsgToClient)
	http.HandleFunc("POST /msg/clients/list", s.SendMsgToClientsList)

	http.HandleFunc("POST /simulation/start", s.StartSimulation)
	http.HandleFunc("GET /simulation/pause", s.StopSimulation)
	http.HandleFunc("GET /simulation/continue", s.StopSimulation)
	http.HandleFunc("GET /simulation/end", s.StopSimulation)

	http.HandleFunc("GET /simulation/next-phase", s.NextPhase)
	http.HandleFunc("GET /simulation/next-epoch", s.NextEpoch)

	err := http.ListenAndServe(":8090", nil) // TODO - Port Get it from config

	if err != nil {
		panic(err)
	}
}

func (s *SimulationEngine) errorHandler() {
	log.Fatal("Error runing program")
}

func StartSimulationEnine() {

	// Orchestrator
	orchestrator := orchestrator.NewOrquestrator()

	se := &SimulationEngine{Orchestrator: orchestrator}
	se.Orchestrator.NextEpoch()

	wg := sync.WaitGroup{}

	defer se.errorHandler()
	wg.Add(2)

	go se.GrpcConnect()
	log.Println("Listening for RPC on 127.0.0.1:8080")

	go se.HttpConnect()
	log.Println("Listening for HTTP on 127.0.0.1:8090")

	wg.Wait()
}

// adaptador para slog (la lib es agnóstica del logger)
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
