# Testing Procedures and Guidelines

## 🧪 **Testing Strategy Overview**

The ControlMe Go server testing approach focuses on **real-world functionality** with emphasis on WebSocket communication, authentication flows, and database operations.

## 🚀 **Quick Testing Workflow**

### 1. Server Startup Test
```powershell
# Navigate to project directory
cd D:\Workspace\github.com\TheControlApp\server

# Set test configuration
$env:CONFIG_FILE="config.test.yaml"

# Build and run server
go build -o tmp/server.exe cmd/server/main.go
.\tmp\server.exe
```

**Expected Output**:
```
2025/09/28 21:07:37 ✓ Using SQLite database
2025/09/28 21:07:37 ✓ Database connection established successfully
2025/09/28 21:07:37 ✅ Database migration completed successfully.
{"level":"info","msg":"Starting server on port 8082","time":"2025-09-28T21:07:37-07:00"}
```

### 2. Basic Connectivity Test
```powershell
# Test health endpoint
curl http://localhost:8082/health

# Test Swagger UI (should load in browser)
start http://localhost:8082/swagger/index.html
```

## 🔌 **WebSocket Testing**

### Using Built-in Test Client
```powershell
# Build test client
go build -o bin/test-websocket-auth.exe cmd/tools/test-websocket-auth/main.go

# Run interactive test client
.\bin\test-websocket-auth.exe
```

**Test Scenarios**:
1. **Anonymous Connection**: Connect without authentication
2. **Ping/Pong**: Test connection health
3. **Progressive Authentication**: Login via WebSocket message
4. **Message Routing**: Verify different message types work

### Manual WebSocket Testing

#### Browser Console Test
```javascript
// Connect anonymously
const ws = new WebSocket('ws://localhost:8082/ws/client');

// Setup message handler
ws.onmessage = (event) => {
    console.log('Received:', JSON.parse(event.data));
};

// Send ping
ws.send(JSON.stringify({
    type: 'ping',
    payload: {},
    timestamp: new Date().toISOString()
}));

// Test authentication
ws.send(JSON.stringify({
    type: 'auth_login',
    payload: {
        username: 'testuser',
        password: 'password123'
    },
    timestamp: new Date().toISOString()
}));
```

#### Using wscat (Node.js tool)
```powershell
# Install wscat (requires Node.js)
npm install -g wscat

# Connect to WebSocket
wscat -c ws://localhost:8082/ws/client

# Send test messages
{"type":"ping","payload":{}}
{"type":"auth_login","payload":{"username":"testuser","password":"password123"}}
```

## 🔐 **Authentication Testing**

### REST API Authentication Tests
```powershell
# Test user registration
curl -X POST http://localhost:8082/api/v1/auth/register `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"password123","email":"test@example.com"}'

# Test user login
curl -X POST http://localhost:8082/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"password123"}'
```

**Expected Response**:
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

### WebSocket Authentication Methods

#### Method 1: Header Authentication
```javascript
const ws = new WebSocket('ws://localhost:8082/ws/client', [], {
    headers: {
        'Authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
    }
});
```

#### Method 2: Query Parameter Authentication
```javascript
const token = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...';
const ws = new WebSocket(`ws://localhost:8082/ws/client?token=${token}`);
```

#### Method 3: Progressive Authentication
```javascript
// Connect anonymously first
const ws = new WebSocket('ws://localhost:8082/ws/client');

// Then authenticate via message
ws.send(JSON.stringify({
    type: 'auth_login',
    payload: {
        username: 'testuser',
        password: 'password123'
    }
}));
```

## 🗄️ **Database Testing**

### SQLite Database Tests
```powershell
# Check database file exists
Test-Path "D:\Workspace\github.com\TheControlApp\server\data\controlme.db"

# View database structure (requires sqlite3 CLI)
sqlite3 data/controlme.db ".schema"

# Check user records
sqlite3 data/controlme.db "SELECT id, username, email FROM users;"
```

### Database Migration Tests
```powershell
# Delete database file to test migration
Remove-Item data/controlme.db -ErrorAction SilentlyContinue

# Start server (should recreate database)
$env:CONFIG_FILE="config.test.yaml"; .\tmp\server.exe
```

**Expected Migration Output**:
```
✓ User table migrated successfully
✓ Tag table migrated successfully
✓ Command table migrated successfully
✓ Block table migrated successfully
✓ Report table migrated successfully
```

## 📡 **API Endpoint Testing**

### Comprehensive API Test Script
```powershell
# Test all major endpoints
$baseUrl = "http://localhost:8082"

# Health check
Invoke-RestMethod "$baseUrl/health"

# Register user
$registerData = @{
    username = "testuser$(Get-Random)"
    password = "password123"
    email = "test$(Get-Random)@example.com"
} | ConvertTo-Json

$response = Invoke-RestMethod -Uri "$baseUrl/api/v1/auth/register" -Method POST -Body $registerData -ContentType "application/json"
$token = $response.token

# Login with created user
$loginData = @{
    username = $response.user.username
    password = "password123"
} | ConvertTo-Json

Invoke-RestMethod -Uri "$baseUrl/api/v1/auth/login" -Method POST -Body $loginData -ContentType "application/json"

# Test authenticated endpoint
$headers = @{ Authorization = "Bearer $token" }
Invoke-RestMethod -Uri "$baseUrl/api/v1/commands/pending" -Headers $headers
```

## 🔍 **Error Testing**

### Authentication Error Tests
```powershell
# Test invalid credentials
curl -X POST http://localhost:8082/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{"username":"invalid","password":"wrong"}'

# Test malformed JSON
curl -X POST http://localhost:8082/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{"username":"test"'  # Invalid JSON

# Test missing fields
curl -X POST http://localhost:8082/api/v1/auth/register `
  -H "Content-Type: application/json" `
  -d '{"username":"test"}'  # Missing password and email
```

### WebSocket Error Tests
```javascript
// Test invalid message format
ws.send('invalid json');

// Test invalid message type
ws.send(JSON.stringify({
    type: 'invalid_type',
    payload: {}
}));

// Test authentication with wrong credentials
ws.send(JSON.stringify({
    type: 'auth_login',
    payload: {
        username: 'nonexistent',
        password: 'wrongpassword'
    }
}));
```

## 🚦 **Load Testing**

### Connection Limit Testing
```javascript
// Test multiple connections from same user
const connections = [];
for (let i = 0; i < 5; i++) {
    const ws = new WebSocket('ws://localhost:8082/ws/client?token=' + jwt_token);
    connections.push(ws);
}

// Should see connection limit enforcement after 3 connections (default limit)
```

### Message Rate Testing
```javascript
// Test rapid message sending
const ws = new WebSocket('ws://localhost:8082/ws/client');
ws.onopen = () => {
    for (let i = 0; i < 100; i++) {
        ws.send(JSON.stringify({
            type: 'ping',
            payload: { sequence: i }
        }));
    }
};
```

## 🔧 **Performance Testing**

### Latency Testing
```javascript
// Measure WebSocket round-trip time
const startTime = Date.now();
ws.send(JSON.stringify({
    type: 'ping',
    payload: { timestamp: startTime }
}));

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === 'pong') {
        const latency = Date.now() - message.payload.timestamp;
        console.log(`Round-trip latency: ${latency}ms`);
    }
};
```

### Memory Usage Testing
```powershell
# Monitor server memory usage during testing
Get-Process | Where-Object {$_.ProcessName -eq "server"} | Select-Object ProcessName, WorkingSet, VirtualMemorySize
```

## 🐛 **Debugging Tests**

### Enable Debug Logging
```yaml
# Add to config.test.yaml
logging:
  level: debug
  format: json
```

### WebSocket Connection Debugging
```javascript
// Comprehensive WebSocket event logging
const ws = new WebSocket('ws://localhost:8082/ws/client');

ws.onopen = (event) => {
    console.log('WebSocket opened:', event);
};

ws.onclose = (event) => {
    console.log('WebSocket closed:', event.code, event.reason);
};

ws.onerror = (event) => {
    console.error('WebSocket error:', event);
};

ws.onmessage = (event) => {
    console.log('Message received:', JSON.parse(event.data));
};
```

## ✅ **Test Validation Checklist**

### Basic Functionality
- [ ] Server starts without errors
- [ ] Database migrations complete successfully
- [ ] Health endpoint responds correctly
- [ ] Swagger UI loads properly

### Authentication
- [ ] User registration works
- [ ] User login returns valid JWT token
- [ ] Invalid credentials are rejected
- [ ] JWT token validation works

### WebSocket
- [ ] Anonymous connections accepted
- [ ] Authenticated connections work
- [ ] Progressive authentication functions
- [ ] Ping/pong messaging works
- [ ] Message routing is correct
- [ ] Connection limits enforced
- [ ] Graceful disconnection handling

### Database
- [ ] User creation and retrieval
- [ ] Command storage and retrieval
- [ ] Foreign key relationships work
- [ ] JSON instruction serialization

### Error Handling
- [ ] Invalid JSON handled gracefully
- [ ] Authentication errors return proper codes
- [ ] WebSocket errors don't crash server
- [ ] Database errors handled properly

## 📊 **Test Reporting**

### Test Results Documentation
When running tests, document:
1. **Test Environment**: OS, Go version, database type
2. **Test Configuration**: Which config file used
3. **Test Results**: Pass/fail status for each test
4. **Performance Metrics**: Response times, connection counts
5. **Error Conditions**: Any unexpected errors encountered

### Regression Testing
Before any code changes:
1. Run full test suite
2. Document current baseline performance
3. Make changes
4. Re-run tests
5. Compare results and investigate any regressions

This comprehensive testing approach ensures the ControlMe server maintains reliability and performance across all supported features.