// cmd/api/main.go
package main

import (
	"fmt"
	"github.com/LocalLink/internal/api"
	"github.com/LocalLink/internal/config"
	"github.com/LocalLink/internal/database"
	"log"
	"net/http"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Connect to Database
	dbPool := database.Connect(cfg.DatabaseURL)
	defer dbPool.Close()
	store := database.NewStore(dbPool)

	// 3. Initialize Router
	router := api.NewRouter(store, cfg)

	// 4. Start Server
	serverAddr := ":8080"
	fmt.Printf("Starting server on %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}