# Development Environment Setup

## 🚀 **Quick Start**

### Prerequisites
- **Go 1.21+** installed and in PATH
- **Git** for version control  
- **VS Code** with Go extension (recommended)
- **Windows PowerShell** (current environment)

### One-Minute Setup
```powershell
# Clone and navigate to project
cd D:\Workspace\github.com\TheControlApp\server

# Install dependencies
go mod download

# Build the server
go build -o tmp/server.exe cmd/server/main.go

# Run with test configuration
$env:CONFIG_FILE="config.test.yaml"; .\tmp\server.exe
```

**Server will be running at**: `http://localhost:8082`  
**WebSocket endpoint**: `ws://localhost:8082/ws/client`  
**API Documentation**: `http://localhost:8082/swagger/index.html`

## 📋 **Detailed Setup Guide**

### Step 1: Environment Verification
```powershell
# Check Go installation
go version
# Should output: go version go1.21+ windows/amd64

# Check Git installation  
git --version
# Should output: git version 2.x.x

# Verify workspace location
pwd
# Should be: D:\Workspace\github.com\TheControlApp\server
```

### Step 2: Dependencies Installation
```powershell
# Download all Go modules
go mod download

# Verify dependencies are installed
go mod verify
# Should output: all modules verified

# Optional: Clean module cache if issues occur
go clean -modcache
go mod download
```

### Step 3: Configuration Setup
The project includes multiple configuration files:

**Development Configuration** (Recommended):
```powershell
# Use SQLite test configuration
$env:CONFIG_FILE="config.test.yaml"
```

**Production Configuration**:
```powershell  
# Use PostgreSQL production configuration
$env:CONFIG_FILE="config.yaml"
```

**Docker Configuration**:
```powershell
# Use containerized configuration
$env:CONFIG_FILE="config.docker.yaml"
```

### Step 4: Database Setup

#### SQLite (Development - Recommended)
```powershell
# SQLite database is automatically created
# Location: data/controlme.db
# No additional setup required
```

#### PostgreSQL (Production)
```powershell
# Start PostgreSQL with Docker
docker-compose up -d postgres

# Or install PostgreSQL manually
# Configure connection in config.yaml:
# database:
#   host: localhost
#   port: 5432
#   name: controlme
#   username: postgres
#   password: postgres
```

### Step 5: Build and Run
```powershell
# Build main server
go build -o tmp/server.exe cmd/server/main.go

# Build test tools (optional)
go build -o bin/test-websocket-auth.exe cmd/tools/test-websocket-auth/main.go
go build -o bin/create-test-user.exe cmd/tools/create-test-user/main.go

# Run server with test configuration
$env:CONFIG_FILE="config.test.yaml"; .\tmp\server.exe
```

## 🛠️ **Development Tools**

### Available Build Targets
```powershell
# Main server application
go build -o tmp/server.exe cmd/server/main.go

# WebSocket authentication test client
go build -o bin/test-websocket-auth.exe cmd/tools/test-websocket-auth/main.go

# General test client
go build -o bin/test-client.exe cmd/tools/test-client/main.go

# User creation utility
go build -o bin/create-test-user.exe cmd/tools/create-test-user/main.go

# Command creation utility
go build -o bin/create-commands.exe cmd/tools/create-commands/main.go
```

### Using Mage Build Tool (Optional)
```powershell
# Install Mage (optional build tool)
go install github.com/magefile/mage@latest

# Available Mage targets
mage -l
# build      Build the server
# test       Run all tests
# clean      Clean build artifacts
# docker     Build Docker image

# Build with Mage
mage build

# Run tests with Mage
mage test
```

### Using Air for Hot Reload (Optional)
```powershell
# Install Air for hot reloading
go install github.com/cosmtrek/air@latest

# Run with hot reload (if .air.toml exists)
air

# Or run without config
air -c .air.toml
```

## 🧪 **Testing Setup**

### Running Tests
```powershell
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run specific test packages
go test ./internal/auth/
go test ./internal/database/
go test ./internal/websocket/

# Run tests with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using Test Clients

#### WebSocket Test Client
```powershell
# Build and run WebSocket test client
go build -o bin/test-websocket-auth.exe cmd/tools/test-websocket-auth/main.go
.\bin\test-websocket-auth.exe

# Interactive commands:
# ping - Send ping message
# login <username> <password> - Authenticate  
# quit - Exit client
```

#### Manual Testing with curl
```powershell
# Test user registration
curl -X POST http://localhost:8082/api/v1/auth/register `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"password123","email":"test@example.com"}'

# Test user login
curl -X POST http://localhost:8082/api/v1/auth/login `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"password123"}'

# Test health endpoint
curl http://localhost:8082/health
```

#### WebSocket Testing with wscat
```powershell
# Install wscat (requires Node.js)
npm install -g wscat

# Connect anonymously
wscat -c ws://localhost:8082/ws/client

# Send ping message
{"type":"ping","payload":{}}

# Send authentication message
{"type":"auth_login","payload":{"username":"testuser","password":"password123"}}
```

## 🐳 **Docker Development**

### Using Docker Compose
```powershell
# Start all services (app + database)
docker-compose up

# Start in background
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild services
docker-compose up --build
```

### Docker Development Commands
```powershell
# Build Docker image
docker build -t controlme-server .

# Run with Docker
docker run -p 8081:8081 controlme-server

# Run with custom config
docker run -p 8081:8081 -v ${PWD}/config.docker.yaml:/app/config.yaml controlme-server
```

## 🔧 **IDE Configuration**

### VS Code Setup
**Recommended Extensions**:
- Go (official Go extension)
- REST Client (for API testing)
- GitLens (enhanced Git integration)
- Thunder Client (API testing)

**VS Code Settings** (`.vscode/settings.json`):
```json
{
    "go.toolsManagement.checkForUpdates": "local",
    "go.useLanguageServer": true,
    "go.formatTool": "goimports",
    "go.lintTool": "golangci-lint",
    "go.testFlags": ["-v"],
    "go.testTimeout": "30s"
}
```

**Launch Configuration** (`.vscode/launch.json`):
```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Server",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/server/main.go",
            "env": {
                "CONFIG_FILE": "config.test.yaml"
            },
            "args": []
        }
    ]
}
```

### GoLand Setup
```
1. Open project in GoLand
2. Go to File → Settings → Go → Build Tags & Vendoring
3. Set Build tags: development
4. Configure Run Configuration:
   - Program: cmd/server/main.go
   - Environment: CONFIG_FILE=config.test.yaml
```

## 📊 **Development Workflow**

### Recommended Development Process
```
1. Start Development Environment
   ├── Set CONFIG_FILE=config.test.yaml
   ├── Start server: go run cmd/server/main.go
   └── Verify at: http://localhost:8082/health

2. Make Code Changes
   ├── Follow Go conventions and standards
   ├── Add tests for new functionality
   ├── Update documentation as needed
   └── Verify with test clients

3. Testing
   ├── Run unit tests: go test ./...
   ├── Run integration tests with test clients
   ├── Test WebSocket connections
   └── Verify API endpoints with Swagger UI

4. Code Quality
   ├── Format code: go fmt ./...
   ├── Run linter: golangci-lint run
   ├── Check for vulnerabilities: go mod audit
   └── Update dependencies: go mod tidy

5. Commit and Deploy
   ├── Commit changes with descriptive messages
   ├── Tag releases appropriately
   ├── Build production binaries
   └── Deploy to target environment
```

### Git Workflow
```powershell
# Create feature branch
git checkout -b feature/new-functionality

# Make changes and commit
git add .
git commit -m "feat: add new functionality"

# Push to remote
git push origin feature/new-functionality

# Create pull request (via GitHub/GitLab)
# Merge after review
```

## 🚨 **Troubleshooting**

### Common Issues and Solutions

#### Port Already in Use
```powershell
# Check what's using port 8082
netstat -ano | findstr :8082

# Kill process using the port
taskkill /PID <process_id> /F

# Or use different port in config.test.yaml
server:
  port: 8083
```

#### Database Connection Issues
```powershell
# For SQLite: Check file permissions
# Ensure data/ directory exists and is writable

# For PostgreSQL: Verify connection settings
# Check if PostgreSQL is running
docker-compose ps
```

#### Module Download Issues
```powershell
# Clear module cache
go clean -modcache

# Set Go proxy (if behind corporate firewall)
$env:GOPROXY="direct"
$env:GOSUMDB="off"

# Download dependencies again
go mod download
```

#### Build Issues
```powershell
# Clean build cache
go clean -cache

# Rebuild with verbose output
go build -v -o tmp/server.exe cmd/server/main.go

# Check for import issues
go mod tidy
```

### Performance Debugging
```powershell
# Enable profiling
go build -race -o tmp/server.exe cmd/server/main.go

# Run with memory profiling
$env:CONFIG_FILE="config.test.yaml"; .\tmp\server.exe -memprofile=mem.prof

# Analyze profile
go tool pprof mem.prof
```

### WebSocket Debugging
```powershell
# Enable WebSocket debug logging
# Add to config.test.yaml:
logging:
  level: debug

# Monitor WebSocket connections
# Check browser dev tools for WebSocket frames
# Use test client with verbose output
```

## 📚 **Additional Resources**

### Documentation Links
- [System Architecture](../architecture/system-overview.md)
- [WebSocket Implementation Guide](../../docs/WEBSOCKET_IMPLEMENTATION.md)
- [Complete API Reference](../../docs/COMPLETE_API_REFERENCE.md)
- [Project Status](../project-state/current-status.md)

### External Resources
- [Go Documentation](https://golang.org/doc/)
- [Gin Framework Guide](https://gin-gonic.com/docs/)
- [GORM Documentation](https://gorm.io/docs/)
- [WebSocket RFC](https://tools.ietf.org/html/rfc6455)

### Community Resources
- [Go Community](https://golang.org/help/)
- [Gin Community](https://github.com/gin-gonic/gin)
- [WebSocket Best Practices](https://www.websocket.org/)

This setup guide should get any developer up and running with the ControlMe Go server in under 5 minutes.