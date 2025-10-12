# API Test Examples

This directory contains JSON test files for testing the ControlMe API endpoints. These files demonstrate various scenarios including successful requests, validation errors, and authentication failures.

## Authentication Endpoints

### Registration Tests (`/api/v1/auth/register`)

#### Successful Registration
- **File**: `test_valid_user.json`
- **Description**: Valid user registration with all required fields
- **Expected Response**: HTTP 201 with user data

#### Validation Errors
- **File**: `test_short_username.json`
- **Description**: Username too short (less than 3 characters)
- **Expected Response**: HTTP 422 with validation error details

- **File**: `test_weak_password.json`
- **Description**: Password too weak (less than 6 characters)
- **Expected Response**: HTTP 422 with validation error details

#### Conflict Error
- **Description**: Attempting to register with an existing username
- **Expected Response**: HTTP 409 with conflict error details
- **Usage**: Register `test_valid_user.json` first, then try again

### Login Tests (`/api/v1/auth/login`)

#### Successful Login
- **File**: `test_login_valid.json`
- **Description**: Valid login credentials
- **Expected Response**: HTTP 200 with user data and JWT token

#### Authentication Errors
- **File**: `test_login_nonexistent.json`
- **Description**: Login attempt with non-existent username
- **Expected Response**: HTTP 401 with unauthorized error details

- **File**: `test_login_wrong_password.json`
- **Description**: Login attempt with incorrect password
- **Expected Response**: HTTP 401 with unauthorized error details

## Error Response Structure

The API follows RFC 7807 Problem Details standard for consistent error reporting:

### Bad Request Errors (HTTP 400)
```json
{
  "type": "bad_request",
  "title": "Bad Request", 
  "status": 400,
  "detail": "Request body is not valid JSON or missing required fields"
}
```

### Validation Errors (HTTP 422)
```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "One or more fields failed validation",
  "errors": [
    {
      "field": "username",
      "message": "Username must be at least 3 characters long",
      "code": "MIN_LENGTH"
    }
  ],
  "help": "Please check the field requirements in the API documentation"
}
```

### Unauthorized Errors (HTTP 401)
```json
{
  "type": "unauthorized",
  "title": "Authentication Failed",
  "status": 401,
  "detail": "Invalid username or password",
  "action": "Please check your credentials and try again"
}
```

### Conflict Errors (HTTP 409)
```json
{
  "type": "conflict",
  "title": "Resource Conflict",
  "status": 409,
  "detail": "Username already exists",
  "instance": {
    "username": "testuser123"
  },
  "action": "Please choose a different username and try again"
}
```

### Server Errors (HTTP 500)
```json
{
  "error": "Internal server error during user creation"
}
```

## Testing Commands

Use these curl commands to test the API endpoints:

```bash
# Test registration with valid data
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d @test_valid_user.json

# Test registration with short username
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d @test_short_username.json

# Test login with valid credentials
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d @test_login_valid.json

# Test login with wrong password
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d @test_login_wrong_password.json
```

## HTTP Status Codes Reference

| Status Code | Type | Usage |
|-------------|------|-------|
| `200 OK` | Success | Successful login, returns user data and JWT token |
| `201 Created` | Success | Successful registration, returns user data |
| `400 Bad Request` | Client Error | Malformed JSON or missing required structure |
| `401 Unauthorized` | Client Error | Authentication failed (wrong credentials) |
| `409 Conflict` | Client Error | Resource conflict (username already exists) |
| `422 Unprocessable Entity` | Client Error | Valid JSON but semantic validation errors |
| `500 Internal Server Error` | Server Error | Unexpected server-side errors |

## Validation Error Codes

| Field | Code | Description | Developer Tip |
|-------|------|-------------|---------------|
| `username` | `MIN_LENGTH` | Username must be at least 3 characters | Use a longer username |
| `username` | `MAX_LENGTH` | Username must be no more than 50 characters | Use a shorter username |
| `username` | `INVALID_FORMAT` | Username has invalid format (spaces, etc.) | Remove spaces and special characters |
| `password` | `MIN_LENGTH` | Password must be at least 6 characters | Use a longer password |
| `password` | `MAX_LENGTH` | Password must be no more than 128 characters | Use a shorter password |

## Notes for Client Developers

1. **Follow RFC 7807 Problem Details** - Check `type` and `status` fields for error handling
2. **Use proper HTTP status codes** - 422 for validation, 409 for conflicts, 401 for auth failures
3. **Handle validation errors** - Check `errors` array for field-specific issues with `code` properties
4. **Malformed JSON** returns 400 Bad Request, validation issues return 422 Unprocessable Entity
5. **All endpoints require `Content-Type: application/json`** header
6. **Successful registration** returns 201 Created with user data (no password hash)
7. **Successful login** returns 200 OK with user data and JWT token
8. **Authentication errors** are always 401 Unauthorized (don't leak user existence info)