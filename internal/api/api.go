package api

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thecontrolapp/server/internal/api/handlers"
	"github.com/thecontrolapp/server/internal/auth"
	"github.com/thecontrolapp/server/internal/config"
	"github.com/thecontrolapp/server/internal/services"
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
	setupRoutes(router, db, hub, cfg)

	return router
}

// setupRoutes configures all the routes for the application (internal)
func setupRoutes(router *gin.Engine, db *gorm.DB, hub *websocket.Hub, cfg *config.Config) {
	// Initialize services
	jwtExpiration := time.Duration(cfg.Auth.JWTExpiration) * time.Second
	authService := auth.NewAuthService(cfg.Auth.JWTSecret, jwtExpiration)
	userService := services.NewUserService(db, authService)
	commandService := services.NewCommandService(db)

	// Initialize handlers from handlers package
	healthHandlers := handlers.NewHealthHandlers()
	userHandlers := handlers.NewUserHandlers(userService)
	authHandlers := handlers.NewAuthHandlers(userService)
	commandHandlers := handlers.NewCommandHandlers(commandService)
	wsHandlers := newWebSocketHandlers(hub, authService.JWTManager, userService, commandService)

	// Health check endpoint
	router.GET("/health", healthHandlers.HealthCheck)

	// Add Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandlers.Login)
			auth.POST("/register", authHandlers.Register)
		}

		// Command routes
		commands := v1.Group("/commands")
		{
			commands.GET("/pending", commandHandlers.GetPendingCommands)
			commands.POST("/complete", commandHandlers.CompleteCommand)
		}

		// User routes
		v1.GET("/users", userHandlers.GetUsers)
		v1.GET("/users/:id", userHandlers.GetUserByID)
	}

	// WebSocket route
	router.GET("/ws/client", wsHandlers.HandleClientWebSocket)
}
