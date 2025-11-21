package client

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocketManager handles WebSocket connection to the ControlApp server
type WebSocketManager struct {
	config    *Config
	conn      *websocket.Conn
	connected bool
	mu        sync.RWMutex
	logger    Logger

	// Channels
	messageChan chan WSMessage
	errorChan   chan error
	closeChan   chan struct{}

	// Heartbeat
	lastPong   time.Time
	pingTicker *time.Ticker

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(config *Config, logger Logger) *WebSocketManager {
	return &WebSocketManager{
		config:      config,
		logger:      logger,
		messageChan: make(chan WSMessage, 100),
		errorChan:   make(chan error, 100),
		closeChan:   make(chan struct{}),
	}
}

// Connect establishes WebSocket connection to the server
func (ws *WebSocketManager) Connect(ctx context.Context, serverURL string, auth *AuthManager) error {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	// Create cancellable context
	ws.ctx, ws.cancel = context.WithCancel(ctx)

	// Parse server URL and build WebSocket URL
	u, err := url.Parse(serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}

	// Convert HTTP(S) to WS(S)
	scheme := "ws"
	if u.Scheme == "https" {
		scheme = "wss"
	}

	wsURL := fmt.Sprintf("%s://%s/ws/client", scheme, u.Host)
	wsU, err := url.Parse(wsURL)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL: %w", err)
	}

	// Add JWT token if available
	if auth != nil {
		if token := auth.GetToken(); token != "" {
			q := wsU.Query()
			q.Set("token", token)
			wsU.RawQuery = q.Encode()
		}
	}

	ws.logger.Info("Connecting to WebSocket", "url", wsU.String())

	// Dial WebSocket connection
	conn, _, err := websocket.DefaultDialer.Dial(wsU.String(), nil)
	if err != nil {
		return &NetworkError{
			Operation: "websocket_connect",
			URL:       wsU.String(),
			Cause:     err,
		}
	}

	ws.conn = conn
	ws.connected = true
	ws.lastPong = time.Now()

	// Start goroutines for message handling
	go ws.readPump()
	go ws.writePump()
	go ws.heartbeat()

	ws.logger.Info("WebSocket connected successfully")
	return nil
}

// Disconnect closes the WebSocket connection
func (ws *WebSocketManager) Disconnect() {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	if ws.cancel != nil {
		ws.cancel()
	}

	if ws.conn != nil {
		ws.logger.Info("Disconnecting WebSocket")
		ws.conn.Close()
		ws.conn = nil
	}

	ws.connected = false

	if ws.pingTicker != nil {
		ws.pingTicker.Stop()
	}

	close(ws.closeChan)
}

// IsConnected returns whether the WebSocket is connected
func (ws *WebSocketManager) IsConnected() bool {
	ws.mu.RLock()
	defer ws.mu.RUnlock()
	return ws.connected && ws.conn != nil
}

// SendMessage sends a message through the WebSocket
func (ws *WebSocketManager) SendMessage(msgType string, payload interface{}) error {
	ws.mu.RLock()
	defer ws.mu.RUnlock()

	if !ws.connected || ws.conn == nil {
		return ErrNotConnected
	}

	msg := WSMessage{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return ws.conn.WriteJSON(msg)
}

// Messages returns the message channel
func (ws *WebSocketManager) Messages() <-chan WSMessage {
	return ws.messageChan
}

// Errors returns the error channel
func (ws *WebSocketManager) Errors() <-chan error {
	return ws.errorChan
}

// Authenticate sends authentication token if connected
func (ws *WebSocketManager) Authenticate(token string) error {
	if !ws.IsConnected() {
		return ErrNotConnected
	}

	return ws.SendMessage("auth", map[string]string{
		"token": token,
	})
}

// readPump reads messages from the WebSocket connection
func (ws *WebSocketManager) readPump() {
	defer func() {
		ws.mu.Lock()
		ws.connected = false
		ws.mu.Unlock()

		if ws.conn != nil {
			ws.conn.Close()
		}
	}()

	// Configure connection limits and timeouts
	ws.conn.SetReadLimit(512 * 1024) // 512KB max message size
	ws.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	ws.conn.SetPongHandler(func(string) error {
		ws.lastPong = time.Now()
		ws.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		select {
		case <-ws.ctx.Done():
			return
		default:
		}

		var msg WSMessage
		err := ws.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ws.logger.Error("WebSocket error", "error", err)
				ws.errorChan <- err
			}
			return
		}

		ws.logger.Debug("Received WebSocket message", "type", msg.Type)

		// Handle special message types
		switch msg.Type {
		case "pong":
			ws.lastPong = time.Now()
		default:
			// Send to message channel
			select {
			case ws.messageChan <- msg:
			case <-ws.ctx.Done():
				return
			default:
				ws.logger.Warn("Message channel full, dropping message", "type", msg.Type)
			}
		}
	}
}

// writePump handles writing messages to the WebSocket
func (ws *WebSocketManager) writePump() {
	ticker := time.NewTicker(54 * time.Second) // Ping every 54 seconds
	defer func() {
		ticker.Stop()
		if ws.conn != nil {
			ws.conn.Close()
		}
	}()

	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-ticker.C:
			ws.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := ws.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				ws.logger.Debug("Failed to send ping", "error", err)
				return
			}
		}
	}
}

// heartbeat monitors connection health
func (ws *WebSocketManager) heartbeat() {
	ws.pingTicker = time.NewTicker(30 * time.Second)
	defer ws.pingTicker.Stop()

	for {
		select {
		case <-ws.ctx.Done():
			return
		case <-ws.pingTicker.C:
			// Check if we received a pong recently
			if time.Since(ws.lastPong) > 90*time.Second {
				ws.logger.Warn("WebSocket heartbeat timeout")
				ws.errorChan <- fmt.Errorf("heartbeat timeout")
				return
			}
		}
	}
}
