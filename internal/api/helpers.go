package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/errors"
)

// BindJSON binds request JSON and returns typed struct or error response
func BindJSON[T any](c *gin.Context) (*T, *errors.AppError) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.BadRequest("Request body is not valid JSON or missing required fields")
	}
	return &req, nil
}

// GetUserID extracts user ID from context (set by JWT middleware)
func GetUserID(c *gin.Context) (uuid.UUID, *errors.AppError) {
	val, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.Unauthorized("User not authenticated", "")
	}
	userID, ok := val.(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.Internal("Invalid user ID in context")
	}
	return userID, nil
}

// JSON sends a JSON response with standard structure
func JSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

// Error sends an error response
func Error(c *gin.Context, err *errors.AppError) {
	c.JSON(err.Status, err)
}
