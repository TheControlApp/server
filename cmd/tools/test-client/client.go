package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

type Client struct {
	serverURL string
	httpURL   string
	conn      *websocket.Conn
}

type AuthResponse struct {
	Token string `json:"token"`
}

func NewClient(serverURL string) *Client {
	// Convert wss:// to https:// for HTTP requests
	httpURL := strings.Replace(serverURL, "wss://", "https://", 1)
	httpURL = strings.Replace(httpURL, "ws://", "http://", 1)

	return &Client{
		serverURL: serverURL,
		httpURL:   httpURL,
	}
}

func (c *Client) Login(username, password string) (string, error) {
	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.httpURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", err
	}

	return authResp.Token, nil
}

func (c *Client) Register(username, password string) (string, error) {
	registerData := map[string]interface{}{
		"username":      username,
		"password":      password,
		"screen_name":   username, // Use username as screen_name for test client
		"random_opt_in": false,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(c.httpURL+"/api/v1/auth/register", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("registration failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Registration doesn't return a token, so we need to login
	return c.Login(username, password)
}

func (c *Client) Connect(token string) error {
	// Parse the URL and add token as query parameter
	u, err := url.Parse(c.serverURL + "/ws/client")
	if err != nil {
		return err
	}

	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()

	// Configure dialer with larger buffer sizes
	dialer := websocket.Dialer{
		ReadBufferSize:  10 * 1024 * 1024, // 10MB
		WriteBufferSize: 10 * 1024 * 1024, // 10MB
	}

	// Connect to WebSocket
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
