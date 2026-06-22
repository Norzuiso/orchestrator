package main

import (
	"fmt"
	"log"
	"net"

	"github.com/Norzuiso/orchestrator/internal/register"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"google.golang.org/grpc"
)

func main() {

	grpcServer := grpc.NewServer()

	pool := &register.Pool{
		Connection: make(map[string]*register.Connection),
	}

	pb.RegisterBroadcastServer(grpcServer, pool)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}

	fmt.Println("Server started at port :8080")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Error creating the server: %v", err)
	}

}
