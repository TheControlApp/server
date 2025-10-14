# Integration Test Client

This directory contains a comprehensive integration test client written in Go that validates the documentation accuracy by testing the actual server implementation.

## What It Tests

The integration test verifies:

### ✅ **REST API Endpoints**
- User registration (`POST /api/v1/auth/register`)
- User login (`POST /api/v1/auth/login`) 
- JWT token generation and validation
- RFC 7807 compliant error responses

### ✅ **WebSocket Functionality**
- WebSocket connection establishment (`/ws/client`)
- Authentication via query parameter
- Authentication via message (`auth` type)
- Ping/pong heartbeat functionality
- Message routing and delivery
- Connection management

### ✅ **Documentation Accuracy**
- Validates that documented APIs actually work
- Tests authentication methods described in docs
- Verifies error response formats match RFC 7807
- Confirms WebSocket message formats are correct

## How It Works

1. **User Setup**: Creates or uses existing test users (`test1`, `test2`)
2. **Authentication**: Obtains JWT tokens via REST API
3. **WebSocket Connections**: Establishes two WebSocket connections
4. **Message Exchange**: Tests various message types and routing
5. **Results Analysis**: Provides detailed validation report

## Running the Test

### Prerequisites

1. **Server Running**: Make sure the ControlApp server is running on `localhost:8080`
   ```bash
   # From the server root directory
   go run cmd/server/main.go
   # OR
   air  # if using air for hot reload
   ```

2. **Go Environment**: Go 1.19+ installed

### Run the Integration Test

```bash
# Navigate to the examples directory
cd docs-new/examples

# Download dependencies
go mod tidy

# Run the integration test
go run integration-test.go
```

### Expected Output

```
🚀 Starting ControlApp Integration Test
======================================

📋 Phase 1: User Setup and Authentication
[15:04:05.123][TestClient1] Setting up user...
[15:04:05.234][TestClient1] ✅ User setup complete - Token: eyJhbGciOiJIUzI1NiIs...
[15:04:05.345][TestClient2] Setting up user...
[15:04:05.456][TestClient2] ✅ User setup complete - Token: eyJhbGciOiJIUzI1NiIs...

🔌 Phase 2: WebSocket Connections
[15:04:05.567][TestClient1] Connecting to WebSocket...
[15:04:05.678][TestClient1] ✅ WebSocket connected
[15:04:05.789][TestClient2] Connecting to WebSocket...
[15:04:05.890][TestClient2] ✅ WebSocket connected

🔐 Phase 3: WebSocket Authentication
[15:04:05.991][TestClient1] Authenticating WebSocket connection...
[15:04:06.102][TestClient1] ✅ Authentication message sent
[15:04:06.213][TestClient2] Authenticating WebSocket connection...
[15:04:06.324][TestClient2] ✅ Authentication message sent

💬 Phase 4: Command Exchange Test
🏓 Test 1: Ping/Pong Exchange
📋 Test 2: Command Message Exchange  
📢 Test 3: Broadcast Message Test

📊 Phase 5: Results Analysis
📈 Test Results Summary
======================

👤 TestClient1 (test1)
   📧 Total messages received: 3
   📊 Message types:
      - auth_success: 1
      - pong: 1
      - command: 1
   ✅ WebSocket authentication: SUCCESS
   ✅ Ping/Pong functionality: SUCCESS

👤 TestClient2 (test2)
   📧 Total messages received: 3
   📊 Message types:
      - auth_success: 1
      - pong: 1
      - broadcast: 1
   ✅ WebSocket authentication: SUCCESS
   ✅ Ping/Pong functionality: SUCCESS

🎯 Overall Assessment
====================
✅ TestClient1: JWT authentication successful
✅ TestClient1: WebSocket connection successful
✅ TestClient2: JWT authentication successful
✅ TestClient2: WebSocket connection successful

🎉 ALL TESTS PASSED! Documentation is accurate.

📋 Documentation Validation Results
===================================
✅ REST API /api/v1/auth/register: WORKING
✅ REST API /api/v1/auth/login: WORKING
✅ WebSocket connection /ws/client: WORKING
✅ WebSocket auth via query param: WORKING
✅ WebSocket auth via message: WORKING
✅ RFC 7807 error responses: WORKING

⏳ Test running... Press Ctrl+C to stop
```

## Test Features

### 🔄 **Adaptive User Management**
- Attempts login first (handles existing users)
- Falls back to registration if user doesn't exist
- Handles duplicate username conflicts gracefully
- Uses default credentials for consistency

### 🔗 **Comprehensive WebSocket Testing**
- Tests both authentication methods (query param + message)
- Validates ping/pong heartbeat functionality
- Simulates real command exchange between clients
- Tests broadcast message functionality

### 📊 **Detailed Analysis**
- Counts and categorizes all received messages
- Validates authentication success
- Checks heartbeat functionality
- Provides overall pass/fail assessment

### 🛡️ **Error Handling Validation**
- Tests RFC 7807 error response format
- Handles network failures gracefully
- Validates authentication errors
- Tests malformed request handling

## Troubleshooting

### Common Issues

**Connection Refused**
```
❌ Failed to setup user TestClient1: Post "http://localhost:8080/api/v1/auth/login": dial tcp [::1]:8080: connectex: No connection could be made
```
**Solution**: Make sure the server is running on port 8080

**Authentication Failures**
```
❌ Login failed, attempting registration: unauthorized: Invalid username or password
```
**Solution**: This is normal - the test will try registration and then login

**WebSocket Connection Issues**
```
❌ WebSocket error: websocket: bad handshake
```
**Solution**: Check that WebSocket endpoint `/ws/client` is available

### Debug Mode

To see more detailed output, you can modify the test to include debug logging:

```go
// Add this near the top of main()
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## Files in This Directory

- `integration-test.go` - Main integration test client
- `go.mod` - Go module definition with dependencies
- `websocket-client.js` - JavaScript client example (browser/Node.js)
- `README.md` - This documentation

## Integration with Documentation

This test validates the examples and claims made in:

- `../api/rest-api.md` - REST API endpoints and responses
- `../api/websocket-api.md` - WebSocket connection and authentication
- `../api/authentication.md` - JWT token handling
- `../error-handling.md` - RFC 7807 error response format
- `../MIGRATION_GUIDE.md` - Breaking changes and migration info

The test serves as living documentation that ensures all examples and API descriptions remain accurate as the server implementation evolves.