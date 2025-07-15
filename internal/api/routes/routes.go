package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thecontrolapp/controlme-go/internal/api/handlers"
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

	// Initialize handlers
	userHandlers := handlers.NewUserHandlers(userService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "controlme-go",
			"message": "Server running with modern authentication",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Login endpoint - TODO"})
			})
			auth.POST("/register", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Register endpoint - TODO"})
			})
		}

		// Command routes
		commands := v1.Group("/commands")
		{
			commands.GET("/pending", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Pending commands - TODO"})
			})
			commands.POST("/complete", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "Complete command - TODO"})
			})
		}

		// User routes
		v1.GET("/users", userHandlers.GetUsers)
		v1.GET("/users/:id", userHandlers.GetUserByID)
		v1.POST("/users", userHandlers.CreateUser)
	}

	// WebSocket routes
	router.GET("/ws/client", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "WebSocket client endpoint - TODO"})
	})

	router.GET("/ws/web", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "WebSocket web endpoint - TODO"})
	})
}
