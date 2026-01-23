# Test Commander

A test tool for sending commands to ControlApp clients to verify functionality.

## Purpose

This tool acts as a "commander" that can send various test commands to connected clients via WebSocket. It automatically registers an account if it doesn't exist.

## Features

- Auto-register if account doesn't exist
- Connect as an authenticated commander user
- Send individual or broadcast commands
- Test all command types (std_*, kink_*)
- Interactive menu for easy testing
- Pre-built command templates
- Support for targeted (specific user) and broadcast (all users) commands

## Usage

### Quick Start

```bash
# Just run it - it will register the account if needed
go run cmd/tools/test-commander/main.go -username=commander -password=commander123

# Or specify a different server
go run cmd/tools/test-commander/main.go -url=https://ctrlapp.merith.xyz -username=myuser -password=mypass
```

The tool will:
1. Try to login with provided credentials
2. If account doesn't exist, automatically register it
3. Connect to WebSocket and present interactive menu

### Testing Workflow

**Terminal 1 - Start a client:**
```bash
go run cmd/client/ctrlapp-console/main.go -username=client -password=client123
```

**Terminal 2 - Start the commander:**
```bash
go run cmd/tools/test-commander/main.go -username=commander -password=commander123
```

Select commands from the interactive menu and verify the client receives and processes them correctly.

## Example Commands

### Ping
```json
{
  "instructions": [
    {
      "type": "std_ping",
      "content": {"timestamp": "2024-01-22T10:00:00Z"}
    }
  ],
  "tags": "test"
}
```

### Popup Message
```json
{
  "instructions": [
    {
      "type": "std_popup",
      "content": {
        "body": "This is a test popup",
        "title": "Test",
        "timeout": 30
      }
    }
  ],
  "tags": "test"
}
```

### Timer
```json
{
  "instructions": [
    {
      "type": "std_timer",
      "content": {
        "duration": 60,
        "title": "Test Timer",
        "show_progress": true
      }
    }
  ],
  "tags": "test"
}
```
