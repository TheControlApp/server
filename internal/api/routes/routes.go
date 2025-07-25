package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/thecontrolapp/controlme-go/internal/api/handlers"
	"github.com/thecontrolapp/controlme-go/internal/api/responses"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/services"
	"github.com/thecontrolapp/controlme-go/internal/websocket"
	"gorm.io/gorm"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, db *gorm.DB, hub *websocket.Hub, cfg *config.Config) {
	// Initialize services
	jwtExpiration := time.Duration(cfg.Auth.JWTExpiration) * time.Hour
	authService := auth.NewAuthService(cfg.Auth.JWTSecret, jwtExpiration)
	userService := services.NewUserService(db, authService)
	commandService := services.NewCommandService(db)

	// Initialize handlers
	userHandlers := handlers.NewUserHandlers(userService)
	authHandlers := handlers.NewAuthHandlers(userService)
	commandHandlers := handlers.NewCommandHandlers(commandService)
	wsHandlers := handlers.NewWebSocketHandlers(hub)

	// Health check endpoint
	// Health godoc
	// @Summary      Health check
	// @Description  Check if the server is running
	// @Tags         health
	// @Accept       json
	// @Produce      json
	// @Success      200  {object}  responses.HealthResponse
	// @Router       /health [get]
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, responses.HealthResponse{
			Status:  "ok",
			Message: "Server is running",
		})
	})

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
		v1.POST("/users", userHandlers.CreateUser)
	}

	// WebSocket routes
	router.GET("/ws/client", wsHandlers.HandleClientWebSocket)
	router.GET("/ws/web", wsHandlers.HandleWebWebSocket)
}
