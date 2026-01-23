package services

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/auth"
	"github.com/thecontrolapp/server/internal/errors"
	"github.com/thecontrolapp/server/internal/models"
	"gorm.io/gorm"
)

type UserService struct {
	db   *gorm.DB
	Auth *auth.AuthService
}

func NewUserService(db *gorm.DB, authService *auth.AuthService) *UserService {
	return &UserService{db: db, Auth: authService}
}

type CreateUserRequest struct {
	LoginName   string `json:"login_name"`
	ScreenName  string `json:"screen_name"`
	Password    string `json:"password"`
	RandomOptIn bool   `json:"random_opt_in"`
}

// AuthenticateUser authenticates a user and returns user + token
func (us *UserService) AuthenticateUser(username, password string) (*models.User, string, *errors.AppError) {
	var user models.User
	if err := us.db.Where("login_name = ? OR screen_name = ?", username, username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, "", errors.Unauthorized("Invalid username or password", "Please check your credentials and try again")
		}
		return nil, "", errors.Internal("Database error during authentication")
	}

	if err := us.Auth.PasswordManager.VerifyPassword(password, user.Password); err != nil {
		return nil, "", errors.Unauthorized("Invalid username or password", "Please check your credentials and try again")
	}

	token, err := us.Auth.JWTManager.GenerateToken(user.ID)
	if err != nil {
		return nil, "", errors.Internal("Failed to generate authentication token")
	}

	user.LoginDate = time.Now()
	us.db.Save(&user)

	return &user, token, nil
}

// CreateUser creates a new user with validation
func (us *UserService) CreateUser(req CreateUserRequest) (*models.User, *errors.AppError) {
	// Validate inputs
	var validationErrors []errors.FieldError
	if err := errors.ValidateField("username", req.LoginName, 3, 50); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := errors.ValidateField("screen_name", req.ScreenName, 3, 50); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if err := errors.ValidateField("password", req.Password, 6, 128); err != nil {
		validationErrors = append(validationErrors, *err)
	}
	if len(validationErrors) > 0 {
		return nil, errors.ValidationFailed(validationErrors)
	}

	// Check for existing user
	var existingUser models.User
	err := us.db.Where("login_name = ? OR screen_name = ?", req.LoginName, req.ScreenName).First(&existingUser).Error
	if err == nil {
		return nil, errors.Conflict("Username already exists", "Please choose a different username and try again", map[string]string{"username": req.LoginName})
	} else if err != gorm.ErrRecordNotFound {
		return nil, errors.Internal("Database error while checking username")
	}

	hashedPassword, err := us.Auth.PasswordManager.HashPassword(req.Password)
	if err != nil {
		return nil, errors.Internal("Failed to hash password")
	}

	user := models.User{
		LoginName:   req.LoginName,
		ScreenName:  req.ScreenName,
		Password:    hashedPassword,
		RandomOptIn: req.RandomOptIn,
		Role:        "user",
	}

	if err := us.db.Create(&user).Error; err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return nil, errors.Conflict("Username already exists", "Please choose a different username", map[string]string{"username": req.LoginName})
		}
		return nil, errors.Internal("Database error during user creation")
	}

	return &user, nil
}

// GetAllUsers returns all users
func (us *UserService) GetAllUsers() ([]models.User, *errors.AppError) {
	var users []models.User
	if err := us.db.Find(&users).Error; err != nil {
		return nil, errors.Internal("Database error retrieving users")
	}
	return users, nil
}

// GetUserByID retrieves a user by ID
func (us *UserService) GetUserByID(id uuid.UUID) (*models.User, *errors.AppError) {
	var user models.User
	if err := us.db.First(&user, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("User not found")
		}
		return nil, errors.Internal("Database error retrieving user")
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (us *UserService) GetUserByUsername(username string) (*models.User, *errors.AppError) {
	var user models.User
	if err := us.db.Where("login_name = ? OR screen_name = ?", username, username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NotFound("User not found")
		}
		return nil, errors.Internal("Database error retrieving user")
	}
	return &user, nil
}
