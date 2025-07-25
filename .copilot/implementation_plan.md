# Implementation Priority Plan

## Current Status: PHASE 1 READY ✅

### Infrastructure Completed
- ✅ **Database Schema**: All 4 core tables migrated successfully
  - `users` - User accounts with UUID primary keys
  - `commands` - Command definitions with type, content, status tracking
  - `control_app_cmds` - Command assignments with sender/receiver relationships
  - `chat_logs` - Chat messages between users
- ✅ **Server Running**: Port 8080 with all endpoints functional
- ✅ **REST API**: Authentication, user management, command management endpoints
- ✅ **WebSocket Infrastructure**: Basic endpoints `/ws/client` and `/ws/web` available
- ✅ **Database Migrations**: GORM AutoMigrate working with PostgreSQL 15
- ✅ **Docker Environment**: PostgreSQL + Go server containerized and stable

### Models Currently Implemented
```go
// Core working models (migrated successfully)
- User: Full user management with authentication fields
- Command: Command storage with type, content, data, status
- ControlAppCmd: Command assignments linking senders, receivers, commands
- ChatLog: Message storage between users

// Legacy models (present but not in migration)
- Group, GroupMember, Relationship, Block, Report, Invite
```

## Clarifications Received ✅

### File Management
- **Any file type uploads** with CSAM + virus scanning
- **Hash-based deduplication** (single copy per unique hash)
- **Filename preservation** (hash + original name sent to clients)
- **Storage structure**: `/storage/files/{hash_prefix}/{full_hash}` + metadata

### Authentication & Connection
- **24-hour session tokens** from `/api/v1/auth/login`
- **WebSocket auth** via Basic Auth header or query param
- **Connection queue delivery** on reconnect for offline period

### Data Persistence
- **Profile management**: REST API only
- **Chat messages**: 50/50 decision, store on both client/server if implemented
- **Command history**: 2-week retention on server + client
- **File retention**: 2 weeks for command-related files

### Broadcast & Groups
- **Content categories** for command filtering (e.g., "feet", "censored", "general", etc.)
- **User preferences**: Users can block/allow specific categories
- **Broadcast types**: direct user, categories, all-cast
- **Category filtering**: Commands tagged with categories, users control what they receive

### Connection Management
- **Heartbeat system**: 3 missed beats = disconnect
- **Error handling**: Log to file, terminate connection
- **Command queues**: Persistent storage for offline users

## Implementation Phases

### Phase 1: Core Infrastructure (CURRENT FOCUS) 🎯
**Status: Foundation Complete - Ready for WebSocket Implementation**

1. ✅ **Database foundation** - All core tables migrated
2. ✅ **Basic server** - Running with REST endpoints
3. 🔄 **WebSocket endpoint enhancement** - Implement actual WebSocket logic at `/api/ws`
4. 🔄 **Session token authentication** - Integrate JWT auth with WebSocket connections
5. 🔄 **Basic command queue** - Implement in-memory then persistent queues
6. 🔄 **Heartbeat mechanism** - Connection management and cleanup

### Phase 2: File System (Week 3-4)
1. **File upload API** with hash-based storage
2. **CSAM scanning integration** (external service)
3. **Virus scanning setup** (ClamAV or similar)
4. **Deduplication system** implementation
5. **File metadata tracking**

### Phase 3: Command & Broadcast System (Week 5-6)
1. **Command creation/assignment** via WebSocket
2. **Content category system** (feet, censored, general, etc.)
3. **User preference management** (block/allow categories)
4. **Category-based filtering** for command delivery
5. **Broadcast functionality** (user/category/all-cast)
6. **Command completion tracking**
7. **Queue delivery for reconnecting users**

### Phase 4: Polish & Optional Features (Week 7)
1. **Chat system** (if decision made to implement)
2. **Advanced moderation tools**
3. **Admin dashboard** for group/user management
4. **Performance optimization**
5. **Comprehensive logging**

## Ready for Phase 1 Implementation ✅

### Immediate Next Steps:
1. **WebSocket Message Protocol** - Define JSON message structure for commands
2. **Connection Manager** - Track active WebSocket connections with user mapping
3. **Command Queue Service** - Persistent storage for offline message delivery
4. **Authentication Integration** - Session token validation for WebSocket connections

### Technical Foundation Solid:
- Database schema matches architecture requirements
- Server infrastructure stable and tested
- REST API provides necessary user/auth management
- Legacy models isolated (not breaking current functionality)
- Docker environment proven working

**Ready to begin Phase 1 WebSocket implementation immediately!** 🚀
