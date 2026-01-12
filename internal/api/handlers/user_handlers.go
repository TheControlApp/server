package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/api"
	"github.com/thecontrolapp/server/internal/errors"
	"github.com/thecontrolapp/server/internal/services"
)

type UserHandlers struct {
	Service *services.UserService
}

func NewUserHandlers(service *services.UserService) *UserHandlers {
	return &UserHandlers{Service: service}
}

// GetUsers godoc
//
//	@Summary		Get all users
//	@Description	Retrieves a list of all users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	api.UsersResponse
//	@Failure		500	{object}	errors.AppError
//	@Router			/users [get]
func (h *UserHandlers) GetUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		api.Error(c, err)
		return
	}
	api.JSON(c, http.StatusOK, api.UsersResponse{Users: users})
}

// GetUserByID godoc
//
//	@Summary		Get a user by ID
//	@Description	Retrieves a user by their ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"User ID"
//	@Success		200	{object}	api.UserResponse
//	@Failure		400	{object}	errors.AppError
//	@Failure		404	{object}	errors.AppError
//	@Router			/users/{id} [get]
func (h *UserHandlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		api.Error(c, errors.BadRequest("Invalid user ID format"))
		return
	}
	user, appErr := h.Service.GetUserByID(userID)
	if appErr != nil {
		api.Error(c, appErr)
		return
	}
	api.JSON(c, http.StatusOK, api.UserResponse{User: *user})
}
