# Implementation Priority Plan

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

### Phase 1: Core Infrastructure (Week 1-2)
1. **Single WebSocket endpoint** at `/api/ws`
2. **Session token authentication** system
3. **Basic command queue** storage and delivery
4. **Heartbeat mechanism** implementation

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

## Ready to Start Implementation

All major architectural questions have been answered. The system design is:
- **WebSocket-first** with minimal REST API
- **Secure file handling** with scanning and deduplication  
- **Persistent command queues** for offline users
- **Flexible broadcast system** with content categories and user filtering preferences
- **Robust connection management** with heartbeats

Ready to begin Phase 1 implementation?
