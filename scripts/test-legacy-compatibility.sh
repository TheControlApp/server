#!/bin/bash

# Test script to verify legacy endpoint compatibility
# This script tests the updated Go handlers against the legacy ASP.NET format

echo "Testing ControlMe Legacy API Compatibility"
echo "=========================================="

# Configuration
SERVER_URL="http://localhost:8080"
TEST_USER="testuser"
TEST_PASS="testpass"
VERSION="012"

echo "Testing AppCommand endpoint..."

# Test Outstanding command
echo "1. Testing Outstanding command..."
curl -s "${SERVER_URL}/AppCommand.aspx?usernm=${TEST_USER}&pwd=${TEST_PASS}&vrs=${VERSION}&cmd=Outstanding" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  | grep -o '\[.*\]' || echo "FAIL: Outstanding command"

# Test Content command  
echo "2. Testing Content command..."
curl -s "${SERVER_URL}/AppCommand.aspx?usernm=${TEST_USER}&pwd=${TEST_PASS}&vrs=${VERSION}&cmd=Content" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  | grep -o '\[.*\]' || echo "FAIL: Content command"

# Test GetContent endpoint
echo "3. Testing GetContent endpoint..."
curl -s "${SERVER_URL}/GetContent.aspx?usernm=${TEST_USER}&pwd=${TEST_PASS}&vrs=${VERSION}" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  | grep -o '<asp:Label.*>.*</asp:Label>' || echo "FAIL: GetContent endpoint"

# Test GetCount endpoint
echo "4. Testing GetCount endpoint..."
curl -s "${SERVER_URL}/GetCount.aspx?usernm=${TEST_USER}&pwd=${TEST_PASS}&vrs=${VERSION}" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  | grep -o '<asp:Label.*>.*</asp:Label>' || echo "FAIL: GetCount endpoint"

# Test version validation
echo "5. Testing version validation..."
curl -s "${SERVER_URL}/AppCommand.aspx?usernm=${TEST_USER}&pwd=${TEST_PASS}&vrs=999&cmd=Outstanding" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  | grep -o 'Wrong version' || echo "FAIL: Version validation"

echo "Legacy API compatibility test complete!"
echo "Check the output above for any FAIL messages."
