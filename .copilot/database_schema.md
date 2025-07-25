# Database Schema Documentation

## Overview
The ControlMe database uses PostgreSQL with UUID primary keys and GORM ORM for Go integration. The schema supports a flexible command system with multi-instruction commands, tag-based content filtering, and user management.

## Current Active Models

### 1. User Model
Represents user accounts and authentication information.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    screen_name VARCHAR(50) NOT NULL,
    login_name VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(300) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50),
    random_opt_in BOOLEAN DEFAULT FALSE,
    anon_cmd BOOLEAN DEFAULT FALSE,
    verified BOOLEAN DEFAULT FALSE,
    verified_code INTEGER DEFAULT 0,
    thumbs_up INTEGER DEFAULT 0,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    login_date TIMESTAMP
);

CREATE INDEX idx_users_login_name ON users(login_name);
CREATE INDEX idx_users_email ON users(email);
```

**Fields:**
- `id`: UUID primary key
- `screen_name`: Display name for the user
- `login_name`: Unique login identifier
- `email`: User's email address (unique)
- `password`: Hashed password (never returned in JSON)
- `role`: User role (e.g., "admin", "user")
- `random_opt_in`: Whether user opts into random command broadcasts
- `anon_cmd`: Whether user allows anonymous commands
- `verified`: Account verification status
- `verified_code`: Email verification code
- `thumbs_up`: User rating/karma score
- `login_date`: Last login timestamp (set by BeforeCreate hook)

### 2. Command Model
Represents command assignments with multiple instructions.

```sql
CREATE TABLE commands (
    id UUID PRIMARY KEY,
    instructions TEXT NOT NULL, -- JSON array of instruction objects
    sender_id UUID NOT NULL REFERENCES users(id),
    receiver_id UUID REFERENCES users(id), -- NULL for broadcasts
    tags TEXT, -- JSON array of tag names
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_commands_sender_id ON commands(sender_id);
CREATE INDEX idx_commands_receiver_id ON commands(receiver_id);
CREATE INDEX idx_commands_status ON commands(status);
CREATE INDEX idx_commands_created_at ON commands(created_at);
```

**Fields:**
- `id`: UUID primary key
- `instructions`: JSON string containing array of instruction objects
- `sender_id`: Foreign key to user who sent the command
- `receiver_id`: Foreign key to specific target user (NULL for broadcasts)
- `tags`: JSON array of tag names for content filtering
- `status`: Command lifecycle status ("pending", "delivered", "completed")

**Example Instructions JSON:**
```json
[
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
      "file_hash": "abc123...",
      "file_name": "task.pdf"
    }
  }
]
```

**Example Tags JSON:**
```json
["chastity", "daily", "general"]
```

### 3. Tag Model
Represents content categories for filtering.

```sql
CREATE TABLE tags (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_tags_name ON tags(name);
```

**Fields:**
- `id`: UUID primary key
- `name`: Unique tag name (e.g., "chastity", "feet", "general")
- `description`: Human-readable description of the tag

**Standard Tags:**
- `general`: Safe-for-work content
- `adult`: Explicit content
- `chastity`: Chastity-related content
- `feet`: Foot-related content
- `roleplay`: Roleplay scenarios
- `humiliation`: Humiliation-based content
- `public`: Public setting commands
- `private`: Private/intimate commands

### 4. Block Model
Represents user blocking relationships.

```sql
CREATE TABLE blocks (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    blocked_id UUID NOT NULL REFERENCES users(id),
    reason TEXT,
    created_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_blocks_unique ON blocks(user_id, blocked_id);
CREATE INDEX idx_blocks_user_id ON blocks(user_id);
CREATE INDEX idx_blocks_blocked_id ON blocks(blocked_id);
```

**Fields:**
- `id`: UUID primary key
- `user_id`: User who is doing the blocking
- `blocked_id`: User being blocked
- `reason`: Optional reason for blocking

### 5. Report Model
Represents user reports for moderation.

```sql
CREATE TABLE reports (
    id UUID PRIMARY KEY,
    reporter_id UUID NOT NULL REFERENCES users(id),
    reported_id UUID NOT NULL REFERENCES users(id),
    reason TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP
);

CREATE INDEX idx_reports_reporter_id ON reports(reporter_id);
CREATE INDEX idx_reports_reported_id ON reports(reported_id);
CREATE INDEX idx_reports_status ON reports(status);
CREATE INDEX idx_reports_created_at ON reports(created_at);
```

**Fields:**
- `id`: UUID primary key
- `reporter_id`: User making the report
- `reported_id`: User being reported
- `reason`: Description of the issue
- `status`: Report status ("pending", "reviewed", "resolved", "dismissed")

## Planned Future Models

### File Metadata Model (Phase 2)
```sql
CREATE TABLE file_metadata (
    id UUID PRIMARY KEY,
    hash VARCHAR(128) NOT NULL UNIQUE,
    scan_status VARCHAR(20) DEFAULT 'pending',
    scan_results TEXT, -- JSON
    uploader_id UUID REFERENCES users(id),
    created_at TIMESTAMP,
    expires_at TIMESTAMP
);

CREATE INDEX idx_file_metadata_hash ON file_metadata(hash);
CREATE INDEX idx_file_metadata_uploader_id ON file_metadata(uploader_id);
CREATE INDEX idx_file_metadata_expires_at ON file_metadata(expires_at);
```

### User Preferences Model (Phase 3)
```sql
CREATE TABLE user_preferences (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    allowed_tags TEXT, -- JSON array
    blocked_tags TEXT, -- JSON array
    allow_anonymous BOOLEAN DEFAULT TRUE,
    allow_broadcasts BOOLEAN DEFAULT TRUE,
    settings TEXT, -- JSON object for additional settings
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE UNIQUE INDEX idx_user_preferences_user_id ON user_preferences(user_id);
```

### Session Management Model (Phase 1)
```sql
CREATE TABLE user_sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    token_hash VARCHAR(128) NOT NULL,
    expires_at TIMESTAMP,
    created_at TIMESTAMP,
    last_activity TIMESTAMP
);

CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_token_hash ON user_sessions(token_hash);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);
```

## GORM Model Definitions

The Go models use GORM tags for database mapping:

```go
// User model with GORM tags
type User struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
    ScreenName   string    `gorm:"size:50;not null" json:"screen_name"`
    LoginName    string    `gorm:"size:50;not null;unique" json:"login_name"`
    Email        string    `gorm:"size:300;not null;unique" json:"email"`
    Password     string    `gorm:"size:255;not null" json:"-"`
    // ... other fields
}
```

## Migration Strategy

### Current Migration (Completed)
- Users table with authentication fields
- Commands table with JSON instruction storage
- Basic indexes for performance

### Phase 1 Migration (Next)
- Add user_sessions table
- Add file_metadata table
- Update commands table with additional indexes

### Phase 2 Migration
- Add user_preferences table
- Add command_queue table for offline delivery
- Add performance optimization indexes

## Performance Considerations

### Indexes
- Primary keys on all UUID fields
- Foreign key indexes for joins
- Status and timestamp indexes for filtering
- Unique constraints where appropriate

### JSON Fields
- PostgreSQL JSON/JSONB support for flexible instruction storage
- Consider JSONB for better query performance on tag filtering
- Use GIN indexes on JSONB fields when needed

### Cleanup Jobs
- Regular cleanup of expired commands (2-week retention)
- File cleanup based on command references
- Session cleanup for expired tokens

## Query Examples

### Find Commands for User
```sql
SELECT c.*, u.screen_name as sender_name 
FROM commands c 
JOIN users u ON c.sender_id = u.id 
WHERE c.receiver_id = $1 
   OR (c.receiver_id IS NULL AND c.tags && $2)
ORDER BY c.created_at DESC;
```

### Check if User is Blocked
```sql
SELECT EXISTS(
    SELECT 1 FROM blocks 
    WHERE user_id = $1 AND blocked_id = $2
);
```

### Get User's Tag Preferences
```sql
SELECT allowed_tags, blocked_tags 
FROM user_preferences 
WHERE user_id = $1;
```

## Security Considerations

### Password Storage
- Passwords are hashed using bcrypt
- Never returned in JSON responses (`json:"-"` tag)

### UUID Usage
- All primary keys use UUIDs to prevent enumeration attacks
- Foreign keys properly constrained with CASCADE options

### Data Retention
- Commands auto-expire after 2 weeks
- File cleanup based on command references
- User data retention follows GDPR requirements

### Access Control
- All database access through authenticated sessions
- Row-level security could be implemented for multi-tenant scenarios
- Audit logging for sensitive operations
