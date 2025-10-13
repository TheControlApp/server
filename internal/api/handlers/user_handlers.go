package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/api/responses"
	"github.com/thecontrolapp/server/internal/services"
)

type UserHandlers struct {
	Service *services.UserService
}

func NewUserHandlers(service *services.UserService) *UserHandlers {
	return &UserHandlers{Service: service}
}

// UserHandler provides modern RESTful user endpoints
// GetUsers godoc
// @Summary      Get all users
// @Description  Retrieves a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.UsersResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /users [get]
func (h *UserHandlers) GetUsers(c *gin.Context) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewInternalServerError("Failed to fetch users"))
		return
	}
	c.JSON(http.StatusOK, responses.UsersResponse{Users: users})
}

// GetUserByID godoc
// @Summary      Get a user by ID
// @Description  Retrieves a user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200  {object}  responses.UserResponse
// @Failure      400  {object}  responses.ValidationErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /users/{id} [get]
func (h *UserHandlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	userID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, responses.NewInvalidFormatError("id", "UUID"))
		return
	}
	user, err := h.Service.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, responses.NewNotFoundError("User not found"))
		return
	}
	c.JSON(http.StatusOK, responses.UserResponse{User: *user})
}
