package main

import (
	"log"
	"os"

	"github.com/guilherme096/data-sync/internal/api"
)

func main() {
	// 1. Configuration (Env vars, defaults)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	// 2. Setup Dependencies (Database, Core Service)
	// repo := trino.NewRepository(...)
	// service := core.NewService(repo)

	// 3. Start API Server
	srv := api.NewServer(":" + port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
