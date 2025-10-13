# REST API Reference

## Overview

The ControlApp Server provides a RESTful API for user management, authentication, and command operations. All API endpoints are prefixed with `/api/v1/` and return JSON responses following RFC 7807 Problem Details for HTTP APIs.

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

The API uses JWT (JSON Web Token) based authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Error Handling

All error responses follow RFC 7807 Problem Details standard:

```json
{
  "type": "string",           // Error type identifier
  "title": "string",          // Human-readable error title  
  "status": 400,              // HTTP status code
  "detail": "string",         // Specific error description
  "action": "string",         // Suggested action (optional)
  "help": "string"            // Additional guidance (optional)
}
```

### Common Error Types

- `bad_request` (400) - Request is malformed or invalid
- `unauthorized` (401) - Authentication required or failed
- `forbidden` (403) - Access denied for authenticated user
- `not_found` (404) - Requested resource not found
- `conflict` (409) - Resource conflict (e.g., username already exists)
- `validation_error` (422) - Input validation failed
- `internal_server_error` (500) - Server encountered an error

### Validation Error Format

```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation",
  "help": "Please check the field requirements in the API documentation",
  "errors": [
    {
      "field": "username",
      "message": "Username must be at least 3 characters long",
      "code": "MIN_LENGTH"
    }
  ]
}
```

## Endpoints

### Health Check

#### GET /health

Check if the server is running.

**Response (200 OK):**
```json
{
  "status": "ok",
  "message": "Server is running"
}
```

---

## Authentication Endpoints

### Register User

#### POST /auth/register

Create a new user account.

**Request Body:**
```json
{
  "username": "string",      // Required, min 3 chars, must be unique
  "password": "string",      // Required, min 6 chars
  "screen_name": "string",   // Required, display name
  "random_opt_in": boolean   // Optional, default false
}
```

**Success Response (201 Created):**
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "login_name": "testuser",
    "screen_name": "Test User",
    "role": "user",
    "random_opt_in": false,
    "anon_cmd": false,
    "verified": false,
    "verified_code": 0,
    "thumbs_up": 0,
    "created_at": "2025-10-12T18:00:00Z",
    "updated_at": "2025-10-12T18:00:00Z",
    "login_date": "2025-10-12T18:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON or missing required fields
- `409 Conflict`: Username already exists
- `422 Validation Error`: Field validation failed
- `500 Internal Server Error`: Server error during registration

### Login User

#### POST /auth/login

Authenticate an existing user and receive a JWT token.

**Request Body:**
```json
{
  "username": "string",      // Required
  "password": "string"       // Required
}
```

**Success Response (200 OK):**
```json
{
  "message": "Login successful",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "login_name": "testuser",
    "screen_name": "Test User",
    "role": "user",
    "random_opt_in": false,
    "anon_cmd": false,
    "verified": false,
    "verified_code": 0,
    "thumbs_up": 0,
    "created_at": "2025-10-12T18:00:00Z",
    "updated_at": "2025-10-12T18:00:00Z",
    "login_date": "2025-10-12T18:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON or missing required fields
- `401 Unauthorized`: Invalid username or password
- `500 Internal Server Error`: Server error during authentication

---

## User Management Endpoints

### Get All Users

#### GET /users

Retrieve a list of all users in the system.

**Success Response (200 OK):**
```json
{
  "users": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "login_name": "testuser",
      "screen_name": "Test User",
      "role": "user",
      "random_opt_in": false,
      "anon_cmd": false,
      "verified": false,
      "verified_code": 0,
      "thumbs_up": 0,
      "created_at": "2025-10-12T18:00:00Z",
      "updated_at": "2025-10-12T18:00:00Z",
      "login_date": "2025-10-12T18:00:00Z"
    }
  ]
}
```

**Error Responses:**
- `500 Internal Server Error`: Failed to fetch users

### Get User by ID

#### GET /users/{id}

Retrieve a specific user by their UUID.

**Path Parameters:**
- `id` (string, required): User UUID

**Success Response (200 OK):**
```json
{
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "login_name": "testuser",
    "screen_name": "Test User",
    "role": "user",
    "random_opt_in": false,
    "anon_cmd": false,
    "verified": false,
    "verified_code": 0,
    "thumbs_up": 0,
    "created_at": "2025-10-12T18:00:00Z",
    "updated_at": "2025-10-12T18:00:00Z",
    "login_date": "2025-10-12T18:00:00Z"
  }
}
```

**Error Responses:**
- `400 Bad Request`: Invalid UUID format
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

---

## Command Management Endpoints

### Get Pending Commands

#### GET /commands/pending

Retrieve pending commands for a specific user.

**Query Parameters:**
- `user_id` (string, required): User UUID

**Success Response (200 OK):**
```json
{
  "commands": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "instructions": [
        {
          "type": "std_popup",
          "content": {
            "body": "This is a test command",
            "button": "OK"
          }
        }
      ],
      "sender_id": "123e4567-e89b-12d3-a456-426614174001",
      "receiver_id": "123e4567-e89b-12d3-a456-426614174000",
      "tags": "",
      "status": "pending",
      "created_at": "2025-10-12T18:00:00Z",
      "updated_at": "2025-10-12T18:00:00Z",
      "sender": {
        "id": "123e4567-e89b-12d3-a456-426614174001",
        "login_name": "sender",
        "screen_name": "Command Sender"
      },
      "receiver": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "login_name": "receiver",
        "screen_name": "Command Receiver"
      }
    }
  ]
}
```

**Error Responses:**
- `400 Bad Request`: Missing or invalid user_id parameter
- `422 Validation Error`: Invalid UUID format
- `500 Internal Server Error`: Failed to fetch commands

### Complete Command

#### POST /commands/complete

Mark a specific command as completed.

**Query Parameters:**
- `user_id` (string, required): User UUID
- `command_id` (string, required): Command UUID

**Success Response (200 OK):**
```json
{
  "message": "Command completed successfully"
}
```

**Error Responses:**
- `400 Bad Request`: Missing required parameters
- `422 Validation Error`: Invalid UUID format
- `500 Internal Server Error`: Failed to complete command

---

## Testing with Provided JSON Files

The server includes pre-built test JSON files for easy API testing:

### Test Files Available

#### `api_test_registration.json`
```json
{
  "username": "api_test_user",
  "screen_name": "API Test User", 
  "password": "apitest123",
  "random_opt_in": false
}
```

#### `api_test_login.json`
```json
{
  "username": "api_test_user",
  "password": "apitest123"
}
```

#### `test_empty_registration.json`
```json
{}
```

### Using Test Files

#### Registration Test
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d @api_test_registration.json
```

#### Login Test
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d @api_test_login.json
```

#### Validation Error Test
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d @test_empty_registration.json
```

#### Extract JWT Token for Testing
```bash
# Save token to environment variable
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d @api_test_login.json | jq -r '.token')

# Use token in authenticated requests
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/v1/protected-endpoint
```

### PowerShell Testing
```powershell
# Registration
$body = Get-Content "api_test_registration.json" -Raw
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/register" -Method POST -Body $body -ContentType "application/json"

# Login
$body = Get-Content "api_test_login.json" -Raw
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/v1/auth/login" -Method POST -Body $body -ContentType "application/json"
$token = $response.token
```

---

## Data Models

### User Model

```json
{
  "id": "string (UUID)",           // Unique identifier
  "login_name": "string",          // Username for login
  "screen_name": "string",         // Display name
  "role": "string",                // User role (default: "user")
  "random_opt_in": "boolean",      // Opted in for random commands
  "anon_cmd": "boolean",           // Allow anonymous commands
  "verified": "boolean",           // Account verification status
  "verified_code": "integer",      // Verification code
  "thumbs_up": "integer",          // User rating
  "created_at": "string (ISO8601)", // Account creation time
  "updated_at": "string (ISO8601)", // Last update time
  "login_date": "string (ISO8601)" // Last login time
}
```

### Command Model

```json
{
  "id": "string (UUID)",           // Unique identifier
  "instructions": "array",         // Array of instruction objects
  "sender_id": "string (UUID)",    // User who sent the command
  "receiver_id": "string (UUID)",  // Target user (optional)
  "tags": "string",                // JSON array of tag names
  "status": "string",              // pending, delivered, completed
  "created_at": "string (ISO8601)", // Creation time
  "updated_at": "string (ISO8601)", // Last update time
  "sender": "User",                // Sender user object
  "receiver": "User"               // Receiver user object (optional)
}
```

### Instruction Model

```json
{
  "type": "string",               // Instruction type identifier
  "content": "object"             // Arbitrary instruction data
}
```

## Rate Limiting

Currently, no rate limiting is implemented. This may be added in future versions.

## Versioning

The API uses URL versioning. Current version is `v1`. Future versions will be available at `/api/v2`, etc.

## CORS

CORS is configured to allow requests from:
- `http://localhost:3000`
- `http://localhost:8080` 
- `http://127.0.0.1:3000`
- `http://127.0.0.1:8080`

Additional origins can be configured in the server settings.