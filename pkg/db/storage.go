package db

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/tidwall/buntdb"
)

type DistributedStorage struct {
	nodes []*Storage
	mutex sync.RWMutex
}

type Object struct {
	ID       string            `json:"id"`
	Object   []byte            `json:"object"`
	Metadata map[string]string `json:"metadata"`
}

var ErrNotFound = errors.New("object not found")
var ErrObjectNotFound = errors.New("object not found")
var ErrVectorNotFound = errors.New("vector not found")

func NewDistributedStorage(nodeAddresses []string) (*DistributedStorage, error) {
	var nodes []*Storage

	for _, address := range nodeAddresses {
		dbPath := fmt.Sprintf("data/node_%s.db", address)
		storage, err := NewStorage(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create storage for node %s: %v", address, err)
		}
		nodes = append(nodes, storage)
	}

	return &DistributedStorage{
		nodes: nodes,
	}, nil
}

func (ds *DistributedStorage) InsertVector(vector *Vector) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	nodeIndex := ds.getNodeIndex(vector.ID)
	return ds.nodes[nodeIndex].InsertVector(vector)
}

func (ds *DistributedStorage) GetVector(id string) (*Vector, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	nodeIndex := ds.getNodeIndex(id)
	return ds.nodes[nodeIndex].GetVector(id)
}

func (ds *DistributedStorage) DeleteVector(id string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	nodeIndex := ds.getNodeIndex(id)
	return ds.nodes[nodeIndex].DeleteVector(id)
}

func (ds *DistributedStorage) GetAllVectors() ([]*Vector, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	var vectors []*Vector
	for _, node := range ds.nodes {
		nodeVectors, err := node.GetAllVectors()
		if err != nil {
			return nil, fmt.Errorf("failed to get vectors from node: %v", err)
		}
		vectors = append(vectors, nodeVectors...)
	}

	return vectors, nil
}

func (ds *DistributedStorage) UpdateVectorMetadata(id string, metadata map[string]string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	nodeIndex := ds.getNodeIndex(id)
	return ds.nodes[nodeIndex].UpdateVectorMetadata(id, metadata)
}

func (ds *DistributedStorage) InsertObject(object *Object) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	nodeIndex := ds.getNodeIndex(object.ID)
	return ds.nodes[nodeIndex].InsertObject(object)
}

func (ds *DistributedStorage) GetObject(id string) (*Object, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	nodeIndex := ds.getNodeIndex(id)
	return ds.nodes[nodeIndex].GetObject(id)
}

func (ds *DistributedStorage) DeleteObject(id string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	nodeIndex := ds.getNodeIndex(id)
	return ds.nodes[nodeIndex].DeleteObject(id)
}

func (ds *DistributedStorage) GetAllObjects() ([]*Object, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	var objects []*Object
	for _, node := range ds.nodes {
		nodeObjects, err := node.GetAllObjects()
		if err != nil {
			return nil, fmt.Errorf("failed to get objects from node: %v", err)
		}
		objects = append(objects, nodeObjects...)
	}

	return objects, nil
}

func (ds *DistributedStorage) UpdateObjectMetadata(id string, metadata map[string]string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	nodeIndex := ds.getNodeIndex(id)
	return ds.nodes[nodeIndex].UpdateObjectMetadata(id, metadata)
}

func (ds *DistributedStorage) Close() error {
	for _, node := range ds.nodes {
		err := node.Close()
		if err != nil {
			return fmt.Errorf("failed to close node storage: %v", err)
		}
	}
	return nil
}

func (ds *DistributedStorage) getNodeIndex(id string) int {
	hash := sha256.New()

	hash.Write([]byte(id))

	hashValue := hash.Sum(nil)

	hashUint64 := binary.BigEndian.Uint64(hashValue)

	bestNodeIndex := 0
	minDistance := uint64(math.MaxUint64)

	for i := range ds.nodes {
		nodeIndexBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(nodeIndexBytes, uint64(i))

		nodeHash := sha256.New()

		nodeHash.Write(nodeIndexBytes)

		nodeHashValue := nodeHash.Sum(nil)

		nodeHashUint64 := binary.BigEndian.Uint64(nodeHashValue)

		distance := hashUint64 ^ nodeHashUint64

		if distance < minDistance {
			bestNodeIndex = i
			minDistance = distance
		}
	}

	return bestNodeIndex
}

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

func (s *Storage) InsertVector(vector *Vector) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	serializedData, err := json.Marshal(vector)
	if err != nil {
		return fmt.Errorf("failed to marshal vector: %v", err)
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(vector.ID, string(serializedData), nil)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to insert vector: %v", err)
	}

	return nil
}

func (s *Storage) GetVector(id string) (*Vector, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var data string
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return ErrVectorNotFound
			}
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

func (s *Storage) DeleteVector(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return ErrVectorNotFound
			}
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete vector: %v", err)
	}

	return nil
}

func (s *Storage) GetAllVectors() ([]*Vector, error) {
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

func (s *Storage) UpdateVectorMetadata(id string, metadata map[string]string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var vector Vector
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return ErrVectorNotFound
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

func (s *Storage) InsertObject(object *Object) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	serializedData, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %v", err)
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(object.ID, string(serializedData), nil)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to insert object: %v", err)
	}

	return nil
}

func (s *Storage) GetObject(id string) (*Object, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var data string
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return ErrObjectNotFound
			}
			return err
		}
		data = val
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}

	var object Object
	err = json.Unmarshal([]byte(data), &object)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal object: %v", err)
	}

	return &object, nil
}

func (s *Storage) DeleteObject(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	err := s.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(id)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return ErrObjectNotFound
			}
			return err
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	return nil
}

func (s *Storage) GetAllObjects() ([]*Object, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var objects []*Object
	err := s.db.View(func(tx *buntdb.Tx) error {
		err := tx.Ascend("", func(key, value string) bool {
			var object Object
			err := json.Unmarshal([]byte(value), &object)
			if err != nil {
				return false
			}
			objects = append(objects, &object)
			return true
		})
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get all objects: %v", err)
	}

	return objects, nil
}

func (s *Storage) UpdateObjectMetadata(id string, metadata map[string]string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	var object Object
	err := s.db.View(func(tx *buntdb.Tx) error {
		val, err := tx.Get(id)
		if err != nil {
			if err == buntdb.ErrNotFound {
				return ErrObjectNotFound
			}
			return fmt.Errorf("failed to get object: %v", err)
		}

		err = json.Unmarshal([]byte(val), &object)
		if err != nil {
			return fmt.Errorf("failed to unmarshal object: %v", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	for key, value := range metadata {
		object.Metadata[key] = value
	}

	data, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("failed to marshal object: %v", err)
	}

	err = s.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(id, string(data), nil)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to update object: %v", err)
	}

	return nil
}

func (s *Storage) GetVectors(ids []string) ([]*Vector, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var results []string
	err := s.db.View(func(tx *buntdb.Tx) error {
		for _, id := range ids {
			val, err := tx.Get(id)
			if err != nil {
				if err == buntdb.ErrNotFound {
					continue
				}
				return fmt.Errorf("failed to get vector with ID %s: %v", id, err)
			}
			results = append(results, val)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %v", err)
	}

	vectors := make([]*Vector, 0, len(results))
	for _, data := range results {
		var vector Vector
		err := json.Unmarshal([]byte(data), &vector)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal vector: %v", err)
		}
		vectors = append(vectors, &vector)
	}

	return vectors, nil
}

func (ds *DistributedStorage) GetVectors(ids []string) ([]*Vector, error) {
	ds.mutex.RLock()
	defer ds.mutex.RUnlock()

	var vectors []*Vector
	for _, id := range ids {
		nodeIndex := ds.getNodeIndex(id)
		vector, err := ds.nodes[nodeIndex].GetVector(id)
		if err != nil {
			if errors.Is(err, ErrVectorNotFound) {
				continue
			}
			return nil, fmt.Errorf("failed to get vector with ID %s: %v", id, err)
		}
		vectors = append(vectors, vector)
	}

	return vectors, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
