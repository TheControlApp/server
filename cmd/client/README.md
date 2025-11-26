# ControlApp Clients

This directory contains the various client implementations for ControlApp, organized by client type as requested.

## Client Structure

```
cmd/client/
├── ctrlapp-console/     # Minimal console client for testing and basic communication
└── ctrlapp-windows/     # Full GUI client with system tray support (Windows-first)
```

### Console Client (`ctrlapp-console/`)

**Purpose**: Minimal console client for basic communication and testing
**Current State**: ✅ Complete and building successfully

**Features**:
- Connection management and authentication
- Unified command support (std_ping, std_message, std_notification)
- Console-based output for all commands
- Credential storage and auto-login capability
- Simple command-line interface

**Use Cases**:
- Testing server connectivity and authentication
- Basic command sending and receiving
- Reference implementation for shared client functionality
- Headless operation and automation

**Build**: `go build -o bin/ctrlapp-console.exe ./cmd/client/ctrlapp-console`

### Windows GUI Client (`ctrlapp-windows/`)

**Purpose**: Full GUI client with system tray daemon functionality
**Current State**: ✅ Complete structure, building successfully

**Features**:
- **User Registration**: Full user registration support with display name, username, and password
- **Auto-login**: Attempts credential-based auto-login on startup
- **System Tray Integration**: Minimizes to Windows system tray with context menu
- **Background Mode**: Runs as daemon when minimized, processing commands
- **Command Center**: GUI window for monitoring activity and sending test commands
- **Settings Panel**: Configuration interface for command consent settings
- **Credential Management**: Secure storage and caching of authentication tokens
- **Unified Commands**: Uses consolidated std_* command set (ping, message, notification, etc.)
- **Thread Safety**: Proper Fyne UI thread management for background operations

**Key Components**:
- Login window with consent notice
- Background daemon mode with system tray
- Command center for activity monitoring
- Settings window for consent configuration
- Automatic credential storage and retrieval

**Supported Commands**:
- Core (always allowed): std_ping, std_message, std_notification
- User configurable: std_open_url, std_download_file, std_play_audio, std_display_image, std_timer
- Consent required: std_change_wallpaper, std_lock_screen, std_execute_file

**Build**: `go build -o bin/ctrlapp-windows.exe ./cmd/client/ctrlapp-windows`

**Dependencies**: 
- Fyne v2.7.1 for cross-platform GUI
- Windows-first design with system tray support

## Shared Components

Both clients utilize shared components from `internal/client/`:

### Credential Store (`internal/client/credentials.go`)
- JWT token storage with 1-week validity
- Secure file-based storage in `~/.controlapp/` directory
- Automatic token refresh and validation
- Cross-platform path handling

### Client Library (`internal/client/`)
- Connection management and WebSocket handling
- Authentication and user management
- Command processing and response handling
- Event system for client state changes
- Cross-platform command handlers

## Command Specification

All clients now use the unified command specification documented in `docs/UNIFIED_COMMANDS.md`. The previous separation between `kink_*` and `std_*` commands has been consolidated into a single standard command set with appropriate consent levels.

## Development Priorities

1. ✅ **Console Client**: Complete - minimal viable client for testing
2. ✅ **Windows GUI Client**: Complete structure - full GUI with system tray
3. 🔄 **Future**: Android client, Linux GUI client (using build tags)

## Consent Model

**Consent-by-Running**: Users consent to command execution by running the client application. The consent model is clearly communicated through:

1. **Explicit Notice**: Login window displays consent agreement
2. **Configurable Settings**: Users can configure which command types to allow
3. **Risk Classification**: Commands are categorized by risk level
4. **User Control**: Settings panel allows granular control over command consent

## Architecture Notes

- **Windows-First Development**: Primary focus on Windows compatibility with cross-platform preparation
- **Shared Libraries**: Maximum code reuse between client types
- **Build Tags**: Prepared for platform-specific features (system tray, notifications, etc.)
- **Security**: JWT-based authentication with secure credential storage
- **Modularity**: Clean separation between UI, networking, and business logic

## Usage

### Console Client
```bash
# Build and run console client
go build -o bin/ctrlapp-console.exe ./cmd/client/ctrlapp-console
./bin/ctrlapp-console.exe
```

### Windows GUI Client
```bash
# Build and run Windows GUI client
go build -o bin/ctrlapp-windows.exe ./cmd/client/ctrlapp-windows
./bin/ctrlapp-windows.exe
```

The GUI client will:
1. Attempt auto-login with stored credentials
2. If successful, start in background/system tray mode
3. If no valid credentials, show login window
4. Right-click system tray icon for menu options

## Configuration

Both clients store configuration in:
- **Windows**: `%USERPROFILE%\.controlapp\`
- **Linux/macOS**: `~/.controlapp/`

Files:
- `credentials.json` - Encrypted JWT tokens and authentication data
- `settings.json` - User consent settings and preferences (future)