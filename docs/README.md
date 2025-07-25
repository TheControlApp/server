# TheControlApp Server Documentation

## Overview
Real-time command delivery platform for desktop applications. Users send commands with instructions to other users via WebSocket connections.

## Quick Start
```bash
docker-compose up
```
Server: http://localhost:8080  
WebSocket: ws://localhost:8080/api/ws  

## Core Concepts
- **Commands** - JSON messages with instruction arrays
- **Instructions** - Individual tasks (popup, file download, etc.)
- **Tags** - Content filtering and user preferences
- **WebSocket First** - Primary communication method

## Documentation
- **[REST API](./api/rest.md)** - Authentication and file endpoints
- **[WebSocket API](./api/websocket.md)** - Real-time messaging
- **[Instructions](./api/instructions.md)** - Command instruction types
- **[C# Client](./client/csharp.md)** - Desktop integration guide

## Basic Command Example
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
