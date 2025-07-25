# ControlMe Go Server Makefile

.PHONY: help build run test swagger docker-up docker-down clean

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the server binary"
	@echo "  run          - Run the server locally"
	@echo "  test         - Run tests"
	@echo "  swagger      - Generate Swagger documentation"
	@echo "  swagger-serve - Generate docs and run server"
	@echo "  docker-up    - Start with Docker Compose"
	@echo "  docker-down  - Stop Docker Compose"
	@echo "  clean        - Clean build artifacts"

# Build the server
build:
	go build -o server.exe ./cmd/server

# Run the server locally
run:
	go run ./cmd/server

# Run tests
test:
	go test ./...

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go -o docs/swagger
	@echo "✅ Swagger documentation generated at docs/swagger/"
	@echo "🌐 Access it at: http://localhost:8080/swagger/index.html"

# Generate swagger and run server
swagger-serve: swagger run

# Install swagger tool
install-swagger:
	go install github.com/swaggo/swag/cmd/swag@latest

# Docker commands
docker-up:
	docker-compose up --build

docker-down:
	docker-compose down -v

# Clean build artifacts
clean:
	rm -f server.exe
	rm -rf tmp/

# Development helpers
dev-setup: install-swagger
	@echo "✅ Development environment setup complete"

# Format code
fmt:
	go fmt ./...

# Lint code  
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Full development cycle
dev: swagger-serve
