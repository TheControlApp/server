#!/bin/bash

echo "=== Testing ControlMe Go Legacy Endpoints ==="
echo

# Test health endpoint first
echo "1. Testing health endpoint..."
curl -s http://localhost:8080/health | head -c 200
echo
echo

# Test GetCount.aspx with legacy auth
echo "2. Testing GetCount.aspx (should fail with missing parameters)..."
curl -s http://localhost:8080/GetCount.aspx | head -c 200
echo
echo

# Test GetOptions.aspx with legacy auth
echo "3. Testing GetOptions.aspx (should fail with missing parameters)..."
curl -s http://localhost:8080/GetOptions.aspx | head -c 200
echo
echo

# Test with valid user parameters (need encrypted password)
echo "4. Testing GetCount.aspx with parameters (will fail without valid encrypted password)..."
curl -s "http://localhost:8080/GetCount.aspx?User=testdom&Password=invalid" | head -c 200
echo
echo

# Test AppCommand.aspx
echo "5. Testing AppCommand.aspx (GET request)..."
curl -s http://localhost:8080/AppCommand.aspx | head -c 200
echo
echo

# Test GetContent.aspx
echo "6. Testing GetContent.aspx..."
curl -s http://localhost:8080/GetContent.aspx | head -c 200
echo
echo

# Test legacy web pages
echo "7. Testing Default.aspx..."
curl -s http://localhost:8080/Default.aspx | head -c 200
echo
echo

echo "8. Testing Messages.aspx..."
curl -s http://localhost:8080/Messages.aspx | head -c 200
echo
echo

echo "=== Legacy Endpoint Tests Complete ==="

# Create encrypted password for testing
echo
echo "=== Creating test command for encrypted password ==="
echo "To test with real authentication, we need to:"
echo "1. Get an encrypted password using the legacy crypto"
echo "2. Use that encrypted password in the requests"
echo
echo "Run this Go code to get encrypted password for 'testpass1':"
echo 'package main'
echo 'import ('
echo '    "fmt"'
echo '    "github.com/thecontrolapp/controlme-go/internal/auth"'
echo ')'
echo 'func main() {'
echo '    crypto := auth.NewLegacyCrypto("default-key-32-bytes-long-for-aes")'
echo '    encrypted, err := crypto.Encrypt("testpass1")'
echo '    if err != nil { panic(err) }'
echo '    fmt.Println("Encrypted password:", encrypted)'
echo '}'
