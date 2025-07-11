# ControlMe Go Backend

A modern, secure, and scalable rewrite of the ControlMe platform in Go.

## Project Structure

```
controlme-go/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/         # HTTP handlers
│   │   ├── middleware/       # HTTP middleware
│   │   └── routes/          # Route definitions
│   ├── auth/                # Authentication logic
│   ├── config/              # Configuration management
│   ├── database/            # Database models and migrations
│   ├── legacy/              # Legacy API compatibility layer
│   ├── models/              # Data models
│   ├── services/            # Business logic
│   └── websocket/           # WebSocket hub and handlers
├── migrations/              # Database migration files
├── configs/                 # Configuration files
├── scripts/                 # Build and deployment scripts
├── docker/                  # Docker configuration
├── docs/                    # Documentation
└── tests/                   # Test files
```

## Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 15+
- Redis (optional, for WebSocket scaling)

### Installation

1. Clone the repository
2. Install dependencies: `go mod download`
3. Set up PostgreSQL database
4. Copy `configs/config.example.yaml` to `configs/config.yaml`
5. Run migrations: `go run cmd/migrate/main.go up`
6. Start the server: `go run cmd/server/main.go`

## Development

### Running the Server
```bash
go run cmd/server/main.go
```

### Running Tests
```bash
go test ./...
```

### Database Migrations
```bash
# Create new migration
go run cmd/migrate/main.go create migration_name

# Run migrations
go run cmd/migrate/main.go up

# Rollback migrations
go run cmd/migrate/main.go down
```

## API Endpoints

### New API (v1)
- `POST /api/v1/auth/login` - User authentication
- `GET /api/v1/commands/pending` - Get pending commands
- `POST /api/v1/commands/complete` - Mark command as completed
- `WebSocket /ws/client` - Real-time command delivery
- `WebSocket /ws/web` - Web client real-time features

### Legacy API (Exact Compatibility)
- `GET /AppCommand.aspx` - Legacy command polling
- `GET /GetContent.aspx` - Legacy command fetching
- `GET /GetCount.aspx` - Legacy command count
- `POST /ProcessComplete.aspx` - Legacy command completion

## Features

- ✅ Real-time command delivery via WebSocket
- ✅ Legacy client compatibility (exact API match)
- ✅ Modern authentication with JWT
- ✅ Secure password hashing with bcrypt
- ✅ Cross-platform deployment
- ✅ Comprehensive logging and monitoring
- ✅ Gradual migration support

## License

This project is for educational and research purposes only.
