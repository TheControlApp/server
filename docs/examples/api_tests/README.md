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
- **Expected Response**: HTTP 400 with `INVALID_USERNAME` error code

- **File**: `test_weak_password.json`
- **Description**: Password too weak (less than 6 characters)
- **Expected Response**: HTTP 400 with `PASSWORD_TOO_WEAK` error code

#### Duplicate User Error
- **Description**: Attempting to register with an existing username
- **Expected Response**: HTTP 409 with `DUPLICATE_USERNAME` error code
- **Usage**: Register `test_valid_user.json` first, then try again

### Login Tests (`/api/v1/auth/login`)

#### Successful Login
- **File**: `test_login_valid.json`
- **Description**: Valid login credentials
- **Expected Response**: HTTP 200 with user data and JWT token

#### Authentication Errors
- **File**: `test_login_nonexistent.json`
- **Description**: Login attempt with non-existent username
- **Expected Response**: HTTP 401 with `USER_NOT_FOUND` error code

- **File**: `test_login_wrong_password.json`
- **Description**: Login attempt with incorrect password
- **Expected Response**: HTTP 401 with `INVALID_PASSWORD` error code

## Error Response Structure

The API now returns detailed error information to help client developers:

### Validation Errors (HTTP 400)
```json
{
  "error": "Validation failed",
  "details": [
    {
      "field": "username",
      "message": "Username is required and must be valid"
    }
  ]
}
```

### Detailed Errors (HTTP 400/401/409)
```json
{
  "error": "Registration failed",
  "code": "DUPLICATE_USERNAME",
  "message": "A user with this username already exists. Please choose a different username.",
  "details": {
    "username": "testuser123"
  }
}
```

### Simple Errors (HTTP 500)
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

## Error Codes Reference

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `INVALID_USERNAME` | 400 | Username doesn't meet requirements (length, format) |
| `PASSWORD_TOO_WEAK` | 400 | Password doesn't meet minimum security requirements |
| `DUPLICATE_USERNAME` | 409 | Username already exists in the system |
| `USER_NOT_FOUND` | 401 | No user found with the provided username |
| `INVALID_PASSWORD` | 401 | The password provided is incorrect |

## Notes for Client Developers

1. **Always check the `code` field** in error responses for programmatic handling
2. **Use `details` object** for field-specific error information
3. **Empty or malformed JSON** returns validation errors with field-level details
4. **All endpoints require `Content-Type: application/json`** header
5. **Successful registration** returns user data without password hash
6. **Successful login** returns both user data and JWT token for authenticated requests