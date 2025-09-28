# Test WebSocket Client Script
# This script demonstrates how to test the WebSocket connection

# First, you need to get a JWT token by registering/logging in
Write-Host "=== ControlMe WebSocket Test Client ===" -ForegroundColor Cyan
Write-Host ""

# Check if server is running
Write-Host "1. Testing if server is running..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/health" -Method GET -ErrorAction Stop
    Write-Host "✅ Server is running!" -ForegroundColor Green
} catch {
    Write-Host "❌ Server is not running. Please start the server first:" -ForegroundColor Red
    Write-Host "   cd server && go run cmd/server/main.go" -ForegroundColor Gray
    exit 1
}

Write-Host ""
Write-Host "2. To test the WebSocket client, you need a JWT token." -ForegroundColor Yellow
Write-Host "   You can get one by:" -ForegroundColor Gray
Write-Host ""
Write-Host "   a) Register a new user:" -ForegroundColor Gray
Write-Host '   Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -ContentType "application/json" -Body \'{"login_name":"testuser","email":"test@example.com","password":"testpass123","screen_name":"Test User"}\'' -ForegroundColor DarkGray
Write-Host ""
Write-Host "   b) Login with existing user:" -ForegroundColor Gray  
Write-Host '   Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body \'{"login_name":"testuser","password":"testpass123"}\'' -ForegroundColor DarkGray
Write-Host ""
Write-Host "3. Then run the test client:" -ForegroundColor Yellow
Write-Host "   .\test-client.exe -token=YOUR_JWT_TOKEN_HERE" -ForegroundColor Gray
Write-Host ""

# Try to register a test user for demonstration
Write-Host "Attempting to register a test user for you..." -ForegroundColor Yellow

$registerData = @{
    login_name = "wstestuser"
    email = "wstest@example.com" 
    password = "testpass123"
    screen_name = "WebSocket Test User"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-WebRequest -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -ContentType "application/json" -Body $registerData -ErrorAction Stop
    Write-Host "✅ Test user registered successfully!" -ForegroundColor Green
} catch {
    if ($_.Exception.Response.StatusCode -eq 409) {
        Write-Host "ℹ️  Test user already exists, trying to login..." -ForegroundColor Blue
    } else {
        Write-Host "❌ Failed to register test user: $($_.Exception.Message)" -ForegroundColor Red
    }
}

# Try to login and get token
Write-Host "Attempting to login and get JWT token..." -ForegroundColor Yellow

$loginData = @{
    login_name = "wstestuser" 
    password = "testpass123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -ContentType "application/json" -Body $loginData -ErrorAction Stop
    $token = $loginResponse.token
    
    Write-Host "✅ Login successful! Got JWT token." -ForegroundColor Green
    Write-Host ""
    Write-Host "🚀 Starting test client with token..." -ForegroundColor Cyan
    Write-Host ""
    
    # Start the test client
    & ".\test-client.exe" -token=$token
    
} catch {
    Write-Host "❌ Failed to login: $($_.Exception.Message)" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please manually get a JWT token and run:" -ForegroundColor Yellow
    Write-Host ".\test-client.exe -token=YOUR_JWT_TOKEN" -ForegroundColor Gray
}