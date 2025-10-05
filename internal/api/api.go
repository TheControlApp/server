package api

import (
	"github.com/gin-gonic/gin"
	"github.com/thecontrolapp/server/internal/api/routes"
	"github.com/thecontrolapp/server/internal/config"
	"github.com/thecontrolapp/server/internal/websocket"
	"gorm.io/gorm"
)

// SetupRouter configures and returns the main HTTP router
func SetupRouter(db *gorm.DB, hub *websocket.Hub, cfg *config.Config) *gin.Engine {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, db, hub, cfg)

	return router
}
