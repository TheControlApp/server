# Client Integration Guide

## Overview
This guide provides examples for integrating with the ControlMe API across different platforms and programming languages.

## Authentication Flow

### JavaScript/Web Client
```javascript
class ControlMeClient {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
    this.authToken = null;
    this.websocket = null;
  }

  async login(loginName, password) {
    const response = await fetch(`${this.baseUrl}/api/v1/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ login_name: loginName, password })
    });
    
    const data = await response.json();
    if (data.token) {
      this.authToken = data.token;
      this.connectWebSocket();
    }
    return data;
  }

  connectWebSocket() {
    this.websocket = new WebSocket(`${this.baseUrl.replace('http', 'ws')}/api/ws`);
    this.websocket.onopen = () => {
      // Send authentication
      this.websocket.send(JSON.stringify({
        type: 'auth',
        token: this.authToken
      }));
    };
    
    this.websocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };
  }

  handleMessage(message) {
    switch (message.type) {
      case 'command':
        this.handleCommand(message.data);
        break;
      case 'error':
        console.error('WebSocket error:', message.message);
        break;
    }
  }

  handleCommand(command) {
    // Execute instructions sequentially
    command.instructions.forEach((instruction, index) => {
      setTimeout(() => {
       .executeInstruction(instruction, command);
      }, index * 1000);
    });
  }

  executeInstruction(instruction, command) {
    const handlers = {
      'popup-msg': this.handlePopupMessage.bind(this),
      'popup-web': this.handlePopupWeb.bind(this),
      'popup-video': this.handlePopupVideo.bind(this),
      'download-file': this.handleFileDownload.bind(this),
      'custom-notification': this.handleCustomNotification.bind(this)
    };

    const handler = handlers[instruction.type];
    if (handler) {
      handler(instruction.content, command);
    } else {
      console.warn(`Unknown instruction type: ${instruction.type}`);
    }
  }

  handlePopupMessage(content, command) {
    const modal = document.createElement('div');
    modal.className = 'command-modal';
    modal.innerHTML = `
      <div class="modal-content">
        <p>${content.body}</p>
        <button onclick="this.parentElement.parentElement.remove()">
          ${content.button || 'OK'}
        </button>
        <small>From: ${command.sender.screen_name}</small>
      </div>
    `;
    document.body.appendChild(modal);
  }

  handlePopupWeb(content, command) {
    window.open(content, '_blank');
  }

  handlePopupVideo(content, command) {
    const video = document.createElement('video');
    video.src = content;
    video.controls = true;
    video.autoplay = true;
    
    const modal = document.createElement('div');
    modal.className = 'video-modal';
    modal.appendChild(video);
    document.body.appendChild(modal);
  }

  async handleFileDownload(content, command) {
    const { file_hash, file_name } = content;
    const downloadUrl = `${this.baseUrl}/api/v1/files?filehash=${file_hash}&filename=${encodeURIComponent(file_name)}`;
    
    try {
      const response = await fetch(downloadUrl, {
        headers: { 'Authorization': `Bearer ${this.authToken}` }
      });
      
      if (response.ok) {
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = file_name;
        a.click();
        
        URL.revokeObjectURL(url);
      }
    } catch (error) {
      console.error('Download failed:', error);
    }
  }

  sendCommand(instructions, receiverUsername = null, tags = []) {
    if (!this.websocket) return;
    
    const command = {
      type: 'send_command',
      data: {
        instructions,
        tags
      }
    };
    
    if (receiverUsername) {
      command.data.receiver = receiverUsername;
    }
    
    this.websocket.send(JSON.stringify(command));
  }

  async uploadFile(file) {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await fetch(`${this.baseUrl}/api/v1/files`, {
      method: 'POST',
      headers: { 'Authorization': `Bearer ${this.authToken}` },
      body: formData
    });
    
    return await response.json();
  }
}
```

## Python Client

```python
import asyncio
import websockets
import json
import aiohttp
import aiofiles

class ControlMeClient:
    def __init__(self, base_url):
        self.base_url = base_url
        self.auth_token = None
        self.websocket = None
        
    async def login(self, login_name, password):
        async with aiohttp.ClientSession() as session:
            async with session.post(
                f"{self.base_url}/api/v1/auth/login",
                json={"login_name": login_name, "password": password}
            ) as response:
                data = await response.json()
                if "token" in data:
                    self.auth_token = data["token"]
                    await self.connect_websocket()
                return data
    
    async def connect_websocket(self):
        ws_url = self.base_url.replace("http", "ws") + "/api/ws"
        headers = {"Authorization": f"Bearer {self.auth_token}"}
        
        self.websocket = await websockets.connect(ws_url, extra_headers=headers)
        
        # Start message handler
        asyncio.create_task(self.message_handler())
    
    async def message_handler(self):
        async for message in self.websocket:
            data = json.loads(message)
            await self.handle_message(data)
    
    async def handle_message(self, message):
        if message["type"] == "command":
            await self.handle_command(message["data"])
        elif message["type"] == "error":
            print(f"WebSocket error: {message['message']}")
    
    async def handle_command(self, command):
        for instruction in command["instructions"]:
            await self.execute_instruction(instruction, command)
            await asyncio.sleep(1)  # 1 second delay between instructions
    
    async def execute_instruction(self, instruction, command):
        handlers = {
            "popup-msg": self.handle_popup_message,
            "popup-web": self.handle_popup_web,
            "download-file": self.handle_file_download,
            "custom-hardware": self.handle_hardware_control,
        }
        
        handler = handlers.get(instruction["type"])
        if handler:
            await handler(instruction["content"], command)
        else:
            print(f"Unknown instruction type: {instruction['type']}")
    
    async def handle_popup_message(self, content, command):
        print(f"\n=== MESSAGE FROM {command['sender']['screen_name']} ===")
        print(content["body"])
        input(f"Press Enter to {content.get('button', 'continue')}...")
    
    async def handle_popup_web(self, content, command):
        import webbrowser
        webbrowser.open(content)
    
    async def handle_file_download(self, content, command):
        file_hash = content["file_hash"]
        file_name = content["file_name"]
        
        url = f"{self.base_url}/api/v1/files?filehash={file_hash}&filename={file_name}"
        headers = {"Authorization": f"Bearer {self.auth_token}"}
        
        async with aiohttp.ClientSession() as session:
            async with session.get(url, headers=headers) as response:
                if response.status == 200:
                    async with aiofiles.open(f"downloads/{file_name}", "wb") as f:
                        async for chunk in response.content.iter_chunked(8192):
                            await f.write(chunk)
                    print(f"Downloaded: {file_name}")
    
    async def send_command(self, instructions, receiver_username=None, tags=None):
        if not self.websocket:
            return
        
        command = {
            "type": "send_command",
            "data": {
                "instructions": instructions,
                "tags": tags or []
            }
        }
        
        if receiver_username:
            command["data"]["receiver"] = receiver_username
        
        await self.websocket.send(json.dumps(command))
    
    async def upload_file(self, file_path):
        async with aiohttp.ClientSession() as session:
            with open(file_path, "rb") as f:
                data = aiohttp.FormData()
                data.add_field("file", f, filename=file_path.split("/")[-1])
                
                headers = {"Authorization": f"Bearer {self.auth_token}"}
                
                async with session.post(
                    f"{self.base_url}/api/v1/files",
                    data=data,
                    headers=headers
                ) as response:
                    return await response.json()

# Usage example
async def main():
    client = ControlMeClient("http://localhost:8080")
    
    # Login
    result = await client.login("username", "password")
    if "token" in result:
        print("Logged in successfully!")
        
        # Send a command
        await client.send_command([
            {
                "type": "popup-msg",
                "content": {
                    "body": "Hello from Python client!",
                    "button": "Got it!"
                }
            }
        ], tags=["general"])
        
        # Keep connection alive
        await asyncio.sleep(60)

if __name__ == "__main__":
    asyncio.run(main())
```

## React Native Client

```javascript
import React, { useEffect, useState } from 'react';
import { Alert, Linking } from 'react-native';
import RNFetchBlob from 'rn-fetch-blob';

class ControlMeReactNativeClient {
  constructor(baseUrl) {
    this.baseUrl = baseUrl;
    this.authToken = null;
    this.websocket = null;
  }

  async login(loginName, password) {
    try {
      const response = await fetch(`${this.baseUrl}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ login_name: loginName, password })
      });
      
      const data = await response.json();
      if (data.token) {
        this.authToken = data.token;
        this.connectWebSocket();
      }
      return data;
    } catch (error) {
      console.error('Login failed:', error);
      return { error: error.message };
    }
  }

  connectWebSocket() {
    const wsUrl = this.baseUrl.replace('http', 'ws') + '/api/ws';
    this.websocket = new WebSocket(wsUrl, [], {
      headers: { Authorization: `Bearer ${this.authToken}` }
    });
    
    this.websocket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      this.handleMessage(message);
    };
  }

  handleMessage(message) {
    if (message.type === 'command') {
      this.handleCommand(message.data);
    }
  }

  async handleCommand(command) {
    for (let i = 0; i < command.instructions.length; i++) {
      await this.executeInstruction(command.instructions[i], command);
      if (i < command.instructions.length - 1) {
        await new Promise(resolve => setTimeout(resolve, 1000));
      }
    }
  }

  async executeInstruction(instruction, command) {
    const { type, content } = instruction;
    
    switch (type) {
      case 'popup-msg':
        Alert.alert(
          `Message from ${command.sender.screen_name}`,
          content.body,
          [{ text: content.button || 'OK' }]
        );
        break;
        
      case 'popup-web':
        Linking.openURL(content);
        break;
        
      case 'download-file':
        await this.handleFileDownload(content);
        break;
        
      case 'custom-vibration':
        if (content.pattern) {
          Vibration.vibrate(content.pattern);
        }
        break;
        
      default:
        console.warn(`Unknown instruction type: ${type}`);
    }
  }

  async handleFileDownload(content) {
    const { file_hash, file_name } = content;
    const downloadUrl = `${this.baseUrl}/api/v1/files?filehash=${file_hash}&filename=${encodeURIComponent(file_name)}`;
    
    try {
      const dirs = RNFetchBlob.fs.dirs;
      const downloadPath = `${dirs.DownloadDir}/${file_name}`;
      
      const response = await RNFetchBlob.fetch('GET', downloadUrl, {
        Authorization: `Bearer ${this.authToken}`
      });
      
      await RNFetchBlob.fs.writeFile(downloadPath, response.data, 'base64');
      
      Alert.alert('Download Complete', `File saved: ${file_name}`);
    } catch (error) {
      Alert.alert('Download Failed', error.message);
    }
  }

  sendCommand(instructions, receiverUsername = null, tags = []) {
    if (!this.websocket) return;
    
    const command = {
      type: 'send_command',
      data: { instructions, tags }
    };
    
    if (receiverUsername) {
      command.data.receiver = receiverUsername;
    }
    
    this.websocket.send(JSON.stringify(command));
  }
}

export default ControlMeReactNativeClient;
```

## Custom Hardware Client (Arduino/ESP32)

```cpp
#include <WiFi.h>
#include <WebSocketsClient.h>
#include <ArduinoJson.h>

class ControlMeHardwareClient {
private:
  WebSocketsClient webSocket;
  String authToken;
  String baseUrl;
  
public:
  ControlMeHardwareClient(String url) : baseUrl(url) {}
  
  bool login(String loginName, String password) {
    // HTTP login request (simplified)
    // Set authToken on success
    return true;
  }
  
  void connectWebSocket() {
    webSocket.begin(baseUrl, 8080, "/api/ws");
    webSocket.onEvent([this](WStype_t type, uint8_t * payload, size_t length) {
      this->webSocketEvent(type, payload, length);
    });
    
    // Send auth token after connection
    webSocket.setAuthorization(("Bearer " + authToken).c_str());
  }
  
  void webSocketEvent(WStype_t type, uint8_t * payload, size_t length) {
    switch(type) {
      case WStype_TEXT:
        handleMessage((char*)payload);
        break;
    }
  }
  
  void handleMessage(String message) {
    DynamicJsonDocument doc(2048);
    deserializeJson(doc, message);
    
    if (doc["type"] == "command") {
      handleCommand(doc["data"]);
    }
  }
  
  void handleCommand(JsonObject command) {
    JsonArray instructions = command["instructions"];
    
    for (JsonObject instruction : instructions) {
      String type = instruction["type"];
      JsonObject content = instruction["content"];
      
      if (type == "custom-vibration") {
        handleVibration(content);
      } else if (type == "custom-led") {
        handleLED(content);
      } else if (type == "custom-servo") {
        handleServo(content);
      }
      
      delay(1000); // Wait between instructions
    }
  }
  
  void handleVibration(JsonObject content) {
    int intensity = content["intensity"];
    int duration = content["duration"];
    
    // Control vibration motor
    analogWrite(VIBRATION_PIN, intensity);
    delay(duration);
    analogWrite(VIBRATION_PIN, 0);
  }
  
  void handleLED(JsonObject content) {
    JsonArray pattern = content["pattern"];
    String color = content["color"];
    
    // Control RGB LED strip
    for (int i = 0; i < pattern.size(); i++) {
      setLEDColor(color);
      delay(pattern[i]);
      setLEDColor("off");
      delay(100);
    }
  }
  
  void handleServo(JsonObject content) {
    int angle = content["angle"];
    int speed = content["speed"];
    
    // Control servo motor
    moveServoToAngle(angle, speed);
  }
};
```

## Integration Best Practices

### Error Handling
- Always handle unknown instruction types gracefully
- Implement retry logic for WebSocket connections
- Log errors for debugging but don't crash on unknown commands
- Validate content structure before processing

### Security
- Store auth tokens securely (keychain/encrypted storage)
- Validate file downloads before opening
- Implement user confirmation for sensitive actions
- Rate limit outgoing commands

### User Experience
- Show sender information for all commands
- Provide options to block/report users
- Allow users to configure which instruction types they accept
- Implement offline queuing for commands

### Performance
- Process instructions asynchronously when possible
- Cache frequently accessed data
- Implement proper cleanup for media resources
- Use appropriate timeouts for network operations

### Platform-Specific Considerations

#### Web Browsers
- Handle popup blockers
- Respect user gesture requirements for media
- Use service workers for offline functionality
- Implement proper CORS handling

#### Mobile Apps
- Request appropriate permissions (camera, files, notifications)
- Handle app backgrounding gracefully
- Implement push notifications for offline commands
- Optimize for battery usage

#### Desktop Applications
- Handle multiple monitor setups
- Respect system notification settings
- Implement proper file system permissions
- Consider accessibility requirements

#### Hardware/IoT Devices
- Implement proper error recovery
- Handle network disconnections
- Use appropriate power management
- Implement safety limits for physical controls
