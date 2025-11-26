# WebSocket Implementation Guide

## Overview

The ControlMe server implements a flexible WebSocket architecture that supports both anonymous and authenticated connections, enabling real-time command distribution and progressive authentication.

## Architecture

### Connection Types

1. **Anonymous Connections**
   - No authentication required
   - Receive broadcast messages only
   - Can upgrade to authenticated session
   - Limited functionality

2. **Authenticated Connections**
   - JWT token required (header or query parameter)
   - Full access to user-specific commands
   - Can send commands and receive responses
   - Complete functionality

### Progressive Authentication Flow

```
1. Client connects anonymously → ws://localhost:8080/ws/client
2. Server accepts connection, assigns anonymous session
3. Client sends auth_login or auth_token message
4. Server validates credentials and upgrades session
5. Client now has full authenticated access
```

## Implementation Details

### Server-Side Components

#### WebSocket Handler (`internal/api/handlers/websocket_handlers.go`)
- Handles connection upgrades
- Manages authentication (token extraction from headers/query)
- Routes messages to appropriate handlers
- Manages connection lifecycle

#### Hub (`internal/websocket/hub.go`)
- Central connection manager
- Maintains client registry (anonymous + authenticated)
- Handles message broadcasting
- Enforces connection limits
- Manages user sessions

#### Client Structure
```go
type Client struct {
    UserID     *uint             // nil for anonymous clients
    Username   string            // empty for anonymous clients  
    Token      string            // JWT token (if authenticated)
    Connection *websocket.Conn   // WebSocket connection
    Hub        *Hub              // Reference to hub
    Send       chan []byte       // Outgoing message channel
    // ... other fields
}
```

### Authentication Methods

#### 1. Header-Based Authentication
```http
GET /ws/client HTTP/1.1
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### 2. Query Parameter Authentication  
```http
GET /ws/client?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9... HTTP/1.1
```

#### 3. Progressive Authentication
```javascript
// 1. Connect anonymously
const ws = new WebSocket('ws://localhost:8080/ws/client');

// 2. Send authentication message after connection
ws.send(JSON.stringify({
    type: 'auth_login',
    payload: {
        username: 'testuser',
        password: 'password123'
    }
}));

// 3. Handle authentication response
ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === 'auth_success') {
        console.log('Authentication successful:', message.payload);
        // Now can send authenticated commands
    }
};
```

## Message Routing

### Inbound Message Flow
```
WebSocket Message → Handler → Parse Type → Route to Processor → Response
```

### Message Processors

#### Authentication Processor
- Handles `auth_login` and `auth_token` messages
- Validates credentials via UserService
- Updates client session state
- Sends authentication response

#### Command Processor  
- Handles `execute_command` messages
- Validates user permissions
- Executes commands via CommandService
- Streams results back to client

#### System Processor
- Handles `ping`, status requests, etc.
- Provides connection health monitoring
- Returns system information

### Outbound Message Flow
```
Service → Hub → Target Clients → WebSocket Connection → Client
```

## Security Considerations

### Authentication Security
- JWT tokens validated on connection and message processing
- Token expiration handled gracefully (session downgrade)
- Rate limiting prevents brute force attacks
- Invalid tokens result in anonymous session (not disconnection)

### Authorization Levels
- **Anonymous:** Broadcast messages only
- **Authenticated:** User-specific commands + broadcasts  
- **Admin:** All messages + administrative commands

### Connection Security
- Configurable connection limits per user
- IP-based rate limiting
- Message size limits
- Automatic cleanup of stale connections

## Error Handling

### Connection Errors
```json
{
  "type": "error",
  "payload": {
    "message": "Connection limit exceeded",
    "code": "CONNECTION_LIMIT",
    "details": {
      "current_connections": 5,
      "max_allowed": 3
    }
  }
}
```

### Authentication Errors
```json
{
  "type": "auth_error", 
  "payload": {
    "message": "Invalid credentials",
    "code": "INVALID_CREDENTIALS",
    "retry_allowed": true
  }
}
```

### Command Errors
```json
{
  "type": "command_error",
  "payload": {
    "command_id": 123,
    "message": "Command execution failed", 
    "code": "EXECUTION_ERROR",
    "details": {
      "exit_code": 1,
      "stderr": "command not found"
    }
  }
}
```

## Client Implementation Examples

### JavaScript/Browser Client
```javascript
class ControlMeClient {
    constructor(url) {
        this.ws = new WebSocket(url);
        this.authenticated = false;
        this.setupEventHandlers();
    }
    
    setupEventHandlers() {
        this.ws.onopen = () => {
            console.log('Connected to ControlMe server');
            this.sendPing();
        };
        
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            this.handleMessage(message);
        };
        
        this.ws.onclose = () => {
            console.log('Disconnected from server');
            this.attemptReconnect();
        };
    }
    
    authenticate(username, password) {
        this.send('auth_login', {
            username: username,
            password: password
        });
    }
    
    send(type, payload) {
        const message = {
            type: type,
            payload: payload,
            timestamp: new Date().toISOString()
        };
        this.ws.send(JSON.stringify(message));
    }
    
    handleMessage(message) {
        switch(message.type) {
            case 'auth_success':
                this.authenticated = true;
                console.log('Authentication successful');
                break;
            case 'pong':
                console.log('Pong received');
                break;
            case 'command_result':
                console.log('Command result:', message.payload);
                break;
            default:
                console.log('Unknown message type:', message.type);
        }
    }
    
    sendPing() {
        this.send('ping', {});
    }
    
    executeCommand(commandId, parameters = {}) {
        if (!this.authenticated) {
            console.error('Must be authenticated to execute commands');
            return;
        }
        
        this.send('execute_command', {
            command_id: commandId,
            parameters: parameters
        });
    }
}

// Usage
const client = new ControlMeClient('ws://localhost:8080/ws/client');
client.authenticate('testuser', 'password123');
```

### Go Client Example
```go
package main

import (
    "encoding/json"
    "log"
    "github.com/gorilla/websocket"
)

type Client struct {
    conn *websocket.Conn
    authenticated bool
}

type Message struct {
    Type      string      `json:"type"`
    Payload   interface{} `json:"payload"`
    Timestamp string      `json:"timestamp"`
}

func (c *Client) Connect(url string) error {
    var err error
    c.conn, _, err = websocket.DefaultDialer.Dial(url, nil)
    if err != nil {
        return err
    }
    
    go c.readMessages()
    return nil
}

func (c *Client) Send(msgType string, payload interface{}) error {
    msg := Message{
        Type:      msgType,
        Payload:   payload,
        Timestamp: time.Now().Format(time.RFC3339),
    }
    
    return c.conn.WriteJSON(msg)
}

func (c *Client) Authenticate(username, password string) error {
    return c.Send("auth_login", map[string]string{
        "username": username,
        "password": password,
    })
}

func (c *Client) readMessages() {
    for {
        var msg Message
        err := c.conn.ReadJSON(&msg)
        if err != nil {
            log.Printf("Read error: %v", err)
            break
        }
        
        c.handleMessage(msg)
    }
}

func (c *Client) handleMessage(msg Message) {
    switch msg.Type {
    case "auth_success":
        c.authenticated = true
        log.Println("Authentication successful")
    case "pong":
        log.Println("Pong received")
    case "command_result":
        log.Printf("Command result: %+v", msg.Payload)
    default:
        log.Printf("Unknown message type: %s", msg.Type)
    }
}
```

## Testing

### Manual Testing with WebSocket Clients

1. **wscat (Node.js)**
```bash
# Install wscat
npm install -g wscat

# Connect anonymously
wscat -c ws://localhost:8080/ws/client

# Send ping
{"type":"ping","payload":{}}

# Authenticate
{"type":"auth_login","payload":{"username":"testuser","password":"password123"}}
```

2. **Browser Developer Console**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/client');
ws.onmessage = (e) => console.log(JSON.parse(e.data));
ws.send('{"type":"ping","payload":{}}');
```

### Integration Testing

The server includes comprehensive WebSocket tests in `cmd/tools/test-websocket-auth/` that verify:
- Anonymous connections
- Authentication flows  
- Message routing
- Error handling
- Connection limits
- Session management

## Deployment Considerations

### Production Configuration
- Enable TLS for secure WebSocket connections (wss://)
- Configure appropriate connection limits
- Set up monitoring for WebSocket connections
- Implement proper logging for debugging
- Configure rate limiting and DDoS protection

### Scaling
- WebSocket connections are stateful and tied to specific server instances
- Consider using Redis for pub/sub if scaling across multiple servers
- Implement connection pooling and load balancing strategies
- Monitor memory usage for connection management

### Monitoring
- Track active connection counts
- Monitor authentication success/failure rates  
- Log message processing times
- Alert on connection limit thresholds
- Track WebSocket upgrade success rates