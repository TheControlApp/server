# ControlMe Go Development Makefile

.PHONY: help setup dev build test clean docker-up docker-down docker-clean install-tools lint fmt vet seed test-legacy logs install-air run

# Default target
help:
	@echo "ControlMe Go Development Commands:"
	@echo "  setup       - Set up development environment"
	@echo "  dev         - Start development environment with hot reloading"
	@echo "  build       - Build the server binary"
	@echo "  test        - Run all tests"
	@echo "  lint        - Run linter"
	@echo "  fmt         - Format code"
	@echo "  vet         - Run go vet"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-up   - Start Docker services"
	@echo "  docker-down - Stop Docker services"
	@echo "  docker-clean - Clean Docker volumes"
	@echo "  install-tools - Install development tools"
	@echo "  seed        - Run database seed data"
	@echo "  test-legacy - Test legacy endpoint compatibility"
	@echo "  logs        - Show Docker service logs"

# Set up development environment
setup:
	@echo "🚀 Setting up development environment..."
	@./scripts/setup.sh

# Start development environment
dev:
	@echo "🚀 Starting development environment..."
	@./scripts/dev.sh

# Build the server
build:
	@echo "🔨 Building server..."
	@mkdir -p bin
	@go build -o bin/server cmd/server/main.go
	@echo "✅ Build complete: bin/server"

# Run the server directly
run:
	@echo "🚀 Running server..."
	@go run cmd/server/main.go

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Run linter
lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "🔍 Running go vet..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/ tmp/* build-errors.log
	@go clean
	@echo "✅ Clean complete"

# Docker commands
docker-up:
	@echo "🐳 Starting Docker services..."
	@./scripts/docker.sh up

docker-down:
	@echo "🛑 Stopping Docker services..."
	@./scripts/docker.sh down

docker-clean:
	@echo "🧹 Cleaning Docker volumes..."
	@./scripts/docker.sh clean

# Show Docker service logs
logs:
	@echo "📋 Docker service logs:"
	@docker-compose logs -f

# Install development tools
install-tools:
	@echo "🔧 Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest
	@echo "✅ Tools installed"

# Run seed data
seed:
	@echo "🌱 Running seed data..."
	@go run cmd/tools/seed-data/main.go

# Test legacy endpoints
test-legacy:
	@echo "🧪 Testing legacy endpoints..."
	@./scripts/test-legacy-endpoints.sh
