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
	"fyne.io/fyne/v2/widget"

	"github.com/thecontrolapp/server/internal/client"
)

type App struct {
	fyneApp  fyne.App
	window   fyne.Window
	client   *client.Client
	
	// UI Components
	serverEntry    *widget.Entry
	usernameEntry  *widget.Entry
	passwordEntry  *widget.Entry
	statusLabel    *widget.Label
	connectBtn     *widget.Button
	loginBtn       *widget.Button
	disconnectBtn  *widget.Button
	
	// Log display
	logArea        *widget.Entry
	
	// Connection state
	connected      bool
	authenticated  bool
}

func main() {
	myApp := app.New()
	myApp.Metadata().Name = "ControlApp Windows Client"
	myApp.Metadata().Version = "1.0.0"

	window := myApp.NewWindow("ControlApp Client")
	window.Resize(fyne.NewSize(800, 600))

	clientApp := &App{
		fyneApp: myApp,
		window:  window,
	}

	clientApp.initClient()
	clientApp.setupUI()
	clientApp.startEventLoop()

	window.ShowAndRun()
}

func (a *App) initClient() {
	config := client.DefaultConfig()
	config.ServerURL = "ws://localhost:3080/ws"
	
	// Enable Windows-specific kink commands
	config.AllowedCommands = append(config.AllowedCommands, 
		"kink_open_link",
		"kink_download_file", 
		"kink_play_audio",
		"kink_tts",
		"kink_popup_image",
		"kink_change_wallpaper",
	)
	
	config.Logger = client.NewDefaultLogger()
	
	a.client = client.NewClient(config)
	
	// Register Windows handlers for kink commands
	a.client.RegisterWindowsHandlers()
}

func (a *App) setupUI() {
	// Connection section
	a.serverEntry = widget.NewEntry()
	a.serverEntry.SetText("ws://localhost:3080/ws")
	a.serverEntry.SetPlaceHolder("Server URL")

	a.connectBtn = widget.NewButton("Connect", a.handleConnect)
	a.disconnectBtn = widget.NewButton("Disconnect", a.handleDisconnect)
	a.disconnectBtn.Disable()

	connectionBox := container.NewHBox(
		widget.NewLabel("Server:"),
		a.serverEntry,
		a.connectBtn,
		a.disconnectBtn,
	)

	// Authentication section
	a.usernameEntry = widget.NewEntry()
	a.usernameEntry.SetPlaceHolder("Username")
	
	a.passwordEntry = widget.NewPasswordEntry()
	a.passwordEntry.SetPlaceHolder("Password")

	a.loginBtn = widget.NewButton("Login", a.handleLogin)
	a.loginBtn.Disable()

	registerBtn := widget.NewButton("Register", a.handleRegister)
	registerBtn.Disable()

	authBox := container.NewHBox(
		widget.NewLabel("Auth:"),
		a.usernameEntry,
		a.passwordEntry,
		a.loginBtn,
		registerBtn,
	)

	// Status section
	a.statusLabel = widget.NewLabel("Disconnected")
	statusBox := container.NewHBox(
		widget.NewLabel("Status:"),
		a.statusLabel,
	)

	// Kink command testing section
	testBox := a.createTestCommandsBox()

	// Log area
	a.logArea = widget.NewMultiLineEntry()
	a.logArea.SetText("=== ControlApp Windows Client Log ===\n")
	a.logArea.Wrapping = fyne.TextWrapWord
	logScroll := container.NewScroll(a.logArea)
	logScroll.SetMinSize(fyne.NewSize(0, 200))

	// Main layout
	content := container.NewVBox(
		connectionBox,
		authBox,  
		statusBox,
		widget.NewSeparator(),
		widget.NewLabel("Kink Command Testing:"),
		testBox,
		widget.NewSeparator(),
		widget.NewLabel("Log:"),
		logScroll,
	)

	a.window.SetContent(container.NewScroll(content))
}

func (a *App) createTestCommandsBox() *fyne.Container {
	// Message test
	messageEntry := widget.NewEntry()
	messageEntry.SetPlaceHolder("Enter message text...")
	messageBtn := widget.NewButton("Send Message", func() {
		a.sendKinkMessage(messageEntry.Text)
	})

	// Link test
	linkEntry := widget.NewEntry()
	linkEntry.SetText("https://example.com")
	linkBtn := widget.NewButton("Open Link", func() {
		a.sendKinkOpenLink(linkEntry.Text)
	})

	// TTS test
	ttsEntry := widget.NewEntry()
	ttsEntry.SetPlaceHolder("Text to speak...")
	ttsBtn := widget.NewButton("Speak Text", func() {
		a.sendKinkTTS(ttsEntry.Text)
	})

	// Download test
	downloadEntry := widget.NewEntry()
	downloadEntry.SetText("https://httpbin.org/json")
	downloadBtn := widget.NewButton("Download File", func() {
		a.sendKinkDownload(downloadEntry.Text)
	})

	return container.NewVBox(
		container.NewHBox(messageEntry, messageBtn),
		container.NewHBox(linkEntry, linkBtn),
		container.NewHBox(ttsEntry, ttsBtn),
		container.NewHBox(downloadEntry, downloadBtn),
	)
}

func (a *App) handleConnect() {
	url := a.serverEntry.Text
	if url == "" {
		url = "ws://localhost:3080/ws"
	}

	a.logf("Connecting to %s...", url)
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	go func() {
		err := a.client.Connect(ctx, url)
		if err != nil {
			a.logf("Connection failed: %v", err)
			return
		}
		
		a.connected = true
		a.updateUI()
		a.logf("Connected successfully!")
	}()
}

func (a *App) handleDisconnect() {
	a.logf("Disconnecting...")
	
	go func() {
		err := a.client.Disconnect()
		if err != nil {
			a.logf("Disconnect error: %v", err)
		}
		
		a.connected = false
		a.authenticated = false
		a.updateUI()
		a.logf("Disconnected")
	}()
}

func (a *App) handleLogin() {
	username := a.usernameEntry.Text
	password := a.passwordEntry.Text
	
	if username == "" || password == "" {
		a.logf("Please enter username and password")
		return
	}

	a.logf("Logging in as %s...", username)
	
	go func() {
		err := a.client.Login(username, password)
		if err != nil {
			a.logf("Login failed: %v", err)
			return
		}
		
		a.authenticated = true
		a.updateUI()
		
		user := a.client.GetUser()
		if user != nil {
			a.logf("Logged in successfully! Welcome, %s (ID: %d)", user.ScreenName, user.ID)
		} else {
			a.logf("Logged in successfully!")
		}
	}()
}

func (a *App) handleRegister() {
	username := a.usernameEntry.Text
	password := a.passwordEntry.Text
	
	if username == "" || password == "" {
		a.logf("Please enter username and password")
		return
	}

	screenName := username // Use username as screen name for simplicity
	
	a.logf("Registering user %s...", username)
	
	go func() {
		err := a.client.Register(screenName, username, password)
		if err != nil {
			a.logf("Registration failed: %v", err)
			return
		}
		
		a.logf("Registered successfully! User: %s", username)
	}()
}

func (a *App) sendKinkMessage(message string) {
	if !a.authenticated {
		a.logf("Must be connected and authenticated to send commands")
		return
	}
	
	if message == "" {
		a.logf("Please enter a message")
		return
	}

	cmd := client.Command{
		ID:   fmt.Sprintf("msg-%d", time.Now().UnixNano()),
		Type: "kink_message",
		Content: map[string]interface{}{
			"message": message,
			"title":   "ControlApp GUI Test",
			"style":   "info",
		},
		ReceivedAt: time.Now(),
	}

	a.logf("Sending kink_message: %s", message)
	
	go func() {
		err := a.client.SendCommand(cmd)
		if err != nil {
			a.logf("Failed to send message: %v", err)
		} else {
			a.logf("Message sent successfully!")
		}
	}()
}

func (a *App) sendKinkOpenLink(url string) {
	if !a.authenticated {
		a.logf("Must be connected and authenticated to send commands")
		return
	}
	
	if url == "" {
		a.logf("Please enter a URL")
		return
	}

	cmd := client.Command{
		ID:   fmt.Sprintf("link-%d", time.Now().UnixNano()),
		Type: "kink_open_link", 
		Content: map[string]interface{}{
			"url": url,
		},
		ReceivedAt: time.Now(),
	}

	a.logf("Sending kink_open_link: %s", url)
	
	go func() {
		err := a.client.SendCommand(cmd)
		if err != nil {
			a.logf("Failed to send open link: %v", err)
		} else {
			a.logf("Open link sent successfully!")
		}
	}()
}

func (a *App) sendKinkTTS(text string) {
	if !a.authenticated {
		a.logf("Must be connected and authenticated to send commands")
		return
	}
	
	if text == "" {
		a.logf("Please enter text to speak")
		return
	}

	cmd := client.Command{
		ID:   fmt.Sprintf("tts-%d", time.Now().UnixNano()),
		Type: "kink_tts",
		Content: map[string]interface{}{
			"text":   text,
			"voice":  "default",
			"rate":   0,
			"volume": 80,
		},
		ReceivedAt: time.Now(),
	}

	a.logf("Sending kink_tts: %s", text)
	
	go func() {
		err := a.client.SendCommand(cmd)
		if err != nil {
			a.logf("Failed to send TTS: %v", err)
		} else {
			a.logf("TTS sent successfully!")
		}
	}()
}

func (a *App) sendKinkDownload(url string) {
	if !a.authenticated {
		a.logf("Must be connected and authenticated to send commands")
		return
	}
	
	if url == "" {
		a.logf("Please enter a URL to download")
		return
	}

	cmd := client.Command{
		ID:   fmt.Sprintf("dl-%d", time.Now().UnixNano()),
		Type: "kink_download_file",
		Content: map[string]interface{}{
			"url":      url,
			"filename": "", // Auto-generate
		},
		ReceivedAt: time.Now(),
	}

	a.logf("Sending kink_download_file: %s", url)
	
	go func() {
		err := a.client.SendCommand(cmd)
		if err != nil {
			a.logf("Failed to send download: %v", err)
		} else {
			a.logf("Download sent successfully!")
		}
	}()
}

func (a *App) updateUI() {
	if a.connected {
		a.statusLabel.SetText("Connected")
		a.connectBtn.Disable()
		a.disconnectBtn.Enable()
		a.loginBtn.Enable()
	} else {
		a.statusLabel.SetText("Disconnected")
		a.connectBtn.Enable()
		a.disconnectBtn.Disable()
		a.loginBtn.Disable()
	}

	if a.authenticated {
		a.statusLabel.SetText("Connected & Authenticated")
		a.loginBtn.Disable()
	}
}

func (a *App) startEventLoop() {
	go func() {
		for {
			select {
			case event := <-a.client.Events():
				a.handleEvent(event)
			case cmd := <-a.client.Commands():
				a.handleIncomingCommand(cmd)
			case err := <-a.client.Errors():
				a.logf("Client error: %v", err)
			}
		}
	}()
}

func (a *App) handleEvent(event client.Event) {
	switch event.Type {
	case client.EventConnected:
		a.logf("Event: Connected to server")
	case client.EventDisconnected:
		a.logf("Event: Disconnected from server")
	case client.EventAuthenticated:
		a.logf("Event: Authenticated successfully")
	case client.EventCommandReceived:
		a.logf("Event: Command received")
	case client.EventCommandCompleted:
		a.logf("Event: Command completed successfully")
	case client.EventCommandFailed:
		if errorMsg, ok := event.Data["error"].(string); ok {
			a.logf("Event: Command failed - %s", errorMsg)
		} else {
			a.logf("Event: Command failed")
		}
	}
}

func (a *App) handleIncomingCommand(cmd client.Command) {
	a.logf("Received command: %s (ID: %s)", cmd.Type, cmd.ID)
	
	// Commands are automatically processed by the client
	// This is just for logging/notification purposes
	
	// For certain commands, we might want to show user confirmation
	if cmd.Type == "kink_run_file" || cmd.Type == "kink_lock_screen" {
		dialog.ShowConfirm(
			"Command Confirmation",
			fmt.Sprintf("Allow command: %s?", cmd.Type),
			func(confirmed bool) {
				if confirmed {
					a.logf("User confirmed command: %s", cmd.Type)
				} else {
					a.logf("User denied command: %s", cmd.Type)
					// TODO: Send command rejection back to server
				}
			},
			a.window,
		)
	}
}

func (a *App) logf(format string, args ...interface{}) {
	timestamp := time.Now().Format("15:04:05")
	message := fmt.Sprintf("[%s] %s\n", timestamp, fmt.Sprintf(format, args...))
	
	// Update log area (must be done on UI thread)
	a.logArea.SetText(a.logArea.Text + message)
	
	// Also log to console
	log.Printf(format, args...)
}