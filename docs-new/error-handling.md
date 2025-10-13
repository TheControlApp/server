# Error Handling Guide

This guide explains how errors are handled in the ControlApp API and provides comprehensive examples for implementing proper error handling in your applications.

## RFC 7807 Compliance

The ControlApp API follows the [RFC 7807 Problem Details for HTTP APIs](https://tools.ietf.org/html/rfc7807) standard. This provides a machine-readable format for specifying errors in HTTP API responses, ensuring consistency and interoperability.

## Error Response Structure

All error responses follow this consistent structure:

```json
{
  "type": "string",           // Error type identifier (required)
  "title": "string",          // Human-readable error title (required)
  "status": 400,              // HTTP status code (required)
  "detail": "string",         // Specific error description (required)
  "instance": "string",       // Request path (optional)
  "action": "string",         // Suggested action for developer (optional)
  "help": "string"            // Additional guidance (optional)
}
```

### Field Descriptions

- **type**: A URI reference that identifies the problem type. When dereferenced, it should provide human-readable documentation.
- **title**: A short, human-readable summary of the problem type.
- **status**: The HTTP status code for this occurrence of the problem.
- **detail**: A human-readable explanation specific to this occurrence of the problem.
- **instance**: A URI reference that identifies the specific occurrence of the problem (typically the request path).
- **action**: Suggested action the client can take to resolve the issue (optional).
- **help**: Additional help text or reference to documentation (optional).

## Common Error Types

### 400 Bad Request
Used when the request is malformed or contains invalid data.

```json
{
  "type": "bad_request",
  "title": "Bad Request",
  "status": 400,
  "detail": "Request body is not valid JSON",
  "instance": "/api/v1/auth/login",
  "action": "Ensure request body contains valid JSON"
}
```

**Common causes:**
- Invalid JSON syntax
- Missing Content-Type header
- Malformed request body

### 401 Unauthorized
Used when authentication is required or has failed.

```json
{
  "type": "unauthorized", 
  "title": "Unauthorized Access",
  "status": 401,
  "detail": "Invalid username or password",
  "instance": "/api/v1/auth/login",
  "action": "Please check your credentials and try again"
}
```

**Common causes:**
- Invalid credentials
- Missing Authorization header
- Expired JWT token

### 403 Forbidden
Used when the authenticated user doesn't have permission to access the resource.

```json
{
  "type": "forbidden",
  "title": "Access Forbidden", 
  "status": 403,
  "detail": "You do not have permission to access this resource",
  "instance": "/api/v1/admin/users",
  "action": "Contact your administrator for access"
}
```

**Common causes:**
- Insufficient user permissions
- Trying to access admin-only resources
- User account is blocked

### 404 Not Found
Used when the requested resource doesn't exist.

```json
{
  "type": "not_found",
  "title": "Resource Not Found",
  "status": 404,
  "detail": "User with ID '550e8400-e29b-41d4-a716-446655440000' not found",
  "instance": "/api/v1/users/550e8400-e29b-41d4-a716-446655440000",
  "action": "Verify the resource ID and try again"
}
```

**Common causes:**
- Non-existent resource ID
- Deleted resource
- Typo in URL path

### 409 Conflict
Used when the request conflicts with the current state of the resource.

```json
{
  "type": "conflict",
  "title": "Resource Conflict", 
  "status": 409,
  "detail": "Username 'testuser' already exists",
  "instance": "/api/v1/auth/register",
  "action": "Please choose a different username"
}
```

**Common causes:**
- Duplicate username/email
- Resource already exists
- Concurrent modification conflicts

### 422 Validation Error
Used when input validation fails. Includes detailed field-level error information.

```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation",
  "instance": "/api/v1/auth/register",
  "help": "Please check the field requirements in the API documentation",
  "errors": [
    {
      "field": "username",
      "message": "Username must be at least 3 characters long",
      "code": "MIN_LENGTH"
    },
    {
      "field": "password",
      "message": "Password must be at least 6 characters long", 
      "code": "MIN_LENGTH"
    }
  ]
}
```

#### Validation Error Codes

Common validation error codes include:

- `REQUIRED`: Required field is missing
- `MIN_LENGTH`: Field value is too short
- `MAX_LENGTH`: Field value is too long
- `INVALID_FORMAT`: Field value doesn't match expected format
- `INVALID_EMAIL`: Email format is invalid
- `WEAK_PASSWORD`: Password doesn't meet strength requirements

### 500 Internal Server Error
Used when an unexpected server error occurs.

```json
{
  "type": "internal_server_error",
  "title": "Internal Server Error",
  "status": 500, 
  "detail": "Database connection failed",
  "instance": "/api/v1/users",
  "action": "Please try again later or contact support if the problem persists"
}
```

**Common causes:**
- Database connection issues
- Unhandled exceptions
- Configuration problems

## Client Implementation Guidelines

### Error Handling Best Practices

1. **Always check the status code** first to determine the category of error
2. **Use the type field** to programmatically handle specific error types
3. **Display the detail message** to users for specific error information
4. **Follow the action guidance** when provided to help users resolve issues
5. **Log the complete error response** for debugging purposes
6. **Implement retry logic** for temporary errors (5xx status codes)

### Example Implementation (JavaScript)

```javascript
class ControlAppAPI {
    async makeRequest(url, options = {}) {
        try {
            const response = await fetch(url, {
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers
                },
                ...options
            });
            
            return await this.handleResponse(response);
        } catch (error) {
            throw new Error(`Network error: ${error.message}`);
        }
    }
    
    async handleResponse(response) {
        if (response.ok) {
            return await response.json();
        }
        
        const error = await response.json();
        
        // Handle specific error types
        switch (error.type) {
            case 'validation_error':
                this.handleValidationError(error);
                break;
                
            case 'unauthorized':
                this.handleUnauthorizedError(error);
                break;
                
            case 'not_found':
                this.handleNotFoundError(error);
                break;
                
            case 'conflict':
                this.handleConflictError(error);
                break;
                
            case 'internal_server_error':
                this.handleServerError(error);
                break;
                
            default:
                this.handleGenericError(error);
        }
        
        throw error;
    }
    
    handleValidationError(error) {
        console.log('Validation failed:', error.detail);
        
        if (error.errors) {
            error.errors.forEach(fieldError => {
                console.log(`${fieldError.field}: ${fieldError.message}`);
                this.displayFieldError(fieldError.field, fieldError.message);
            });
        }
    }
    
    handleUnauthorizedError(error) {
        console.log('Authentication required:', error.detail);
        // Redirect to login or refresh token
        this.redirectToLogin();
    }
    
    handleNotFoundError(error) {
        console.log('Resource not found:', error.detail);
        this.showNotFoundMessage(error.detail);
    }
    
    handleConflictError(error) {
        console.log('Resource conflict:', error.detail);
        this.showConflictMessage(error.detail, error.action);
    }
    
    handleServerError(error) {
        console.error('Server error:', error.detail);
        this.showServerErrorMessage(error.action);
    }
    
    handleGenericError(error) {
        console.error('API error:', error.detail);
        this.showErrorMessage(error.detail);
    }
}
```

### Example Implementation (Python)

```python
import requests
import logging

class ControlAppAPI:
    def __init__(self, base_url, token=None):
        self.base_url = base_url
        self.token = token
        self.session = requests.Session()
        
        if token:
            self.session.headers.update({'Authorization': f'Bearer {token}'})
    
    def make_request(self, method, endpoint, **kwargs):
        url = f"{self.base_url}{endpoint}"
        
        try:
            response = self.session.request(method, url, **kwargs)
            return self.handle_response(response)
        except requests.exceptions.RequestException as e:
            raise Exception(f"Network error: {e}")
    
    def handle_response(self, response):
        if response.ok:
            return response.json()
        
        try:
            error = response.json()
        except ValueError:
            # Non-JSON error response
            raise Exception(f"HTTP {response.status_code}: {response.text}")
        
        error_type = error.get('type', 'unknown_error')
        
        # Handle specific error types
        if error_type == 'validation_error':
            self.handle_validation_error(error)
        elif error_type == 'unauthorized':
            self.handle_unauthorized_error(error)
        elif error_type == 'not_found':
            self.handle_not_found_error(error)
        elif error_type == 'conflict':
            self.handle_conflict_error(error)
        elif error_type == 'internal_server_error':
            self.handle_server_error(error)
        else:
            self.handle_generic_error(error)
        
        raise Exception(error['detail'])
    
    def handle_validation_error(self, error):
        logging.warning(f"Validation failed: {error['detail']}")
        
        if 'errors' in error:
            for field_error in error['errors']:
                logging.warning(f"{field_error['field']}: {field_error['message']}")
    
    def handle_unauthorized_error(self, error):
        logging.warning(f"Authentication failed: {error['detail']}")
        # Clear token and redirect to login
        self.token = None
        self.session.headers.pop('Authorization', None)
    
    def handle_not_found_error(self, error):
        logging.info(f"Resource not found: {error['detail']}")
    
    def handle_conflict_error(self, error):
        logging.warning(f"Resource conflict: {error['detail']}")
    
    def handle_server_error(self, error):
        logging.error(f"Server error: {error['detail']}")
    
    def handle_generic_error(self, error):
        logging.error(f"API error: {error['detail']}")
```

### WebSocket Error Handling

WebSocket connections can also receive error messages. These follow a similar structure:

```json
{
    "type": "error",
    "data": {
        "type": "authentication_failed",
        "title": "Authentication Failed",
        "status": 401,
        "detail": "Invalid token provided",
        "action": "Please login again and reconnect"
    }
}
```

Example WebSocket error handling:

```javascript
ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    
    if (message.type === 'error') {
        const error = message.data;
        console.error('WebSocket error:', error.detail);
        
        if (error.type === 'authentication_failed') {
            // Reconnect with new token
            this.reconnectWithNewToken();
        }
    }
};
```

## Testing Error Responses

You can test various error responses using curl:

```bash
# Test validation error
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "ab", "password": "123"}'

# Test unauthorized error  
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login_name": "nonexistent", "password": "wrong"}'

# Test not found error
curl -X GET http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000

# Test bad request error
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "invalid json"

# Test forbidden error (without auth)
curl -X GET http://localhost:8080/api/v1/admin/users
```

## Debugging Tips

1. **Check the instance field** to identify which endpoint caused the error
2. **Look for action fields** that provide specific guidance
3. **Use the type field** for programmatic error handling
4. **Log complete error responses** for debugging
5. **Check HTTP status codes** first to categorize errors
6. **Implement exponential backoff** for retry logic on 5xx errors

## See Also

- [REST API Reference](rest-api.md) - Complete API endpoint documentation
- [Authentication Guide](authentication.md) - JWT token handling
- [WebSocket API Reference](websocket-api.md) - Real-time communication
- [RFC 7807](https://tools.ietf.org/html/rfc7807) - Problem Details for HTTP APIs standard