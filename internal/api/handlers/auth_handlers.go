package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecontrolapp/server/internal/api/responses"
	"github.com/thecontrolapp/server/internal/services"
)

type AuthHandlers struct {
	UserService *services.UserService
}

func NewAuthHandlers(userService *services.UserService) *AuthHandlers {
	return &AuthHandlers{UserService: userService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	ScreenName  string `json:"screen_name" binding:"required"`
	RandomOptIn bool   `json:"random_opt_in"`
}

// Login authenticates a user and returns a JWT token
// Login godoc
//
//	@Summary		User login
//	@Description	Authenticates a user and returns a JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		LoginRequest	true	"User credentials"
//	@Success		200			{object}	responses.AuthResponse
//	@Failure		400			{object}	responses.ErrorResponse
//	@Failure		401			{object}	responses.ErrorResponse
//	@Failure		500			{object}	responses.ErrorResponse
//	@Router			/auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewBadRequestError(
			"Request body is not valid JSON or missing required fields",
		))
		return
	}

	user, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) || errors.Is(err, services.ErrUnauthorized) {
			c.JSON(http.StatusUnauthorized, responses.NewUnauthorizedError(
				"Invalid username or password",
				"Please check your credentials and try again",
			))
			return
		}

		c.JSON(http.StatusInternalServerError, responses.NewInternalServerError("Internal server error during authentication"))
		return
	}

	token, err := h.UserService.Auth.JWTManager.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.NewInternalServerError("Failed to generate authentication token"))
		return
	}

	c.JSON(http.StatusOK, responses.AuthResponse{
		BaseResponse: responses.BaseResponse{
			Message: "Login successful",
		},
		User:  *user,
		Token: token,
	})
} // Register creates a new user account
// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		RegisterRequest	true	"User registration details"
//	@Success		201		{object}	responses.UserResponse
//	@Failure		400		{object}	responses.ErrorResponse
//	@Failure		409		{object}	responses.ConflictErrorResponse
//	@Failure		422		{object}	responses.ValidationErrorResponse
//	@Failure		500		{object}	responses.ErrorResponse
//	@Router			/auth/register [post]
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.NewBadRequestError(
			"Request body is not valid JSON or missing required fields",
		))
		return
	}

	userReq := services.CreateUserRequest{
		ScreenName:  req.ScreenName,
		LoginName:   req.Username,
		Password:    req.Password,
		RandomOptIn: req.RandomOptIn,
	}

	user, err := h.UserService.CreateUser(userReq)
	if err != nil {
		// Check if it's a validation error
		if validationErr, ok := err.(*services.ValidationError); ok {
			c.JSON(http.StatusUnprocessableEntity, responses.NewValidationError([]responses.ValidationError{{
				Field:   validationErr.Field,
				Message: validationErr.Message,
				Code:    validationErr.Code,
			}}))
			return
		}

		// Check if it's a conflict error
		if conflictErr, ok := err.(*services.ConflictError); ok {
			c.JSON(http.StatusConflict, responses.NewConflictError(
				conflictErr.Message,
				map[string]string{conflictErr.Resource: conflictErr.Value},
				"Please choose a different username and try again",
			))
			return
		}

		// Generic server error for unexpected issues
		c.JSON(http.StatusInternalServerError, responses.NewInternalServerError("Internal server error during user creation"))
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{User: *user})
}
