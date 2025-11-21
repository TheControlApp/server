# ControlApp Client Development Guide

## Overview

This guide provides everything needed to build 3rd party clients for the ControlApp platform. The ControlApp server provides a flexible, modern API that supports both REST and WebSocket communication, enabling rich interactive applications.

## 🚀 Quick Start

### Minimum Viable Client

A basic ControlApp client needs to:

1. **Connect to the WebSocket API**
2. **Handle authentication** (or operate anonymously)  
3. **Process incoming commands**
4. **Send command results back**

```javascript
// Minimal JavaScript client
const client = new WebSocket('ws://localhost:3080/ws/client');

client.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === 'command') {
        processCommand(message.payload);
    }
};

function processCommand(command) {
    // Execute the command and send result back
    client.send(JSON.stringify({
        type: 'command_result',
        payload: { command_id: command.id, status: 'completed' }
    }));
}
```

## 🏗️ Architecture Overview

### Communication Protocols

| Protocol | Purpose | Authentication | Use Cases |
|----------|---------|----------------|-----------|
| **WebSocket** | Real-time commands | Optional (progressive) | Interactive clients, command execution |
| **REST API** | User management | Required (JWT) | Registration, login, user profiles |

### Client Types

| Type | Authentication | Capabilities | Examples |
|------|----------------|--------------|----------|
| **Anonymous** | None | Receive broadcasts only | Public displays, demos |
| **Authenticated** | JWT Token | Send/receive commands | Personal clients, apps |
| **Administrative** | Admin JWT | Manage users, moderate | Admin panels, moderation tools |

## 🔐 Authentication

### 1. REST API Authentication

First, authenticate via REST to get a JWT token:

```javascript
// Register a new user
const registerResponse = await fetch('http://localhost:3080/api/v1/auth/register', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        screen_name: 'MyClient',
        login_name: 'myclient_user',
        password: 'secure_password123'
    })
});

// Login to get JWT token
const loginResponse = await fetch('http://localhost:3080/api/v1/auth/login', {
    method: 'POST', 
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        login_name: 'myclient_user',
        password: 'secure_password123'
    })
});

const { token, user } = await loginResponse.json();
```

### 2. WebSocket Authentication

Three methods to authenticate WebSocket connections:

#### Method A: Query Parameter (Recommended)
```javascript
const ws = new WebSocket(`ws://localhost:3080/ws/client?token=${jwtToken}`);
```

#### Method B: Authorization Header
```javascript
const ws = new WebSocket('ws://localhost:3080/ws/client', [], {
    headers: { 'Authorization': `Bearer ${jwtToken}` }
});
```

#### Method C: Progressive Authentication
```javascript
// 1. Connect anonymously first
const ws = new WebSocket('ws://localhost:3080/ws/client');

// 2. Send authentication message after connection
ws.onopen = () => {
    ws.send(JSON.stringify({
        type: 'auth_token',
        payload: { token: jwtToken }
    }));
};

// 3. Handle authentication response
ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    if (message.type === 'auth_success') {
        console.log('Now authenticated!', message.payload);
    }
};
```

## 📨 Message Protocol

### Message Structure

All WebSocket messages follow this JSON structure:

```typescript
interface Message {
    type: string;           // Message type identifier
    payload: any;           // Message content (varies by type)
    timestamp?: string;     // Optional ISO timestamp
    message_id?: string;    // Optional unique identifier
}
```

### Core Message Types

| Type | Direction | Purpose | Auth Required |
|------|-----------|---------|---------------|
| `ping` | Client → Server | Keep connection alive | No |
| `pong` | Server → Client | Connection keepalive response | No |
| `auth_token` | Client → Server | Authenticate with JWT | No |
| `auth_success` | Server → Client | Authentication successful | No |
| `auth_error` | Server → Client | Authentication failed | No |
| `command` | Server → Client | Execute command | Yes (target) |
| `command_result` | Client → Server | Command execution result | Yes |
| `broadcast` | Server → Client | Broadcast message | No |
| `error` | Server → Client | Error notification | No |

## 🎯 Standard Command Set

### Core Commands (All clients should support)

#### 1. Display Commands

**`std_popup`** - Show popup message
```json
{
    "type": "std_popup",
    "content": {
        "body": "Hello! Please click OK to continue.",
        "button": "OK",
        "timeout": 30
    }
}
```

**`std_notification`** - Show notification
```json
{
    "type": "std_notification", 
    "content": {
        "title": "New Message",
        "body": "You have a new command to execute",
        "icon": "info",
        "duration": 5
    }
}
```

**`std_display_text`** - Display text content
```json
{
    "type": "std_display_text",
    "content": {
        "text": "Welcome to ControlApp!",
        "format": "markdown",
        "style": "info"
    }
}
```

#### 2. Interaction Commands

**`std_form_input`** - Collect user input
```json
{
    "type": "std_form_input",
    "content": {
        "title": "User Survey",
        "fields": [
            {
                "name": "rating",
                "label": "How do you feel today?", 
                "type": "select",
                "options": ["Great", "Good", "Okay", "Not great"],
                "required": true
            },
            {
                "name": "comments",
                "label": "Additional comments",
                "type": "textarea",
                "required": false
            }
        ],
        "submit_text": "Submit Response"
    }
}
```

**`std_choice`** - Present multiple choice
```json
{
    "type": "std_choice",
    "content": {
        "question": "What would you like to do?",
        "options": [
            { "id": "option1", "text": "Continue current task" },
            { "id": "option2", "text": "Start new task" },
            { "id": "option3", "text": "Take a break" }
        ],
        "allow_multiple": false,
        "timeout": 60
    }
}
```

#### 3. Timing Commands

**`std_timer`** - Start countdown timer
```json
{
    "type": "std_timer",
    "content": {
        "duration": 300,
        "title": "Focus Session",
        "message": "Time to concentrate!",
        "show_progress": true,
        "sound_alert": true
    }
}
```

**`std_schedule`** - Schedule future action
```json
{
    "type": "std_schedule",
    "content": {
        "execute_at": "2024-12-25T10:00:00Z",
        "command": {
            "type": "std_notification",
            "content": {
                "title": "Merry Christmas!",
                "body": "Hope you have a wonderful day!"
            }
        }
    }
}
```

#### 4. Media Commands

**`std_open_url`** - Open URL in browser/app
```json
{
    "type": "std_open_url",
    "content": {
        "url": "https://example.com",
        "target": "_blank",
        "confirm": true
    }
}
```

**`std_play_sound`** - Play audio notification
```json
{
    "type": "std_play_sound",
    "content": {
        "sound": "notification.wav",
        "volume": 0.8,
        "loop": false
    }
}
```

#### 5. System Commands

**`std_system_info`** - Request system information
```json
{
    "type": "std_system_info",
    "content": {
        "requested_fields": ["os", "version", "uptime", "memory"]
    }
}
```

### Optional Extended Commands

#### File Operations
- `ext_download_file` - Download file from URL
- `ext_upload_file` - Upload file to server  
- `ext_file_info` - Get file information

#### Advanced UI
- `ext_custom_ui` - Render custom UI component
- `ext_progress_bar` - Show progress indication
- `ext_table_display` - Display tabular data

#### Integration
- `ext_webhook_call` - Make HTTP webhook call
- `ext_api_request` - Make API request
- `ext_database_query` - Execute database query (admin only)

## 🛠️ Client Implementation Examples

### JavaScript/Web Client

```javascript
class ControlAppClient {
    constructor(serverUrl, token = null) {
        this.serverUrl = serverUrl;
        this.token = token;
        this.ws = null;
        this.authenticated = false;
        this.commandHandlers = new Map();
        this.setupStandardCommands();
    }
    
    // Connection management
    async connect() {
        const wsUrl = this.token 
            ? `${this.serverUrl}/ws/client?token=${this.token}`
            : `${this.serverUrl}/ws/client`;
            
        this.ws = new WebSocket(wsUrl);
        this.setupEventHandlers();
    }
    
    setupEventHandlers() {
        this.ws.onopen = () => {
            console.log('Connected to ControlApp server');
            this.startHeartbeat();
        };
        
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            this.handleMessage(message);
        };
        
        this.ws.onclose = () => {
            console.log('Disconnected from server');
            this.stopHeartbeat();
        };
        
        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }
    
    // Message handling
    handleMessage(message) {
        console.log('Received message:', message);
        
        switch(message.type) {
            case 'command':
                this.executeCommand(message.payload);
                break;
            case 'auth_success':
                this.authenticated = true;
                console.log('Authentication successful');
                break;
            case 'pong':
                console.log('Server heartbeat received');
                break;
            default:
                console.log('Unknown message type:', message.type);
        }
    }
    
    // Command execution
    executeCommand(command) {
        const handler = this.commandHandlers.get(command.type);
        if (handler) {
            handler(command).then(result => {
                this.sendCommandResult(command.id, 'completed', result);
            }).catch(error => {
                this.sendCommandResult(command.id, 'failed', { error: error.message });
            });
        } else {
            console.warn('Unknown command type:', command.type);
            this.sendCommandResult(command.id, 'failed', { error: 'Unknown command type' });
        }
    }
    
    sendCommandResult(commandId, status, result = {}) {
        this.send('command_result', {
            command_id: commandId,
            status: status,
            result: result,
            timestamp: new Date().toISOString()
        });
    }
    
    // Standard command implementations
    setupStandardCommands() {
        // Popup command
        this.commandHandlers.set('std_popup', async (command) => {
            const { body, button = 'OK', timeout = 30 } = command.content;
            
            return new Promise((resolve) => {
                const result = window.confirm(body);
                resolve({ button_clicked: result ? button : 'cancelled' });
            });
        });
        
        // Notification command  
        this.commandHandlers.set('std_notification', async (command) => {
            const { title, body, duration = 5 } = command.content;
            
            if ('Notification' in window) {
                const permission = await Notification.requestPermission();
                if (permission === 'granted') {
                    const notification = new Notification(title, { body });
                    setTimeout(() => notification.close(), duration * 1000);
                }
            }
            
            return { displayed: true };
        });
        
        // Timer command
        this.commandHandlers.set('std_timer', async (command) => {
            const { duration, title, message } = command.content;
            
            return new Promise((resolve) => {
                let remaining = duration;
                const interval = setInterval(() => {
                    remaining--;
                    console.log(`${title}: ${remaining}s remaining`);
                    
                    if (remaining <= 0) {
                        clearInterval(interval);
                        alert(message || 'Timer completed!');
                        resolve({ completed: true, duration });
                    }
                }, 1000);
            });
        });
        
        // Add more standard command handlers...
    }
    
    // Utility methods
    send(type, payload) {
        const message = {
            type,
            payload,
            timestamp: new Date().toISOString()
        };
        
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
        } else {
            console.error('WebSocket not connected');
        }
    }
    
    startHeartbeat() {
        this.heartbeatInterval = setInterval(() => {
            this.send('ping', {});
        }, 30000); // Every 30 seconds
    }
    
    stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
        }
    }
    
    // Authentication
    async authenticate(username, password) {
        const response = await fetch(`${this.serverUrl}/api/v1/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ login_name: username, password })
        });
        
        if (response.ok) {
            const { token } = await response.json();
            this.token = token;
            
            // Send auth message if already connected
            if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                this.send('auth_token', { token });
            }
            
            return token;
        } else {
            throw new Error('Authentication failed');
        }
    }
}

// Usage
const client = new ControlAppClient('ws://localhost:3080');
await client.connect();
await client.authenticate('username', 'password');
```

### Python Client

```python
import asyncio
import websockets
import json
import requests
from typing import Dict, Any, Callable

class ControlAppClient:
    def __init__(self, server_url: str, token: str = None):
        self.server_url = server_url
        self.token = token
        self.ws = None
        self.authenticated = False
        self.command_handlers: Dict[str, Callable] = {}
        self.setup_standard_commands()
    
    async def connect(self):
        """Connect to the WebSocket server"""
        ws_url = f"{self.server_url}/ws/client"
        if self.token:
            ws_url += f"?token={self.token}"
            
        self.ws = await websockets.connect(ws_url)
        await self.listen_for_messages()
    
    async def listen_for_messages(self):
        """Listen for incoming messages"""
        async for message_str in self.ws:
            try:
                message = json.loads(message_str)
                await self.handle_message(message)
            except json.JSONDecodeError as e:
                print(f"Failed to parse message: {e}")
    
    async def handle_message(self, message: Dict[str, Any]):
        """Handle incoming messages"""
        message_type = message.get('type')
        payload = message.get('payload', {})
        
        if message_type == 'command':
            await self.execute_command(payload)
        elif message_type == 'auth_success':
            self.authenticated = True
            print("Authentication successful")
        elif message_type == 'pong':
            print("Server heartbeat received")
        else:
            print(f"Unknown message type: {message_type}")
    
    async def execute_command(self, command: Dict[str, Any]):
        """Execute a received command"""
        command_type = command.get('type')
        handler = self.command_handlers.get(command_type)
        
        if handler:
            try:
                result = await handler(command)
                await self.send_command_result(
                    command['id'], 'completed', result
                )
            except Exception as e:
                await self.send_command_result(
                    command['id'], 'failed', {'error': str(e)}
                )
        else:
            print(f"Unknown command type: {command_type}")
            await self.send_command_result(
                command['id'], 'failed', {'error': 'Unknown command type'}
            )
    
    async def send_command_result(self, command_id: str, status: str, result: Dict[str, Any] = {}):
        """Send command execution result back to server"""
        await self.send('command_result', {
            'command_id': command_id,
            'status': status,
            'result': result,
            'timestamp': datetime.now().isoformat()
        })
    
    def setup_standard_commands(self):
        """Set up handlers for standard commands"""
        
        async def handle_popup(command):
            content = command['content']
            body = content.get('body', '')
            button = content.get('button', 'OK')
            
            print(f"POPUP: {body}")
            input(f"Press Enter to click '{button}'...")
            return {'button_clicked': button}
        
        async def handle_notification(command):
            content = command['content']
            title = content.get('title', 'Notification')
            body = content.get('body', '')
            
            print(f"NOTIFICATION - {title}: {body}")
            return {'displayed': True}
        
        async def handle_timer(command):
            content = command['content']
            duration = content.get('duration', 60)
            title = content.get('title', 'Timer')
            
            print(f"Starting timer: {title} ({duration}s)")
            await asyncio.sleep(duration)
            print(f"Timer completed: {title}")
            
            return {'completed': True, 'duration': duration}
        
        self.command_handlers.update({
            'std_popup': handle_popup,
            'std_notification': handle_notification, 
            'std_timer': handle_timer
        })
    
    async def send(self, message_type: str, payload: Dict[str, Any]):
        """Send message to server"""
        message = {
            'type': message_type,
            'payload': payload,
            'timestamp': datetime.now().isoformat()
        }
        
        if self.ws:
            await self.ws.send(json.dumps(message))
    
    async def authenticate(self, username: str, password: str):
        """Authenticate with username/password"""
        response = requests.post(f"{self.server_url}/api/v1/auth/login", json={
            'login_name': username,
            'password': password
        })
        
        if response.status_code == 200:
            self.token = response.json()['token']
            if self.ws:
                await self.send('auth_token', {'token': self.token})
            return self.token
        else:
            raise Exception('Authentication failed')

# Usage
async def main():
    client = ControlAppClient('ws://localhost:3080')
    await client.connect()
    await client.authenticate('username', 'password')

asyncio.run(main())
```

## 🧪 Testing Your Client

### 1. Manual Testing

Use our integration test server to send test commands:

```bash
# Start the ControlApp server
cd server && air

# Run the integration test (creates test users)
cd cmd/tools/integration-test && ./integration-test.exe

# Connect your client and authenticate with:
# Username: integtest1, Password: password123
```

### 2. Command Testing

Send test commands via the WebSocket console or REST API:

```javascript
// Send a test popup command
fetch('http://localhost:3080/api/v1/commands', {
    method: 'POST',
    headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
    },
    body: JSON.stringify({
        receiver_id: 'user-uuid-here',
        instructions: [{
            type: 'std_popup',
            content: {
                body: 'Test message from server!',
                button: 'Got it!'
            }
        }]
    })
});
```

### 3. Error Handling Testing

Test how your client handles various error scenarios:

- Invalid authentication tokens
- Network disconnections
- Unknown command types
- Malformed messages
- Server unavailability

## 📖 Best Practices

### 1. Error Handling
```javascript
// Always handle WebSocket errors gracefully
ws.onerror = (error) => {
    console.error('WebSocket error:', error);
    // Implement reconnection logic
    setTimeout(() => this.reconnect(), 5000);
};

// Validate command structure before execution
executeCommand(command) {
    if (!command.id || !command.type || !command.content) {
        console.error('Invalid command structure:', command);
        return;
    }
    // Process valid command...
}
```

### 2. User Experience
```javascript
// Provide feedback for long-running commands
async handle_timer(command) {
    const duration = command.content.duration;
    
    // Show progress to user
    for (let i = duration; i > 0; i--) {
        updateUI(`Timer: ${i}s remaining`);
        await sleep(1000);
    }
    
    showNotification('Timer completed!');
}
```

### 3. Security
```javascript
// Sanitize user inputs
function sanitizeInput(input) {
    return input.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '');
}

// Validate URLs before opening
function validateURL(url) {
    try {
        const parsed = new URL(url);
        return ['http:', 'https:'].includes(parsed.protocol);
    } catch {
        return false;
    }
}
```

### 4. Performance
```javascript
// Implement command queuing for busy periods
class CommandQueue {
    constructor() {
        this.queue = [];
        this.processing = false;
    }
    
    async add(command) {
        this.queue.push(command);
        if (!this.processing) {
            await this.process();
        }
    }
    
    async process() {
        this.processing = true;
        while (this.queue.length > 0) {
            const command = this.queue.shift();
            await this.executeCommand(command);
        }
        this.processing = false;
    }
}
```

## 🎨 UI/UX Guidelines

### Visual Consistency
- Use consistent colors and fonts
- Follow platform UI guidelines (Material Design, Human Interface Guidelines)
- Provide clear visual feedback for all interactions

### Accessibility
- Support keyboard navigation
- Provide screen reader compatibility
- Use sufficient color contrast
- Include alternative text for images

### Responsiveness
- Handle different screen sizes
- Optimize for mobile devices
- Provide offline functionality where possible

## 📦 Deployment

### Desktop Applications
- **Electron** (JavaScript/TypeScript)
- **Tauri** (Rust + Web)
- **Flutter** (Dart)
- **Qt** (C++)

### Mobile Applications
- **React Native** (JavaScript/TypeScript)
- **Flutter** (Dart)
- **Native iOS/Android**

### Web Applications
- **React/Vue/Angular** (SPA)
- **Progressive Web App** (PWA)
- **WebRTC** for peer-to-peer features

### IoT/Embedded
- **Python** (Raspberry Pi)
- **C/C++** (Arduino, ESP32)
- **Go** (Cross-platform)

## 🤝 Contributing

### Submitting New Standard Commands

1. **Propose the command** in GitHub Issues
2. **Define the specification** (type, content structure, expected behavior)
3. **Implement reference handler** in JavaScript and Python
4. **Create test cases** 
5. **Update documentation**
6. **Submit Pull Request**

### Sharing Your Client

We'd love to showcase community clients! Share your implementation:

1. **Create a README** with setup instructions
2. **Include screenshots/demos**
3. **Document unique features**
4. **Submit to our showcase page**

## 📚 Additional Resources

- **[WebSocket API Reference](WEBSOCKET_API.md)** - Complete WebSocket message documentation
- **[REST API Reference](REST_API.md)** - HTTP endpoint documentation  
- **[Error Handling Guide](ERROR_HANDLING.md)** - RFC 7807 error response reference
- **[Server Setup Guide](../README.md)** - How to run the ControlApp server
- **[Example Clients](../examples/)** - Reference implementations

## 💬 Community & Support

- **GitHub Issues** - Bug reports and feature requests
- **GitHub Discussions** - General questions and community chat
- **Discord Server** - Real-time community support (coming soon)

---

*Happy coding! We can't wait to see what amazing clients you build with ControlApp!* 🚀