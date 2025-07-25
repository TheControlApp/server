package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecontrolapp/controlme-go/internal/api/responses"
	"github.com/thecontrolapp/controlme-go/internal/services"
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
	Email       string `json:"email" binding:"required"`
	RandomOptIn bool   `json:"random_opt_in" binding:"required"`
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
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Invalid request"})
		return
	}

	user, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, responses.ErrorResponse{Error: "Invalid credentials"})
		return
	}

	token, err := h.UserService.Auth.JWTManager.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to generate token"})
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
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /auth/register [post]
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, responses.ErrorResponse{Error: "Invalid request"})
		return
	}

	userReq := services.CreateUserRequest{
		ScreenName:  req.ScreenName,
		LoginName:   req.Username,
		Password:    req.Password,
		Email:       req.Email,
		RandomOptIn: req.RandomOptIn,
	}

	user, err := h.UserService.CreateUser(userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse{Error: "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, responses.UserResponse{User: *user})
}
