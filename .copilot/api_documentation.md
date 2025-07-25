# ControlMe API Documentation

## Overview
The ControlMe API provides a WebSocket-first architecture for real-time command delivery with a minimal REST API for authentication, file management, and profile operations.

## Authentication
All API endpoints require JWT authentication via Bearer token:
```
Authorization: Bearer <jwt_token>
```

## REST API Endpoints

### Authentication
```
POST /api/v1/auth/login
POST /api/v1/auth/register
```

### User Management
```
GET /api/v1/users              # List users
GET /api/v1/users/:id          # Get specific user
PUT /api/v1/users/profile      # Update profile
```

### File Management

#### File Upload
```
POST /api/v1/files
Content-Type: multipart/form-data
Authorization: Bearer <jwt_token>

Form Data:
- file: <binary_data>
```

**Success Response:**
```json
{
  "status": "success",
  "file_hash": "abc123def456..."
}
```
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

#### File Download
```
GET /api/v1/files?filehash=<hash>&filename=<desired_filename>
Authorization: Bearer <jwt_token>

Response:
- Content-Type: <original_mime_type>
- Content-Disposition: attachment; filename="<desired_filename>"
- <binary_file_data>
```

### Health Check
```
GET /health                    # Server health status
```

## WebSocket API

### Connection
```
WS /api/ws
Authorization: Bearer <jwt_token> (via header or query param)
```

### Message Types

#### Command Reception (Server → Client)
```json
{
  "type": "command",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "Complete your daily task",
          "button": "Got it!"
        }
      },
      {
        "type": "download-file",
        "content": {
          "file_hash": "abc123def456...",
          "file_name": "task.pdf"
        }
      }
    ],
    "sender": "master123", // sender's username
    "receiver": "user456", // receiver's username (null for broadcasts)
    "tags": ["chastity", "daily"],
    "status": "pending",
    "created_at": "2025-07-24T15:30:00Z"
  }
}
```

#### Command Sending (Client → Server)
```json
{
  "type": "send_command",
  "data": {
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "Take a 15-minute break",
          "button": "Done"
        }
      }
    ],
    "receiver": "user456", // optional - omit for broadcast
    "tags": ["general", "break"]
  }
}
```

#### Broadcast Command (No specific receiver)
```json
{
  "type": "send_command",
  "data": {
    "instructions": [
      {
        "type": "announcement",
        "content": {
          "title": "System Maintenance",
          "body": "Server will be down for 10 minutes at 3 PM",
          "priority": "high"
        }
      }
    ],
    "tags": ["announcements", "system"]
  }
}
```

#### Command Status Update
```json
{
  "type": "command_status",
  "data": {
    "command_id": "uuid",
    "status": "completed",
    "completed_at": "2025-07-24T15:45:00Z"
  }
}
```

#### Heartbeat
```json
{
  "type": "heartbeat",
  "timestamp": "2025-07-24T15:30:00Z"
}
```

## Instruction Types

### Standard Types
- **`popup-msg`**: Show popup dialog with body and button text
- **`popup-web`**: Open URL in browser/webview
- **`popup-video`**: Open video player with URL
- **`download-file`**: Trigger file download

### Custom Types (Extensible)
- **`custom-vibration`**: Hardware vibration patterns
- **`custom-iot`**: IoT device control
- **`custom-ar`**: Augmented reality instructions
- **`custom-biometric`**: Biometric sensor readings
- **`custom-game`**: Game-specific commands
- **`custom-notification`**: Platform-specific notifications
- **`custom-camera`**: Camera control instructions
- **`custom-gps`**: GPS/location instructions

### Instruction Content Examples

#### Popup Message
```json
{
  "type": "popup-msg",
  "content": {
    "body": "Please complete your daily task",
    "button": "Got it!",
    "priority": "normal"
  }
}
```

#### File Download
```json
{
  "type": "download-file",
  "content": {
    "file_hash": "abc123def456...",
    "file_name": "Daily Task Instructions.pdf"
  }
}
```

#### Custom Vibration
```json
{
  "type": "custom-vibration",
  "content": {
    "pattern": [500, 200, 500],
    "intensity": 75,
    "duration": 2000
  }
}
```

## Error Handling

### WebSocket Errors
```json
{
  "type": "error",
  "error": "authentication_failed",
  "message": "Invalid or expired token"
}

{
  "type": "error",
  "error": "command_failed",
  "message": "Unable to deliver command to target user",
  "command_id": "uuid"
}
```

### Connection Management
- **Heartbeat Interval**: 30 seconds
- **Missed Heartbeat Limit**: 3 consecutive misses
- **Automatic Reconnection**: Client responsibility
- **Queue Delivery**: Commands queued for offline users

## Security Features

### File Security
- CSAM detection on all uploads
- Virus scanning before storage
- Hash-based deduplication
- 2-week retention policy
- Authentication required for all file operations

### Connection Security
- JWT token authentication
- Heartbeat monitoring
- Automatic disconnection of inactive connections
- Comprehensive error logging
- User-controlled content filtering via tags

## Rate Limiting
- Commands: 10 per minute per user
- File uploads: 5 per minute per user
- File downloads: 20 per minute per user
- WebSocket connections: 1 per user

## Content Filtering

### Tag-Based System
Users can opt-in/out of content categories:
- `general` - Safe-for-work content
- `adult` - Explicit content
- `chastity` - Chastity-related content
- `feet` - Foot-related content
- `roleplay` - Roleplay scenarios
- `humiliation` - Humiliation-based content
- `public` - Public setting commands
- `private` - Private/intimate commands

### Broadcasting Rules
- Commands with `receiver_id` go to specific user
- Commands without `receiver_id` broadcast to users subscribed to matching tags
- Users can block specific senders via Block system
- Reported users may have delivery restrictions
