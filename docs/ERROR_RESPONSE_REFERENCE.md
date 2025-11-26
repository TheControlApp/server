# API Error Response Reference

This document provides comprehensive information about the error response format used throughout the ControlApp API.

## RFC 7807 Compliance

The API follows the [RFC 7807 Problem Details for HTTP APIs](https://tools.ietf.org/html/rfc7807) standard. This provides a machine-readable format for specifying errors in HTTP API responses.

## Error Response Structure

All error responses follow this consistent structure:

```json
{
  "type": "string",           // Error type identifier (required)
  "title": "string",          // Human-readable error title (required)
  "status": 400,              // HTTP status code (required)
  "detail": "string",         // Specific error description (required)
  "action": "string",         // Suggested action for developer (optional)
  "help": "string"            // Additional guidance (optional)
}
```

### Field Descriptions

- **type**: A URI reference that identifies the problem type. When dereferenced, it should provide human-readable documentation.
- **title**: A short, human-readable summary of the problem type.
- **status**: The HTTP status code for this occurrence of the problem.
- **detail**: A human-readable explanation specific to this occurrence of the problem.
- **action**: Suggested action the client can take to resolve the issue (optional).
- **help**: Additional help text or reference to documentation (optional).

## Error Types by Status Code

### 400 Bad Request
Used when the request is malformed or contains invalid data.

```json
{
  "type": "bad_request",
  "title": "Bad Request",
  "status": 400,
  "detail": "Request body is not valid JSON or missing required fields"
}
```

### 401 Unauthorized
Used when authentication is required or has failed.

```json
{
  "type": "unauthorized", 
  "title": "Unauthorized Access",
  "status": 401,
  "detail": "Invalid username or password",
  "action": "Please check your credentials and try again"
}
```

### 403 Forbidden
Used when the authenticated user doesn't have permission to access the resource.

```json
{
  "type": "forbidden",
  "title": "Access Forbidden", 
  "status": 403,
  "detail": "You do not have permission to access this resource",
  "action": "Contact your administrator for access"
}
```

### 404 Not Found
Used when the requested resource doesn't exist.

```json
{
  "type": "not_found",
  "title": "Resource Not Found",
  "status": 404,
  "detail": "User not found"
}
```

### 409 Conflict
Used when the request conflicts with the current state of the resource.

```json
{
  "type": "conflict",
  "title": "Resource Conflict", 
  "status": 409,
  "detail": "Username already exists",
  "action": "Please choose a different username",
  "instance": {
    "username": "testuser"
  }
}
```

### 422 Validation Error
Used when input validation fails. Includes detailed field-level error information.

```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation",
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

- `MIN_LENGTH`: Field value is too short
- `MAX_LENGTH`: Field value is too long
- `INVALID_FORMAT`: Field value doesn't match expected format
- `REQUIRED`: Required field is missing

### 500 Internal Server Error
Used when an unexpected server error occurs.

```json
{
  "type": "internal_server_error",
  "title": "Internal Server Error",
  "status": 500, 
  "detail": "Database connection failed",
  "action": "Please try again later or contact support if the problem persists"
}
```

## Client Implementation Guidelines

### Error Handling Best Practices

1. **Always check the status code** first to determine the category of error
2. **Use the type field** to programmatically handle specific error types
3. **Display the detail message** to users for specific error information
4. **Follow the action guidance** when provided to help users resolve issues
5. **Log the complete error response** for debugging purposes

### Example Client Code (JavaScript)

```javascript
async function handleApiResponse(response) {
  if (!response.ok) {
    const error = await response.json();
    
    switch (error.type) {
      case 'validation_error':
        // Handle validation errors
        displayFieldErrors(error.errors);
        break;
        
      case 'unauthorized':
        // Redirect to login
        redirectToLogin();
        break;
        
      case 'not_found':
        // Show not found message
        showNotFoundMessage(error.detail);
        break;
        
      default:
        // Generic error handling
        showErrorMessage(error.detail);
    }
    
    return null;
  }
  
  return await response.json();
}
```

### Example Client Code (Python)

```python
import requests

def handle_api_response(response):
    if not response.ok:
        error = response.json()
        
        if error['type'] == 'validation_error':
            for field_error in error.get('errors', []):
                print(f"Error in {field_error['field']}: {field_error['message']}")
        
        elif error['type'] == 'unauthorized':
            print("Authentication failed. Please login again.")
            
        else:
            print(f"API Error: {error['detail']}")
            
        return None
    
    return response.json()
```

## Testing Error Responses

You can test the various error responses using the provided test JSON files:

```bash
# Test validation error
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d @docs/examples/api_tests/test_short_username.json

# Test unauthorized error  
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d @docs/examples/api_tests/test_login_nonexistent.json

# Test not found error
curl -X GET http://localhost:8080/api/v1/users/550e8400-e29b-41d4-a716-446655440000

# Test bad request error
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d "{}"
```

## Swagger Documentation

For interactive API documentation and testing, visit:
http://localhost:8080/swagger/index.html

The Swagger documentation includes complete examples of all error responses and allows you to test the API directly from your browser.