package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/services"
)

type CommandHandlers struct {
	Service *services.CommandService
}

func NewCommandHandlers(service *services.CommandService) *CommandHandlers {
	return &CommandHandlers{Service: service}
}

// GetPendingCommands gets pending commands for a user
func (h *CommandHandlers) GetPendingCommands(c *gin.Context) {
	// TODO: Get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	command, err := h.Service.GetNextCommand(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch commands"})
		return
	}

	if command == nil {
		c.JSON(http.StatusOK, gin.H{"command": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"command": command})
}

// CompleteCommand marks a command as completed
func (h *CommandHandlers) CompleteCommand(c *gin.Context) {
	// TODO: Get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = h.Service.CompleteCommand(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete command"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Command completed successfully"})
}
