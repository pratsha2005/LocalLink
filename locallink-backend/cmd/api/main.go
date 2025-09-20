package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LocalLink/internal/api"
	"github.com/LocalLink/internal/config"
	"github.com/LocalLink/internal/database"
	"github.com/LocalLink/internal/websocket"
)

func main() {
	cfg := config.Load()

	dbPool := database.Connect(cfg.DatabaseURL)
	defer dbPool.Close()
	store := database.NewStore(dbPool)

	hub := websocket.NewHub()
	go hub.Run()

	router := api.NewRouter(store, cfg, hub)

	serverAddr := ":8080"
	fmt.Printf("Starting server on %s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}