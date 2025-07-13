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
	db            *gorm.DB
	userService   *services.UserService
	cmdService    *services.CommandService
	legacyService *services.LegacyService
	authService   *auth.AuthService
	config        *config.Config
}

// NewLegacyHandlers creates a new legacy handlers instance
func NewLegacyHandlers(db *gorm.DB, userService *services.UserService, cmdService *services.CommandService, legacyService *services.LegacyService, authService *auth.AuthService, cfg *config.Config) *LegacyHandlers {
	return &LegacyHandlers{
		db:            db,
		userService:   userService,
		cmdService:    cmdService,
		legacyService: legacyService,
		authService:   authService,
		config:        cfg,
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
		// Authenticate user using legacy service (supports DES, AES, plain, and bcrypt)
		loginResult, err := h.legacyService.USP_Login(usernm, pwd)
		if err != nil {
			result = err.Error()
		} else {
			// Handle different command types
			switch cmd {
			case "Outstanding":
				// Get outstanding command count and sender info (USP_GetOutstanding)
				outstanding, err := h.legacyService.USP_GetOutstanding(loginResult.ID)
				if err != nil {
					result = err.Error()
				} else {
					// Format exactly like USP_GetOutstanding: [count],[whonext],[verified],[thumbs]
					result = fmt.Sprintf("[%d],[%s],[%s],[%s]", outstanding.Count, outstanding.WhoNext, outstanding.Verified, outstanding.Thumbs)
				}

			case "Content":
				// Get next command content (USP_GetAppContent)
				content, err := h.legacyService.USP_GetAppContent(loginResult.ID)
				if err != nil {
					result = err.Error()
				} else if content != nil {
					// Format: [SenderId],[Content],[SenderName]
					result = fmt.Sprintf("[%s],[%s],[%s]", content.SenderID, content.Content, content.SenderName)
					// Mark command as completed (USP_CmdComplete)
					_ = h.legacyService.USP_CmdComplete(loginResult.ID)
				}

			case "Delete":
				// Delete outstanding commands (USP_DeleteOutstanding)
				err := h.legacyService.USP_DeleteOutstanding(loginResult.ID)
				if err != nil {
					result = err.Error()
				}

			default:
				result = "Unknown command: " + cmd
			}
		}
	} else {
		result = "Wrong version."
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
		// Decrypt password using legacy crypto
		decryptedPwd, err := h.authService.LegacyCrypto.Decrypt(pwd)
		if err == nil {
			// Authenticate user using legacy database
			loginResult, err := h.legacyService.USP_Login(usernm, decryptedPwd)
			if err == nil && loginResult != nil {
				// Get app content using legacy stored procedure
				content, err := h.legacyService.USP_GetAppContent(loginResult.ID)
				if err == nil && content != nil {
					// Format response exactly like USP_GetAppContent stored procedure
					senderID = content.SenderID
					result = content.Content
					verified = content.SenderName

					// Mark command as completed (USP_CmdComplete)
					_ = h.legacyService.USP_CmdComplete(loginResult.ID)
				}
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
		// Decrypt password using legacy crypto
		decryptedPwd, err := h.authService.LegacyCrypto.Decrypt(pwd)
		if err == nil {
			// Authenticate user using legacy database
			loginResult, err := h.legacyService.USP_Login(usernm, decryptedPwd)
			if err == nil && loginResult != nil {
				// Get outstanding command info using legacy stored procedure
				outstanding, err := h.legacyService.USP_GetOutstanding(loginResult.ID)
				if err == nil {
					result = fmt.Sprintf("%d", outstanding.Count)
					next = outstanding.WhoNext
					vari = outstanding.Verified
				}
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
