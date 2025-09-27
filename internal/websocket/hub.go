package websocket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// User connections map for targeted messaging
	userConnections map[uuid.UUID]map[*Client]bool

	// Token connections map for enforcing one connection per token
	tokenConnections map[string]*Client

	// Maximum connections per user (0 = unlimited)
	maxConnectionsPerUser int

	// Message cache for broadcast optimization
	messageCache map[string][]byte
}

// Client represents a WebSocket client
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// User ID
	userID uuid.UUID

	// JWT Token used for authentication
	token string

	// Buffered channel of outbound messages
	send chan []byte

	// Reference to the hub
	hub *Hub
}

// Message represents a WebSocket message
type Message struct {
	Type      string      `json:"type"`
	ID        string      `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	From      uuid.UUID   `json:"from"`
	To        uuid.UUID   `json:"to"`
	Data      interface{} `json:"data"`
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:               make(map[*Client]bool),
		broadcast:             make(chan []byte),
		register:              make(chan *Client),
		unregister:            make(chan *Client),
		userConnections:       make(map[uuid.UUID]map[*Client]bool),
		tokenConnections:      make(map[string]*Client),
		maxConnectionsPerUser: 0, // 0 = unlimited, can be configured
		messageCache:          make(map[string][]byte),
	}
}

// Run starts the WebSocket hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Check if token already has an active connection
			if existingClient, exists := h.tokenConnections[client.token]; exists {
				// Close the existing connection
				logrus.WithFields(logrus.Fields{
					"user_id":     existingClient.userID,
					"new_user_id": client.userID,
				}).Info("Replacing existing WebSocket connection for token")

				// Remove existing client
				delete(h.clients, existingClient)
				close(existingClient.send)
				existingClient.conn.Close()

				// Remove from user connections
				if connections, userExists := h.userConnections[existingClient.userID]; userExists {
					delete(connections, existingClient)
					if len(connections) == 0 {
						delete(h.userConnections, existingClient.userID)
					}
				}
			}

			// Check connection limits per user
			if h.maxConnectionsPerUser > 0 {
				if connections, exists := h.userConnections[client.userID]; exists {
					if len(connections) >= h.maxConnectionsPerUser {
						logrus.WithFields(logrus.Fields{
							"user_id":     client.userID,
							"max_allowed": h.maxConnectionsPerUser,
						}).Warn("Maximum connections per user exceeded")
						close(client.send)
						client.conn.Close()
						continue
					}
				}
			}

			// Register the new client
			h.clients[client] = true
			h.tokenConnections[client.token] = client

			// Add to user connections
			if _, exists := h.userConnections[client.userID]; !exists {
				h.userConnections[client.userID] = make(map[*Client]bool)
			}
			h.userConnections[client.userID][client] = true

			logrus.WithFields(logrus.Fields{
				"user_id":       client.userID,
				"total_clients": len(h.clients),
			}).Info("Client connected")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)

				// Remove from token connections
				delete(h.tokenConnections, client.token)

				// Remove from user connections
				if connections, exists := h.userConnections[client.userID]; exists {
					delete(connections, client)
					// Remove user from map if no more connections
					if len(connections) == 0 {
						delete(h.userConnections, client.userID)
					}
				}

				logrus.WithFields(logrus.Fields{
					"user_id":       client.userID,
					"total_clients": len(h.clients),
				}).Info("Client disconnected")
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// SendToUser sends a message to a specific user
func (h *Hub) SendToUser(userID uuid.UUID, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message")
		return
	}

	if connections, exists := h.userConnections[userID]; exists {
		for client := range connections {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(h.clients, client)
				delete(connections, client)
			}
		}
		// Clean up empty connection maps
		if len(connections) == 0 {
			delete(h.userConnections, userID)
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message")
		return
	}

	h.broadcast <- data
}

// GetConnectedUsers returns a list of connected user IDs
func (h *Hub) GetConnectedUsers() []uuid.UUID {
	users := make([]uuid.UUID, 0, len(h.userConnections))
	for userID := range h.userConnections {
		users = append(users, userID)
	}
	return users
}

// IsUserConnected checks if a user is connected
func (h *Hub) IsUserConnected(userID uuid.UUID) bool {
	_, exists := h.userConnections[userID]
	return exists
}

// GetUserConnections returns the number of connections for a user
func (h *Hub) GetUserConnections(userID uuid.UUID) int {
	if connections, exists := h.userConnections[userID]; exists {
		return len(connections)
	}
	return 0
}

// NewClient creates a new WebSocket client with authentication
func NewClient(conn *websocket.Conn, userID uuid.UUID, token string, hub *Hub) *Client {
	return &Client{
		conn:   conn,
		userID: userID,
		token:  token,
		send:   make(chan []byte, 256),
		hub:    hub,
	}
}

// RegisterClient registers a new authenticated client with the hub
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// SetMaxConnectionsPerUser sets the maximum number of connections allowed per user
func (h *Hub) SetMaxConnectionsPerUser(max int) {
	h.maxConnectionsPerUser = max
}

// GetStats returns hub statistics
func (h *Hub) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_clients":      len(h.clients),
		"total_users":        len(h.userConnections),
		"active_tokens":      len(h.tokenConnections),
		"max_per_user":       h.maxConnectionsPerUser,
		"message_cache_size": len(h.messageCache),
	}
}

// CleanupStaleConnections removes connections that are no longer valid
func (h *Hub) CleanupStaleConnections() {
	for client := range h.clients {
		// Send a ping to test if connection is alive
		if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			h.unregister <- client
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection
func (c *Client) WritePump() {
	const (
		writeWait  = 10 * time.Second    // Time allowed to write a message to peer
		pongWait   = 60 * time.Second    // Time allowed to read the next pong message from peer
		pingPeriod = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait
	)

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logrus.WithError(err).Error("Failed to write message to WebSocket")
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logrus.WithError(err).Debug("Failed to send ping message")
				return
			}
		}
	}
}

// ReadPump pumps messages from the websocket connection to the hub
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Configure connection timeouts and ping/pong handling
	const (
		writeWait      = 10 * time.Second    // Time allowed to write a message to peer
		pongWait       = 60 * time.Second    // Time allowed to read the next pong message from peer
		pingPeriod     = (pongWait * 9) / 10 // Send pings to peer with this period. Must be less than pongWait
		maxMessageSize = 512                 // Maximum message size allowed from peer
	)

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithError(err).Error("WebSocket error")
			}
			break
		}

		// Handle incoming messages (TODO: implement message handling)
		logrus.WithFields(logrus.Fields{
			"user_id": c.userID,
			"message": string(message),
		}).Debug("Received WebSocket message")

		// Echo message back for now (TODO: implement proper message routing)
		c.send <- message
	}
}
