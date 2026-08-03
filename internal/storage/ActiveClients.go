package storage

import (
	"fmt"

	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func (s *Storage) ActiveClientsSave(client *pb.Client) error {
	return s.Db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(utils.BucketClients))
		if err != nil {
			return err
		}
		key := utils.Int64ToBytes(client.Id)
		value, err := proto.Marshal(client)
		if err != nil {
			return err
		}
		return bucket.Put(key, value)
	})
}

func (s *Storage) ActiveClientsGet(id int64) (*pb.Client, error) {
	client := &pb.Client{}
	err := s.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketClients))
		if bucket == nil {
			return fmt.Errorf("Bucket does not exist")
		}

		key := utils.Int64ToBytes(id)
		value := bucket.Get(key)

		if value == nil {
			return fmt.Errorf("Client not found")
		}

		if err := proto.Unmarshal(value, client); err != nil {
			return err
		}

		return nil
	})
	return client, err
}
