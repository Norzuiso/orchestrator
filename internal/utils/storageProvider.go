package utils

import (
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
)

type StorageProvider interface {
	ActiveClientsSave(client *pb.Client) (*pb.Client, error)
	ActiveClientsGet(id int64) (*pb.Client, error)
	ClientsResponseSave(id int64, msg *pb.Message) error
	ClientsResponseGet(id int64) (*pb.Message, error)
	ClientToClientSave(id int64, conns *pb.ClientConnectionList) error
	ClientToClientGet(id int64) (*pb.ClientConnectionList, error)
}
