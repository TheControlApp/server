# Complete API Reference

> **Interactive Docs**: [Swagger UI](http://localhost:8080/swagger/index.html) for hands-on testing

## 📚 Table of Contents
- [REST API](#rest-api)
- [WebSocket API](#websocket-api)  
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Examples](#examples)

---

## 🌐 REST API

**Base URL**: `http://localhost:8080/api/v1`

### Authentication Endpoints

#### `POST /auth/register`
Create a new user account.

**Request:**
```json
{
  "username": "testuser",
  "password": "password123", 
  "screen_name": "Test User",
  "random_opt_in": false
}
```

**Response (201):**
```json
{
  "user": {
    "id": "uuid",
    "login_name": "testuser",
    "screen_name": "Test User",
    "role": "user",
    "created_at": "2025-11-25T12:00:00Z"
  }
}
```

#### `POST /auth/login`
Authenticate user and receive JWT token.

**Request:**
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "login_name": "testuser", 
    "screen_name": "Test User"
  }
}
```

### User Endpoints

#### `GET /users`
List all users in the system.

**Headers:**
- `Authorization: Bearer <token>` (optional for now)

**Response (200):**
```json
{
  "users": [
    {
      "id": "uuid",
      "login_name": "user1",
      "screen_name": "User One",
      "role": "user"
    }
  ]
}
```

#### `GET /users/{id}`
Get specific user by ID.

**Response (200):**
```json
{
  "user": {
    "id": "uuid",
    "login_name": "testuser",
    "screen_name": "Test User",
    "role": "user",
    "created_at": "2025-11-25T12:00:00Z"
  }
}
```

### Command Endpoints

#### `GET /commands/pending`
Get pending commands for a user.

**Parameters:**
- `user_id` (query, required) - UUID of the user

**Example:**
```bash
GET /commands/pending?user_id=123e4567-e89b-12d3-a456-426614174000
```

**Response (200):**
```json
{
  "commands": [
    {
      "id": "uuid",
      "instructions": [
        {
          "type": "std_popup",
          "content": {
            "body": "Hello World!",
            "button": "OK"
          }
        }
      ],
      "sender_id": "uuid",
      "receiver_id": "uuid", 
      "tags": "general,test",
      "status": "pending",
      "created_at": "2025-11-25T12:00:00Z"
    }
  ]
}
```

#### `POST /commands/complete`
Mark a command as completed.

**Parameters:**
- `user_id` (query, required) - UUID of the user completing
- `command_id` (query, required) - UUID of the command

**Response (200):**
```json
{
  "message": "Command completed successfully"
}
```

### System Endpoints

#### `GET /health`
Check server health and status.

**Response (200):**
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

---

## 📡 WebSocket API

**Connection**: `ws://localhost:8080/ws/client`

### Connection
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/client');

// With authentication via query parameter
const ws = new WebSocket('ws://localhost:8080/ws/client?token=jwt_token');
```

### Authentication (Optional)
```javascript
// Authenticate after connection
ws.send(JSON.stringify({
  type: 'auth',
  token: 'your_jwt_token_here'
}));

// Success response
{
  "type": "auth_success",
  "message": "Authentication successful",
  "user_id": "uuid"
}
```

### Sending Commands
```javascript
// Command structure
{
  "instructions": [
    {
      "type": "instruction_type",
      "content": { /* instruction-specific data */ }
    }
  ],
  "tags": "general,category",
  "receiver_id": "uuid_or_null_for_broadcast"
}

// Example: Send popup message
ws.send(JSON.stringify({
  instructions: [{
    type: 'std_popup',
    content: {
      body: 'Hello World!',
      button: 'OK'
    }
  }],
  tags: 'general,greeting',
  receiver_id: null // broadcast to all
}));
```

### Instruction Types

| Type | Description | Content Fields |
|------|-------------|----------------|
| `std_popup` | Show popup dialog | `body`, `button` |
| `std_timer` | Start countdown timer | `duration`, `message` |
| `display_text` | Display text message | `text`, `style` |
| `notification` | Show notification | `title`, `body`, `icon` |
| `open_url` | Open URL in browser | `url`, `target` |
| `download_file` | Download file | `url`, `filename` |
| `form_input` | Request user input | `fields`, `title` |

### System Messages
```javascript
// Ping/Pong
ws.send(JSON.stringify({ type: 'ping' }));
// Response: { "type": "pong", "timestamp": 1234567890 }

// Error response
{
  "type": "error", 
  "message": "Error description"
}
```

---

## 🔐 Authentication

### JWT Tokens
Use JWT tokens for authenticated requests:
```bash
Authorization: Bearer <your_jwt_token>
```

### Token Lifecycle
1. **Register/Login** → Receive JWT token
2. **Include in requests** → Authorization header or WebSocket auth
3. **Token expiration** → Re-authenticate when token expires

**Example Usage:**
```bash
# Get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "pass"}' | jq -r '.token')

# Use token
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/users
```

---

## ❌ Error Handling

All errors follow **RFC 7807 Problem Details** standard.

### Error Response Structure
```json
{
  "type": "error_type",
  "title": "Human Readable Title", 
  "status": 400,
  "detail": "Specific error description",
  "action": "Suggested user action",
  "help": "Additional guidance"
}
```

### Validation Errors (422)
```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation",
  "errors": [
    {
      "field": "username",
      "message": "Username must be at least 3 characters",
      "code": "MIN_LENGTH"
    }
  ]
}
```

### Common Error Types
- `bad_request` (400) - Malformed request
- `unauthorized` (401) - Authentication required  
- `forbidden` (403) - Access denied
- `not_found` (404) - Resource not found
- `conflict` (409) - Resource conflict (e.g., username taken)
- `validation_error` (422) - Input validation failed
- `internal_server_error` (500) - Server error

---

## 🧪 Examples

### Complete Registration Flow
```bash
# 1. Register user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "password": "securepass123",
    "screen_name": "New User"
  }'

# 2. Login to get token  
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser", 
    "password": "securepass123"
  }'
```

### WebSocket Command Example
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/client');

ws.onopen = () => {
  // Send a popup command
  ws.send(JSON.stringify({
    instructions: [{
      type: 'std_popup',
      content: {
        body: 'Welcome to ControlApp!',
        button: 'Got it'
      }
    }],
    tags: 'welcome,onboarding',
    receiver_id: null // broadcast
  }));
};

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};
```

### Error Handling Example
```javascript
fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'x', password: 'y' })
})
.then(response => {
  if (!response.ok) {
    return response.json().then(error => {
      console.error('Error type:', error.type);
      console.error('Details:', error.detail);
      if (error.errors) {
        error.errors.forEach(err => {
          console.error(`Field ${err.field}: ${err.message}`);
        });
      }
    });
  }
  return response.json();
});
```

---

## 🔗 Additional Resources

- **[Swagger UI](http://localhost:8080/swagger/index.html)** - Interactive API testing
- **[Integration Tests](cmd/tools/integration-test/)** - Automated validation tool  
- **[WebSocket Examples](examples/websocket-client.js)** - Client implementation examples

---

**💡 Pro Tip**: Use the Swagger UI for hands-on API exploration and the integration test tool to validate your setup!