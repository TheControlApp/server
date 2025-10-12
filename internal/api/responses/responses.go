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

// ErrorResponse represents a simple error response
type ErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}

// ValidationError represents a field-specific validation error
type ValidationError struct {
	Field   string `json:"field" example:"username"`
	Message string `json:"message" example:"Username must be at least 3 characters"`
	Code    string `json:"code" example:"TOO_SHORT"`
}

// SimpleErrorResponse represents a simplified error response for basic use cases
type SimpleErrorResponse struct {
	Error   string `json:"error" example:"Username is too short"`
	Code    string `json:"code,omitempty" example:"VALIDATION_ERROR"`
	Message string `json:"message,omitempty" example:"Please check your input and try again"`
}

// ValidationErrorResponse represents validation errors for form fields (HTTP 422)
type ValidationErrorResponse struct {
	Type   string            `json:"type" example:"validation_error"`
	Title  string            `json:"title" example:"Validation Failed"`
	Status int               `json:"status" example:"422"`
	Detail string            `json:"detail" example:"One or more fields failed validation"`
	Errors []ValidationError `json:"errors"`
	// Developer-friendly addition
	Help string `json:"help,omitempty" example:"Check the API documentation for field requirements"`
}

// ConflictErrorResponse represents resource conflict errors (HTTP 409)
type ConflictErrorResponse struct {
	Type     string      `json:"type" example:"conflict"`
	Title    string      `json:"title" example:"Resource Conflict"`
	Status   int         `json:"status" example:"409"`
	Detail   string      `json:"detail" example:"Username already exists"`
	Instance interface{} `json:"instance,omitempty"`
	// Developer-friendly addition
	Action string `json:"action,omitempty" example:"Please choose a different username"`
}

// UnauthorizedErrorResponse represents authentication errors (HTTP 401)
type UnauthorizedErrorResponse struct {
	Type   string `json:"type" example:"unauthorized"`
	Title  string `json:"title" example:"Authentication Failed"`
	Status int    `json:"status" example:"401"`
	Detail string `json:"detail" example:"Invalid credentials provided"`
	// Developer-friendly addition
	Action string `json:"action,omitempty" example:"Please check your username and password"`
}

// BadRequestErrorResponse represents malformed request errors (HTTP 400)
type BadRequestErrorResponse struct {
	Type   string `json:"type" example:"bad_request"`
	Title  string `json:"title" example:"Bad Request"`
	Status int    `json:"status" example:"400"`
	Detail string `json:"detail" example:"Request body is not valid JSON"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Server is running"`
}

// Helper functions for creating common error responses

// NewSimpleError creates a simple error response
func NewSimpleError(error, code, message string) SimpleErrorResponse {
	return SimpleErrorResponse{
		Error:   error,
		Code:    code,
		Message: message,
	}
}

// NewValidationError creates a validation error response with helpful info
func NewValidationError(errors []ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{
		Type:   "validation_error",
		Title:  "Validation Failed",
		Status: 422,
		Detail: "One or more fields failed validation",
		Errors: errors,
		Help:   "Please check the field requirements in the API documentation",
	}
}
