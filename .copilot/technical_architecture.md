# Technical Architecture Summary
## WebSocket Authentication Security Implementation

### Current Architecture
```
Client Application
    ↓ HTTP POST
API Server (Gin Framework)
    ↓ JWT Token Generation
Authentication Service
    ↓ Token Validation
WebSocket Hub
    ↓ Connection Management
WebSocket Clients
```

### Key Components

#### 1. Authentication Flow
- **Registration**: `POST /api/v1/auth/register` → User creation
- **Login**: `POST /api/v1/auth/login` → JWT token issuance
- **WebSocket**: `GET /ws/client` → Token-based connection

#### 2. Security Model
- **JWT Tokens**: HS256 signature, 1-week expiry
- **Token Claims**: user_id, exp, iat
- **WebSocket Auth**: Bearer token in Authorization header
- **Session Management**: One connection per token

#### 3. Database Schema (User Model)
```go
type User struct {
    ID          uuid.UUID
    LoginName   string
    ScreenName  string  
    Email       string
    Password    string  // bcrypt hashed
    Role        string
    RandomOptIn bool
    // ... other fields
}
```

#### 4. Key Files and Their Roles
- `cmd/server/main.go` - Application entry point
- `internal/api/handlers/auth_handlers.go` - Registration/login endpoints
- `internal/api/handlers/websocket_handlers.go` - WebSocket upgrade with auth
- `internal/websocket/hub.go` - Connection management and token tracking
- `internal/auth/auth.go` - JWT token operations
- `internal/services/user_service.go` - User CRUD operations
- `internal/config/config.go` - Application configuration

#### 5. Security Vulnerability Fixed
**Before**: WebSocket accepted user IDs directly
```go
// VULNERABLE: Direct user ID parameter
userID := c.Query("user_id") 
```

**After**: WebSocket requires JWT token validation
```go
// SECURE: Token extraction and validation
token := extractToken(c.GetHeader("Authorization"))
claims, err := validateToken(token)
userID := claims.UserID
```

### Configuration
- **Server Port**: 8080
- **JWT Secret**: From config file
- **JWT Expiry**: 604800 seconds (1 week)
- **Database**: SQLite fallback, PostgreSQL preferred

### Development Environment
- **Hot Reload**: Air configuration in `.air.toml`
- **Testing**: busybox tools for HTTP requests
- **Database**: SQLite file at `data/controlme.db`

### API Endpoints Tested
- `GET /health` - Health check (200 OK)
- `POST /api/v1/auth/register` - User registration (201 Created)
- `POST /api/v1/auth/login` - User authentication (200 OK with token)
- `GET /ws/client` - WebSocket upgrade (101 Switching Protocols with valid token)

### Error Responses
- **401 Unauthorized**: Missing or invalid authentication token
- **400 Bad Request**: Invalid request format or validation errors
- **500 Internal Server Error**: Server-side processing errors