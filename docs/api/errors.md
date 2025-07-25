# Error Handling

## Error Format

**REST API:**
```json
{
  "error": "error_code",
  "message": "Description",
  "field": "field_name"
}
```

**WebSocket:**
```json
{
  "type": "error",
  "error": "error_code",
  "message": "Description"
}
```

## Common Error Codes

### Authentication
- `authentication_failed` - Invalid credentials or token
- `token_expired` - JWT token expired
- `access_denied` - Insufficient permissions

### Validation  
- `invalid_request` - Malformed request data
- `missing_field` - Required field missing
- `invalid_field` - Field format invalid

### User Management
- `user_not_found` - Username doesn't exist
- `user_exists` - Username/email already taken
- `command_blocked` - Content filtered by preferences

### Rate Limiting
- `rate_limit_exceeded` - Too many requests
- `daily_limit_exceeded` - Daily quota reached

### Files
- `file_not_found` - File hash invalid
- `file_too_large` - Exceeds size limit
- `invalid_file_type` - File type not allowed

## HTTP Status Codes
- **200** - Success
- **400** - Bad request
- **401** - Authentication required
- **403** - Access denied
- **404** - Not found
- **429** - Rate limited
- **500** - Server error
  for (let attempt = 1; attempt <= maxRetries; attempt++) {
    try {
      return await apiCall(url, options);
    } catch (error) {
      if (error.code === 'rate_limit_exceeded') {
        const retryAfter = error.details.retry_after || 60;
        await sleep(retryAfter * 1000);
        continue;
      }
      
      if (attempt === maxRetries || !isRetryableError(error.code)) {
        throw error;
      }
      
      // Exponential backoff
      await sleep(Math.pow(2, attempt) * 1000);
    }
  }
}

function isRetryableError(code) {
  return ['network_error', 'server_error', 'timeout'].includes(code);
}
```

### Circuit Breaker Pattern

```javascript
class CircuitBreaker {
  constructor(threshold = 5, timeout = 60000) {
    this.threshold = threshold;
    this.timeout = timeout;
    this.failureCount = 0;
    this.lastFailureTime = null;
    this.state = 'CLOSED'; // CLOSED, OPEN, HALF_OPEN
  }
  
  async execute(fn) {
    if (this.state === 'OPEN') {
      if (Date.now() - this.lastFailureTime > this.timeout) {
        this.state = 'HALF_OPEN';
      } else {
        throw new Error('Circuit breaker is OPEN');
      }
    }
    
    try {
      const result = await fn();
      this.onSuccess();
      return result;
    } catch (error) {
      this.onFailure();
      throw error;
    }
  }
  
  onSuccess() {
    this.failureCount = 0;
    this.state = 'CLOSED';
  }
  
  onFailure() {
    this.failureCount++;
    this.lastFailureTime = Date.now();
    
    if (this.failureCount >= this.threshold) {
      this.state = 'OPEN';
    }
  }
}
```

## Debugging Tips

### Enable Debug Logging

```javascript
// Client-side debug logging
const DEBUG = localStorage.getItem('debug') === 'true';

function debugLog(category, message, data = null) {
  if (DEBUG) {
    console.log(`[${category}] ${message}`, data);
  }
}

// Usage
debugLog('API', 'Making request', { url, method });
debugLog('WebSocket', 'Received message', message);
debugLog('Error', 'API call failed', error);
```

### Error Reporting

```javascript
function reportError(error, context = {}) {
  const errorReport = {
    timestamp: new Date().toISOString(),
    error: {
      code: error.code,
      message: error.message,
      status: error.status,
      stack: error.stack
    },
    context: {
      url: window.location.href,
      userAgent: navigator.userAgent,
      ...context
    }
  };
  
  // Send to error tracking service
  console.error('Error Report:', errorReport);
}
```

This comprehensive error handling guide ensures consistent error management across all client implementations and provides clear recovery strategies for common failure scenarios.
