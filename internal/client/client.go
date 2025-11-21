package client

import (
	"context"
	"time"

	"github.com/thecontrolapp/server/internal/models"
)

// Client represents the core ControlApp client
type Client struct {
	config   *Config
	auth     *AuthManager
	ws       *WebSocketManager
	commands *CommandProcessor
	logger   Logger

	// Client state
	connected     bool
	authenticated bool
	user          *models.User

	// Event channels
	eventChan   chan Event
	commandChan chan Command
	errorChan   chan error

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// ClientInterface defines the core client functionality
type ClientInterface interface {
	// Connection management
	Connect(ctx context.Context, serverURL string) error
	Disconnect() error
	IsConnected() bool

	// Authentication
	Login(username, password string) error
	Register(screenName, loginName, password string) error
	Logout() error
	IsAuthenticated() bool

	// Command handling
	SendCommand(cmd Command) error
	RegisterCommandHandler(cmdType string, handler CommandHandler)

	// Event handling
	Events() <-chan Event
	Commands() <-chan Command
	Errors() <-chan error

	// Configuration
	SetConfig(config *Config) error
	GetConfig() *Config

	// State
	GetUser() *models.User
	GetConnectionState() ConnectionState
}

// ConnectionState represents the client's connection state
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateAuthenticating
	StateAuthenticated
	StateError
)

func (s ConnectionState) String() string {
	switch s {
	case StateDisconnected:
		return "Disconnected"
	case StateConnecting:
		return "Connecting"
	case StateConnected:
		return "Connected"
	case StateAuthenticating:
		return "Authenticating"
	case StateAuthenticated:
		return "Authenticated"
	case StateError:
		return "Error"
	default:
		return "Unknown"
	}
}

// Event represents a client event
type Event struct {
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

// EventType constants
const (
	EventConnected        = "connected"
	EventDisconnected     = "disconnected"
	EventAuthenticated    = "authenticated"
	EventAuthFailed       = "auth_failed"
	EventCommandReceived  = "command_received"
	EventCommandCompleted = "command_completed"
	EventCommandFailed    = "command_failed"
	EventError            = "error"
)

// Logger interface for flexible logging
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

// NewClient creates a new ControlApp client
func NewClient(config *Config) *Client {
	ctx, cancel := context.WithCancel(context.Background())

	logger := config.Logger
	if logger == nil {
		logger = NewDefaultLogger()
	}

	client := &Client{
		config:      config,
		logger:      logger,
		eventChan:   make(chan Event, 100),
		commandChan: make(chan Command, 100),
		errorChan:   make(chan error, 100),
		ctx:         ctx,
		cancel:      cancel,
	}

	client.auth = NewAuthManager(config, logger)
	client.ws = NewWebSocketManager(config, logger)
	client.commands = NewCommandProcessor(logger)

	// Register core command handlers
	client.registerCoreHandlers()

	return client
}

// Connect establishes connection to the ControlApp server
func (c *Client) Connect(ctx context.Context, serverURL string) error {
	c.logger.Info("Connecting to ControlApp server", "url", serverURL)

	// Update config with server URL
	c.config.ServerURL = serverURL

	// Set state to connecting
	c.connected = false
	c.emitEvent(EventConnected, map[string]interface{}{
		"server_url": serverURL,
	})

	// Connect WebSocket
	err := c.ws.Connect(ctx, serverURL, c.auth)
	if err != nil {
		c.logger.Error("Failed to connect WebSocket", "error", err)
		c.emitError(err)
		return err
	}

	c.connected = true

	// Start message processing
	go c.processMessages()

	c.emitEvent(EventConnected, map[string]interface{}{
		"server_url": serverURL,
		"connected":  true,
	})

	return nil
}

// Disconnect closes the connection to the server
func (c *Client) Disconnect() error {
	c.logger.Info("Disconnecting from ControlApp server")

	c.cancel() // Cancel context

	if c.ws != nil {
		c.ws.Disconnect()
	}

	c.connected = false
	c.authenticated = false

	c.emitEvent(EventDisconnected, nil)

	return nil
}

// IsConnected returns whether the client is connected
func (c *Client) IsConnected() bool {
	return c.connected && c.ws != nil && c.ws.IsConnected()
}

// Login authenticates with the server
func (c *Client) Login(username, password string) error {
	c.logger.Info("Logging in", "username", username)

	err := c.auth.Login(username, password)
	if err != nil {
		c.logger.Error("Login failed", "error", err)
		c.emitEvent(EventAuthFailed, map[string]interface{}{
			"username": username,
			"error":    err.Error(),
		})
		return err
	}

	c.authenticated = true
	c.user = c.auth.GetUser()

	// Reconnect WebSocket with authentication
	if c.IsConnected() {
		c.ws.Authenticate(c.auth.GetToken())
	}

	c.emitEvent(EventAuthenticated, map[string]interface{}{
		"user": c.user,
	})

	return nil
}

// Register creates a new user account
func (c *Client) Register(screenName, loginName, password string) error {
	c.logger.Info("Registering new user", "login_name", loginName)

	err := c.auth.Register(screenName, loginName, password)
	if err != nil {
		c.logger.Error("Registration failed", "error", err)
		return err
	}

	c.logger.Info("Registration successful", "login_name", loginName)
	return nil
}

// Logout clears authentication
func (c *Client) Logout() error {
	c.logger.Info("Logging out")

	c.auth.Logout()
	c.authenticated = false
	c.user = nil

	// Reconnect WebSocket without authentication
	if c.IsConnected() {
		// For now, just disconnect - could support anonymous connections
		c.Disconnect()
	}

	return nil
}

// IsAuthenticated returns whether the client is authenticated
func (c *Client) IsAuthenticated() bool {
	return c.authenticated && c.auth != nil && c.auth.IsAuthenticated()
}

// SendCommand sends a command to another user
func (c *Client) SendCommand(cmd Command) error {
	if !c.IsConnected() {
		return ErrNotConnected
	}

	return c.ws.SendMessage("command", cmd)
}

// RegisterCommandHandler registers a handler for a specific command type
func (c *Client) RegisterCommandHandler(cmdType string, handler CommandHandler) {
	c.commands.RegisterHandler(cmdType, handler)
}

// Events returns the event channel
func (c *Client) Events() <-chan Event {
	return c.eventChan
}

// Commands returns the command channel
func (c *Client) Commands() <-chan Command {
	return c.commandChan
}

// Errors returns the error channel
func (c *Client) Errors() <-chan error {
	return c.errorChan
}

// SetConfig updates the client configuration
func (c *Client) SetConfig(config *Config) error {
	c.config = config
	return nil
}

// GetConfig returns the current configuration
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetUser returns the current user (now uses server's User model)
func (c *Client) GetUser() *models.User {
	return c.user
}

// GetConnectionState returns the current connection state
func (c *Client) GetConnectionState() ConnectionState {
	if !c.IsConnected() {
		return StateDisconnected
	}

	if !c.IsAuthenticated() {
		return StateConnected
	}

	return StateAuthenticated
}

// processMessages processes incoming WebSocket messages
func (c *Client) processMessages() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case msg := <-c.ws.Messages():
			c.handleMessage(msg)
		case err := <-c.ws.Errors():
			c.logger.Error("WebSocket error", "error", err)
			c.emitError(err)
		}
	}
}

// handleMessage processes a WebSocket message
func (c *Client) handleMessage(msg WSMessage) {
	c.logger.Debug("Received message", "type", msg.Type)

	switch msg.Type {
	case "command":
		// Parse command and execute
		var cmd Command
		if err := ParsePayload(msg.Payload, &cmd); err != nil {
			c.logger.Error("Failed to parse command", "error", err)
			return
		}

		c.commandChan <- cmd
		c.handleCommand(cmd)

	case "pong":
		// Handle ping response
		c.logger.Debug("Received pong")

	default:
		c.logger.Warn("Unknown message type", "type", msg.Type)
	}
}

// handleCommand executes a received command
func (c *Client) handleCommand(cmd Command) {
	c.logger.Info("Executing command", "type", cmd.Type, "id", cmd.ID)

	c.emitEvent(EventCommandReceived, map[string]interface{}{
		"command_id":   cmd.ID,
		"command_type": cmd.Type,
	})

	// Process command
	result := c.commands.ProcessCommand(cmd)

	// Send result back to server
	if c.IsConnected() {
		c.ws.SendMessage("command_result", result)
	}

	// Emit event
	if result.Status == "completed" {
		c.emitEvent(EventCommandCompleted, map[string]interface{}{
			"command_id": cmd.ID,
			"result":     result.Result,
		})
	} else {
		c.emitEvent(EventCommandFailed, map[string]interface{}{
			"command_id": cmd.ID,
			"error":      result.Error,
		})
	}
}

// registerCoreHandlers registers the core command handlers
func (c *Client) registerCoreHandlers() {
	// std_ping - basic connectivity test
	c.RegisterCommandHandler("std_ping", func(cmd Command) CommandResult {
		return CommandResult{
			CommandID: cmd.ID,
			Status:    "completed",
			Result: map[string]interface{}{
				"pong_timestamp": time.Now().Format(time.RFC3339Nano),
				"latency_ms":     0, // Will be calculated by server
			},
		}
	})

	// Placeholder for other core handlers - will be implemented by UI layer
	c.RegisterCommandHandler("std_popup", func(cmd Command) CommandResult {
		return CommandResult{
			CommandID: cmd.ID,
			Status:    "failed",
			Error:     "std_popup requires platform-specific implementation",
		}
	})

	c.RegisterCommandHandler("std_notification", func(cmd Command) CommandResult {
		return CommandResult{
			CommandID: cmd.ID,
			Status:    "failed",
			Error:     "std_notification requires platform-specific implementation",
		}
	})

	c.RegisterCommandHandler("std_display_text", func(cmd Command) CommandResult {
		return CommandResult{
			CommandID: cmd.ID,
			Status:    "failed",
			Error:     "std_display_text requires platform-specific implementation",
		}
	})
}

// emitEvent sends an event to the event channel
func (c *Client) emitEvent(eventType string, data map[string]interface{}) {
	event := Event{
		Type:      eventType,
		Timestamp: time.Now(),
		Data:      data,
	}

	select {
	case c.eventChan <- event:
	default:
		c.logger.Warn("Event channel full, dropping event", "type", eventType)
	}
}

// emitError sends an error to the error channel
func (c *Client) emitError(err error) {
	select {
	case c.errorChan <- err:
	default:
		c.logger.Warn("Error channel full, dropping error", "error", err)
	}
}
