package services

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"gorm.io/gorm"
)

type CommandService struct {
	db *gorm.DB
}

func NewCommandService(db *gorm.DB) *CommandService {
	return &CommandService{db: db}
}

func (cs *CommandService) GetPendingCommands(userID uuid.UUID) ([]models.Command, error) {
	var commands []models.Command
	err := cs.db.Where("receiver_id = ? AND status = ?", userID, "pending").
		Preload("Sender").
		Preload("Receiver").
		Find(&commands).Error
	return commands, err
}

func (cs *CommandService) GetPendingCommandCount(userID uuid.UUID) (int64, error) {
	var count int64
	err := cs.db.Model(&models.Command{}).Where("receiver_id = ? AND status = ?", userID, "pending").Count(&count).Error
	return count, err
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
