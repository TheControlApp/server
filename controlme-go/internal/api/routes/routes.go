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
	authService := auth.NewAuthService(cfg.Legacy.CryptoKey, cfg.Auth.JWTSecret, jwtExpiration)
	userService := services.NewUserService(db, authService)
	cmdService := services.NewCommandService(db)
	legacyService := services.NewLegacyService(db, authService)

	// Initialize handlers
	legacyHandlers := handlers.NewLegacyHandlers(db, userService, cmdService, legacyService, authService, cfg)
	userHandlers := handlers.NewUserHandlers(userService)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "controlme-go",
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

	// Legacy routes (exact compatibility)
	setupLegacyRoutes(router, legacyHandlers)
}

// setupLegacyRoutes sets up the legacy ASP.NET compatible routes
func setupLegacyRoutes(router *gin.Engine, handlers *handlers.LegacyHandlers) {
	// Legacy command endpoints
	router.GET("/AppCommand.aspx", handlers.AppCommand)
	router.GET("/GetContent.aspx", handlers.GetContent)
	router.GET("/GetCount.aspx", handlers.GetCount)
	router.POST("/ProcessComplete.aspx", handlers.ProcessComplete)
	router.POST("/DeleteOut.aspx", handlers.DeleteOut)
	router.GET("/GetOptions.aspx", handlers.GetOptions)

	// Legacy web interface endpoints (placeholder for now)
	router.GET("/Default.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Default page - TODO")
	})

	router.GET("/Messages.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Messages page - TODO")
	})

	router.POST("/Upload.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Upload endpoint - TODO")
	})

	// Legacy auth pages (placeholder for now)
	router.GET("/Pages/Login.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Login page - TODO")
	})

	router.POST("/Pages/Login.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Login processing - TODO")
	})

	router.GET("/Pages/Register.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Register page - TODO")
	})

	router.POST("/Pages/Register.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Register processing - TODO")
	})

	router.GET("/Pages/ControlPC.aspx", func(c *gin.Context) {
		c.String(200, "Legacy ControlPC page - TODO")
	})
}
