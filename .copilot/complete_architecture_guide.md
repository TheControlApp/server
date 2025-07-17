# Complete Architecture & Implementation Guide

## Project Overview
ControlMe Go Backend - A modern, secure, and scalable rewrite of the ControlMe platform in Go, providing a WebSocket-first API for real-time command and communication systems.

## Core Architecture Decisions

### API Design Philosophy
- **WebSocket-First**: Primary communication through single `/api/ws` endpoint
- **Minimal REST**: Only for authentication, profile management, and file operations
- **Real-time Priority**: Commands, notifications, and status updates via WebSocket
- **Persistent Queues**: Offline command delivery when users reconnect

### Authentication & Security
- **Session Tokens**: 24-hour JWT tokens from `/api/v1/auth/login`
- **WebSocket Auth**: Session token via Basic Auth header or query parameter
- **File Security**: CSAM scanning and virus detection on all uploads
- **Connection Security**: Heartbeat system with 3-miss termination

### File Management System
- **Hash-based Storage**: Single copy per unique file hash with deduplication
- **Filename Preservation**: Store hash + original filename mapping for client delivery
- **Storage Structure**: `/storage/files/{hash_prefix}/{full_hash}` + metadata JSON
- **Retention Policy**: 2-week retention for command-related files
- **Security Scanning**: CSAM detection and virus scanning on upload

### Content Category System
- **User-Controlled Filtering**: Users set preferences to block/allow content categories
- **Category Examples**: general, censored, adult, feet, extreme, roleplay, humiliation, public, private
- **Command Tagging**: All commands tagged with appropriate categories
- **Delivery Filtering**: Server filters commands based on recipient preferences
- **Default Behavior**: Commands default to "general" category if none specified

### Command System
- **Persistent Storage**: Commands stored for 2 weeks with status tracking
- **Queue Delivery**: Offline users receive backlog on reconnection
- **Status Tracking**: pending → delivered → completed → expired
- **File Attachments**: Commands can include file references (hash + filename)
- **Broadcast Types**: direct user, category-based, all-cast (with group membership)

## Database Schema

### Core Tables
```sql
-- Users
CREATE TABLE users (
    id UUID PRIMARY KEY,
    screen_name VARCHAR(50) NOT NULL,
    login_name VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(300) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    verified BOOLEAN DEFAULT false,
    verified_code INTEGER DEFAULT 0,
    thumbs_up INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    login_date TIMESTAMP DEFAULT NOW()
);

-- Commands
CREATE TABLE commands (
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    command_type VARCHAR(50) NOT NULL,
    categories TEXT[] DEFAULT '{"general"}',
    data TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '2 weeks'
);

-- Command Assignments
CREATE TABLE command_assignments (
    id UUID PRIMARY KEY,
    command_id UUID NOT NULL REFERENCES commands(id),
    recipient_id UUID NOT NULL REFERENCES users(id),
    delivered_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

-- User Preferences
CREATE TABLE user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id),
    blocked_categories TEXT[] DEFAULT '{}',
    allowed_categories TEXT[] DEFAULT '{"general"}',
    updated_at TIMESTAMP DEFAULT NOW()
);

-- File Metadata
CREATE TABLE file_metadata (
    hash VARCHAR(64) PRIMARY KEY,
    size BIGINT NOT NULL,
    mime_type VARCHAR(100),
    scan_status VARCHAR(20) DEFAULT 'pending',
    scan_results JSONB,
    upload_date TIMESTAMP DEFAULT NOW()
);

-- Command Files
CREATE TABLE command_files (
    command_id UUID REFERENCES commands(id),
    file_hash VARCHAR(64) REFERENCES file_metadata(hash),
    original_filename VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (command_id, file_hash)
);

-- User Groups (for all-cast permissions)
CREATE TABLE user_groups (
    user_id UUID REFERENCES users(id),
    group_name VARCHAR(50) NOT NULL,
    assigned_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (user_id, group_name)
);

-- Chat Messages (optional)
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL REFERENCES users(id),
    recipient_id UUID NOT NULL REFERENCES users(id),
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## API Endpoints

### REST API (Minimal)
```
POST /api/v1/auth/login
POST /api/v1/auth/register
GET/PUT /api/v1/users/profile
POST /api/v1/files/upload
GET /api/v1/files/download/{hash}/{filename}
GET /api/v1/health
```

### WebSocket Messages

#### Authentication & Connection
```json
// Connection with session token via Basic Auth or query param
// WebSocket URL: /api/ws?token=<session_token>

{
  "type": "heartbeat",
  "payload": {
    "timestamp": "2025-07-17T10:00:00Z"
  }
}
```

#### Command Operations
```json
{
  "type": "command.create",
  "payload": {
    "commandType": "message|task|system",
    "content": "Command text",
    "data": "Additional metadata",
    "targetUser": "user_id_or_username",
    "categories": ["general", "roleplay"],
    "broadcast": "user|category|allcast",
    "files": [{"hash": "file_hash", "filename": "original_name.jpg"}]
  }
}

{
  "type": "command.queue",
  "payload": {
    "since": "2025-07-15T10:00:00Z"
  }
}

{
  "type": "command.complete",
  "payload": {
    "commandId": "uuid"
  }
}
```

#### User Preferences
```json
{
  "type": "preferences.set",
  "payload": {
    "blockedCategories": ["feet", "extreme"],
    "allowedCategories": ["general", "censored", "roleplay"]
  }
}
```

#### Chat Operations (Optional)
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

#### System Operations
```json
{
  "type": "user.status",
  "payload": {
    "status": "online|away|busy|offline"
  }
}

{
  "type": "system.notification",
  "payload": {
    "notificationType": "info|warning|error",
    "message": "Notification text"
  }
}
```

## Implementation Phases

### Phase 1: Core Infrastructure (Week 1-2)
- Single WebSocket endpoint consolidation at `/api/ws`
- Session token authentication system
- Basic command queue storage and delivery
- Heartbeat mechanism implementation
- Connection management with error logging

### Phase 2: File System (Week 3-4)
- File upload API with hash-based storage
- CSAM scanning integration (external service)
- Virus scanning setup (ClamAV or similar)
- Deduplication system implementation
- File metadata tracking and retention policies

### Phase 3: Command & Broadcast System (Week 5-6)
- Command creation/assignment via WebSocket
- Content category system implementation
- User preference management (block/allow categories)
- Category-based filtering for command delivery
- Broadcast functionality (user/category/all-cast)
- Command completion tracking
- Queue delivery for reconnecting users

### Phase 4: Polish & Optional Features (Week 7)
- Chat system (if decision made to implement)
- Advanced moderation tools
- Admin dashboard for user/category management
- Performance optimization
- Comprehensive logging and monitoring

## Key Implementation Notes

### File Storage Strategy
- Store files as `/storage/files/{first2chars_of_hash}/{full_hash}`
- Metadata stored separately as `{hash}.json` with filename mappings
- Deduplication at storage level, filename preservation at client level
- 2-week cleanup process for expired command files

### Category Filtering Logic
- Commands blocked if ANY category is in user's blocked list
- Commands delivered only if ALL categories are in user's allowed list
- Default to "general" category for untagged commands
- Admin override capabilities for moderation purposes

### Connection Management
- WebSocket heartbeat every 30 seconds
- 3 missed heartbeats = automatic disconnection
- Graceful reconnection with queue delivery
- Error logging to files with connection termination

### Security Considerations
- All file uploads scanned for CSAM and viruses
- Session tokens expire after 24 hours
- Command retention limited to 2 weeks
- User preferences control content exposure
- Admin oversight on broadcast permissions

### Offline User Handling
- Commands queued in database for offline users
- Queue delivered on reconnection based on last_seen timestamp
- Commands expire after 2 weeks automatically
- File references maintained during queue period

## Next Steps for Implementation
1. Begin with Phase 1 WebSocket consolidation
2. Implement session token authentication
3. Build basic command queue system
4. Add heartbeat and connection management
5. Progress through phases with testing at each stage

## Dependencies & External Services
- PostgreSQL database
- CSAM scanning service (PhotoDNA, AWS Rekognition, or similar)
- Virus scanning (ClamAV or cloud-based solution)
- File storage system (local filesystem or cloud storage)
- JWT library for session tokens
- WebSocket library (Gorilla WebSocket or similar)

This document contains all architectural decisions, database schemas, API specifications, and implementation guidance needed to continue development of the ControlMe Go Backend system.
