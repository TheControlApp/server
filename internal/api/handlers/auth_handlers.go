package handlers

import (
	"errors"
	"net/http"
	"strings"

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
// @Summary      User login
// @Description  Authenticates a user and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "User credentials"
// @Success      200  {object}  responses.AuthResponse
// @Failure      400  {object}  responses.ValidationErrorResponse
// @Failure      401  {object}  responses.DetailedErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		var validationErrors []responses.ValidationError

		errorMsg := err.Error()
		if strings.Contains(errorMsg, "username") {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "username",
				Message: "Username is required",
			})
		}
		if strings.Contains(errorMsg, "password") {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "password",
				Message: "Password is required",
			})
		}

		if len(validationErrors) == 0 {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "request",
				Message: "Invalid JSON format or missing required fields",
			})
		}

		c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
		return
	}

	user, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			c.JSON(http.StatusUnauthorized, responses.DetailedErrorResponse{
				Error:   "Authentication failed",
				Code:    "USER_NOT_FOUND",
				Message: "No user found with the provided username",
			})
			return
		}

		if errors.Is(err, services.ErrInvalidPassword) {
			c.JSON(http.StatusUnauthorized, responses.DetailedErrorResponse{
				Error:   "Authentication failed",
				Code:    "INVALID_PASSWORD",
				Message: "The password provided is incorrect",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error: "Internal server error during authentication",
		})
		return
	}

	token, err := h.UserService.Auth.JWTManager.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to generate authentication token"})
		return
	}

	c.JSON(http.StatusOK, responses.AuthResponse{
		Message: "Login successful",
		User:    *user,
		Token:   token,
	})
}

// Register creates a new user account
// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body RegisterRequest true "User registration details"
// @Success      201  {object}  responses.UserResponse
// @Failure      400  {object}  responses.ValidationErrorResponse
// @Failure      409  {object}  responses.DetailedErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Parse JSON binding errors to provide field-specific feedback
		var validationErrors []responses.ValidationError

		errorMsg := err.Error()
		if strings.Contains(errorMsg, "username") || strings.Contains(errorMsg, "login_name") {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "username",
				Message: "Username is required and must be valid",
			})
		}
		if strings.Contains(errorMsg, "screen_name") {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "screen_name",
				Message: "Screen name is required",
			})
		}
		if strings.Contains(errorMsg, "password") {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "password",
				Message: "Password is required",
			})
		}
		if strings.Contains(errorMsg, "random_opt_in") {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "random_opt_in",
				Message: "Random opt-in field is required",
			})
		}

		if len(validationErrors) == 0 {
			validationErrors = append(validationErrors, responses.ValidationError{
				Field:   "request",
				Message: "Invalid JSON format or missing required fields",
			})
		}

		c.JSON(http.StatusBadRequest, responses.ValidationErrorResponse{
			Error:   "Validation failed",
			Details: validationErrors,
		})
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
		// Handle specific error types
		if errors.Is(err, services.ErrDuplicateUsername) {
			c.JSON(http.StatusConflict, responses.DetailedErrorResponse{
				Error:   "Registration failed",
				Code:    "DUPLICATE_USERNAME",
				Message: "A user with this username already exists. Please choose a different username.",
				Details: map[string]string{"username": req.Username},
			})
			return
		}

		if errors.Is(err, services.ErrInvalidUsername) {
			c.JSON(http.StatusBadRequest, responses.DetailedErrorResponse{
				Error:   "Registration failed",
				Code:    "INVALID_USERNAME",
				Message: err.Error(),
				Details: map[string]string{"username": req.Username},
			})
			return
		}

		if errors.Is(err, services.ErrPasswordTooWeak) {
			c.JSON(http.StatusBadRequest, responses.DetailedErrorResponse{
				Error:   "Registration failed",
				Code:    "PASSWORD_TOO_WEAK",
				Message: err.Error(),
			})
			return
		}

		// Generic server error for unexpected issues
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{
			Error: "Internal server error during user creation",
		})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{User: *user})
}
