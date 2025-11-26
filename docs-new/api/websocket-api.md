# WebSocket API Reference

## Overview

The ControlApp Server provides a real-time WebSocket API for command distribution and client communication. The WebSocket endpoint supports both authenticated and anonymous connections, with message routing capabilities for targeted and broadcast messaging.

## Connection

### WebSocket Endpoint

```
ws://localhost:8080/ws/client
```

### Authentication

WebSocket connections support three authentication methods:

1. **Authorization Header** (recommended):
   ```
   Authorization: Bearer <your-jwt-token>
   ```

2. **Query Parameter**:
   ```
   ws://localhost:8080/ws/client?token=<your-jwt-token>
   ```

3. **Message-based Authentication** (after connection):
   ```json
   {
     "type": "auth",
     "token": "<your-jwt-token>"
   }
   ```

### Connection Types

- **Authenticated Connection**: Full access to all features, user-specific messaging
- **Anonymous Connection**: Limited access, broadcast messaging only

## Connection Management

### Connection Limits

- **Token-based**: One connection per JWT token (new connections replace existing ones)
- **User-based**: Configurable maximum connections per user (default: unlimited)
- **Anonymous**: No limits (subject to server capacity)

### Connection Lifecycle

1. **Connection Established**: Client connects via WebSocket upgrade
2. **Authentication** (optional): Client sends auth message or uses header/query auth
3. **Message Exchange**: Bidirectional message communication
4. **Heartbeat**: Automatic ping/pong to maintain connection
5. **Disconnection**: Graceful or unexpected connection termination

## Message Protocol

### Message Structure

All WebSocket messages are JSON-formatted with consistent structure:

```json
{
  "type": "string",           // Message type identifier
  "timestamp": "ISO8601",     // Server-generated timestamp (outbound)
  "data": {},                 // Message-specific payload
  "id": "string"              // Unique message ID (optional)
}
```

## System Messages

### Authentication Messages

#### Auth Request
**Only supported authentication method via message:**
```json
{
  "type": "auth",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Note**: The server only supports JWT token authentication via the `auth` message type. Progressive authentication with username/password is not implemented.

#### Auth Success Response
```json
{
  "type": "auth_success",
  "message": "Authentication successful",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

#### Auth Error Response
```json
{
  "type": "error",
  "message": "Invalid or expired token"
}
```

#### Deprecated Authentication Methods
The following authentication methods are **NOT supported**:
- `auth_login` with username/password (documented in old versions but never implemented)
- Progressive authentication flows

### Heartbeat Messages

#### Ping (Client to Server)
```json
{
  "type": "ping"
}
```

#### Pong (Server to Client)
```json
{
  "type": "pong",
  "timestamp": 1697130000
}
```

### Error Messages

#### General Error
```json
{
  "type": "error",
  "message": "Descriptive error message"
}
```

## Command Messages

### Command Structure

Commands are the primary message type for instruction delivery:

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "instructions": [
    {
      "type": "instruction_type",
      "content": {}
    }
  ],
  "sender_id": "123e4567-e89b-12d3-a456-426614174001",
  "receiver_id": "123e4567-e89b-12d3-a456-426614174000",
  "tags": "general,test",
  "status": "pending",
  "created_at": "2025-10-12T18:00:00Z",
  "updated_at": "2025-10-12T18:00:00Z"
}
```

### Sending Commands (Client to Server)

#### Complete Command
```json
{
  "instructions": [
    {
      "type": "std_popup",
      "content": {
        "body": "This is a test message",
        "button": "OK"
      }
    }
  ],
  "receiver_id": null,
  "tags": "general"
}
```

#### Legacy Format (Backward Compatibility)
```json
{
  "type": "std_popup",
  "content": {
    "body": "This is a test message",
    "button": "OK"
  }
}
```

### Command Fields

- **`instructions`** (array, required): Array of instruction objects
- **`receiver_id`** (UUID, optional): Target user ID (null for broadcast)
- **`tags`** (string, required): Comma-separated tags for categorization
- **`sender_id`** (UUID, auto-set): Sending user ID (server-controlled)
- **`status`** (string, auto-set): Command status (server-controlled)
- **`id`** (UUID, auto-set): Unique command identifier (server-controlled)

### Command Validation

Commands must pass the following validation:

1. **Instructions Array**: Must contain at least one instruction
2. **Tags**: Must be provided (non-empty string)
3. **Instruction Types**: Each instruction must have a valid `type` field
4. **Instruction Content**: Each instruction must have a `content` field
5. **Authentication**: Some instruction types require authentication

### Authentication Requirements

The following instruction types require authentication:
- All types except `ping` and `std_test`
- Anonymous clients can only send basic test commands

## Instruction Types

### Standard Instructions

#### Popup Message
```json
{
  "type": "std_popup",
  "content": {
    "body": "Message text",
    "button": "Button text"
  }
}
```

#### Timer
```json
{
  "type": "std_timer",
  "content": {
    "duration": 60,
    "title": "Timer Title"
  }
}
```

#### Display Text
```json
{
  "type": "display_text",
  "content": {
    "text": "Text to display",
    "format": "plain"
  }
}
```

#### Notification
```json
{
  "type": "notification",
  "content": {
    "title": "Notification Title",
    "body": "Notification message"
  }
}
```

#### Open URL
```json
{
  "type": "open_url",
  "content": {
    "url": "https://example.com"
  }
}
```

#### Download File
```json
{
  "type": "download_file",
  "content": {
    "file_hash": "sha256hash",
    "file_name": "example.pdf"
  }
}
```

### Form Input
```json
{
  "type": "form_input",
  "content": {
    "fields": [
      {
        "name": "username",
        "label": "Username",
        "type": "text",
        "required": true
      }
    ],
    "submit_to": "api_endpoint"
  }
}
```

### Custom Instructions

Custom instruction types are supported by providing arbitrary content:

```json
{
  "type": "custom_instruction",
  "content": {
    "any_field": "any_value",
    "nested": {
      "data": "supported"
    }
  }
}
```

## Message Routing

### Broadcast Messages
- **Receiver ID**: `null` or omitted
- **Delivery**: Sent to all connected clients
- **Use Case**: General announcements, public commands

### Targeted Messages
- **Receiver ID**: Specific user UUID
- **Delivery**: Sent only to specified user's connections
- **Use Case**: Private commands, direct communication

### Connection Behavior

#### Multiple Connections
- Users can have multiple connections (configurable limit)
- Messages sent to all of a user's active connections
- Token replacement: New connection with same token replaces existing

#### Connection States
- **Active**: Receiving all applicable messages
- **Stale**: Connection lost, automatic cleanup
- **Rate Limited**: Temporary message throttling (if implemented)

## Connection Examples

### JavaScript Client

```javascript
// Basic connection
const ws = new WebSocket('ws://localhost:8080/ws/client');

// With authentication
const wsAuth = new WebSocket('ws://localhost:8080/ws/client?token=' + jwtToken);

ws.onopen = function() {
    console.log('Connected to WebSocket');
    
    // Authenticate after connection
    ws.send(JSON.stringify({
        type: 'auth',
        token: 'your-jwt-token'
    }));
};

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('Received:', message);
};

// Send a command
ws.send(JSON.stringify({
    instructions: [{
        type: 'std_popup',
        content: {
            body: 'Hello World!',
            button: 'OK'
        }
    }],
    tags: 'general'
}));
```

### Python Client

```python
import websocket
import json

def on_message(ws, message):
    data = json.loads(message)
    print(f"Received: {data}")

def on_open(ws):
    print("Connected to WebSocket")
    
    # Send authentication
    auth_msg = {
        "type": "auth",
        "token": "your-jwt-token"
    }
    ws.send(json.dumps(auth_msg))

# Connect with authentication header
headers = {"Authorization": "Bearer your-jwt-token"}
ws = websocket.WebSocketApp("ws://localhost:8080/ws/client",
                          header=headers,
                          on_open=on_open,
                          on_message=on_message)
ws.run_forever()
```

## Error Handling

### Common Error Scenarios

1. **Invalid JSON**: Malformed message structure
2. **Authentication Required**: Command requires authenticated user
3. **Invalid Command**: Missing required fields
4. **Connection Limit**: Maximum connections exceeded
5. **Token Conflict**: Token already in use by another connection

### Error Response Format

```json
{
  "type": "error",
  "message": "Specific error description"
}
```

### Error Recovery

- **Connection Errors**: Implement automatic reconnection with exponential backoff
- **Authentication Errors**: Re-authenticate with fresh token
- **Message Errors**: Validate message structure before sending

## Performance Considerations

### Message Limits

- **Message Size**: Maximum 512 bytes per message
- **Connection Timeout**: 60 seconds for pong response
- **Write Timeout**: 10 seconds for message delivery

### Optimization

- **Message Caching**: Server caches broadcast messages for optimization
- **Connection Pooling**: Efficient management of multiple connections
- **Heartbeat**: Automatic cleanup of stale connections

## Security

### Authentication
- JWT token validation for authenticated features
- Token expiration handling
- Secure token transmission

### Authorization
- Instruction-level access control
- User-based command filtering
- Anonymous connection restrictions

### Connection Security
- CORS configuration for web clients
- Origin validation
- Rate limiting (configurable)

## Development and Testing

### WebSocket Testing Tools

1. **Browser DevTools**: WebSocket frame inspection
2. **Postman**: WebSocket connection testing
3. **wscat**: Command-line WebSocket client
4. **Custom Test Clients**: See `/examples` directory

### Debug Mode

Enable debug logging to see detailed WebSocket message flow:

```bash
LOG_LEVEL=debug go run cmd/server/main.go
```

## Migration Notes

### Legacy Compatibility

The server supports legacy message formats for backward compatibility:

```json
// Legacy format (deprecated)
{
  "type": "std_popup",
  "content": {"body": "message"}
}

// New format (recommended)
{
  "instructions": [{"type": "std_popup", "content": {"body": "message"}}],
  "tags": "general"
}
```

### Breaking Changes

- Version 1.0: Introduction of structured command format
- Future versions will phase out legacy message support