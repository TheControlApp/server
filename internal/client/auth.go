package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/thecontrolapp/server/internal/models"
)

// AuthManager handles authentication with the ControlApp server
type AuthManager struct {
	config    *Config
	client    *http.Client
	token     string
	user      *models.User
	expiresAt time.Time
	logger    Logger
}

// NewAuthManager creates a new authentication manager
func NewAuthManager(config *Config, logger Logger) *AuthManager {
	return &AuthManager{
		config: config,
		client: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		logger: logger,
	}
}

// Login authenticates with the server using username and password
func (a *AuthManager) Login(username, password string) error {
	a.logger.Info("Attempting to login", "username", username)

	req := LoginRequest{
		LoginName: username,
		Password:  password,
	}

	resp, err := a.makeAuthRequest("POST", "/api/v1/auth/login", req)
	if err != nil {
		return fmt.Errorf("login request failed: %w", err)
	}

	var authResp AuthResponse
	if err := json.Unmarshal(resp, &authResp); err != nil {
		return fmt.Errorf("failed to parse login response: %w", err)
	}

	a.token = authResp.Token
	a.user = &authResp.User

	// Set a default expiration of 1 hour (JWT parsing can be added later)
	a.expiresAt = time.Now().Add(time.Hour)

	a.logger.Info("Login successful", "user_id", a.user.ID)
	return nil
}

// Register creates a new user account
func (a *AuthManager) Register(screenName, loginName, password string) error {
	a.logger.Info("Attempting to register", "login_name", loginName)

	req := RegisterRequest{
		ScreenName: screenName,
		LoginName:  loginName,
		Password:   password,
	}

	_, err := a.makeAuthRequest("POST", "/api/v1/auth/register", req)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	a.logger.Info("Registration successful", "login_name", loginName)
	return nil
}

// Logout clears the authentication state
func (a *AuthManager) Logout() {
	a.logger.Info("Logging out")
	a.token = ""
	a.user = nil
	a.expiresAt = time.Time{}
}

// GetToken returns the current JWT token if valid
func (a *AuthManager) GetToken() string {
	if a.isTokenExpired() {
		return ""
	}
	return a.token
}

// GetUser returns the current authenticated user
func (a *AuthManager) GetUser() *models.User {
	return a.user
}

// IsAuthenticated returns whether the user is authenticated with a valid token
func (a *AuthManager) IsAuthenticated() bool {
	return a.token != "" && !a.isTokenExpired() && a.user != nil
}

// makeAuthRequest makes an HTTP request to the authentication API
func (a *AuthManager) makeAuthRequest(method, endpoint string, body interface{}) ([]byte, error) {
	// Build full URL - convert WebSocket URL to HTTP
	serverURL := strings.TrimSuffix(a.config.ServerURL, "/ws/client")
	serverURL = strings.Replace(serverURL, "ws://", "http://", 1)
	serverURL = strings.Replace(serverURL, "wss://", "https://", 1)

	fullURL := serverURL + endpoint

	// Marshal request body
	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	// Create request
	req, err := http.NewRequest(method, fullURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if a.token != "" {
		req.Header.Set("Authorization", "Bearer "+a.token)
	}

	// Make request
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, &NetworkError{
			Operation: method + " " + endpoint,
			URL:       fullURL,
			Cause:     err,
		}
	}
	defer resp.Body.Close()

	// Read response
	respBody := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			respBody = append(respBody, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	// Check status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// isTokenExpired checks if the current token has expired
func (a *AuthManager) isTokenExpired() bool {
	if a.expiresAt.IsZero() {
		return true
	}
	// Add 30 second buffer to avoid edge cases
	return time.Now().Add(30 * time.Second).After(a.expiresAt)
}
