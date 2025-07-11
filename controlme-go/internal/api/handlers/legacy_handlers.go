package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
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

	// Validate required parameters
	if usernm == "" || pwd == "" {
		result = "Missing required parameters"
	} else if vrs != "012" && vrs != "" {
		result = "Unsupported version"
	} else {
		// Authenticate user
		user, err := h.userService.AuthenticateLegacyUser(usernm, pwd)
		if err != nil {
			result = "Authentication failed"
		} else {
			// Handle different command types
			switch cmd {
			case "Outstanding":
				// Get outstanding command count and sender info (like USP_GetOutstanding)
				count, err := h.cmdService.GetPendingCommandCount(user.ID)
				if err != nil {
					result = "Failed to get outstanding commands"
				} else {
					// Format: [count],[whonext],[verified],[thumbs]
					// For now, simplified version - could enhance with relationship/group logic
					whonext := "User" // Simplified - in legacy this checks relationships/groups
					verified := "0"   // User.Verified field
					thumbs := "0"     // User.ThumbsUp field
					result = fmt.Sprintf("[%d],[%s],[%s],[%s]", count, whonext, verified, thumbs)
				}

			case "Content":
				// Get next command content (like USP_GetAppContent)
				assignment, err := h.cmdService.GetNextCommand(user.ID)
				if err != nil {
					result = "Failed to get content"
				} else if assignment == nil {
					result = "" // No commands
				} else {
					// Format: [SenderId],[Content] or [SenderId],[Content],[SenderName]
					senderName := assignment.Sender.ScreenName
					if senderName == "" {
						result = fmt.Sprintf("[%s],[%s]", assignment.SenderID.String(), assignment.Command.Data)
					} else {
						result = fmt.Sprintf("[%s],[%s],[%s]", assignment.SenderID.String(), assignment.Command.Data, senderName)
					}
					// Mark command as completed
					_ = h.cmdService.CompleteCommand(user.ID)
				}

			default:
				// Legacy behavior for sending commands (original AppCommand logic)
				from := c.Query("From")
				to := c.Query("To")
				data := c.Query("Data")
				password := c.Query("Password")

				if from == "" || to == "" || data == "" || password == "" {
					result = "Missing command parameters"
				} else {
					// Get recipient
					recipient, err := h.userService.GetUserByUsername(to)
					if err != nil {
						result = "Recipient not found"
					} else {
						// Decrypt command data
						decryptedData, err := h.authService.LegacyCrypto.Decrypt(data)
						if err != nil {
							result = "Failed to decrypt command data"
						} else {
							// Parse command data (format: "command|additional_data")
							parts := strings.SplitN(decryptedData, "|", 2)
							commandType := parts[0]
							commandContent := ""
							if len(parts) > 1 {
								commandContent = parts[1]
							}

							// Create command
							command, err := h.cmdService.CreateCommand(commandType, commandContent, data)
							if err != nil {
								result = "Failed to create command"
							} else {
								// Assign command to recipient
								_, err = h.cmdService.AssignCommandToUser(user.ID, recipient.ID, command.ID, nil)
								if err != nil {
									result = "Failed to assign command"
								} else {
									result = "OK"
								}
							}
						}
					}
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
            <span id="result">%s</span>
        </div>
    </form>
</body>
</html>
`, result)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// GetContent handles legacy /GetContent.aspx endpoint
// Original: Gets pending commands for a client
func (h *LegacyHandlers) GetContent(c *gin.Context) {
	// Get query parameters (exact legacy format)
	usernm := c.Query("usernm")
	pwd := c.Query("pwd")
	vrs := c.Query("vrs")

	// Initialize response variables
	senderID := ""
	result := ""
	verified := ""

	// Validate required parameters and version
	if usernm == "" || pwd == "" {
		result = "Missing required parameters"
	} else if vrs != "012" && vrs != "" {
		result = "Unsupported version"
	} else {
		// Authenticate user (decrypt password first)
		user, err := h.userService.AuthenticateLegacyUser(usernm, pwd)
		if err != nil {
			result = "Authentication failed"
		} else {
			// Get next pending command for user
			assignment, err := h.cmdService.GetNextCommand(user.ID)
			if err != nil {
				result = "Failed to get commands"
			} else if assignment == nil {
				// No commands found, leave result empty (legacy behavior)
			} else {
				// Format response exactly like legacy stored procedure
				// Legacy stores user IDs as integers, but we use UUIDs - need to convert or use a mapping
				senderID = assignment.SenderID.String()
				result = assignment.Command.Data

				// Mark command as completed (legacy behavior: exec USP_CmdComplete)
				_ = h.cmdService.CompleteCommand(user.ID)
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
            <span id="SenderId">%s</span>
            <span id="Result">%s</span>
            <span id="Varified">%s</span>
        </div>
    </form>
</body>
</html>
`, senderID, result, verified)

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

// GetCount handles legacy /GetCount.aspx endpoint
// Original: Gets count of pending commands
func (h *LegacyHandlers) GetCount(c *gin.Context) {
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

	// Get command count
	count, err := h.cmdService.GetPendingCommandCount(authenticatedUser.ID)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to get command count")
		return
	}

	// Return count as string (legacy format)
	c.String(http.StatusOK, strconv.FormatInt(count, 10))
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
