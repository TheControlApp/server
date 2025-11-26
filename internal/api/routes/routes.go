package routes

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

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, db *gorm.DB, hub *websocket.Hub, cfg *config.Config) {
	// Initialize services
	jwtExpiration := time.Duration(cfg.Auth.JWTExpiration) * time.Second
	authService := auth.NewAuthService(cfg.Auth.JWTSecret, jwtExpiration)
	userService := services.NewUserService(db, authService)
	commandService := services.NewCommandService(db)

	// Initialize handlers
	healthHandlers := handlers.NewHealthHandlers()
	userHandlers := handlers.NewUserHandlers(userService)
	authHandlers := handlers.NewAuthHandlers(userService)
	commandHandlers := handlers.NewCommandHandlers(commandService)
	wsHandlers := handlers.NewWebSocketHandlers(hub, authService.JWTManager, userService)

	// Health check endpoint
	router.GET("/health", healthHandlers.HealthCheck)

	// Add Swagger route
	// The URL points to the auto-generated swagger.json file.
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

	// WebSocket route - single endpoint for all clients
	// WebSocket godoc
	//	@Summary		WebSocket connection endpoint
	//	@Description	Establishes WebSocket connection for real-time command distribution. Supports anonymous and authenticated connections.
	//	@Tags			websocket
	//	@Accept			json
	//	@Produce		json
	//	@Param			Authorization	header		string	false	"Bearer token for authentication"
	//	@Param			token			query		string	false	"Token for authentication via query parameter"
	//	@Success		101				{string}	string	"Switching Protocols"
	//	@Router			/ws/client [get]
	router.GET("/ws/client", wsHandlers.HandleClientWebSocket)
}
