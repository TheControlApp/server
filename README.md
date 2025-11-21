# ControlMe Go Backend

A modern, secure, and scalable rewrite of the ControlMe platform in Go, providing a clean modern API for applications.

## 🚀 Features

- ✅ **Modern API**: RESTful API with JWT authentication
- ✅ **RFC 7807 Error Handling**: Standardized, developer-friendly error responses
- ✅ **Real-time Communication**: WebSocket support for instant messaging
- ✅ **Secure**: Modern authentication, bcrypt password hashing, HTTPS support
- ✅ **Interactive Documentation**: Swagger/OpenAPI documentation
- ✅ **Scalable**: Docker-based deployment
- ✅ **Cross-platform**: Runs on Linux, macOS, and Windows
- ✅ **Well-tested**: Comprehensive test coverage

## 🔥 New: RFC 7807 Compliant Error Handling

The API now provides structured, consistent error responses following the RFC 7807 Problem Details standard:

```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation",
  "help": "Please check the field requirements in the API documentation",
  "errors": [
    {
      "field": "username", 
      "message": "Username must be at least 3 characters long",
      "code": "MIN_LENGTH"
    }
  ]
}
```

**Benefits for client developers:**
- 🎯 **Consistent structure** across all error responses
- 🛠️ **Machine-readable** error types and codes
- 📝 **Detailed validation** errors with field-specific information
- 🔗 **Actionable guidance** for resolving issues
- 📖 **Complete documentation** with examples

## 📁 Project Structure

```
controlme-go/
├── cmd/
│   ├── server/              # Main application entry point
│   └── tools/               # Development and maintenance tools
├── internal/
│   ├── api/
│   │   ├── handlers/        # HTTP request handlers
│   │   ├── responses/       # RFC 7807 response types
│   │   └── routes/          # Route definitions
│   ├── auth/                # Authentication logic
│   ├── config/              # Configuration management
│   ├── database/            # Database connection and setup
│   ├── middleware/          # HTTP middleware
│   ├── models/              # Data models (GORM)
│   ├── services/            # Business logic layer
│   └── websocket/           # WebSocket hub and handlers
├── configs/                 # Configuration files
├── scripts/                 # Development and deployment scripts
├── docker/                  # Docker configuration
└── docs/                    # Documentation
```

## 🛠️ Quick Start

### Prerequisites

- **Go 1.21+**
- **Docker & Docker Compose**
- **Make** (optional, but recommended)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd controlme-go
   ```

2. **Set up the development environment**
   ```bash
   make setup
   ```

3. **Start the development server**
   ```bash
   make dev
   ```

The server will be available at `http://localhost:3080`

### Manual Setup (without Make)

1. **Start Docker services**
   ```bash
   ./scripts/docker.sh up
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Build and run**
   ```bash
   go build -o bin/server cmd/server/main.go
   ./bin/server
   ```

## 🔧 Development

### Available Commands

```bash
make help          # Show all available commands
make setup         # Set up development environment
make dev           # Start development server with hot reload
make build         # Build the server binary
make test          # Run all tests
make lint          # Run code linter
make fmt           # Format code
make clean         # Clean build artifacts
make docker-up     # Start Docker services (includes Swagger UI)
make docker-down   # Stop Docker services
make seed          # Run database seed data
make swagger       # Generate Swagger documentation
make swagger-serve # Generate docs and start server
```

### Docker Services

The Docker Compose setup includes:
- **PostgreSQL Database** (port 5432)
- **Go Server** with hot reload (port 3080)
- **Swagger UI** for API documentation (port 3080)
- **Nginx** reverse proxy for production (port 80/443)

#### Access Points
- **API Server**: http://localhost:3080
- **Built-in Swagger**: http://localhost:3080/swagger/index.html
- **Dedicated Swagger UI**: http://localhost:3080
- **Health Check**: http://localhost:3080/health

### Configuration

Copy the example configuration and modify as needed:

```bash
cp configs/config.example.yaml configs/config.yaml
```

Key configuration options:
- Database connection settings
- Server port and host
- JWT secret keys
- CORS settings
- Log levels

### Testing

Run the full test suite:
```bash
make test
```

Test specific packages:
```bash
go test ./internal/auth/...
go test ./internal/api/handlers/...
```

## 🌐 API Documentation

### Modern API (v1)

Base URL: `http://localhost:3080/api/v1`

All endpoints return RFC 7807 compliant error responses with structured error details and optional developer help. See [docs/ERROR_RESPONSE_REFERENCE.md](docs/ERROR_RESPONSE_REFERENCE.md) for complete error handling documentation.

#### Authentication
- `POST /auth/login` - User authentication
- `POST /auth/refresh` - Refresh JWT token

#### Commands
- `GET /commands/pending` - Get pending commands for user
- `POST /commands/complete` - Mark command as completed
- `POST /commands/create` - Create new command

#### Users
- `GET /users/profile` - Get user profile
- `PUT /users/profile` - Update user profile

#### WebSocket
- `WS /ws/client` - Universal WebSocket connection (all clients)

## 🎯 Building Custom Clients

ControlApp is designed to support rich 3rd party clients! Whether you're building a desktop app, mobile client, web interface, or IoT integration, we provide comprehensive resources:

### 📖 **Client Development Resources**
- **[docs/CLIENT_DEVELOPMENT_GUIDE.md](docs/CLIENT_DEVELOPMENT_GUIDE.md)** - Complete guide for building 3rd party clients
- **[docs/STANDARD_COMMANDS.md](docs/STANDARD_COMMANDS.md)** - Official command set that all clients should support
- **[docs/WEBSOCKET_IMPLEMENTATION.md](docs/WEBSOCKET_IMPLEMENTATION.md)** - WebSocket API implementation details

### 🛠️ **Reference Implementations**
- **JavaScript/Web** - Full-featured browser client with all standard commands
- **Python** - Async client library perfect for automation and scripting
- **Go** - High-performance client for system integration
- **Examples** - Platform-specific implementations (React Native, Flutter, Desktop)

### ⚡ **Quick Start for Client Developers**
```javascript
// Minimal client in just a few lines
const client = new WebSocket('ws://localhost:3080/ws/client?token=your-jwt');
client.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === 'command') {
        executeCommand(message.payload);
    }
};
```

## 📚 Complete Documentation

For comprehensive documentation including API references, implementation guides, and examples:

- **[docs/](docs/)** - Complete documentation directory
- **[docs/COMPLETE_API_REFERENCE.md](docs/COMPLETE_API_REFERENCE.md)** - Full REST + WebSocket API documentation  
- **[docs/ERROR_RESPONSE_REFERENCE.md](docs/ERROR_RESPONSE_REFERENCE.md)** - RFC 7807 error handling guide
- **[docs/WEBSOCKET_IMPLEMENTATION.md](docs/WEBSOCKET_IMPLEMENTATION.md)** - WebSocket implementation guide
- **[docs/API_SWAGGER.md](docs/API_SWAGGER.md)** - OpenAPI/Swagger documentation

## 🏗️ Architecture

### Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT with bcrypt password hashing
- **Real-time**: WebSocket with message hub
- **Deployment**: Docker Compose

### Key Components

1. **Modern API Layer**: RESTful API with proper HTTP methods and status codes
2. **Authentication Service**: JWT-based authentication with bcrypt password hashing
3. **WebSocket Hub**: Real-time message broadcasting and client management
4. **Command Service**: Business logic for command creation, assignment, and completion
5. **User Service**: User management, authentication, and profile handling

## 🚢 Deployment

### Docker Deployment

1. **Production deployment**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

2. **Environment variables**
   ```bash
   export DB_HOST=your-db-host
   export DB_PASSWORD=your-secure-password
   export JWT_SECRET=your-jwt-secret
   ```

### Manual Deployment

1. **Build for production**
   ```bash
   CGO_ENABLED=0 GOOS=linux go build -o controlme-server cmd/server/main.go
   ```

2. **Run with environment configuration**
   ```bash
   export ENVIRONMENT=production
   ./controlme-server
   ```

## 🧪 Testing

### Unit Tests
```bash
go test ./internal/...
```

### Integration Tests
```bash
go test -tags=integration ./...
```

### Load Testing
```bash
# TODO: Add load testing instructions
```

## �️ Development Tools

The `cmd/tools/` directory contains helpful development and testing utilities:

### Integration Test Tool
Comprehensive API validation tool that tests the server against its documentation:

```bash
# Run the integration test
go run cmd/tools/integration-test/main.go
```

**What it tests:**
- ✅ REST API endpoints (`/auth/register`, `/auth/login`)
- ✅ JWT token generation and validation
- ✅ WebSocket connections and authentication
- ✅ RFC 7807 error response format
- ✅ All documented authentication methods

**Use cases:**
- 🔍 Validate API changes don't break existing functionality
- 📚 Ensure documentation stays accurate with implementation
- 🚀 CI/CD pipeline integration for automated testing
- 🐛 Debug authentication and WebSocket issues

The tool creates test users (`test1`, `test2`) and runs a full integration test suite, providing detailed output and appropriate exit codes for automated environments.

### Other Tools
- **`test-client/`** - Interactive WebSocket client for manual testing
- **`test-websocket-auth/`** - WebSocket authentication testing utility

See individual tool README files for detailed usage instructions.

## �📊 Monitoring & Logging

### Health Check
```bash
curl http://localhost:3080/health
```

### Metrics
- Application metrics available at `/metrics` (when enabled)
- Docker container metrics via `docker stats`

### Logging
- Structured JSON logging via logrus
- Log levels: debug, info, warn, error
- Configurable log output (stdout, file)

## 🔒 Security

- **Password Security**: bcrypt hashing with salt
- **JWT Authentication**: Secure token-based authentication
- **HTTPS Support**: TLS/SSL configuration available
- **CORS**: Configurable cross-origin resource sharing
- **Rate Limiting**: Built-in request rate limiting
- **Input Validation**: Comprehensive input sanitization

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linter (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

### Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write comprehensive tests
- Document public APIs
- Follow semantic commit messages

## 📝 License

This project is for educational and research purposes only.

## 🆘 Support

- **Documentation**: Check this README and inline code documentation
- **Issues**: Open an issue on GitHub

## 🗺️ Roadmap

### ✅ Phase 1: Modern Authentication (Complete)
- JWT-based authentication
- Bcrypt password hashing
- RESTful API design
- WebSocket communication

### 🔄 Phase 2: Enhanced Features (In Progress)
- Enhanced security features
- Improved error handling
- Comprehensive testing
- Performance optimization

### 📋 Phase 3: Advanced Features (Planned)
- Microservices architecture
- Advanced monitoring
- Load balancing
- Multi-tenant support
- API versioning strategy

---

**Last Updated**: July 2025  
**Version**: 1.0.0  
**Go Version**: 1.21+
