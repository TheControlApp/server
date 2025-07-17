# Architecture Changes - Removing Dom/Sub System

## Summary of Changes
Removing the Dom/Sub role differentiation and relationship management system to simplify the application architecture.

## Impact Analysis

### Database Models
**User Model Changes:**
- Remove `Role` field (no more "dom"/"sub" roles)
- Keep basic user fields: LoginName, ScreenName, Email, Password
- Maintain verification and thumbs up system
- Remove relationship/invite related fields

**Command Model Changes:**
- `ControlAppCmd` becomes user-agnostic assignment
- Remove `GroupRefID` (no more group relationships)
- Keep basic Sender -> Recipient command flow
- All users can send/receive commands

### API Changes
**Removed Endpoints:**
- Invitation system endpoints
- Relationship management endpoints
- Role-based access controls

**Simplified Endpoints:**
- `/api/v1/commands/create` - Any user can create commands
- `/api/v1/commands/pending` - Get commands assigned to user
- `/api/v1/commands/send` - Send command to specific user
- `/api/v1/users` - Basic user management

### WebSocket Changes
**Simplified Channels:**
- Remove role-based channel separation
- Single user channel per connected user
- Broadcast commands to specific recipients
- General chat/notification channels

### Command Types (Simplified)
1. **message**: Basic text communication between users
2. **task**: Actionable items any user can assign
3. **system**: System operations and notifications
4. **feedback**: Rating and thumbs up system

## Benefits
- Simplified user management
- Reduced complexity in authorization logic
- More flexible user interactions
- Easier to understand and maintain
- Focus on core command/task functionality

## Migration Notes
- Existing Dom/Sub users become regular users
- Command assignments remain but without role restrictions
- Relationships/invites will be removed from database
- WebSocket connections simplified to user-based channels
