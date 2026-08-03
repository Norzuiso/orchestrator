package storage

import (
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Storage struct {
	Db *bolt.DB
}

func (s *Storage) OpenDb() {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{Timeout: 5 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	s.Db = db
	defer db.Close()
}
