package servers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/orchestrator"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type ClientToClientService struct {
	pb.UnimplementedBroadcastServer

	ClientStreams map[int64]*models.Connection
	Orchestrator  *orchestrator.Orchestrator

	StateProvider utils.StateProvider
}

func NewClientToClientService(o *orchestrator.Orchestrator) *ClientToClientService {
	return &ClientToClientService{
		ClientStreams: make(map[int64]*models.Connection),
		Orchestrator:  o,
	}
}

func (c *ClientToClientService) ConnectClient(ctx context.Context, req *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	client := req.GetClient()

	if client.GetId() == 0 {
		client.Id = int64(len(c.Orchestrator.ActiveClients) + 1)
	}

	if client.GetName() == "" {
		client.Name = strconv.FormatInt(client.Id, 10)
	}

	if client.GetSeed() == 0 {
		client.Seed = c.Orchestrator.GetClientSeed(client.Id)
	}

	c.Orchestrator.ActiveClients[client.Id] = client
	c.Orchestrator.ClientToClientConnection[client.Id] = make([]*pb.ClientConnection, 0)

	res := &pb.ConnectionResponse{Client: client}

	return res, nil
}

func (c *ClientToClientService) RegisterConnection(ctx context.Context, req *pb.RegisterConnectionRequest) (*pb.RegisterConnectionResponse, error) {
	if _, ok := c.Orchestrator.ActiveClients[req.FromId]; !ok {
		return nil, errors.New("Client not foud on active clients")
	}
	if _, ok := c.Orchestrator.ClientToClientConnection[req.FromId]; !ok {
		return nil, errors.New("Client not foud on cliet to client connections clients")
	}

	c.Orchestrator.ClientToClientConnection[req.FromId] = append(c.Orchestrator.ClientToClientConnection[req.FromId], req.To...)
	log.Println(c.Orchestrator.ClientToClientConnection)
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
	if c.StateProvider.GetCurrentState().GetStateName() != utils.WaitingConnectionsStr {
		return fmt.Errorf("Simulator is not accepting new connections.")
	}
	return c.validateMsgTypeOnState(msg)
}

func (c *ClientToClientService) validateMsgTypeOnState(msg *pb.Message) error {
	// check if msgType is allow it
	if !c.StateProvider.GetCurrentState().IsMsgTypeAllowIt(msg) {
		return fmt.Errorf("MessageType: %v:%v is not allow it", msg.MessageType.String(), msg.MessageType.Number())
	}
	return nil
}

func (c *ClientToClientService) ClientToClientMessage(stream pb.Broadcast_ClientToClientMessageServer) error {

	// Read first msg
	senderId, conn, err := openStream(stream, c)
	if err != nil {
		return err
	}

	// READ MESSAGE
	go func() {
		for {
			msg, err := stream.Recv() // This function block until we get a new msg
			if err != nil {
				return
			}

			if err := c.validateMsgTypeOnState(msg); err != nil {
				done := make(chan error, 1)
				conn.Outbox <- models.OutboxItem{Msg: utils.BuildPhaseErrorMsg(msg, err), Done: done}
				<-done
				continue
			}
			log.Println("Read: ", msg.String())

			if err = c.StateProvider.GetCurrentState().ReadMsg(msg, c.ClientStreams[senderId]); err != nil {
				done := make(chan error, 1)
				conn.Outbox <- models.OutboxItem{Msg: utils.BuildPhaseErrorMsg(msg, err), Done: done}
				<-done
			}
		}
	}()

	// Send msg to client
	go func() {
		clientConnection := c.ClientStreams[senderId]
		for item := range clientConnection.Outbox {
			item.Msg.Epoch = c.Orchestrator.Epoch
			err := stream.Send(item.Msg)
			item.Done <- err
			if err != nil {
				c.ClientStreams[senderId].ErrCh <- err
				return
			}
			log.Println("Send: ", item.Msg.String())
		}
	}()
	<-c.ClientStreams[senderId].ErrCh
	delete(c.ClientStreams, senderId)

	return nil
}

func openStream(stream pb.Broadcast_ClientToClientMessageServer, c *ClientToClientService) (int64, *models.Connection, error) {
	msg, err := stream.Recv()
	if err != nil {
		return 0, nil, err
	}
	if err := c.validateFirstMsg(msg); err != nil {
		return 0, nil, err
	}

	senderId := msg.GetSenderId()

	conn := &models.Connection{Stream: stream, Outbox: make(chan models.OutboxItem)}
	c.ClientStreams[senderId] = conn

	stream.Send(&pb.Message{
		SenderId:    0,
		MessageType: pb.MessageType_MESSAGE_TYPE_DEFAULT,
		Content:     "Connection created",
		Epoch:       c.Orchestrator.Epoch,
		Seed:        c.Orchestrator.GetClientSeed(senderId),
	})
	log.Println(msg.String())

	c.ClientStreams[senderId].ErrCh = make(chan error, 2)
	return senderId, conn, nil
}
