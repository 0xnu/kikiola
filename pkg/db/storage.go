package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/tidwall/buntdb"
)

type Storage struct {
	db    *buntdb.DB
	mutex sync.RWMutex
}

func NewStorage(dbPath string) (*Storage, error) {
	err := os.MkdirAll(filepath.Dir(dbPath), os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	db, err := buntdb.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	storage := &Storage{
		db: db,
	}

	return storage, nil
}

func (s *Storage) Insert(vector *Vector) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	data, err := json.Marshal(vector)
	if err != nil {
		return fmt.Errorf("failed to marshal vector: %v", err)
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(vector.ID, string(data), nil)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to insert vector: %v", err)
	}

	return nil
}

func (s *Storage) Get(id string) (*Vector, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var data string
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			return err
		}
		data = val
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get vector: %v", err)
	}

	var vector Vector
	err = json.Unmarshal([]byte(data), &vector)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal vector: %v", err)
	}

	return &vector, nil
}

func (s *Storage) Delete(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to delete vector: %v", err)
	}

	return nil
}

func (s *Storage) GetAll() ([]*Vector, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var vectors []*Vector
	err := s.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			var vector Vector
			err := json.Unmarshal([]byte(value), &vector)
			if err != nil {
				return false
			}
			vectors = append(vectors, &vector)
			return true
		})
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get all vectors: %v", err)
	}

	return vectors, nil
}

func (s *Storage) Update(id string, metadata map[string]string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var vector Vector
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return fmt.Errorf("vector not found")
			}
			return fmt.Errorf("failed to get vector: %v", err)
		}

		err = json.Unmarshal([]byte(val), &vector)
		if err != nil {
			return fmt.Errorf("failed to unmarshal vector: %v", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	for key, value := range metadata {
		vector.Metadata[key] = value
	}

	data, err := json.Marshal(vector)
	if err != nil {
		return fmt.Errorf("failed to marshal vector: %v", err)
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(id, string(data), nil)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to update vector: %v", err)
	}

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
