package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/0xnu/kikiola/pkg/db"
	"github.com/0xnu/kikiola/pkg/index"
	"github.com/0xnu/kikiola/pkg/server"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "data/vectors.db"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3400"
	}

	hostAddress := "localhost"

	nodeAddresses := generateNodeAddresses(hostAddress, 3401, 3420)

	storage, err := db.NewDistributedStorage(nodeAddresses)
	if err != nil {
		log.Fatalf("Failed to initialize distributed storage: %v", err)
	}
	defer storage.Close()

	index, err := index.NewIndex(storage)
	if err != nil {
		log.Fatalf("Failed to initialize index: %v", err)
	}

	server := server.NewServer(storage, index)

	log.Printf("Starting server on %s:%s...", hostAddress, port)
	go func() {
		if err := server.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func generateNodeAddresses(hostAddress string, startPort, endPort int) []string {
	var nodeAddresses []string
	for port := startPort; port <= endPort; port++ {
		nodeAddress := hostAddress + ":" + strconv.Itoa(port)
		nodeAddresses = append(nodeAddresses, nodeAddress)
	}
	return nodeAddresses
}
