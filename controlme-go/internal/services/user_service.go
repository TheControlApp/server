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
	auth *auth.AuthService
}

// NewUserService creates a new user service
func NewUserService(db *gorm.DB, authService *auth.AuthService) *UserService {
	return &UserService{
		db:   db,
		auth: authService,
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
	err = us.auth.PasswordManager.VerifyPassword(password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	// Update login date
	user.LoginDate = time.Now()
	us.db.Save(&user)

	return &user, nil
}

// AuthenticateLegacyUser authenticates a user with legacy encrypted password
func (us *UserService) AuthenticateLegacyUser(username, encryptedPassword string) (*models.User, error) {
	// Decrypt the password using legacy crypto
	decryptedPassword, err := us.auth.LegacyCrypto.Decrypt(encryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt password: %w", err)
	}

	// Authenticate with decrypted password
	return us.AuthenticateUser(username, decryptedPassword)
}

// CreateUserRequest is used for creating a new user via modern API
type CreateUserRequest struct {
	LoginName  string `json:"login_name" binding:"required"`
	ScreenName string `json:"screen_name" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// CreateUser creates a new user with the modern API
func (us *UserService) CreateUser(req CreateUserRequest) (*models.User, error) {
	user := models.User{
		LoginName:  req.LoginName,
		ScreenName: req.ScreenName,
		Password:   req.Password, // Use modern hash
	}
	err := us.db.Create(&user).Error
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

// CommandService handles command-related operations
type CommandService struct {
	db *gorm.DB
}

// NewCommandService creates a new command service
func NewCommandService(db *gorm.DB) *CommandService {
	return &CommandService{
		db: db,
	}
}

// CreateCommand creates a new command
func (cs *CommandService) CreateCommand(commandType, content, data string) (*models.Command, error) {
	command := models.Command{
		ID:        uuid.New(),
		Type:      commandType,
		Content:   content,
		Data:      data,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := cs.db.Create(&command).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	return &command, nil
}

// AssignCommandToUser assigns a command to a user (creates ControlAppCmd)
func (cs *CommandService) AssignCommandToUser(senderID, subID, commandID uuid.UUID, groupRefID *uuid.UUID) (*models.ControlAppCmd, error) {
	assignment := models.ControlAppCmd{
		ID:         uuid.New(),
		SenderID:   senderID,
		SubID:      subID,
		CommandID:  commandID,
		GroupRefID: groupRefID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	err := cs.db.Create(&assignment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to assign command: %w", err)
	}

	return &assignment, nil
}

// GetPendingCommandsForUser retrieves pending commands for a specific user
func (cs *CommandService) GetPendingCommandsForUser(userID uuid.UUID) ([]models.ControlAppCmd, error) {
	var assignments []models.ControlAppCmd

	err := cs.db.Preload("Command").Preload("Sender").
		Where("sub_id = ? AND commands.status = ?", userID, "pending").
		Joins("JOIN commands ON control_app_cmds.command_id = commands.id").
		Find(&assignments).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get pending commands: %w", err)
	}

	return assignments, nil
}

// GetPendingCommandCount gets the count of pending commands for a user
func (cs *CommandService) GetPendingCommandCount(userID uuid.UUID) (int64, error) {
	var count int64

	err := cs.db.Model(&models.ControlAppCmd{}).
		Where("sub_id = ? AND EXISTS (SELECT 1 FROM commands WHERE commands.id = control_app_cmds.command_id AND commands.status = ?)",
			userID, "pending").
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count pending commands: %w", err)
	}

	return count, nil
}

// GetNextCommand gets the next command for a user (simulates USP_GetAppContent)
func (cs *CommandService) GetNextCommand(userID uuid.UUID) (*models.ControlAppCmd, error) {
	var assignment models.ControlAppCmd

	err := cs.db.Preload("Command").Preload("Sender").
		Where("sub_id = ? AND EXISTS (SELECT 1 FROM commands WHERE commands.id = control_app_cmds.command_id AND commands.status = ?)",
			userID, "pending").
		Order("created_at ASC").
		First(&assignment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No pending commands
		}
		return nil, fmt.Errorf("failed to get next command: %w", err)
	}

	return &assignment, nil
}

// CompleteCommand marks a command as completed
func (cs *CommandService) CompleteCommand(userID uuid.UUID) error {
	// Get the oldest pending command for this user
	var assignment models.ControlAppCmd

	err := cs.db.Where("sub_id = ? AND EXISTS (SELECT 1 FROM commands WHERE commands.id = control_app_cmds.command_id AND commands.status = ?)",
		userID, "pending").
		Order("created_at ASC").
		First(&assignment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No pending commands to complete
		}
		return fmt.Errorf("failed to find command to complete: %w", err)
	}

	// Mark the command as completed
	err = cs.db.Model(&models.Command{}).
		Where("id = ?", assignment.CommandID).
		Update("status", "completed").Error

	if err != nil {
		return fmt.Errorf("failed to mark command as completed: %w", err)
	}

	return nil
}

// DeleteOutstandingCommands deletes all outstanding commands for a user
func (cs *CommandService) DeleteOutstandingCommands(userID uuid.UUID) error {
	// Delete all pending command assignments for this user
	err := cs.db.Where("sub_id = ?", userID).Delete(&models.ControlAppCmd{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete outstanding commands: %w", err)
	}

	return nil
}

// GetUserRelationships gets the relationships for a user (dom/sub relationships)
func (cs *CommandService) GetUserRelationships(userID uuid.UUID) ([]models.Relationship, error) {
	var relationships []models.Relationship

	err := cs.db.Preload("Dom").Preload("Sub").
		Where("dom_id = ? OR sub_id = ?", userID, userID).
		Find(&relationships).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get relationships: %w", err)
	}

	return relationships, nil
}

// MarkCommandCompleted marks a specific command as completed
func (cs *CommandService) MarkCommandCompleted(commandID uuid.UUID, userID uuid.UUID) error {
	// Verify the command belongs to this user and update status
	err := cs.db.Model(&models.Command{}).
		Where("id = ? AND EXISTS (SELECT 1 FROM control_app_cmds WHERE control_app_cmds.command_id = ? AND control_app_cmds.sub_id = ?)",
			commandID, commandID, userID).
		Update("status", "completed").Error

	if err != nil {
		return fmt.Errorf("failed to mark command as completed: %w", err)
	}

	return nil
}

// DeletePendingCommandsForUser deletes all pending commands for a user
func (cs *CommandService) DeletePendingCommandsForUser(userID uuid.UUID) error {
	// First get all pending command IDs for this user
	var commandIDs []uuid.UUID
	err := cs.db.Model(&models.ControlAppCmd{}).
		Select("command_id").
		Where("sub_id = ?", userID).
		Pluck("command_id", &commandIDs).Error

	if err != nil {
		return fmt.Errorf("failed to get command IDs: %w", err)
	}

	// Delete the command assignments
	err = cs.db.Where("sub_id = ?", userID).Delete(&models.ControlAppCmd{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete command assignments: %w", err)
	}

	// Update command status to cancelled for these commands
	if len(commandIDs) > 0 {
		err = cs.db.Model(&models.Command{}).
			Where("id IN ? AND status = ?", commandIDs, "pending").
			Update("status", "cancelled").Error
		if err != nil {
			return fmt.Errorf("failed to update command status: %w", err)
		}
	}

	return nil
}
