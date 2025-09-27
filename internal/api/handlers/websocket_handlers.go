package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	wshub "github.com/thecontrolapp/controlme-go/internal/websocket"
)

type WebSocketHandlers struct {
	Hub        *wshub.Hub
	JWTManager *auth.JWTManager
}

func NewWebSocketHandlers(hub *wshub.Hub, jwtManager *auth.JWTManager) *WebSocketHandlers {
	return &WebSocketHandlers{
		Hub:        hub,
		JWTManager: jwtManager,
	}
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
func (h *WebSocketHandlers) HandleClientWebSocket(c *gin.Context) {
	h.handleWebSocketConnection(c)
}

// handleWebSocketConnection handles the WebSocket upgrade and authentication
func (h *WebSocketHandlers) handleWebSocketConnection(c *gin.Context) {
	token := h.extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authentication token"})
		return
	}
	userID, err := h.validateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := wshub.NewClient(conn, userID, token, h.Hub)
	h.Hub.RegisterClient(client)
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
