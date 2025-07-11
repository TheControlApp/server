package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
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
	userConnections map[uuid.UUID][]*Client
}

// Client represents a WebSocket client
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// User ID
	userID uuid.UUID

	// Client type (web, desktop)
	clientType string

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
		clients:         make(map[*Client]bool),
		broadcast:       make(chan []byte),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		userConnections: make(map[uuid.UUID][]*Client),
	}
}

// Run starts the WebSocket hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			
			// Add to user connections
			if _, exists := h.userConnections[client.userID]; !exists {
				h.userConnections[client.userID] = []*Client{}
			}
			h.userConnections[client.userID] = append(h.userConnections[client.userID], client)
			
			logrus.WithFields(logrus.Fields{
				"user_id":     client.userID,
				"client_type": client.clientType,
				"total_clients": len(h.clients),
			}).Info("Client connected")

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				
				// Remove from user connections
				if connections, exists := h.userConnections[client.userID]; exists {
					for i, conn := range connections {
						if conn == client {
							h.userConnections[client.userID] = append(connections[:i], connections[i+1:]...)
							break
						}
					}
					// Remove user from map if no more connections
					if len(h.userConnections[client.userID]) == 0 {
						delete(h.userConnections, client.userID)
					}
				}
				
				logrus.WithFields(logrus.Fields{
					"user_id":     client.userID,
					"client_type": client.clientType,
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
		for _, client := range connections {
			select {
			case client.send <- data:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

// SendToUserByType sends a message to a specific user's client type
func (h *Hub) SendToUserByType(userID uuid.UUID, clientType string, message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal message")
		return
	}

	if connections, exists := h.userConnections[userID]; exists {
		for _, client := range connections {
			if client.clientType == clientType {
				select {
				case client.send <- data:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
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
