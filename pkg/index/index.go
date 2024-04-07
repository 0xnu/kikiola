package index

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/agnivade/levenshtein"
)

type Index struct {
	storage *db.DistributedStorage
	index   map[string][]*db.Vector
	mutex   sync.RWMutex
}

func NewIndex(storage *db.DistributedStorage) (*Index, error) {
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

func (i *Index) Insert(vector *db.Vector) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	err := i.storage.InsertVector(vector)
	if err != nil {
		return err
	}

	for _, value := range vector.Embedding {
		key := i.getKey(value)
		i.index[key] = append(i.index[key], vector)
	}

	return nil
}

func (i *Index) Delete(id string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	vector, err := i.storage.GetVector(id)
	if err != nil {
		return err
	}

	for _, value := range vector.Embedding {
		key := i.getKey(value)
		i.index[key] = removeVector(i.index[key], vector)
	}

	err = i.storage.DeleteVector(id)
	if err != nil {
		return err
	}

	return nil
}

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

	sort.Slice(candidates, func(i, j int) bool {
		simI, _ := cosineSimilarity(*vector, *candidates[i])
		simJ, _ := cosineSimilarity(*vector, *candidates[j])
		return simI > simJ
	})

	uniqueCandidates := make([]*db.Vector, 0, len(candidates))
	seenIDs := make(map[string]bool)

	for _, candidate := range candidates {
		if !seenIDs[candidate.ID] {
			uniqueCandidates = append(uniqueCandidates, candidate)
			seenIDs[candidate.ID] = true
		}
	}

	if len(uniqueCandidates) > k {
		uniqueCandidates = uniqueCandidates[:k]
	}

	results, err := i.storage.GetVectors(getIDs(uniqueCandidates))
	if err != nil {
		return nil, fmt.Errorf("failed to get vectors: %v", err)
	}

	Rerank(results, vector.Text)

	if len(results) > k {
		results = results[:k]
	}

	return results, nil
}

func getIDs(vectors []*db.Vector) []string {
	ids := make([]string, len(vectors))
	for i, vector := range vectors {
		ids[i] = vector.ID
	}
	return ids
}

func (i *Index) buildIndex() error {
	vectors, err := i.storage.GetAllVectors()
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

func (i *Index) getKey(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func removeVector(vectors []*db.Vector, vector *db.Vector) []*db.Vector {
	for i, v := range vectors {
		if v.ID == vector.ID {
			return append(vectors[:i], vectors[i+1:]...)
		}
	}
	return vectors
}

func cosineSimilarity(v1, v2 db.Vector) (float64, error) {
	if v1.Compressed != v2.Compressed {
		return 0, errors.New("cannot calculate similarity between compressed and uncompressed vectors")
	}

	if v1.Compressed {
		if len(v1.Embedding) != len(v2.Embedding) {
			return 0, errors.New("embedding dimensions mismatch")
		}

		dotProduct := 0.0
		normV1 := 0.0
		normV2 := 0.0

		for i := range v1.Embedding {
			v1Val := v1.QuantizationParams.Dequantize(v1.Embedding[i])
			v2Val := v2.QuantizationParams.Dequantize(v2.Embedding[i])
			dotProduct += v1Val * v2Val
			normV1 += v1Val * v1Val
			normV2 += v2Val * v2Val
		}

		if normV1 == 0 || normV2 == 0 {
			return 0, nil
		}

		return dotProduct / (math.Sqrt(normV1) * math.Sqrt(normV2)), nil
	}

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

func Rerank(vectors []*db.Vector, searchQuery string) {
	for _, vector := range vectors {
		vector.Relevance = calculateRelevanceScore(vector, searchQuery)
	}
	sort.Slice(vectors, func(i, j int) bool {
		return vectors[i].Relevance > vectors[j].Relevance
	})
}

func calculateRelevanceScore(vector *db.Vector, searchQuery string) float64 {
	score := 0.0

	if strings.Contains(vector.Text, searchQuery) {
		score += 1.0
	}
	for _, value := range vector.Metadata {
		if strings.Contains(value, searchQuery) {
			score += 0.5
		}
	}

	distance := levenshtein.ComputeDistance(vector.Text, searchQuery)
	similarity := 1.0 - float64(distance)/float64(max(len(vector.Text), len(searchQuery)))
	score += similarity

	return score
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
