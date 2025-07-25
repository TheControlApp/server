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
// GetPendingCommands godoc
// @Summary      Get pending commands for a user
// @Description  Retrieves pending commands for a given user
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        user_id query string true "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /commands/pending [get]
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

	commands, err := h.Service.GetPendingCommands(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch commands"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"commands": commands})
}

// CompleteCommand marks a command as completed
// CompleteCommand godoc
// @Summary      Mark a command as completed
// @Description  Marks a specific command as completed
// @Tags         commands
// @Accept       json
// @Produce      json
// @Param        command_id query string true "Command ID"
// @Param        user_id query string true "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /commands/complete [post]
func (h *CommandHandlers) CompleteCommand(c *gin.Context) {
	// TODO: Get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required"})
		return
	}

	commandIDStr := c.Query("command_id")
	if commandIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "command_id required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	commandID, err := uuid.Parse(commandIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid command ID"})
		return
	}

	err = h.Service.CompleteCommand(commandID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete command"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Command completed successfully"})
}
