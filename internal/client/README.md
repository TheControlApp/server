# ControlApp Client Package

The `internal/client` package provides a complete client implementation for connecting to ControlApp servers and handling kink commands with proper consent and safety mechanisms.

## Package Structure

- **`client.go`** - Main client interface and core functionality
- **`auth.go`** - Authentication management (JWT, login/register)
- **`websocket.go`** - WebSocket connection handling
- **`commands.go`** - Command processing and routing
- **`types.go`** - Type definitions and data structures
- **`logger.go`** - Logging interface and default implementation
- **`errors.go`** - Error definitions and handling
- **`windows.go`** - Windows-specific kink command implementations

## Features

### Core Client Functionality
- Server connection via WebSocket
- JWT-based authentication
- Real-time command processing
- Event-driven architecture
- Graceful error handling and recovery

### Authentication System
- User login and registration
- JWT token management
- Automatic token refresh
- Session persistence

### Command System
- Consent-based command execution
- Command validation and routing
- Audit logging for all commands
- Platform-specific implementations
- Safety checks and limits

### Kink Command Support
All kink commands are implemented with safety and consent as primary concerns:

- **Tier 1 (Basic)**: Messages, notifications, basic interactions
- **Tier 2 (Standard)**: File operations, media playback, system interactions  
- **Tier 3 (Advanced)**: System-level operations, scheduling, high-risk commands

## Usage

### Basic Client Setup
```go
import "github.com/thecontrolapp/server/internal/client"

// Create client with configuration
config := client.DefaultConfig()
config.ServerURL = "ws://localhost:3080/ws"

// Initialize client
c := client.NewClient(config)

// Connect to server
ctx := context.Background()
err := c.Connect(ctx, config.ServerURL)
if err != nil {
    log.Fatal("Connection failed:", err)
}

// Authenticate
err = c.Login("username", "password")
if err != nil {
    log.Fatal("Authentication failed:", err)
}
```

### Event Handling
```go
// Listen for events
go func() {
    for event := range c.Events() {
        switch event.Type {
        case client.EventCommandReceived:
            log.Printf("Command received: %+v", event.Data)
        case client.EventConnected:
            log.Println("Connected to server")
        case client.EventAuthenticated:
            log.Println("Successfully authenticated")
        }
    }
}()

// Listen for commands
go func() {
    for cmd := range c.Commands() {
        log.Printf("Processing command: %s", cmd.Type)
        // Commands are automatically processed by the client
    }
}()

// Listen for errors
go func() {
    for err := range c.Errors() {
        log.Printf("Client error: %v", err)
    }
}()
```

### Custom Command Handlers
```go
// Register a custom command handler
c.RegisterCommandHandler("custom_command", func(cmd client.Command) client.CommandResult {
    // Process the command
    return client.CommandResult{
        Status: "completed",
        Result: map[string]interface{}{
            "processed_at": time.Now(),
        },
    }
})
```

## Configuration

### Default Configuration
```go
config := client.DefaultConfig()
// Provides safe defaults:
// - Basic consent level
// - Limited command whitelist
// - File type restrictions
// - Size limits
// - Audit logging enabled
```

### Custom Configuration
```go
config := &client.Config{
    ServerURL: "ws://your-server:3080/ws",
    
    // Allow specific commands
    AllowedCommands: []string{
        "std_ping",
        "kink_message",
        "kink_open_link",
        "kink_play_audio",
    },
    
    // Block dangerous commands
    BlockedCommands: []string{
        "kink_run_file",
    },
    
    // Require confirmation for risky commands
    RequireConfirm: []string{
        "kink_lock_screen",
        "kink_change_wallpaper",
    },
    
    // File handling limits
    MaxFileSize: 50, // MB
    AllowedFileTypes: []string{".jpg", ".png", ".mp3"},
    DownloadFolder: "downloads",
    
    // Audio settings
    AudioVolume: 80,
    AllowTTS: true,
}
```

## Platform Support

### Windows Implementation
The `windows.go` file provides Windows-specific implementations for all kink commands:

- **Message Boxes**: Native Windows MessageBox API
- **File Operations**: PowerShell and Windows APIs
- **Audio Playback**: Windows Media Player COM objects
- **System Operations**: User32.dll API calls
- **Wallpaper Changes**: SystemParametersInfo API

### Cross-Platform Considerations
- Core client package is platform-agnostic
- Platform-specific implementations in separate files
- Fallback handlers for unsupported platforms
- Consistent error handling across platforms

## Safety and Security

### Consent Framework
- Users must explicitly enable each command type
- Per-command consent levels (Basic, Standard, Full)
- Real-time consent checking before execution
- Emergency disable functionality

### Security Measures
- All inputs are validated and sanitized
- File operations are sandboxed
- Command rate limiting
- Audit logging for all operations
- Encrypted WebSocket connections (WSS)

### Error Handling
- Graceful degradation on failures
- Detailed error messages for debugging
- Automatic reconnection on connection loss
- Circuit breaker patterns for reliability

## Development

### Building
```bash
go build ./internal/client
```

### Testing
```bash
go test ./internal/client
```

### Adding New Commands
1. Define command payload structure in `types.go`
2. Add validation logic in `commands.go`
3. Implement platform-specific handler (e.g., in `windows.go`)
4. Register handler in client initialization
5. Add tests and documentation

### Platform Support
To add support for a new platform:
1. Create platform-specific file (e.g., `linux.go`, `macos.go`)
2. Implement the required kink command handlers
3. Use build tags to conditionally compile
4. Ensure consistent error handling and return types

## Dependencies

- **WebSocket**: `github.com/gorilla/websocket`
- **Windows APIs**: `golang.org/x/sys/windows`
- **Server Models**: Uses `internal/models` for User type
- **Standard Library**: Uses only Go standard library for core functionality

## Architecture

The client follows a modular, event-driven architecture:

```
Client Core
├── WebSocket Manager (Real-time communication)
├── Auth Manager (JWT authentication)
├── Command Processor (Command routing and execution)
├── Event System (Async event handling)
└── Platform Handlers (OS-specific implementations)
```

This design provides:
- **Separation of Concerns**: Each component has a single responsibility
- **Testability**: Components can be mocked and tested independently
- **Extensibility**: New features can be added without breaking existing code
- **Maintainability**: Clear interfaces and well-defined responsibilities