package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
)

func main() {
	// Create a new storage instance
	storage, err := db.NewStorage("data/benchmark.db")
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	// Create a new index instance
	indexer, err := index.NewIndex(storage)
	if err != nil {
		panic(err)
	}

	// Generate 1,000,000 random vectors
	numVectors := 1000000
	vectorDim := 128
	vectors := make([]*db.Vector, numVectors)
	for i := 0; i < numVectors; i++ {
		embedding := make([]float64, vectorDim)
		for j := 0; j < vectorDim; j++ {
			embedding[j] = rand.Float64()
		}
		vectors[i] = &db.Vector{
			ID:        fmt.Sprintf("vector%d", i),
			Embedding: embedding,
			Metadata: map[string]string{
				"benchmark": "true",
			},
			Compressed: true,
			QuantizationParams: &db.QuantizationParams{
				Min:  -1.0,
				Max:  1.0,
				Bits: 8,
			},
			PruningMask:   make([]bool, vectorDim),
			SparseIndices: make([]int, 0),
		}
		// Randomly prune some elements of the vector
		numPruned := rand.Intn(vectorDim / 2)
		for j := 0; j < numPruned; j++ {
			index := rand.Intn(vectorDim)
			vectors[i].PruningMask[index] = true
			vectors[i].Embedding[index] = 0.0
		}
		// Convert the vector to sparse representation
		for j := 0; j < vectorDim; j++ {
			if vectors[i].Embedding[j] != 0.0 {
				vectors[i].SparseIndices = append(vectors[i].SparseIndices, j)
			}
		}
	}

	// Measure the embedding time
	start := time.Now()
	for _, vector := range vectors {
		err := indexer.Insert(vector)
		if err != nil {
			panic(err)
		}
	}
	elapsed := time.Since(start)

	// Print the benchmark results
	fmt.Printf("Number of vectors: %d\n", numVectors)
	fmt.Printf("Vector dimension: %d\n", vectorDim)
	fmt.Printf("Embedding time: %s\n", elapsed)
	fmt.Printf("Embedding speed: %.2f vectors/sec\n", float64(numVectors)/elapsed.Seconds())
}
