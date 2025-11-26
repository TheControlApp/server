package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Message represents the WebSocket message format
type Message struct {
	Type      string      `json:"type"`
	ID        string      `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Data      interface{} `json:"data"`
}

// Client represents our test WebSocket client
type TestClient struct {
	conn      *websocket.Conn
	userID    string
	token     string
	serverURL string
}

func main() {
	var (
		serverURL = flag.String("url", "ws://localhost:8080", "WebSocket server URL")
		token     = flag.String("token", "", "JWT token for authentication")
		userID    = flag.String("user", "", "User ID (optional, will generate if not provided)")
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("Token is required. Use -token=your_jwt_token")
	}

	if *userID == "" {
		*userID = uuid.New().String()
		fmt.Printf("Generated User ID: %s\n", *userID)
	}

	client := &TestClient{
		userID:    *userID,
		token:     *token,
		serverURL: *serverURL,
	}

	// Connect to server
	if err := client.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	fmt.Println("🎉 Connected to WebSocket server!")
	fmt.Println("Commands:")
	fmt.Println("  ping     - Send a ping message")
	fmt.Println("  status   - Send status update")
	fmt.Println("  raw      - Send raw JSON message")
	fmt.Println("  quit     - Exit the client")
	fmt.Println()

	// Start message receiver in background
	go client.ReadMessages()

	// Handle Ctrl+C gracefully
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// Command input loop
	scanner := bufio.NewScanner(os.Stdin)

	go func() {
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				return
			}

			command := strings.TrimSpace(scanner.Text())
			if command == "" {
				continue
			}

			switch strings.ToLower(command) {
			case "quit", "exit", "q":
				fmt.Println("Goodbye!")
				interrupt <- os.Interrupt
				return
			case "ping":
				client.SendPing()
			case "status":
				client.SendStatusUpdate()
			case "raw":
				client.SendRawMessage()
			default:
				fmt.Printf("Unknown command: %s\n", command)
				fmt.Println("Available commands: ping, status, raw, quit")
			}
		}
	}()

	// Wait for interrupt
	<-interrupt
	fmt.Println("\nShutting down client...")
}

// Connect establishes WebSocket connection
func (c *TestClient) Connect() error {
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %v", err)
	}

	// Add WebSocket path
	u.Path = "/ws/client"

	// Add token as query parameter
	q := u.Query()
	q.Set("token", c.token)
	u.RawQuery = q.Encode()

	fmt.Printf("Connecting to: %s\n", u.String())

	// Connect with headers
	headers := make(map[string][]string)
	headers["Authorization"] = []string{"Bearer " + c.token}

	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("connection failed with status %d: %v", resp.StatusCode, err)
		}
		return fmt.Errorf("connection failed: %v", err)
	}

	c.conn = conn
	return nil
}

// Close closes the WebSocket connection
func (c *TestClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// ReadMessages reads incoming messages from server
func (c *TestClient) ReadMessages() {
	for {
		_, messageBytes, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		// Pretty print the received message
		var message map[string]interface{}
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			fmt.Printf("📨 Raw message: %s\n", string(messageBytes))
			continue
		}

		// Format and display
		messageType, _ := message["type"].(string)
		timestamp, _ := message["timestamp"].(string)

		fmt.Printf("\n📨 Received: %s (at %s)\n", messageType, timestamp)

		// Pretty print the full message
		prettyJSON, err := json.MarshalIndent(message, "   ", "  ")
		if err == nil {
			fmt.Printf("   %s\n", string(prettyJSON))
		} else {
			fmt.Printf("   %s\n", string(messageBytes))
		}

		fmt.Print("> ") // Re-prompt for next command
	}
}

// SendPing sends a ping message to server
func (c *TestClient) SendPing() {
	message := Message{
		Type:      "ping",
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		From:      c.userID,
		To:        "server",
		Data: map[string]interface{}{
			"client_info": map[string]interface{}{
				"version":  "1.0.0",
				"type":     "test-cli",
				"platform": "cli",
			},
			"sequence": time.Now().Unix(),
		},
	}

	c.sendMessage(message)
	fmt.Println("📤 Sent ping message")
}

// SendStatusUpdate sends a status update message
func (c *TestClient) SendStatusUpdate() {
	message := Message{
		Type:      "status_update",
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		From:      c.userID,
		To:        "server",
		Data: map[string]interface{}{
			"status":       "online",
			"availability": "available",
			"last_seen":    time.Now().Format(time.RFC3339),
		},
	}

	c.sendMessage(message)
	fmt.Println("📤 Sent status update")
}

// SendRawMessage allows sending custom JSON
func (c *TestClient) SendRawMessage() {
	fmt.Print("Enter JSON message: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return
	}

	rawJSON := strings.TrimSpace(scanner.Text())
	if rawJSON == "" {
		fmt.Println("Empty message, cancelled.")
		return
	}

	// Validate JSON
	var testMsg interface{}
	if err := json.Unmarshal([]byte(rawJSON), &testMsg); err != nil {
		fmt.Printf("Invalid JSON: %v\n", err)
		return
	}

	// Send raw message
	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(rawJSON)); err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
		return
	}

	fmt.Println("📤 Sent raw message")
}

// sendMessage sends a structured message to the server
func (c *TestClient) sendMessage(message Message) {
	data, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("Failed to marshal message: %v\n", err)
		return
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
		return
	}
}
