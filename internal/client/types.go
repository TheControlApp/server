package client

import (
	"encoding/json"
	"time"

	"github.com/thecontrolapp/server/internal/models"
)

// Command represents a ControlApp command (reusing server's Instruction model)
type Command struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Content    map[string]interface{} `json:"content"`
	ReceivedAt time.Time              `json:"received_at"`
	SenderID   string                 `json:"sender_id,omitempty"`
}

// CommandFromInstruction converts a server Instruction to client Command
func CommandFromInstruction(instr models.Instruction, id, senderID string) Command {
	content := make(map[string]interface{})

	// Convert interface{} to map[string]interface{} if possible
	if instrContent, ok := instr.Content.(map[string]interface{}); ok {
		content = instrContent
	} else {
		// Try to marshal and unmarshal to convert
		if data, err := json.Marshal(instr.Content); err == nil {
			json.Unmarshal(data, &content)
		}
	}

	return Command{
		ID:         id,
		Type:       instr.Type,
		Content:    content,
		ReceivedAt: time.Now(),
		SenderID:   senderID,
	}
}

// CommandResult represents the result of command execution
type CommandResult struct {
	CommandID string        `json:"command_id"`
	Status    string        `json:"status"` // "completed", "failed", "cancelled", "timeout", "partial"
	Result    interface{}   `json:"result,omitempty"`
	Error     string        `json:"error,omitempty"`
	Duration  time.Duration `json:"duration"`
	Metadata  interface{}   `json:"metadata,omitempty"`
}

// CommandHandler defines the function signature for command handlers
type CommandHandler func(cmd Command) CommandResult

// Config represents client configuration
type Config struct {
	// Server connection
	ServerURL     string `yaml:"server_url"`
	Timeout       int    `yaml:"timeout"`
	ReconnectWait int    `yaml:"reconnect_wait"`

	// Application info
	AppName string `yaml:"app_name"`
	Version string `yaml:"version"`

	// Logging
	LogLevel string `yaml:"log_level"`
	LogFile  string `yaml:"log_file"`
	Logger   Logger `yaml:"-"`

	// Security & Consent
	AllowedCommands []string `yaml:"allowed_commands"`
	BlockedCommands []string `yaml:"blocked_commands"`
	RequireConfirm  []string `yaml:"require_confirmation"`
	EmergencyHotkey string   `yaml:"emergency_hotkey"`

	// File handling
	DownloadFolder   string   `yaml:"download_folder"`
	MaxFileSize      int64    `yaml:"max_file_size_mb"`
	AllowedFileTypes []string `yaml:"allowed_file_types"`

	// Audio settings
	AudioVolume int    `yaml:"audio_volume"`
	AllowTTS    bool   `yaml:"allow_tts"`
	TTSVoice    string `yaml:"tts_voice"`

	// Display settings
	AllowWallpaper  bool   `yaml:"allow_wallpaper"`
	WallpaperStyle  string `yaml:"wallpaper_style"`
	RestorePrevious bool   `yaml:"restore_previous_wallpaper"`
}

// DefaultConfig returns a safe default configuration
func DefaultConfig() *Config {
	return &Config{
		ServerURL:     "ws://localhost:3080",
		Timeout:       30,
		ReconnectWait: 5,
		AppName:       "ControlApp Client",
		Version:       "1.0.0",
		LogLevel:      "info",

		// Safe defaults - user must explicitly enable
		AllowedCommands: []string{
			"std_ping",
			"kink_message",
		},
		BlockedCommands: []string{
			"kink_run_file", // Blocked by default for security
		},
		RequireConfirm: []string{
			"kink_run_file",
			"kink_lock_screen",
			"kink_change_wallpaper",
		},

		DownloadFolder:   "downloads",
		MaxFileSize:      50, // 50MB
		AllowedFileTypes: []string{".jpg", ".jpeg", ".png", ".gif", ".mp3", ".wav"},

		AudioVolume: 80,
		AllowTTS:    true,
		TTSVoice:    "default",

		AllowWallpaper:  false, // User must explicitly enable
		WallpaperStyle:  "fill",
		RestorePrevious: true,
	}
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type      string      `json:"type"`
	Payload   interface{} `json:"payload"`
	Timestamp string      `json:"timestamp"`
}

// LoginRequest represents a login request (compatible with server)
type LoginRequest struct {
	LoginName string `json:"login_name"`
	Password  string `json:"password"`
}

// RegisterRequest represents a registration request (compatible with server)
type RegisterRequest struct {
	ScreenName string `json:"screen_name"`
	LoginName  string `json:"login_name"`
	Password   string `json:"password"`
}

// AuthResponse represents an authentication response (compatible with server)
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"` // Use server's User model
}

// ConsentLevel represents different levels of user consent
type ConsentLevel int

const (
	ConsentDenied   ConsentLevel = iota
	ConsentBasic                 // Basic commands only (ping, message)
	ConsentStandard              // Most commands allowed
	ConsentFull                  // All commands allowed (with confirmation)
)

func (c ConsentLevel) String() string {
	switch c {
	case ConsentDenied:
		return "Denied"
	case ConsentBasic:
		return "Basic"
	case ConsentStandard:
		return "Standard"
	case ConsentFull:
		return "Full"
	default:
		return "Unknown"
	}
}

// CommandAudit represents an audit log entry
type CommandAudit struct {
	CommandID     string    `json:"command_id"`
	CommandType   string    `json:"command_type"`
	SenderID      string    `json:"sender_id"`
	ReceiverID    string    `json:"receiver_id"`
	ExecutedAt    time.Time `json:"executed_at"`
	Status        string    `json:"status"`
	UserConsent   bool      `json:"user_consent"`
	SafetyChecks  bool      `json:"safety_checks_passed"`
	ErrorMessage  string    `json:"error_message,omitempty"`
	ExecutionTime int64     `json:"execution_time_ms"`
}

// ParsePayload parses a payload into the target structure
func ParsePayload(payload interface{}, target interface{}) error {
	// Convert to JSON and back to parse into target struct
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}

// IsCommandAllowed checks if a command type is allowed by configuration
func (c *Config) IsCommandAllowed(commandType string) bool {
	// Check blocked list first
	for _, blocked := range c.BlockedCommands {
		if blocked == commandType {
			return false
		}
	}

	// Check allowed list
	for _, allowed := range c.AllowedCommands {
		if allowed == commandType {
			return true
		}
	}

	return false
}

// RequiresConfirmation checks if a command requires user confirmation
func (c *Config) RequiresConfirmation(commandType string) bool {
	for _, cmd := range c.RequireConfirm {
		if cmd == commandType {
			return true
		}
	}
	return false
}
