package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/Norzuiso/orchestrator/internal/orquestrator"
	"github.com/Norzuiso/orchestrator/internal/register"
	"github.com/Norzuiso/orchestrator/internal/servers"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
)

// adaptador para slog (la lib es agnóstica del logger)
func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	opts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}
	orquestrator := orquestrator.NewOrquestrator(":8081")
	server := servers.NewClientToClientService(1001, orquestrator)

	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(logging.StreamServerInterceptor(interceptorLogger(logger), opts...)),
		grpc.ChainUnaryInterceptor(logging.UnaryServerInterceptor(interceptorLogger(logger), opts...)),
		grpc.StatsHandler(&register.StatsHandler{}),
	)

	pb.RegisterBroadcastServer(grpcServer, server)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	fmt.Println("Server started at port :8080")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server: %v", err)
	}

}
