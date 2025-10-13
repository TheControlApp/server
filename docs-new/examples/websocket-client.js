/**
 * ControlMe WebSocket Client Example (Updated)
 * 
 * This example demonstrates how to connect to the ControlMe server
 * and handle the various message formats using the current API.
 * 
 * Updated to work with the actual server implementation.
 */

class ControlMeClient {
    constructor(serverUrl, token = null) {
        this.serverUrl = serverUrl;
        this.token = token;
        this.ws = null;
        this.userId = null;
        this.isConnected = false;
        this.isAuthenticated = false;
        
        // Event handlers
        this.onCommand = null;
        this.onError = null;
        this.onStatusUpdate = null;
        this.onConnectionStatusChange = null;
    }

    /**
     * Connect to the WebSocket server
     * Can connect with or without authentication
     */
    connect() {
        let wsUrl = `${this.serverUrl}/ws/client`;
        
        // Add token as query parameter if provided
        if (this.token) {
            wsUrl += `?token=${encodeURIComponent(this.token)}`;
        }
        
        try {
            this.ws = new WebSocket(wsUrl);
            
            this.ws.onopen = (event) => {
                console.log('Connected to ControlMe server');
                this.isConnected = true;
                
                // If we didn't authenticate via query param, try message-based auth
                if (this.token && !this.token.includes('anonymous')) {
                    this.authenticateWithToken();
                }
                
                // Start heartbeat
                this.startHeartbeat();
                
                if (this.onConnectionStatusChange) {
                    this.onConnectionStatusChange('connected');
                }
            };
            
            this.ws.onmessage = (event) => {
                try {
                    const message = JSON.parse(event.data);
                    this.handleMessage(message);
                } catch (error) {
                    console.error('Failed to parse message:', error, event.data);
                }
            };
            
            this.ws.onclose = (event) => {
                console.log('Disconnected from server:', event.code, event.reason);
                this.isConnected = false;
                this.isAuthenticated = false;
                
                if (this.onConnectionStatusChange) {
                    this.onConnectionStatusChange('disconnected');
                }
                
                // Stop heartbeat
                this.stopHeartbeat();
            };
            
            this.ws.onerror = (error) => {
                console.error('WebSocket error:', error);
                if (this.onError) {
                    this.onError({
                        type: 'websocket_error',
                        detail: 'Connection error occurred',
                        instance: this.serverUrl
                    });
                }
            };
            
        } catch (error) {
            console.error('Failed to connect:', error);
            if (this.onError) {
                this.onError({
                    type: 'connection_failed',
                    detail: error.message,
                    instance: this.serverUrl
                });
            }
        }
    }

    /**
     * Authenticate using JWT token via message
     */
    authenticateWithToken() {
        if (!this.token) {
            console.warn('No token provided for authentication');
            return;
        }
        
        const authMessage = {
            type: 'auth',
            token: this.token
        };
        
        this.sendMessage(authMessage);
        console.log('Sent authentication message');
    }

    /**
     * Start heartbeat to keep connection alive
     */
    startHeartbeat() {
        this.heartbeatInterval = setInterval(() => {
            if (this.isConnected) {
                this.sendMessage({ type: 'ping' });
            }
        }, 30000); // Ping every 30 seconds
    }

    /**
     * Stop heartbeat
     */
    stopHeartbeat() {
        if (this.heartbeatInterval) {
            clearInterval(this.heartbeatInterval);
            this.heartbeatInterval = null;
        }
    }

    /**
     * Handle incoming messages from server
     */
    handleMessage(message) {
        console.log('Received message:', message);

        // Handle different message types based on actual server implementation
        switch (message.type) {
            case 'pong':
                // Heartbeat response
                console.log('Received pong');
                break;
                
            case 'auth_success':
                this.handleAuthSuccess(message);
                break;
                
            case 'auth_error':
                this.handleAuthError(message);
                break;
                
            case 'command':
                this.handleCommandMessage(message);
                break;
                
            case 'command_status':
                this.handleCommandStatus(message);
                break;
                
            case 'error':
                this.handleError(message);
                break;
                
            case 'broadcast':
                this.handleBroadcast(message);
                break;
                
            default:
                console.warn('Unknown message type:', message.type, message);
        }
    }

    /**
     * Handle authentication success
     */
    handleAuthSuccess(message) {
        this.isAuthenticated = true;
        if (message.user_id) {
            this.userId = message.user_id;
        }
        console.log('Authentication successful', message);
        
        if (this.onConnectionStatusChange) {
            this.onConnectionStatusChange('authenticated');
        }
    }

    /**
     * Handle authentication error
     */
    handleAuthError(message) {
        this.isAuthenticated = false;
        console.error('Authentication failed:', message);
        
        if (this.onError) {
            this.onError({
                type: 'authentication_failed',
                detail: message.error || 'Authentication failed',
                instance: '/ws/client'
            });
        }
    }

    /**
     * Handle command message from server
     */
    handleCommandMessage(message) {
        console.log('Received command:', message);
        
        if (this.onCommand) {
            this.onCommand(message);
        }
        
        // Auto-acknowledge command receipt
        this.sendCommandResponse(message.id, 'received', {
            timestamp: new Date().toISOString()
        });
    }

    /**
     * Handle command status updates
     */
    handleCommandStatus(message) {
        console.log('Command status update:', message);
        
        if (this.onStatusUpdate) {
            this.onStatusUpdate(message);
        }
    }

    /**
     * Handle broadcast messages
     */
    handleBroadcast(message) {
        console.log('Received broadcast:', message);
        
        // You can add custom handling for broadcasts here
    }

    /**
     * Handle error messages
     */
    handleError(message) {
        console.error('Server error:', message);
        
        if (this.onError) {
            // Convert to RFC 7807 format if not already
            const error = {
                type: message.type || 'server_error',
                title: message.title || 'Server Error',
                status: message.status || 500,
                detail: message.detail || message.error || 'An error occurred',
                instance: message.instance || '/ws/client'
            };
            
            this.onError(error);
        }
    }

    /**
     * Send command response to server
     */
    sendCommandResponse(commandId, status, data = {}) {
        const response = {
            type: 'command_response',
            command_id: commandId,
            status: status,
            data: data,
            timestamp: new Date().toISOString()
        };
        
        this.sendMessage(response);
        console.log('Sent command response:', response);
    }

    /**
     * Send a command to another client (if authenticated)
     */
    sendCommand(targetUserId, commandData) {
        if (!this.isAuthenticated) {
            console.error('Must be authenticated to send commands');
            return;
        }
        
        const command = {
            type: 'command',
            target_user_id: targetUserId,
            data: commandData,
            timestamp: new Date().toISOString()
        };
        
        this.sendMessage(command);
        console.log('Sent command:', command);
    }

    /**
     * Send message to server
     */
    sendMessage(message) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const messageStr = JSON.stringify(message);
            this.ws.send(messageStr);
            console.log('Sent message:', message);
        } else {
            console.error('WebSocket not connected');
            if (this.onError) {
                this.onError({
                    type: 'websocket_not_connected',
                    detail: 'Cannot send message - WebSocket not connected',
                    instance: '/ws/client'
                });
            }
        }
    }

    /**
     * Get connection status
     */
    getStatus() {
        return {
            connected: this.isConnected,
            authenticated: this.isAuthenticated,
            userId: this.userId,
            readyState: this.ws ? this.ws.readyState : WebSocket.CLOSED
        };
    }

    /**
     * Disconnect from server
     */
    disconnect() {
        this.stopHeartbeat();
        
        if (this.ws) {
            this.ws.close(1000, 'Client disconnect');
        }
        
        this.isConnected = false;
        this.isAuthenticated = false;
        this.userId = null;
    }
}

/**
 * Example usage for authenticated client:
 */
/*
const authenticatedClient = new ControlMeClient('ws://localhost:8080', 'your-jwt-token-here');

authenticatedClient.onError = (error) => {
    console.error('Client error:', error);
};

authenticatedClient.onCommand = (command) => {
    console.log('Received command:', command);
    // Handle the command and send response
    authenticatedClient.sendCommandResponse(command.id, 'completed', {
        result: 'Command executed successfully'
    });
};

authenticatedClient.onConnectionStatusChange = (status) => {
    console.log('Connection status:', status);
};

authenticatedClient.connect();
*/

/**
 * Example usage for anonymous client:
 */
/*
const anonymousClient = new ControlMeClient('ws://localhost:8080');

anonymousClient.onError = (error) => {
    console.error('Anonymous client error:', error);
};

anonymousClient.onConnectionStatusChange = (status) => {
    console.log('Anonymous connection status:', status);
};

anonymousClient.connect();
*/

/**
 * Example with error handling and reconnection:
 */
/*
class ReconnectingControlMeClient extends ControlMeClient {
    constructor(serverUrl, token, maxReconnectAttempts = 5) {
        super(serverUrl, token);
        this.maxReconnectAttempts = maxReconnectAttempts;
        this.reconnectAttempts = 0;
        this.reconnectDelay = 1000; // Start with 1 second
    }
    
    connect() {
        super.connect();
        
        // Override onclose to add reconnection logic
        const originalOnClose = this.ws.onclose;
        this.ws.onclose = (event) => {
            originalOnClose.call(this, event);
            
            if (this.reconnectAttempts < this.maxReconnectAttempts) {
                this.reconnectAttempts++;
                console.log(`Reconnecting in ${this.reconnectDelay}ms (attempt ${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
                
                setTimeout(() => {
                    this.connect();
                }, this.reconnectDelay);
                
                // Exponential backoff
                this.reconnectDelay = Math.min(this.reconnectDelay * 2, 30000);
            } else {
                console.error('Max reconnection attempts reached');
            }
        };
        
        // Reset reconnection counter on successful connection
        const originalOnOpen = this.ws.onopen;
        this.ws.onopen = (event) => {
            originalOnOpen.call(this, event);
            this.reconnectAttempts = 0;
            this.reconnectDelay = 1000;
        };
    }
}
*/

// Export for use in Node.js environments
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { ControlMeClient };
}