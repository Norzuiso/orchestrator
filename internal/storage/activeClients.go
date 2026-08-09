package storage

import (
	"fmt"

	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"go.etcd.io/bbolt"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func (s *Storage) ActiveClientsSave(client *pb.Client) (*pb.Client, error) {
	err := s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(utils.BucketClients))
		if err != nil {
			return err
		}

		if client.GetId() == 0 {
			id, _ := bucket.NextSequence()
			client.Id = int64(id)
		}

		key := utils.Int64ToBytes(client.Id)
		value, err := proto.Marshal(client)
		if err != nil {
			return err
		}
		return bucket.Put(key, value)
	})
	return client, err
}

func (s *Storage) ActiveClientsGet(id int64) (*pb.Client, error) {
	client := &pb.Client{}
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketClients))
		if bucket == nil {
			return fmt.Errorf("Bucket does not exist")
		}

		key := utils.Int64ToBytes(id)
		value := bucket.Get(key)

		if value == nil {
			return fmt.Errorf("BucketClients not found")
		}

		if err := proto.Unmarshal(value, client); err != nil {
			return err
		}

		return nil
	})
	return client, err
}

func (s *Storage) GetAllActiveClients() ([]*pb.Client, error) {
	clients := make([]*pb.Client, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketClients))
		if bucket == nil {
			return fmt.Errorf("Bucket does not exist")
		}

		bucket.ForEach(func(k, v []byte) error {
			client := &pb.Client{}

			if err := proto.Unmarshal(v, client); err != nil {
				return err
			}
			clients = append(clients, client)
			return nil
		})

		return nil
	})

	return clients, err
}
