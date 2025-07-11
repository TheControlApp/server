package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"github.com/thecontrolapp/controlme-go/internal/services"
	"gorm.io/gorm"
)

// LegacyHandlers contains handlers for legacy ASP.NET endpoints
type LegacyHandlers struct {
	db          *gorm.DB
	userService *services.UserService
	cmdService  *services.CommandService
	authService *auth.AuthService
	config      *config.Config
}

// NewLegacyHandlers creates a new legacy handlers instance
func NewLegacyHandlers(db *gorm.DB, userService *services.UserService, cmdService *services.CommandService, authService *auth.AuthService, cfg *config.Config) *LegacyHandlers {
	return &LegacyHandlers{
		db:          db,
		userService: userService,
		cmdService:  cmdService,
		authService: authService,
		config:      cfg,
	}
}

// AppCommand handles legacy /AppCommand.aspx endpoint
// Original: Multiple functions based on cmd parameter
func (h *LegacyHandlers) AppCommand(c *gin.Context) {
	// Get query parameters (legacy format)
	usernm := c.Query("usernm")
	pwd := c.Query("pwd")
	vrs := c.Query("vrs")
	cmd := c.Query("cmd")

	// Initialize result
	result := ""

	// Validate version first (legacy behavior)
	if vrs == "012" {
		// Authenticate user
		user, err := h.userService.AuthenticateLegacyUser(usernm, pwd)
		if err != nil {
			result = err.Error()
		} else {
			// Handle different command types
			switch cmd {
			case "Outstanding":
				// Get outstanding command count and sender info (like USP_GetOutstanding)
				// 1. Delete commands from blocked users
				h.db.Where("sender_id IN (SELECT blockee_id FROM blocks WHERE blocker_id = ?) AND sub_id = ?", user.ID, user.ID).Delete(&models.ControlAppCmd{})
				
				// 2. Delete anonymous commands if user doesn't allow them
				if !user.AnonCmd {
					h.db.Where("sender_id = ? AND sub_id = ?", "00000000-0000-0000-0000-000000000001", user.ID).Delete(&models.ControlAppCmd{})
				}

				// 3. Update login date
				user.LoginDate = time.Now()
				h.db.Save(&user)

				// 4. Get command count and next sender info
				var count int64
				err = h.db.Model(&models.ControlAppCmd{}).Where("sub_id = ?", user.ID).Count(&count).Error
				if err != nil {
					result = "Failed to get command count"
				} else {
					whonext := "User"
					verified := "0"
					if user.Verified {
						verified = "1"
					}
					thumbs := strconv.Itoa(user.ThumbsUp)

					// Get next sender info
					var assignment models.ControlAppCmd
					err = h.db.Preload("Sender").
						Where("sub_id = ?", user.ID).
						Order("created_at ASC").
						First(&assignment).Error

					if err == nil {
						// Check if it's a group command
						if assignment.GroupRefID != nil {
							var group models.Group
							err := h.db.First(&group, "id = ?", assignment.GroupRefID).Error
							if err == nil {
								whonext = "Group: " + group.Name
							} else {
								whonext = "Group"
							}
						} else if assignment.SenderID.String() == "00000000-0000-0000-0000-000000000001" {
							whonext = "Anon"
						} else {
							// Check if sender is in a relationship with the user
							var relationship models.Relationship
							err = h.db.Where("dom_id = ? AND sub_id = ?", assignment.SenderID, user.ID).First(&relationship).Error
							if err == nil {
								whonext = assignment.Sender.ScreenName
							} else {
								whonext = "User"
							}
						}
					}

					// Format exactly like USP_GetOutstanding: [count],[whonext],[verified],[thumbs]
					result = fmt.Sprintf("[%d],[%s],[%s],[%s]", count, whonext, verified, thumbs)
				}

			case "Content":
				// Get next command content (like USP_GetAppContent)
				// 1. Delete commands from blocked users
				h.db.Where("sender_id IN (SELECT blockee_id FROM blocks WHERE blocker_id = ?) AND sub_id = ?", user.ID, user.ID).Delete(&models.ControlAppCmd{})
				
				// 2. Delete anonymous commands if user doesn't allow them
				if !user.AnonCmd {
					h.db.Where("sender_id = ? AND sub_id = ?", "00000000-0000-0000-0000-000000000001", user.ID).Delete(&models.ControlAppCmd{})
				}

				// 3. Get next command
				var assignment models.ControlAppCmd
				err = h.db.Preload("Sender").Preload("Command").
					Where("sub_id = ?", user.ID).
					Order("created_at ASC").
					First(&assignment).Error

				if err == nil {
					// Format: [SenderId],[Content] or [SenderId],[Content],[SenderName]
					senderName := assignment.Sender.ScreenName
					if senderName == "" {
						result = fmt.Sprintf("[%s],[%s]", assignment.SenderID.String(), assignment.Command.Content)
					} else {
						result = fmt.Sprintf("[%s],[%s],[%s]", assignment.SenderID.String(), assignment.Command.Content, senderName)
					}
					// Mark command as completed (exec USP_CmdComplete)
					_ = h.cmdService.CompleteCommand(user.ID)
				}

			default:
				// Handle other command types (Accept, Reject, Thumbs, etc.)
				if len(cmd) > 6 {
					switch cmd[:6] {
					case "Accept":
						// USP_AcceptInvite
						inviteFrom := cmd[6:]
						// TODO: Implement invite acceptance logic
						result = "Invite accepted: " + inviteFrom
					case "Reject":
						// USP_DeleteInvite
						inviteFrom := cmd[6:]
						// TODO: Implement invite rejection logic
						result = "Invite rejected: " + inviteFrom
					case "Thumbs":
						// USP_thumbsup
						thumbsValue := cmd[6:]
						// TODO: Implement thumbs up logic
						result = "Thumbs up: " + thumbsValue
					default:
						result = "Unknown command: " + cmd
					}
				} else {
					result = "Unknown command: " + cmd
				}
			}
		}
	}

	// Return HTML response matching original ASP.NET page structure
	html := fmt.Sprintf(`<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml">
<head runat="server">
    <title></title>
</head>
<body>
    <form id="form1" runat="server">
        <div>
            <asp:Label ID="result" runat="server">%s</asp:Label>
        </div>
    </form>
</body>
</html>
`, result)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// GetContent handles legacy /GetContent.aspx endpoint
// Original: Gets pending commands for a client (USP_GetAppContent)
func (h *LegacyHandlers) GetContent(c *gin.Context) {
	// Get query parameters (exact legacy format)
	usernm := c.Query("usernm")
	pwd := c.Query("pwd")
	vrs := c.Query("vrs")

	// Initialize response variables
	senderID := ""
	result := ""
	verified := ""

	// Validate version first (legacy behavior)
	if vrs == "012" {
		// Authenticate user (decrypt password first)
		user, err := h.userService.AuthenticateLegacyUser(usernm, pwd)
		if err == nil && user != nil {
			// Implement exact USP_GetAppContent logic
			// 1. Delete commands from blocked users
			h.db.Where("sender_id IN (SELECT blockee_id FROM blocks WHERE blocker_id = ?) AND sub_id = ?", user.ID, user.ID).Delete(&models.ControlAppCmd{})
			
			// 2. Delete anonymous commands if user doesn't allow them
			if !user.AnonCmd {
				h.db.Where("sender_id = ? AND sub_id = ?", "00000000-0000-0000-0000-000000000001", user.ID).Delete(&models.ControlAppCmd{})
			}

			// 3. Get next command with sender info (TOP 1 ORDER BY SendDate)
			var assignment models.ControlAppCmd
			err = h.db.Preload("Command").Preload("Sender").
				Where("sub_id = ?", user.ID).
				Order("created_at ASC").
				First(&assignment).Error

			if err == nil {
				// Format response exactly like USP_GetAppContent stored procedure
				senderID = assignment.SenderID.String()
				result = assignment.Command.Content // Use Content field, not Data
				verified = assignment.Sender.ScreenName // SenderName from stored procedure
				
				// Mark command as completed (exec USP_CmdComplete)
				_ = h.cmdService.CompleteCommand(user.ID)
			}
		}
	}

	// Return HTML response matching exact ASP.NET page structure
	html := fmt.Sprintf(`<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml">
<head runat="server">
    <title></title>
</head>
<body>
    <form id="form1" runat="server">
        <div>
            <asp:Label ID="SenderId" runat="server">%s</asp:Label>
            <asp:Label ID="Result" runat="server">%s</asp:Label>
            <asp:Label ID="Varified" runat="server">%s</asp:Label>
        </div>
    </form>
</body>
</html>
`, senderID, result, verified)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// GetCount handles legacy /GetCount.aspx endpoint
// Original: Gets count of pending commands (USP_GetOutstanding)
func (h *LegacyHandlers) GetCount(c *gin.Context) {
	// Get query parameters (exact legacy format)
	usernm := c.Query("usernm")
	pwd := c.Query("pwd")
	vrs := c.Query("vrs")

	// Initialize response variables
	result := ""
	next := ""
	vari := ""

	// Validate version first (legacy behavior)
	if vrs == "012" {
		// Authenticate user (decrypt password first)
		user, err := h.userService.AuthenticateLegacyUser(usernm, pwd)
		if err == nil && user != nil {
			// Implement exact USP_GetOutstanding logic
			// 1. Delete commands from blocked users
			h.db.Where("sender_id IN (SELECT blockee_id FROM blocks WHERE blocker_id = ?) AND sub_id = ?", user.ID, user.ID).Delete(&models.ControlAppCmd{})
			
			// 2. Delete anonymous commands if user doesn't allow them
			if !user.AnonCmd {
				h.db.Where("sender_id = ? AND sub_id = ?", "00000000-0000-0000-0000-000000000001", user.ID).Delete(&models.ControlAppCmd{})
			}

			// 3. Update login date
			user.LoginDate = time.Now()
			h.db.Save(&user)

			// 4. Get command count
			var count int64
			err = h.db.Model(&models.ControlAppCmd{}).Where("sub_id = ?", user.ID).Count(&count).Error
			if err == nil {
				result = strconv.FormatInt(count, 10)
			}

			// 5. Get "whonext" info - who is the next sender
			var assignment models.ControlAppCmd
			err = h.db.Preload("Sender").
				Where("sub_id = ?", user.ID).
				Order("created_at ASC").
				First(&assignment).Error

			if err == nil {
				// Check if it's a group command
				if assignment.GroupRefID != nil {
					var group models.Group
					err := h.db.First(&group, "id = ?", assignment.GroupRefID).Error
					if err == nil {
						next = "Group: " + group.Name
					} else {
						next = "Group"
					}
				} else if assignment.SenderID.String() == "00000000-0000-0000-0000-000000000001" {
					next = "Anon"
				} else {
					// Check if sender is in a relationship with the user
					var relationship models.Relationship
					err = h.db.Where("dom_id = ? AND sub_id = ?", assignment.SenderID, user.ID).First(&relationship).Error
					if err == nil {
						next = assignment.Sender.ScreenName
					} else {
						next = "User"
					}
				}
			} else {
				next = "User" // Default if no commands
			}

			// 6. Set verified status
			if user.Verified {
				vari = "1"
			} else {
				vari = "0"
			}
		}
	}

	// Return HTML response matching exact ASP.NET page structure
	html := fmt.Sprintf(`<!DOCTYPE html>

<html xmlns="http://www.w3.org/1999/xhtml">
<head runat="server">
    <title></title>
</head>
<body>
    <form id="form1" runat="server">
        <div>
            <asp:Label ID="result" runat="server">%s</asp:Label>
            <asp:Label ID="next" runat="server">%s</asp:Label>
            <asp:Label ID="vari" runat="server">%s</asp:Label>
        </div>
    </form>
</body>
</html>
`, result, next, vari)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// ProcessComplete handles legacy /ProcessComplete.aspx endpoint
// Original: Marks a command as completed
func (h *LegacyHandlers) ProcessComplete(c *gin.Context) {
	// Get form parameters
	user := c.PostForm("User")
	password := c.PostForm("Password")
	commandID := c.PostForm("CommandID")

	if user == "" || password == "" || commandID == "" {
		c.String(http.StatusBadRequest, "Missing required parameters")
		return
	}

	// Authenticate user
	authenticatedUser, err := h.userService.AuthenticateLegacyUser(user, password)
	if err != nil {
		c.String(http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Parse command ID
	cmdUUID, err := uuid.Parse(commandID)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid command ID")
		return
	}

	// Mark command as completed
	err = h.cmdService.MarkCommandCompleted(cmdUUID, authenticatedUser.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to mark command as completed")
		return
	}

	// Return success
	c.String(http.StatusOK, "OK")
}

// DeleteOut handles legacy /DeleteOut.aspx endpoint
// Original: Deletes outstanding commands
func (h *LegacyHandlers) DeleteOut(c *gin.Context) {
	// Get form parameters
	user := c.PostForm("User")
	password := c.PostForm("Password")

	if user == "" || password == "" {
		c.String(http.StatusBadRequest, "Missing required parameters")
		return
	}

	// Authenticate user
	authenticatedUser, err := h.userService.AuthenticateLegacyUser(user, password)
	if err != nil {
		c.String(http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Delete outstanding commands for user
	err = h.cmdService.DeletePendingCommandsForUser(authenticatedUser.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to delete commands")
		return
	}

	// Return success
	c.String(http.StatusOK, "OK")
}

// GetOptions handles legacy /GetOptions.aspx endpoint
// Original: Gets user options/settings
func (h *LegacyHandlers) GetOptions(c *gin.Context) {
	// Get query parameters (exact legacy format)
	user := c.Query("usernm")
	password := c.Query("pwd")
	version := c.Query("vrs")

	if user == "" || password == "" {
		c.String(http.StatusBadRequest, "Missing required parameters")
		return
	}

	// Check version for legacy compatibility
	if version != "012" && version != "" {
		c.String(http.StatusBadRequest, "Unsupported version")
		return
	}

	// Authenticate user
	authenticatedUser, err := h.userService.AuthenticateLegacyUser(user, password)
	if err != nil {
		c.String(http.StatusUnauthorized, "Authentication failed")
		return
	}

	// Get user settings (for now, return basic info)
	// In the original, this returned various user preferences
	options := map[string]string{
		"ScreenName": authenticatedUser.ScreenName,
		"Role":       authenticatedUser.Role,
		"Verified":   strconv.FormatBool(authenticatedUser.Verified),
	}

	// Format as key=value pairs (legacy format)
	var response strings.Builder
	for key, value := range options {
		if response.Len() > 0 {
			response.WriteString("\n")
		}
		response.WriteString(key)
		response.WriteString("=")
		response.WriteString(value)
	}

	c.String(http.StatusOK, response.String())
}
