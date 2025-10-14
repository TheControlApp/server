# Integration Test Suite

A comprehensive test suite that validates the ControlApp server API and WebSocket functionality.

## Overview

This integration test suite provides thorough validation of:
- **Server Health**: Basic connectivity and health endpoint
- **User Management**: Registration and authentication flows  
- **WebSocket Communication**: Real-time messaging and connection handling
- **Authentication Methods**: JWT tokens via multiple methods
- **Error Handling**: RFC 7807 compliant error responses

## Features

### 🎯 **Comprehensive Testing**
- Multi-phase test execution with clear progress reporting
- Detailed test results with success/failure indicators
- Error collection and reporting for debugging
- Automatic cleanup of test resources

### 🏗️ **Clean Architecture**
- Modular design with separate concerns
- Robust error handling and recovery
- Configurable test parameters
- Thread-safe message handling

### 📊 **Rich Reporting**
- Real-time progress updates
- Comprehensive final report
- Performance metrics (test duration)
- API validation summary

## Usage

### Prerequisites
- Server running on `localhost:8080`
- Go 1.19+ installed

### Running Tests

```bash
# From the integration test directory
cd cmd/tools/integration-test

# Run the complete test suite
go run main.go client.go

# Alternative: run with verbose output
go run main.go client.go 2>&1 | tee test-results.log
```

### Example Output

```
🚀 ControlApp Server Integration Test Suite
==========================================
Comprehensive validation of server API and WebSocket functionality.
Tests user registration, authentication, and real-time communication.

🔍 Phase 1: Server Health Check
================================
✅ Server is accessible and healthy

👤 Phase 2: User Management Tests
=================================
📝 Testing user: TestClient1
   ✅ Registration successful
   ✅ Login successful - Token acquired
   📋 User ID: 12345678-1234-5678-9abc-123456789012

🔌 Phase 3: WebSocket Connection Tests
=====================================
🔗 Testing WebSocket connection for TestClient1
   ✅ WebSocket connected successfully
   ✅ WebSocket authenticated successfully

📊 Integration Test Report
==========================
⏱️  Test Duration: 2.345s
🖥️  Server Health: WORKING
🎯 Overall Result: WORKING

✅ All tests passed! Server implementation is working correctly.
```

## Test Phases

### Phase 1: Server Health Check
- Verifies server accessibility at configured URL
- Tests the `/health` endpoint
- Ensures basic HTTP connectivity

### Phase 2: User Management Tests
- Tests user registration endpoint
- Validates login functionality
- Verifies JWT token generation
- Handles existing user scenarios gracefully

### Phase 3: WebSocket Connection Tests
- Establishes WebSocket connections
- Tests authentication via query parameters
- Validates message listener functionality
- Verifies connection stability

### Phase 4: Authentication Method Tests
- JWT via query parameters
- JWT via WebSocket messages
- RFC 7807 error response validation
- Multi-method authentication support

### Phase 5: Real-time Communication Tests
- Ping/pong heartbeat functionality
- Message exchange between clients
- Broadcast message distribution
- Connection lifecycle management

## Configuration

The test suite uses sensible defaults but can be customized:

```go
// Default configuration
config := &Config{
    ServerURL: "http://localhost:8080",
    WSBaseURL: "ws://localhost:8080", 
    Timeout:   10 * time.Second,
}
```

## Architecture

### Core Components

- **TestSuite**: Main orchestrator managing test execution
- **TestClient**: Individual client simulation with full API support
- **TestResults**: Comprehensive result tracking and reporting
- **Config**: Centralized configuration management

### File Structure

```
integration-test/
├── main.go          # Test suite orchestration and reporting
├── client.go        # TestClient implementation
└── README.md        # This documentation
```

## Error Handling

The test suite includes robust error handling:

- **Graceful Degradation**: Continues testing even if some phases fail
- **Detailed Error Reporting**: Captures and reports specific failure reasons
- **Resource Cleanup**: Ensures proper cleanup even on failures
- **Timeout Protection**: Prevents hanging on network issues

## Validation Areas

### API Endpoints Tested
- `GET /health` - Server health check
- `POST /api/v1/auth/register` - User registration  
- `POST /api/v1/auth/login` - User authentication
- `GET /ws/client` - WebSocket connection

### WebSocket Features Tested
- Connection establishment
- Authentication methods (query param + message)
- Message routing and delivery
- Ping/pong heartbeat mechanism
- Multi-client communication

### Error Response Testing
- RFC 7807 compliance validation
- Proper HTTP status codes
- Structured error message format
- Developer-friendly error details

## Development Workflow

This tool supports the development workflow by:

1. **Validating Changes**: Run after any API modifications
2. **Documentation Sync**: Ensures docs match implementation
3. **Regression Testing**: Catches breaking changes early
4. **Integration Verification**: Confirms end-to-end functionality

## Troubleshooting

### Common Issues

**Server Not Running**
```
❌ Server health check failed
💡 Tip: Make sure the server is running with: go run cmd/server/main.go
```

**Authentication Failures**
- Check if test users already exist with different passwords
- Verify JWT token generation is working
- Ensure database connectivity

**WebSocket Connection Issues**
- Verify WebSocket endpoint is accessible
- Check for firewall or proxy interference
- Ensure proper authentication token format

### Debug Mode

For detailed debugging, check server logs while running tests:

```bash
# Terminal 1: Run server with debug logging
go run cmd/server/main.go

# Terminal 2: Run integration tests
go run main.go client.go
```

## Contributing

When modifying the integration tests:

1. **Maintain Backwards Compatibility**: Ensure existing functionality still works
2. **Add Comprehensive Logging**: Include detailed progress and error messages  
3. **Update Documentation**: Keep this README current with any changes
4. **Test Thoroughly**: Verify changes work in both success and failure scenarios

## See Also

- [API Documentation](../../../docs-new/api/) - Complete API reference
- [WebSocket Guide](../../../docs-new/api/websocket-api.md) - WebSocket implementation details
- [Error Handling](../../../docs-new/error-handling.md) - RFC 7807 error responses