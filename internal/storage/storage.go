package storage

import (
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Storage struct {
	db *bolt.DB
}

func (s *Storage) OpenDb(path string) {
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	s.db = db
}

func (s *Storage) CloseDb() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
