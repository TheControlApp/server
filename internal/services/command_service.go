package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/server/internal/models"
	"gorm.io/gorm"
)

type CommandService struct {
	db *gorm.DB
}

func NewCommandService(db *gorm.DB) *CommandService {
	return &CommandService{db: db}
}

// CreateCommand stores a new command in the database
func (cs *CommandService) CreateCommand(command *models.Command) error {
	if command.ID == uuid.Nil {
		command.ID = uuid.New()
	}

	if command.Status == "" {
		command.Status = "pending"
	}

	command.CreatedAt = time.Now()
	command.UpdatedAt = time.Now()

	return cs.db.Create(command).Error
}

// GetPendingCommandsForUser retrieves all pending commands for a specific user or broadcast commands
// Returns commands in FIFO order (oldest first)
func (cs *CommandService) GetPendingCommandsForUser(userID uuid.UUID) ([]models.Command, error) {
	var commands []models.Command

	// Get commands that are either:
	// 1. Targeted to this specific user (receiver_id = userID)
	// 2. Broadcast commands (receiver_id IS NULL)
	// Order by created_at ASC for FIFO
	err := cs.db.Where("receiver_id = ? OR receiver_id IS NULL", userID).
		Where("status = ?", "pending").
		Order("created_at ASC").
		Find(&commands).Error

	return commands, err
}

// GetPendingCommands (legacy compatibility) - retrieves pending commands for specific user only
func (cs *CommandService) GetPendingCommands(userID uuid.UUID) ([]models.Command, error) {
	var commands []models.Command
	err := cs.db.Where("receiver_id = ? AND status = ?", userID, "pending").
		Preload("Sender").
		Preload("Receiver").
		Order("created_at ASC").
		Find(&commands).Error
	return commands, err
}

func (cs *CommandService) GetPendingCommandCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := cs.db.Model(&models.Command{}).
		Where("(receiver_id = ? OR receiver_id IS NULL) AND status = ?", userID, "pending").
		Count(&count).Error
	return count, err
}

// MarkCommandDelivered marks a command as delivered and deletes if targeted
func (cs *CommandService) MarkCommandDelivered(commandID uuid.UUID, userID uuid.UUID) error {
	var command models.Command
	if err := cs.db.First(&command, "id = ?", commandID).Error; err != nil {
		return err
	}

	// If command is targeted to specific user, delete it
	if command.ReceiverID != nil && *command.ReceiverID == userID {
		return cs.db.Delete(&command).Error
	}

	// For broadcast commands, update status to delivered
	return cs.db.Model(&command).Update("status", "delivered").Error
}

func (cs *CommandService) CompleteCommand(commandID uuid.UUID, userID uuid.UUID) error {
	result := cs.db.Model(&models.Command{}).Where("id = ? AND receiver_id = ?", commandID, userID).Update("status", "completed")
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("command not found")
	}
	return nil
}

// GetCommandByID gets a command by ID with relationships loaded
func (cs *CommandService) GetCommandByID(commandID uuid.UUID) (*models.Command, error) {
	var command models.Command
	err := cs.db.Where("id = ?", commandID).
		Preload("Sender").
		Preload("Receiver").
		First(&command).Error
	if err != nil {
		return nil, fmt.Errorf("command not found: %w", err)
	}
	return &command, nil
}
