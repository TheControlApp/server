package api

import (
	"github.com/thecontrolapp/server/internal/models"
)

// Success responses with data
type AuthResponse struct {
	Message string      `json:"message,omitempty"`
	User    models.User `json:"user"`
	Token   string      `json:"token,omitempty"`
}

type UserResponse struct {
	User models.User `json:"user"`
}

type UsersResponse struct {
	Users []models.User `json:"users"`
}

type CommandsResponse struct {
	Commands []models.Command `json:"commands"`
}

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
