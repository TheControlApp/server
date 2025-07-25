# WebSocket API

## Connection
```
WS /api/ws
Authorization: Bearer <jwt_token> (header or query param)
```

## Message Format
All messages are JSON:
```json
{
  "type": "message_type",
  "data": { }
}
```

## Message Types

### Receive Commands (Server → Client)
```json
{
  "type": "command",
  "data": {
    "id": "uuid",
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "Hello!",
          "button": "OK"
        }
      }
    ],
    "sender": "username",
    "receiver": "you",
    "tags": ["general"]
  }
}
```

### Send Commands (Client → Server)  
```json
{
  "type": "send_command",
  "data": {
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "Message",
          "button": "OK"
        }
      }
    ],
    "receiver": "username",
    "tags": ["general"]
  }
}
```

### Command Status Updates
```json
{
  "type": "command_status",
  "data": {
    "command_id": "uuid",
    "status": "received|completed|failed"
  }
}
```

### Heartbeat
```json
{
  "type": "heartbeat",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Error Messages
```json
{
  "type": "error",
  "error": "error_code",
  "message": "Description"
}
```

## Common Errors
- `authentication_failed` - Invalid token
- `rate_limit_exceeded` - Too many commands
- `user_not_found` - Invalid receiver
- `command_blocked` - Content filtered
        }
      },
      {
        "type": "download-file",
        "content": {
          "file_hash": "abc123def456...",
          "file_name": "task.pdf"
        }
      }
    ],
    "sender": "master123",
    "receiver": "user456", 
    "tags": ["chastity", "daily"],
    "status": "pending",
    "created_at": "2025-07-24T15:30:00Z"
  }
}
```

### 2. Command Sending (Client → Server)

Send commands to specific users:

```json
{
  "type": "send_command",
  "data": {
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "Take a 15-minute break",
          "button": "Done"
        }
      }
    ],
    "receiver": "user456",
    "tags": ["general", "break"]
  }
}
```

### 3. Broadcast Commands

Send commands to users subscribed to specific tags:

```json
{
  "type": "send_command",
  "data": {
    "instructions": [
      {
        "type": "announcement",
        "content": {
          "title": "System Maintenance",
          "body": "Server will be down for 10 minutes at 3 PM",
          "priority": "high"
        }
      }
    ],
    "tags": ["announcements", "system"]
  }
}
```

### 4. Command Status Updates

Update command completion status:

```json
{
  "type": "command_status",
  "data": {
    "command_id": "550e8400-e29b-41d4-a716-446655440001",
    "status": "completed",
    "completed_at": "2025-07-24T15:45:00Z"
  }
}
```

### 5. Heartbeat

Maintain connection with periodic heartbeats:

```json
{
  "type": "heartbeat",
  "timestamp": "2025-07-24T15:30:00Z"
}
```

## Connection Management

### Heartbeat System
- **Interval**: 30 seconds
- **Timeout**: 3 missed heartbeats
- **Client**: Must respond to server heartbeats
- **Server**: Automatically disconnects inactive connections

### Reconnection Strategy
```javascript
function connectWebSocket() {
  const ws = new WebSocket('ws://localhost:8080/api/ws?token=' + authToken);
  
  ws.onopen = () => {
    console.log('Connected');
    heartbeatInterval = setInterval(() => {
      ws.send(JSON.stringify({
        type: 'heartbeat',
        timestamp: new Date().toISOString()
      }));
    }, 30000);
  };
  
  ws.onclose = (event) => {
    clearInterval(heartbeatInterval);
    if (event.code !== 1000) {
      // Reconnect after 5 seconds
      setTimeout(connectWebSocket, 5000);
    }
  };
}
```

### Queue Delivery
- Commands sent while offline are queued
- Delivered immediately upon reconnection
- 2-week retention period for undelivered commands

## Error Handling

### Error Message Format
```json
{
  "type": "error",
  "error": "error_code",
  "message": "Human-readable error message",
  "command_id": "uuid" // optional, if related to specific command
}
```

### Common Errors
```json
{
  "type": "error",
  "error": "authentication_failed",
  "message": "Invalid or expired token"
}

{
  "type": "error",
  "error": "command_failed",
  "message": "Unable to deliver command to target user",
  "command_id": "550e8400-e29b-41d4-a716-446655440001"
}

{
  "type": "error",
  "error": "rate_limit_exceeded",
  "message": "Too many commands sent. Please wait before sending more."
}

{
  "type": "error",
  "error": "user_not_found",
  "message": "Target user 'username' does not exist"
}
```

## Broadcasting Rules

### Direct Messages
- Commands with `receiver` field go to specific user
- Sender must not be blocked by receiver
- Receiver must be online or command is queued

### Tag-Based Broadcasting
- Commands without `receiver` field broadcast to tag subscribers
- Users control which tags they receive via preferences
- Multiple tags = command sent to users subscribed to ANY tag

### Content Filtering
Available tags for filtering:
- `general` - Safe-for-work content
- `adult` - Explicit content  
- `chastity` - Chastity-related content
- `feet` - Foot-related content
- `roleplay` - Roleplay scenarios
- `humiliation` - Humiliation-based content
- `public` - Public setting commands
- `private` - Private/intimate commands

## Security Features

### Authentication
- JWT token required for all connections
- Tokens expire after 24 hours
- Invalid tokens result in immediate disconnection

### Rate Limiting
- 10 commands per minute per user
- Burst allowance of 5 additional commands
- Temporary blocks for abuse

### Content Moderation
- File hash scanning for CSAM/malware
- User reporting system
- Automatic flagging of suspicious content
- Admin intervention capabilities
