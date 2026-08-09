package storage

import (
	"database/sql"
	"time"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"
	_ "modernc.org/sqlite"
)

type LogStorage struct {
	Db *sql.DB
}

func (s *LogStorage) CloseDb() error {
	if s.Db != nil {
		return s.Db.Close()
	}

	return nil
}

func (s *LogStorage) OpenDb(path string) error {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return err
	}

	if _, err := db.Exec("PRAGMA busy_timeout=5000;"); err != nil {
		return err
	}
	db.SetMaxOpenConns(1)
	s.Db = db

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS  logs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp DATATIME NOT NULL,
		direction TEXT,
		state TEXT,
		sender_id INTEGER,
		content TEXT,
		message_type  TEXT,
		epoch INTEGER,
		seed INTEGER,
		attributes TEXT
	)
	`)
	return err
}

func (s *LogStorage) LogMessage(direction string, msg *pb.Message, state string) error {
	_, err := s.Db.Exec(`
		INSERT INTO logs (timestamp, direction, sender_id, state, content, message_type, epoch, seed, attributes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, time.Now(), direction, msg.GetSenderId(), state, msg.GetContent(), msg.MessageType.String(), msg.GetEpoch(), msg.GetSeed(), msg.String())
	return err
}
