# Test WebSocket Client

A simple CLI client for testing WebSocket communication with the ControlMe server.

## Features

- 🔌 Connect to WebSocket server with JWT authentication
- 📤 Send ping messages to test connectivity
- 📊 Send status updates
- 🔧 Send raw JSON messages for testing
- 📨 Pretty-print received messages
- 🎯 Interactive command-line interface

## Usage

### Build the client:
```bash
cd cmd/tools/test-client
go build -o test-client.exe
```

### Run the client:
```bash
# Basic usage with JWT token
./test-client.exe -token=your_jwt_token_here

# Specify server URL and user ID
./test-client.exe -url=ws://localhost:8080 -token=your_jwt_token -user=your_user_id
```

### Available Commands:

Once connected, you can use these commands:

- `ping` - Send a ping message to test connectivity
- `status` - Send a status update message
- `raw` - Send a custom JSON message (you'll be prompted to enter JSON)
- `quit` - Exit the client

### Example Session:

```
$ ./test-client.exe -token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
Generated User ID: 550e8400-e29b-41d4-a716-446655440000
Connecting to: ws://localhost:8080/ws/client?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
🎉 Connected to WebSocket server!
Commands:
  ping     - Send a ping message
  status   - Send status update
  raw      - Send raw JSON message
  quit     - Exit the client

> ping
📤 Sent ping message

📨 Received: pong (at 2025-09-27T16:30:15.123Z)
   {
     "type": "pong",
     "id": "response-12345",
     "timestamp": "2025-09-27T16:30:15.123Z",
     "from": "server",
     "to": "550e8400-e29b-41d4-a716-446655440000",
     "data": {
       "original_sequence": 1695844215
     }
   }

> status
📤 Sent status update

> quit
Goodbye!
Shutting down client...
```

## Message Format

The client sends messages in the standard WebSocket message format:

```json
{
  "type": "ping",
  "id": "unique-message-id",
  "timestamp": "2025-09-27T16:30:15.123Z",
  "from": "user-id",
  "to": "server",
  "data": {
    "client_info": {
      "version": "1.0.0",
      "type": "test-cli",
      "platform": "cli"
    },
    "sequence": 1695844215
  }
}
```

## Testing Different Message Types

### Ping Message:
```bash
> ping
```

### Status Update:
```bash
> status
```

### Custom Message:
```bash
> raw
Enter JSON message: {"type":"custom_test","id":"test-123","timestamp":"2025-09-27T16:30:15.123Z","from":"user-id","to":"server","data":{"test":true}}
```

## Troubleshooting

### Connection Issues:
- Ensure the server is running on the specified URL
- Verify your JWT token is valid and not expired
- Check that the token has proper permissions

### Authentication Errors:
- The client will show HTTP status codes for connection failures
- Check server logs for authentication issues
- Ensure the JWT token is properly formatted

### Message Issues:
- The client validates JSON before sending raw messages
- Check that message format matches expected server schema
- Use the ping command first to verify basic connectivity