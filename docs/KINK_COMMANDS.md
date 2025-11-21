# ControlApp Kink Command Specification

This document defines the kink command system for ControlApp - a consensual techdom application. All commands operate within a framework of explicit consent, safety, and user control.

## Safety & Consent Framework

### Consent Levels
- **Basic**: Safe commands only (message, notifications)
- **Standard**: Most commands allowed (file downloads, wallpapers, audio)
- **Full**: All commands including system-level operations (file execution, screen lock)

### Safety Requirements
1. **Explicit Consent**: Users must explicitly enable each command type
2. **Audit Logging**: All commands are logged with timestamps and results
3. **Emergency Controls**: Users can disable all commands via emergency hotkey
4. **File Type Restrictions**: Only whitelisted file types allowed
5. **Size Limits**: File downloads and uploads have configurable size limits

## Core Kink Commands (Tier 1)

### kink_message
Display a message dialog or notification to the user.

**Payload:**
```json
{
  "message": "Your message text here",
  "title": "Optional title",
  "style": "info|warning|error|question",
  "duration": 5000
}
```

**Consent Level:** Basic
**Safety Level:** Safe

### kink_open_link
Open a URL in the user's default browser.

**Payload:**
```json
{
  "url": "https://example.com",
  "confirm": true
}
```

**Consent Level:** Standard
**Safety Level:** Medium (URL validation required)

### kink_download_file
Download a file to the user's designated download folder.

**Payload:**
```json
{
  "url": "https://example.com/file.jpg",
  "filename": "optional_custom_name.jpg",
  "download_folder": "/path/to/downloads"
}
```

**Consent Level:** Standard
**Safety Level:** Medium (file type and size restrictions)

## Enhanced Kink Commands (Tier 2)

### kink_change_wallpaper
Change the user's desktop wallpaper.

**Payload:**
```json
{
  "image_path": "/path/to/image.jpg",
  "style": "fill|fit|stretch|tile|center",
  "restore_previous": true
}
```

**Consent Level:** Full
**Safety Level:** Medium (file validation required)

### kink_play_audio
Play an audio file or text-to-speech.

**Payload:**
```json
{
  "audio_path": "/path/to/audio.mp3",
  "volume": 0.8,
  "loop": false
}
```

**Consent Level:** Standard
**Safety Level:** Medium (volume limits, file type validation)

### kink_tts
Convert text to speech and play it.

**Payload:**
```json
{
  "text": "Text to speak",
  "voice": "default",
  "rate": 0,
  "volume": 80
}
```

**Consent Level:** Standard
**Safety Level:** Safe

### kink_popup_image
Display an image in a popup window.

**Payload:**
```json
{
  "image_path": "/path/to/image.jpg",
  "title": "Image Title",
  "duration": 0,
  "fullscreen": false
}
```

**Consent Level:** Standard
**Safety Level:** Medium (file validation required)

## Advanced Kink Commands (Tier 3)

### kink_run_file
Execute a file or program (HIGHEST RISK).

**Payload:**
```json
{
  "file_path": "/path/to/executable",
  "arguments": ["--arg1", "value1"],
  "work_dir": "/working/directory",
  "wait": true,
  "env": {"VAR": "value"}
}
```

**Consent Level:** Full + Explicit Per-Use Confirmation
**Safety Level:** HIGH RISK (requires explicit whitelist)

### kink_lock_screen
Lock the user's screen/workstation.

**Payload:**
```json
{
  "delay": 0,
  "message": "Optional lock message"
}
```

**Consent Level:** Full
**Safety Level:** Medium (system-level operation)

### kink_timer_task
Schedule a delayed task execution.

**Payload:**
```json
{
  "delay_seconds": 300,
  "command_type": "kink_message",
  "command_payload": {...},
  "task_id": "unique_task_id"
}
```

**Consent Level:** Full
**Safety Level:** High (combines scheduling with other commands)

## Extended Kink Commands (Tier 4)

### kink_scheduled_task
Create a recurring scheduled task.

**Payload:**
```json
{
  "schedule": "0 */2 * * *",
  "command_type": "kink_message", 
  "command_payload": {...},
  "task_name": "reminder_task",
  "expires_at": "2024-12-31T23:59:59Z"
}
```

**Consent Level:** Full + Admin Approval
**Safety Level:** HIGH RISK (persistent system changes)

### kink_mood_check
Request user mood/status input.

**Payload:**
```json
{
  "prompt": "How are you feeling?",
  "options": ["Good", "Stressed", "Tired", "Excited"],
  "require_response": true,
  "timeout": 300
}
```

**Consent Level:** Standard
**Safety Level:** Safe (input collection only)

### kink_permission_request
Request elevated permissions for future commands.

**Payload:**
```json
{
  "command_types": ["kink_run_file", "kink_lock_screen"],
  "duration_minutes": 60,
  "reason": "Planned session activities"
}
```

**Consent Level:** Full + Interactive Confirmation
**Safety Level:** High (permission elevation)

## Implementation Guidelines

### Client-Side Implementation
1. **Command Validation**: Validate all payloads before execution
2. **Consent Checking**: Verify user consent before executing any command
3. **Audit Logging**: Log all commands with detailed metadata
4. **Error Handling**: Graceful failure with detailed error messages
5. **Safety Limits**: Enforce file size, duration, and frequency limits

### Server-Side Requirements
1. **Authentication**: All commands require valid authentication
2. **Rate Limiting**: Prevent command spam/abuse
3. **Audit Trail**: Maintain comprehensive audit logs
4. **User Permissions**: Respect user-configured permission levels
5. **Session Management**: Commands tied to active sessions

### Security Considerations
1. **Input Sanitization**: All user inputs must be sanitized
2. **Path Validation**: File paths must be validated and sandboxed
3. **Command Whitelisting**: Only approved commands allowed
4. **Emergency Shutdown**: Users can disable all commands immediately
5. **Consent Verification**: Regular re-confirmation of consent levels

## Configuration Example

```yaml
# Client configuration for kink commands
kink_commands:
  consent_level: "standard"
  
  # Per-command consent
  allowed_commands:
    - "kink_message"
    - "kink_open_link"
    - "kink_download_file"
    - "kink_play_audio"
    - "kink_tts"
    - "kink_popup_image"
  
  blocked_commands:
    - "kink_run_file"  # Explicitly blocked for safety
  
  require_confirmation:
    - "kink_lock_screen"
    - "kink_change_wallpaper"
  
  # Safety limits
  max_file_size_mb: 50
  allowed_file_types: [".jpg", ".jpeg", ".png", ".gif", ".mp3", ".wav"]
  download_folder: "downloads"
  
  # Audio settings
  max_volume: 80
  audio_duration_limit: 300  # 5 minutes
  
  # Emergency controls
  emergency_hotkey: "Ctrl+Alt+Shift+E"
  auto_disable_minutes: 120  # Auto-disable after 2 hours
```

## Error Codes

- `CONSENT_DENIED`: User has not consented to this command type
- `SAFETY_VIOLATION`: Command violates safety policies
- `FILE_NOT_ALLOWED`: File type or size not permitted
- `RATE_LIMITED`: Too many commands in short time period
- `INVALID_PAYLOAD`: Command payload validation failed
- `EXECUTION_FAILED`: Command execution encountered an error
- `PERMISSION_DENIED`: Insufficient permissions for command
- `EMERGENCY_STOP`: Commands disabled by emergency control

---

**Remember**: This system is designed for consensual adult activities. Always respect boundaries, maintain clear communication, and prioritize safety above all else.