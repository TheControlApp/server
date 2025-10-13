# ControlApp Server Documentation

## Overview

The ControlApp Server is a modern, secure, and scalable command distribution platform built with Go. It provides real-time WebSocket communication for delivering instructions between users, with comprehensive REST API endpoints for user management and authentication.

## Key Features

- **Real-time Communication**: WebSocket-based command distribution with support for both authenticated and anonymous connections
- **Secure Authentication**: JWT token-based authentication with bcrypt password hashing
- **Flexible Command System**: JSON-based instruction format supporting multiple command types
- **User Management**: Complete user registration, authentication, and profile management
- **Database Flexibility**: Support for both PostgreSQL (production) and SQLite (development)
- **RFC 7807 Compliance**: Structured error responses following HTTP API best practices
- **Comprehensive Documentation**: Complete API reference, examples, and integration guides

## Quick Start

### Prerequisites

- Go 1.19 or later
- PostgreSQL 12+ (optional, SQLite fallback available)
- Air (for development hot reload) - optional

### Installation

```bash
# Clone the repository
git clone https://github.com/TheControlApp/server.git
cd server

# Install dependencies
go mod download

# Copy and configure settings
cp configs/config.example.yaml config.yaml
# Edit config.yaml with your database settings

# Run with hot reload (development)
air

# Or run directly
go run cmd/server/main.go
```

### Configuration

```yaml
# config.yaml
server:
  port: 8080
  environment: "development"

database:
  type: "postgres"  # or "sqlite"
  host: "localhost"
  port: 5432
  name: "controlapp"
  username: "postgres"
  password: "your_password"

auth:
  jwt_secret: "your-secret-key-here"
  jwt_expiration: 86400  # 24 hours
```

### First API Call

```bash
# Health check
curl http://localhost:8080/health

# Register a user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "screen_name": "Test User"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## Documentation Structure

### 📚 API Documentation
- **[REST API Reference](api/rest-api.md)** - Complete HTTP endpoint documentation
- **[WebSocket API Reference](api/websocket-api.md)** - Real-time communication protocol
- **[Authentication Guide](api/authentication.md)** - JWT tokens, security, and auth flows

### 🏗️ Architecture Documentation
- **[Services & Business Logic](architecture/services.md)** - Service layer architecture and patterns
- **[Database Models](reference/data-models.md)** - Complete schema and relationships

### 🛠️ Developer Resources
- **[Getting Started Guide](getting-started/quick-start.md)** - Setup and first steps
- **[API Examples](examples/)** - Code examples and integration patterns
- **[Swagger/OpenAPI](swagger/swagger.yaml)** - Interactive API documentation
- **[API Validation Results](API_VALIDATION_RESULTS.md)** - ✅ Verified documentation accuracy

### 🧪 Testing Resources
- **`api_test_registration.json`** - Ready-to-use registration test data
- **`api_test_login.json`** - Ready-to-use login credentials
- **`test_empty_registration.json`** - For testing validation errors

## Core Concepts

### Commands and Instructions

Commands are the fundamental unit of communication in the ControlApp system. Each command contains one or more instructions that specify actions to be performed.

```json
{
  "instructions": [
    {
      "type": "std_popup",
      "content": {
        "body": "Hello, World!",
        "button": "OK"
      }
    }
  ],
  "receiver_id": "123e4567-e89b-12d3-a456-426614174000",
  "tags": "general"
}
```

### User Roles and Authentication

- **Anonymous Users**: Can connect via WebSocket, limited command access
- **Registered Users**: Full access to commands, personal command history
- **Admin Users**: Extended privileges (future enhancement)

### Real-time Communication

The WebSocket endpoint (`/ws/client`) provides:
- Bidirectional real-time messaging
- Command broadcasting to all connected clients
- Targeted messaging to specific users
- Authentication via multiple methods (header, query, message-based)

## API Endpoints

### Health & Status
- `GET /health` - Server health check

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User authentication

### User Management
- `GET /api/v1/users` - List all users
- `GET /api/v1/users/{id}` - Get user by ID

### Command Management
- `GET /api/v1/commands/pending` - Get pending commands for user
- `POST /api/v1/commands/complete` - Mark command as completed

### WebSocket
- `GET /ws/client` - WebSocket connection endpoint

## Development

### Project Structure

```
server/
├── cmd/server/           # Application entry point
├── internal/
│   ├── api/             # HTTP handlers and routes
│   ├── auth/            # Authentication services
│   ├── config/          # Configuration management
│   ├── database/        # Database setup and migrations
│   ├── middleware/      # HTTP middleware
│   ├── models/          # Data models
│   ├── services/        # Business logic layer
│   └── websocket/       # WebSocket hub and clients
├── docs/                # Documentation
├── configs/             # Configuration files
└── docs-new/            # Updated documentation
```

### Development Workflow

1. **Hot Reload Development**
   ```bash
   air  # Automatically rebuilds on file changes
   ```

2. **Database Migrations**
   ```bash
   # Automatic on startup, or manual via:
   go run cmd/server/main.go
   ```

3. **Testing**
   ```bash
   go test ./...
   ```

4. **API Documentation**
   - Swagger UI: http://localhost:8080/swagger/index.html
   - Generated from code comments and updated manually

### Docker Development

```bash
# Development with Docker Compose
docker-compose up

# Production build
docker-compose -f docker-compose.prod.yml up
```

## Production Deployment

### Environment Variables

```bash
# Server Configuration
PORT=8080
ENVIRONMENT=production

# Database
DB_TYPE=postgres
DB_HOST=your-db-host
DB_PORT=5432
DB_NAME=controlapp
DB_USERNAME=your-username
DB_PASSWORD=your-password

# Authentication
JWT_SECRET=your-super-secret-key
JWT_EXPIRATION=86400
```

### Security Considerations

1. **HTTPS**: Always use HTTPS in production
2. **JWT Secrets**: Use strong, randomly generated secrets
3. **Database Security**: Secure database access and use SSL/TLS
4. **Rate Limiting**: Implement rate limiting for API endpoints
5. **CORS**: Configure CORS for your frontend domains
6. **Input Validation**: All inputs are validated at service layer

### Performance Optimization

1. **Database**: Use PostgreSQL with proper indexing
2. **Connection Pooling**: Configure database connection limits
3. **WebSocket Scaling**: Consider Redis for multi-instance deployments
4. **Caching**: Implement Redis for session and data caching
5. **Load Balancing**: Use reverse proxy (nginx) for load distribution

## Integration Examples

### JavaScript/React Client

```javascript
// REST API Integration
const authService = {
  async login(username, password) {
    const response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    });
    const data = await response.json();
    return data;
  }
};

// WebSocket Integration
const ws = new WebSocket('ws://localhost:8080/ws/client');
ws.onmessage = (event) => {
  const command = JSON.parse(event.data);
  // Handle command...
};
```

### Python Client

```python
import requests
import websocket
import json

# REST API
response = requests.post('http://localhost:8080/api/v1/auth/login', 
                        json={'username': 'user', 'password': 'pass'})
data = response.json()

# WebSocket
def on_message(ws, message):
    command = json.loads(message)
    print(f"Received command: {command}")

ws = websocket.WebSocketApp("ws://localhost:8080/ws/client",
                           on_message=on_message)
ws.run_forever()
```

### Go Client

```go
import (
    "github.com/gorilla/websocket"
    "net/http"
)

// REST API client
client := &http.Client{}
resp, err := client.Post("http://localhost:8080/api/v1/auth/login", 
                        "application/json", 
                        strings.NewReader(`{"username":"user","password":"pass"}`))

// WebSocket client
c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws/client", nil)
defer c.Close()

for {
    _, message, err := c.ReadMessage()
    if err != nil {
        break
    }
    // Handle message...
}
```

## Support and Contributing

### Getting Help

1. **Documentation**: Start with this documentation and API reference
2. **Examples**: Check the `/examples` directory for code samples
3. **Issues**: Report bugs and request features on GitHub
4. **Discussions**: Join the community discussions

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Update documentation if needed
5. Submit a pull request

### Development Guidelines

- Follow Go best practices and conventions
- Write tests for new functionality
- Update documentation for API changes
- Use meaningful commit messages
- Ensure code passes linting and tests

## Roadmap

### Planned Features

- [ ] Command persistence and history
- [ ] User blocking and reporting system
- [ ] Advanced command scheduling
- [ ] Admin dashboard and moderation tools
- [ ] Rate limiting and abuse prevention
- [ ] Advanced analytics and metrics
- [ ] Mobile app APIs
- [ ] Plugin system for custom instructions

### Performance & Scaling

- [ ] Redis integration for caching and sessions
- [ ] Database read replicas
- [ ] Horizontal scaling with load balancers
- [ ] WebSocket clustering with Redis pub/sub
- [ ] Performance monitoring and alerting

## License

This project is licensed under the Apache 2.0 License. See the LICENSE file for details.

## Changelog

### Version 1.0.0 (Current)
- Initial release with core functionality
- REST API for user management and authentication
- WebSocket real-time command distribution
- PostgreSQL and SQLite database support
- JWT authentication with bcrypt password hashing
- Comprehensive documentation and examples