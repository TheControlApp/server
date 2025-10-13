package responses

import (
	"fmt"

	"github.com/thecontrolapp/server/internal/models"
)

// PaymentRequiredErrorResponse represents a payment required error (HTTP 402)
type PaymentRequiredErrorResponse struct {
	ProblemDetails
	DeveloperHelp
	PricingURL string `json:"pricing_url,omitempty"`
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

// ErrorResponse represents a simple error response
type ErrorResponse struct {
	Error string `json:"error" example:"Internal server error"`
}

// ValidationError represents a field-specific validation error
type ValidationError struct {
	Field   string `json:"field" example:"username"`
	Message string `json:"message" example:"Username must be at least 3 characters"`
	Code    string `json:"code" example:"MIN_LENGTH"`
}

// SimpleErrorResponse represents a simplified error response for basic use cases
type SimpleErrorResponse struct {
	Error   string `json:"error" example:"Username is too short"`
	Code    string `json:"code,omitempty" example:"VALIDATION_ERROR"`
	Message string `json:"message,omitempty" example:"Please check your input and try again"`
}

// ValidationErrorResponse represents validation errors for form fields (HTTP 422)
type ValidationErrorResponse struct {
	ProblemDetails
	DeveloperHelp
	Errors []ValidationError `json:"errors"`
}

// ConflictErrorResponse represents resource conflict errors (HTTP 409)
type ConflictErrorResponse struct {
	ProblemDetails
	DeveloperHelp
	Instance interface{} `json:"instance,omitempty"`
}

// UnauthorizedErrorResponse represents authentication errors (HTTP 401)
type UnauthorizedErrorResponse struct {
	ProblemDetails
	DeveloperHelp
}

// BadRequestErrorResponse represents malformed request errors (HTTP 400)
type BadRequestErrorResponse struct {
	ProblemDetails
}

// NotFoundErrorResponse represents not found errors (HTTP 404)
type NotFoundErrorResponse struct {
	ProblemDetails
}

// ForbiddenErrorResponse represents forbidden access errors (HTTP 403)
type ForbiddenErrorResponse struct {
	ProblemDetails
	DeveloperHelp
}

// InternalServerErrorResponse represents internal server errors (HTTP 500)
type InternalServerErrorResponse struct {
	ProblemDetails
	DeveloperHelp
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

// NewUnauthorizedError creates an unauthorized error response
func NewUnauthorizedError(detail, action string) UnauthorizedErrorResponse {
	return UnauthorizedErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "unauthorized",
			Title:  "Authentication Failed",
			Status: 401,
			Detail: detail,
		},
		DeveloperHelp: DeveloperHelp{
			Action: action,
		},
	}
}

// NewBadRequestError creates a bad request error response
func NewBadRequestError(detail string) BadRequestErrorResponse {
	return BadRequestErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "bad_request",
			Title:  "Bad Request",
			Status: 400,
			Detail: detail,
		},
	}
}

// NewInternalServerError creates a server error response
func NewInternalServerError(detail string) InternalServerErrorResponse {
	return InternalServerErrorResponse{
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
func NewNotFoundError(detail string) NotFoundErrorResponse {
	return NotFoundErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "not_found",
			Title:  "Resource Not Found",
			Status: 404,
			Detail: detail,
		},
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
func NewForbiddenError(detail, action string) ForbiddenErrorResponse {
	return ForbiddenErrorResponse{
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

// NewPaymentRequiredError creates a payment required error response
func NewPaymentRequiredError(detail, action, pricingURL string) PaymentRequiredErrorResponse {
	return PaymentRequiredErrorResponse{
		ProblemDetails: ProblemDetails{
			Type:   "payment_required",
			Title:  "Payment Required",
			Status: 402,
			Detail: detail,
		},
		DeveloperHelp: DeveloperHelp{
			Action: action,
		},
		PricingURL: pricingURL,
	}
}
