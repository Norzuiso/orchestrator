package main

import (
	"log"
	"net"

	"github.com/Norzuiso/orchestrator/internal/register"
	registerv1 "github.com/Norzuiso/protocol/gen/go/orchestrator/v1"
	"google.golang.org/grpc"
)

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer := grpc.NewServer()
	registerv1.RegisterRegistryServiceServer(grpcServer, register.NewServer())

	log.Println("Server listening port :50051")
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
