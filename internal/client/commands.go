package client

import (
	"fmt"
	"sync"
	"time"
)

// CommandProcessor handles command execution and routing
type CommandProcessor struct {
	handlers map[string]CommandHandler
	mu       sync.RWMutex
	logger   Logger
	config   *Config

	// Command history
	history   []CommandAudit
	historyMu sync.RWMutex
}

// NewCommandProcessor creates a new command processor
func NewCommandProcessor(logger Logger) *CommandProcessor {
	cp := &CommandProcessor{
		handlers: make(map[string]CommandHandler),
		logger:   logger,
		history:  make([]CommandAudit, 0),
	}

	// Register core command handlers
	cp.registerCoreHandlers()

	return cp
}

// RegisterHandler registers a handler for a specific command type
func (cp *CommandProcessor) RegisterHandler(cmdType string, handler CommandHandler) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.handlers[cmdType] = handler
	cp.logger.Info("Registered command handler", "type", cmdType)
}

// ProcessCommand processes a command and returns the result
func (cp *CommandProcessor) ProcessCommand(cmd Command) CommandResult {
	start := time.Now()

	// Check if handler exists
	cp.mu.RLock()
	handler, exists := cp.handlers[cmd.Type]
	cp.mu.RUnlock()

	if !exists {
		result := CommandResult{
			CommandID: cmd.ID,
			Status:    "failed",
			Error:     fmt.Sprintf("Unknown command type: %s", cmd.Type),
			Duration:  time.Since(start),
		}
		cp.auditCommand(cmd, result, false, false)
		return result
	}

	// Check consent if config is available
	if cp.config != nil {
		if !cp.config.IsCommandAllowed(cmd.Type) {
			result := CommandResult{
				CommandID: cmd.ID,
				Status:    "failed",
				Error:     fmt.Sprintf("Command type %s not allowed by user consent", cmd.Type),
				Duration:  time.Since(start),
			}
			cp.auditCommand(cmd, result, false, false)
			return result
		}
	}

	// Execute command
	cp.logger.Info("Executing command", "type", cmd.Type, "id", cmd.ID)

	result := handler(cmd)
	result.CommandID = cmd.ID
	result.Duration = time.Since(start)

	// Audit the command execution
	cp.auditCommand(cmd, result, true, result.Status == "completed")

	cp.logger.Info("Command executed",
		"type", cmd.Type,
		"id", cmd.ID,
		"status", result.Status,
		"duration", result.Duration)

	return result
}

// SetConfig sets the configuration for consent checking
func (cp *CommandProcessor) SetConfig(config *Config) {
	cp.config = config
}

// GetHistory returns the command execution history
func (cp *CommandProcessor) GetHistory() []CommandAudit {
	cp.historyMu.RLock()
	defer cp.historyMu.RUnlock()

	// Return a copy
	history := make([]CommandAudit, len(cp.history))
	copy(history, cp.history)
	return history
}

// registerCoreHandlers registers the basic command handlers
func (cp *CommandProcessor) registerCoreHandlers() {
	// std_ping - connectivity test
	cp.RegisterHandler("std_ping", func(cmd Command) CommandResult {
		return CommandResult{
			Status: "completed",
			Result: map[string]interface{}{
				"pong_timestamp": time.Now().Format(time.RFC3339Nano),
				"latency_ms":     0, // Will be calculated by server
			},
		}
	})

	// Placeholder handlers for kink commands - will be overridden by platform-specific implementations
	kinkCommands := []string{
		"kink_message",
		"kink_open_link",
		"kink_download_file",
		"kink_change_wallpaper",
		"kink_run_file",
		"kink_play_audio",
		"kink_tts",
		"kink_popup_image",
		"kink_timer_task",
		"kink_lock_screen",
		"kink_scheduled_task",
		"kink_mood_check",
		"kink_permission_request",
	}

	for _, cmdType := range kinkCommands {
		cmdType := cmdType // Capture for closure
		cp.RegisterHandler(cmdType, func(cmd Command) CommandResult {
			return CommandResult{
				Status: "failed",
				Error:  fmt.Sprintf("%s requires platform-specific implementation", cmdType),
			}
		})
	}
}

// auditCommand records a command execution in the audit log
func (cp *CommandProcessor) auditCommand(cmd Command, result CommandResult, userConsent, safetyChecks bool) {
	audit := CommandAudit{
		CommandID:     cmd.ID,
		CommandType:   cmd.Type,
		SenderID:      cmd.SenderID,
		ExecutedAt:    time.Now(),
		Status:        result.Status,
		UserConsent:   userConsent,
		SafetyChecks:  safetyChecks,
		ExecutionTime: result.Duration.Milliseconds(),
	}

	if result.Error != "" {
		audit.ErrorMessage = result.Error
	}

	cp.historyMu.Lock()
	cp.history = append(cp.history, audit)

	// Keep only last 100 entries to avoid memory issues
	if len(cp.history) > 100 {
		cp.history = cp.history[len(cp.history)-100:]
	}
	cp.historyMu.Unlock()
}
