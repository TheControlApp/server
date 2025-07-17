# WebSocket-First Architecture Plan

## Core Design Philosophy
- **WebSocket Primary**: Most application functionality flows through the single `/api/ws` endpoint
- **REST API Minimal**: Only for specific needs like file operations, health checks, profile management
- **Real-time First**: Commands, chat, status updates all happen via WebSocket for immediate responsiveness
- **Persistent Queues**: Commands persist for offline users and are delivered when they reconnect

## Endpoint Strategy

### Keep Existing REST Endpoints
- `/api/v1/auth/login` - User authentication (returns 24-hour session token)
- `/api/v1/auth/register` - User registration  
- `/api/ws` - Single consolidated WebSocket endpoint (auth via session token)

### REST API Usage
- `/api/v1/users/profile` - Profile management (REST only)
- `/api/v1/files/upload` - File upload with CSAM/virus scanning
- `/api/v1/files/download/{hash}/{filename}` - File download by hash+name
- `/api/v1/health` - Health checks for monitoring

## File Management System

### File Upload & Storage
- **Hash-based deduplication**: Single copy per unique file hash
- **Filename preservation**: Store hash + original filename mapping
- **Security scanning**: CSAM detection and virus scanning on upload
- **Client delivery**: Send both hash and expected filename to clients

### File Structure
```
/storage/
  /files/
    /{first2chars_of_hash}/
      /{full_hash} - actual file content
  /metadata/
    /{hash}.json - filename mappings, upload info, scan results
```

## Authentication & Connection Flow

### Login Process
1. User calls `/api/v1/auth/login` with credentials
2. Server returns JWT session token (24-hour validity)
3. Client connects to `/api/ws` with session token (Basic Auth header or query param)
4. Server validates token and establishes WebSocket connection
5. Server sends queued commands/messages for offline period

### Connection Management
- **Heartbeat**: 3 missed heartbeats = connection termination
- **Error handling**: Connection errors logged to file, connection terminated
- **Reconnection**: Client should auto-reconnect with valid session token

## Command Queue System

### Persistent Command Storage
- Commands stored in database with recipient, timestamp, delivery status
- **Queue retention**: 2 weeks for moderation and delivery
- **Offline delivery**: When user reconnects, deliver all undelivered commands
- **File retention**: Command-related files kept for 2 weeks

### Command Status Tracking
- `pending` - Created but not delivered
- `delivered` - Sent to client via WebSocket
- `completed` - Marked complete by recipient
- `expired` - Older than 2 weeks, archived

## Content Category & Filtering System

### Category-Based Command Filtering
- Commands tagged with content categories (e.g., "feet", "censored", "general", "adult", etc.)
- Users can set preferences to block/allow specific categories
- Server filters commands before delivery based on user preferences
- Categories defined by server but user-controlled filtering

### User Preference Management
```json
{
  "type": "preferences.update",
  "payload": {
    "blockedCategories": ["feet", "extreme"],
    "allowedCategories": ["general", "censored", "mild"]
  }
}
```

## Broadcast System

### Broadcast Types
1. **Direct user**: Target specific user ID
2. **Categories**: Send to users who allow specific content categories
3. **All-cast**: Broadcast to all users (requires group membership verification)

### Category Management
- Categories defined and managed by server
- Users control their own filtering preferences
- Commands must be tagged with appropriate categories
- Server enforces category filtering during delivery

## WebSocket Message Types

### Command Operations
```json
{
  "type": "command.create",
  "payload": {
    "commandType": "message|task|system",
    "content": "Command text",
    "data": "Additional metadata",
    "targetUser": "user_id_or_username",
    "categories": ["general", "censored"], // content categories
    "broadcast": "user|category|allcast",
    "files": [{"hash": "file_hash", "filename": "original_name.jpg"}]
  }
}

{
  "type": "command.queue",
  "payload": {
    "since": "2025-07-15T10:00:00Z" // get commands since timestamp
  }
}

{
  "type": "command.complete",
  "payload": {
    "commandId": "uuid"
  }
}
```

### Chat Operations (Optional Feature)
```json
{
  "type": "chat.message",
  "payload": {
    "targetUser": "user_id",
    "message": "Text message",
    "messageType": "private|broadcast",
    "files": [{"hash": "file_hash", "filename": "image.png"}]
  }
}
```

### System Operations
```json
{
  "type": "heartbeat",
  "payload": {
    "timestamp": "2025-07-17T10:00:00Z"
  }
}

{
  "type": "user.status",
  "payload": {
    "status": "online|away|busy|offline"
  }
}
```

## Data Storage Requirements

### Database Tables
- **commands**: id, sender_id, content, command_type, categories[], status, created_at, expires_at
- **command_assignments**: command_id, recipient_id, delivered_at, completed_at
- **command_files**: command_id, file_hash, original_filename
- **file_metadata**: hash, size, mime_type, scan_status, upload_date
- **user_preferences**: user_id, blocked_categories[], allowed_categories[]
- **user_groups**: user_id, group_name, assigned_at (for all-cast permissions)
- **chat_messages**: (if implemented) sender_id, recipient_id, message, created_at

### File Storage
- Hash-based storage with deduplication
- Metadata tracking for filenames and scan results
- 2-week retention policy for command-related files
