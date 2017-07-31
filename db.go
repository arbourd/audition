package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

// Database interfaces the bolt.DB database.
type Database interface {
	ListMessages() ([]*Message, error)
	GetMessage(id int) (*Message, error)
	CreateMessage(*Message) (*Message, error)
	DeleteMessage(id int) error
}

// Message is the model for the message object in the bolt.DB database.
type Message struct {
	ID          int    `json:"id"`
	Message     string `json:"message"`
	IsPalidrome bool   `json:"isPalindrome"`
	CreatedAt   string `json:"createdAt"`
}

// NotFoundError is a custom error returned when the item is not found in the database.
type NotFoundError struct {
	Key   string
	Value string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("could not find with %v: %v", e.Key, e.Value)
}

func (m *Message) determinePalidrome() bool {
	r := regexp.MustCompile(`[\W]+`)
	str := strings.ToLower(r.ReplaceAllString(m.Message, ""))

	for i := 0; i < len(str)/2; i++ {
		if str[i] != str[len(str)-i-1] {
			return false
		}
	}
	return true
}

// Store maintains the bolt.DB database's state.
type Store struct {
	db *bolt.DB
}

// InitializeStore creates or opens the bolt.DB and returns a Store struct containing it.
func InitializeStore(db *bolt.DB) (*Store, error) {
	tx, err := db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	_, err = tx.CreateBucketIfNotExists([]byte("messages"))
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Store{db: db}, err
}

// ListMessages gets all messages from the database
func (s *Store) ListMessages() ([]*Message, error) {
	var msgs = make([]*Message, 0)

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		err := b.ForEach(func(k, v []byte) error {
			var msg Message
			if err := json.Unmarshal(v, &msg); err != nil {
				return err
			}

			msgs = append(msgs, &msg)
			return nil
		})

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return msgs, nil
}

// GetMessage gets a single message by ID.
func (s *Store) GetMessage(id int) (*Message, error) {
	var msg Message

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		v := b.Get(itob(id))
		if len(v) == 0 {
			return &NotFoundError{
				Key:   "ID",
				Value: strconv.Itoa(id),
			}
		}

		if err := json.Unmarshal(v, &msg); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// CreateMessage stores a Message type obj in the database.
func (s *Store) CreateMessage(msg *Message) (*Message, error) {
	msg.CreatedAt = time.Now().Format(time.RFC3339)
	msg.IsPalidrome = msg.determinePalidrome()

	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))

		// Error cannot be returned inside of Update function, can be ignored
		id, _ := b.NextSequence()
		msg.ID = int(id)

		buf, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		return b.Put(itob(msg.ID), buf)
	})

	if err != nil {
		return nil, err
	}

	return msg, nil
}

// DeleteMessage removes a message from the database.
func (s *Store) DeleteMessage(id int) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("messages"))
		v := b.Get(itob(id))
		if len(v) == 0 {
			return &NotFoundError{
				Key:   "ID",
				Value: strconv.Itoa(id),
			}
		}

		return b.Delete(itob(id))
	})

	return err
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
