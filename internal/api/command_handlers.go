package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/errors"
	"github.com/thecontrolapp/server/internal/services"
)

type CommandHandlers struct {
	Service *services.CommandService
}

func newCommandHandlers(service *services.CommandService) *CommandHandlers {
	return &CommandHandlers{Service: service}
}

// GetPendingCommands godoc
//
//	@Summary		Get pending commands for a user
//	@Description	Retrieves pending commands for a given user
//	@Tags			commands
//	@Accept			json
//	@Produce		json
//	@Param			user_id	query		string	true	"User ID"
//	@Success		200		{object}	CommandsResponse
//	@Failure		400		{object}	errors.AppError
//	@Router			/commands/pending [get]
func (h *CommandHandlers) GetPendingCommands(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		Error(c, errors.BadRequest("user_id is required"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		Error(c, errors.BadRequest("Invalid user_id format"))
		return
	}

	commands, dbErr := h.Service.GetPendingCommands(userID)
	if dbErr != nil {
		Error(c, errors.Internal("Failed to fetch commands"))
		return
	}

	JSON(c, http.StatusOK, CommandsResponse{Commands: commands})
}

// CompleteCommand godoc
//
//	@Summary		Mark a command as completed
//	@Description	Marks a specific command as completed
//	@Tags			commands
//	@Accept			json
//	@Produce		json
//	@Param			command_id	query	string	true	"Command ID"
//	@Param			user_id		query	string	true	"User ID"
//	@Success		200			{object}	map[string]string
//	@Failure		400			{object}	errors.AppError
//	@Router			/commands/complete [post]
func (h *CommandHandlers) CompleteCommand(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		Error(c, errors.BadRequest("user_id is required"))
		return
	}

	commandIDStr := c.Query("command_id")
	if commandIDStr == "" {
		Error(c, errors.BadRequest("command_id is required"))
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		Error(c, errors.BadRequest("Invalid user_id format"))
		return
	}

	commandID, err := uuid.Parse(commandIDStr)
	if err != nil {
		Error(c, errors.BadRequest("Invalid command_id format"))
		return
	}

	if dbErr := h.Service.CompleteCommand(commandID, userID); dbErr != nil {
		Error(c, errors.Internal("Failed to complete command"))
		return
	}

	JSON(c, http.StatusOK, gin.H{"message": "Command completed successfully"})
}
