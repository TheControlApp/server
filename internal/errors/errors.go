package errors

import (
	"fmt"
	"net/http"
)

// AppError represents a structured application error with HTTP context
type AppError struct {
	Type     string            `json:"type"`
	Title    string            `json:"title"`
	Status   int               `json:"status"`
	Detail   string            `json:"detail"`
	Action   string            `json:"action,omitempty"`
	Fields   []FieldError      `json:"errors,omitempty"`
	Instance map[string]string `json:"instance,omitempty"`
}

// FieldError represents a validation error on a specific field
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Title, e.Detail)
}

// Common constructors
func BadRequest(detail string) *AppError {
	return &AppError{Type: "bad_request", Title: "Bad Request", Status: http.StatusBadRequest, Detail: detail}
}

func Unauthorized(detail, action string) *AppError {
	return &AppError{Type: "unauthorized", Title: "Unauthorized Access", Status: http.StatusUnauthorized, Detail: detail, Action: action}
}

func NotFound(detail string) *AppError {
	return &AppError{Type: "not_found", Title: "Resource Not Found", Status: http.StatusNotFound, Detail: detail}
}

func Conflict(detail, action string, instance map[string]string) *AppError {
	return &AppError{Type: "conflict", Title: "Resource Conflict", Status: http.StatusConflict, Detail: detail, Action: action, Instance: instance}
}

func ValidationFailed(fields []FieldError) *AppError {
	return &AppError{Type: "validation_error", Title: "Validation Failed", Status: http.StatusUnprocessableEntity, Detail: "One or more fields failed validation", Fields: fields}
}

func Internal(detail string) *AppError {
	return &AppError{Type: "internal_error", Title: "Internal Server Error", Status: http.StatusInternalServerError, Detail: detail}
}

// Validation helpers
func ValidateField(field, value string, min, max int) *FieldError {
	if len(value) < min {
		return &FieldError{Field: field, Message: fmt.Sprintf("%s must be at least %d characters", field, min), Code: "MIN_LENGTH"}
	}
	if len(value) > max {
		return &FieldError{Field: field, Message: fmt.Sprintf("%s must be no more than %d characters", field, max), Code: "MAX_LENGTH"}
	}
	return nil
}
