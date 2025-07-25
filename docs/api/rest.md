# REST API

All endpoints require JWT authentication: `Authorization: Bearer <token>`

## Authentication

### Login
```http
POST /api/v1/auth/login
```
**Body:** 
```json
{
  "login_name": "username",
  "password": "password"
}
```
**Returns:** 
```json
{
  "token": "jwt_token",
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
**Body:** 
```json
{
  "screen_name": "Display Name",
  "login_name": "username",
  "email": "user@example.com",
  "password": "password"
}
```
**Returns:** Success/failure message

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
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |

## Content Types

Supported file types for upload:
- Images: jpg, png, gif, webp (max 10MB)
- Documents: pdf, txt, doc, docx (max 50MB)
- Videos: mp4, webm, mov (max 100MB)
- Archives: zip, rar, 7z (max 50MB)
