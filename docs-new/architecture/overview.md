# Architecture Overview

## System Overview

The ControlApp Server is a real-time command delivery platform built in Go that enables users to send structured commands with instructions to other users via WebSocket connections. The system supports both authenticated and anonymous connections with flexible message routing capabilities.

## Core Concepts

### Commands
Commands are JSON messages containing arrays of instructions sent via WebSocket. They represent high-level tasks that can be executed by client applications.

```json
{
  "id": "cmd-123",
  "name": "Morning Routine",
  "description": "Daily startup tasks",
  "instructions": [
    {
      "type": "std_popup",
      "content": {
        "title": "Good Morning!",
        "body": "Ready to start your day?"
      }
    },
    {
      "type": "std_timer", 
      "content": {
        "title": "Coffee Break",
        "duration": 300
      }
    }
  ]
}
```

### Instructions
Instructions are individual actions within commands that client applications can execute. The system supports standard instruction types like popups, timers, notifications, and forms.

**Standard Instruction Types:**
- `std_popup` - Display alert/confirmation dialogs
- `std_notification` - Show system notifications
- `std_timer` - Execute timed delays
- `std_input` - Display input forms
- `std_download` - File download requests

### Broadcasting
Commands can be sent to all connected clients (broadcast) or targeted to specific users. This enables both group coordination and individual task assignment.

### Tags and Blocks
- **Tags** - Categorize and organize commands for easier management
- **Blocks** - Group related instructions within commands for logical organization

## System Architecture

### High-Level Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Clients   │    │  Desktop Apps   │    │  Mobile Apps    │
│                 │    │                 │    │                 │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │    WebSocket Gateway    │
                    │   ws://server/ws/client │
                    └────────────┬────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │     WebSocket Hub       │
                    │  Connection Management  │
                    │   Message Routing       │
                    └────────────┬────────────┘
                                 │
            ┌────────────────────┼────────────────────┐
            │                    │                    │
   ┌────────▼────────┐  ┌────────▼────────┐  ┌────────▼────────┐
   │   HTTP API      │  │  Authentication │  │   Database      │
   │   (REST)        │  │     Service     │  │   (PostgreSQL/  │
   │                 │  │                 │  │    SQLite)      │
   └─────────────────┘  └─────────────────┘  └─────────────────┘
```

### Component Architecture

#### 1. WebSocket Layer (`internal/websocket/`)

**Hub** - Central connection manager
- Maintains registry of all client connections
- Routes messages between clients
- Enforces connection limits and policies
- Handles connection lifecycle events

**Client** - Individual connection wrapper
- Represents a single WebSocket connection
- Stores authentication state and user information
- Manages outgoing message queue
- Handles connection-specific operations

#### 2. API Layer (`internal/api/`)

**Handlers** - HTTP and WebSocket request processors
- REST API endpoints for CRUD operations
- WebSocket connection upgrade and message handling
- Request validation and response formatting
- Error handling and logging

**Routes** - URL routing and middleware configuration
- RESTful route definitions
- Authentication middleware
- CORS and security headers
- Rate limiting and request logging

#### 3. Service Layer (`internal/services/`)

**UserService** - User management and authentication
- User registration and login
- JWT token generation and validation
- Password hashing and verification
- User profile management

**CommandService** - Command lifecycle management
- Command creation and storage
- Instruction processing and validation
- Command execution tracking
- Result aggregation and reporting

#### 4. Data Layer (`internal/models/` & `internal/database/`)

**Models** - Data structure definitions
- User, Command, Instruction entities
- Database relationships and constraints
- Validation rules and business logic
- JSON serialization configuration

**Database** - Persistence and query layer
- GORM ORM for database operations
- Migration management
- Connection pooling and configuration
- Support for PostgreSQL and SQLite

## Connection Management

### Connection Types

#### Anonymous Connections
- No authentication required
- Receive broadcast messages only
- Limited functionality
- Can upgrade to authenticated session

```javascript
// Anonymous connection
const ws = new WebSocket('ws://localhost:8080/ws/client');
```

#### Authenticated Connections
- JWT token required
- Full access to user-specific features
- Can send and receive targeted messages
- Complete API functionality

```javascript
// Authenticated connection via query parameter
const ws = new WebSocket('ws://localhost:8080/ws/client?token=jwt-token');

// Or via message after connection
ws.send(JSON.stringify({
    type: 'auth',
    token: 'jwt-token'
}));
```

### Connection Lifecycle

1. **Connection Establishment**
   - Client initiates WebSocket upgrade request
   - Server validates request and upgrades protocol
   - Connection registered in Hub

2. **Authentication** (Optional)
   - Token validation via header, query parameter, or message
   - User session creation and permission assignment
   - Connection upgrade from anonymous to authenticated

3. **Message Exchange**
   - Bidirectional message communication
   - Message type routing and processing
   - Error handling and response generation

4. **Heartbeat Monitoring**
   - Automatic ping/pong messages
   - Connection health monitoring
   - Stale connection cleanup

5. **Disconnection**
   - Graceful closure handling
   - Resource cleanup and deregistration
   - Session state preservation (if needed)

## Message Routing

### Inbound Message Flow

```
WebSocket Message → Authentication Check → Type Routing → Handler Processing → Response Generation
```

**Message Types:**
- `auth` - Authentication requests
- `ping` - Connection health checks
- `command` - Command execution requests
- `broadcast` - Broadcast message requests

### Outbound Message Flow

```
Service/Handler → Hub → Target Selection → Client Filtering → WebSocket Delivery
```

**Routing Strategies:**
- **Broadcast** - All connected clients
- **User-specific** - Clients belonging to specific user
- **Anonymous-only** - Anonymous clients only
- **Authenticated-only** - Authenticated clients only

## Security Architecture

### Authentication & Authorization

**JWT Token-based Authentication**
- Stateless token validation
- Configurable token expiration
- Role-based access control
- Automatic token refresh handling

**Authorization Levels:**
- **Anonymous** - Broadcast messages only
- **Authenticated** - User-specific operations
- **Admin** - System administration (future)

### Security Measures

**Connection Security:**
- Configurable connection limits per user
- IP-based rate limiting
- Message size and frequency limits
- Automatic cleanup of inactive connections

**Data Security:**
- Password hashing with bcrypt
- JWT token signing and validation
- SQL injection prevention via ORM
- XSS protection in API responses

**WebSocket Security:**
- Origin validation for browser connections
- Protocol-level message validation
- Error message sanitization
- Connection state isolation

## Data Flow

### Command Execution Flow

```
1. User creates command via REST API
   ↓
2. Command stored in database with instructions
   ↓  
3. Command sent to target users via WebSocket
   ↓
4. Client applications receive and execute instructions
   ↓
5. Execution results sent back via WebSocket
   ↓
6. Results aggregated and stored
   ↓
7. Command completion status updated
```

### Real-time Communication Flow

```
Client A                Hub                Client B
   │                     │                     │
   ├─── send message ───→│                     │
   │                     ├─── route message ──→│
   │                     │                     ├─── process message
   │                     │←─── send response ──┤
   │←─── route response ─┤                     │
   ├─── handle response  │                     │
```

## Scalability Considerations

### Horizontal Scaling
- Stateless application design
- Database connection pooling
- Session state externalization
- Load balancer compatibility

### Performance Optimization
- Connection pooling and reuse
- Message batching for broadcasts
- Efficient JSON serialization
- Database query optimization

### Resource Management
- Memory-efficient connection handling
- Automatic cleanup of stale resources
- Configurable limits and timeouts
- Graceful degradation under load

## Error Handling Strategy

### Error Categories
- **Client Errors** (4xx) - Invalid requests, authentication failures
- **Server Errors** (5xx) - Internal failures, database issues
- **WebSocket Errors** - Connection issues, protocol violations
- **Business Logic Errors** - Command execution failures

### Error Response Format
All errors follow RFC 7807 Problem Details standard:

```json
{
  "type": "validation_error",
  "title": "Validation Failed", 
  "status": 422,
  "detail": "Username must be at least 3 characters",
  "instance": "/api/v1/auth/register"
}
```

### Error Recovery
- Automatic retry for transient failures
- Graceful degradation for non-critical features
- Client reconnection handling
- Data consistency preservation

## Development Principles

### Code Organization
- Clean architecture with clear layer separation
- Dependency injection for testability
- Interface-based design for flexibility
- Consistent error handling patterns

### Testing Strategy
- Unit tests for business logic
- Integration tests for API endpoints
- WebSocket connection testing
- Database migration testing

### Documentation Standards
- Comprehensive API documentation
- Code comments for complex logic
- Architecture decision records
- Deployment and configuration guides

## Future Architecture Considerations

### Planned Enhancements
- Microservice decomposition
- Event sourcing for command history
- Redis for session management
- Message queue for reliable delivery

### Monitoring & Observability
- Application metrics collection
- Distributed tracing
- Centralized logging
- Health check endpoints

### Security Enhancements
- OAuth2/OIDC integration
- Rate limiting per user/IP
- Audit logging
- Content filtering and validation