package index

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/0xnu/kikiola/pkg/db"
)

// Index represents an index for efficient vector retrieval.
type Index struct {
	storage *db.Storage
	index   map[string][]*db.Vector
	mutex   sync.RWMutex
}

// NewIndex creates a new instance of Index.
func NewIndex(storage *db.Storage) (*Index, error) {
	index := &Index{
		storage: storage,
		index:   make(map[string][]*db.Vector),
	}

	err := index.buildIndex()
	if err != nil {
		return nil, err
	}

	return index, nil
}

// Insert inserts a vector into the index.
func (i *Index) Insert(vector *db.Vector) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// Insert the vector into the storage
	err := i.storage.Insert(vector)
	if err != nil {
		return err
	}

	// Update the index
	for _, value := range vector.Embedding {
		key := i.getKey(value)
		i.index[key] = append(i.index[key], vector)
	}

	return nil
}

// Delete removes a vector from the index.
func (i *Index) Delete(id string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	// Get the vector from the storage
	vector, err := i.storage.Get(id)
	if err != nil {
		return err
	}

	// Remove the vector from the index
	for _, value := range vector.Embedding {
		key := i.getKey(value)
		i.index[key] = removeVector(i.index[key], vector)
	}

	// Delete the vector from the storage
	err = i.storage.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

// Search searches for the nearest neighbors of a given vector using cosine similarity.
func (i *Index) Search(vector *db.Vector, k int) ([]*db.Vector, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	if k <= 0 {
		return nil, errors.New("invalid value of k")
	}

	var candidates []*db.Vector
	for _, value := range vector.Embedding {
		key := i.getKey(value)
		candidates = append(candidates, i.index[key]...)
	}

	// Sort candidates by cosine similarity to the query vector
	sort.Slice(candidates, func(i, j int) bool {
		simI, _ := cosineSimilarity(*vector, *candidates[i])
		simJ, _ := cosineSimilarity(*vector, *candidates[j])
		return simI > simJ
	})

	// Return the top k candidates
	if len(candidates) > k {
		candidates = candidates[:k]
	}

	return candidates, nil
}

// buildIndex builds the index from the vectors in the storage.
func (i *Index) buildIndex() error {
	vectors, err := i.storage.GetAll()
	if err != nil {
		return err
	}

	for _, vector := range vectors {
		for _, value := range vector.Embedding {
			key := i.getKey(value)
			i.index[key] = append(i.index[key], vector)
		}
	}

	return nil
}

// getKey generates a key for a given value.
func (i *Index) getKey(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// removeVector removes a vector from a slice of vectors.
func removeVector(vectors []*db.Vector, vector *db.Vector) []*db.Vector {
	for i, v := range vectors {
		if v.ID == vector.ID {
			return append(vectors[:i], vectors[i+1:]...)
		}
	}
	return vectors
}

// cosineSimilarity calculates the cosine similarity between two vectors.
func cosineSimilarity(v1, v2 db.Vector) (float64, error) {
	if len(v1.Embedding) != len(v2.Embedding) {
		return 0, errors.New("embedding dimensions mismatch")
	}

	dotProduct := 0.0
	normV1 := 0.0
	normV2 := 0.0

	for i := range v1.Embedding {
		dotProduct += v1.Embedding[i] * v2.Embedding[i]
		normV1 += v1.Embedding[i] * v1.Embedding[i]
		normV2 += v2.Embedding[i] * v2.Embedding[i]
	}

	if normV1 == 0 || normV2 == 0 {
		return 0, nil
	}

	return dotProduct / (math.Sqrt(normV1) * math.Sqrt(normV2)), nil
}
