package register

import (
	"context"
	"io"
	"log"
	"sync"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type Connection struct {
	pb.UnimplementedBroadcastServer
	stream pb.Broadcast_PairToPairMessageServer
	id     string
	error  chan error
}

var ServiceConnections map[string][]string = make(map[string][]string)

type Pool struct {
	pb.UnimplementedBroadcastServer
	Connection map[string]*Connection
}

func (p *Pool) PairToPairMessage(stream pb.Broadcast_PairToPairMessageServer) error {
	wait := sync.WaitGroup{}
	done := make(chan int)
	wait.Add(1)
	go func(stream pb.Broadcast_PairToPairMessageServer) {
		in, err := p.AddConnetion(stream)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		recivers := ServiceConnections[in.SenderId]
		log.Printf("Inicio llamada de: %v\n Envia a: %v\n", in.SenderId, recivers)
		for _, reciver := range recivers {
			if _, ok := p.Connection[reciver]; ok {
				p.Connection[reciver].stream.Send(in)
			} else {
				stream.Send(&pb.Message{
					SenderId: "Server",
					Content:  reciver + "is not created",
				})
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

func (p *Pool) AddConnetion(stream pb.Broadcast_PairToPairMessageServer) (*pb.Message, error) {
	in, err := stream.Recv()
	if err == io.EOF {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	if _, ok := p.Connection[in.SenderId]; ok {
		return nil, err
	}
	p.Connection[in.SenderId] = &Connection{
		stream: stream,
		error:  make(chan error),
		id:     in.SenderId,
	}
	return in, nil
}

func (p *Pool) RegisterConnection(ctx context.Context, pconn *pb.Connect) (*pb.ConnectResponse, error) {
	ServiceConnections[pconn.User.Id] = pconn.UserConnectionIds
	log.Println(ServiceConnections)

	return &pb.ConnectResponse{ResponseMsg: "Todo bien, ya quedo registrado"}, nil

}

func (p *Pool) CreateStream(pconn *pb.Connect, stream pb.Broadcast_CreateStreamServer) error {

	conn := &Connection{
		//	stream: stream,
		error: make(chan error),
		id:    pconn.User.Id,
	}
	p.Connection[pconn.User.Id] = conn
	return <-conn.error
}

func (p *Pool) RegisterBroadcastServer(stream *pb.Broadcast_PairToPairMessageServer) {

}

func SendMsg(ReceiverId string, SenderId string, msg *pb.Message, conn *Connection) (*pb.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)
	log.Printf("Se esta mandando a: %v desde: %v\n", ReceiverId, SenderId)
	wait.Add(1)
	go func(msg *pb.Message, conn *Connection) {
		defer wait.Done()
		err := conn.stream.SendMsg(msg)
		if err != nil {
			log.Printf("Error with Stream: %v - Error: %v\n", conn.stream, err)
			conn.error <- err
		}
	}(msg, conn)
	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
	return &pb.Close{}, nil
}

func (p *Pool) BroadcastMessage(ctx context.Context, msg *pb.Message) (*pb.Close, error) {
	for _, connId := range ServiceConnections[msg.SenderId] {
		conn := p.Connection[connId]
		SendMsg(connId, msg.SenderId, msg, conn)
	}
	return &pb.Close{}, nil
}
