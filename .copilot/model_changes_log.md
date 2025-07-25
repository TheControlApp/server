# Model Changes and Updates Log

## Overview
This document tracks all changes made to the database models and their implications for the API and client implementations.

## Recent Model Changes (July 24, 2025)

### 1. Client-Facing API Simplification

**Username-Based References:**
- Commands now use `username` strings instead of UUID references for sender/receiver
- Clients send/receive `"sender": "master123"` instead of `"sender_id": "uuid"`
- Server handles UUID lookups internally via username→UUID mapping
- Simplified client integration - no UUID management required

**Simplified File Upload:**
```javascript
// Before: Complex form data
const formData = new FormData();
formData.append('file', file);
formData.append('original_name', file.name);
formData.append('content_type', file.type);
formData.append('size', file.size);

// After: Just the file
const formData = new FormData();
formData.append('file', file);
```

**Simplified File Storage:**
- Database only stores file hash and metadata
- No original filename or content type storage
- Client specifies download filename in instruction content
- Server provides raw file bytes, client handles presentation

### 1. Model Cleanup and Removal

**Removed Models:**
- ❌ **ChatLog** - Chat functionality removed for current implementation
- ❌ **GroupMember** - Client-side tag preferences replace group membership
- ❌ **Relationship** - User relationships not needed with tag-based system
- ❌ **Invite** - Invitation system removed from scope

**Rationale:** These models were part of a more complex social system that conflicts with the simplified tag-based content filtering approach.

### 2. Command Model Consolidation

**Before (Two separate models):**
```go
// Old Command - just templates
type Command struct {
    ID      uuid.UUID
    Type    string     // single type
    Content string     // single content
    Data    string     // redundant field
    Status  string
}

// Old ControlAppCmd - assignments
type ControlAppCmd struct {
    ID        uuid.UUID
    SenderID  uuid.UUID
    SubID     uuid.UUID  // receiver
    CommandID uuid.UUID  // reference to Command
    GroupRefID *uuid.UUID
}
```

**After (Single consolidated model):**
```go
// New Command - complete command assignments
type Command struct {
    ID           uuid.UUID
    Instructions string     // JSON array of multiple instructions
    SenderID     uuid.UUID  // direct sender reference
    ReceiverID   *uuid.UUID // optional receiver (null = broadcast)
    Tags         string     // JSON array of tag names
    Status       string
}
```

**Benefits:**
- ✅ Simpler database queries (single table)
- ✅ Better JSON serialization for API responses
- ✅ Support for multi-step commands
- ✅ Direct tag-based broadcasting
- ✅ Eliminated redundant Data field

### 3. Group to Tag Model Rename

**Before:**
```go
type Group struct {
    ID          uuid.UUID
    Name        string
    Description string
    OwnerID     uuid.UUID  // Groups had owners
}
```

**After:**
```go
type Tag struct {
    ID          uuid.UUID
    Name        string    // "chastity", "feet", "general", etc.
    Description string
    // No owner - tags are system-wide categories
}
```

**Changes:**
- ✅ Removed OwnerID - tags are global categories
- ✅ Added unique constraint on Name
- ✅ Tags represent content categories, not user groups

### 4. Instructions Structure Evolution

**Client JSON Structure:**
```json
{
  "id": "uuid",
  "instructions": [
    {
      "type": "popup-msg",
      "content": {
        "body": "Complete task",
        "button": "Done"
      }
    },
    {
      "type": "download-file", 
      "content": {
        "file_hash": "abc123...",
        "file_name": "task.pdf" // Client defines download name
      }
    }
  ],
  "sender": "master123", // username string
  "receiver": "user456", // username string (null for broadcasts)
  "tags": ["chastity", "daily"],
  "status": "pending"
}
```

**Key Features:**
- 🎯 **Sequential Execution**: Instructions execute in array order
- 🎯 **Arbitrary Types**: Type field is completely flexible for custom clients
- 🎯 **Structured Content**: Each instruction has type-specific content
- 🎯 **File Integration**: File references use hash-based download API

## Database Migration Impact

### Current Migration Status
```go
// Currently migrated models
db.AutoMigrate(
    &User{},
    &Command{},    // Updated structure
    &Tag{},        // Renamed from Group
    &Block{},
    &Report{},
)
```

### Migration Steps Required
1. **Drop Removed Tables**: ChatLog, GroupMember, Relationship, Invite
2. **Update Command Table**: 
   - Add Instructions field (TEXT)
   - Add Tags field (TEXT) 
   - Remove Type, Content, Data fields
   - Update ReceiverID to nullable
3. **Rename Group to Tag**:
   - Rename table
   - Remove OwnerID field
   - Add unique constraint on Name

### Database Schema Changes
```sql
-- Drop removed tables
DROP TABLE IF EXISTS chat_logs;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS relationships;
DROP TABLE IF EXISTS invites;

-- Update commands table
ALTER TABLE commands 
    ADD COLUMN instructions TEXT NOT NULL,
    ADD COLUMN tags TEXT,
    ALTER COLUMN receiver_id DROP NOT NULL,
    DROP COLUMN type,
    DROP COLUMN content,
    DROP COLUMN data;

-- Rename and update groups table
ALTER TABLE groups RENAME TO tags;
ALTER TABLE tags DROP COLUMN owner_id;
ALTER TABLE tags ADD CONSTRAINT tags_name_unique UNIQUE (name);
```

## API Changes

### WebSocket Message Format Changes

**Before:**
```json
{
  "type": "command",
  "command_id": "uuid",
  "command_type": "popup",
  "content": "message text",
  "sender": "user123"
}
```

**After:**
```json
{
  "type": "command", 
  "data": {
    "id": "uuid",
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "message text",
          "button": "OK"
        }
      }
    ],
    "sender": {
      "id": "uuid",
      "screen_name": "user123"
    },
    "tags": ["general"]
  }
}
```

### File API Integration

**New File Reference Format:**
```json
{
  "type": "download-file",
  "content": {
    "file_hash": "abc123def456...",
    "file_name": "task-instructions.pdf",
    "file_size": 2048576,
    "content_type": "application/pdf"
  }
}
```

**Download URL Format:**
```
GET /api/v1/files?filehash=abc123def456&filename=task-instructions.pdf
```

## Client Integration Impact

### JavaScript Client Changes

**Before:**
```javascript
// Single instruction handling
function handleCommand(command) {
  switch(command.type) {
    case 'popup':
      showPopup(command.content);
      break;
  }
}
```

**After:**
```javascript
// Multi-instruction handling
function handleCommand(command) {
  command.instructions.forEach((instruction, index) => {
    setTimeout(() => {
      executeInstruction(instruction, command);
    }, index * 1000);
  });
}

function executeInstruction(instruction, command) {
  switch(instruction.type) {
    case 'popup-msg':
      showPopup(instruction.content);
      break;
    case 'download-file':
      downloadFile(instruction.content);
      break;
    default:
      console.warn(`Unknown type: ${instruction.type}`);
  }
}
```

### Mobile Client Changes

**File Download Handling:**
```javascript
// React Native file download
async function handleFileDownload(content) {
  const { file_hash, file_name } = content;
  const downloadUrl = `/api/v1/files?filehash=${file_hash}&filename=${file_name}`;
  
  const response = await RNFetchBlob.fetch('GET', downloadUrl, {
    'Authorization': `Bearer ${authToken}`
  });
  
  // Save to device
  const filePath = `${RNFetchBlob.fs.dirs.DownloadDir}/${file_name}`;
  await RNFetchBlob.fs.writeFile(filePath, response.data, 'base64');
}
```

## Breaking Changes Summary

### For Server Development
1. **Database Migration Required**: Models have significantly changed
2. **API Response Format**: Commands now return instruction arrays
3. **WebSocket Message Format**: Updated message structure
4. **File API**: New hash-based file download system

### For Client Development
1. **Instruction Processing**: Must handle arrays of instructions
2. **Type System**: Instruction types are now arbitrary/extensible
3. **File Downloads**: New file hash + filename URL format
4. **Tag Filtering**: Tags are now arrays, not single values

### Backward Compatibility
- ❌ **API Responses**: Complete change in command structure
- ❌ **WebSocket Messages**: New message format
- ❌ **Database Schema**: Major table restructuring
- ✅ **Authentication**: JWT token system unchanged
- ✅ **Basic REST Endpoints**: User management endpoints unchanged

## Testing Requirements

### Database Testing
- [ ] Test model migrations work correctly
- [ ] Verify foreign key constraints
- [ ] Test JSON field storage and retrieval
- [ ] Validate unique constraints on tags

### API Testing  
- [ ] Test multi-instruction command creation
- [ ] Verify file upload/download flow
- [ ] Test broadcast vs direct message delivery
- [ ] Validate tag-based filtering

### Client Testing
- [ ] Test instruction sequence execution
- [ ] Verify unknown instruction type handling
- [ ] Test file download integration
- [ ] Validate WebSocket reconnection

## Migration Checklist

### Phase 1 (Immediate)
- [x] Update Go models
- [x] Remove unused models
- [x] Consolidate Command + ControlAppCmd
- [x] Rename Group to Tag
- [ ] Update database migration
- [ ] Test model changes

### Phase 2 (API Updates)
- [ ] Update WebSocket handlers
- [ ] Implement multi-instruction processing
- [ ] Add file API endpoints
- [ ] Update REST API responses

### Phase 3 (Client Updates)
- [ ] Update client integration examples
- [ ] Test cross-platform compatibility
- [ ] Validate instruction type extensibility
- [ ] Document breaking changes

This consolidation significantly simplifies the system architecture while adding powerful multi-instruction and extensible type capabilities! 🚀
