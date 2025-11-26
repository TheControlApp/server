package responses

import (
	"fmt"

	"github.com/thecontrolapp/server/internal/models"
)

// PaymentRequiredErrorResponse represents a payment requ// NewUnauthorizedError creates an unauthorized error response
func NewUnauthorizedError(detail, action string) ErrorResponse {
	return ErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "unauthorized",
			Title:  "Unauthorized Access",
			Status: 401,
			Detail: detail,
		},
		DeveloperHelp: DeveloperHelp{
			Action: action,
		},
	}
}

// Response structures to reduce duplication

// BaseResponse provides common fields for successful responses
type BaseResponse struct {
	Message string `json:"message,omitempty" example:"Operation completed successfully"`
}

// ProblemDetails represents RFC 7807 Problem Details base structure
type ProblemDetails struct {
	Type   string `json:"type" example:"validation_error"`
	Title  string `json:"title" example:"Validation Failed"`
	Status int    `json:"status" example:"422"`
	Detail string `json:"detail" example:"One or more fields failed validation"`
}

// DeveloperHelp provides optional developer assistance fields
type DeveloperHelp struct {
	Help   string `json:"help,omitempty" example:"Check the API documentation for field requirements"`
	Action string `json:"action,omitempty" example:"Please choose a different username"`
}

// Specific response types

// AuthResponse represents the response for authentication endpoints
type AuthResponse struct {
	BaseResponse
	User  models.User `json:"user"`
	Token string      `json:"token,omitempty" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
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
	BaseResponse
}

// ValidationError represents a field-specific validation error
type ValidationError struct {
	Field   string `json:"field" example:"username"`
	Message string `json:"message" example:"Username must be at least 3 characters"`
	Code    string `json:"code" example:"MIN_LENGTH"`
}

// ValidationErrorResponse represents validation errors for form fields (HTTP 422)
type ValidationErrorResponse struct {
	ProblemDetails
	DeveloperHelp
	Errors []ValidationError `json:"errors"`
}

// ErrorResponse represents all standard API errors with optional developer help
type ErrorResponse struct {
	ProblemDetails
	DeveloperHelp `json:",omitempty"`
}

// ConflictErrorResponse has additional Instance field
type ConflictErrorResponse struct {
	ProblemDetails
	DeveloperHelp
	Instance interface{} `json:"instance,omitempty"`
} // HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"Server is running"`
}

// Helper functions for creating common error responses

// NewValidationError creates a validation error response with helpful info
func NewValidationError(errors []ValidationError) ValidationErrorResponse {
	return ValidationErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "validation_error",
			Title:  "Validation Failed",
			Status: 422,
			Detail: "One or more fields failed validation",
		},
		DeveloperHelp: DeveloperHelp{
			Help: "Please check the field requirements in the API documentation",
		},
		Errors: errors,
	}
}

// NewConflictError creates a conflict error response
func NewConflictError(detail string, instance interface{}, action string) ConflictErrorResponse {
	return ConflictErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "conflict",
			Title:  "Resource Conflict",
			Status: 409,
			Detail: detail,
		},
		DeveloperHelp: DeveloperHelp{
			Action: action,
		},
		Instance: instance,
	}
}

// NewBadRequestError creates a bad request error response
func NewBadRequestError(detail string) ErrorResponse {
	return ErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "bad_request",
			Title:  "Bad Request",
			Status: 400,
			Detail: detail,
		},
		// No DeveloperHelp for simple bad requests
	}
}

// NewInternalServerError creates a server error response
func NewInternalServerError(detail string) ErrorResponse {
	return ErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "internal_server_error",
			Title:  "Internal Server Error",
			Status: 500,
			Detail: detail,
		},
		DeveloperHelp: DeveloperHelp{
			Action: "Please try again later or contact support if the problem persists",
		},
	}
}

// NewNotFoundError creates a not found error response
func NewNotFoundError(detail string) ErrorResponse {
	return ErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "not_found",
			Title:  "Resource Not Found",
			Status: 404,
			Detail: detail,
		},
		// No DeveloperHelp for simple not found errors
	}
}

// NewRequiredFieldError creates a validation error for missing required fields
func NewRequiredFieldError(field string) ValidationErrorResponse {
	return NewValidationError([]ValidationError{{
		Field:   field,
		Message: fmt.Sprintf("%s is required", field),
		Code:    "REQUIRED",
	}})
}

// NewInvalidFormatError creates a validation error for invalid field formats
func NewInvalidFormatError(field, expectedFormat string) ValidationErrorResponse {
	return NewValidationError([]ValidationError{{
		Field:   field,
		Message: fmt.Sprintf("%s must be a valid %s", field, expectedFormat),
		Code:    "INVALID_FORMAT",
	}})
}

// NewForbiddenError creates a forbidden access error response
func NewForbiddenError(detail, action string) ErrorResponse {
	return ErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "forbidden",
			Title:  "Access Forbidden",
			Status: 403,
			Detail: detail,
		},
		DeveloperHelp: DeveloperHelp{
			Action: action,
		},
	}
}
