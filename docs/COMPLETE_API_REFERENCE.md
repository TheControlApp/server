# Complete API Reference

## REST API Endpoints

### Authentication

#### POST /api/v1/auth/register
Register a new user account.

**Request Body:**
```json
{
  "username": "string (required, min 3 chars)",
  "password": "string (required, min 6 chars)",
  "email": "string (required, valid email)"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### POST /api/v1/auth/login
Authenticate an existing user.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "user": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Commands

#### GET /api/v1/commands
Get all commands for authenticated user.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "success": true,
  "commands": [
    {
      "id": 1,
      "name": "test_command",
      "description": "A test command",
      "instructions": [
        {
          "type": "shell",
          "command": "echo 'Hello World'",
          "args": [],
          "timeout": 30
        }
      ],
      "user_id": 1,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /api/v1/commands
Create a new command.

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "string (required)",
  "description": "string (optional)",
  "instructions": [
    {
      "type": "shell|powershell|cmd|python|node|custom",
      "command": "string (required)",
      "args": ["string", "array", "optional"],
      "timeout": 30,
      "working_directory": "string (optional)",
      "environment": {
        "KEY": "value"
      }
    }
  ]
}
```

### Users

#### GET /api/v1/admin/users
Get all users (admin only).

**Headers:**
```
Authorization: Bearer <admin_token>
```

## WebSocket API

### Connection Endpoint

#### GET /ws/client
Establishes a WebSocket connection for real-time command distribution.

**Connection Options:**

1. **Anonymous Connection:**
   ```
   ws://localhost:8080/ws/client
   ```

2. **Authenticated Connection (Header):**
   ```
   ws://localhost:8080/ws/client
   Headers: Authorization: Bearer <token>
   ```

3. **Authenticated Connection (Query):**
   ```
   ws://localhost:8080/ws/client?token=<jwt_token>
   ```

### Connection Flow

1. **Anonymous Session:**
   - Connect without authentication
   - Receive broadcast messages only
   - Can upgrade to authenticated session via WebSocket messages

2. **Authenticated Session:**
   - Connect with JWT token (header or query parameter)
   - Receive all messages (broadcasts + user-specific)
   - Can send commands and receive responses

3. **Progressive Authentication:**
   - Start with anonymous connection
   - Send login message to upgrade session
   - Receive authentication confirmation

### Message Types

All WebSocket messages follow this format:
```json
{
  "type": "message_type",
  "payload": { /* type-specific data */ },
  "timestamp": "2024-01-01T00:00:00.000Z"
}
```

#### Client → Server Messages

##### Authentication Messages

**Login Request:**
```json
{
  "type": "auth_login",
  "payload": {
    "username": "string",
    "password": "string"
  }
}
```

**Token Authentication:**
```json
{
  "type": "auth_token",
  "payload": {
    "token": "jwt_token_string"
  }
}
```

##### Command Messages

**Ping:**
```json
{
  "type": "ping",
  "payload": {}
}
```

**Execute Command:**
```json
{
  "type": "execute_command",
  "payload": {
    "command_id": 1,
    "parameters": {
      "key": "value"
    }
  }
}
```

#### Server → Client Messages

##### Authentication Responses

**Authentication Success:**
```json
{
  "type": "auth_success",
  "payload": {
    "user": {
      "id": 1,
      "username": "testuser",
      "email": "test@example.com"
    },
    "session_id": "uuid",
    "capabilities": ["execute_commands", "receive_broadcasts"]
  }
}
```

**Authentication Error:**
```json
{
  "type": "auth_error",
  "payload": {
    "message": "Invalid credentials",
    "code": "INVALID_CREDENTIALS"
  }
}
```

##### System Messages

**Pong:**
```json
{
  "type": "pong",
  "payload": {}
}
```

**Connection Status:**
```json
{
  "type": "connection_status",
  "payload": {
    "status": "connected|authenticated|disconnected",
    "session_type": "anonymous|authenticated",
    "connected_at": "2024-01-01T00:00:00.000Z"
  }
}
```

##### Command Messages

**Command Result:**
```json
{
  "type": "command_result",
  "payload": {
    "command_id": 1,
    "status": "success|error|timeout",
    "output": "command output text",
    "error": "error message if any",
    "execution_time": 1.23,
    "exit_code": 0
  }
}
```

**Broadcast Command:**
```json
{
  "type": "broadcast_command",
  "payload": {
    "command": {
      "id": 1,
      "name": "system_update",
      "instructions": [
        {
          "type": "shell",
          "command": "echo 'System update available'",
          "args": [],
          "timeout": 30
        }
      ]
    },
    "sender": "admin",
    "broadcast_time": "2024-01-01T00:00:00.000Z"
  }
}
```

##### Error Messages

**General Error:**
```json
{
  "type": "error",
  "payload": {
    "message": "Error description",
    "code": "ERROR_CODE",
    "details": {
      "additional": "context"
    }
  }
}
```

### Session Management

- **Anonymous Sessions:** Limited to receiving broadcast messages
- **Authenticated Sessions:** Full access to user-specific commands and responses
- **Session Upgrade:** Anonymous sessions can authenticate via WebSocket messages
- **Token Validation:** JWT tokens are validated on connection and periodically refreshed
- **Connection Limits:** Configurable per-user connection limits
- **Heartbeat:** Automatic ping/pong for connection health monitoring

### Error Handling

Common WebSocket error scenarios:

1. **Invalid Token:** Connection rejected or downgraded to anonymous
2. **Expired Token:** Session downgraded, client should re-authenticate
3. **Rate Limiting:** Temporary message rejection with retry-after information
4. **Invalid Message Format:** Error response with format requirements
5. **Permission Denied:** Error response for unauthorized operations

### Best Practices

1. **Always handle authentication errors gracefully**
2. **Implement exponential backoff for reconnections**  
3. **Validate message formats before sending**
4. **Handle connection drops and implement auto-reconnect**
5. **Use ping/pong for connection health monitoring**
6. **Process messages asynchronously to avoid blocking**