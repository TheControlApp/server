# System Architecture Overview

## 🏗️ **High-Level Architecture**

The ControlMe Go server implements a **WebSocket-first architecture** with REST API authentication support, designed for real-time command distribution and progressive user authentication.

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Client Apps   │    │   Web Browsers   │    │  Mobile Apps    │
└─────────────────┘    └──────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────────────┐
                    │     Load Balancer       │
                    │     (Future: Nginx)     │
                    └─────────────────────────┘
                                 │
                    ┌─────────────────────────┐
                    │   ControlMe Go Server   │
                    │   (Gin + WebSocket)     │
                    └─────────────────────────┘
                                 │
                    ┌─────────────────────────┐
                    │      Database           │
                    │  (SQLite/PostgreSQL)    │
                    └─────────────────────────┘
```

## 🔌 **WebSocket-First Design**

### Core Philosophy
- **Primary Communication**: WebSocket for real-time messaging
- **Authentication Helper**: REST API for user management
- **Progressive Enhancement**: Anonymous → Authenticated sessions
- **Single Endpoint**: Unified WebSocket endpoint for all client types

### Connection Flow
```
1. Client initiates WebSocket connection
   ├── Anonymous: ws://server/ws/client
   ├── With Token: ws://server/ws/client?token=jwt
   └── With Header: Authorization: Bearer jwt

2. Server validates authentication (optional)
   ├── Valid Token: Authenticated session
   ├── Invalid Token: Graceful fallback to anonymous
   └── No Token: Anonymous session

3. Client can upgrade session
   ├── Send auth_login message
   ├── Send auth_token message
   └── Receive auth_success/auth_error response
```

## 🏛️ **Technical Architecture**

### Application Layer Structure
```
┌─────────────────────────────────────────────────────────────┐
│                     Application Layer                       │
├─────────────────────────────────────────────────────────────┤
│  REST API Handlers     │         WebSocket Handlers        │
│  - auth_handlers.go    │         - websocket_handlers.go   │
│  - command_handlers.go │         - hub.go (connection mgmt) │
│  - user_handlers.go    │         - client.go (session mgmt) │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                     Business Logic Layer                    │
├─────────────────────────────────────────────────────────────┤
│  Services              │         Authentication            │
│  - user_service.go     │         - auth.go (JWT handling)  │
│  - command_service.go  │         - middleware.go (HTTP)    │
│  - websocket_service   │         - session validation     │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│                       Data Layer                           │
├─────────────────────────────────────────────────────────────┤
│  Models                │         Database                  │
│  - models.go           │         - database.go             │
│  - User, Command, etc  │         - migrations.go           │
│  - JSON serialization │         - connection pooling      │
└─────────────────────────────────────────────────────────────┘
```

## 🔐 **Authentication Architecture**

### Multi-Method Authentication System
```
┌─────────────────────────────────────────────────────────────┐
│                  Authentication Methods                     │
├─────────────────────────────────────────────────────────────┤
│  REST API              │         WebSocket                 │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ POST /auth/     │   │  │ Connection Authentication   │  │
│  │ - register      │   │  │ - Header: Bearer <token>    │  │
│  │ - login         │   │  │ - Query: ?token=<jwt>       │  │
│  │ Returns JWT     │   │  │ - Progressive via messages  │  │
│  └─────────────────┘   │  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                                 │
┌─────────────────────────────────────────────────────────────┐
│                   Session Management                        │
├─────────────────────────────────────────────────────────────┤
│  Anonymous Sessions    │      Authenticated Sessions       │
│  - Receive broadcasts  │      - Full command access        │
│  - Limited functionality│      - User-specific messages    │
│  - Can upgrade later   │      - Command execution rights   │
└─────────────────────────────────────────────────────────────┘
```

### JWT Token Flow
```
1. User Registration/Login (REST API)
   ├── POST /api/v1/auth/register
   ├── POST /api/v1/auth/login
   └── Returns: JWT token + user data

2. Token Validation
   ├── Parse and verify JWT signature
   ├── Check token expiration
   ├── Extract user claims (ID, username, etc.)
   └── Return user context or error

3. Session Creation
   ├── Authenticated: Full access client
   ├── Anonymous: Limited access client
   └── Progressive: Upgrade anonymous to authenticated
```

## 📡 **WebSocket Hub Architecture**

### Hub-Based Connection Management
```go
type Hub struct {
    // Client management
    clients    map[*Client]bool           // All connected clients
    broadcast  chan []byte                // Broadcast channel
    register   chan *Client               // Client registration
    unregister chan *Client               // Client unregistration
    
    // User tracking
    userClients map[uint][]*Client        // Clients per user ID
    anonClients []*Client                 // Anonymous clients
    
    // Connection limits
    maxConnectionsPerUser int             // Per-user connection limit
}
```

### Client Session Structure
```go
type Client struct {
    // Authentication state
    UserID     *uint                      // nil for anonymous
    Username   string                     // empty for anonymous
    Token      string                     // JWT token if authenticated
    
    // Connection management  
    Connection *websocket.Conn            // WebSocket connection
    Hub        *Hub                       // Reference to hub
    Send       chan []byte                // Outgoing message queue
    
    // Session metadata
    ConnectedAt time.Time                 // Connection timestamp
    LastPing    time.Time                 // Last ping received
}
```

### Message Routing System
```
┌─────────────────────────────────────────────────────────────┐
│                    Message Flow                             │
├─────────────────────────────────────────────────────────────┤
│  Incoming Message      │         Processing                │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ Raw WebSocket   │───│→ │ JSON Parse                  │  │
│  │ Message         │   │  │ ├── Extract type            │  │
│  └─────────────────┘   │  │ ├── Validate structure      │  │
│                        │  │ └── Route to handler        │  │
│                        │  └─────────────────────────────┘  │
│                        │                │                  │
│  ┌─────────────────────────────────────┼─────────────────┐ │
│  │              Message Handlers       │                 │ │
│  │  ┌─────────────┐ ┌─────────────┐   │ ┌─────────────┐ │ │
│  │  │ auth_login  │ │ auth_token  │   │ │    ping     │ │ │
│  │  │ auth_error  │ │ auth_success│   │ │    pong     │ │ │
│  │  └─────────────┘ └─────────────┘   │ └─────────────┘ │ │
│  └─────────────────────────────────────┼─────────────────┘ │
│                        │                │                  │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ Response        │←──│─ │ Generate Response           │  │
│  │ Message         │   │  │ ├── Create response object  │  │
│  └─────────────────┘   │  │ ├── JSON serialize          │  │
│                        │  │ └── Send via client channel │  │
│                        │  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

## 🗄️ **Database Architecture**

### Multi-Database Support
```
┌─────────────────────────────────────────────────────────────┐
│                  Database Abstraction                       │
├─────────────────────────────────────────────────────────────┤
│  GORM ORM Layer        │         Database Drivers          │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ Model Definitions│   │  │ SQLite (Development)        │  │
│  │ - User          │   │  │ - File: data/controlme.db   │  │
│  │ - Command       │   │  │ - No connection pooling     │  │
│  │ - Instruction   │   │  └─────────────────────────────┘  │
│  │ - Tag, Block    │   │  ┌─────────────────────────────┐  │
│  │ - Report        │   │  │ PostgreSQL (Production)     │  │
│  └─────────────────┘   │  │ - Connection pooling        │  │
│                        │  │ - Advanced features         │  │
│                        │  │ - Scalable performance      │  │
│                        │  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Data Model Relationships
```
┌─────────────────┐         ┌─────────────────┐
│      Users      │         │    Commands     │
├─────────────────┤         ├─────────────────┤
│ ID (PK)         │←────────│ UserID (FK)     │
│ Username        │         │ Name            │
│ Email           │         │ Description     │
│ Password (hash) │         │ Instructions[]  │
│ CreatedAt       │         │ CreatedAt       │
│ UpdatedAt       │         │ UpdatedAt       │
└─────────────────┘         └─────────────────┘
         │                           │
         │                           │
┌─────────────────┐         ┌─────────────────┐
│     Blocks      │         │      Tags       │
├─────────────────┤         ├─────────────────┤
│ UserID (FK)     │         │ Name            │
│ BlockedID (FK)  │         │ Description     │
│ CreatedAt       │         │ Color           │
└─────────────────┘         └─────────────────┘
         │
┌─────────────────┐
│     Reports     │
├─────────────────┤
│ ReporterID (FK) │
│ ReportedID (FK) │
│ Reason          │
│ CreatedAt       │
└─────────────────┘
```

### JSON Instruction Storage
```go
// Instructions stored as JSON array in database
type Command struct {
    Instructions []Instruction `gorm:"type:json" json:"instructions"`
}

type Instruction struct {
    Type              string            `json:"type"`               // shell, powershell, cmd, etc.
    Command           string            `json:"command"`            // Command to execute
    Args              []string          `json:"args,omitempty"`     // Command arguments
    Timeout           int               `json:"timeout,omitempty"`  // Execution timeout
    WorkingDirectory  string            `json:"working_directory,omitempty"`
    Environment       map[string]string `json:"environment,omitempty"`
}
```

## 🚦 **Request/Response Flow**

### REST API Flow
```
1. HTTP Request
   ├── Client sends HTTP request
   ├── Gin router matches endpoint
   ├── Middleware processes request (CORS, logging)
   ├── Handler validates input
   ├── Service layer processes business logic
   ├── Database layer performs data operations
   ├── Response marshalled to JSON
   └── HTTP response sent to client

2. Error Handling
   ├── Validation errors → 400 Bad Request
   ├── Authentication errors → 401 Unauthorized
   ├── Authorization errors → 403 Forbidden
   ├── Not found errors → 404 Not Found
   ├── Server errors → 500 Internal Server Error
   └── Proper error message formatting
```

### WebSocket Message Flow
```
1. Message Reception
   ├── WebSocket receives raw message
   ├── JSON parsing and validation
   ├── Message type extraction
   ├── Route to appropriate handler
   └── Generate response message

2. Message Broadcasting
   ├── Hub receives broadcast message
   ├── Determine target clients (all/authenticated/specific user)
   ├── Queue message for each target client
   ├── Client goroutines send messages
   └── Handle send failures and cleanup

3. Session Management
   ├── Client connection established
   ├── Registration with hub
   ├── Authentication status determination
   ├── Connection health monitoring
   ├── Graceful disconnection handling
   └── Resource cleanup
```

## 🔒 **Security Architecture**

### Authentication Security
```
┌─────────────────────────────────────────────────────────────┐
│                    Security Layers                          │
├─────────────────────────────────────────────────────────────┤
│  Input Validation      │         Authentication            │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ JSON validation │   │  │ JWT signature validation    │  │
│  │ Schema checking │   │  │ Token expiration checking   │  │
│  │ Size limits     │   │  │ User session validation     │  │
│  │ Type validation │   │  │ Permission level checking   │  │
│  └─────────────────┘   │  └─────────────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│  Connection Security   │         Data Protection           │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ Connection limits│   │  │ Password hashing (bcrypt)   │  │
│  │ Rate limiting   │   │  │ SQL injection prevention    │  │
│  │ Origin validation│   │  │ XSS prevention             │  │
│  │ Timeout handling │   │  │ Data sanitization          │  │
│  └─────────────────┘   │  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Authorization Levels
```
Anonymous Users:
├── Can connect to WebSocket
├── Receive broadcast messages only  
├── Cannot send commands
├── Cannot access user-specific data
└── Can upgrade to authenticated session

Authenticated Users:
├── Full WebSocket access
├── Can send and receive all message types
├── Access to user-specific commands
├── Can execute commands
└── Receive targeted messages

Admin Users (Future):
├── All authenticated user permissions
├── User management capabilities
├── System administration commands
├── Monitoring and analytics access
└── Configuration management
```

## ⚡ **Performance Architecture**  

### Concurrency Model
```go
// Goroutine per client for message handling
func (c *Client) writePump() {
    // Dedicated goroutine for sending messages
    for {
        select {
        case message := <-c.Send:
            // Send message to WebSocket connection
        case <-ticker.C:
            // Send ping to keep connection alive
        }
    }
}

func (c *Client) readPump() {
    // Dedicated goroutine for receiving messages
    for {
        // Read message from WebSocket connection
        // Process and route message
        // Handle errors and cleanup
    }
}
```

### Connection Management
```
┌─────────────────────────────────────────────────────────────┐
│                Connection Pool Management                    │
├─────────────────────────────────────────────────────────────┤
│  WebSocket Connections │         Database Connections       │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ Per-user limits │   │  │ Connection pooling (GORM)   │  │
│  │ Connection reuse│   │  │ Idle connection cleanup     │  │
│  │ Graceful cleanup│   │  │ Query optimization          │  │
│  │ Memory management│   │  │ Index utilization          │  │
│  └─────────────────┘   │  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Scalability Considerations
```
Current Architecture (Single Instance):
├── In-memory session management
├── Local WebSocket connections
├── Single database connection
└── Suitable for moderate load

Future Scaling (Multi-Instance):
├── Redis for session storage
├── Message queue for broadcasting
├── Load balancer with sticky sessions
├── Database read replicas
└── Microservice decomposition
```

## 🔧 **Configuration Architecture**

### Environment-Based Configuration
```yaml
# Development (config.test.yaml)
environment: development
server:
  port: 8082
database:
  type: sqlite
  path: data/controlme.db
auth:
  jwt_secret: "development-secret"

# Production (config.yaml)  
environment: production
server:
  port: 8081
database:
  type: postgres
  host: postgres-server
  name: controlme_prod
auth:
  jwt_secret: "${JWT_SECRET_FROM_ENV}"
```

### Configuration Hierarchy
```
1. Environment Variables (highest priority)
   ├── CONFIG_FILE
   ├── JWT_SECRET
   └── DATABASE_URL

2. Configuration Files
   ├── config.yaml (default)
   ├── config.test.yaml (testing)
   ├── config.docker.yaml (containerized)
   └── config.example.yaml (template)

3. Application Defaults (lowest priority)
   ├── Default port: 8081
   ├── Default database: SQLite
   └── Default connection limits
```

## 📊 **Monitoring and Observability**

### Logging Architecture
```
┌─────────────────────────────────────────────────────────────┐
│                     Logging System                          │
├─────────────────────────────────────────────────────────────┤
│  Structured Logging    │         Log Levels                │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ JSON formatting │   │  │ DEBUG: Detailed info        │  │
│  │ Field consistency│   │  │ INFO: General operations    │  │
│  │ Timestamp sync  │   │  │ WARN: Potential issues      │  │
│  │ Context tracking│   │  │ ERROR: Error conditions     │  │
│  └─────────────────┘   │  │ FATAL: Application crashes  │  │
│                        │  └─────────────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│  Log Categories        │         Output Destinations       │
│  ┌─────────────────┐   │  ┌─────────────────────────────┐  │
│  │ HTTP requests   │   │  │ Console (development)       │  │
│  │ WebSocket ops   │   │  │ File rotation (production)  │  │
│  │ Database queries│   │  │ External systems (ELK)      │  │
│  │ Authentication  │   │  │ Monitoring tools            │  │
│  └─────────────────┘   │  └─────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Health Monitoring
```go
// Health check endpoint
GET /health
{
    "status": "healthy",
    "database": "connected",
    "websocket_connections": 42,
    "uptime": "2h34m12s",
    "version": "1.0.0"
}
```

This architecture provides a solid foundation for the real-time command distribution platform with room for future scaling and enhancement.