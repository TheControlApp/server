# WebSocket Message Format Reference

This document defines the complete JSON data structures used for WebSocket communication between clients and the ControlMe server.

## Table of Contents

- [Message Structure Overview](#message-structure-overview)
- [Server-to-Client Messages](#server-to-client-messages)
- [Client-to-Server Messages](#client-to-server-messages)
- [Instruction Types](#instruction-types)
- [Command Status Flow](#command-status-flow)
- [Error Handling](#error-handling)

## Message Structure Overview

All WebSocket messages follow a standard envelope format:

```json
{
  "type": "string",        // Message type identifier
  "id": "uuid",           // Unique message ID
  "timestamp": "iso8601", // Message timestamp
  "from": "uuid",         // Sender user ID
  "to": "uuid",           // Recipient user ID (optional for broadcasts)
  "data": {}              // Message-specific payload
}
```

## Server-to-Client Messages

### 1. Command Assignment

When the server assigns a new command to a client:

```json
{
  "type": "command_assignment",
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "timestamp": "2025-09-27T10:30:00Z",
  "from": "dom-user-uuid",
  "to": "sub-user-uuid",
  "data": {
    "command": {
      "id": "cmd-uuid",
      "instructions": [
        {
          "type": "std_popup",
          "content": {
            "title": "Task Assignment",
            "body": "Please complete this task within 30 minutes.",
            "button": "Acknowledge"
          }
        }
      ],
      "sender_id": "dom-user-uuid",
      "receiver_id": "sub-user-uuid",
      "status": "pending",
      "created_at": "2025-09-27T10:30:00Z"
    }
  }
}
```

**Instructions Array Structure:**

The `instructions` field contains a JSON array of instruction objects:

```javascript
// Access instructions directly as an array
const instructions = command.instructions;
// Process each instruction
instructions.forEach(instruction => {
  console.log(instruction.type, instruction.content);
});
```

**Array Structure:**
```json
[
  {
    "type": "std_popup",
    "content": {
      "title": "Task Assignment",
      "body": "Please complete this task within 30 minutes.",
      "button": "Acknowledge"
    }
  }
]
```

### 2. Command Status Update

When a command status changes:

```json
{
  "type": "command_status",
  "id": "msg-uuid",
  "timestamp": "2025-09-27T10:35:00Z",
  "from": "server-uuid",
  "to": "dom-user-uuid",
  "data": {
    "command_id": "cmd-uuid",
    "status": "completed",
    "completed_by": "sub-user-uuid",
    "completed_at": "2025-09-27T10:35:00Z"
  }
}
```

### 3. Connection Status

Server connection and authentication confirmations:

```json
{
  "type": "connection_status",
  "id": "msg-uuid",
  "timestamp": "2025-09-27T10:00:00Z",
  "from": "server-uuid",
  "to": "user-uuid",
  "data": {
    "status": "authenticated",
    "user_id": "user-uuid",
    "session_id": "session-uuid"
  }
}
```

### 4. Error Response

When an error occurs:

```json
{
  "type": "error",
  "id": "msg-uuid",
  "timestamp": "2025-09-27T10:30:00Z",
  "from": "server-uuid",
  "to": "user-uuid",
  "data": {
    "error_code": "INVALID_COMMAND",
    "error_message": "Command format is invalid",
    "reference_id": "original-msg-uuid"
  }
}
```

## Client-to-Server Messages

### 1. Command Completion

When a client completes a command:

```json
{
  "type": "command_completion",
  "id": "msg-uuid",
  "timestamp": "2025-09-27T10:35:00Z",
  "from": "sub-user-uuid",
  "to": "dom-user-uuid",
  "data": {
    "command_id": "cmd-uuid",
    "status": "completed",
    "response": {
      "success": true,
      "data": {
        "user_input": "Task completed successfully",
        "screenshot": "base64-encoded-image"
      }
    }
  }
}
```

### 2. Heartbeat/Ping

Client keepalive messages:

```json
{
  "type": "ping",
  "id": "msg-uuid",
  "timestamp": "2025-09-27T10:30:00Z",
  "from": "user-uuid",
  "to": "server-uuid",
  "data": {
    "client_info": {
      "version": "1.0.0",
      "platform": "web"
    }
  }
}
```

### 3. Status Update

Client status changes:

```json
{
  "type": "status_update",
  "id": "msg-uuid",
  "timestamp": "2025-09-27T10:30:00Z",
  "from": "user-uuid",
  "to": "server-uuid",
  "data": {
    "status": "online",
    "availability": "available"
  }
}
```

## Multi-Instruction Commands

Commands can contain multiple instructions that are executed in sequence. This allows for complex workflows and multi-step tasks.

### Multi-Instruction Example

```json
{
  "instructions": [
    {
      "type": "std_notification",
      "content": {
        "title": "Starting Task",
        "body": "Beginning workflow"
      }
    },
    {
      "type": "std_popup", 
      "content": {
        "title": "Step 1",
        "body": "Complete preparation",
        "button": "Done"
      }
    },
    {
      "type": "std_timer",
      "content": {
        "duration": 300,
        "title": "Work Session"
      }
    },
    {
      "type": "std_input",
      "content": {
        "title": "Report",
        "fields": [
          {
            "name": "result",
            "label": "Outcome", 
            "type": "text"
          }
        ]
      }
    }
  ]
}
```

**Instructions Array:**
```json
[
  {
    "type": "std_notification",
    "content": {
      "title": "Starting Task",
      "body": "Beginning workflow"
    }
  },
  {
    "type": "std_popup", 
    "content": {
      "title": "Step 1",
      "body": "Complete preparation",
      "button": "Done"
    }
  },
  {
    "type": "std_timer",
    "content": {
      "duration": 300,
      "title": "Work Session"
    }
  },
  {
    "type": "std_input",
    "content": {
      "title": "Report",
      "fields": [
        {
          "name": "result",
          "label": "Outcome", 
          "type": "text"
        }
      ]
    }
  }
]
```

### Execution Order

Instructions should be executed in **sequential order** (array index 0, 1, 2, etc.). Each instruction should be completed before moving to the next, unless the instruction type specifically allows parallel execution.

### Completion Tracking

For multi-instruction commands, clients should track:
- Number of instructions completed
- Total execution time
- Any intermediate results
- Final completion status

## Instruction Types

Commands contain arrays of instructions. Each instruction has a `type` and `content` structure:

### Instruction Structure
```json
{
  "type": "instruction_type_name",
  "content": { /* arbitrary JSON data */ }
}
```

**Important:** The `content` field is completely flexible and can contain any JSON structure. Different instruction types define their own content schemas, but clients can also define custom instruction types with their own content formats.

### Standard Instruction Types

#### 1. `std_popup` - Display Modal Dialog

```json
{
  "type": "std_popup",
  "content": {
    "title": "Task Title",
    "body": "Detailed task description or message",
    "button": "Button Text",
    "timeout": 300,           // Optional: auto-close after seconds
    "required": true          // Optional: must acknowledge
  }
}
```

#### 2. `std_notification` - System Notification

```json
{
  "type": "std_notification",
  "content": {
    "title": "Notification Title",
    "body": "Notification message",
    "priority": "high",       // low, normal, high, urgent
    "sound": true,           // Optional: play sound
    "persistent": false      // Optional: stays until dismissed
  }
}
```

#### 3. `std_timer` - Countdown Timer

```json
{
  "type": "std_timer",
  "content": {
    "duration": 1800,        // Duration in seconds
    "title": "Timer Title",
    "description": "Timer description",
    "show_progress": true,   // Optional: show progress bar
    "auto_complete": false   // Optional: auto-complete when done
  }
}
```

#### 4. `std_input` - User Input Form

```json
{
  "type": "std_input",
  "content": {
    "title": "Input Required",
    "description": "Please provide the following information",
    "fields": [
      {
        "name": "response",
        "label": "Your Response",
        "type": "text",         // text, number, select, checkbox
        "required": true,
        "options": ["Option 1", "Option 2"]  // For select type
      }
    ],
    "submit_button": "Submit",
    "cancel_button": "Cancel"
  }
}
```

### Extended Instruction Types

#### 5. `download_file` - File Download

```json
{
  "type": "download_file",
  "content": {
    "file_hash": "sha256-hash",
    "file_name": "document.pdf",
    "file_size": 1024000,
    "download_url": "https://server.com/api/files/download/hash",
    "description": "Download this file and save to Desktop"
  }
}
```

#### 6. `open_url` - Open Website

```json
{
  "type": "open_url",
  "content": {
    "url": "https://example.com",
    "title": "Website Title",
    "description": "Visit this website",
    "new_window": true,
    "wait_for_close": false
  }
}
```

#### 7. `display_text` - Text Display

```json
{
  "type": "display_text",
  "content": {
    "text": "Text to display",
    "format": "plain",        // plain, html, markdown
    "title": "Display Title",
    "duration": 0,           // 0 = until dismissed
    "position": "center"     // center, top, bottom
  }
}
```

### Custom Instruction Types

Clients can define their own instruction types with arbitrary content structures:

```json
{
  "type": "custom_game_command",
  "content": {
    "game_id": "tetris",
    "level": 5,
    "settings": {
      "speed": "fast",
      "music": true,
      "effects": ["particle", "glow"]
    },
    "player_data": {
      "high_score": 15000,
      "achievements": ["first_clear", "speed_demon"]
    }
  }
}
```

```json
{
  "type": "hardware_control",
  "content": {
    "device": "smart_lights",
    "action": "color_sequence",
    "parameters": {
      "colors": ["#FF0000", "#00FF00", "#0000FF"],
      "duration_ms": 2000,
      "repeat": 3,
      "brightness": 0.8
    }
  }
}
```

The server will pass through any instruction type and content structure to the client unchanged, allowing for maximum flexibility in client implementations.

## Server Implementation Notes

**Go Server Implementation:**

The Go server uses proper typed structs for instructions:

```go
type Command struct {
    ID           uuid.UUID     `json:"id"`
    Instructions []Instruction `gorm:"serializer:json;type:text" json:"instructions"`
    // ... other fields
}

type Instruction struct {
    Type    string      `json:"type"`    // Instruction type (std_popup, etc.)
    Content interface{} `json:"content"` // Arbitrary data - any JSON structure
}
```

**Key Features:**
- Go clients can work directly with `[]Instruction` slices - no string parsing needed
- GORM automatically serializes the slice to/from JSON in the database
- The `Content` field uses `interface{}` for maximum flexibility
- Server passes through all instruction types and content unchanged
- Database storage is handled transparently by GORM's JSON serializer

**Usage in Go:**
```go
// Creating commands in Go
cmd := &models.Command{
    Instructions: []models.Instruction{
        {
            Type: "std_popup",
            Content: map[string]interface{}{
                "title": "Hello",
                "body": "World",
            },
        },
        {
            Type: "custom_type",
            Content: YourCustomStruct{
                Field1: "value1",
                Field2: 42,
            },
        },
    },
}
```

## Command Status Flow

Commands progress through these statuses:

1. **`pending`** - Command assigned, awaiting client acknowledgment
2. **`delivered`** - Client has received and acknowledged command
3. **`in_progress`** - Client is actively working on command (optional)
4. **`completed`** - Command successfully completed
5. **`failed`** - Command failed or was rejected
6. **`cancelled`** - Command was cancelled by sender

## Error Handling

### Client Error Response Format

When a client encounters an error processing a command:

```json
{
  "type": "command_completion",
  "id": "msg-uuid",
  "timestamp": "2025-09-27T10:35:00Z",
  "from": "sub-user-uuid",
  "to": "dom-user-uuid",
  "data": {
    "command_id": "cmd-uuid",
    "status": "failed",
    "error": {
      "error_code": "UNSUPPORTED_INSTRUCTION",
      "error_message": "Instruction type 'custom_type' is not supported",
      "instruction_index": 2
    }
  }
}
```

### Server Error Codes

Common error codes the server may send:

- `AUTHENTICATION_FAILED` - Invalid or expired JWT token
- `INVALID_MESSAGE_FORMAT` - Message doesn't match expected JSON structure
- `COMMAND_NOT_FOUND` - Referenced command ID doesn't exist
- `PERMISSION_DENIED` - User lacks permission for requested action
- `RATE_LIMITED` - Too many messages sent too quickly
- `SERVER_ERROR` - Internal server error

## Implementation Notes

### Message IDs

- All messages must include a unique `id` field (UUID format recommended)
- Responses should reference the original message `id` when applicable
- Use `reference_id` in error responses to link to the failed message

### Timestamps

- All timestamps use ISO 8601 format: `YYYY-MM-DDTHH:mm:ssZ`
- Server timestamps are authoritative
- Client timestamps should match their local timezone

### User IDs

- All user IDs are UUIDs
- `from` field identifies the message sender
- `to` field identifies the recipient (omit for broadcast messages)
- Server messages use a special server UUID for the `from` field

### Content Validation

- Server validates all instruction `content` structures
- Unknown instruction types are rejected with `UNSUPPORTED_INSTRUCTION` error
- Missing required fields result in `INVALID_CONTENT` error
- Clients should validate incoming instruction types before processing

## Examples

See the `/docs/examples/` directory for complete working examples of:
- Client WebSocket connection and authentication
- Sending and receiving commands
- Handling different instruction types
- Error handling and recovery