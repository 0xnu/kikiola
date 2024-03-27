package main

import (
	"log"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
	"github.com/0xnu/kikiola/pkg/server"
)

func main() {
	// Initialise the storage
	storage, err := db.NewStorage("data/vectors.db")
	if err != nil {
		log.Fatalf("Failed to initialise storage: %v", err)
	}
	defer storage.Close()

	// Initialise the index
	index, err := index.NewIndex(storage)
	if err != nil {
		log.Fatalf("Failed to initialise index: %v", err)
	}

	// Initialise the server
	server := server.NewServer(storage, index)

	// Start the server
	log.Println("Starting server on port 3400...")
	if err := server.Start(":3400"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
