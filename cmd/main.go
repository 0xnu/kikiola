package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
	"github.com/0xnu/kikiola/pkg/server"
)

func main() {
	// Read configuration from environment variables
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/vectors.db"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3400"
	}

	// Initialize the storage
	storage, err := db.NewStorage(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	defer storage.Close()

	// Initialize the index
	index, err := index.NewIndex(storage)
	if err != nil {
		log.Fatalf("Failed to initialize index: %v", err)
	}

	// Initialize the server
	server := server.NewServer(storage, index)

	// Start the server
	log.Printf("Starting server on port %s...", port)
	go func() {
		if err := server.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
