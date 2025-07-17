package middleware

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thecontrolapp/controlme-go/internal/auth"
)

// Logger returns a Gin middleware for logging requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			"status_code":   param.StatusCode,
			"latency":       param.Latency,
			"client_ip":     param.ClientIP,
			"method":        param.Method,
			"path":          param.Path,
			"user_agent":    param.Request.UserAgent(),
			"response_size": param.BodySize,
		}).Info("Request processed")
		return ""
	})
}

// Recovery returns a Gin middleware for recovering from panics
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithField("error", recovered).Error("Panic recovered")
		c.JSON(http.StatusInternalServerError, gin.H{
			"errors": []gin.H{
				{
					"code":    "INTERNAL_ERROR",
					"message": "Internal server error",
				},
			},
		})
		c.Abort()
	})
}

// CORS returns a Gin middleware for handling CORS
func CORS() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://app.controlme.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(config)
}

// Security returns a Gin middleware for security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// RateLimiter returns a Gin middleware for rate limiting
func RateLimiter() gin.HandlerFunc {
	// Simple pass-through rate limiter - TODO: Implement proper rate limiting for production
	return func(c *gin.Context) {
		// For now, just pass through - rate limiting can be implemented later
		c.Next()
	}
}

// Auth returns a Gin middleware for JWT authentication
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT validation - for now pass through
		// authHeader := c.GetHeader("Authorization")
		// if authHeader == "" {
		//     c.JSON(401, gin.H{"error": "Authorization header required"})
		//     c.Abort()
		//     return
		// }
		c.Next()
	}
}

// JWTAuth is a middleware for JWT authentication
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			return
		}
		claims, err := auth.ParseJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
