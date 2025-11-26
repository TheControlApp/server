# Complete API Reference

## Error Response Format

This API follows the RFC 7807 Problem Details standard for HTTP APIs. All error responses include structured information to help developers understand and handle errors effectively.

### Standard Error Response Structure

```json
{
  "type": "string",           // Error type identifier
  "title": "string",          // Human-readable error title
  "status": 400,              // HTTP status code
  "detail": "string",         // Specific error description
  "action": "string",         // Suggested action (optional)
  "help": "string"            // Additional guidance (optional)
}
```

### Common Error Types

- `bad_request` (400) - Request is malformed or invalid
- `unauthorized` (401) - Authentication required or failed
- `forbidden` (403) - Access denied for authenticated user
- `not_found` (404) - Requested resource not found
- `validation_error` (422) - Input validation failed
- `internal_server_error` (500) - Server encountered an error

### Validation Error Response

Validation errors include detailed field-level information:

```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation",
  "help": "Please check the field requirements in the API documentation",
  "errors": [
    {
      "field": "username",
      "message": "Username must be at least 3 characters long",
      "code": "MIN_LENGTH"
    }
  ]
}
```

## REST API Endpoints

### Authentication

#### POST /api/v1/auth/register
Register a new user account.

**Request Body:**
```json
{
  "username": "string (required, min 3 chars)",
  "password": "string (required, min 6 chars)", 
  "screen_name": "string (required)",
  "random_opt_in": "boolean (optional)"
}
```

**Success Response (201):**
```json
{
  "user": {
    "id": "uuid",
    "login_name": "testuser", 
    "screen_name": "Test User",
    "role": "user",
    "random_opt_in": false,
    "anon_cmd": false,
    "verified": false,
    "created_at": "2023-10-12T18:00:00Z",
    "updated_at": "2023-10-12T18:00:00Z"
  }
}
```

**Error Responses:**
- **400 Bad Request**: Invalid JSON or missing required fields
- **409 Conflict**: Username already exists
- **422 Validation Error**: Field validation failed
- **500 Internal Server Error**: Server error during registration

#### POST /api/v1/auth/login
Authenticate an existing user.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)"
}
```

**Success Response (200):**
```json
{
  "user": {
    "id": "uuid",
    "login_name": "testuser",
    "screen_name": "Test User", 
    "role": "user",
    "random_opt_in": false,
    "anon_cmd": false,
    "verified": false,
    "created_at": "2023-10-12T18:00:00Z",
    "updated_at": "2023-10-12T18:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- **400 Bad Request**: Invalid JSON or missing required fields
- **401 Unauthorized**: Invalid username or password  
- **500 Internal Server Error**: Server error during authentication

### Users

#### GET /api/v1/users
Get all users in the system.

**Success Response (200):**
```json
{
  "users": [
    {
      "id": "uuid",
      "login_name": "testuser",
      "screen_name": "Test User",
      "role": "user", 
      "random_opt_in": false,
      "anon_cmd": false,
      "verified": false,
      "created_at": "2023-10-12T18:00:00Z",
      "updated_at": "2023-10-12T18:00:00Z"
    }
  ]
}
```

**Error Responses:**
- **500 Internal Server Error**: Server error retrieving users

#### GET /api/v1/users/{id}
Get a specific user by ID.

**Path Parameters:**
- `id` (string, required): User UUID

**Success Response (200):**
```json
{
  "user": {
    "id": "uuid",
    "login_name": "testuser",
    "screen_name": "Test User",
    "role": "user",
    "random_opt_in": false,
    "anon_cmd": false, 
    "verified": false,
    "created_at": "2023-10-12T18:00:00Z",
    "updated_at": "2023-10-12T18:00:00Z"
  }
}
```

**Error Responses:**
- **422 Validation Error**: Invalid UUID format
- **404 Not Found**: User not found

### Commands

#### GET /api/v1/commands/pending
Get pending commands for a user.

**Query Parameters:**
- `user_id` (string, required): User UUID

**Success Response (200):**
```json
{
  "commands": [
    {
      "id": "uuid",
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
      "user_id": "uuid",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

**Error Responses:**
- **422 Validation Error**: Invalid user_id format
- **500 Internal Server Error**: Server error retrieving commands

#### POST /api/v1/commands/complete
Mark a command as completed.

**Query Parameters:**
- `user_id` (string, required): User UUID
- `command_id` (string, required): Command UUID

**Success Response (200):**
```json
{
  "message": "Command completed successfully"
}
```

**Error Responses:**
- **422 Validation Error**: Invalid UUID format for user_id or command_id
- **500 Internal Server Error**: Server error completing command

## Health Check

#### GET /health
Check server health status.

**Success Response (200):**
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

## Example Error Responses

### Bad Request (400)
```json
{
  "type": "bad_request",
  "title": "Bad Request", 
  "status": 400,
  "detail": "Request body is not valid JSON or missing required fields"
}
```

### Unauthorized (401)
```json
{
  "type": "unauthorized",
  "title": "Unauthorized Access",
  "status": 401,
  "detail": "Invalid username or password",
  "action": "Please check your credentials and try again"
}
```

### Validation Error (422)
```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation", 
  "help": "Please check the field requirements in the API documentation",
  "errors": [
    {
      "field": "username",
      "message": "Username must be at least 3 characters long",
      "code": "MIN_LENGTH"
    },
    {
      "field": "password", 
      "message": "Password must be at least 6 characters long",
      "code": "MIN_LENGTH"
    }
  ]
}
```

### Conflict (409)
```json
{
  "type": "conflict",
  "title": "Resource Conflict",
  "status": 409,
  "detail": "Username already exists",
  "action": "Please choose a different username",
  "instance": {
    "username": "testuser"
  }
}
```

### Not Found (404)
```json
{
  "type": "not_found",
  "title": "Resource Not Found",
  "status": 404,
  "detail": "User not found"
}
```

### Internal Server Error (500)
```json
{
  "type": "internal_server_error",
  "title": "Internal Server Error", 
  "status": 500,
  "detail": "Database connection failed",
  "action": "Please try again later or contact support if the problem persists"
}
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