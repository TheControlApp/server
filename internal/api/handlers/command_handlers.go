package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/api/responses"
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
// @Success      200  {object}  responses.CommandsResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /commands/pending [get]
func (h *CommandHandlers) GetPendingCommands(c *gin.Context) {
	// TODO: Get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "user_id required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	commands, err := h.Service.GetPendingCommands(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to fetch commands"})
		return
	}

	c.JSON(http.StatusOK, responses.CommandsResponse{Commands: commands})
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
// @Success      200  {object}  responses.MessageResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /commands/complete [post]
func (h *CommandHandlers) CompleteCommand(c *gin.Context) {
	// TODO: Get user ID from JWT token
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "user_id required"})
		return
	}

	commandIDStr := c.Query("command_id")
	if commandIDStr == "" {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "command_id required"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Invalid user ID"})
		return
	}

	commandID, err := uuid.Parse(commandIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Invalid command ID"})
		return
	}

	err = h.Service.CompleteCommand(commandID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to complete command"})
		return
	}

	c.JSON(http.StatusOK, responses.MessageResponse{Message: "Command completed successfully"})
}
