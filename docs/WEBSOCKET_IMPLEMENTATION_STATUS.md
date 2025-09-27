# WebSocket Implementation Status

This document tracks the current implementation status of WebSocket message handling in the ControlMe server.

## ✅ Completed Features

- [x] **WebSocket Connection Management** - Hub handles client connections
- [x] **JWT Authentication** - Token validation on connection
- [x] **One Session Per Token** - Prevents multiple connections with same token
- [x] **Origin Validation** - Security for allowed domains
- [x] **Connection Cleanup** - Proper disconnection handling
- [x] **Message Envelope Format** - Standard JSON structure defined
- [x] **Client Type Removal** - Simplified to single endpoint
- [x] **Documentation** - Complete message format reference

## 🚧 In Development

- [ ] **Message Routing** - Currently echoes back, needs proper routing
- [ ] **Command Processing** - Parse and execute instruction types
- [ ] **Command Status Updates** - Send status back to command sender
- [ ] **Error Handling** - Proper error message format and routing

## 📋 TODO Implementation Tasks

### 1. Message Router Implementation

**File:** `internal/websocket/hub.go` - `ReadPump()` function

**Current Code:**
```go
// Handle incoming messages (TODO: implement message handling)
logrus.WithFields(logrus.Fields{
    "user_id": c.userID,
    "message": string(message),
}).Debug("Received WebSocket message")

// Echo message back for now (TODO: implement proper message routing)
c.send <- message
```

**Needs:**
- Parse incoming JSON messages
- Route messages based on type
- Validate message format
- Handle different message types appropriately

### 2. Command Processing

**Files:** 
- `internal/websocket/hub.go` - Add command processing methods
- `internal/services/command_service.go` - Integrate with existing command service

**Needs:**
- Process `command_completion` messages from clients
- Update command status in database
- Notify command senders of status changes
- Handle command validation and error responses

### 3. Message Types to Implement

#### Client-to-Server Messages:
- `command_completion` - Client reports command status
- `ping` - Heartbeat/keepalive messages
- `status_update` - Client availability status

#### Server-to-Client Messages:
- `command_assignment` - Send new commands to clients
- `command_status` - Notify senders of command updates
- `error` - Error responses
- `pong` - Heartbeat responses

### 4. Integration Points

**Command Service Integration:**
- Connect WebSocket command completion to `CommandService.CompleteCommand()`
- Send new commands via WebSocket when assigned
- Update command status in real-time

**User Service Integration:**
- User availability status
- User preference handling
- Permission validation

### 5. Error Handling

**Error Types to Implement:**
- `AUTHENTICATION_FAILED` - Invalid JWT token
- `INVALID_MESSAGE_FORMAT` - Malformed JSON
- `COMMAND_NOT_FOUND` - Invalid command ID
- `PERMISSION_DENIED` - User lacks permission
- `RATE_LIMITED` - Too many messages
- `SERVER_ERROR` - Internal errors

## 🎯 Implementation Priority

1. **High Priority** - Message routing and basic command completion
2. **Medium Priority** - Command assignment via WebSocket
3. **Low Priority** - Advanced features like status updates and heartbeat

## 🧪 Testing Requirements

### Unit Tests Needed:
- [ ] Message parsing and validation
- [ ] Command completion flow
- [ ] Error handling scenarios
- [ ] Authentication edge cases

### Integration Tests Needed:
- [ ] End-to-end command flow (REST API → WebSocket → Completion)
- [ ] Multiple client scenarios
- [ ] Error recovery and reconnection
- [ ] Token expiration handling

## 📊 Current Message Flow

```
[REST API] → [Database] ← [WebSocket Client]
     ↓            ↑              ↓
  Creates      Reads         Completes
  Command    Command        Command
     ↓            ↑              ↓
     └──────── NOT CONNECTED ────┘
```

**Target Message Flow:**
```
[REST API] → [Database] → [WebSocket Hub] → [WebSocket Client]
     ↓            ↑              ↓              ↓
  Creates      Updates      Notifies       Completes
  Command    Command        Client         Command
     ↓            ↑              ↑              ↓
     └────────────┴──────────────┴──────────────┘
```

## 🔧 Implementation Notes

### Message Size Limits
- Client-to-server: 512 bytes (configured in ReadPump)
- Server-to-client: No limit (but should be reasonable)

### Connection Limits
- One connection per JWT token
- Configurable max connections per user (currently unlimited)

### Performance Considerations
- Message caching for broadcasts
- Connection pooling
- Rate limiting per user

### Security Considerations
- All messages must include valid user IDs
- Command permissions must be validated
- Origin validation for WebSocket connections
- Input sanitization for all message content