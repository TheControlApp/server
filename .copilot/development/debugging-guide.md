# Debugging Guide and Troubleshooting

## 🔍 **Quick Debugging Checklist**

When encountering issues, check these in order:

1. **Server Status**: Is the server running and responding?
2. **Configuration**: Is the correct config file being used?
3. **Database**: Is the database file accessible and migrations complete?
4. **Network**: Are the expected ports available and accessible?
5. **Authentication**: Are JWT tokens valid and not expired?
6. **Message Format**: Are WebSocket messages properly formatted?

## 🚨 **Common Issues and Solutions**

### Server Won't Start

#### Issue: Port Already in Use
```
{"level":"fatal","msg":"Failed to start server: listen tcp :8082: bind: Only one usage of each socket address..."}
```

**Solution**:
```powershell
# Check what's using the port
netstat -ano | findstr :8082

# Kill the process using the port
taskkill /PID <process_id> /F

# Or change port in config.test.yaml
server:
  port: 8083
```

#### Issue: Database Connection Failed
```
[error] failed to initialize database, got error failed to connect to `host=localhost user=postgres database=controlme`
```

**Root Cause**: Server trying to connect to PostgreSQL instead of SQLite

**Solution**:
```powershell
# Ensure using test configuration
$env:CONFIG_FILE="config.test.yaml"

# Verify config file exists and has SQLite settings
Get-Content config.test.yaml | Select-String "sqlite"
```

#### Issue: Configuration File Not Found
```
Error: config file not found
```

**Solution**:
```powershell
# Check if config file exists
Test-Path config.test.yaml

# Verify environment variable
echo $env:CONFIG_FILE

# Set correct path
$env:CONFIG_FILE="config.test.yaml"
```

### WebSocket Connection Issues

#### Issue: WebSocket Connection Refused
```javascript
WebSocket connection to 'ws://localhost:8082/ws/client' failed: Error during WebSocket handshake
```

**Debugging Steps**:
```powershell
# 1. Verify server is running
curl http://localhost:8082/health

# 2. Check server logs for WebSocket errors
# Look for connection upgrade messages

# 3. Test with simple connection
# Use browser console or wscat
```

**Common Causes**:
- Server not running on expected port
- CORS issues (check browser dev tools)
- WebSocket endpoint path incorrect
- Authentication headers malformed

#### Issue: WebSocket Authentication Fails
```json
{"type":"auth_error","payload":{"message":"Invalid credentials"}}
```

**Debugging Process**:
1. **Verify Credentials**: Test with REST API login first
2. **Check Token Format**: Ensure JWT token is valid
3. **Validate Message Format**: Ensure proper JSON structure
4. **Check Server Logs**: Look for authentication error details

```javascript
// Debug authentication step by step
// 1. First get valid token from REST API
fetch('http://localhost:8082/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        username: 'testuser',
        password: 'password123'
    })
}).then(r => r.json()).then(data => {
    console.log('Token:', data.token);
    
    // 2. Then use token in WebSocket
    const ws = new WebSocket(`ws://localhost:8082/ws/client?token=${data.token}`);
});
```

### Database Issues

#### Issue: Database Migrations Fail
```
[error] failed to migrate database: table already exists
```

**Solution**:
```powershell
# Delete database file and restart server
Remove-Item data/controlme.db
$env:CONFIG_FILE="config.test.yaml"; .\tmp\server.exe
```

#### Issue: Permission Denied on Database File
```
[error] failed to open database: unable to open database file
```

**Solutions**:
```powershell
# Check file permissions
Get-Acl data/controlme.db

# Ensure data directory exists
New-Item -ItemType Directory -Path "data" -Force

# Check disk space
Get-PSDrive C
```

### API Endpoint Issues

#### Issue: 404 Not Found on API Endpoints
```
GET http://localhost:8082/api/v1/auth/login 404 (Not Found)
```

**Debugging**:
```powershell
# Check if server is running with correct routes
curl http://localhost:8082/health

# Verify Swagger UI shows all routes
start http://localhost:8082/swagger/index.html

# Check server startup logs for route registration
```

#### Issue: JSON Parsing Errors
```json
{"error":"invalid character 'u' looking for beginning of value"}
```

**Common Causes**:
- Malformed JSON in request body
- Missing Content-Type header
- Incorrect data encoding

**Solution**:
```powershell
# Ensure proper headers and JSON format
curl -X POST http://localhost:8082/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{"username":"test","password":"test123"}'
```

## 🔧 **Advanced Debugging Techniques**

### Enable Debug Logging
```yaml
# Add to config.test.yaml
logging:
  level: debug
  format: json
```

**What Debug Logging Shows**:
- WebSocket connection attempts
- Authentication token validation
- Database query execution
- Message routing decisions
- Error stack traces

### Database Query Debugging
```powershell
# Enable GORM debug mode
$env:GORM_DEBUG="true"
$env:CONFIG_FILE="config.test.yaml"; .\tmp\server.exe
```

**What You'll See**:
```
[0.500ms] [rows:1] SELECT * FROM `users` WHERE username = 'testuser' AND deleted_at IS NULL
```

### WebSocket Message Tracing
```javascript
// Enable comprehensive WebSocket logging
const ws = new WebSocket('ws://localhost:8082/ws/client');

ws.addEventListener('open', (event) => {
    console.log('🟢 WebSocket opened', event);
});

ws.addEventListener('close', (event) => {
    console.log('🔴 WebSocket closed', event.code, event.reason);
});

ws.addEventListener('error', (event) => {
    console.error('❌ WebSocket error', event);
});

ws.addEventListener('message', (event) => {
    try {
        const data = JSON.parse(event.data);
        console.log('📨 Message received', data);
    } catch (e) {
        console.error('❌ Invalid JSON received', event.data);
    }
});

// Monitor sent messages
const originalSend = ws.send;
ws.send = function(data) {
    console.log('📤 Sending message', data);
    return originalSend.call(this, data);
};
```

### Performance Debugging

#### Memory Usage Monitoring
```powershell
# Monitor server memory usage
while ($true) {
    $process = Get-Process | Where-Object {$_.ProcessName -eq "server"}
    if ($process) {
        Write-Host "Memory: $([math]::Round($process.WorkingSet/1MB, 2)) MB"
    }
    Start-Sleep 5
}
```

#### Connection Count Monitoring
```javascript
// Monitor WebSocket connection count
setInterval(() => {
    fetch('http://localhost:8082/health')
        .then(r => r.json())
        .then(data => {
            console.log('Health check:', data);
        });
}, 10000);
```

## 🐛 **Debugging Specific Features**

### Authentication Flow Debugging

#### Progressive Authentication
```javascript
// Step-by-step progressive auth debugging
const ws = new WebSocket('ws://localhost:8082/ws/client');

ws.onopen = () => {
    console.log('1. ✅ Connected anonymously');
    
    // Send auth message
    ws.send(JSON.stringify({
        type: 'auth_login',
        payload: {
            username: 'testuser',
            password: 'password123'
        },
        timestamp: new Date().toISOString()
    }));
    console.log('2. 📤 Sent auth_login message');
};

ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('3. 📨 Received response:', message);
    
    if (message.type === 'auth_success') {
        console.log('4. ✅ Authentication successful');
    } else if (message.type === 'auth_error') {
        console.log('4. ❌ Authentication failed:', message.payload.message);
    }
};
```

#### JWT Token Validation
```javascript
// Decode JWT token to check contents
function decodeJWT(token) {
    const parts = token.split('.');
    const payload = JSON.parse(atob(parts[1]));
    console.log('JWT Payload:', payload);
    console.log('Expires:', new Date(payload.exp * 1000));
    console.log('Is Expired:', Date.now() > payload.exp * 1000);
    return payload;
}

// Use with token from login response
const token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...";
decodeJWT(token);
```

### Message Routing Debugging
```javascript
// Test all message types systematically
const messageTypes = [
    { type: 'ping', payload: {} },
    { type: 'auth_login', payload: { username: 'test', password: 'test' } },
    { type: 'invalid_type', payload: {} }
];

messageTypes.forEach((msg, index) => {
    setTimeout(() => {
        console.log(`Sending message ${index + 1}:`, msg);
        ws.send(JSON.stringify(msg));
    }, index * 1000);
});
```

## 🔍 **Log Analysis Techniques**

### Server Log Patterns
Look for these patterns in server logs:

**Successful Operations**:
```
✓ Using SQLite database
✓ Database connection established successfully
{"level":"info","msg":"Starting server on port 8082"}
[GIN] 2025/09/28 - 21:07:37 | 200 | GET "/health"
```

**Error Patterns**:
```
[error] failed to initialize database
{"level":"fatal","msg":"Failed to start server"}
[GIN] 2025/09/28 - 21:07:37 | 401 | POST "/api/v1/auth/login"
```

**WebSocket Activity**:
```
[GIN] 2025/09/28 - 21:07:37 | 101 | GET "/ws/client"
WebSocket connection established for user: testuser
WebSocket connection closed for user: testuser
```

### Client-Side Debug Console
```javascript
// Comprehensive client-side debugging
class DebugWebSocket {
    constructor(url) {
        this.ws = new WebSocket(url);
        this.messageId = 0;
        this.setupLogging();
    }
    
    setupLogging() {
        this.ws.onopen = (event) => {
            console.log(`🟢 [${new Date().toISOString()}] WebSocket opened`, event);
        };
        
        this.ws.onclose = (event) => {
            console.log(`🔴 [${new Date().toISOString()}] WebSocket closed`, {
                code: event.code,
                reason: event.reason,
                wasClean: event.wasClean
            });
        };
        
        this.ws.onerror = (event) => {
            console.error(`❌ [${new Date().toISOString()}] WebSocket error`, event);
        };
        
        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                console.log(`📨 [${new Date().toISOString()}] Received:`, data);
            } catch (e) {
                console.error(`❌ [${new Date().toISOString()}] Invalid JSON:`, event.data);
            }
        };
    }
    
    send(message) {
        const msgWithId = {
            ...message,
            messageId: ++this.messageId,
            timestamp: new Date().toISOString()
        };
        
        console.log(`📤 [${new Date().toISOString()}] Sending:`, msgWithId);
        this.ws.send(JSON.stringify(msgWithId));
    }
}

// Usage
const debugWs = new DebugWebSocket('ws://localhost:8082/ws/client');
debugWs.send({ type: 'ping', payload: {} });
```

## 🛠️ **Development Tools**

### Useful PowerShell Functions
```powershell
# Add to PowerShell profile for quick debugging
function Test-ControlMeServer {
    try {
        $response = Invoke-RestMethod "http://localhost:8082/health" -TimeoutSec 5
        Write-Host "✅ Server is running: $($response.status)" -ForegroundColor Green
    } catch {
        Write-Host "❌ Server is not responding" -ForegroundColor Red
    }
}

function Restart-ControlMeServer {
    # Kill existing server
    Get-Process | Where-Object {$_.ProcessName -eq "server"} | Stop-Process -Force
    
    # Set environment and restart
    $env:CONFIG_FILE="config.test.yaml"
    Start-Process -FilePath ".\tmp\server.exe" -WindowStyle Minimized
    
    Write-Host "🔄 Server restarted" -ForegroundColor Yellow
}

function Show-ControlMeLogs {
    # Monitor server output (if running in background)
    Get-Content ".\server.log" -Wait -Tail 10
}
```

### Browser DevTools Tips
1. **Network Tab**: Monitor WebSocket frames and HTTP requests
2. **Console Tab**: Run JavaScript debugging code
3. **Application Tab**: Check WebSocket connection status
4. **Security Tab**: Verify HTTPS/WSS if using TLS

### VS Code Debugging
```json
// .vscode/launch.json for debugging
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server/main.go",
            "env": {
                "CONFIG_FILE": "config.test.yaml",
                "GIN_MODE": "debug"
            },
            "args": []
        }
    ]
}
```

This debugging guide should help quickly identify and resolve most issues encountered during development or deployment of the ControlMe server.