#!/bin/bash

# Development setup script for ControlMe Go

set -e

echo "🚀 Setting up ControlMe Go development environment..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
if [[ $(echo "$GO_VERSION 1.21" | awk '{print ($1 < $2)}') -eq 1 ]]; then
    echo "❌ Go version $GO_VERSION is too old. Please install Go 1.21+ first."
    exit 1
fi

# Create necessary directories
echo "📁 Creating directories..."
mkdir -p bin logs tmp

# Install Go dependencies
echo "� Installing Go dependencies..."
go mod tidy

# Start Docker services
echo "🐳 Starting Docker services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 5

# Run database migrations/setup
echo "� Setting up database..."
go run cmd/tools/seed-data/main.go

# Build the server
echo "🔨 Building server..."
go build -o bin/server cmd/server/main.go

echo "✅ Development environment setup complete!"
echo ""
echo "🚀 To start the server:"
echo "  make dev  # or ./bin/server"
echo ""
echo "🧪 To run tests:"
echo "  make test"
echo ""
echo "🌐 API endpoints will be available at:"
echo "  http://localhost:8080/health"
echo "  http://localhost:8080/api/v1/..."
echo "  http://localhost:8080/AppCommand.aspx (legacy)"
