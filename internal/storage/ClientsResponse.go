package storage

import (
	"fmt"

	"github.com/Norzuiso/orchestrator/internal/utils"
	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	"go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func (s *Storage) ClientsResponseSave(id int64, msg *pb.Message) error {
	return s.Db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(utils.BucketClientResponse))
		if err != nil {
			return err
		}

		key := utils.Int64ToBytes(id)

		value, err := proto.Marshal(msg)
		if err != nil {
			return err
		}

		return bucket.Put(key, value)
	})
}

func (s *Storage) ClientsResponseGet(id int64) (*pb.Message, error) {
	msg := &pb.Message{}

	err := s.Db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketClientResponse))
		if bucket == nil {
			return fmt.Errorf("Bucket does not exist")
		}

		key := utils.Int64ToBytes(id)
		value := bucket.Get(key)

		if value == nil {
			return fmt.Errorf("No event expected from client: %v", id)
		}

		if err := proto.Unmarshal(value, msg); err != nil {
			return err
		}

		return nil
	})

	return msg, err
}
