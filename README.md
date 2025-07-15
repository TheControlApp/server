# ControlMe Go Backend

A modern, secure, and scalable rewrite of the ControlMe platform in Go, providing a clean modern API for applications.

## 🚀 Features

- ✅ **Modern API**: RESTful API with JWT authentication
- ✅ **Real-time Communication**: WebSocket support for instant messaging
- ✅ **Secure**: Modern authentication, bcrypt password hashing, HTTPS support
- ✅ **Scalable**: Docker-based deployment, Redis caching
- ✅ **Cross-platform**: Runs on Linux, macOS, and Windows
- ✅ **Well-tested**: Comprehensive test coverage

## 📁 Project Structure

```
controlme-go/
├── cmd/
│   ├── server/              # Main application entry point
│   └── tools/               # Development and maintenance tools
├── internal/
│   ├── api/
│   │   ├── handlers/        # HTTP request handlers
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

The server will be available at `http://localhost:8080`

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
make docker-up     # Start Docker services
make docker-down   # Stop Docker services
make seed          # Run database seed data
```

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

Base URL: `http://localhost:8080/api/v1`

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
- `WS /ws/client` - Client WebSocket connection
- `WS /ws/web` - Web client WebSocket connection

## 🏗️ Architecture

### Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis
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

## 📊 Monitoring & Logging

### Health Check
```bash
curl http://localhost:8080/health
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
