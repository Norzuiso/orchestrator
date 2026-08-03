package storage

import (
	"fmt"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"

	"github.com/Norzuiso/orchestrator/internal/utils"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

func (s *Storage) ClientToClientSave(id int64, conns *pb.ClientConnectionList) error {
	return s.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(utils.BucketClientToClient))
		if err != nil {
			return err
		}

		key := utils.Int64ToBytes(id)
		value, err := proto.Marshal(conns) // ahora sí es UN solo mensaje, serializable normal
		if err != nil {
			return err
		}
		return bucket.Put(key, value)
	})

}

func (s *Storage) ClientToClientGet(id int64) (*pb.ClientConnectionList, error) {
	conns := &pb.ClientConnectionList{}
	err := s.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(utils.BucketClientToClient))
		if bucket == nil {
			return fmt.Errorf("Bucket does not exist")
		}

		key := utils.Int64ToBytes(id)
		value := bucket.Get(key)

		if value == nil {
			return fmt.Errorf("Client not found")
		}

		if err := proto.Unmarshal(value, conns); err != nil {
			return err
		}

		return nil
	})
	return conns, err
}
