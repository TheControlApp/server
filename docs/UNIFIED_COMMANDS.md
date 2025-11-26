# ControlApp Unified Command Specification

This document defines the unified command set for ControlApp first-party clients. Commands are consolidated from both standard and kink functionality into a cohesive set.

## Core Commands (Level 1) - All Clients Must Implement

### **std_ping**
Test connectivity and measure latency.
```json
{
  "type": "std_ping",
  "content": {
    "timestamp": "2024-12-20T10:30:00Z"
  }
}
```

### **std_message**
Display a message to the user (consolidates std_popup + kink_message).
```json
{
  "type": "std_message",
  "content": {
    "message": "Your message text here",
    "title": "Optional title",
    "style": "info|warning|error|question",
    "duration": 5000,
    "modal": true
  }
}
```

### **std_notification**
Show a non-blocking notification.
```json
{
  "type": "std_notification",
  "content": {
    "title": "Notification Title",
    "body": "Notification message",
    "icon": "info|warning|error",
    "duration": 5,
    "sound": true
  }
}
```

## Standard Commands (Level 2) - GUI Clients Should Implement

### **std_open_url**
Open a URL in the default browser (consolidates std_open_url + kink_open_link).
```json
{
  "type": "std_open_url",
  "content": {
    "url": "https://example.com",
    "title": "Link title",
    "confirm": true,
    "confirmation_message": "Open this link?"
  }
}
```

### **std_download_file**
Download a file to the designated folder (consolidates ext_file_download + kink_download_file).
```json
{
  "type": "std_download_file",
  "content": {
    "url": "https://example.com/file.jpg",
    "filename": "optional_custom_name.jpg",
    "download_folder": "downloads",
    "show_progress": true,
    "verify_checksum": false
  }
}
```

### **std_play_audio**
Play an audio file or text-to-speech (consolidates kink_play_audio + kink_tts).
```json
{
  "type": "std_play_audio",
  "content": {
    "source_type": "file|tts",
    "audio_path": "/path/to/audio.mp3",
    "tts_text": "Text to speak",
    "volume": 0.8,
    "voice": "default",
    "rate": 0,
    "loop": false
  }
}
```

### **std_display_image**
Display an image (consolidates kink_popup_image).
```json
{
  "type": "std_display_image",
  "content": {
    "image_path": "/path/to/image.jpg",
    "title": "Image Title",
    "duration": 0,
    "fullscreen": false,
    "closable": true
  }
}
```

### **std_timer**
Start a countdown timer.
```json
{
  "type": "std_timer",
  "content": {
    "duration": 300,
    "title": "Timer Title",
    "message": "Timer message",
    "show_progress": true,
    "completion_message": "Timer completed!"
  }
}
```

## Advanced Commands (Level 3) - Optional/Consent Required

### **std_change_wallpaper**
Change the desktop wallpaper (Windows/Linux only).
```json
{
  "type": "std_change_wallpaper",
  "content": {
    "image_path": "/path/to/image.jpg",
    "style": "fill|fit|stretch|tile|center",
    "restore_previous": true
  }
}
```

### **std_lock_screen**
Lock the user's screen/workstation.
```json
{
  "type": "std_lock_screen",
  "content": {
    "delay": 0,
    "message": "Screen locked by ControlApp"
  }
}
```

### **std_choice**
Present multiple choice selection.
```json
{
  "type": "std_choice",
  "content": {
    "question": "What would you like to do?",
    "options": [
      { "id": "option1", "text": "Option 1" },
      { "id": "option2", "text": "Option 2" }
    ],
    "allow_multiple": false,
    "timeout": 60
  }
}
```

### **std_form_input**
Collect structured user input.
```json
{
  "type": "std_form_input",
  "content": {
    "title": "Input Form",
    "fields": [
      {
        "name": "field1",
        "label": "Field Label",
        "type": "text|textarea|select|checkbox",
        "required": true
      }
    ]
  }
}
```

## High-Risk Commands (Level 4) - Explicit Consent + Confirmation Required

### **std_execute_file**
Execute a file or program (HIGHEST RISK - consolidates kink_run_file).
```json
{
  "type": "std_execute_file",
  "content": {
    "file_path": "/path/to/executable",
    "arguments": ["--arg1", "value1"],
    "work_dir": "/working/directory",
    "wait": true,
    "env": {"VAR": "value"}
  }
}
```

## Client Implementation Priorities

### Console Client (Minimal)
```go
SupportedCommands = []string{
    "std_ping",
    "std_message",      // Console output only
    "std_notification", // Console output only
}
```

### GUI Client (Full-Featured)
```go
DefaultAllowedCommands = []string{
    // Core (always allowed)
    "std_ping", "std_message", "std_notification",
    
    // Standard (user can disable)
    "std_open_url", "std_download_file", "std_play_audio", 
    "std_display_image", "std_timer",
    
    // Advanced (user must explicitly enable)
    "std_choice", "std_form_input",
}

RequireExplicitConsent = []string{
    "std_change_wallpaper", "std_lock_screen", 
}

BlockedByDefault = []string{
    "std_execute_file", // Too dangerous - requires special permission
}
```

## Consent Management

### Client Configuration
```yaml
consent:
  # Core commands (always enabled)
  core_enabled: true
  
  # Standard commands (user configurable)
  allow_url_opening: true
  allow_file_downloads: true
  allow_audio_playback: true
  allow_image_display: true
  allow_timers: true
  
  # Advanced commands (explicit consent required)
  allow_wallpaper_changes: false
  allow_screen_locking: false
  allow_user_input: false
  
  # High-risk (special permission + confirmation)
  allow_file_execution: false

# Security settings
security:
  download_folder: "downloads"
  max_file_size_mb: 50
  allowed_file_types: [".jpg", ".jpeg", ".png", ".gif", ".mp3", ".wav"]
  confirm_high_risk: true
  emergency_hotkey: "Ctrl+Alt+Shift+E"
```

This consolidation removes redundant commands while keeping the essential functionality needed for a consensual control system.