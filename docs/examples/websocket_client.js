/**
 * ControlMe WebSocket Client Example
 * 
 * This example demonstrates how to connect to the ControlMe server
 * and handle the various message formats.
 */

class ControlMeClient {
    constructor(serverUrl, token) {
        this.serverUrl = serverUrl;
        this.token = token;
        this.ws = null;
        this.userId = null;
        this.isConnected = false;
        
        // Event handlers
        this.onCommand = null;
        this.onError = null;
        this.onStatusUpdate = null;
    }

    /**
     * Connect to the WebSocket server
     */
    connect() {
        const wsUrl = `${this.serverUrl}/ws/client?token=${this.token}`;
        
        try {
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = (event) => {
                console.log('Connected to ControlMe server');
                this.isConnected = true;
            };
            
            this.ws.onmessage = (event) => {
                this.handleMessage(JSON.parse(event.data));
            };
            
            this.ws.onclose = (event) => {
                console.log('Disconnected from server:', event.code, event.reason);
                this.isConnected = false;
            };
            
            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                if (this.onError) {
                    this.onError(error);
                }
            };
            
        } catch (error) {
            console.error('Failed to connect:', error);
        }
    }

    /**
     * Handle incoming messages from server
     */
    handleMessage(message) {
        console.log('Received message:', message);

        switch (message.type) {
            case 'connection_status':
                this.handleConnectionStatus(message);
                break;
                
            case 'command_assignment':
                this.handleCommandAssignment(message);
                break;
                
            case 'command_status':
                this.handleCommandStatus(message);
                break;
                
            case 'error':
                this.handleError(message);
                break;
                
            default:
                console.warn('Unknown message type:', message.type);
        }
    }

    /**
     * Handle connection status messages
     */
    handleConnectionStatus(message) {
        if (message.data.status === 'authenticated') {
            this.userId = message.data.user_id;
            console.log('Authenticated as user:', this.userId);
        }
    }

    /**
     * Handle command assignment from server
     */
    handleCommandAssignment(message) {
        const command = message.data.command;
        console.log('Received command:', command.id);
        
        try {
            // Access the instructions array directly
            const instructions = command.instructions;
            console.log(`Command has ${instructions.length} instruction(s)`);
            
            // Execute instructions sequentially
            this.executeInstructionsSequentially(instructions, command, 0);
            
        } catch (error) {
            console.error('Failed to process command instructions:', error);
            this.sendCommandCompletion(command.id, 'failed', {
                error: 'Failed to process instructions'
            });
        }
    }

    /**
     * Execute instructions in sequential order
     */
    async executeInstructionsSequentially(instructions, command, index) {
        if (index >= instructions.length) {
            // All instructions completed
            console.log('All instructions completed for command:', command.id);
            this.sendCommandCompletion(command.id, 'completed', {
                instructions_completed: instructions.length,
                execution_time: Date.now() - this.commandStartTime
            });
            return;
        }

        const instruction = instructions[index];
        console.log(`Executing instruction ${index + 1}/${instructions.length}: ${instruction.type}`);
        
        try {
            await this.processInstructionAsync(instruction, command, index);
            // Move to next instruction
            setTimeout(() => {
                this.executeInstructionsSequentially(instructions, command, index + 1);
            }, 100); // Small delay between instructions
            
        } catch (error) {
            console.error(`Failed to execute instruction ${index}:`, error);
            this.sendCommandCompletion(command.id, 'failed', {
                error: `Failed at instruction ${index}: ${error.message}`,
                instructions_completed: index,
                failed_instruction: instruction.type
            });
        }
    }

    /**
     * Process individual instruction (async version for sequential execution)
     */
    async processInstructionAsync(instruction, command, index) {
        // Store command start time for the first instruction
        if (index === 0) {
            this.commandStartTime = Date.now();
        }

        return new Promise((resolve, reject) => {
            switch (instruction.type) {
                case 'std_popup':
                    this.handlePopupInstructionAsync(instruction, command)
                        .then(resolve)
                        .catch(reject);
                    break;
                    
                case 'std_notification':
                    this.handleNotificationInstructionAsync(instruction, command)
                        .then(resolve)
                        .catch(reject);
                    break;
                    
                case 'std_timer':
                    this.handleTimerInstructionAsync(instruction, command)
                        .then(resolve)
                        .catch(reject);
                    break;
                    
                case 'std_input':
                    this.handleInputInstructionAsync(instruction, command)
                        .then(resolve)
                        .catch(reject);
                    break;
                    
                default:
                    reject(new Error(`Unsupported instruction type: ${instruction.type}`));
            }
        });
    }

    /**
     * Handle popup instruction (async version)
     */
    async handlePopupInstructionAsync(instruction, command) {
        const content = instruction.content;
        
        return new Promise((resolve, reject) => {
            // Show popup dialog
            const result = confirm(`${content.title}\n\n${content.body}\n\nClick OK to acknowledge.`);
            
            if (result) {
                console.log('Popup acknowledged:', content.title);
                resolve({
                    button_clicked: content.button || 'OK'
                });
            } else {
                reject(new Error('User cancelled popup'));
            }
        });
    }

    /**
     * Handle notification instruction (async version)
     */
    async handleNotificationInstructionAsync(instruction, command) {
        const content = instruction.content;
        
        return new Promise((resolve) => {
            // Show browser notification if supported
            if ('Notification' in window) {
                if (Notification.permission === 'granted') {
                    new Notification(content.title, {
                        body: content.body,
                        icon: '/favicon.ico'
                    });
                    console.log('Notification shown:', content.title);
                    resolve({ notification_shown: true });
                } else if (Notification.permission !== 'denied') {
                    Notification.requestPermission().then(permission => {
                        if (permission === 'granted') {
                            new Notification(content.title, {
                                body: content.body,
                                icon: '/favicon.ico'
                            });
                        }
                        console.log('Notification permission:', permission);
                        resolve({ notification_shown: permission === 'granted' });
                    });
                } else {
                    console.log('Notifications denied');
                    resolve({ notification_shown: false });
                }
            } else {
                console.log('Notifications not supported');
                resolve({ notification_shown: false });
            }
        });
    }

    /**
     * Handle timer instruction (async version)
     */
    async handleTimerInstructionAsync(instruction, command) {
        const content = instruction.content;
        const duration = content.duration; // seconds
        
        console.log(`Starting timer: ${content.title} for ${duration} seconds`);
        
        return new Promise((resolve) => {
            // Start countdown
            let remaining = duration;
            const timer = setInterval(() => {
                remaining--;
                if (remaining % 60 === 0 || remaining <= 10) {
                    console.log(`Timer: ${remaining} seconds remaining`);
                }
                
                if (remaining <= 0) {
                    clearInterval(timer);
                    console.log('Timer completed!');
                    resolve({
                        duration_completed: duration,
                        completed_method: 'auto'
                    });
                }
            }, 1000);
        });
    }

    /**
     * Handle input form instruction (async version)
     */
    async handleInputInstructionAsync(instruction, command) {
        const content = instruction.content;
        const formData = {};
        
        console.log(`Input form: ${content.title}`);
        console.log(content.description);
        
        return new Promise((resolve, reject) => {
            // Simple prompt-based input for demo
            try {
                content.fields.forEach(field => {
                    let value;
                    
                    if (field.type === 'select') {
                        const options = field.options.join(', ');
                        value = prompt(`${field.label} (Options: ${options})`);
                    } else if (field.type === 'number') {
                        value = prompt(field.label);
                        if (value) value = parseFloat(value);
                    } else {
                        value = prompt(field.label);
                    }
                    
                    if (field.required && !value) {
                        throw new Error(`Required field '${field.name}' not provided`);
                    }
                    
                    formData[field.name] = value;
                });
                
                console.log('Form completed:', formData);
                resolve({
                    form_data: formData
                });
                
            } catch (error) {
                reject(error);
            }
        });
    }

    /**
     * Handle command status updates
     */
    handleCommandStatus(message) {
        console.log('Command status update:', message.data);
        
        if (this.onStatusUpdate) {
            this.onStatusUpdate(message.data);
        }
    }

    /**
     * Handle error messages
     */
    handleError(message) {
        console.error('Server error:', message.data);
        
        if (this.onError) {
            this.onError(message.data);
        }
    }

    /**
     * Send command completion to server
     */
    sendCommandCompletion(commandId, status, responseData = {}) {
        const message = {
            type: 'command_completion',
            id: this.generateId(),
            timestamp: new Date().toISOString(),
            from: this.userId,
            to: null, // Server will route appropriately
            data: {
                command_id: commandId,
                status: status,
                response: {
                    success: status === 'completed',
                    timestamp: new Date().toISOString(),
                    data: responseData
                }
            }
        };
        
        this.sendMessage(message);
    }

    /**
     * Send message to server
     */
    sendMessage(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.ws.send(JSON.stringify(message));
            console.log('Sent message:', message);
        } else {
            console.error('WebSocket not connected');
        }
    }

    /**
     * Generate unique message ID
     */
    generateId() {
        return 'client-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
    }

    /**
     * Disconnect from server
     */
    disconnect() {
        if (this.ws) {
            this.ws.close();
        }
    }
}

// Example usage:
/*
const client = new ControlMeClient('ws://localhost:8080', 'your-jwt-token-here');

client.onError = (error) => {
    console.error('Client error:', error);
};

client.onStatusUpdate = (status) => {
    console.log('Command status updated:', status);
};

client.connect();
*/