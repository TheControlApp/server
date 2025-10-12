package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/auth"
	"github.com/thecontrolapp/server/internal/models"
	"gorm.io/gorm"
)

// Standard error types following HTTP conventions
var (
	ErrUserNotFound = errors.New("user not found")
	ErrUnauthorized = errors.New("authentication failed")
	ErrConflict     = errors.New("resource conflict")
	ErrValidation   = errors.New("validation error")
	ErrBadRequest   = errors.New("bad request")
)

// ValidationError represents a field validation error with standard code
type ValidationError struct {
	Field   string
	Message string
	Code    string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ConflictError represents a resource conflict
type ConflictError struct {
	Resource string
	Value    string
	Message  string
}

// Error implements the error interface
func (e *ConflictError) Error() string {
	return e.Message
}

// ValidateUsername checks if a username meets requirements
func ValidateUsername(username string) *ValidationError {
	if len(username) < 3 {
		return &ValidationError{
			Field:   "username",
			Message: "Username must be at least 3 characters long",
			Code:    "MIN_LENGTH", // More intuitive than TOO_SHORT
		}
	}
	if len(username) > 50 {
		return &ValidationError{
			Field:   "username",
			Message: "Username must be no more than 50 characters long",
			Code:    "MAX_LENGTH", // More intuitive than TOO_LONG
		}
	}
	if strings.TrimSpace(username) != username {
		return &ValidationError{
			Field:   "username",
			Message: "Username cannot have leading or trailing spaces",
			Code:    "INVALID_FORMAT",
		}
	}
	return nil
}

// ValidatePassword checks if a password meets requirements
func ValidatePassword(password string) *ValidationError {
	if len(password) < 6 {
		return &ValidationError{
			Field:   "password",
			Message: "Password must be at least 6 characters long",
			Code:    "MIN_LENGTH", // Consistent with username
		}
	}
	if len(password) > 128 {
		return &ValidationError{
			Field:   "password",
			Message: "Password must be no more than 128 characters long",
			Code:    "MAX_LENGTH", // Consistent with username
		}
	}
	return nil
}

// UserService handles user-related operations
type UserService struct {
	db   *gorm.DB
	Auth *auth.AuthService
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB, authService *auth.AuthService) *UserService {
	return &UserService{
		db:   db,
		Auth: authService,
	}
}

// AuthenticateUser authenticates a user with username and password
func (us *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User

	// Try to find user by login name or screen name
	err := us.db.Where("login_name = ? OR screen_name = ?", username, username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	err = us.Auth.PasswordManager.VerifyPassword(password, user.Password)
	if err != nil {
		return nil, ErrUnauthorized
	}

	// Update login date
	user.LoginDate = time.Now()
	us.db.Save(&user)

	return &user, nil
}

// CreateUserRequest is used for creating a new user via modern API
type CreateUserRequest struct {
	LoginName   string `json:"login_name" binding:"required"`
	ScreenName  string `json:"screen_name" binding:"required"`
	Password    string `json:"password" binding:"required"`
	RandomOptIn bool   `json:"random_opt_in" binding:"required"`
}

// CreateUser creates a new user with the modern API
func (us *UserService) CreateUser(req CreateUserRequest) (*models.User, error) {
	// Validate username requirements
	if err := ValidateUsername(req.LoginName); err != nil {
		return nil, err
	}

	// Validate screen name (same rules as username for now)
	if err := ValidateUsername(req.ScreenName); err != nil {
		return nil, &ValidationError{
			Field:   "screen_name",
			Message: err.Message,
			Code:    err.Code,
		}
	}

	// Validate password requirements
	if err := ValidatePassword(req.Password); err != nil {
		return nil, err
	}

	// Check if username already exists (either login_name or screen_name)
	var existingUser models.User
	err := us.db.Where("login_name = ? OR screen_name = ?", req.LoginName, req.ScreenName).First(&existingUser).Error
	if err == nil {
		// User exists
		return nil, &ConflictError{
			Resource: "username",
			Value:    req.LoginName,
			Message:  "Username already exists",
		}
	} else if err != gorm.ErrRecordNotFound {
		// Database error
		return nil, fmt.Errorf("database error while checking username: %w", err)
	}

	hashedPassword, err := us.Auth.PasswordManager.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		LoginName:   req.LoginName,
		ScreenName:  req.ScreenName,
		Password:    hashedPassword,
		RandomOptIn: req.RandomOptIn,
		Role:        "user",
	}
	err = us.db.Create(&user).Error
	if err != nil {
		// Check for database constraint violations
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
			strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, &ConflictError{
				Resource: "username",
				Value:    req.LoginName,
				Message:  "Username already exists",
			}
		}
		return nil, fmt.Errorf("database error during user creation: %w", err)
	}
	return &user, nil
} // GetAllUsers returns all users
func (us *UserService) GetAllUsers() ([]models.User, error) {
	var users []models.User
	err := us.db.Find(&users).Error
	return users, err
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	err := us.db.First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username (login name or screen name)
func (us *UserService) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := us.db.Where("login_name = ? OR screen_name = ?", username, username).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return &user, nil
}
