package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type      string      `json:"type"`
	Token     string      `json:"token,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func main() {
	fmt.Println("WebSocket Authentication Test Client")
	fmt.Println("=====================================")

	// Connect to WebSocket server (anonymous connection)
	serverURL := "ws://localhost:8080/ws"
	fmt.Printf("Connecting to %s...\n", serverURL)

	conn, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Failed to connect to WebSocket server:", err)
	}
	defer conn.Close()

	fmt.Println("✅ Connected successfully (anonymous session)")
	fmt.Println()
	fmt.Println("Available commands:")
	fmt.Println("  ping           - Send a ping message")
	fmt.Println("  auth <token>   - Authenticate with JWT token")
	fmt.Println("  quit           - Exit the client")
	fmt.Println()

	// Start message reading goroutine
	go readMessages(conn)

	// Handle user input
	scanner := bufio.NewScanner(os.Stdin)

	// Handle Ctrl+C
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	for {
		fmt.Print("> ")

		select {
		case <-interrupt:
			fmt.Println("\nReceived interrupt, closing connection...")
			return
		default:
			if scanner.Scan() {
				input := strings.TrimSpace(scanner.Text())
				if input == "" {
					continue
				}

				parts := strings.SplitN(input, " ", 2)
				command := parts[0]

				switch command {
				case "quit", "exit":
					fmt.Println("Goodbye!")
					return

				case "ping":
					sendPing(conn)

				case "auth":
					if len(parts) < 2 {
						fmt.Println("Usage: auth <token>")
						continue
					}
					sendAuth(conn, parts[1])

				default:
					fmt.Printf("Unknown command: %s\n", command)
				}
			}
		}
	}
}

func readMessages(conn *websocket.Conn) {
	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("WebSocket error: %v\n", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			fmt.Printf("📨 Raw message: %s\n", string(messageBytes))
			continue
		}

		switch msg.Type {
		case "pong":
			fmt.Printf("🏓 Pong received (timestamp: %d)\n", msg.Timestamp)
		case "auth_success":
			fmt.Printf("✅ Authentication successful! %s\n", msg.Message)
		case "error":
			fmt.Printf("❌ Error: %s\n", msg.Message)
		default:
			fmt.Printf("📨 Received: %s\n", string(messageBytes))
		}
	}
}

func sendPing(conn *websocket.Conn) {
	msg := Message{
		Type:      "ping",
		Timestamp: time.Now().Unix(),
	}

	if err := conn.WriteJSON(msg); err != nil {
		fmt.Printf("Error sending ping: %v\n", err)
		return
	}

	fmt.Println("🏓 Ping sent")
}

func sendAuth(conn *websocket.Conn, token string) {
	msg := Message{
		Type:  "auth",
		Token: token,
	}

	if err := conn.WriteJSON(msg); err != nil {
		fmt.Printf("Error sending auth: %v\n", err)
		return
	}

	fmt.Printf("🔐 Authentication request sent with token: %s...\n", token[:min(len(token), 20)])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
