# WebSocket Message Examples

This file contains real JSON examples of messages that flow between the ControlMe server and clients.

## Connection Examples

### Client Connects and Authenticates

**1. Client initiates WebSocket connection:**
```
ws://localhost:8080/ws/client?token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**2. Server confirms authentication:**
```json
{
  "type": "connection_status",
  "id": "conn-12345",
  "timestamp": "2025-09-27T14:30:15.123Z",
  "from": "00000000-0000-0000-0000-000000000000",
  "to": "user-550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "status": "authenticated",
    "user_id": "user-550e8400-e29b-41d4-a716-446655440001",
    "message": "WebSocket connection established"
  }
}
```

## Command Flow Examples

### Example 1: Single Instruction Command

**Server sends command to client:**
```json
{
  "type": "command_assignment",
  "id": "cmd-msg-12345",
  "timestamp": "2025-09-27T14:35:22.456Z",
  "from": "dom-550e8400-e29b-41d4-a716-446655440002",
  "to": "sub-550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "command": {
      "id": "550e8400-e29b-41d4-a716-446655440003",
      "instructions": [
        {
          "type": "std_popup",
          "content": {
            "title": "Daily Check-in",
            "body": "Please confirm you are ready to start your tasks for today.",
            "button": "Ready!"
          }
        }
      ],
      "sender_id": "dom-550e8400-e29b-41d4-a716-446655440002",
      "receiver_id": "sub-550e8400-e29b-41d4-a716-446655440001",
      "status": "pending",
      "created_at": "2025-09-27T14:35:22.456Z",
      "sender": {
        "id": "dom-550e8400-e29b-41d4-a716-446655440002",
        "screen_name": "MasterUser",
        "login_name": "master_dom"
      }
    }
  }
}
```

**Instructions Array Structure:**
```json
[
  {
    "type": "std_popup",
    "content": {
      "title": "Daily Check-in",
      "body": "Please confirm you are ready to start your tasks for today.",
      "button": "Ready!"
    }
  }
]
```

**Client acknowledges command completion:**
```json
{
  "type": "command_completion",
  "id": "completion-12345",
  "timestamp": "2025-09-27T14:36:45.789Z",
  "from": "sub-550e8400-e29b-41d4-a716-446655440001",
  "to": "dom-550e8400-e29b-41d4-a716-446655440002",
  "data": {
    "command_id": "550e8400-e29b-41d4-a716-446655440003",
    "status": "completed",
    "response": {
      "success": true,
      "timestamp": "2025-09-27T14:36:45.789Z",
      "data": {
        "button_clicked": "Ready!"
      }
    }
  }
}
```

### Example 2: Timer Command

**Server sends timer command:**
```json
{
  "type": "command_assignment",
  "id": "timer-cmd-12345",
  "timestamp": "2025-09-27T15:00:00.000Z",
  "from": "dom-550e8400-e29b-41d4-a716-446655440002",
  "to": "sub-550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "command": {
      "id": "timer-550e8400-e29b-41d4-a716-446655440004",
      "instructions": [
        {
          "type": "std_timer",
          "content": {
            "duration": 900,
            "title": "Break Time", 
            "description": "Take a 15-minute break. Relax and recharge.",
            "show_progress": true,
            "auto_complete": true
          }
        }
      ],
      "sender_id": "dom-550e8400-e29b-41d4-a716-446655440002",
      "receiver_id": "sub-550e8400-e29b-41d4-a716-446655440001",
      "status": "pending",
      "created_at": "2025-09-27T15:00:00.000Z",
      "sender": {
        "id": "dom-550e8400-e29b-41d4-a716-446655440002",
        "screen_name": "MasterUser",
        "login_name": "master_dom"
      }
    }
  }
}
```

**Client completes timer:**
```json
{
  "type": "command_completion",
  "id": "timer-completion-12345",
  "timestamp": "2025-09-27T15:15:00.000Z",
  "from": "sub-550e8400-e29b-41d4-a716-446655440001",
  "to": "dom-550e8400-e29b-41d4-a716-446655440002",
  "data": {
    "command_id": "timer-550e8400-e29b-41d4-a716-446655440004",
    "status": "completed",
    "response": {
      "success": true,
      "timestamp": "2025-09-27T15:15:00.000Z",
      "data": {
        "duration_completed": 900,
        "completed_method": "auto"
      }
    }
  }
}
```

### Example 4: Multi-Instruction Command

**Server sends complex multi-step command:**
```json
{
  "type": "command_assignment",
  "id": "multi-cmd-12345",
  "timestamp": "2025-09-27T16:30:00.000Z",
  "from": "dom-550e8400-e29b-41d4-a716-446655440002",
  "to": "sub-550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "command": {
      "id": "multi-550e8400-e29b-41d4-a716-446655440006",
      "instructions": [
        {
          "type": "std_notification",
          "content": {
            "title": "Task Sequence Starting",
            "body": "Beginning multi-step task sequence",
            "priority": "normal"
          }
        },
        {
          "type": "std_popup",
          "content": {
            "title": "Step 1: Preparation",
            "body": "Please ensure you have a quiet workspace and close all distracting applications.",
            "button": "Ready"
          }
        },
        {
          "type": "std_timer",
          "content": {
            "duration": 600,
            "title": "Focus Session",
            "description": "Work on your assigned project for the next 10 minutes",
            "show_progress": true,
            "auto_complete": false
          }
        },
        {
          "type": "std_input",
          "content": {
            "title": "Session Report",
            "description": "How did the focus session go?",
            "fields": [
              {
                "name": "productivity",
                "label": "Productivity Level (1-10)",
                "type": "number",
                "required": true
              },
              {
                "name": "distractions",
                "label": "Number of Distractions",
                "type": "number",
                "required": true
              },
              {
                "name": "notes",
                "label": "Additional Notes",
                "type": "text",
                "required": false
              }
            ],
            "submit_button": "Submit Report"
          }
        },
        {
          "type": "std_notification",
          "content": {
            "title": "Task Complete",
            "body": "Thank you for completing the focus session!",
            "priority": "high"
          }
        }
      ],
      "sender_id": "dom-550e8400-e29b-41d4-a716-446655440002",
      "receiver_id": "sub-550e8400-e29b-41d4-a716-446655440001",
      "status": "pending",
      "created_at": "2025-09-27T16:30:00.000Z",
      "sender": {
        "id": "dom-550e8400-e29b-41d4-a716-446655440002",
        "screen_name": "MasterUser",
        "login_name": "master_dom"
      }
    }
  }
}
```

**Instructions Array (parsed from JSON string):**
```json
[
  {
    "type": "std_notification",
    "content": {
      "title": "Task Sequence Starting",
      "body": "Beginning multi-step task sequence",
      "priority": "normal"
    }
  },
  {
    "type": "std_popup",
    "content": {
      "title": "Step 1: Preparation",
      "body": "Please ensure you have a quiet workspace and close all distracting applications.",
      "button": "Ready"
    }
  },
  {
    "type": "std_timer",
    "content": {
      "duration": 600,
      "title": "Focus Session",
      "description": "Work on your assigned project for the next 10 minutes",
      "show_progress": true,
      "auto_complete": false
    }
  },
  {
    "type": "std_input",
    "content": {
      "title": "Session Report",
      "description": "How did the focus session go?",
      "fields": [
        {
          "name": "productivity",
          "label": "Productivity Level (1-10)",
          "type": "number",
          "required": true
        },
        {
          "name": "distractions",
          "label": "Number of Distractions",
          "type": "number",
          "required": true
        },
        {
          "name": "notes",
          "label": "Additional Notes",
          "type": "text",
          "required": false
        }
      ],
      "submit_button": "Submit Report"
    }
  },
  {
    "type": "std_notification",
    "content": {
      "title": "Task Complete",
      "body": "Thank you for completing the focus session!",
      "priority": "high"
    }
  }
]
```

**Client completes multi-instruction command:**
```json
{
  "type": "command_completion",
  "id": "multi-completion-12345",
  "timestamp": "2025-09-27T16:45:30.000Z",
  "from": "sub-550e8400-e29b-41d4-a716-446655440001",
  "to": "dom-550e8400-e29b-41d4-a716-446655440002",
  "data": {
    "command_id": "multi-550e8400-e29b-41d4-a716-446655440006",
    "status": "completed",
    "response": {
      "success": true,
      "timestamp": "2025-09-27T16:45:30.000Z",
      "data": {
        "instructions_completed": 5,
        "execution_time_seconds": 930,
        "final_form_data": {
          "productivity": 8,
          "distractions": 2,
          "notes": "Good focus session, only interrupted by two phone calls."
        }
      }
    }
  }
}
```

### Example 3: Input Form Command

**Server requests user input:**
```json
{
  "type": "command_assignment",
  "id": "input-cmd-12345",
  "timestamp": "2025-09-27T16:00:00.000Z",
  "from": "dom-550e8400-e29b-41d4-a716-446655440002",
  "to": "sub-550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "command": {
      "id": "input-550e8400-e29b-41d4-a716-446655440005",
      "instructions": [
        {
          "type": "std_input",
          "content": {
            "title": "Daily Report",
            "description": "Please provide your daily progress report",
            "fields": [
              {
                "name": "tasks_completed",
                "label": "Tasks Completed Today",
                "type": "number",
                "required": true
              },
              {
                "name": "mood_rating",
                "label": "How are you feeling?",
                "type": "select",
                "required": true,
                "options": ["Excellent", "Good", "Fair", "Struggling"]
              },
              {
                "name": "notes",
                "label": "Additional Notes",
                "type": "text",
                "required": false
              }
            ],
            "submit_button": "Submit Report",
            "cancel_button": "Skip"
          }
        }
      ],
      "sender_id": "dom-550e8400-e29b-41d4-a716-446655440002",
      "receiver_id": "sub-550e8400-e29b-41d4-a716-446655440001",
      "status": "pending",
      "created_at": "2025-09-27T16:00:00.000Z",
      "sender": {
        "id": "dom-550e8400-e29b-41d4-a716-446655440002",
        "screen_name": "MasterUser",
        "login_name": "master_dom"
      }
    }
  }
}
```

**Client submits form data:**
```json
{
  "type": "command_completion",
  "id": "input-completion-12345",
  "timestamp": "2025-09-27T16:05:30.000Z",
  "from": "sub-550e8400-e29b-41d4-a716-446655440001",
  "to": "dom-550e8400-e29b-41d4-a716-446655440002",
  "data": {
    "command_id": "input-550e8400-e29b-41d4-a716-446655440005",
    "status": "completed",
    "response": {
      "success": true,
      "timestamp": "2025-09-27T16:05:30.000Z",
      "data": {
        "form_data": {
          "tasks_completed": 5,
          "mood_rating": "Good",
          "notes": "Had a productive day, completed all assigned tasks."
        }
      }
    }
  }
}
```

## Error Examples

### Invalid Command Format

**Client sends malformed message:**
```json
{
  "type": "command_completion",
  "id": "bad-msg-12345",
  "timestamp": "2025-09-27T17:00:00.000Z",
  "from": "sub-550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "command_id": "non-existent-command"
  }
}
```

**Server responds with error:**
```json
{
  "type": "error",
  "id": "error-response-12345",
  "timestamp": "2025-09-27T17:00:01.000Z",
  "from": "00000000-0000-0000-0000-000000000000",
  "to": "sub-550e8400-e29b-41d4-a716-446655440001",
  "data": {
    "error_code": "COMMAND_NOT_FOUND",
    "error_message": "Command with ID 'non-existent-command' does not exist",
    "reference_id": "bad-msg-12345"
  }
}
```

### Authentication Error

**Client with invalid token:**
```json
{
  "type": "error",
  "id": "auth-error-12345",
  "timestamp": "2025-09-27T17:05:00.000Z",
  "from": "00000000-0000-0000-0000-000000000000",
  "to": "unknown",
  "data": {
    "error_code": "AUTHENTICATION_FAILED",
    "error_message": "Invalid or expired JWT token",
    "details": "Token signature verification failed"
  }
}
```

## Processing Instructions Array

**Important:** The `instructions` field contains a JSON array of instruction objects that can be accessed directly:

### Sequential Execution (Recommended)

For multi-instruction commands, execute instructions in order:

```javascript
// JavaScript example - Sequential execution
const command = message.data.command;
const instructions = command.instructions; // Already a JSON array

async function executeInstructions(instructions) {
  const results = [];
  
  for (let i = 0; i < instructions.length; i++) {
    const instruction = instructions[i];
    console.log(`Executing ${i+1}/${instructions.length}: ${instruction.type}`);
    
    try {
      const result = await executeInstruction(instruction);
      results.push(result);
    } catch (error) {
      console.error(`Failed at instruction ${i}:`, error);
      throw error;
    }
  }
  
  return results;
}

async function executeInstruction(instruction) {
  switch(instruction.type) {
    case 'std_popup':
      return await showPopup(instruction.content);
    case 'std_timer':
      return await startTimer(instruction.content);
    case 'std_input':
      return await showInputForm(instruction.content);
    // ... handle other types
  }
}
```

### Parallel Execution (Advanced)

Only for instructions that can run simultaneously:

```javascript
// JavaScript example - Parallel execution (use carefully)
const command = message.data.command;
const instructions = command.instructions; // Already a JSON array

// Only for compatible instruction types
const promises = instructions.map(instruction => {
  return executeInstruction(instruction);
});

const results = await Promise.all(promises);
```

### Python Example

```python
# Python example - Sequential execution
import asyncio

async def execute_command(command):
    instructions = command['instructions']  # Already a list
    results = []
    
    for i, instruction in enumerate(instructions):
        print(f"Executing {i+1}/{len(instructions)}: {instruction['type']}")
        
        try:
            if instruction['type'] == 'std_popup':
                result = await show_popup(instruction['content'])
            elif instruction['type'] == 'std_timer':
                result = await start_timer(instruction['content'])
            elif instruction['type'] == 'std_input':
                result = await show_input_form(instruction['content'])
            
            results.append(result)
            
        except Exception as error:
            print(f"Failed at instruction {i}: {error}")
            raise error
    
    return results
```

## Real Message Size Examples

**Small popup message:** ~800 bytes
**Complex form message:** ~2,500 bytes  
**Multi-instruction command:** ~4,000 bytes

The WebSocket connection has a maximum message size limit of 512 bytes for client-to-server messages, but server-to-client messages can be larger for command delivery.

## Custom Content Structures

**Important:** The `content` field in each instruction is completely arbitrary and can contain any JSON structure you need. The server acts as a passthrough - it stores and forwards the content without validation or modification.

**Examples of custom content structures:**

```json
// Gaming instruction
{
  "type": "play_game",
  "content": {
    "game": "simon_says",
    "sequence": [1, 3, 2, 4, 1],
    "timeout_seconds": 30,
    "reward_points": 100
  }
}

// IoT device control
{
  "type": "device_control", 
  "content": {
    "device_id": "bedroom_lights",
    "actions": [
      {"command": "brightness", "value": 75},
      {"command": "color", "value": "#FF69B4"},
      {"command": "fade_duration", "value": 2000}
    ]
  }
}

// Media playback
{
  "type": "media_command",
  "content": {
    "action": "play_playlist",
    "source": "spotify",
    "playlist_id": "workout_mix_2024",
    "shuffle": true,
    "volume": 0.6
  }
}
```

Clients should implement handling for their supported instruction types and gracefully ignore or report unsupported types.