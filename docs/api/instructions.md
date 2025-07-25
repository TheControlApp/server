# Instruction Types

Instructions are the building blocks of commands. Each has a `type` and `content` field.

## Standard Types

### `popup-msg`
Display modal message with button.
```json
{
  "type": "popup-msg",
  "content": {
    "body": "Your message here",
    "button": "OK"
  }
}
```

### `download-file`
Trigger file download using file_hash and file_name.
```json
{
  "type": "download-file", 
  "content": {
    "file_hash": "abc123...",
    "file_name": "document.pdf"
  }
}
```

### `display-text`
Show formatted text content.
```json
{
  "type": "display-text",
  "content": {
    "text": "Content here",
    "format": "plain"
  }
}
```

### `timer`
Start countdown timer with duration and title.
```json
{
  "type": "timer",
  "content": {
    "duration": 300,
    "title": "Break Time"
  }
}
```

### `notification`
System notification with title and body.
```json
{
  "type": "notification",
  "content": {
    "title": "Alert",
    "body": "Something happened"
  }
}
```

### `open-url`
Open URL in browser.
```json
{
  "type": "open-url",
  "content": {
    "url": "https://example.com"
  }
}
```

### `form-input`
Collect user input with title and fields array.
```json
{
  "type": "form-input",
  "content": {
    "title": "Feedback",
    "fields": [
      {
        "name": "rating",
        "label": "Rate 1-5",
        "type": "number"
      }
    ]
  }
}
```

## Extended Types

### `announcement`
Broadcast message with title, body, and priority.
```json
{
  "type": "announcement",
  "content": {
    "title": "System Update",
    "body": "Maintenance tonight",
    "priority": "normal"
  }
}
```

### `schedule-task`
Schedule future command execution.
```json
{
  "type": "schedule-task",
  "content": {
    "execute_at": "2024-01-01T10:00:00Z",
    "command": { }
  }
}
```

## Implementation

Handle instructions by checking the instruction type and processing the content accordingly.

## Custom Instructions

Add custom types by defining type name/content structure, adding handler logic, and documenting for other developers.

Show formatted text content to the user.

**Format:**
```json
{
  "type": "display-text",
  "content": {
    "text": "Remember to complete your evening routine before bed.",
    "format": "plain"
  }
}
```

**Formats:**
- `plain` - Plain text
- `markdown` - Markdown formatting
- `html` - Basic HTML (sanitized)

### 4. Timer (`timer`)

Create a countdown timer for tasks.

**Format:**
```json
{
  "type": "timer",
  "content": {
    "duration": 300,
    "title": "Break Time",
    "message": "Take a 5-minute break"
  }
}
```

**Implementation:**
- `duration` in seconds
- Display countdown prominently
- Optional notification when complete

### 5. URL Launch (`open-url`)

Open a URL in browser or in-app.

**Format:**
```json
{
  "type": "open-url",
  "content": {
    "url": "https://example.com/task",
    "display": "external"
  }
}
```

**Display Options:**
- `external` - Open in system browser
- `inline` - Open in app webview
- `tab` - Open in new tab (web apps)

### 6. Notification (`notification`)

Send a system notification.

**Format:**
```json
{
  "type": "notification",
  "content": {
    "title": "Task Reminder",
    "body": "You have a pending task to complete",
    "priority": "normal"
  }
}
```

**Priority Levels:**
- `low` - Silent notification
- `normal` - Standard notification
- `high` - Urgent notification with sound

### 7. Form Input (`form-input`)

Request input from the user.

**Format:**
```json
{
  "type": "form-input",
  "content": {
    "title": "Daily Check-in",
    "fields": [
      {
        "name": "mood",
        "label": "How are you feeling today?",
        "type": "select",
        "options": ["Great", "Good", "Okay", "Not great"],
        "required": true
      },
      {
        "name": "notes",
        "label": "Additional notes",
        "type": "textarea",
        "required": false
      }
    ],
    "submit_to": "master123"
  }
}
```

**Field Types:**
- `text` - Single line text
- `textarea` - Multi-line text
- `select` - Dropdown selection
- `radio` - Radio buttons
- `checkbox` - Multiple checkboxes
- `number` - Numeric input

### 8. Image Display (`display-image`)

Show an image to the user.

**Format:**
```json
{
  "type": "display-image",
  "content": {
    "file_hash": "def456abc789...",
    "caption": "Today's outfit requirement",
    "width": 400
  }
}
```

### 9. Audio Play (`play-audio`)

Play an audio file.

**Format:**
```json
{
  "type": "play-audio",
  "content": {
    "file_hash": "ghi789jkl012...",
    "autoplay": true,
    "loop": false,
    "volume": 0.8
  }
}
```

## Extended Instruction Types

### 10. Schedule Task (`schedule-task`)

Schedule a future task or reminder.

**Format:**
```json
{
  "type": "schedule-task",
  "content": {
    "title": "Evening Routine",
    "scheduled_for": "2025-07-25T20:00:00Z",
    "instructions": [
      {
        "type": "popup-msg",
        "content": {
          "body": "Time for your evening routine!",
          "button": "Start"
        }
      }
    ]
  }
}
```

### 11. Vibration Pattern (`vibrate`)

Control vibration on supported devices.

**Format:**
```json
{
  "type": "vibrate",
  "content": {
    "pattern": [200, 100, 200, 100, 500],
    "intensity": 75
  }
}
```

### 12. Photo Capture (`capture-photo`)

Request photo from device camera.

**Format:**
```json
{
  "type": "capture-photo",
  "content": {
    "purpose": "Verification photo",
    "front_camera": true,
    "submit_to": "master123"
  }
}
```

### 13. Location Request (`request-location`)

Request current location from user.

**Format:**
```json
{
  "type": "request-location",
  "content": {
    "purpose": "Check-in location",
    "accuracy": "high",
    "submit_to": "master123"
  }
}
```

### 14. Device Control (`device-control`)

Control connected smart devices.

**Format:**
```json
{
  "type": "device-control",
  "content": {
    "device_type": "smart_plug",
    "device_id": "bedroom_lamp",
    "action": "turn_on",
    "parameters": {
      "brightness": 50
    }
  }
}
```

### 15. Announcement (`announcement`)

System-wide announcements with priority.

**Format:**
```json
{
  "type": "announcement",
  "content": {
    "title": "Server Maintenance",
    "body": "The server will be down for maintenance from 2-3 PM EST",
    "priority": "high",
    "expires_at": "2025-07-25T15:00:00Z"
  }
}
```

## Custom Instruction Types

Clients can implement custom instruction types by:

1. **Using the `custom` type:**
```json
{
  "type": "custom",
  "content": {
    "custom_type": "your_custom_type",
    "data": {
      // Your custom data structure
    }
  }
}
```

2. **Registering handlers:**
```csharp
// Add to your switch statement
case "your_custom_type":
    // Handle your custom instruction
    Console.WriteLine("Custom instruction received");
    break;
```

## Implementation Guidelines

### Basic Instruction Handling

```csharp
// Simple instruction processing
private void HandleInstruction(JObject instruction)
{
    var type = instruction["type"]?.ToString();
    var content = instruction["content"] as JObject;
    
    switch (type)
    {
        case "popup-msg":
            var message = content["body"]?.ToString();
            Console.WriteLine($"POPUP: {message}");
            break;
            
        case "notification":
            var title = content["title"]?.ToString();
            Console.WriteLine($"NOTIFICATION: {title}");
            break;
            
        case "timer":
            var duration = content["duration"]?.Value<int>() ?? 0;
            Console.WriteLine($"TIMER: {duration} seconds");
            break;
            
        default:
            Console.WriteLine($"Unknown type: {type}");
            break;
    }
}
```

### Processing Commands

```csharp
// Handle incoming command message
private void ProcessMessage(string json)
{
    var message = JObject.Parse(json);
    
    if (message["type"]?.ToString() == "command")
    {
        var instructions = message["data"]["instructions"] as JArray;
        
        foreach (var instruction in instructions)
        {
            HandleInstruction(instruction as JObject);
        }
    }
}
```

### Security Considerations

- Validate URLs before opening them
- Sanitize any HTML content 
- Verify file hashes are valid format
- Request permissions for camera/location appropriately
- Implement basic rate limiting for file downloads

### Testing

Use the create-commands tool to test instructions:

```bash
go run cmd/tools/create-commands/main.go
```
