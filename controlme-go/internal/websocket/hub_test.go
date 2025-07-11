package websocket

import (
	"testing"

	"github.com/google/uuid"
)

func TestHub_UserConnectionLifecycle(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	client := &Client{
		userID:     userID,
		clientType: "web",
		send:       make(chan []byte, 1),
		hub:        hub,
	}

	// Register client
	hub.register <- client

	if !hub.IsUserConnected(userID) {
		t.Errorf("User should be connected after registration")
	}
	if hub.GetUserConnections(userID) != 1 {
		t.Errorf("User should have 1 connection")
	}

	// Unregister client
	hub.unregister <- client

	if hub.IsUserConnected(userID) {
		t.Errorf("User should not be connected after unregister")
	}
	if hub.GetUserConnections(userID) != 0 {
		t.Errorf("User should have 0 connections after unregister")
	}
}
