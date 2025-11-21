# ControlApp Client Architecture - Go Shared Core

## 🏗️ Project Structure

```
client/                                 # New client directory (sibling to server)
├── go.mod                             # Client Go module
├── go.sum
├── README.md                          # Client-specific documentation
├── 
├── internal/                          # Shared client core (95% reuse)
│   └── client/                        # Core client package
│       ├── client.go                  # Main client struct and interface
│       ├── auth.go                    # Authentication (login, register, JWT)
│       ├── websocket.go               # WebSocket connection management
│       ├── commands.go                # Command processing and handlers
│       ├── types.go                   # Shared types and structures
│       ├── config.go                  # Configuration management
│       ├── errors.go                  # Client-specific error types
│       └── utils.go                   # Shared utilities
│
├── cmd/                               # Platform-specific entry points
│   ├── windows/                       # Windows desktop client
│   │   ├── main.go                    # Windows entry point
│   │   ├── ui/                        # Windows-specific UI
│   │   └── platform/                  # Windows-specific features
│   ├── macos/                         # macOS desktop client
│   │   ├── main.go                    # macOS entry point
│   │   ├── ui/                        # macOS-specific UI
│   │   └── platform/                  # macOS-specific features
│   ├── linux/                         # Linux desktop client
│   │   ├── main.go                    # Linux entry point
│   │   ├── ui/                        # Linux-specific UI
│   │   └── platform/                  # Linux-specific features
│   └── android/                       # Android mobile client
│       ├── main.go                    # Android entry point
│       ├── ui/                        # Mobile-specific UI
│       └── platform/                  # Android-specific features
│
├── pkg/                               # Public client libraries
│   ├── fyne/                          # Fyne UI components (shared across desktop)
│   │   ├── components/                # Reusable Fyne widgets
│   │   ├── themes/                    # Custom themes
│   │   └── layouts/                   # Custom layouts
│   └── mobile/                        # Mobile-specific utilities
│
├── web/                               # Web client (TypeScript)
│   ├── package.json
│   ├── src/
│   └── dist/
│
├── extension/                         # Chrome extension (TypeScript)
│   ├── manifest.json
│   ├── src/
│   └── dist/
│
├── build/                             # Build scripts and configurations
│   ├── windows.sh                     # Windows build script
│   ├── macos.sh                       # macOS build script
│   ├── linux.sh                       # Linux build script
│   ├── android.sh                     # Android build script
│   ├── docker/                        # Docker build environments
│   └── ci/                            # CI/CD configurations
│
└── assets/                            # Shared assets
    ├── icons/                         # Application icons
    ├── images/                        # Images and graphics
    └── fonts/                         # Custom fonts
```

## 🔧 Core Client Package (`internal/client`)

### Main Client Interface
```go
// internal/client/client.go
package client

import (
    "context"
    "time"
)

// Client represents the core ControlApp client
type Client struct {
    config     *Config
    auth       *AuthManager
    ws         *WebSocketManager
    commands   *CommandProcessor
    logger     Logger
    
    // Client state
    connected    bool
    authenticated bool
    user         *User
    
    // Event channels
    eventChan    chan Event
    commandChan  chan Command
    errorChan    chan error
}

// ClientInterface defines the core client functionality
type ClientInterface interface {
    // Connection management
    Connect(ctx context.Context, serverURL string) error
    Disconnect() error
    IsConnected() bool
    
    // Authentication
    Login(username, password string) error
    Register(screenName, loginName, password string) error
    Logout() error
    IsAuthenticated() bool
    
    // Command handling
    SendCommand(cmd Command) error
    RegisterCommandHandler(cmdType string, handler CommandHandler)
    
    // Event handling
    Events() <-chan Event
    Commands() <-chan Command
    Errors() <-chan error
    
    // Configuration
    SetConfig(config *Config) error
    GetConfig() *Config
}

// NewClient creates a new ControlApp client
func NewClient(config *Config) *Client {
    return &Client{
        config:      config,
        auth:        NewAuthManager(config),
        ws:          NewWebSocketManager(config),
        commands:    NewCommandProcessor(),
        eventChan:   make(chan Event, 100),
        commandChan: make(chan Command, 100),
        errorChan:   make(chan error, 100),
    }
}
```

### Authentication Manager
```go
// internal/client/auth.go
package client

import (
    "encoding/json"
    "net/http"
    "time"
)

type AuthManager struct {
    config    *Config
    client    *http.Client
    token     string
    user      *User
    expiresAt time.Time
}

type LoginRequest struct {
    LoginName string `json:"login_name"`
    Password  string `json:"password"`
}

type RegisterRequest struct {
    ScreenName string `json:"screen_name"`
    LoginName  string `json:"login_name"`
    Password   string `json:"password"`
}

type AuthResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

func (a *AuthManager) Login(username, password string) error {
    req := LoginRequest{
        LoginName: username,
        Password:  password,
    }
    
    resp, err := a.makeAuthRequest("POST", "/api/v1/auth/login", req)
    if err != nil {
        return err
    }
    
    var authResp AuthResponse
    if err := json.Unmarshal(resp, &authResp); err != nil {
        return err
    }
    
    a.token = authResp.Token
    a.user = &authResp.User
    // Parse JWT to get expiration time
    a.expiresAt = a.parseTokenExpiration(authResp.Token)
    
    return nil
}

func (a *AuthManager) Register(screenName, loginName, password string) error {
    req := RegisterRequest{
        ScreenName: screenName,
        LoginName:  loginName,
        Password:   password,
    }
    
    _, err := a.makeAuthRequest("POST", "/api/v1/auth/register", req)
    return err
}

func (a *AuthManager) GetToken() string {
    if a.isTokenExpired() {
        return ""
    }
    return a.token
}

func (a *AuthManager) isTokenExpired() bool {
    return time.Now().After(a.expiresAt)
}
```

### WebSocket Manager
```go
// internal/client/websocket.go
package client

import (
    "encoding/json"
    "net/url"
    "sync"
    "time"
    
    "github.com/gorilla/websocket"
)

type WebSocketManager struct {
    config     *Config
    auth       *AuthManager
    conn       *websocket.Conn
    connected  bool
    mu         sync.RWMutex
    
    // Channels
    messageChan chan WSMessage
    errorChan   chan error
    closeChan   chan struct{}
    
    // Heartbeat
    lastPong    time.Time
    pingTicker  *time.Ticker
}

type WSMessage struct {
    Type      string      `json:"type"`
    Payload   interface{} `json:"payload"`
    Timestamp string      `json:"timestamp"`
}

func (ws *WebSocketManager) Connect(serverURL string, auth *AuthManager) error {
    ws.mu.Lock()
    defer ws.mu.Unlock()
    
    u, err := url.Parse(serverURL + "/ws/client")
    if err != nil {
        return err
    }
    
    // Add JWT token if authenticated
    if token := auth.GetToken(); token != "" {
        q := u.Query()
        q.Set("token", token)
        u.RawQuery = q.Encode()
    }
    
    conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        return err
    }
    
    ws.conn = conn
    ws.connected = true
    ws.lastPong = time.Now()
    
    // Start message handling goroutines
    go ws.readMessages()
    go ws.writeMessages()
    go ws.heartbeat()
    
    return nil
}

func (ws *WebSocketManager) SendMessage(msgType string, payload interface{}) error {
    msg := WSMessage{
        Type:      msgType,
        Payload:   payload,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    
    ws.mu.RLock()
    defer ws.mu.RUnlock()
    
    if !ws.connected {
        return ErrNotConnected
    }
    
    return ws.conn.WriteJSON(msg)
}

func (ws *WebSocketManager) readMessages() {
    defer func() {
        ws.mu.Lock()
        ws.connected = false
        ws.mu.Unlock()
        close(ws.closeChan)
    }()
    
    for {
        var msg WSMessage
        err := ws.conn.ReadJSON(&msg)
        if err != nil {
            ws.errorChan <- err
            return
        }
        
        // Handle special message types
        switch msg.Type {
        case "pong":
            ws.lastPong = time.Now()
        default:
            ws.messageChan <- msg
        }
    }
}
```

### Command Processor
```go
// internal/client/commands.go
package client

import (
    "encoding/json"
    "fmt"
    "sync"
)

type CommandProcessor struct {
    handlers map[string]CommandHandler
    mu       sync.RWMutex
}

type CommandHandler func(cmd Command) CommandResult

type Command struct {
    ID           string      `json:"id"`
    Type         string      `json:"type"`
    Content      interface{} `json:"content"`
    ReceivedAt   time.Time   `json:"received_at"`
}

type CommandResult struct {
    CommandID string      `json:"command_id"`
    Status    string      `json:"status"`    // "completed", "failed", "cancelled"
    Result    interface{} `json:"result"`
    Error     string      `json:"error,omitempty"`
    Duration  time.Duration `json:"duration"`
}

func NewCommandProcessor() *CommandProcessor {
    cp := &CommandProcessor{
        handlers: make(map[string]CommandHandler),
    }
    
    // Register core command handlers
    cp.registerCoreHandlers()
    
    return cp
}

func (cp *CommandProcessor) RegisterHandler(cmdType string, handler CommandHandler) {
    cp.mu.Lock()
    defer cp.mu.Unlock()
    cp.handlers[cmdType] = handler
}

func (cp *CommandProcessor) ProcessCommand(cmd Command) CommandResult {
    cp.mu.RLock()
    handler, exists := cp.handlers[cmd.Type]
    cp.mu.RUnlock()
    
    if !exists {
        return CommandResult{
            CommandID: cmd.ID,
            Status:    "failed",
            Error:     fmt.Sprintf("Unknown command type: %s", cmd.Type),
        }
    }
    
    start := time.Now()
    result := handler(cmd)
    result.Duration = time.Since(start)
    
    return result
}

func (cp *CommandProcessor) registerCoreHandlers() {
    // std_ping
    cp.RegisterHandler("std_ping", func(cmd Command) CommandResult {
        return CommandResult{
            CommandID: cmd.ID,
            Status:    "completed",
            Result: map[string]interface{}{
                "pong_timestamp": time.Now().Format(time.RFC3339Nano),
                "latency_ms":     0, // Will be calculated by caller
            },
        }
    })
    
    // std_popup - placeholder (platform-specific implementation)
    cp.RegisterHandler("std_popup", func(cmd Command) CommandResult {
        return CommandResult{
            CommandID: cmd.ID,
            Status:    "failed",
            Error:     "std_popup requires platform-specific implementation",
        }
    })
    
    // Add more core handlers...
}
```

## 🖥️ Windows Client Implementation

### Windows Entry Point
```go
// cmd/windows/main.go
package main

import (
    "context"
    "log"
    "os"
    
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/widget"
    
    "github.com/thecontrolapp/client/internal/client"
    windowsui "github.com/thecontrolapp/client/cmd/windows/ui"
)

func main() {
    // Create ControlApp client
    config := &client.Config{
        ServerURL: "ws://localhost:3080",
        AppName:   "ControlApp Windows Client",
        Version:   "1.0.0",
    }
    
    controlClient := client.NewClient(config)
    
    // Create Fyne app
    fyneApp := app.NewWithID("com.controlapp.windows")
    fyneApp.SetMetadata(&fyne.AppMetadata{
        Name: "ControlApp",
        Icon: resourceIconPng,
    })
    
    // Create main window with Windows-specific UI
    mainWindow := windowsui.NewMainWindow(fyneApp, controlClient)
    
    // Setup platform-specific command handlers
    setupWindowsCommandHandlers(controlClient)
    
    // Start the app
    mainWindow.ShowAndRun()
}

func setupWindowsCommandHandlers(client *client.Client) {
    // Windows-specific popup implementation
    client.RegisterCommandHandler("std_popup", func(cmd client.Command) client.CommandResult {
        content := cmd.Content.(map[string]interface{})
        body := content["body"].(string)
        
        // Use Windows native dialog or Fyne dialog
        dialog := widget.NewModalPopUp(
            widget.NewLabel(body),
            windowsui.GetCurrentWindow().Canvas(),
        )
        dialog.Show()
        
        return client.CommandResult{
            CommandID: cmd.ID,
            Status:    "completed",
            Result:    map[string]interface{}{"button_clicked": "OK"},
        }
    })
    
    // Add more Windows-specific handlers...
}
```

### Windows UI Package
```go
// cmd/windows/ui/main_window.go
package ui

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    
    "github.com/thecontrolapp/client/internal/client"
)

type MainWindow struct {
    app    fyne.App
    window fyne.Window
    client *client.Client
    
    // UI components
    connectBtn   *widget.Button
    statusLabel  *widget.Label
    commandList  *widget.List
    logText      *widget.Entry
}

func NewMainWindow(app fyne.App, client *client.Client) *MainWindow {
    window := app.NewWindow("ControlApp")
    window.Resize(fyne.NewSize(800, 600))
    
    mw := &MainWindow{
        app:    app,
        window: window,
        client: client,
    }
    
    mw.setupUI()
    mw.setupEventHandlers()
    
    return mw
}

func (mw *MainWindow) setupUI() {
    // Connection section
    mw.connectBtn = widget.NewButton("Connect", mw.onConnectClick)
    mw.statusLabel = widget.NewLabel("Disconnected")
    
    connectionSection := container.NewHBox(
        mw.connectBtn,
        mw.statusLabel,
    )
    
    // Commands section
    mw.commandList = widget.NewList(
        func() int { return 0 }, // Will be updated with actual commands
        func() fyne.CanvasObject { return widget.NewLabel("") },
        func(id widget.ListItemID, o fyne.CanvasObject) {},
    )
    
    // Log section
    mw.logText = widget.NewMultiLineEntry()
    mw.logText.SetText("Ready to connect...")
    
    // Layout
    content := container.NewVSplit(
        container.NewVBox(
            connectionSection,
            widget.NewSeparator(),
            widget.NewLabel("Commands:"),
            mw.commandList,
        ),
        container.NewVBox(
            widget.NewLabel("Log:"),
            container.NewScroll(mw.logText),
        ),
    )
    
    mw.window.SetContent(content)
}

func (mw *MainWindow) onConnectClick() {
    if !mw.client.IsConnected() {
        err := mw.client.Connect(context.Background(), "ws://localhost:3080")
        if err != nil {
            mw.logText.SetText(mw.logText.Text + "\nConnection failed: " + err.Error())
            return
        }
        mw.connectBtn.SetText("Disconnect")
        mw.statusLabel.SetText("Connected")
    } else {
        mw.client.Disconnect()
        mw.connectBtn.SetText("Connect")
        mw.statusLabel.SetText("Disconnected")
    }
}

func (mw *MainWindow) ShowAndRun() {
    mw.window.ShowAndRun()
}
```

## 🚀 Implementation Plan

### Phase 1: Core Client Package (1 week)
1. **Set up Go module structure**
2. **Implement `internal/client` package:**
   - Basic client interface and struct
   - Authentication manager
   - WebSocket connection manager
   - Command processor with core handlers
   - Configuration and error types
3. **Unit tests for core functionality**

### Phase 2: Windows Client (1 week) 
1. **Set up Fyne development environment**
2. **Implement Windows client:**
   - Main entry point (`cmd/windows/main.go`)
   - Windows-specific UI package
   - Platform-specific command handlers
   - Windows build configuration
3. **Integration testing with server**

### Phase 3: Cross-platform expansion
1. **macOS client** (reuse 95% of core + Windows UI patterns)
2. **Linux client** (same approach)
3. **Android client** (mobile-optimized UI, same core)

## 🔧 Build Commands

```bash
# Setup
mkdir client && cd client
go mod init github.com/thecontrolapp/client

# Install dependencies
go get fyne.io/fyne/v2/app
go get fyne.io/fyne/v2/widget
go get github.com/gorilla/websocket

# Build Windows client
go build -o controlapp-windows.exe ./cmd/windows

# Cross-compile for other platforms
GOOS=darwin go build -o controlapp-macos ./cmd/macos
GOOS=linux go build -o controlapp-linux ./cmd/linux

# Android (with fyne)
fyne package -os android -appID com.controlapp.client ./cmd/android
```

This architecture gives you:
- ✅ **95% code reuse** via `internal/client` package
- ✅ **Platform flexibility** - each client can have platform-specific adaptations
- ✅ **Shared server patterns** - reuse Go expertise and patterns
- ✅ **Simple builds** - standard Go toolchain
- ✅ **Easy testing** - standard Go testing for core functionality

Ready to start with the **Windows client implementation**? 🚀