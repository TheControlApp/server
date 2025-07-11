#!/bin/bash

# Development setup script for ControlMe Go

echo "🚀 Setting up ControlMe Go development environment..."

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "❌ PostgreSQL is not installed. Please install PostgreSQL first."
    exit 1
fi

# Check if PostgreSQL is running
if ! pg_isready -h localhost -p 5432 &> /dev/null; then
    echo "❌ PostgreSQL is not running. Please start PostgreSQL first."
    exit 1
fi

# Create database if it doesn't exist
echo "📊 Creating database..."
sudo -u postgres psql -c "CREATE DATABASE controlme;" 2>/dev/null || echo "Database already exists"

# Run go mod tidy
echo "📦 Installing Go dependencies..."
go mod tidy

# Build the server
echo "🔨 Building server..."
go build -o bin/server cmd/server/main.go

# Create bin directory if it doesn't exist
mkdir -p bin

echo "✅ Development environment setup complete!"
echo ""
echo "To start the server:"
echo "  ./bin/server"
echo ""
echo "To run tests:"
echo "  go test ./..."
echo ""
echo "API endpoints will be available at:"
echo "  http://localhost:8080/health"
echo "  http://localhost:8080/api/v1/..."
echo "  http://localhost:8080/AppCommand.aspx (legacy)"
