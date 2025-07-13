package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"gorm.io/gorm"
)

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

// GetPendingCommandCount gets the count of pending commands for a user
// Implements logic from USP_GetOutstanding stored procedure
func (cs *CommandService) GetPendingCommandCount(userID uuid.UUID) (int64, error) {
	var count int64

	// Delete commands from blocked users (from USP_GetOutstanding)
	cs.db.Where("sender_id IN (SELECT blockee_id FROM blocks WHERE blocker_id = ?) AND sub_id = ?", userID, userID).Delete(&models.ControlAppCmd{})

	// Check user's AnonCmd setting
	var user models.User
	err := cs.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return 0, fmt.Errorf("user not found: %w", err)
	}

	// If user doesn't allow anonymous commands, delete them (using special anonymous UUID)
	if !user.AnonCmd {
		// Use a special UUID that represents anonymous sender (sender_id = -1 in legacy)
		anonymousID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		cs.db.Where("sender_id = ? AND sub_id = ?", anonymousID, userID).Delete(&models.ControlAppCmd{})
	}

	// Count remaining commands
	err = cs.db.Model(&models.ControlAppCmd{}).Where("sub_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count commands: %w", err)
	}

	// Update user's login date
	cs.db.Model(&user).Update("login_date", time.Now())

	return count, nil
}

// GetNextCommand gets the next command for a user
// Implements logic from USP_GetAppContent stored procedure
func (cs *CommandService) GetNextCommand(userID uuid.UUID) (*models.ControlAppCmd, error) {
	// Delete commands from blocked users (from USP_GetAppContent)
	cs.db.Where("sender_id IN (SELECT blockee_id FROM blocks WHERE blocker_id = ?) AND sub_id = ?", userID, userID).Delete(&models.ControlAppCmd{})

	// Delete anonymous commands if user doesn't allow them
	var user models.User
	err := cs.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if !user.AnonCmd {
		// Use a special UUID that represents anonymous sender (sender_id = -1 in legacy)
		anonymousID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
		cs.db.Where("sender_id = ? AND sub_id = ?", anonymousID, userID).Delete(&models.ControlAppCmd{})
	}

	// Get the next command with relationships
	var assignment models.ControlAppCmd
	err = cs.db.Preload("Sender").Preload("Command").
		Where("sub_id = ?", userID).
		Order("created_at ASC").
		First(&assignment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No commands found
		}
		return nil, fmt.Errorf("failed to get next command: %w", err)
	}

	return &assignment, nil
}

// CompleteCommand marks a command as completed
// Implements logic from USP_CmdComplete stored procedure
func (cs *CommandService) CompleteCommand(userID uuid.UUID) error {
	// Delete the oldest command for this user
	var assignment models.ControlAppCmd
	err := cs.db.Where("sub_id = ?", userID).Order("created_at ASC").First(&assignment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No commands to complete
		}
		return fmt.Errorf("failed to find command to complete: %w", err)
	}

	// Delete the command assignment
	err = cs.db.Delete(&assignment).Error
	if err != nil {
		return fmt.Errorf("failed to delete command assignment: %w", err)
	}

	return nil
}

// MarkCommandCompleted marks a specific command as completed
func (cs *CommandService) MarkCommandCompleted(commandID uuid.UUID, userID uuid.UUID) error {
	// Find the command assignment
	var assignment models.ControlAppCmd
	err := cs.db.Where("command_id = ? AND sub_id = ?", commandID, userID).First(&assignment).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("command not found")
		}
		return fmt.Errorf("failed to find command: %w", err)
	}

	// Delete the command assignment
	err = cs.db.Delete(&assignment).Error
	if err != nil {
		return fmt.Errorf("failed to delete command assignment: %w", err)
	}

	return nil
}

// DeletePendingCommandsForUser deletes all pending commands for a user
// Implements logic from USP_DeleteOutstanding stored procedure
func (cs *CommandService) DeletePendingCommandsForUser(userID uuid.UUID) error {
	err := cs.db.Where("sub_id = ?", userID).Delete(&models.ControlAppCmd{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete pending commands: %w", err)
	}
	return nil
}

// CreateCommand creates a new command
func (cs *CommandService) CreateCommand(commandType, content, data string) (*models.Command, error) {
	command := models.Command{
		Type:    commandType,
		Content: content,
		Data:    data,
		Status:  "pending",
	}

	err := cs.db.Create(&command).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create command: %w", err)
	}

	return &command, nil
}

// AssignCommandToUser assigns a command to a user
func (cs *CommandService) AssignCommandToUser(senderID, recipientID, commandID uuid.UUID, groupRefID *uuid.UUID) (*models.ControlAppCmd, error) {
	assignment := models.ControlAppCmd{
		SenderID:   senderID,
		SubID:      recipientID,
		CommandID:  commandID,
		GroupRefID: groupRefID,
	}

	err := cs.db.Create(&assignment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to assign command: %w", err)
	}

	return &assignment, nil
}

// GetWhosNext gets the next sender info for a user (implements USP_GetOutstanding logic)
func (cs *CommandService) GetWhosNext(userID uuid.UUID) (string, error) {
	// Get the next command's sender or group
	var assignment models.ControlAppCmd
	err := cs.db.Preload("Sender").
		Where("sub_id = ?", userID).
		Order("created_at ASC").
		First(&assignment).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "User", nil // Default if no commands
		}
		return "", fmt.Errorf("failed to get next command: %w", err)
	}

	// Check if it's a group command
	if assignment.GroupRefID != nil {
		var group models.Group
		err := cs.db.First(&group, "id = ?", assignment.GroupRefID).Error
		if err != nil {
			return "Group", nil
		}
		return "Group: " + group.Name, nil
	}

	// Check if it's anonymous
	if assignment.SenderID.String() == "00000000-0000-0000-0000-000000000001" { // Special UUID for anonymous
		return "Anon", nil
	}

	// Check if sender is in a relationship with the user
	var relationship models.Relationship
	err = cs.db.Where("dom_id = ? AND sub_id = ?", assignment.SenderID, userID).First(&relationship).Error
	if err == nil {
		return assignment.Sender.ScreenName, nil
	}

	return "User", nil
}

// AcceptInvite accepts an invitation from a dom (implements USP_AcceptInvite logic)
func (cs *CommandService) AcceptInvite(subID uuid.UUID, domName string) error {
	var domUser models.User
	err := cs.db.Where("screen_name = ?", domName).First(&domUser).Error
	if err != nil {
		return fmt.Errorf("dom user not found: %w", err)
	}

	// Create relationship
	relationship := models.Relationship{
		UserID:    domUser.ID,
		RelatedID: subID,
		Type:      "control",
		Status:    "active",
	}

	err = cs.db.Create(&relationship).Error
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	// Delete the invite
	err = cs.db.Where("dom_id = ? AND sub_id = ?", domUser.ID, subID).Delete(&models.Invite{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

// RejectInvite rejects an invitation from a dom (implements USP_DeleteInvite logic)
func (cs *CommandService) RejectInvite(subID uuid.UUID, domName string) error {
	var domUser models.User
	err := cs.db.Where("screen_name = ?", domName).First(&domUser).Error
	if err != nil {
		return fmt.Errorf("dom user not found: %w", err)
	}

	// Delete the invite
	err = cs.db.Where("dom_id = ? AND sub_id = ?", domUser.ID, subID).Delete(&models.Invite{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

// GiveThumbsUp gives a thumbs up to a sender (implements USP_thumbsup logic)
func (cs *CommandService) GiveThumbsUp(userID uuid.UUID, senderID uuid.UUID) error {
	if userID == senderID {
		return fmt.Errorf("cannot give thumbs up to yourself")
	}

	// Check if sender exists
	var sender models.User
	err := cs.db.First(&sender, "id = ?", senderID).Error
	if err != nil {
		return fmt.Errorf("sender not found: %w", err)
	}

	// Increment thumbs up count
	err = cs.db.Model(&sender).Update("thumbs_up", gorm.Expr("thumbs_up + 1")).Error
	if err != nil {
		return fmt.Errorf("failed to update thumbs up: %w", err)
	}

	return nil
}
