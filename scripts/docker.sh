#!/bin/bash

# Docker management script for ControlMe Go

set -e

COMPOSE_FILE="docker-compose.yml"

case "${1:-help}" in
    "up")
        echo "🚀 Starting ControlMe database services..."
        docker-compose up -d postgres
        echo "⏳ Waiting for services to be ready..."
        sleep 5
        echo "📊 Database services started successfully!"
        echo "  PostgreSQL: localhost:5432 (controlme/postgres/postgres)"
        ;;
    
    "down")
        echo "🛑 Stopping ControlMe services..."
        docker-compose down
        echo "✅ Services stopped"
        ;;
    
    "restart")
        echo "🔄 Restarting ControlMe services..."
        docker-compose down
        docker-compose up -d postgres
        echo "✅ Services restarted"
        ;;
     "logs")
        echo "📋 Showing service logs..."
        docker-compose logs -f postgres
        ;;
    
    "status")
        echo "📊 Service status:"
        docker-compose ps
        ;;
    
    "clean")
        echo "🧹 Cleaning up containers and volumes..."
        docker-compose down -v --remove-orphans
        echo "✅ Cleanup complete"
        ;;
    
    "reset")
        echo "🔄 Resetting database (WARNING: This will delete all data!)..."
        read -p "Are you sure? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            docker-compose down -v
            docker volume rm controlme-go_postgres_data 2>/dev/null || true
            docker-compose up -d postgres
            echo "✅ Database reset complete"
        else
            echo "❌ Reset cancelled"
        fi
        ;;
    
    "psql")
        echo "🐘 Connecting to PostgreSQL..."
        docker-compose exec postgres psql -U postgres -d controlme
        ;;
    
    "help"|*)
        echo "ControlMe Docker Management"
        echo ""
        echo "Usage: $0 <command>"
        echo ""
        echo "Commands:"
        echo "  up       Start database services (PostgreSQL)"
        echo "  down     Stop all services"
        echo "  restart  Restart all services"
        echo "  logs     Show service logs"
        echo "  status   Show service status"
        echo "  clean    Stop services and remove volumes"
        echo "  reset    Reset database (deletes all data!)"
        echo "  psql     Connect to PostgreSQL"
        echo "  help     Show this help message"
        ;;
esac
