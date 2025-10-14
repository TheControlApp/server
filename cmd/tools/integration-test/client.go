package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// UserExistsError represents a user already exists condition
type UserExistsError struct {
	Username string
}

func (e *UserExistsError) Error() string {
	return fmt.Sprintf("user %s already exists", e.Username)
}

// TestClient represents a test client with full functionality
type TestClient struct {
	Name      string
	LoginName string
	Password  string
	Token     string
	User      User
	WSConn    *websocket.Conn
	Messages  []WSMessage
	mu        sync.Mutex
	listening bool
	stopChan  chan struct{}
}

// NewTestClient creates a new test client
func NewTestClient(name, loginName, password string) *TestClient {
	return &TestClient{
		Name:      name,
		LoginName: loginName,
		Password:  password,
		Messages:  make([]WSMessage, 0),
		stopChan:  make(chan struct{}),
	}
}

// Log outputs a timestamped log message for this client
func (tc *TestClient) Log(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05.000")
	prefix := fmt.Sprintf("[%s][%s]", timestamp, tc.Name)
	fmt.Printf(prefix+" "+format+"\n", args...)
}

// AddMessage safely adds a message to the client's message list
func (tc *TestClient) AddMessage(msg WSMessage) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.Messages = append(tc.Messages, msg)
}

// GetMessages safely returns a copy of all messages
func (tc *TestClient) GetMessages() []WSMessage {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	return append([]WSMessage{}, tc.Messages...)
}

// Register attempts to register the user
func (tc *TestClient) Register(serverURL string) error {
	payload := map[string]interface{}{
		"username":      tc.LoginName,
		"screen_name":   tc.Name,
		"password":      tc.Password,
		"random_opt_in": false,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal registration data: %v", err)
	}

	resp, err := http.Post(serverURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("registration request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		// Registration successful
		return nil
	} else if resp.StatusCode == 409 {
		// User already exists, return a special error type
		return &UserExistsError{Username: tc.LoginName}
	}

	return fmt.Errorf("registration failed with status %d", resp.StatusCode)
}

// Login authenticates the user and retrieves a JWT token
func (tc *TestClient) Login(serverURL string) error {
	fmt.Printf("   🔧 Attempting login with username: %s\n", tc.LoginName)

	payload := map[string]interface{}{
		"username": tc.LoginName,
		"password": tc.Password,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login data: %v", err)
	}

	resp, err := http.Post(serverURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("login request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("login failed with status %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %v", err)
	}

	tc.Token = authResp.Token
	tc.User = authResp.User

	fmt.Printf("   ✅ Login successful for username: %s, got token\n", tc.LoginName)

	return nil
}

// ConnectWebSocket establishes a WebSocket connection with authentication
func (tc *TestClient) ConnectWebSocket(wsBaseURL string) error {
	// Connect with token in query parameter
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws/client"}
	q := u.Query()
	q.Set("token", tc.Token)
	u.RawQuery = q.Encode()

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %v", err)
	}

	tc.WSConn = conn
	return nil
}

// AuthenticateWebSocket sends authentication via WebSocket message
func (tc *TestClient) AuthenticateWebSocket() error {
	if tc.WSConn == nil {
		return fmt.Errorf("no websocket connection")
	}

	authMsg := WSMessage{
		Type:  "auth",
		Token: tc.Token,
	}

	if err := tc.WSConn.WriteJSON(authMsg); err != nil {
		return fmt.Errorf("failed to send auth message: %v", err)
	}

	return nil
}

// StartMessageListener starts listening for WebSocket messages
func (tc *TestClient) StartMessageListener() {
	if tc.WSConn == nil || tc.listening {
		return
	}

	tc.listening = true
	tc.Log("🎧 Starting message listener...")

	for {
		select {
		case <-tc.stopChan:
			tc.Log("🔇 Message listener stopped")
			return
		default:
			var msg WSMessage
			err := tc.WSConn.ReadJSON(&msg)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					tc.Log("🔇 WebSocket closed normally")
				} else if strings.Contains(err.Error(), "use of closed network connection") {
					// Expected during cleanup, don't log as error
					tc.Log("🔇 Connection closed during cleanup")
				} else {
					tc.Log("❌ WebSocket read error: %v", err)
				}
				return
			}

			tc.Log("📨 Received: %s", msg.Type)
			if msg.Message != "" {
				tc.Log("   Message: %s", msg.Message)
			}
			if msg.UserID != "" {
				tc.Log("   User ID: %s", msg.UserID)
			}

			tc.AddMessage(msg)

			// Handle auth success
			if msg.Type == "auth_success" {
				tc.Log("🔓 Authentication successful!")
			}
		}
	}
}

// SendMessage sends a message via WebSocket
func (tc *TestClient) SendMessage(message map[string]interface{}) error {
	if tc.WSConn == nil {
		return fmt.Errorf("no websocket connection")
	}

	return tc.WSConn.WriteJSON(message)
}

// SendPing sends a ping message
func (tc *TestClient) SendPing() error {
	if tc.WSConn == nil {
		return fmt.Errorf("no websocket connection")
	}

	timestamp := fmt.Sprintf("%f", float64(time.Now().UnixNano())/1e9)
	return tc.WSConn.WriteMessage(websocket.PingMessage, []byte(timestamp))
}

// Cleanup closes connections and stops listeners
func (tc *TestClient) Cleanup() {
	if tc.listening {
		// First signal the listener to stop
		close(tc.stopChan)
		tc.listening = false

		// Give the listener a moment to stop gracefully
		time.Sleep(100 * time.Millisecond)
	}

	if tc.WSConn != nil {
		// Send a close message to the server before closing
		tc.WSConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		tc.WSConn.Close()
		tc.WSConn = nil
	}
}
