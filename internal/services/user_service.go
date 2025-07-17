package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"gorm.io/gorm"
)

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
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	err = us.Auth.PasswordManager.VerifyPassword(password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password")
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
	Email       string `json:"email" binding:"required"`
	RandomOptIn bool   `json:"random_opt_in" binding:"required"`
}

// CreateUser creates a new user with the modern API
func (us *UserService) CreateUser(req CreateUserRequest) (*models.User, error) {
	hashedPassword, err := us.Auth.PasswordManager.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		LoginName:   req.LoginName,
		ScreenName:  req.ScreenName,
		Password:    hashedPassword,
		Email:       req.Email,
		RandomOptIn: req.RandomOptIn,
		Role:        "user",
	}
	err = us.db.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers returns all users
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
