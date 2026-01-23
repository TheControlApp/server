package api

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thecontrolapp/server/internal/auth"
	"github.com/thecontrolapp/server/internal/errors"
)

// Logger middleware with structured logging
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logrus.WithFields(logrus.Fields{
			"status": param.StatusCode,
			"ms":     param.Latency.Milliseconds(),
			"ip":     param.ClientIP,
			"method": param.Method,
			"path":   param.Path,
		}).Info("request")
		return ""
	})
}

// Recovery middleware with standardized error response
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithField("panic", recovered).Error("Panic recovered")
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.Internal("Internal server error"))
	})
}

// CORS returns a Gin middleware for handling CORS
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173", "https://app.controlme.io"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

// Security middleware adds security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Header("Server", "")
		c.Next()
	}
}

// JWTAuth middleware validates JWT tokens and sets user_id in context
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.Unauthorized("Missing authorization token", "Include Authorization header with valid JWT token"))
			return
		}
		claims, err := auth.ParseJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.Unauthorized("Invalid authorization token", "Ensure token is valid and not expired"))
			return
		}
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
