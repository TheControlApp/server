#!/bin/bash

# ControlMe Go Server Development Startup Script
# This script starts the necessary services and runs the Go server with hot reloading

set -e

echo "🚀 Starting ControlMe Go Development Environment"

# Change to the project directory
cd /workspace/server/controlme-go

# Function to cleanup on exit
cleanup() {
    echo "🧹 Cleaning up..."
    docker-compose down
    exit 0
}

# Register cleanup function
trap cleanup EXIT INT TERM

# Start Docker services
echo "🐳 Starting Docker services (PostgreSQL, Redis)..."
docker-compose up -d postgres redis

# Wait for services to be ready
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check if services are healthy
echo "🔍 Checking service health..."
docker-compose ps

# Create tmp directory for Air
mkdir -p tmp

# Start the Go server with Air for hot reloading
echo "🔥 Starting Go server with Air hot reloading..."
echo "📝 Server will be available at: http://localhost:8080"
echo "💾 Database (PostgreSQL): localhost:5432"
echo "🔄 Redis: localhost:6379"
echo ""
echo "🎯 Press Ctrl+C to stop all services"
echo ""

# Start Air
air
