package register

import (
	"context"
	"sync"

	registerv1 "github.com/Norzuiso/protocol/gen/go/orchestrator/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	registerv1.UnimplementedRegistryServiceServer
	mu       sync.Mutex
	services map[string]Service
}

type Service struct {
	Id       uuid.UUID
	Name     string
	Endpoint string
	Category string
}

func NewServer() *Server {
	return &Server{services: make(map[string]Service)}
}

func (s *Server) Register(ctx context.Context, req *registerv1.RegistryRequest) (*registerv1.RegistryResponse, error) {
	if req.Endpoint == "" {
		return nil, status.Error(codes.InvalidArgument, "Endpoint should not be empty")
	}
	newUuid, err := uuid.NewUUID()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error generating service id")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.services[req.Endpoint]; ok {
		return nil, status.Errorf(codes.InvalidArgument, "Endpoint should be unique. %s is already register", req.Endpoint)
	}

	s.services[req.Endpoint] = Service{
		Id:       newUuid,
		Name:     req.Name,
		Endpoint: req.Endpoint,
		Category: req.Category,
	}
	return &registerv1.RegistryResponse{ResponseMessage: "Service register"}, nil
}
