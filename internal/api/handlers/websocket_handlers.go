package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/thecontrolapp/server/internal/auth"
	"github.com/thecontrolapp/server/internal/models"
	"github.com/thecontrolapp/server/internal/services"
	wshub "github.com/thecontrolapp/server/internal/websocket"
)

type WebSocketHandlers struct {
	Hub         *wshub.Hub
	JWTManager  *auth.JWTManager
	UserService *services.UserService
}

func NewWebSocketHandlers(hub *wshub.Hub, jwtManager *auth.JWTManager, userService *services.UserService) *WebSocketHandlers {
	handlers := &WebSocketHandlers{
		Hub:         hub,
		JWTManager:  jwtManager,
		UserService: userService,
	}

	// Set this handler as the message handler for the hub
	hub.SetMessageHandler(handlers)

	return handlers
}

// WebSocket upgrader with proper security settings
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowedOrigins := []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		}
		if origin == "" {
			return true
		}
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				return true
			}
		}
		return false
	},
}

// HandleClientWebSocket handles WebSocket connections for all clients
// Supports both authenticated and anonymous connections
func (h *WebSocketHandlers) HandleClientWebSocket(c *gin.Context) {
	// Try to extract token (optional)
	token := h.extractToken(c)
	var userID uuid.UUID
	var authenticated bool

	// If token provided, try to authenticate
	if token != "" {
		if uid, err := h.validateToken(token); err == nil {
			userID = uid
			authenticated = true
		}
		// If token is invalid, we continue with anonymous connection
	}

	// Upgrade to WebSocket (always allow, regardless of auth status)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	// Create client (authenticated or anonymous)
	client := wshub.NewClient(conn, userID, token, h.Hub)
	client.SetAuthenticated(authenticated)

	// Register client and start message handling
	h.Hub.RegisterClient(client)

	log.Printf("WebSocket client connected - UserID: %v, Authenticated: %v", userID, authenticated)

	// Start pumps
	go client.WritePump()
	go client.ReadPump()
}

func (h *WebSocketHandlers) extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		if strings.HasPrefix(authHeader, "Bearer ") {
			return strings.TrimPrefix(authHeader, "Bearer ")
		}
		return authHeader
	}
	return c.Query("token")
}

func (h *WebSocketHandlers) validateToken(tokenString string) (uuid.UUID, error) {
	claims, err := h.JWTManager.ValidateToken(tokenString)
	if err != nil {
		return uuid.Nil, err
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	return userID, nil
}

// HandleMessage implements the MessageHandler interface
func (h *WebSocketHandlers) HandleMessage(client *wshub.Client, message []byte) {
	// First, try to parse as a system message (auth, ping)
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error parsing message: %v", err)
		h.sendErrorMessage(client, "Invalid JSON message")
		return
	}

	// Handle system message types first
	if msgType, ok := msg["type"].(string); ok && h.isSystemMessage(msgType) {
		switch msgType {
		case "auth":
			h.handleAuthMessage(client, msg)
		case "ping":
			h.handlePingMessage(client, msg)
		}
		return
	}

	// Try to parse as a Command message
	h.handleCommandMessage(client, message)
}

// handleAuthMessage processes authentication messages within WebSocket connection
func (h *WebSocketHandlers) handleAuthMessage(client *wshub.Client, msg map[string]interface{}) {
	token, ok := msg["token"].(string)
	if !ok || token == "" {
		h.sendErrorMessage(client, "Missing token in auth message")
		return
	}

	// Validate the token
	userID, err := h.validateToken(token)
	if err != nil {
		h.sendErrorMessage(client, "Invalid or expired token")
		return
	}

	// Update client authentication status
	client.SetUserID(userID)
	client.SetToken(token)
	client.SetAuthenticated(true)

	// Send success response
	response := map[string]interface{}{
		"type":    "auth_success",
		"message": "Authentication successful",
		"user_id": userID.String(),
	}
	h.sendMessage(client, response)

	log.Printf("WebSocket client authenticated - UserID: %v", userID)
}

// handlePingMessage processes ping messages
func (h *WebSocketHandlers) handlePingMessage(client *wshub.Client, msg map[string]interface{}) {
	response := map[string]interface{}{
		"type":      "pong",
		"timestamp": time.Now().Unix(),
	}
	h.sendMessage(client, response)
}

// isSystemMessage checks if a message type is a system message (not a command)
func (h *WebSocketHandlers) isSystemMessage(msgType string) bool {
	systemTypes := map[string]bool{
		"ping": true,
		"auth": true,
	}
	return systemTypes[msgType]
}

// sendMessage sends a JSON message to the client
func (h *WebSocketHandlers) sendMessage(client *wshub.Client, data interface{}) {
	messageBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	select {
	case client.Send() <- messageBytes:
	default:
		close(client.Send())
		// Use the hub's unregister channel
		client.GetHub().UnregisterClient(client)
	}
}

// sendErrorMessage sends an error message to the client
func (h *WebSocketHandlers) sendErrorMessage(client *wshub.Client, message string) {
	errorMsg := map[string]interface{}{
		"type":    "error",
		"message": message,
	}
	h.sendMessage(client, errorMsg)
}

// handleCommandMessage processes command messages and creates proper Command structures
func (h *WebSocketHandlers) handleCommandMessage(client *wshub.Client, message []byte) {
	// Try to parse as legacy format first (for backward compatibility)
	var legacyMsg map[string]interface{}
	if err := json.Unmarshal(message, &legacyMsg); err == nil {
		if msgType, ok := legacyMsg["type"].(string); ok {
			// Convert legacy message to proper Command structure
			instruction := models.Instruction{
				Type:    msgType,
				Content: legacyMsg["content"],
			}

			command := models.Command{
				ID:           uuid.New(),
				Instructions: []models.Instruction{instruction},
				SenderID:     client.GetUserID(),
				ReceiverID:   nil,       // No specific receiver for legacy messages
				Tags:         "general", // Default tag for legacy messages
				Status:       "pending",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}

			h.broadcastCommand(client, command)
			return
		}
	}

	// Try to parse as proper Command structure
	var command models.Command
	if err := json.Unmarshal(message, &command); err != nil {
		log.Printf("Error parsing command: %v", err)
		h.sendErrorMessage(client, "Invalid command format")
		return
	}

	// Validate required fields
	if len(command.Instructions) == 0 {
		h.sendErrorMessage(client, "Command must contain at least one instruction")
		return
	}

	if command.Tags == "" {
		h.sendErrorMessage(client, "Tags are required")
		return
	}

	// Validate each instruction
	for i, instruction := range command.Instructions {
		if instruction.Type == "" {
			h.sendErrorMessage(client, fmt.Sprintf("Instruction %d is missing required 'type' field", i))
			return
		}
		if instruction.Content == nil {
			h.sendErrorMessage(client, fmt.Sprintf("Instruction %d is missing required 'content' field", i))
			return
		}
	}

	// Set server-controlled fields
	command.ID = uuid.New()
	command.SenderID = client.GetUserID()
	command.Status = "pending"
	command.CreatedAt = time.Now()
	command.UpdatedAt = time.Now()

	// Ensure ReceiverID is properly set (nil if not specified)
	if command.ReceiverID != nil && *command.ReceiverID == uuid.Nil {
		command.ReceiverID = nil
	}

	// Clear relationship data to avoid garbage in broadcast
	command.Sender = models.User{}
	command.Receiver = nil

	// Check authentication requirements
	if !client.IsAuthenticated() && h.requiresAuthenticationForCommand(&command) {
		h.sendErrorMessage(client, "Authentication required for this command")
		return
	}

	log.Printf("Processing command with %d instructions from user: %v", len(command.Instructions), client.GetUserID())

	h.broadcastCommand(client, command)
}

// broadcastCommand broadcasts a command to all connected clients
func (h *WebSocketHandlers) broadcastCommand(client *wshub.Client, command models.Command) {
	// Create a clean broadcast structure without relationship data
	broadcastCommand := struct {
		ID           uuid.UUID            `json:"id"`
		Instructions []models.Instruction `json:"instructions"`
		SenderID     uuid.UUID            `json:"sender_id"`
		ReceiverID   *uuid.UUID           `json:"receiver_id"`
		Tags         string               `json:"tags"`
		Status       string               `json:"status"`
		CreatedAt    time.Time            `json:"created_at"`
		UpdatedAt    time.Time            `json:"updated_at"`
	}{
		ID:           command.ID,
		Instructions: command.Instructions,
		SenderID:     command.SenderID,
		ReceiverID:   command.ReceiverID,
		Tags:         command.Tags,
		Status:       command.Status,
		CreatedAt:    command.CreatedAt,
		UpdatedAt:    command.UpdatedAt,
	}

	commandBytes, err := json.Marshal(broadcastCommand)
	if err != nil {
		log.Printf("Error marshaling command: %v", err)
		h.sendErrorMessage(client, "Failed to process command")
		return
	}

	h.Hub.BroadcastRaw(commandBytes)
}

// requiresAuthenticationForCommand checks if a command requires authentication
func (h *WebSocketHandlers) requiresAuthenticationForCommand(command *models.Command) bool {
	// For now, require authentication for all commands except basic ones
	publicInstructionTypes := map[string]bool{
		"ping":     true,
		"std_test": true,
	}

	for _, instruction := range command.Instructions {
		if !publicInstructionTypes[instruction.Type] {
			return true
		}
	}
	return false
}
