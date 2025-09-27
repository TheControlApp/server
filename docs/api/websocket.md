# WebSocket API

**🔗 Code References:**
- [WebSocket Hub Implementation](../../internal/websocket/hub.go) - Connection management and message routing
- [WebSocket Handlers](../../internal/api/handlers/websocket_handlers.go) - Connection setup and authentication
- [API Routes](../../internal/api/routes/routes.go) - WebSocket endpoint definitions
- [Mini Client Example](../../client/mini-client.html) - Complete working implementation

## Connection

**Endpoint:** (defined in [routes.go](../../internal/api/routes/routes.go))
```
ws://localhost:8080/ws/client   (single endpoint for all clients)
```

**Authentication Methods:**
1. **Authorization Header** (preferred):
   ```
   Authorization: Bearer <jwt_token>
   ```
2. **Query Parameter** (fallback):
   ```
   ?token=<jwt_token>
   ```

**Implementation:** `HandleClientWebSocket()` in [websocket_handlers.go](../../internal/api/handlers/websocket_handlers.go)

**Client Type Detection:** Server automatically detects client type (web/desktop) based on User-Agent header

**Security Features:**
- Origin validation for allowed domains (`localhost:3000`, `localhost:8080`, etc.)
- JWT token validation using the auth service
- One connection per JWT token enforcement via `tokenConnections` map
- Automatic client cleanup on disconnect

## Message Format

All messages use the standard JSON envelope format defined in [hub.go](../../internal/websocket/hub.go):

```go
type Message struct {
	Type      string      `json:"type"`
	ID        string      `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	From      uuid.UUID   `json:"from"`
	To        uuid.UUID   `json:"to"`
	Data      interface{} `json:"data"`
}
```

**JSON Example:**
```json
{
  "type": "command",
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "timestamp": "2025-09-16T12:00:00Z",
  "from": "sender_user_id",
  "to": "receiver_user_id",
  "data": {
    "command_type": "std_popup_text",
    "payload": { /* command-specific data */ }
  }
}
```

## Message Types

### 1. Command Messages (Server → Client)

Server sends commands to be executed by the client:

```json
{
  "type": "command",
  "id": "550e8400-e29b-41d4-a716-446655440001",
  "timestamp": "2025-09-16T12:00:00Z",
  "from": "sender_user_id",
  "to": "receiver_user_id",
  "data": {
    "command_type": "std_popup_text",
    "payload": {
      "title": "Daily Check-in",
      "message": "How are you feeling today?",
      "buttons": [
        {
          "id": "good",
          "text": "Good",
          "style": "primary"
        },
        {
          "id": "bad",
          "text": "Not Great",
          "style": "secondary"
        }
      ]
    },
    "metadata": {
      "priority": "normal",
      "expires_at": "2025-09-16T13:00:00Z"
    }
  }
}
```

### 2. Acknowledgment Messages (Client → Server)

Client confirms command completion:

```json
{
  "type": "ack",
  "id": "response_uuid",
  "timestamp": "2025-09-16T12:00:30Z",
  "from": "receiver_user_id",
  "to": "sender_user_id",
  "data": {
    "original_message_id": "550e8400-e29b-41d4-a716-446655440001",
    "status": "completed",
    "response": {
      "action": "good",
      "completion_time": "2025-09-16T12:00:25Z"
    }
  }
}
```

### 3. Heartbeat Messages

Connection keep-alive (handled automatically by the hub):

```json
{
  "type": "ping",
  "id": "ping_uuid",
  "timestamp": "2025-09-16T12:00:00Z"
}

{
  "type": "pong", 
  "id": "pong_uuid",
  "timestamp": "2025-09-16T12:00:00Z"
}
```

### 4. Error Messages

Server error notifications:

```json
{
  "type": "error",
  "id": "error_uuid",
  "timestamp": "2025-09-16T12:00:00Z",
  "data": {
    "error_code": "authentication_failed",
    "message": "Invalid or expired JWT token",
    "details": {}
  }
}
```

## Standard Command Types

**Note:** These are proposed standard command types for client compatibility. The actual command handling depends on client implementation.

### std_popup_url
Opens a URL in the default browser:
```json
{
  "command_type": "std_popup_url",
  "payload": {
    "url": "https://example.com",
    "title": "Check this out!",
    "description": "Optional link description"
  }
}
```

### std_popup_video  
Displays a video in a popup player:
```json
{
  "command_type": "std_popup_video",
  "payload": {
    "video_source": {
      "type": "url",
      "url": "https://example.com/video.mp4"
    },
    "title": "Video Title",
    "controls": {
      "autoplay": false,
      "loop": false,
      "volume": 0.8
    }
  }
}
```

### std_popup_text
Shows a text dialog with optional buttons:
```json
{
  "command_type": "std_popup_text", 
  "payload": {
    "title": "Confirmation",
    "message": "Are you ready to proceed?",
    "buttons": [
      {
        "id": "yes",
        "text": "Yes, Continue",
        "style": "primary"
      },
      {
        "id": "no",
        "text": "Not Yet",
        "style": "secondary"
      }
    ]
  }
}
```

### std_download_file
Downloads a file via API reference:
```json
{
  "command_type": "std_download_file",
  "payload": {
    "file_reference": {
      "file_id": "file_uuid",
      "access_token": "temp_token",
      "api_endpoint": "/api/files/download"
    },
    "filename": "document.pdf",
    "file_size": 1024000,
    "mime_type": "application/pdf"
  }
}
```

## Error Handling

### Common Error Codes
- `authentication_failed` - Invalid/expired JWT token
- `user_not_found` - Target user does not exist  
- `rate_limit_exceeded` - Too many messages sent
- `invalid_message_format` - Malformed JSON message
- `command_not_supported` - Unknown command type

### Connection Issues
- **Invalid Token:** Connection immediately closed with error
- **Origin Validation:** Non-allowed origins are rejected
- **Duplicate Connection:** New connection replaces existing one for same token

## Connection Management

### Authentication Flow
1. Connect to WebSocket endpoint with JWT token
2. Server validates token and user credentials
3. Client registered in hub with user ID and token
4. Begin message exchange

### Connection Limits
- **One connection per JWT token** (enforced by hub)
- **Automatic cleanup** of stale connections  
- **Origin validation** for security

### Heartbeat System
Connection health is maintained automatically by the WebSocket implementation. No manual ping/pong required.

## 🛠️ Client Implementation Guide

### Basic Connection Setup

#### JavaScript Implementation
```javascript
class ControlMeClient {
  constructor(token, serverUrl = 'ws://localhost:8080') {
    this.token = token;
    this.serverUrl = serverUrl;
    this.ws = null;
    this.userId = null;
  }
  
  connect() {
    // Use Authorization header for token
    this.ws = new WebSocket(`${this.serverUrl}/ws/client`, [], {
      headers: {
        'Authorization': `Bearer ${this.token}`
      }
    });
    
    // Or use query parameter as fallback
    // this.ws = new WebSocket(`${this.serverUrl}/ws/client?token=${this.token}`);
    
    this.ws.onopen = () => {
      console.log('Connected to ControlMe server');
    };
    
    this.ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };
    
    this.ws.onclose = (event) => {
      console.log('Disconnected from server, code:', event.code);
    };
    
    this.ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };
  }
  
  handleMessage(message) {
    switch(message.type) {
      case 'command':
        this.executeCommand(message);
        break;
      case 'error':
        console.error('Server error:', message.data);
        break;
      default:
        console.log('Unknown message type:', message.type);
    }
  }
  
  executeCommand(message) {
    const { command_type, payload } = message.data;
    
    switch(command_type) {
      case 'std_popup_text':
        this.showTextPopup(payload, message.id);
        break;
      case 'std_popup_url':
        window.open(payload.url, '_blank');
        this.sendAck(message.id, 'completed', { action: 'opened' });
        break;
      case 'std_popup_video':
        this.showVideoPopup(payload, message.id);
        break;
      case 'std_download_file':
        this.initiateDownload(payload, message.id);
        break;
      default:
        console.log(`Unsupported command type: ${command_type}`);
        this.sendAck(message.id, 'failed', { error: 'Unsupported command' });
    }
  }
  
  showTextPopup(payload, messageId) {
    const result = confirm(`${payload.title}\n\n${payload.message}`);
    this.sendAck(messageId, 'completed', { 
      action: result ? 'accepted' : 'declined' 
    });
  }
  
  sendAck(originalId, status, response = {}) {
    const ack = {
      type: 'ack',
      id: this.generateUUID(),
      timestamp: new Date().toISOString(),
      from: this.userId,
      to: null, // Server will route correctly
      data: {
        original_message_id: originalId,
        status: status,
        response: {
          ...response,
          completion_time: new Date().toISOString()
        }
      }
    };
    
    this.ws.send(JSON.stringify(ack));
  }
  
  generateUUID() {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
      const r = Math.random() * 16 | 0;
      const v = c == 'x' ? r : (r & 0x3 | 0x8);
      return v.toString(16);
    });
  }
}

// Usage
const client = new ControlMeClient('your_jwt_token_here');
client.connect();
```

#### Python Implementation
```python
import asyncio
import websockets
import json
import uuid
from datetime import datetime

class ControlMeClient:
    def __init__(self, token, server_url="ws://localhost:8080"):
        self.token = token
        self.server_url = server_url
        self.websocket = None
        self.user_id = None
    
    async def connect(self):
        uri = f"{self.server_url}/ws/client"
        headers = {"Authorization": f"Bearer {self.token}"}
        
        try:
            self.websocket = await websockets.connect(uri, extra_headers=headers)
            print("Connected to ControlMe server")
            await self.listen()
        except Exception as e:
            print(f"Connection error: {e}")
    
    async def listen(self):
        try:
            async for message in self.websocket:
                data = json.loads(message)
                await self.handle_message(data)
        except websockets.exceptions.ConnectionClosed:
            print("Connection closed")
        except Exception as e:
            print(f"Error listening: {e}")
    
    async def handle_message(self, message):
        if message['type'] == 'command':
            await self.execute_command(message)
        elif message['type'] == 'error':
            print(f"Server error: {message['data']}")
    
    async def execute_command(self, message):
        command_type = message['data']['command_type']
        payload = message['data']['payload']
        
        try:
            if command_type == 'std_popup_url':
                import webbrowser
                webbrowser.open(payload['url'])
                await self.send_ack(message['id'], 'completed', {'action': 'opened'})
            
            elif command_type == 'std_popup_text':
                print(f"{payload['title']}: {payload['message']}")
                response = input("Response (y/n): ").lower() == 'y'
                await self.send_ack(message['id'], 'completed', {
                    'action': 'accepted' if response else 'declined'
                })
            
            else:
                print(f"Unsupported command type: {command_type}")
                await self.send_ack(message['id'], 'failed', {
                    'error': 'Unsupported command type'
                })
                
        except Exception as e:
            await self.send_ack(message['id'], 'failed', {'error': str(e)})
    
    async def send_ack(self, original_id, status, response=None):
        ack_message = {
            'type': 'ack',
            'id': str(uuid.uuid4()),
            'timestamp': datetime.utcnow().isoformat() + 'Z',
            'from': self.user_id,
            'to': None,
            'data': {
                'original_message_id': original_id,
                'status': status,
                'response': {
                    **(response or {}),
                    'completion_time': datetime.utcnow().isoformat() + 'Z'
                }
            }
        }
        
        await self.websocket.send(json.dumps(ack_message))

# Usage
async def main():
    client = ControlMeClient("your_jwt_token_here")
    await client.connect()

if __name__ == "__main__":
    asyncio.run(main())
```

### Authentication Integration

1. **Get JWT Token:** Use [REST API login](rest.md#authentication)
2. **Connect WebSocket:** Pass token via Authorization header or query parameter
3. **Handle Messages:** Process incoming commands and send acknowledgments

### Message Flow Example

1. **Client connects** with JWT token
2. **Server sends command:**
   ```json
   {
     "type": "command",
     "id": "cmd-123",
     "timestamp": "2025-09-16T12:00:00Z",
     "from": "sender_user_id",
     "to": "your_user_id",
     "data": {
       "command_type": "std_popup_text",
       "payload": {
         "title": "Hello", 
         "message": "Click OK to continue"
       }
     }
   }
   ```
3. **Client executes** command and shows popup
4. **Client sends acknowledgment:**
   ```json
   {
     "type": "ack",
     "id": "ack-456", 
     "timestamp": "2025-09-16T12:00:05Z",
     "from": "your_user_id",
     "to": "sender_user_id",
     "data": {
       "original_message_id": "cmd-123",
       "status": "completed",
       "response": {
         "action": "accepted",
         "completion_time": "2025-09-16T12:00:05Z"
       }
     }
   }
   ```

For complete working examples, see [Mini Client](../../client/mini-client.html) which demonstrates all functionality in a single HTML file.

## Summary

The ControlMe WebSocket API provides real-time bidirectional communication using:

- **JWT Authentication:** Required for all connections
- **Standard Message Format:** JSON envelope with type, id, timestamp, from, to, data
- **Command System:** Server sends commands, client sends acknowledgments  
- **Standard Command Types:** URL, video, text, and file download commands
- **Error Handling:** Comprehensive error codes and connection management
- **Security:** Origin validation, token enforcement, automatic cleanup

See the [Mini Client](../../client/mini-client.html) for a complete working implementation.