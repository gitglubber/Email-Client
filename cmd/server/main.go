package main

import (
	"fmt"
	"log"

	"github.com/gitglubber/email-client/internal/api"
	"github.com/gitglubber/email-client/internal/config"
)

func main() {
	printBanner()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	fmt.Println("Checking configuration...")
	validation := cfg.Validate()
	validation.PrintValidation()

	if !validation.Valid {
		fmt.Println("\n💡 Tip: Run './scripts/setup-oauth.sh' for guided setup")
		log.Fatal("Cannot start server with invalid configuration")
	}

	// Setup router
	router := api.SetupRouter(cfg)

	// Start server
	fmt.Println("\n🚀 Starting AI Email Client Server...")
	fmt.Println("=====================================")
	fmt.Printf("Backend:  http://localhost:%s\n", cfg.Port)
	fmt.Printf("Frontend: %s\n", cfg.FrontendURL)
	fmt.Println("=====================================")
	fmt.Println("\nPress Ctrl+C to stop the server\n")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func printBanner() {
	banner := `
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║     AI-Powered Email & Calendar Client                   ║
║     Built with Go + React + OpenAI                       ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}
