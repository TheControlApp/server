#!/bin/bash

# Generate Swagger/OpenAPI documentation

set -e

echo "🔥 Generating Swagger documentation..."

# Ensure swag is installed
if ! command -v swag &> /dev/null; then
    echo "❌ swag is not installed. Please install it first:"
    echo "  go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

# Generate the docs
swag init -g cmd/server/main.go

echo "✅ Swagger documentation generated successfully in the 'docs' directory."