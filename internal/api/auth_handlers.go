package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thecontrolapp/server/internal/services"
)

type AuthHandlers struct {
	UserService *services.UserService
}

func newAuthHandlers(userService *services.UserService) *AuthHandlers {
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

// Login godoc
//
//	@Summary		User login
//	@Description	Authenticates a user and returns a JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		LoginRequest	true	"User credentials"
//	@Success		200			{object}	responses.AuthResponse
//	@Failure		400			{object}	errors.AppError
//	@Failure		401			{object}	errors.AppError
//	@Router			/auth/login [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	req, err := BindJSON[LoginRequest](c)
	if err != nil {
		Error(c, err)
		return
	}

	user, token, appErr := h.UserService.AuthenticateUser(req.Username, req.Password)
	if appErr != nil {
		Error(c, appErr)
		return
	}

	JSON(c, http.StatusOK, AuthResponse{
		Message: "Login successful",
		User:    *user,
		Token:   token,
	})
}

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Creates a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		RegisterRequest	true	"User registration details"
//	@Success		201		{object}	responses.UserResponse
//	@Failure		400		{object}	errors.AppError
//	@Failure		409		{object}	errors.AppError
//	@Failure		422		{object}	errors.AppError
//	@Router			/auth/register [post]
func (h *AuthHandlers) Register(c *gin.Context) {
	req, err := BindJSON[RegisterRequest](c)
	if err != nil {
		Error(c, err)
		return
	}

	user, appErr := h.UserService.CreateUser(services.CreateUserRequest{
		ScreenName:  req.ScreenName,
		LoginName:   req.Username,
		Password:    req.Password,
		RandomOptIn: req.RandomOptIn,
	})
	if appErr != nil {
		Error(c, appErr)
		return
	}

	JSON(c, http.StatusCreated, UserResponse{User: *user})
}
