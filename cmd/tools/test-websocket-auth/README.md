# WebSocket Authentication Test Client

This is a test CLI client that demonstrates the new WebSocket authentication flow.

## Features

- **Anonymous Connection**: Connects to WebSocket without authentication
- **Progressive Authentication**: Can authenticate after connection using JWT tokens
- **Ping/Pong**: Test basic connectivity
- **Interactive CLI**: Easy command interface

## Usage

1. **Start the server** (make sure it's running on localhost:8080)

2. **Run the test client**:
   ```bash
   go run cmd/tools/test-websocket-auth/main.go
   ```

3. **Available Commands**:
   - `ping` - Send a ping message to test connectivity
   - `auth <token>` - Authenticate with a JWT token
   - `quit` - Exit the client

## Authentication Flow

### Phase 1: Anonymous Connection
1. Client connects to WebSocket endpoint without any authentication
2. Connection is established in anonymous mode
3. Client can send public messages like `ping`

### Phase 2: Progressive Authentication (Optional)
1. Client sends `auth` message with JWT token
2. Server validates the token
3. If valid, client is upgraded to authenticated status
4. Client can now send authenticated messages

## Example Session

```
WebSocket Authentication Test Client
=====================================
Connecting to ws://localhost:8080/ws...
✅ Connected successfully (anonymous session)

Available commands:
  ping           - Send a ping message
  auth <token>   - Authenticate with JWT token
  quit           - Exit the client

> ping
🏓 Ping sent
🏓 Pong received (timestamp: 1703123456)

> auth eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
🔐 Authentication request sent with token: eyJhbGciOiJIUzI1NiIsInR5cCI...
✅ Authentication successful! Authentication successful

> quit
Goodbye!
```

## Getting JWT Tokens

To get JWT tokens for testing, you can:

1. **Use the REST API** to register/login:
   ```bash
   # Register a new user
   curl -X POST http://localhost:8080/api/auth/register \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'

   # Login to get token
   curl -X POST http://localhost:8080/api/auth/login \
     -H "Content-Type: application/json" \
     -d '{"email":"test@example.com","password":"password123"}'
   ```

2. **Use the create-test-user tool** (if available):
   ```bash
   go run cmd/tools/create-test-user/main.go
   ```

## Message Format

All WebSocket messages use JSON format:

### Ping Message
```json
{
  "type": "ping",
  "timestamp": 1703123456
}
```

### Auth Message
```json
{
  "type": "auth",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Response Messages
```json
{
  "type": "pong",
  "timestamp": 1703123456
}
```

```json
{
  "type": "auth_success",
  "message": "Authentication successful",
  "user_id": "123e4567-e89b-12d3-a456-426614174000"
}
```

```json
{
  "type": "error",
  "message": "Authentication required for this message type"
}
```