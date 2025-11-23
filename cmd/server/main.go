package main

import (
	"log"

	"github.com/gitglubber/email-client/internal/api"
	"github.com/gitglubber/email-client/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup router
	router := api.SetupRouter(cfg)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	log.Printf("Frontend URL: %s", cfg.FrontendURL)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
