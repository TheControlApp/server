package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/thecontrolapp/server/internal/models"
)

// AuthResponse represents the authentication response from the server
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Commander represents our test commander client
type Commander struct {
	conn      *websocket.Conn
	token     string
	user      *models.User
	serverURL string
	wsURL     string
}

// Login authenticates with the server and retrieves a JWT token
func (c *Commander) Login(username, password string) error {
	loginURL := fmt.Sprintf("%s/api/v1/auth/login", c.serverURL)

	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login data: %w", err)
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.token = authResp.Token
	c.user = &authResp.User

	return nil
}

// Register creates a new account on the server
func (c *Commander) Register(screenName, loginName, password string) error {
	registerURL := fmt.Sprintf("%s/api/v1/auth/register", c.serverURL)

	registerData := map[string]interface{}{
		"screen_name":   screenName,
		"username":      loginName,
		"password":      password,
		"random_opt_in": false,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		return fmt.Errorf("failed to marshal register data: %w", err)
	}

	resp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send register request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// Connect establishes WebSocket connection
func (c *Commander) Connect() error {
	// Convert HTTP URL to WebSocket URL
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}

	wsScheme := "ws"
	if u.Scheme == "https" {
		wsScheme = "wss"
	}

	c.wsURL = fmt.Sprintf("%s://%s/ws/client?token=%s", wsScheme, u.Host, c.token)

	// Configure dialer with larger read buffer
	dialer := websocket.Dialer{
		ReadBufferSize:  10 * 1024 * 1024, // 10MB
		WriteBufferSize: 10 * 1024 * 1024, // 10MB
	}

	// Connect to WebSocket
	conn, _, err := dialer.Dial(c.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Set read limit to 100MB to handle large server responses
	conn.SetReadLimit(100 * 1024 * 1024)

	c.conn = conn
	return nil
}

// Close closes the WebSocket connection
func (c *Commander) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// SendCommand sends a command to the WebSocket
func (c *Commander) SendCommand(cmd models.Command) error {
	// Ensure required fields are set
	if cmd.ID == uuid.Nil {
		cmd.ID = uuid.New()
	}
	if cmd.SenderID == uuid.Nil {
		cmd.SenderID = c.user.ID
	}
	if cmd.Tags == "" {
		cmd.Tags = "test"
	}
	if cmd.Status == "" {
		cmd.Status = "pending"
	}
	if cmd.CreatedAt.IsZero() {
		cmd.CreatedAt = time.Now()
	}
	if cmd.UpdatedAt.IsZero() {
		cmd.UpdatedAt = time.Now()
	}

	// Marshal and send
	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	return nil
}
