package main

import (
	"fmt"
	"log"
	"os"

	"github.com/guilherme096/data-sync/internal/api"
	"github.com/guilherme096/data-sync/internal/trino"
	"github.com/guilherme096/data-sync/pkg/data-sync/discovery"
	"github.com/guilherme096/data-sync/pkg/data-sync/storage"
	"github.com/guilherme096/data-sync/pkg/data-sync/sync"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	trinoHost := os.Getenv("TRINO_HOST")
	if trinoHost == "" {
		trinoHost = "localhost"
	}

	trinoPort := os.Getenv("TRINO_PORT")
	if trinoPort == "" {
		trinoPort = "8080"
	}

	trinoUser := os.Getenv("TRINO_USER")
	if trinoUser == "" {
		trinoUser = "trino"
	}

	trinoCatalog := os.Getenv("TRINO_CATALOG")
	trinoSchema := os.Getenv("TRINO_SCHEMA")

	connString := fmt.Sprintf("http://%s@%s:%s", trinoUser, trinoHost, trinoPort)
	if trinoCatalog != "" {
		connString += "?catalog=" + trinoCatalog
		if trinoSchema != "" {
			connString += "&schema=" + trinoSchema
		}
	}

	engine, err := trino.NewEngine(connString)
	if err != nil {
		log.Fatalf("Failed to create Trino engine: %v", err)
	}
	defer engine.Close()

	metadataDiscovery := discovery.NewTrinoMetadataDiscovery(engine)

	// in-memory
	metadataStorage := storage.NewMemoryMetadataStorage()

	syncService := sync.NewMetadataSync(metadataDiscovery, metadataStorage)

	log.Println("Performing initial metadata sync...")
	if err := syncService.SyncAll(); err != nil {
		log.Printf("Warning: initial sync failed: %v", err)
	} else {
		log.Println("Initial metadata sync completed successfully")
	}

	srv := api.NewServer(":"+port, engine, metadataStorage, syncService)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
