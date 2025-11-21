package client

import (
	"errors"
	"fmt"
)

// Common errors
var (
	ErrNotConnected      = errors.New("not connected to server")
	ErrNotAuthenticated  = errors.New("not authenticated")
	ErrInvalidConfig     = errors.New("invalid configuration")
	ErrCommandTimeout    = errors.New("command execution timeout")
	ErrCommandNotAllowed = errors.New("command not allowed by user consent")
	ErrInvalidCommand    = errors.New("invalid command format")
	ErrFileNotFound      = errors.New("file not found")
	ErrDownloadFailed    = errors.New("file download failed")
	ErrPermissionDenied  = errors.New("permission denied")
)

// CommandError represents a command execution error
type CommandError struct {
	CommandID   string `json:"command_id"`
	CommandType string `json:"command_type"`
	Code        string `json:"code"`
	Message     string `json:"message"`
	Recoverable bool   `json:"recoverable"`
}

func (e *CommandError) Error() string {
	return fmt.Sprintf("command %s (%s) failed: %s", e.CommandType, e.CommandID, e.Message)
}

// NewCommandError creates a new command error
func NewCommandError(commandID, commandType, code, message string, recoverable bool) *CommandError {
	return &CommandError{
		CommandID:   commandID,
		CommandType: commandType,
		Code:        code,
		Message:     message,
		Recoverable: recoverable,
	}
}

// ConsentError represents a consent-related error
type ConsentError struct {
	CommandType string `json:"command_type"`
	Reason      string `json:"reason"`
}

func (e *ConsentError) Error() string {
	return fmt.Sprintf("consent denied for command %s: %s", e.CommandType, e.Reason)
}

// NetworkError represents a network-related error
type NetworkError struct {
	Operation string `json:"operation"`
	URL       string `json:"url,omitempty"`
	Cause     error  `json:"-"`
}

func (e *NetworkError) Error() string {
	if e.URL != "" {
		return fmt.Sprintf("network error during %s to %s: %v", e.Operation, e.URL, e.Cause)
	}
	return fmt.Sprintf("network error during %s: %v", e.Operation, e.Cause)
}

func (e *NetworkError) Unwrap() error {
	return e.Cause
}
