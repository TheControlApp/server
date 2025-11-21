# ControlApp Standard Commands Specification

## Overview

This document defines the **official standard command set** that all ControlApp clients should implement. Commands are categorized by complexity level to help developers prioritize implementation.

## Command Categories

### 🟢 **CORE** (Level 1) - Must Implement
Essential commands that every client should support for basic functionality.

### 🟡 **STANDARD** (Level 2) - Should Implement  
Common commands that provide good user experience across most clients.

### 🟠 **EXTENDED** (Level 3) - May Implement
Advanced commands for specialized clients or enhanced functionality.

### 🔴 **EXPERIMENTAL** (Level 4) - Optional
Cutting-edge features that may become standard in future versions.

---

## 🟢 CORE Commands (Level 1)

### **std_ping**
Test connectivity and measure latency.

```json
{
    "type": "std_ping",
    "content": {
        "timestamp": "2024-12-20T10:30:00Z",
        "expect_pong": true
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-123",
        "status": "completed",
        "result": {
            "pong_timestamp": "2024-12-20T10:30:00.123Z",
            "latency_ms": 123
        }
    }
}
```

### **std_popup**
Display a simple popup message with acknowledgment.

```json
{
    "type": "std_popup", 
    "content": {
        "body": "Hello! Please acknowledge this message.",
        "title": "System Message",
        "button": "OK",
        "timeout": 30,
        "priority": "normal"
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-124",
        "status": "completed",
        "result": {
            "button_clicked": "OK",
            "response_time_ms": 2500,
            "timed_out": false
        }
    }
}
```

### **std_notification**
Show a non-blocking notification.

```json
{
    "type": "std_notification",
    "content": {
        "title": "New Message",
        "body": "You have received a new command",
        "icon": "info",
        "duration": 5,
        "sound": true,
        "priority": "normal"
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result", 
    "payload": {
        "command_id": "cmd-125",
        "status": "completed",
        "result": {
            "displayed": true,
            "display_method": "system_notification",
            "duration_shown": 5
        }
    }
}
```

### **std_display_text**
Display formatted text content.

```json
{
    "type": "std_display_text",
    "content": {
        "text": "Welcome to ControlApp!\n\nYou are now connected.",
        "format": "markdown",
        "style": "info",
        "title": "Welcome",
        "closable": true,
        "timeout": 0
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-126", 
        "status": "completed",
        "result": {
            "displayed": true,
            "closed_by": "user",
            "display_duration": 12500
        }
    }
}
```

---

## 🟡 STANDARD Commands (Level 2)

### **std_choice**
Present multiple choice selection.

```json
{
    "type": "std_choice",
    "content": {
        "question": "What would you like to do next?",
        "options": [
            { "id": "continue", "text": "Continue current task", "description": "Keep working on what you're doing" },
            { "id": "new_task", "text": "Start new task", "description": "Begin something different" },
            { "id": "break", "text": "Take a break", "description": "Rest for a while" }
        ],
        "allow_multiple": false,
        "required": true,
        "timeout": 60,
        "default_selection": "continue"
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-127",
        "status": "completed", 
        "result": {
            "selected_options": ["new_task"],
            "selection_time_ms": 4200,
            "timed_out": false
        }
    }
}
```

### **std_form_input**
Collect structured user input.

```json
{
    "type": "std_form_input",
    "content": {
        "title": "User Feedback",
        "description": "Please provide your feedback on the current session",
        "fields": [
            {
                "name": "rating",
                "label": "Overall Rating",
                "type": "select",
                "options": ["Excellent", "Good", "Fair", "Poor"],
                "required": true
            },
            {
                "name": "comments",
                "label": "Additional Comments",
                "type": "textarea",
                "placeholder": "Any additional thoughts...",
                "required": false,
                "max_length": 500
            },
            {
                "name": "recommend",
                "label": "Would you recommend this?",
                "type": "checkbox",
                "required": false
            }
        ],
        "submit_text": "Submit Feedback",
        "cancel_text": "Skip",
        "allow_cancel": true
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-128",
        "status": "completed",
        "result": {
            "submitted": true,
            "form_data": {
                "rating": "Good",
                "comments": "This works really well!",
                "recommend": true
            },
            "completion_time_ms": 45000
        }
    }
}
```

### **std_timer**
Start a countdown timer with progress display.

```json
{
    "type": "std_timer",
    "content": {
        "duration": 300,
        "title": "Focus Session",
        "message": "Time to concentrate on your task!",
        "completion_message": "Great job! Take a short break.",
        "show_progress": true,
        "show_remaining": true,
        "sound_alert": true,
        "allow_pause": true,
        "allow_cancel": true,
        "style": "focus"
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-129",
        "status": "completed",
        "result": {
            "completed": true,
            "total_duration": 300,
            "actual_duration": 298,
            "paused_count": 0,
            "cancelled": false
        }
    }
}
```

### **std_open_url**
Open URL in browser or in-app view.

```json
{
    "type": "std_open_url",
    "content": {
        "url": "https://example.com",
        "title": "External Resource",
        "target": "_blank",
        "confirm": true,
        "confirmation_message": "This will open an external website. Continue?",
        "track_navigation": false
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-130",
        "status": "completed",
        "result": {
            "opened": true,
            "confirmed_by_user": true,
            "open_method": "system_browser"
        }
    }
}
```

### **std_schedule**
Schedule a command for future execution.

```json
{
    "type": "std_schedule",
    "content": {
        "execute_at": "2024-12-25T10:00:00Z",
        "timezone": "UTC",
        "command": {
            "type": "std_notification",
            "content": {
                "title": "Scheduled Reminder", 
                "body": "This is your scheduled reminder!"
            }
        },
        "allow_reschedule": true,
        "confirm_schedule": true
    }
}
```

**Expected Response:**
```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-131",
        "status": "completed",
        "result": {
            "scheduled": true,
            "schedule_id": "sched-456",
            "execute_at": "2024-12-25T10:00:00Z",
            "confirmed_by_user": true
        }
    }
}
```

---

## 🟠 EXTENDED Commands (Level 3)

### **ext_progress_bar**
Display progress indicator for long operations.

```json
{
    "type": "ext_progress_bar",
    "content": {
        "title": "Processing Data",
        "current": 45,
        "maximum": 100,
        "unit": "items",
        "estimated_remaining": "2m 30s",
        "allow_cancel": true,
        "show_details": true,
        "update_frequency": 1000
    }
}
```

### **ext_file_download**
Download file from specified URL.

```json
{
    "type": "ext_file_download",
    "content": {
        "url": "https://example.com/file.pdf",
        "filename": "document.pdf",
        "destination": "downloads",
        "overwrite": false,
        "show_progress": true,
        "verify_checksum": true,
        "checksum": "sha256:abc123..."
    }
}
```

### **ext_custom_ui**
Render custom UI component with HTML/CSS.

```json
{
    "type": "ext_custom_ui", 
    "content": {
        "html": "<div class='custom-widget'>...</div>",
        "css": ".custom-widget { color: blue; }",
        "javascript": "console.log('Custom UI loaded');",
        "sandbox": true,
        "permissions": ["notifications"],
        "size": { "width": 400, "height": 300 },
        "resizable": true
    }
}
```

### **ext_system_info**
Request detailed system information.

```json
{
    "type": "ext_system_info", 
    "content": {
        "requested_fields": [
            "os", "version", "architecture", 
            "memory", "cpu", "uptime", "network"
        ],
        "include_sensitive": false,
        "format": "json"
    }
}
```

### **ext_webhook_call**
Make HTTP webhook request.

```json
{
    "type": "ext_webhook_call",
    "content": {
        "url": "https://api.example.com/webhook",
        "method": "POST",
        "headers": {
            "Content-Type": "application/json",
            "Authorization": "Bearer token123"
        },
        "body": {
            "event": "command_executed",
            "timestamp": "2024-12-20T10:30:00Z"
        },
        "timeout": 30,
        "retry_count": 3
    }
}
```

---

## 🔴 EXPERIMENTAL Commands (Level 4)

### **exp_ai_chat**
Integration with AI chat services.

```json
{
    "type": "exp_ai_chat",
    "content": {
        "service": "openai",
        "model": "gpt-3.5-turbo",
        "prompt": "Help me plan my day",
        "context": {
            "user_preferences": {...},
            "previous_commands": [...]
        },
        "max_tokens": 150,
        "temperature": 0.7
    }
}
```

### **exp_ar_overlay**
Augmented reality overlay display.

```json
{
    "type": "exp_ar_overlay",
    "content": {
        "objects": [
            {
                "type": "text",
                "content": "Look here!",
                "position": { "x": 0.5, "y": 0.3, "z": 1.0 },
                "color": "#ff0000"
            }
        ],
        "duration": 10,
        "interaction_enabled": true
    }
}
```

### **exp_biometric_auth**
Request biometric authentication.

```json
{
    "type": "exp_biometric_auth",
    "content": {
        "methods": ["fingerprint", "face", "voice"],
        "challenge": "Confirm your identity",
        "timeout": 30,
        "fallback_to_password": true
    }
}
```

---

## Command Response Standards

### Status Values
- `"completed"` - Command executed successfully
- `"failed"` - Command execution failed
- `"cancelled"` - Command cancelled by user or system
- `"timeout"` - Command timed out
- `"partial"` - Command partially completed
- `"deferred"` - Command scheduled for later execution

### Error Response Format
When a command fails, include detailed error information:

```json
{
    "type": "command_result",
    "payload": {
        "command_id": "cmd-123",
        "status": "failed",
        "error": {
            "code": "PERMISSION_DENIED",
            "message": "User denied notification permission",
            "details": {
                "permission": "notifications",
                "user_action": "denied"
            },
            "recoverable": true,
            "retry_suggested": false
        }
    }
}
```

### Progress Updates
For long-running commands, send periodic progress updates:

```json
{
    "type": "command_progress",
    "payload": {
        "command_id": "cmd-123",
        "progress": {
            "current": 45,
            "maximum": 100,
            "message": "Processing item 45 of 100",
            "estimated_remaining_ms": 15000
        }
    }
}
```

## Implementation Guidelines

### 1. **Graceful Degradation**
If a client cannot fully support a command, it should:
- Execute what it can
- Report which parts were skipped
- Suggest alternatives if available

### 2. **User Consent**
Commands that access sensitive data or perform potentially disruptive actions should:
- Request explicit user permission
- Explain what will happen
- Allow the user to cancel

### 3. **Accessibility**
All commands should support:
- Screen readers and assistive technologies
- Keyboard navigation
- High contrast modes
- Customizable font sizes

### 4. **Platform Integration**
Commands should integrate naturally with the target platform:
- Follow platform UI guidelines
- Use native controls when possible
- Respect platform preferences (dark mode, etc.)

### 5. **Security Considerations**
- Validate all URLs before opening
- Sanitize HTML content in custom UI
- Limit file access to appropriate directories
- Never execute arbitrary code without user consent

## Testing Command Support

### Capability Advertisement
Clients should advertise their supported commands on connection:

```json
{
    "type": "client_capabilities",
    "payload": {
        "supported_commands": {
            "core": ["std_ping", "std_popup", "std_notification", "std_display_text"],
            "standard": ["std_choice", "std_form_input", "std_timer"],
            "extended": ["ext_progress_bar", "ext_file_download"],
            "experimental": []
        },
        "client_info": {
            "name": "AwesomeControlClient",
            "version": "1.2.0", 
            "platform": "web",
            "capabilities": {
                "notifications": true,
                "file_access": false,
                "system_info": true
            }
        }
    }
}
```

### Command Testing Suite
Use this test suite to verify command implementation:

```javascript
const testCommands = [
    {
        type: "std_ping",
        content: { timestamp: new Date().toISOString() }
    },
    {
        type: "std_popup", 
        content: {
            body: "Test popup message",
            button: "Test OK",
            timeout: 10
        }
    },
    {
        type: "std_notification",
        content: {
            title: "Test Notification",
            body: "This is a test notification",
            duration: 3
        }
    }
    // Add more test commands...
];

// Test each command
for (const command of testCommands) {
    const result = await sendTestCommand(command);
    console.log(`${command.type}: ${result.status}`);
}
```

## Versioning and Compatibility

### Command Versioning
Commands may evolve over time. Clients should handle version differences gracefully:

```json
{
    "type": "std_popup",
    "version": "2.1",
    "content": {
        "body": "This popup uses version 2.1 features",
        "new_feature_v2_1": "enhanced_styling"
    }
}
```

### Backward Compatibility
- New optional fields may be added to existing commands
- Required fields will never be removed
- Major breaking changes will result in new command types
- Deprecation warnings will be provided for 6 months minimum

### Forward Compatibility
Clients should ignore unknown fields in command content to support future enhancements.

---

## Command Registry

This specification is maintained in the ControlApp server repository. To propose new standard commands:

1. **Create GitHub Issue** with command specification
2. **Provide use cases** and justification
3. **Include reference implementation** (JavaScript + Python)
4. **Write comprehensive tests**
5. **Submit Pull Request** with documentation updates

---

*Last Updated: November 2025*  
*Specification Version: 1.0*  
*Compatible with ControlApp Server 1.0+*