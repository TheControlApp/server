#!/bin/bash

# ControlMe Go Server - Development Status
echo "🚀 ControlMe Go Development Environment Status"
echo "============================================="
echo ""

# Check Docker services
echo "📦 Docker Services:"
cd /workspace/server/controlme-go
docker-compose ps --format "table {{.Name}}\t{{.Status}}\t{{.Ports}}"
echo ""

# Check if Go server is running
echo "🔥 Go Server:"
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Server is running on http://localhost:8080"
    echo "   Health Check: $(curl -s http://localhost:8080/health)"
else
    echo "❌ Server is not running"
fi
echo ""

# Check available endpoints
echo "🌐 Available Endpoints:"
echo "   • Health: http://localhost:8080/health"
echo "   • Legacy endpoints: http://localhost:8080/GetCount.aspx?usernm=test&pwd=test&vrs=012"
echo "   • Modern API: http://localhost:8080/api/v1/users"
echo ""

# Check hot reload status
echo "🔥 Hot Reload Status:"
if pgrep -f "air" > /dev/null; then
    echo "✅ Air is running - hot reload active"
else
    echo "❌ Air is not running - no hot reload"
fi
echo ""

# Development commands
echo "🔧 Development Commands:"
echo "   • Start development: make dev"
echo "   • Build server: make build"
echo "   • Run tests: make test"
echo "   • Stop services: make docker-down"
echo ""

echo "📝 Current Focus: Legacy endpoint compatibility with hot reloading"
echo "🎯 Next Steps: Fix database schema and complete authentication system"
