# Database Schema

## Overview
PostgreSQL 15+ with GORM ORM. UUID primary keys and JSONB for flexible data.

## Core Tables

### Users
- ID (UUID primary key)
- Username (unique, 50 chars)
- Display name (100 chars)  
- Email (unique, 255 chars)
- Password hash
- Role (default 'user')
- Preferences (JSONB)
- Timestamps

### Commands
- ID (UUID primary key)
- Sender/receiver user IDs
- Instructions (JSONB array)
- Status (pending/completed/failed)
- Tags (text array)
- Timestamps

### Files
- Hash (64 char primary key)
- Filename and content type
- Size in bytes
- Uploader user ID
- Upload timestamp

### User Relationships
- ID (UUID primary key)
- User and related user IDs
- Relationship type (blocked/friend)
- Creation timestamp

## Key Features
- **UUID Keys** - All primary keys are UUIDs
- **JSONB** - Flexible data storage for preferences and instructions
- **Timestamps** - Automatic creation/update tracking
- **Relationships** - User blocking and friendship system
- **File Deduplication** - Hash-based file storage
