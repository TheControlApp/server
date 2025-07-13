# ControlMe Go Development Makefile

.PHONY: help dev build test clean docker-up docker-down logs

# Default target
help:
	@echo "ControlMe Go Development Commands:"
	@echo "  dev         - Start development environment with hot reloading"
	@echo "  build       - Build the server binary"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start Docker services only"
	@echo "  docker-down - Stop Docker services"
	@echo "  logs        - Show Docker service logs"

# Start development environment
dev:
	@echo "🚀 Starting development environment..."
	@./dev-start.sh

# Build the server
build:
	@echo "🔨 Building server..."
	@go build -o bin/server cmd/server/main.go
	@echo "✅ Build complete: bin/server"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/ tmp/ build-errors.log
	@echo "✅ Clean complete"

# Start Docker services only
docker-up:
	@echo "🐳 Starting Docker services..."
	@docker-compose up -d postgres redis

# Stop Docker services
docker-down:
	@echo "🐳 Stopping Docker services..."
	@docker-compose down

# Show Docker service logs
logs:
	@echo "📋 Docker service logs:"
	@docker-compose logs -f

# Install Air if not present
install-air:
	@echo "📦 Installing Air..."
	@go install github.com/air-verse/air@latest

# Run without hot reloading
run:
	@echo "▶️  Running server..."
	@go run cmd/server/main.go
