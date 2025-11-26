package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/thecontrolapp/server/internal/client"
)

type ControlApp struct {
	app        fyne.App
	mainWindow fyne.Window
	client     *client.Client
	credStore  *client.CredentialStore

	// UI state
	isLoggedIn  bool
	isMinimized bool

	// UI components  
	serverEntry     *widget.Entry
	screenNameEntry *widget.Entry
	usernameEntry   *widget.Entry
	passwordEntry   *widget.Entry
	statusLabel     *widget.Label
	logArea         *widget.Entry	// System tray (Windows)
	systray desktop.App
}

func main() {
	myApp := app.NewWithID("com.controlapp.client")

	clientApp := &ControlApp{
		app:       myApp,
		credStore: client.NewCredentialStore(),
	}

	// Initialize client
	if err := clientApp.initClient(); err != nil {
		log.Fatal("Failed to initialize client:", err)
	}

	// Try auto-login first
	if clientApp.tryAutoLogin() {
		// Auto-login successful, start minimized
		clientApp.startBackgroundMode()
	} else {
		// No valid credentials, show login window
		clientApp.showLoginWindow()
	}

	myApp.Run()
}

func (ca *ControlApp) initClient() error {
	config := client.DefaultConfig()
	config.ServerURL = "ws://localhost:3080/ws/client"

	// Set up supported commands for GUI client
	config.AllowedCommands = []string{
		// Core commands (always allowed)
		"std_ping", "std_message", "std_notification",

		// Standard commands (user configurable)
		"std_open_url", "std_download_file", "std_play_audio",
		"std_display_image", "std_timer",
	}

	// Commands requiring explicit consent (disabled by default)
	config.BlockedCommands = []string{
		"std_change_wallpaper", "std_lock_screen", "std_execute_file",
	}

	config.Logger = client.NewDefaultLogger()

	ca.client = client.NewClient(config)

	return nil
}

func (ca *ControlApp) tryAutoLogin() bool {
	creds, err := ca.credStore.Load()
	if err != nil || !ca.credStore.IsValid(creds) {
		return false
	}

	// Try to connect and authenticate with stored credentials
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := ca.client.Connect(ctx, creds.ServerURL); err != nil {
		ca.logf("Auto-login failed: connection error: %v", err)
		return false
	}

	if err := ca.client.Login(creds.Username, creds.Token); err != nil {
		ca.logf("Auto-login failed: authentication error: %v", err)
		return false
	}

	ca.isLoggedIn = true
	ca.logf("Auto-login successful for user: %s", creds.Username)
	return true
}

func (ca *ControlApp) showLoginWindow() {
	ca.mainWindow = ca.app.NewWindow("ControlApp - Login")
	ca.mainWindow.Resize(fyne.NewSize(500, 400))
	ca.mainWindow.CenterOnScreen()

	// Create login form
	ca.serverEntry = widget.NewEntry()
	ca.serverEntry.SetText("ws://localhost:3080/ws/client")

	ca.screenNameEntry = widget.NewEntry()
	ca.screenNameEntry.SetPlaceHolder("Display Name (only needed for registration)")

	ca.usernameEntry = widget.NewEntry()
	ca.usernameEntry.SetPlaceHolder("Login Username")

	ca.passwordEntry = widget.NewPasswordEntry()
	ca.passwordEntry.SetPlaceHolder("Password")

	loginBtn := widget.NewButton("Login", ca.handleLogin)
	registerBtn := widget.NewButton("Register", ca.handleRegister)

	ca.statusLabel = widget.NewLabel("Enter your credentials to connect")

	// Consent notice
	consentText := widget.NewLabel("CONSENT NOTICE:\n\nBy running this application, you consent to receive and execute commands from authorized users on the connected server.\n\nYou may revoke consent at any time by closing this application.\n\nCommands may include: messages, notifications, file downloads, and other actions based on your configuration.")
	consentText.Wrapping = fyne.TextWrapWord

	// Help text
	helpText := widget.NewLabel("• For login: Enter username and password\n• For registration: Enter display name, username, and password")
	helpText.Wrapping = fyne.TextWrapWord

	// Layout
	form := container.NewVBox(
		widget.NewLabel("Server Connection"),
		ca.serverEntry,
		widget.NewSeparator(),
		widget.NewLabel("Authentication"),
		helpText,
		ca.screenNameEntry,
		ca.usernameEntry,
		ca.passwordEntry,
		container.NewHBox(loginBtn, registerBtn),
		widget.NewSeparator(),
		ca.statusLabel,
		widget.NewSeparator(),
		consentText,
	)

	scroll := container.NewScroll(form)
	ca.mainWindow.SetContent(scroll)

	// Handle window close
	ca.mainWindow.SetCloseIntercept(func() {
		ca.app.Quit()
	})

	ca.mainWindow.Show()
}

func (ca *ControlApp) handleLogin() {
	server := ca.serverEntry.Text
	username := ca.usernameEntry.Text
	password := ca.passwordEntry.Text
	
	// Validation
	if server == "" {
		ca.statusLabel.SetText("Please enter server URL")
		return
	}
	if username == "" {
		ca.statusLabel.SetText("Please enter username")
		return
	}
	if password == "" {
		ca.statusLabel.SetText("Please enter password")
		return
	}

	ca.statusLabel.SetText("Connecting...")

	// Connect and login in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Connect
		if err := ca.client.Connect(ctx, server); err != nil {
			fyne.Do(func() {
				ca.statusLabel.SetText(fmt.Sprintf("Connection failed: %v", err))
			})
			return
		}

		// Login
		if err := ca.client.Login(username, password); err != nil {
			fyne.Do(func() {
				ca.statusLabel.SetText(fmt.Sprintf("Login failed: %v", err))
			})
			return
		}

		// Save credentials
		user := ca.client.GetUser()
		if user != nil {
			creds := &client.Credentials{
				Token:     "temp_token", // TODO: Get actual JWT from auth response
				Username:  username,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 1 week
				ServerURL: server,
			}
			ca.credStore.Save(creds)
		}

		ca.isLoggedIn = true
		fyne.Do(func() {
			ca.statusLabel.SetText("Login successful! Starting background mode...")
		})

		// Wait a moment then switch to background mode
		time.Sleep(2 * time.Second)
		ca.startBackgroundMode()
	}()
}

func (ca *ControlApp) handleRegister() {
	server := ca.serverEntry.Text
	screenName := ca.screenNameEntry.Text
	username := ca.usernameEntry.Text
	password := ca.passwordEntry.Text
	
	// Validation
	if server == "" {
		ca.statusLabel.SetText("Please enter server URL")
		return
	}
	if screenName == "" {
		ca.statusLabel.SetText("Please enter display name")
		return
	}
	if username == "" {
		ca.statusLabel.SetText("Please enter username")
		return
	}
	if password == "" {
		ca.statusLabel.SetText("Please enter password")
		return
	}

	ca.statusLabel.SetText("Registering...")

	// Connect and register in background
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Connect to server first
		if err := ca.client.Connect(ctx, server); err != nil {
			fyne.Do(func() {
				ca.statusLabel.SetText(fmt.Sprintf("Connection failed: %v", err))
			})
			return
		}

		// Register new user
		if err := ca.client.Register(screenName, username, password); err != nil {
			fyne.Do(func() {
				ca.statusLabel.SetText(fmt.Sprintf("Registration failed: %v", err))
			})
			return
		}

		fyne.Do(func() {
			ca.statusLabel.SetText("Registration successful! Now logging in...")
		})

		// After successful registration, automatically login
		if err := ca.client.Login(username, password); err != nil {
			fyne.Do(func() {
				ca.statusLabel.SetText(fmt.Sprintf("Login after registration failed: %v", err))
			})
			return
		}

		// Save credentials for future auto-login
		user := ca.client.GetUser()
		if user != nil {
			creds := &client.Credentials{
				Token:     "temp_token", // TODO: Get actual JWT from auth response
				Username:  username,
				ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 1 week
				ServerURL: server,
			}
			ca.credStore.Save(creds)
		}

		ca.isLoggedIn = true
		fyne.Do(func() {
			ca.statusLabel.SetText("Registration and login successful! Starting background mode...")
		})

		// Wait a moment then switch to background mode
		time.Sleep(2 * time.Second)
		ca.startBackgroundMode()
	}()
}

func (ca *ControlApp) startBackgroundMode() {
	// Close login window if it exists
	if ca.mainWindow != nil {
		ca.mainWindow.Close()
		ca.mainWindow = nil
	}

	ca.isMinimized = true

	// Set up system tray (Windows-specific)
	ca.setupSystemTray()

	// Start command processing
	ca.startEventLoop()

	ca.logf("ControlApp client running in background mode")
}

func (ca *ControlApp) setupSystemTray() {
	if desk, ok := ca.app.(desktop.App); ok {
		ca.systray = desk

		// Create system tray menu
		menu := fyne.NewMenu("ControlApp",
			fyne.NewMenuItem("Show Command Center", ca.showCommandCenter),
			fyne.NewMenuItem("Settings", ca.showSettings),
			fyne.NewMenuItemSeparator(),
			fyne.NewMenuItem("Disconnect", ca.handleDisconnect),
			fyne.NewMenuItem("Exit", func() { ca.app.Quit() }),
		)

		desk.SetSystemTrayMenu(menu)
		ca.logf("System tray initialized")
	}
}

func (ca *ControlApp) showCommandCenter() {
	if ca.mainWindow != nil {
		ca.mainWindow.RequestFocus()
		return
	}

	ca.mainWindow = ca.app.NewWindow("ControlApp - Command Center")
	ca.mainWindow.Resize(fyne.NewSize(800, 600))

	// Status display
	statusInfo := widget.NewLabel(ca.getStatusText())

	// Quick command buttons
	pingBtn := widget.NewButton("Send Ping", func() { ca.sendTestCommand("std_ping") })
	messageBtn := widget.NewButton("Send Test Message", func() { ca.sendTestCommand("std_message") })
	notifyBtn := widget.NewButton("Send Test Notification", func() { ca.sendTestCommand("std_notification") })

	commandBox := container.NewHBox(pingBtn, messageBtn, notifyBtn)

	// Log area
	ca.logArea = widget.NewMultiLineEntry()
	ca.logArea.SetText("=== ControlApp Command Center ===\nClient running in background mode.\n")
	ca.logArea.Wrapping = fyne.TextWrapWord
	logScroll := container.NewScroll(ca.logArea)
	logScroll.SetMinSize(fyne.NewSize(0, 300))

	content := container.NewVBox(
		widget.NewLabel("ControlApp Command Center"),
		widget.NewSeparator(),
		statusInfo,
		widget.NewSeparator(),
		widget.NewLabel("Quick Commands:"),
		commandBox,
		widget.NewSeparator(),
		widget.NewLabel("Activity Log:"),
		logScroll,
	)

	ca.mainWindow.SetContent(content)

	// Handle window close (minimize to tray)
	ca.mainWindow.SetCloseIntercept(func() {
		ca.mainWindow.Hide()
		ca.mainWindow = nil
	})

	ca.mainWindow.Show()
}

func (ca *ControlApp) showSettings() {
	settingsWindow := ca.app.NewWindow("ControlApp - Settings")
	settingsWindow.Resize(fyne.NewSize(600, 500))

	// Consent settings
	consentLabel := widget.NewLabel("Command Consent Settings:")

	// TODO: Create checkboxes for each command type
	urlCheck := widget.NewCheck("Allow opening URLs", nil)
	downloadCheck := widget.NewCheck("Allow file downloads", nil)
	audioCheck := widget.NewCheck("Allow audio playback", nil)
	wallpaperCheck := widget.NewCheck("Allow wallpaper changes", nil)

	consentBox := container.NewVBox(
		consentLabel,
		urlCheck,
		downloadCheck,
		audioCheck,
		wallpaperCheck,
	)

	saveBtn := widget.NewButton("Save Settings", func() {
		// TODO: Save consent settings
		dialog.ShowInformation("Settings", "Settings saved successfully", settingsWindow)
	})

	content := container.NewVBox(
		widget.NewLabel("ControlApp Settings"),
		widget.NewSeparator(),
		consentBox,
		widget.NewSeparator(),
		saveBtn,
	)

	settingsWindow.SetContent(content)
	settingsWindow.Show()
}

func (ca *ControlApp) sendTestCommand(cmdType string) {
	if !ca.client.IsConnected() {
		ca.logf("Cannot send command: not connected")
		return
	}

	var cmd client.Command

	switch cmdType {
	case "std_ping":
		cmd = client.Command{
			ID:   fmt.Sprintf("test-ping-%d", time.Now().UnixNano()),
			Type: "std_ping",
			Content: map[string]interface{}{
				"timestamp": time.Now().Format(time.RFC3339Nano),
			},
			ReceivedAt: time.Now(),
		}
	case "std_message":
		cmd = client.Command{
			ID:   fmt.Sprintf("test-msg-%d", time.Now().UnixNano()),
			Type: "std_message",
			Content: map[string]interface{}{
				"message": "Test message from ControlApp client",
				"title":   "Test Message",
				"style":   "info",
			},
			ReceivedAt: time.Now(),
		}
	case "std_notification":
		cmd = client.Command{
			ID:   fmt.Sprintf("test-notify-%d", time.Now().UnixNano()),
			Type: "std_notification",
			Content: map[string]interface{}{
				"title": "Test Notification",
				"body":  "This is a test notification from ControlApp",
				"icon":  "info",
			},
			ReceivedAt: time.Now(),
		}
	default:
		ca.logf("Unknown test command type: %s", cmdType)
		return
	}

	ca.logf("Sending test command: %s", cmdType)

	go func() {
		if err := ca.client.SendCommand(cmd); err != nil {
			ca.logf("Failed to send %s: %v", cmdType, err)
		} else {
			ca.logf("Test command %s sent successfully", cmdType)
		}
	}()
}

func (ca *ControlApp) handleDisconnect() {
	if ca.client.IsConnected() {
		ca.client.Disconnect()
	}

	// Clear credentials
	ca.credStore.Clear()
	ca.isLoggedIn = false

	// Show login window again
	ca.showLoginWindow()
}

func (ca *ControlApp) startEventLoop() {
	go func() {
		for {
			select {
			case event := <-ca.client.Events():
				ca.handleEvent(event)
			case cmd := <-ca.client.Commands():
				ca.handleIncomingCommand(cmd)
			case err := <-ca.client.Errors():
				ca.logf("Client error: %v", err)
			}
		}
	}()
}

func (ca *ControlApp) handleEvent(event client.Event) {
	switch event.Type {
	case client.EventConnected:
		ca.logf("Connected to server")
	case client.EventDisconnected:
		ca.logf("Disconnected from server")
	case client.EventAuthenticated:
		ca.logf("Authentication successful")
	case client.EventCommandReceived:
		ca.logf("Command received")
	case client.EventCommandCompleted:
		ca.logf("Command completed successfully")
	case client.EventCommandFailed:
		if errorMsg, ok := event.Data["error"].(string); ok {
			ca.logf("Command failed: %s", errorMsg)
		}
	}
}

func (ca *ControlApp) handleIncomingCommand(cmd client.Command) {
	ca.logf("Received command: %s (ID: %s)", cmd.Type, cmd.ID)

	// For high-risk commands, show confirmation dialog
	if ca.isHighRiskCommand(cmd.Type) && ca.mainWindow == nil {
		// If no window is open, show one for confirmation
		ca.showCommandCenter()
	}

	// Commands are automatically processed by the client
	// This is just for logging and user awareness
}

func (ca *ControlApp) isHighRiskCommand(cmdType string) bool {
	highRisk := []string{
		"std_execute_file", "std_lock_screen", "std_change_wallpaper",
	}

	for _, risk := range highRisk {
		if cmdType == risk {
			return true
		}
	}
	return false
}

func (ca *ControlApp) getStatusText() string {
	if !ca.client.IsConnected() {
		return "Status: Disconnected"
	}

	if !ca.client.IsAuthenticated() {
		return "Status: Connected (Not authenticated)"
	}

	user := ca.client.GetUser()
	if user != nil {
		return fmt.Sprintf("Status: Connected as %s (ID: %d)", user.ScreenName, user.ID)
	}

	return "Status: Connected and authenticated"
}

func (ca *ControlApp) logf(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf("[%s] %s", timestamp, fmt.Sprintf(format, args...))

	// Log to console
	log.Println(message)

	// Log to GUI if available (must be done on UI thread)
	if ca.logArea != nil {
		fyne.Do(func() {
			ca.logArea.SetText(ca.logArea.Text + message + "\n")
		})
	}
}
