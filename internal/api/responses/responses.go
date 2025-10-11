package responses

import "github.com/thecontrolapp/server/internal/models"

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	Message string      `json:"message" example:"Login successful"`
	User    models.User `json:"user"`
	Token   string      `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserResponse represents a single user response
type UserResponse struct {
	User models.User `json:"user"`
}

// UsersResponse represents a list of users response
type UsersResponse struct {
	Users []models.User `json:"users"`
}

// CommandsResponse represents a list of commands response
type CommandsResponse struct {
	Commands []models.Command `json:"commands"`
}

// MessageResponse represents a simple message response
type MessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// ValidationError represents a field-specific validation error
type ValidationError struct {
	Field   string `json:"field" example:"username"`
	Message string `json:"message" example:"Username must be at least 3 characters"`
}

// ValidationErrorResponse represents validation errors for form fields
type ValidationErrorResponse struct {
	Error   string            `json:"error" example:"Validation failed"`
	Details []ValidationError `json:"details"`
}

// DetailedErrorResponse represents a detailed error with code and context
type DetailedErrorResponse struct {
	Error   string      `json:"error" example:"Registration failed"`
	Code    string      `json:"code" example:"DUPLICATE_USERNAME"`
	Message string      `json:"message" example:"A user with this username already exists"`
	Details interface{} `json:"details,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Server is running"`
}
