package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
)

func main() {
	nodeAddresses := generateNodeAddresses("localhost", 3401, 3420)

	storage, err := db.NewDistributedStorage(nodeAddresses)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	indexer, err := index.NewIndex(storage)
	if err != nil {
		panic(err)
	}

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

		numPruned := rand.Intn(vectorDim / 2)
		for j := 0; j < numPruned; j++ {
			index := rand.Intn(vectorDim)
			vectors[i].PruningMask[index] = true
			vectors[i].Embedding[index] = 0.0
		}

		for j := 0; j < vectorDim; j++ {
			if vectors[i].Embedding[j] != 0.0 {
				vectors[i].SparseIndices = append(vectors[i].SparseIndices, j)
			}
		}
	}

	start := time.Now()
	for _, vector := range vectors {
		err := indexer.Insert(vector)
		if err != nil {
			panic(err)
		}
	}
	elapsed := time.Since(start)

	fmt.Printf("Number of vectors: %d\n", numVectors)
	fmt.Printf("Vector dimension: %d\n", vectorDim)
	fmt.Printf("Embedding time: %s\n", elapsed)
	fmt.Printf("Embedding speed: %.2f vectors/sec\n", float64(numVectors)/elapsed.Seconds())
}

func generateNodeAddresses(hostAddress string, startPort, endPort int) []string {
	var nodeAddresses []string
	for port := startPort; port <= endPort; port++ {
		nodeAddress := hostAddress + ":" + fmt.Sprintf("%d", port)
		nodeAddresses = append(nodeAddresses, nodeAddress)
	}
	return nodeAddresses
}
