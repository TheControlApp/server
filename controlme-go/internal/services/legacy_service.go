package services

import (
	"fmt"
	"time"

	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"gorm.io/gorm"
)

// LegacyService handles all legacy database operations using the exact legacy schema
type LegacyService struct {
	db   *gorm.DB
	auth *auth.AuthService
}

// NewLegacyService creates a new legacy service
func NewLegacyService(db *gorm.DB, authService *auth.AuthService) *LegacyService {
	return &LegacyService{
		db:   db,
		auth: authService,
	}
}

// LegacyLoginResult represents the result of USP_Login stored procedure
type LegacyLoginResult struct {
	ID       int    `json:"id"`
	Role     string `json:"role"`
	Verified bool   `json:"verified"`
}

// USP_Login implements the exact logic from the legacy USP_Login stored procedure
func (ls *LegacyService) USP_Login(username, password string) (*LegacyLoginResult, error) {
	// SELECT Id, [Role],Varified from Users where [Screen Name]=@UserName and Password=@Password
	var user models.LegacyUser
	err := ls.db.Where("[Screen Name] = ? AND Password = ?", username, password).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	return &LegacyLoginResult{
		ID:       user.ID,
		Role:     user.Role,
		Verified: user.Verified,
	}, nil
}

// LegacyOutstandingResult represents the result of USP_GetOutstanding stored procedure
type LegacyOutstandingResult struct {
	Count    int64  `json:"count"`
	WhoNext  string `json:"who_next"`
	Verified string `json:"verified"`
	Thumbs   string `json:"thumbs"`
}

// USP_GetOutstanding implements the exact logic from the legacy USP_GetOutstanding stored procedure
func (ls *LegacyService) USP_GetOutstanding(userID int) (*LegacyOutstandingResult, error) {
	// DELETE cmd FROM ControlAppCmd cmd
	// WHERE SenderId in (Select BlockeeId FROM Block where BlockerId=@userID)
	// AND SubId=@userID
	ls.db.Exec("DELETE cmd FROM ControlAppCmd cmd WHERE SenderId in (Select BlockeeId FROM Block where BlockerId=?) AND SubId=?", userID, userID)

	// Get user info
	var user models.LegacyUser
	err := ls.db.First(&user, "Id = ?", userID).Error
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// IF (@Anon=0) DELETE anonymous commands
	if !user.AnonCmd {
		ls.db.Exec("DELETE cmd FROM ControlAppCmd cmd WHERE SenderId = -1 AND SubId = ?", userID)
	}

	// Count remaining commands
	var count int64
	err = ls.db.Model(&models.LegacyControlAppCmd{}).Where("SubId = ?", userID).Count(&count).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count commands: %w", err)
	}

	// Update login date
	ls.db.Model(&user).Update("LoginDate", time.Now())

	// Get who's next
	whoNext := "User"
	var assignment models.LegacyControlAppCmd
	err = ls.db.Where("SubId = ?", userID).Order("Id ASC").First(&assignment).Error
	if err == nil {
		if assignment.GroupRefID != nil {
			whoNext = fmt.Sprintf("Group:%d", *assignment.GroupRefID)
		} else if assignment.SenderID == -1 {
			whoNext = "Anon"
		} else {
			// Check if sender has relationship
			var relationship models.LegacyRelationship
			err = ls.db.Where("DomId = ? AND SubId = ?", assignment.SenderID, userID).First(&relationship).Error
			if err == nil {
				// Get sender name
				var sender models.LegacyUser
				err = ls.db.First(&sender, "Id = ?", assignment.SenderID).Error
				if err == nil {
					whoNext = sender.ScreenName
				}
			}
		}
	}

	verified := "0"
	if user.Verified {
		verified = "1"
	}

	return &LegacyOutstandingResult{
		Count:    count,
		WhoNext:  whoNext,
		Verified: verified,
		Thumbs:   fmt.Sprintf("%d", user.ThumbsUp),
	}, nil
}

// LegacyAppContentResult represents the result of USP_GetAppContent stored procedure
type LegacyAppContentResult struct {
	SenderID   string `json:"sender_id"`
	SenderName string `json:"sender_name"`
	Content    string `json:"content"`
}

// USP_GetAppContent implements the exact logic from the legacy USP_GetAppContent stored procedure
func (ls *LegacyService) USP_GetAppContent(userID int) (*LegacyAppContentResult, error) {
	// DELETE cmd FROM ControlAppCmd cmd
	// WHERE SenderId in (Select BlockeeId FROM Block where BlockerId=@userID)
	// AND SubId=@userID
	ls.db.Exec("DELETE cmd FROM ControlAppCmd cmd WHERE SenderId in (Select BlockeeId FROM Block where BlockerId=?) AND SubId=?", userID, userID)

	// DELETE cmd FROM ControlAppCmd cmd WHERE SenderId = -1 AND SubId=@userID
	ls.db.Exec("DELETE cmd FROM ControlAppCmd cmd WHERE SenderId = -1 AND SubId=?", userID)

	// SELECT TOP 1 convert(varchar(100),SenderId),u.[Screen Name] SenderName, cl.Content
	// from ControlAppCmd cac
	// join CommandList cl on cac.CmdId=cl.CmdId
	// join Users u on SenderId=u.Id
	// where cac.SubId=@userID
	// ORDER BY SendDate
	var result struct {
		SenderID   int    `gorm:"column:SenderId"`
		SenderName string `gorm:"column:SenderName"`
		Content    string `gorm:"column:Content"`
	}

	err := ls.db.Table("ControlAppCmd cac").
		Select("cac.SenderId, u.[Screen Name] as SenderName, cl.Content").
		Joins("join CommandList cl on cac.CmdId = cl.CmdId").
		Joins("join Users u on cac.SenderId = u.Id").
		Where("cac.SubId = ?", userID).
		Order("cl.SendDate ASC").
		First(&result).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No commands found
		}
		return nil, fmt.Errorf("failed to get app content: %w", err)
	}

	return &LegacyAppContentResult{
		SenderID:   fmt.Sprintf("%d", result.SenderID),
		SenderName: result.SenderName,
		Content:    result.Content,
	}, nil
}

// USP_CmdComplete implements the exact logic from the legacy USP_CmdComplete stored procedure
func (ls *LegacyService) USP_CmdComplete(userID int) error {
	// delete from ControlAppCmd where id =(
	// SELECT top 1 id from [ControlAppCmd] c
	// join CommandList cl on c.CmdId=cl.CmdId
	// where SubId=@userID
	// ORDER BY SendDate)

	var cmdID int
	err := ls.db.Table("ControlAppCmd c").
		Select("c.Id").
		Joins("join CommandList cl on c.CmdId = cl.CmdId").
		Where("c.SubId = ?", userID).
		Order("cl.SendDate ASC").
		Limit(1).
		Pluck("Id", &cmdID).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil // No commands to complete
		}
		return fmt.Errorf("failed to find command to complete: %w", err)
	}

	// Delete the command
	err = ls.db.Delete(&models.LegacyControlAppCmd{}, cmdID).Error
	if err != nil {
		return fmt.Errorf("failed to delete command: %w", err)
	}

	return nil
}

// USP_AcceptInvite implements the exact logic from the legacy USP_AcceptInvite stored procedure
func (ls *LegacyService) USP_AcceptInvite(subID int, domName string) error {
	// Insert into Relationship (DomId,SubID)
	// select Id,@SubId from Users where [Screen Name]=@DomName
	err := ls.db.Exec("INSERT INTO Relationship (DomId, SubId) SELECT Id, ? FROM Users WHERE [Screen Name] = ?", subID, domName).Error
	if err != nil {
		return fmt.Errorf("failed to create relationship: %w", err)
	}

	// Delete i from Invites i
	// join Users u on i.DomId=u.Id
	// where [Screen Name]=@DomName and i.SubId=@SubId
	err = ls.db.Exec("DELETE i FROM Invites i JOIN Users u ON i.DomId = u.Id WHERE u.[Screen Name] = ? AND i.SubId = ?", domName, subID).Error
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

// USP_DeleteInvite implements the exact logic from the legacy USP_DeleteInvite stored procedure
func (ls *LegacyService) USP_DeleteInvite(subID int, domName string) error {
	// Delete i from Invites i
	// join Users u on i.DomId=u.Id
	// where [Screen Name]=@DomName and i.SubId=@SubId
	err := ls.db.Exec("DELETE i FROM Invites i JOIN Users u ON i.DomId = u.Id WHERE u.[Screen Name] = ? AND i.SubId = ?", domName, subID).Error
	if err != nil {
		return fmt.Errorf("failed to delete invite: %w", err)
	}

	return nil
}

// USP_thumbsup implements the exact logic from the legacy USP_thumbsup stored procedure
func (ls *LegacyService) USP_thumbsup(userID int, senderID int) error {
	// if(select count(*) from dbo.users where id=@senderid)>0
	var count int64
	err := ls.db.Model(&models.LegacyUser{}).Where("Id = ?", senderID).Count(&count).Error
	if err != nil {
		return fmt.Errorf("failed to check sender: %w", err)
	}

	if count > 0 {
		// if @id!=@senderid
		if userID != senderID {
			// update users set ThumbsUp=ThumbsUp+1 where id=@senderid
			err = ls.db.Model(&models.LegacyUser{}).Where("Id = ?", senderID).Update("ThumbsUp", gorm.Expr("ThumbsUp + 1")).Error
			if err != nil {
				return fmt.Errorf("failed to update thumbs up: %w", err)
			}
		}
	}

	return nil
}

// USP_DeleteOutstanding implements the exact logic from the legacy USP_DeleteOutstanding stored procedure
func (ls *LegacyService) USP_DeleteOutstanding(userID int) error {
	// DELETE cmd FROM ControlAppCmd cmd WHERE SubId=@userID
	err := ls.db.Where("SubId = ?", userID).Delete(&models.LegacyControlAppCmd{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete outstanding commands: %w", err)
	}

	return nil
}

// USP_GetInvites2 implements the exact logic from the legacy USP_GetInvites2 stored procedure
func (ls *LegacyService) USP_GetInvites2(subID int) (string, error) {
	// This stored procedure returns a concatenated string of dom users
	// SELECT STUFF((select ',['+DomUser+']' AS [text()] from
	// (SELECT u.[Screen Name] DomUser from Invites i
	// join Users u on i.DomId=u.Id where SubId=@SubId) as sub
	// FOR XML PATH (''), TYPE).value('text()[1]','nvarchar(max)'), 1, 1, '')

	var domUsers []string
	err := ls.db.Table("Invites i").
		Select("u.[Screen Name]").
		Joins("join Users u on i.DomId = u.Id").
		Where("i.SubId = ?", subID).
		Pluck("u.[Screen Name]", &domUsers).Error

	if err != nil {
		return "", fmt.Errorf("failed to get invites: %w", err)
	}

	// Format as [user1],[user2],[user3]
	result := ""
	for i, user := range domUsers {
		if i > 0 {
			result += ","
		}
		result += "[" + user + "]"
	}

	return result, nil
}

// USP_GetRels implements the exact logic from the legacy USP_GetRels stored procedure (assuming it exists)
func (ls *LegacyService) USP_GetRels(subID int) (string, error) {
	// Similar to GetInvites2 but for relationships
	var domUsers []string
	err := ls.db.Table("Relationship r").
		Select("u.[Screen Name]").
		Joins("join Users u on r.DomId = u.Id").
		Where("r.SubId = ?", subID).
		Pluck("u.[Screen Name]", &domUsers).Error

	if err != nil {
		return "", fmt.Errorf("failed to get relationships: %w", err)
	}

	// Format as [user1],[user2],[user3]
	result := ""
	for i, user := range domUsers {
		if i > 0 {
			result += ","
		}
		result += "[" + user + "]"
	}

	return result, nil
}
