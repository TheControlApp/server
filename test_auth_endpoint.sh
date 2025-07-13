#!/bin/bash

# ControlMe Authentication Test Script
# This script tests the authentication endpoint with encrypted passwords

set -e

echo "🔧 ControlMe Authentication Test"
echo "================================"

# Configuration
SERVER_URL="http://localhost:8080"
CRYPTO_KEY="F72F92D9A2E6F0D1B8C4E3A9F7E2D6B4"  # Default from config

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Step 1: Starting the server (if not already running)...${NC}"

# Check if server is running
if ! curl -s "$SERVER_URL/health" > /dev/null 2>&1; then
    echo -e "${YELLOW}Server not running, starting in background...${NC}"
    nohup go run cmd/api/main.go > server.log 2>&1 &
    SERVER_PID=$!
    echo "Server PID: $SERVER_PID"
    
    # Wait for server to start
    echo "Waiting for server to start..."
    for i in {1..30}; do
        if curl -s "$SERVER_URL/health" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ Server is running!${NC}"
            break
        fi
        sleep 1
        echo -n "."
    done
    
    if ! curl -s "$SERVER_URL/health" > /dev/null 2>&1; then
        echo -e "${RED}❌ Server failed to start${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✅ Server is already running${NC}"
fi

echo ""
echo -e "${BLUE}Step 2: Creating test users...${NC}"

# Create a simple test user first using direct database connection
echo "Creating test user with Go tool..."

# Create a minimal test script for user creation
cat > /tmp/create_test_user.go << 'EOF'
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

type TestUser struct {
    ID         uuid.UUID `gorm:"type:uuid;primary_key"`
    ScreenName string    `gorm:"size:50;not null"`
    LoginName  string    `gorm:"size:50;not null;unique"`
    Password   string    `gorm:"size:255;not null"`
    Role       string    `gorm:"size:50"`
    Verified   bool      `gorm:"default:false"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
}

func (u *TestUser) BeforeCreate(tx *gorm.DB) error {
    if u.ID == uuid.Nil {
        u.ID = uuid.New()
    }
    return nil
}

func main() {
    dsn := "host=localhost user=postgres password=postgres dbname=controlme port=5432 sslmode=disable TimeZone=UTC"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }

    // Auto-migrate
    if err := db.AutoMigrate(&TestUser{}); err != nil {
        log.Fatalf("Failed to migrate: %v", err)
    }

    // Hash password
    hash, err := bcrypt.GenerateFromPassword([]byte("test123"), bcrypt.DefaultCost)
    if err != nil {
        log.Fatalf("Failed to hash: %v", err)
    }

    // Create user
    user := TestUser{
        ID:         uuid.New(),
        ScreenName: "testuser",
        LoginName:  "testuser", 
        Password:   string(hash),
        Role:       "user",
        Verified:   true,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    if err := db.Create(&user).Error; err != nil {
        log.Printf("User might already exist: %v", err)
    } else {
        fmt.Println("✅ Created test user: testuser/test123")
    }
}
EOF

# Run the user creation
cd /workspace/server
echo "Running user creation script..."
go run /tmp/create_test_user.go

echo ""
echo -e "${BLUE}Step 3: Encrypting password for testing...${NC}"

# Encrypt the password using crypto tool
echo "Encrypting password 'test123'..."
encrypted_password=$(go run cmd/tools/crypto-test/main.go --password "test123" | grep "Encrypted:" | cut -d' ' -f2)

if [ -z "$encrypted_password" ]; then
    echo -e "${RED}❌ Failed to encrypt password${NC}"
    exit 1
fi

echo -e "${GREEN}✅ Encrypted password: $encrypted_password${NC}"

echo ""
echo -e "${BLUE}Step 4: Testing authentication endpoint...${NC}"

# Test the authentication endpoint
echo "Testing /TestAuth.aspx with encrypted password..."
echo ""

test_url="$SERVER_URL/TestAuth.aspx?usernm=testuser&pwd=${encrypted_password}&vrs=012"

echo -e "${YELLOW}Testing URL:${NC}"
echo "$test_url"
echo ""

echo -e "${YELLOW}Response:${NC}"
echo "----------------------------------------"
curl -s "$test_url"
echo ""
echo "----------------------------------------"

echo ""
echo -e "${BLUE}Step 5: Testing other legacy endpoints...${NC}"

echo ""
echo -e "${YELLOW}Testing GetCount.aspx:${NC}"
curl -s "$SERVER_URL/GetCount.aspx?usernm=testuser&pwd=${encrypted_password}&vrs=012"

echo ""
echo -e "${YELLOW}Testing GetContent.aspx:${NC}"  
curl -s "$SERVER_URL/GetContent.aspx?usernm=testuser&pwd=${encrypted_password}&vrs=012"

echo ""
echo -e "${YELLOW}Testing AppCommand.aspx (Outstanding):${NC}"
curl -s "$SERVER_URL/AppCommand.aspx?usernm=testuser&pwd=${encrypted_password}&vrs=012&cmd=Outstanding"

echo ""
echo ""
echo -e "${GREEN}🎉 Authentication test complete!${NC}"
echo ""
echo -e "${YELLOW}Manual testing URLs:${NC}"
echo "Test Auth: $SERVER_URL/TestAuth.aspx?usernm=testuser&pwd=${encrypted_password}&vrs=012"
echo "Get Count: $SERVER_URL/GetCount.aspx?usernm=testuser&pwd=${encrypted_password}&vrs=012"
echo "Get Content: $SERVER_URL/GetContent.aspx?usernm=testuser&pwd=${encrypted_password}&vrs=012"
echo ""
echo -e "${YELLOW}Server health: $SERVER_URL/health${NC}"

# Clean up temporary file
rm -f /tmp/create_test_user.go

echo ""
echo -e "${GREEN}✅ Test script completed successfully!${NC}"
