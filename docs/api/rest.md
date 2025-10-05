# REST API

All endpoints require JWT authentication: `Authorization: Bearer <token>`

**🔗 Code References:**
- [API Routes Definition](../../internal/api/routes/routes.go) - Complete endpoint mapping
- [Auth Handlers](../../internal/api/handlers/auth_handlers.go) - Authentication implementation
- [User Handlers](../../internal/api/handlers/user_handlers.go) - User management endpoints
- [Command Handlers](../../internal/api/handlers/command_handlers.go) - Command API endpoints
- [Models Definition](../../internal/models/models.go) - Data structures and schemas

## Authentication

**Implementation:** [Auth Handlers](../../internal/api/handlers/auth_handlers.go)

### Login
```http
POST /api/v1/auth/login
```
**Implementation:** `Login()` function in auth_handlers.go
**Body:** 
```json
{
  "username": "username",
  "password": "password"
}
```
**Returns:** 
```json
{
  "token": "jwt_token",
  "message": "Login successful",
  "user": {
    "id": "uuid",
    "screen_name": "Display Name",
    "login_name": "username",
    "email": "user@example.com"
  }
}
```

### Register  
```http
POST /api/v1/auth/register
```
**Implementation:** `Register()` function in auth_handlers.go
**Body:** 
```json
{
  "screen_name": "Display Name",
  "username": "username", 
  "email": "user@example.com",
  "password": "password"
}
```
**Returns:** 
```json
{
  "user": {
    "id": "uuid",
    "screen_name": "Display Name",
    "login_name": "username",
    "email": "user@example.com",
    "role": "user",
    "random_opt_in": false,
    "anon_cmd": false,
    "verified": false,
    "verified_code": 0,
    "thumbs_up": 0,
    "created_at": "2025-10-05T11:02:40.3744811-07:00",
    "updated_at": "2025-10-05T11:02:40.3744811-07:00",
    "login_date": "2025-10-05T11:02:40.3744811-07:00"
  }
}
```

## Users

### List Users
```http
GET /api/v1/users
```
**Returns:** Array of user objects  
**Note:** Auth required but not yet implemented

## Files

### Upload
```http
POST /api/v1/files/upload
```
**Body:** Multipart form data with file field  
**Returns:** 
```json
{
  "file_hash": "abc123..."
}
```

### Download
```http  
GET /api/v1/files/download/{hash}
```
**Returns:** File data with appropriate headers

## Health Check
```http
GET /health
```
**Returns:** 
```json
{
  "status": "ok"
}
```
  "login_name": "username",
  "email": "user@example.com",
  "password": "password"
}
```

**Success Response:**
```json
{
  "message": "User created successfully",
  "user": {
    "id": "uuid",
    "username": "username",
    "display_name": "Display Name",
    "email": "user@example.com"
  }
}
```

**Error Responses:**
```json
// Invalid input
{
  "message": "Invalid username format. Must be 3-50 characters, letters/numbers/underscore/hyphen only"
}

// Username taken
{
  "message": "Username 'username' is already taken"
}

// Email taken
{
  "message": "Email 'user@example.com' is already registered"
}

// Password too weak
{
  "message": "Password must be at least 8 characters long"
}
```

### User Management

**Note: User endpoints require authentication.**

#### List Users
```http
GET /api/v1/users
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "users": [
    {
      "id": "uuid",
      "screen_name": "Display Name",
      "login_name": "username",
      "verified": true
    }
  ]
}
```

#### Get User
```http
GET /api/v1/users/:id
Authorization: Bearer <jwt_token>
```

#### Update Profile
```http
PUT /api/v1/users/profile
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "screen_name": "New Display Name",
  "email": "new@example.com"
}
```

### File Management

#### Upload File
```http
POST /api/v1/files
Content-Type: multipart/form-data
Authorization: Bearer <jwt_token>
```

**Form Data:**
- `file`: Binary file data

**Success Response:**
```json
{
  "status": "success",
  "file_hash": "abc123def456..."
}
```

**Error Responses:**
```json
{
  "status": "error",
  "error": "file_too_large",
  "message": "File exceeds maximum size of 50MB"
}

{
  "status": "error", 
  "error": "file_banned",
  "message": "File hash has been flagged and cannot be uploaded"
}

{
  "status": "error",
  "error": "unauthorized",
  "message": "Invalid or expired authentication token"
}
```

#### Download File
```http
GET /api/v1/files?filehash=<hash>&filename=<desired_filename>
Authorization: Bearer <jwt_token>
```

**Response:**
- Binary file data with appropriate headers
- Content-Disposition: attachment; filename="<desired_filename>"

### Health Check

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "timestamp": "2025-07-24T15:30:00Z"
}
```

## Rate Limits

- Commands: 10 per minute per user
- File uploads: 5 per minute per user  
- File downloads: 20 per minute per user
- Authentication: 5 attempts per minute per IP

## Error Codes

| Code | Description |
|------|-------------|
| 400 | Bad Request - Invalid request format |
| 401 | Unauthorized - Invalid or missing token |
| 403 | Forbidden - Insufficient permissions |
| 404 | Not Found - Resource does not exist |
| 413 | Payload Too Large - File size exceeded |

## 🛠️ Client Development Guide

### Getting Started
1. **Authentication Flow:** Register user → Login → Store JWT token
2. **API Integration:** Use REST endpoints for user/command management  
3. **Real-time Communication:** Connect to WebSocket with JWT token
4. **Command Handling:** Implement handlers for [standard command types](../standards/command_types.md)

### Code Examples

#### JavaScript Authentication
```javascript
// Register new user
const registerResponse = await fetch('/api/v1/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    screen_name: 'Display Name',
    login_name: 'username',
    email: 'user@example.com', 
    password: 'password'
  })
});

// Login and get token
const loginResponse = await fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    login_name: 'username',
    password: 'password'
  })
});

const { token, user } = await loginResponse.json();
localStorage.setItem('jwt_token', token);
```

#### Python Authentication
```python
import requests
import json

# Register user
register_data = {
    "screen_name": "Display Name",
    "login_name": "username", 
    "email": "user@example.com",
    "password": "password"
}
register_response = requests.post('http://localhost:8080/api/v1/auth/register', 
                                 json=register_data)

# Login and get token
login_data = {"login_name": "username", "password": "password"}
login_response = requests.post('http://localhost:8080/api/v1/auth/login', 
                              json=login_data)
token = login_response.json()['token']

# Store token for subsequent requests
headers = {'Authorization': f'Bearer {token}'}
```

### WebSocket Integration
See [WebSocket API documentation](websocket.md) for real-time communication setup.

### Complete Example
For a full working client implementation, see [Mini Client](../../client/mini-client.html) - a complete HTML/JavaScript client that demonstrates all core functionality.
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |

## Content Types

Supported file types for upload:
- Images: jpg, png, gif, webp (max 10MB)
- Documents: pdf, txt, doc, docx (max 50MB)
- Videos: mp4, webm, mov (max 100MB)
- Archives: zip, rar, 7z (max 50MB)
