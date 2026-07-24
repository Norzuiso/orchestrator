package servers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"

	"github.com/Norzuiso/orchestrator/internal/orchestrator"
	register "github.com/Norzuiso/orchestrator/internal/register"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"google.golang.org/grpc/peer"
)

type Connection struct {
	Stream     pb.Broadcast_ClientToClientMessageServer
	MsgHandler *register.MsgHandler
	error      chan error
}

type ClientToClientService struct {
	pb.UnimplementedBroadcastServer
	Seed                      int64
	ActiveClients             map[int64]*pb.Client
	ClientStreams             map[int64]*Connection
	ClientToClientConnections map[int64][]int64
	Orchestrator              *orchestrator.Orchestrator
}

func NewClientToClientService(seed int64, orchestrator *orchestrator.Orchestrator) *ClientToClientService {
	return &ClientToClientService{
		Seed:                      seed,
		ActiveClients:             make(map[int64]*pb.Client),
		ClientStreams:             make(map[int64]*Connection),
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

func (cs *ClientToClientService) ClientToClientMessage(stream pb.Broadcast_ClientToClientMessageServer) error {
	wait := sync.WaitGroup{}
	done := make(chan int)
	wait.Add(1)
	go func(stream pb.Broadcast_ClientToClientMessageServer) {
		for {

			/*
				For each open stream (client)
					1. Read message
					2. Get MessageType
					3. Check if MessageType is allow it on the current ORQUESTRATOR phase
					4. Apply the phase behivor on the Message

			*/

			// READ MESSAGE
			msg, err := stream.Recv()
			if err == io.EOF || err != nil {
				return
			}
			msgHandler := register.NewMessageHandler(msg, stream)
			msgType := msgHandler.GetMessageType()

			if !cs.Orchestrator.GetPhase().IsMsgTypeAllowIt(msgType) {
				msgHandler.SendPhaseError(fmt.Sprintf("Current orchestrator phase: %v only allows msg type: %v", cs.Orchestrator.GetPhase().String(), cs.Orchestrator.GetPhase().GetAllowMsgTypeStr()))
				continue
			}

			/*===================	CHECK	====================*/
			//	This code get the address from the client - We could use this to store information from the client
			log.Println(cs.Orchestrator.CurrentPhase.GetName())
			if cs.Orchestrator.GetPhase().IsWaitingConnection() {
				p, _ := peer.FromContext(stream.Context())
				msgHandler.SendMsgContent(p.Addr.String())
			}
			/*==================================================*/

			sender := msgHandler.SenderId
			if _, ok := cs.ClientStreams[sender]; !ok {
				cs.ClientStreams[sender] = &Connection{Stream: stream, MsgHandler: msgHandler}
			}

			if _, ok := cs.ClientToClientConnections[sender]; !ok {
				msgHandler.SendError(errors.New("Client doesn't have connections"))
			}
			for _, toId := range cs.ClientToClientConnections[sender] {
				if _, ok := cs.ClientStreams[toId]; !ok {
					msgHandler.SendError(fmt.Errorf("Client: %d is not an active client or doesnt exist", toId))
				} else {
					cs.ClientStreams[toId].MsgHandler.SendMsgContent(msgHandler.Content)
				}

			}

		}
	}(stream)
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return nil
}
