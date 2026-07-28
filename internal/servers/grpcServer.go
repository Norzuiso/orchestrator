package servers

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/orchestrator"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type ClientToClientService struct {
	pb.UnimplementedBroadcastServer
	Seed                      int64
	ActiveClients             map[int64]*pb.Client
	ClientStreams             map[int64]*models.Connection
	ClientToClientConnections map[int64][]int64
	Orchestrator              *orchestrator.Orchestrator
	State                     utils.State
}

func NewClientToClientService(seed int64, orchestrator *orchestrator.Orchestrator) *ClientToClientService {
	return &ClientToClientService{
		Seed:                      seed,
		ActiveClients:             make(map[int64]*pb.Client),
		ClientStreams:             make(map[int64]*models.Connection),
		ClientToClientConnections: make(map[int64][]int64),
		Orchestrator:              orchestrator,
	}
}

func (cs *ClientToClientService) ConnectClient(ctx context.Context, req *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	client := req.GetClient()

	if client.GetId() == 0 {
		client.Id = int64(len(cs.ActiveClients) + 1)
	}

	if client.GetName() == "" {
		client.Name = strconv.Itoa(int(client.GetId()))
	}

	if client.GetSeed() == 0 {
		strResult := strconv.Itoa(int(cs.Seed)) + strconv.Itoa(int(client.GetId()))
		newSeed, _ := strconv.Atoi(strResult)
		client.Seed = int64(newSeed)
	}

	cs.ActiveClients[client.Id] = client
	cs.ClientToClientConnections[client.Id] = make([]int64, 0)

	res := &pb.ConnectionResponse{Client: client}

	return res, nil
}

func (cs *ClientToClientService) RegisterConnection(ctx context.Context, req *pb.RegisterConnectionRequest) (*pb.RegisterConnectionResponse, error) {
	if _, ok := cs.ActiveClients[req.FromId]; !ok {
		return nil, errors.New("Client not foud on active clients")
	}
	if _, ok := cs.ClientToClientConnections[req.FromId]; !ok {
		return nil, errors.New("Client not foud on cliet to client connections clients")
	}

	// In the future we are going to validate if the clients from toId exist in ActiveClients
	cs.ClientToClientConnections[req.FromId] = append(cs.ClientToClientConnections[req.FromId], req.GetToId()...)

	res := &pb.RegisterConnectionResponse{Response: "Connections created"}

	return res, nil
}

func (c *ClientToClientService) validateFirstMsg(msg *pb.Message) error {
	senderId := msg.GetSenderId()
	if senderId == 0 {
		return fmt.Errorf("SenderId cannot be 0.")
	}
	// Check if senderId exists on ClientStreams.
	if _, ok := c.ClientStreams[senderId]; ok {
		return fmt.Errorf("Client %v already has an open stream.", senderId)
	}
	// Check if the Simulator state is WaitingConnections and if the msgType is allowit
	if c.State.GetStateName() != utils.WaitingConnectionsStr {
		return fmt.Errorf("Simulator is not accepting new connections.")
	}
	// check if msgType is allow it
	if !c.State.IsMsgTypeAllowIt(msg) {
		return fmt.Errorf("MessageType: %v:%v is not allow it", msg.MessageType.String(), msg.MessageType.Number())
	}
	return nil
}

func (c *ClientToClientService) ClientToClientMessage(stream pb.Broadcast_ClientToClientMessageServer) error {

	// Read first msg
	msg, err := stream.Recv()
	if err != nil {
		return err
	}
	if err := c.validateFirstMsg(msg); err != nil {
		return err
	}

	senderId := msg.GetSenderId()

	conn := &models.Connection{Stream: stream, Outbox: make(chan *pb.Message, 10)}
	c.ClientStreams[senderId] = conn

	errCh := make(chan error, 1)
	// READ MESSAGE
	go func() {
		for {
			msg, err := stream.Recv() // This function block until we get a new msg
			if err != nil {
				errCh <- err
				return
			}
			if err = c.State.ReadMsg(msg, c.ClientStreams[senderId]); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Send msg to client
	go func() {
		clientConnection := c.ClientStreams[senderId]
		for msgQueue := range clientConnection.Outbox {
			err := stream.Send(msgQueue)
			if err != nil {
				errCh <- err
				return
			}
		}
	}()
	<-errCh
	return nil
}
