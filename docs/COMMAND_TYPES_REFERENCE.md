# ControlApp Command Types Reference

This document defines ALL command types that ControlApp clients should implement, organized by priority and functionality.

## Implementation Priority for Clients

### 1. **Console Client** (Bare bones communication)
- Focus on: Core commands for testing connectivity and basic communication
- Commands: `std_ping`, `std_notification` (as console output), `kink_message` (as console output)

### 2. **GUI Clients** (Windows/Linux with systray)
- Focus on: Core + Standard + Basic Kink commands
- Commands: All Core, most Standard, Tier 1-2 Kink commands

### 3. **Mobile Clients** (Android/iOS)
- Focus on: Core + Mobile-appropriate Standard + Safe Kink commands
- Commands: Core + notification-based commands, safe kink commands only

---

## 🟢 CORE Commands (Level 1) - ALL CLIENTS MUST IMPLEMENT

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

### **std_popup**
Display a simple popup message.
```json
{
  "type": "std_popup", 
  "content": {
    "body": "Hello! Please acknowledge this message.",
    "title": "System Message",
    "timeout": 30
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
    "duration": 5
  }
}
```

### **std_display_text**
Display formatted text content.
```json
{
  "type": "std_display_text",
  "content": {
    "text": "Welcome to ControlApp!",
    "title": "Welcome",
    "closable": true
  }
}
```

---

## 🟡 STANDARD Commands (Level 2) - GUI CLIENTS SHOULD IMPLEMENT

### **std_choice**
Present multiple choice selection.
```json
{
  "type": "std_choice",
  "content": {
    "question": "What would you like to do next?",
    "options": [
      { "id": "continue", "text": "Continue current task" },
      { "id": "new_task", "text": "Start new task" }
    ],
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
    "title": "User Feedback",
    "fields": [
      {
        "name": "rating",
        "label": "Overall Rating",
        "type": "select",
        "options": ["Excellent", "Good", "Fair", "Poor"]
      }
    ]
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
    "show_progress": true
  }
}
```

### **std_open_url**
Open URL in browser.
```json
{
  "type": "std_open_url",
  "content": {
    "url": "https://example.com",
    "confirm": true
  }
}
```

---

## 🔵 KINK Commands (Consensual Techdom) - BASED ON USER CONSENT

### **Tier 1: Basic Kink (Consent Level: Basic)**

#### **kink_message**
Display a message dialog (more intrusive than std_notification).
```json
{
  "type": "kink_message",
  "content": {
    "message": "Your dominant wants your attention",
    "title": "Control Message",
    "style": "info|warning|error|question",
    "duration": 5000
  }
}
```

### **Tier 2: Standard Kink (Consent Level: Standard)**

#### **kink_open_link**
Open a URL (same as std_open_url but in kink context).
```json
{
  "type": "kink_open_link",
  "content": {
    "url": "https://example.com",
    "confirm": true
  }
}
```

#### **kink_download_file**
Download a file to designated folder.
```json
{
  "type": "kink_download_file",
  "content": {
    "url": "https://example.com/file.jpg",
    "filename": "optional_custom_name.jpg",
    "download_folder": "/path/to/downloads"
  }
}
```

#### **kink_play_audio**
Play an audio file.
```json
{
  "type": "kink_play_audio",
  "content": {
    "audio_path": "/path/to/audio.mp3",
    "volume": 0.8,
    "loop": false
  }
}
```

#### **kink_tts**
Convert text to speech and play it.
```json
{
  "type": "kink_tts",
  "content": {
    "text": "Text to speak",
    "voice": "default",
    "rate": 0,
    "volume": 80
  }
}
```

#### **kink_popup_image**
Display an image in a popup window.
```json
{
  "type": "kink_popup_image",
  "content": {
    "image_path": "/path/to/image.jpg",
    "title": "Image Title",
    "duration": 0,
    "fullscreen": false
  }
}
```

### **Tier 3: Advanced Kink (Consent Level: Full)**

#### **kink_change_wallpaper**
Change the user's desktop wallpaper.
```json
{
  "type": "kink_change_wallpaper",
  "content": {
    "image_path": "/path/to/image.jpg",
    "style": "fill|fit|stretch|tile|center",
    "restore_previous": true
  }
}
```

#### **kink_lock_screen**
Lock the user's screen/workstation.
```json
{
  "type": "kink_lock_screen",
  "content": {
    "delay": 0,
    "message": "Optional lock message"
  }
}
```

#### **kink_timer_task**
Schedule a delayed task execution.
```json
{
  "type": "kink_timer_task",
  "content": {
    "delay_seconds": 300,
    "command_type": "kink_message",
    "command_payload": {...},
    "task_id": "unique_task_id"
  }
}
```

### **Tier 4: High-Risk Kink (Consent Level: Full + Confirmation)**

#### **kink_run_file**
Execute a file or program (HIGHEST RISK).
```json
{
  "type": "kink_run_file",
  "content": {
    "file_path": "/path/to/executable",
    "arguments": ["--arg1", "value1"],
    "work_dir": "/working/directory",
    "wait": true,
    "env": {"VAR": "value"}
  }
}
```

#### **kink_mood_check**
Request user mood/status input.
```json
{
  "type": "kink_mood_check",
  "content": {
    "prompt": "How are you feeling?",
    "options": ["Good", "Stressed", "Tired", "Excited"],
    "require_response": true,
    "timeout": 300
  }
}
```

---

## 🟠 EXTENDED Commands (Level 3) - OPTIONAL FOR ADVANCED CLIENTS

### **ext_progress_bar**
Display progress indicator.
```json
{
  "type": "ext_progress_bar",
  "content": {
    "title": "Processing Data",
    "current": 45,
    "maximum": 100,
    "allow_cancel": true
  }
}
```

### **ext_file_download**
Advanced file download with progress.
```json
{
  "type": "ext_file_download",
  "content": {
    "url": "https://example.com/file.pdf",
    "filename": "document.pdf",
    "show_progress": true,
    "verify_checksum": true
  }
}
```

---

## Client Implementation Roadmap

### **Phase 1: Console Client**
**Purpose:** Basic communication testing, reference implementation
**Commands to implement:**
- `std_ping` ✅ 
- `std_notification` (as console output)
- `kink_message` (as console output)
- Basic connection, auth, command sending

### **Phase 2: GUI Clients (Windows/Linux)**
**Purpose:** Full desktop experience with systray daemon
**Commands to implement:**
- All Core commands (std_*)
- Basic Kink commands (Tier 1-2)
- System integration (notifications, file handling)
- Systray icon and background operation

### **Phase 3: Mobile Clients (Android)**
**Purpose:** Mobile-appropriate kink experience
**Commands to implement:**
- Core commands adapted for mobile
- Safe kink commands only (message, notification-based)
- Mobile-specific UI patterns
- Background service limitations

### **Phase 4: Advanced Features**
**Purpose:** Full kink command support
**Commands to implement:**
- Advanced Kink commands (Tier 3-4)
- Extended commands (ext_*)
- Platform-specific optimizations
- Advanced security features

---

## Security & Consent Framework

### **Command Whitelisting by Client Type**

#### Console Client
```go
DefaultAllowedCommands = []string{
    "std_ping",
    "std_notification", // Console output only
    "kink_message",     // Console output only
}
```

#### GUI Client (Desktop)
```go
DefaultAllowedCommands = []string{
    // All Core
    "std_ping", "std_popup", "std_notification", "std_display_text",
    // Basic Standard
    "std_choice", "std_open_url",
    // Basic Kink (user must enable)
    "kink_message", "kink_open_link", "kink_download_file",
    "kink_play_audio", "kink_tts", "kink_popup_image",
}

RequireConfirmation = []string{
    "kink_change_wallpaper", "kink_lock_screen", "kink_run_file",
}

BlockedByDefault = []string{
    "kink_run_file", // Too dangerous for default consent
}
```

#### Mobile Client
```go
DefaultAllowedCommands = []string{
    // Core mobile-appropriate
    "std_ping", "std_notification", "std_display_text",
    // Safe kink only
    "kink_message", "kink_open_link",
}

// Most kink commands blocked on mobile for safety
BlockedCommands = []string{
    "kink_run_file", "kink_lock_screen", "kink_change_wallpaper",
    "kink_play_audio", // May be allowed with user permission
}
```

This comprehensive command type definition provides the foundation for building our clients in the correct order with appropriate feature sets for each platform!