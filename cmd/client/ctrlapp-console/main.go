package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/thecontrolapp/server/internal/client"
)

func main() {
	fmt.Println("ControlApp Console Client")
	fmt.Println("========================")

	// Create client with default config
	config := client.DefaultConfig()
	config.ServerURL = "ws://localhost:3080/ws/client"

	// Create console logger and set in config
	config.Logger = client.NewDefaultLogger()

	// Create client
	c := client.NewClient(config)

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		cancel()
	}()

	// Start interactive session
	if err := runInteractiveSession(ctx, c); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func runInteractiveSession(ctx context.Context, c *client.Client) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			fmt.Print("> ")
			if !scanner.Scan() {
				return nil
			}

			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}

			if err := handleCommand(c, line); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	}
}

func handleCommand(c *client.Client, input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}

	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "help":
		showHelp()
	case "connect":
		return handleConnect(c, args)
	case "login":
		return handleLogin(c, args)
	case "register":
		return handleRegister(c, args)
	case "status":
		return handleStatus(c)
	case "ping":
		return handlePing(c)
	case "message":
		return handleMessage(c, args)
	case "notify":
		return handleNotify(c, args)
	case "disconnect":
		return c.Disconnect()
	case "quit", "exit":
		return c.Disconnect()
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", cmd)
	}

	return nil
}

func showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                     - Show this help")
	fmt.Println("  connect [url]           - Connect to server (default: ws://localhost:3080/ws/client)")
	fmt.Println("  login <username> <pass> - Login to server")
	fmt.Println("  register <screen> <login> <pass> - Register new user")
	fmt.Println("  status                  - Show connection status")
	fmt.Println("  ping                    - Send std_ping command")
	fmt.Println("  message <text>          - Send std_message command")
	fmt.Println("  notify <title> <text>   - Send std_notification command")
	fmt.Println("  disconnect              - Disconnect from server")
	fmt.Println("  quit/exit              - Exit client")
}

func handleConnect(c *client.Client, args []string) error {
	url := "ws://localhost:3080/ws/client"
	if len(args) > 0 {
		url = args[0]
	}

	fmt.Printf("Connecting to %s...\n", url)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := c.Connect(ctx, url); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}

	fmt.Println("Connected successfully!")
	return nil
}

func handleLogin(c *client.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: login <username> <password>")
	}

	username := args[0]
	password := args[1]

	fmt.Printf("Logging in as %s...\n", username)

	err := c.Login(username, password)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	user := c.GetUser()
	if user != nil {
		fmt.Printf("Logged in successfully! Welcome, %s (ID: %d)\n", user.ScreenName, user.ID)
	} else {
		fmt.Println("Logged in successfully!")
	}
	return nil
}

func handleRegister(c *client.Client, args []string) error {
	if len(args) < 3 {
		return fmt.Errorf("usage: register <screen_name> <login_name> <password>")
	}

	screenName := args[0]
	loginName := args[1]
	password := args[2]

	fmt.Printf("Registering user %s (%s)...\n", screenName, loginName)

	err := c.Register(screenName, loginName, password)
	if err != nil {
		return fmt.Errorf("registration failed: %w", err)
	}

	fmt.Printf("Registered successfully! User: %s\n", loginName)
	return nil
}

func handleStatus(c *client.Client) error {
	fmt.Println("Client Status:")
	fmt.Printf("  Connected: %v\n", c.IsConnected())
	fmt.Printf("  Authenticated: %v\n", c.IsAuthenticated())
	fmt.Printf("  Connection State: %v\n", c.GetConnectionState())

	user := c.GetUser()
	if user != nil {
		fmt.Printf("  User: %s (ID: %d)\n", user.ScreenName, user.ID)
	}

	return nil
}

func handlePing(c *client.Client) error {
	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
	}

	fmt.Println("Sending ping...")

	cmd := client.Command{
		ID:   fmt.Sprintf("ping-%d", time.Now().UnixNano()),
		Type: "std_ping",
		Content: map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339Nano),
		},
		ReceivedAt: time.Now(),
	}

	err := c.SendCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to send ping: %w", err)
	}

	fmt.Println("Ping sent successfully!")
	return nil
}

func handleMessage(c *client.Client, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: message <text>")
	}

	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
	}

	message := strings.Join(args, " ")

	fmt.Printf("Sending message: %s\n", message)

	cmd := client.Command{
		ID:   fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		Type: "std_message",
		Content: map[string]interface{}{
			"message": message,
			"title":   "Console Test",
			"style":   "info",
		},
		ReceivedAt: time.Now(),
	}

	err := c.SendCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Println("Message sent successfully!")
	return nil
}

func handleNotify(c *client.Client, args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: notify <title> <text>")
	}

	if !c.IsConnected() {
		return fmt.Errorf("not connected to server")
	}

	title := args[0]
	body := strings.Join(args[1:], " ")

	fmt.Printf("Sending notification: %s - %s\n", title, body)

	cmd := client.Command{
		ID:   fmt.Sprintf("notify-%d", time.Now().UnixNano()),
		Type: "std_notification",
		Content: map[string]interface{}{
			"title":    title,
			"body":     body,
			"icon":     "info",
			"duration": 5,
		},
		ReceivedAt: time.Now(),
	}

	err := c.SendCommand(cmd)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	fmt.Println("Notification sent successfully!")
	return nil
}
