package main

import (
	"fmt"
	"log"
	"os"

	"github.com/thecontrolapp/controlme-go/internal/api"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/database"
	"github.com/thecontrolapp/controlme-go/internal/websocket"
)

// @title ControlMe API
// @version 1.0
// @description Modern ControlMe backend API
// @termsOfService http://swagger.io/terms/

// @contact.name ControlMe Support
// @contact.url http://www.controlme.io/support
// @contact.email support@controlme.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	log.Printf("🔧 Starting ControlMe API in %s mode", cfg.Environment)

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Set up API router
	router := api.SetupRouter(db, hub, cfg)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = fmt.Sprintf("%d", cfg.Server.Port)
	}

	log.Printf("🚀 ControlMe API server starting on port %s", port)
	log.Printf("📚 API Documentation: http://localhost:%s/swagger/index.html", port)
	log.Printf("❤️  Health Check: http://localhost:%s/health", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
