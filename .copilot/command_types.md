# Command Types Analysis

## Legacy Command System
Based on analysis of the legacy codebase, the original system supported various command types:

### Core Command Types
1. **Direct Commands**: Text-based instructions sent from Dom to Sub
2. **System Commands**: Built-in operations like "Outstanding", "Content", "Delete"
3. **Relationship Commands**: "Invite", "Accept", "Reject" for Dom/Sub connections
4. **Feedback Commands**: "Thumbs" for rating system

### Legacy Database Structure
- `CommandList`: Stores command content and send date
- `ControlAppCmd`: Links Sender -> Sub with Command assignment
- `CmdLog`: Audit trail of all command activity
- `Users`: Dom/Sub roles, verification, thumbs up ratings

### Key Legacy Operations
- **USP_SendAppCmd**: Send commands to specific users or random users
- **USP_GetAppContent**: Retrieve next command for a user
- **USP_CmdComplete**: Mark command as completed/remove from queue
- **USP_GetOutstanding**: Get count of pending commands

## Modern Go Implementation

### Current Models
```go
type Command struct {
    ID        uuid.UUID
    Type      string    // "message", "task", "pose", "rule", "system"
    Content   string    // Human-readable command text
    Data      string    // Additional metadata/parameters
    Status    string    // "pending", "completed", "cancelled"
}

type ControlAppCmd struct {
    ID         uuid.UUID
    SenderID   uuid.UUID  // Dom user
    SubID      uuid.UUID  // Sub user
    CommandID  uuid.UUID  // Reference to Command
    GroupRefID *uuid.UUID // For group commands
}
```

### Proposed Command Types

1. **message**: Simple text communication
   - Content: "Hello there, pet!"
   - Data: "" (no additional data needed)

2. **task**: Actionable items
   - Content: "Go make me a sandwich"
   - Data: "location:kitchen" or JSON metadata

3. **pose**: Physical positions/actions
   - Content: "Kneel and wait for 5 minutes"
   - Data: "duration:5min"

4. **rule**: Behavioral guidelines
   - Content: "No speaking unless spoken to"
   - Data: "type:silence"

5. **system**: System operations
   - Content: "Get outstanding commands"
   - Data: "operation:outstanding"

### Implementation Notes
- Commands are created once, can be assigned to multiple users
- ControlAppCmd tracks individual assignments
- Status tracking at command level
- Group commands via GroupRefID for broadcasting
- Audit trail maintained through timestamps and relationships

## Migration Strategy
1. Keep similar command patterns but modernize data structure
2. Maintain role-based authorization (Dom can send, Sub receives)
3. Support both individual and group command distribution
4. Add proper validation and error handling
5. Include comprehensive logging and audit trails
