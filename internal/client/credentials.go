package client

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Credentials represents stored authentication credentials
type Credentials struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
	ServerURL string    `json:"server_url"`
}

// CredentialStore manages credential persistence
type CredentialStore struct {
	configDir  string
	configFile string
}

// NewCredentialStore creates a new credential store
func NewCredentialStore() *CredentialStore {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".controlapp")

	return &CredentialStore{
		configDir:  configDir,
		configFile: filepath.Join(configDir, "credentials.json"),
	}
}

// Save stores credentials to disk
func (cs *CredentialStore) Save(creds *Credentials) error {
	// Create config directory if it doesn't exist
	if err := os.MkdirAll(cs.configDir, 0700); err != nil {
		return err
	}

	// Marshal credentials to JSON
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	// Write to file with restricted permissions
	return os.WriteFile(cs.configFile, data, 0600)
}

// Load retrieves credentials from disk
func (cs *CredentialStore) Load() (*Credentials, error) {
	// Check if file exists
	if _, err := os.Stat(cs.configFile); os.IsNotExist(err) {
		return nil, nil // No credentials found
	}

	// Read file
	data, err := os.ReadFile(cs.configFile)
	if err != nil {
		return nil, err
	}

	// Unmarshal JSON
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

// IsValid checks if stored credentials are still valid
func (cs *CredentialStore) IsValid(creds *Credentials) bool {
	if creds == nil {
		return false
	}

	// Check if token has expired (with 1 hour buffer)
	if time.Now().Add(1 * time.Hour).After(creds.ExpiresAt) {
		return false
	}

	// Check if essential fields are present
	if creds.Token == "" || creds.Username == "" {
		return false
	}

	return true
}

// Clear removes stored credentials
func (cs *CredentialStore) Clear() error {
	if _, err := os.Stat(cs.configFile); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to clear
	}

	return os.Remove(cs.configFile)
}

// UpdateToken updates just the token and expiration time
func (cs *CredentialStore) UpdateToken(token string, expiresAt time.Time) error {
	creds, err := cs.Load()
	if err != nil || creds == nil {
		return err
	}

	creds.Token = token
	creds.ExpiresAt = expiresAt

	return cs.Save(creds)
}
