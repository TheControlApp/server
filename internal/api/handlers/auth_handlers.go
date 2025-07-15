package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	ScreenName string `json:"screen_name" binding:"required"`
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	user, err := h.UserService.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// TODO: Generate JWT token
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
		"token":   "jwt-token-placeholder",
	})
}

// Register creates a new user account
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userReq := services.CreateUserRequest{
		ScreenName: req.ScreenName,
		LoginName:  req.Username,
		Password:   req.Password,
	}

	user, err := h.UserService.CreateUser(userReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}
