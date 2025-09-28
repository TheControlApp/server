# TheControlApp Server Documentation

## Overview
Real-time command delivery platform for desktop applications. Users send commands with instructions to other users via WebSocket connections.

## Quick Start

### 🚀 Development Setup
```bash
# Install dependencies
go mod download

# Run with hot reload
air

# Or run directly  
go run cmd/server/main.go
```

### 🐳 Production Setup
```bash
docker-compose up
```

**Server:** http://localhost:8080  
**WebSocket:** ws://localhost:8080/ws/client  
**API Docs:** http://localhost:8080/swagger/index.html  

## 📚 **Documentation**

### **Complete API Reference**
- **[COMPLETE_API_REFERENCE.md](COMPLETE_API_REFERENCE.md)** - Full REST + WebSocket API documentation
- **[WEBSOCKET_IMPLEMENTATION.md](WEBSOCKET_IMPLEMENTATION.md)** - WebSocket implementation guide for developers

### **WebSocket API** 
- **[WEBSOCKET_API.md](WEBSOCKET_API.md)** - Complete WebSocket API reference
- **[WEBSOCKET_QUICK_REF.md](WEBSOCKET_QUICK_REF.md)** - Quick lookup guide
- **[WEBSOCKET_STATUS.md](WEBSOCKET_STATUS.md)** - Implementation status  
- **[WEBSOCKET_COMPLETE.md](WEBSOCKET_COMPLETE.md)** - Full implementation summary

### **Other APIs**
- **[REST API](api/rest.md)** - HTTP endpoints documentation  
- **[Database Schema](database/schema.md)** - Data models and relationships
- **[API Swagger](API_SWAGGER.md)** - OpenAPI documentation

## Core Concepts
- **Commands** - JSON messages with instruction arrays sent via WebSocket
- **Instructions** - Individual actions within commands (popup, timer, download, etc.)
- **Broadcasting** - Commands sent to all connected clients
- **Progressive Authentication** - Connect anonymously, authenticate later
- **Instructions** - Individual tasks (popup, file download, etc.)
- **Tags** - Content filtering and user preferences
- **WebSocket First** - Primary communication method

## Documentation

### 📚 API References
- **[REST API](./api/rest.md)** - Authentication, users, and command endpoints with code examples
- **[WebSocket API](./api/websocket.md)** - Real-time messaging protocol and client implementation guide
- **[API Errors](./api/errors.md)** - Error codes and handling
- **[Swagger Documentation](./swagger/)** - Interactive API explorer

### 🛠️ Client Development  
- **[Standard Command Types](./standards/command_types.md)** - Official command specifications for client compatibility
- **[Mini Client Example](../client/mini-client.html)** - Complete working HTML/JavaScript client
- **[Implementation Examples](./examples/)** - Sample requests and responses

### 🏗️ Server Implementation
- **[Database Schema](./database/schema.md)** - Data models and relationships  
- **[Instructions](./api/instructions.md)** - Legacy instruction types
- **[WebSocket Hub](../internal/websocket/hub.go)** - Connection management implementation

## Basic Command Example

### Standard Command (New Format)
```json
{
  "type": "command",
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2025-09-16T12:00:00Z",
  "from": "sender_user_id",
  "to": "receiver_user_id", 
  "data": {
    "command_type": "std_popup_text",
    "payload": {
      "title": "Hello!",
      "message": "This is a standard popup message",
      "buttons": [
        {
          "id": "ok",
          "text": "OK", 
          "style": "primary"
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

### Legacy Command (Deprecated)
```json
{
  "type": "send_command", 
  "data": {
    "receiver": "username",
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "Hello!",
          "button": "OK"
        }
      }
    ]
  }
}
```

See [Standard Command Types](./standards/command_types.md) for complete specifications.
