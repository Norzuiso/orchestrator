package servers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/orchestrator"
	"github.com/Norzuiso/orchestrator/internal/sse"
	"github.com/Norzuiso/orchestrator/internal/storage"
	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type ClientToClientService struct {
	pb.UnimplementedBroadcastServer

	ClientStreams map[int64]*models.Connection // TODO use mux
	Orchestrator  *orchestrator.Orchestrator

	StateProvider         utils.StateProvider
	StorageProvider       utils.StorageProvider
	LogStorage            *storage.LogStorage
	LogsBroadcaster       *sse.LogsBroadcaster
	OpenStreamBroadcaster *sse.OpenStreamBroadcaster
}

func NewClientToClientService(o *orchestrator.Orchestrator) *ClientToClientService {
	return &ClientToClientService{
		ClientStreams: make(map[int64]*models.Connection),
		Orchestrator:  o,
		LogStorage:    &storage.LogStorage{},
	}
}

func (c *ClientToClientService) ConnectClient(ctx context.Context, req *pb.ConnectionRequest) (*pb.ConnectionResponse, error) {
	client := req.GetClient()

	if client.GetName() == "" {
		client.Name = strconv.FormatInt(client.Id, 10)
	}

	if client.GetSeed() == 0 {
		client.Seed = c.Orchestrator.GetClientSeed(client.Id)
	}
	clientResponse, err := c.StorageProvider.ActiveClientsSave(client)
	if err != nil {
		return nil, err
	}
	res := &pb.ConnectionResponse{Client: clientResponse}
	return res, nil
}

func (c *ClientToClientService) RegisterConnection(ctx context.Context, req *pb.RegisterConnectionRequest) (*pb.RegisterConnectionResponse, error) {

	if _, err := c.StorageProvider.ActiveClientsGet(req.FromId); err != nil {
		return nil, err
	}

	err := c.StorageProvider.ClientToClientSave(req.GetFromId(), req.To)
	if err != nil {
		return nil, err
	}
	return &pb.RegisterConnectionResponse{Response: "Connections created"}, nil
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
		return fmt.Errorf("MessageType: %v:%v is not allow it. Current state: %v", msg.MessageType.String(), msg.MessageType.Number(), c.StateProvider.GetCurrentState().GetStateName())
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

			for {
				if err := c.validateMsgTypeOnState(msg); err == nil {
					break
				}
				select {
				case <-c.StateProvider.StateChangedChan():
					continue
				case <-stream.Context().Done():
					return
				}
			}
			c.LogsBroadcaster.Publish(fmt.Sprint("Read: ", msg))
			err = c.LogStorage.LogMessage("Read", msg, c.StateProvider.GetCurrentState().GetStateName())
			if err != nil {
				log.Println(err)
			}
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
			c.LogsBroadcaster.Publish(fmt.Sprint("Send: ", item.Msg))
			err = c.LogStorage.LogMessage("Send", item.Msg, c.StateProvider.GetCurrentState().GetStateName())

			if err != nil {
				log.Println(err)
			}
		}
	}()

	ctx := stream.Context()
	select {
	// Detect close connection
	case <-ctx.Done():
		c.LogsBroadcaster.Publish(fmt.Sprint("Stream closed: ", &pb.Message{SenderId: senderId, Content: ctx.Err().Error()}))
		err = c.LogStorage.LogMessage("Stream closed", &pb.Message{SenderId: senderId, Content: ctx.Err().Error()}, c.StateProvider.GetCurrentState().GetStateName())
		if err != nil {
			log.Println(err)
		}
		delete(c.ClientStreams, senderId)
	case <-c.ClientStreams[senderId].ErrCh:
		delete(c.ClientStreams, senderId)
		log.Printf("Sender: %v error. ERR: %v", senderId, ctx.Err())

	}
	c.OpenStreamBroadcaster.Publish(senderId)                // Send senderID to get client with not open stream
	c.OpenStreamBroadcaster.RemoveWhere(func(i int64) bool { // Remove the id from history to prevent multiple appear when there are not open streams for client
		return senderId == i
	})

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

	msgResponse := &pb.Message{
		SenderId:    0,
		MessageType: pb.MessageType_MESSAGE_TYPE_DEFAULT,
		Content:     "Connection created",
		Epoch:       c.Orchestrator.Epoch,
		Seed:        c.Orchestrator.GetClientSeed(senderId),
	}
	stream.Send(msgResponse)

	c.ClientStreams[senderId].ErrCh = make(chan error, 2)
	err = c.LogStorage.LogMessage("Stream open", msgResponse, c.StateProvider.GetCurrentState().GetStateName())
	if err != nil {
		log.Println(err)
	}
	c.LogsBroadcaster.Publish(fmt.Sprint("Stream open: ", msg))
	c.OpenStreamBroadcaster.Publish(senderId)
	return senderId, conn, nil
}
