#!/bin/bash

# Development script for ControlMe Go

set -e

echo "🚀 Starting ControlMe Go development environment..."

# Start Docker services
echo "� Starting Docker services..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if database is ready
echo "📊 Checking database connection..."
until docker-compose exec -T postgres pg_isready -U controlme -d controlme > /dev/null 2>&1; do
    echo "Waiting for PostgreSQL..."
    sleep 2
done

# Build the server
echo "🔨 Building server..."
go build -o bin/server cmd/server/main.go

# Start the server with auto-reload if air is available
if command -v air &> /dev/null; then
    echo "🔥 Starting server with hot reload..."
    air
else
    echo "🔥 Starting server..."
    ./bin/server
fi
done

echo "✅ Database is ready!"

# Build the server
echo "🔨 Building server..."
go build -o bin/server cmd/server/main.go

# Run the server
echo "🎯 Starting ControlMe Go server..."
echo "📡 Server will be available at: http://localhost:8080"
echo "🔗 Health check: http://localhost:8080/health"
echo "📱 Legacy endpoints: http://localhost:8080/AppCommand.aspx"
echo ""
echo "Press Ctrl+C to stop the server..."
echo ""

# Start the server
./bin/server
