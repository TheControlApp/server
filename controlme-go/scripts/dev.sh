#!/bin/bash

# Development script for ControlMe Go

set -e

echo "🚀 Starting ControlMe Go development environment..."

# Start Docker services
echo "📊 Starting database services..."
./scripts/docker.sh up

# Wait for services to be ready
echo "⏳ Waiting for database to be ready..."
sleep 3

# Check if database is ready
until docker-compose exec -T postgres pg_isready -U postgres -d controlme > /dev/null 2>&1; do
    echo "Waiting for PostgreSQL..."
    sleep 2
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
