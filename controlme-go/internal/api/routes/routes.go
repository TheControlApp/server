package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/websocket"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(router *gin.Engine, db *gorm.DB, hub *websocket.Hub, cfg *config.Config) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
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
	}

	// WebSocket routes
	router.GET("/ws/client", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "WebSocket client endpoint - TODO"})
	})
	
	router.GET("/ws/web", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "WebSocket web endpoint - TODO"})
	})

	// Legacy routes (exact compatibility)
	setupLegacyRoutes(router, db, hub, cfg)
}

// setupLegacyRoutes sets up the legacy ASP.NET compatible routes
func setupLegacyRoutes(router *gin.Engine, db *gorm.DB, hub *websocket.Hub, cfg *config.Config) {
	// Legacy command endpoints
	router.GET("/AppCommand.aspx", func(c *gin.Context) {
		c.String(200, "Legacy AppCommand endpoint - TODO")
	})
	
	router.GET("/GetContent.aspx", func(c *gin.Context) {
		c.String(200, "Legacy GetContent endpoint - TODO")
	})
	
	router.GET("/GetCount.aspx", func(c *gin.Context) {
		c.String(200, "Legacy GetCount endpoint - TODO")
	})
	
	router.POST("/ProcessComplete.aspx", func(c *gin.Context) {
		c.String(200, "Legacy ProcessComplete endpoint - TODO")
	})
	
	router.POST("/DeleteOut.aspx", func(c *gin.Context) {
		c.String(200, "Legacy DeleteOut endpoint - TODO")
	})
	
	router.GET("/GetOptions.aspx", func(c *gin.Context) {
		c.String(200, "Legacy GetOptions endpoint - TODO")
	})
	
	// Legacy web interface endpoints
	router.GET("/Default.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Default page - TODO")
	})
	
	router.GET("/Messages.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Messages page - TODO")
	})
	
	router.POST("/Upload.aspx", func(c *gin.Context) {
		c.String(200, "Legacy Upload endpoint - TODO")
	})
	
	// Legacy auth pages
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
