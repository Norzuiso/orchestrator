package register

import (
	"context"
	"log"
	"sync"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type Connection struct {
	pb.UnimplementedBroadcastServer
	stream             pb.Broadcast_CreateStreamServer
	id                 string
	error              chan error
	ServiceConnections []string
}

type Pool struct {
	pb.UnimplementedBroadcastServer
	Connection map[string]*Connection
}

func (p *Pool) CreateStream(pconn *pb.Connect, stream pb.Broadcast_CreateStreamServer) error {
	serviceConnections := pconn.UserConnectionIds

	conn := &Connection{
		stream:             stream,
		error:              make(chan error),
		ServiceConnections: serviceConnections,
		id:                 pconn.User.Id,
	}
	p.Connection[pconn.User.Id] = conn
	return <-conn.error
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
	sendConn := p.Connection[msg.SenderId]

	for _, connId := range sendConn.ServiceConnections {
		conn := p.Connection[connId]
		SendMsg(connId, msg.SenderId, msg, conn)
	}
	return &pb.Close{}, nil
}
