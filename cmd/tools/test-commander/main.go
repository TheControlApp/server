package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/thecontrolapp/server/internal/models"
)

// AuthResponse represents the authentication response from the server
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Commander represents our test commander client
type Commander struct {
	conn      *websocket.Conn
	token     string
	user      *models.User
	serverURL string
	wsURL     string
}

func main() {
	var (
		serverURL = flag.String("url", "https://ctrlapp.merith.xyz", "Server base URL")
		username  = flag.String("username", "", "Login username")
		password  = flag.String("password", "", "Login password")
		register  = flag.Bool("register", false, "Register account if it doesn't exist")
	)
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("Username and password are required. Use -username and -password flags")
	}

	commander := &Commander{
		serverURL: *serverURL,
	}

	// Login first
	fmt.Println("🔐 Authenticating...")
	if err := commander.Login(*username, *password); err != nil {
		if *register {
			// If login failed and register flag is set, try to register the account
			fmt.Println("⚠️  Login failed, attempting to register new account...")
			screenName := strings.Title(*username) + " (Commander)"
			if err := commander.Register(screenName, *username, *password); err != nil {
				log.Fatalf("Failed to register: %v", err)
			}
			fmt.Println("✅ Account registered successfully!")

			// Now try to login again
			if err := commander.Login(*username, *password); err != nil {
				log.Fatalf("Failed to login after registration: %v", err)
			}
		} else {
			log.Fatalf("Failed to authenticate: %v\nTip: Use --register flag to create account if it doesn't exist", err)
		}
	}

	fmt.Printf("✅ Logged in as: %s (%s)\n", commander.user.ScreenName, commander.user.LoginName)
	fmt.Printf("   User ID: %s\n\n", commander.user.ID)

	// Connect to WebSocket
	fmt.Println("🔌 Connecting to WebSocket...")
	if err := commander.Connect(); err != nil {
		log.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer commander.Close()

	fmt.Println("✅ Connected to WebSocket server!")
	fmt.Println()

	// Start message receiver
	go commander.ReadMessages()

	// Handle Ctrl+C gracefully
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func() {
		<-interrupt
		fmt.Println("\n👋 Goodbye!")
		commander.Close()
		os.Exit(0)
	}()

	// Start interactive menu
	commander.InteractiveMenu()
}

// Login authenticates with the server and retrieves a JWT token
func (c *Commander) Login(username, password string) error {
	loginURL := fmt.Sprintf("%s/api/v1/auth/login", c.serverURL)

	loginData := map[string]string{
		"username": username,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		return fmt.Errorf("failed to marshal login data: %w", err)
	}

	resp, err := http.Post(loginURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to decode auth response: %w", err)
	}

	c.token = authResp.Token
	c.user = &authResp.User

	return nil
}

// Register creates a new account on the server
func (c *Commander) Register(screenName, loginName, password string) error {
	registerURL := fmt.Sprintf("%s/api/v1/auth/register", c.serverURL)

	registerData := map[string]string{
		"screen_name": screenName,
		"username":    loginName,
		"password":    password,
	}

	jsonData, err := json.Marshal(registerData)
	if err != nil {
		return fmt.Errorf("failed to marshal register data: %w", err)
	}

	resp, err := http.Post(registerURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send register request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed with status: %d", resp.StatusCode)
	}

	return nil
}

// Connect establishes WebSocket connection
func (c *Commander) Connect() error {
	// Convert HTTP URL to WebSocket URL
	u, err := url.Parse(c.serverURL)
	if err != nil {
		return fmt.Errorf("invalid server URL: %w", err)
	}

	wsScheme := "ws"
	if u.Scheme == "https" {
		wsScheme = "wss"
	}

	c.wsURL = fmt.Sprintf("%s://%s/ws/client?token=%s", wsScheme, u.Host, c.token)

	// Configure dialer with larger read buffer
	dialer := websocket.Dialer{
		ReadBufferSize:  1024 * 1024, // 1MB
		WriteBufferSize: 1024 * 1024, // 1MB
	}

	// Connect to WebSocket
	conn, _, err := dialer.Dial(c.wsURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Set read limit to 10MB
	conn.SetReadLimit(10 * 1024 * 1024)

	c.conn = conn
	return nil
}

// Close closes the WebSocket connection
func (c *Commander) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// ReadMessages reads and displays messages from the WebSocket
func (c *Commander) ReadMessages() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			return
		}

		// Pretty print the message
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, message, "", "  "); err == nil {
			fmt.Printf("\n📨 Received message:\n%s\n\n", prettyJSON.String())
		} else {
			fmt.Printf("\n📨 Received: %s\n\n", string(message))
		}
	}
}

// SendCommand sends a command to the WebSocket
func (c *Commander) SendCommand(cmd models.Command) error {
	// Ensure required fields are set
	if cmd.ID == uuid.Nil {
		cmd.ID = uuid.New()
	}
	if cmd.SenderID == uuid.Nil {
		cmd.SenderID = c.user.ID
	}
	if cmd.Tags == "" {
		cmd.Tags = "test"
	}
	if cmd.Status == "" {
		cmd.Status = "pending"
	}
	if cmd.CreatedAt.IsZero() {
		cmd.CreatedAt = time.Now()
	}
	if cmd.UpdatedAt.IsZero() {
		cmd.UpdatedAt = time.Now()
	}

	// Marshal and send
	data, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %w", err)
	}

	fmt.Printf("📤 Sending command: %s\n", cmd.Instructions[0].Type)
	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to send command: %w", err)
	}

	return nil
}

// InteractiveMenu displays an interactive menu for sending commands
func (c *Commander) InteractiveMenu() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		c.PrintMenu()
		fmt.Print("Select command (1-15, or 'q' to quit): ")

		if !scanner.Scan() {
			return
		}

		choice := strings.TrimSpace(scanner.Text())
		if choice == "q" || choice == "quit" {
			fmt.Println("👋 Goodbye!")
			return
		}

		var cmd models.Command
		var err error

		switch choice {
		case "1":
			cmd = c.CreatePingCommand()
		case "2":
			cmd = c.CreateNotificationCommand()
		case "3":
			cmd = c.CreatePopupCommand()
		case "4":
			cmd = c.CreateDisplayTextCommand()
		case "5":
			cmd = c.CreateTimerCommand()
		case "6":
			cmd = c.CreateChoiceCommand()
		case "7":
			cmd = c.CreateOpenURLCommand()
		case "8":
			cmd = c.CreateKinkMessageCommand()
		case "9":
			cmd = c.CreateKinkOpenLinkCommand()
		case "10":
			cmd = c.CreateKinkTTSCommand()
		case "11":
			cmd = c.CreateMultiInstructionCommand()
		case "12":
			cmd, err = c.CreateCustomCommand(scanner)
		case "13":
			cmd, err = c.CreateTargetedCommand(scanner)
		case "14":
			c.TestSequence()
			continue
		case "15":
			c.ListConnectedClients()
			continue
		default:
			fmt.Println("❌ Invalid choice")
			continue
		}

		if err != nil {
			fmt.Printf("❌ Error creating command: %v\n", err)
			continue
		}

		if err := c.SendCommand(cmd); err != nil {
			fmt.Printf("❌ Error sending command: %v\n", err)
		} else {
			fmt.Println("✅ Command sent successfully!")
		}
		fmt.Println()
	}
}

// PrintMenu displays the command menu
func (c *Commander) PrintMenu() {
	fmt.Println("╔════════════════════════════════════════════════════╗")
	fmt.Println("║           Test Commander Menu                      ║")
	fmt.Println("╠════════════════════════════════════════════════════╣")
	fmt.Println("║ Standard Commands:                                 ║")
	fmt.Println("║  1. Ping                                           ║")
	fmt.Println("║  2. Notification                                   ║")
	fmt.Println("║  3. Popup Message                                  ║")
	fmt.Println("║  4. Display Text                                   ║")
	fmt.Println("║  5. Timer                                          ║")
	fmt.Println("║  6. Multiple Choice                                ║")
	fmt.Println("║  7. Open URL                                       ║")
	fmt.Println("╠════════════════════════════════════════════════════╣")
	fmt.Println("║ Kink Commands:                                     ║")
	fmt.Println("║  8. Kink Message                                   ║")
	fmt.Println("║  9. Kink Open Link                                 ║")
	fmt.Println("║ 10. Kink TTS (Text-to-Speech)                      ║")
	fmt.Println("╠════════════════════════════════════════════════════╣")
	fmt.Println("║ Advanced:                                          ║")
	fmt.Println("║ 11. Multi-Instruction Command                      ║")
	fmt.Println("║ 12. Custom JSON Command                            ║")
	fmt.Println("║ 13. Targeted Command (specific user)               ║")
	fmt.Println("║ 14. Test Sequence (send multiple)                  ║")
	fmt.Println("║ 15. List Connected Clients                         ║")
	fmt.Println("╠════════════════════════════════════════════════════╣")
	fmt.Println("║  q. Quit                                           ║")
	fmt.Println("╚════════════════════════════════════════════════════╝")
	fmt.Println()
}

// Command creation functions

func (c *Commander) CreatePingCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_ping",
				Content: map[string]interface{}{
					"timestamp": time.Now().Format(time.RFC3339),
				},
			},
		},
		Tags: "test",
	}
}

func (c *Commander) CreateNotificationCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_notification",
				Content: map[string]interface{}{
					"title":    "Test Notification",
					"body":     "This is a test notification from the commander",
					"duration": 5,
				},
			},
		},
		Tags: "test",
	}
}

func (c *Commander) CreatePopupCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_popup",
				Content: map[string]interface{}{
					"body":    "This is a test popup message. Please acknowledge.",
					"title":   "Test Popup",
					"timeout": 30,
				},
			},
		},
		Tags: "test",
	}
}

func (c *Commander) CreateDisplayTextCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_display_text",
				Content: map[string]interface{}{
					"text":     "Welcome to the ControlApp test!\n\nThis is a display text command.",
					"title":    "Test Display",
					"closable": true,
				},
			},
		},
		Tags: "test",
	}
}

func (c *Commander) CreateTimerCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_timer",
				Content: map[string]interface{}{
					"duration":      60,
					"title":         "Test Timer",
					"show_progress": true,
				},
			},
		},
		Tags: "test",
	}
}

func (c *Commander) CreateChoiceCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_choice",
				Content: map[string]interface{}{
					"question": "How is the client working?",
					"options": []map[string]interface{}{
						{"id": "great", "text": "Great! Everything works"},
						{"id": "good", "text": "Good, minor issues"},
						{"id": "bad", "text": "Bad, not working"},
					},
					"timeout": 60,
				},
			},
		},
		Tags: "test",
	}
}

func (c *Commander) CreateOpenURLCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_open_url",
				Content: map[string]interface{}{
					"url":     "https://github.com/TheControlApp/server",
					"confirm": true,
				},
			},
		},
		Tags: "test",
	}
}

func (c *Commander) CreateKinkMessageCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "kink_message",
				Content: map[string]interface{}{
					"message":  "Your commander wants your attention",
					"title":    "Control Message",
					"style":    "info",
					"duration": 5000,
				},
			},
		},
		Tags: "kink",
	}
}

func (c *Commander) CreateKinkOpenLinkCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "kink_open_link",
				Content: map[string]interface{}{
					"url":     "https://www.example.com",
					"confirm": false,
				},
			},
		},
		Tags: "kink",
	}
}

func (c *Commander) CreateKinkTTSCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "kink_tts",
				Content: map[string]interface{}{
					"text":   "Hello, this is a test of the text to speech system",
					"voice":  "default",
					"volume": 0.8,
				},
			},
		},
		Tags: "kink",
	}
}

func (c *Commander) CreateMultiInstructionCommand() models.Command {
	return models.Command{
		Instructions: []models.Instruction{
			{
				Type: "std_notification",
				Content: map[string]interface{}{
					"title":    "Multi-Instruction Test",
					"body":     "Starting multi-step command sequence",
					"duration": 3,
				},
			},
			{
				Type: "std_display_text",
				Content: map[string]interface{}{
					"text":     "Step 1: This is the first instruction",
					"title":    "Multi-Step Command",
					"closable": true,
				},
			},
			{
				Type: "std_timer",
				Content: map[string]interface{}{
					"duration":      10,
					"title":         "Step 2: Short timer",
					"show_progress": true,
				},
			},
		},
		Tags: "test,multi",
	}
}

func (c *Commander) CreateCustomCommand(scanner *bufio.Scanner) (models.Command, error) {
	fmt.Println("\n📝 Enter custom command JSON (instruction type and content):")
	fmt.Println("Example: {\"type\": \"std_ping\", \"content\": {\"timestamp\": \"2024-01-22T10:00:00Z\"}}")
	fmt.Print("> ")

	if !scanner.Scan() {
		return models.Command{}, fmt.Errorf("failed to read input")
	}

	var instruction models.Instruction
	if err := json.Unmarshal([]byte(scanner.Text()), &instruction); err != nil {
		return models.Command{}, fmt.Errorf("invalid JSON: %w", err)
	}

	return models.Command{
		Instructions: []models.Instruction{instruction},
		Tags:         "custom",
	}, nil
}

func (c *Commander) CreateTargetedCommand(scanner *bufio.Scanner) (models.Command, error) {
	fmt.Print("\n🎯 Enter target user ID: ")

	if !scanner.Scan() {
		return models.Command{}, fmt.Errorf("failed to read input")
	}

	userIDStr := strings.TrimSpace(scanner.Text())
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return models.Command{}, fmt.Errorf("invalid user ID: %w", err)
	}

	cmd := c.CreatePopupCommand()
	cmd.ReceiverID = &userID
	cmd.Instructions[0].Content = map[string]interface{}{
		"body":    "This is a targeted message just for you!",
		"title":   "Private Message",
		"timeout": 30,
	}

	return cmd, nil
}

func (c *Commander) TestSequence() {
	fmt.Println("\n🔄 Running test sequence...")

	commands := []struct {
		name string
		cmd  models.Command
	}{
		{"Ping", c.CreatePingCommand()},
		{"Notification", c.CreateNotificationCommand()},
		{"Display Text", c.CreateDisplayTextCommand()},
	}

	for i, cmdInfo := range commands {
		fmt.Printf("  [%d/%d] Sending %s...\n", i+1, len(commands), cmdInfo.name)
		if err := c.SendCommand(cmdInfo.cmd); err != nil {
			fmt.Printf("  ❌ Failed: %v\n", err)
		} else {
			fmt.Printf("  ✅ Sent\n")
		}
		time.Sleep(2 * time.Second)
	}

	fmt.Println("✅ Test sequence complete!")
}

func (c *Commander) ListConnectedClients() {
	fmt.Println("\n📋 Connected Clients:")
	fmt.Println("(This would require server API endpoint to list connected users)")
	fmt.Println("For now, check server logs for connected clients")
}
