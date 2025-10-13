# Authentication System

## Overview

The ControlApp Server implements a comprehensive authentication system using JWT (JSON Web Tokens) and bcrypt password hashing. The system supports both REST API and WebSocket authentication with flexible user management capabilities.

## Architecture

### Components

1. **AuthService**: Main authentication service combining all functionality
2. **JWTManager**: Handles JWT token generation and validation
3. **PasswordManager**: Manages password hashing and verification
4. **UserService**: User management with authentication integration

### Security Features

- **JWT Token-based Authentication**: Stateless, scalable authentication
- **bcrypt Password Hashing**: Industry-standard password security
- **Token Expiration**: Configurable token lifetime
- **Multiple Authentication Methods**: Header, query parameter, and message-based auth
- **Input Validation**: Comprehensive user input validation
- **Error Handling**: Structured error responses following RFC 7807

## Authentication Flow

### User Registration

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant UserService
    participant Database
    
    Client->>API: POST /auth/register
    API->>UserService: CreateUser(request)
    UserService->>UserService: Validate username
    UserService->>UserService: Validate password
    UserService->>Database: Check username exists
    UserService->>UserService: Hash password
    UserService->>Database: Create user
    Database-->>UserService: User created
    UserService-->>API: User object
    API-->>Client: 201 Created + User data
```

### User Login

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant UserService
    participant AuthService
    participant Database
    
    Client->>API: POST /auth/login
    API->>UserService: AuthenticateUser(username, password)
    UserService->>Database: Find user by username
    UserService->>AuthService: VerifyPassword(password, hash)
    AuthService-->>UserService: Password verified
    UserService->>Database: Update login date
    UserService-->>API: User object
    API->>AuthService: GenerateToken(userID)
    AuthService-->>API: JWT token
    API-->>Client: 200 OK + User data + Token
```

### Token Validation

```mermaid
sequenceDiagram
    participant Client
    participant API
    participant AuthService
    
    Client->>API: Request with Authorization header
    API->>AuthService: ValidateToken(token)
    AuthService->>AuthService: Parse JWT
    AuthService->>AuthService: Verify signature
    AuthService->>AuthService: Check expiration
    AuthService-->>API: Claims object
    API-->>Client: Authorized response
```

## JWT Token Management

### Token Structure

```json
{
  "user_id": "123e4567-e89b-12d3-a456-426614174000",
  "iat": 1697130000,
  "exp": 1697216400
}
```

### Token Generation

```go
// Generate token for user
func (jm *JWTManager) GenerateToken(userID uuid.UUID) (string, error) {
    claims := Claims{
        UserID: userID.String(),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.duration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jm.secretKey)
}
```

### Token Validation

```go
// Validate token and extract claims
func (jm *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return jm.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### Token Configuration

```yaml
# config.yaml
auth:
  jwt_secret: "your-secret-key-here"
  jwt_expiration: 86400  # 24 hours in seconds
```

## Password Security

### Password Hashing

The system uses bcrypt with default cost (currently 10) for password hashing:

```go
// Hash password using bcrypt
func (pm *PasswordManager) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}
```

### Password Verification

```go
// Verify password against hash
func (pm *PasswordManager) VerifyPassword(password, hash string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
```

### Password Requirements

- **Minimum Length**: 6 characters
- **Maximum Length**: 128 characters
- **Character Requirements**: No specific character requirements (flexible for user convenience)

## User Validation

### Username Validation

```go
func ValidateUsername(username string) *ValidationError {
    if len(username) < 3 {
        return &ValidationError{
            Field:   "username",
            Message: "Username must be at least 3 characters long",
            Code:    "MIN_LENGTH",
        }
    }
    if len(username) > 50 {
        return &ValidationError{
            Field:   "username",
            Message: "Username must be no more than 50 characters long",
            Code:    "MAX_LENGTH",
        }
    }
    if strings.TrimSpace(username) != username {
        return &ValidationError{
            Field:   "username",
            Message: "Username cannot have leading or trailing spaces",
            Code:    "INVALID_FORMAT",
        }
    }
    return nil
}
```

### Validation Rules

#### Username (login_name)
- **Length**: 3-50 characters
- **Format**: No leading/trailing spaces
- **Uniqueness**: Must be unique across all users
- **Case Sensitivity**: Case-sensitive

#### Screen Name (display_name)
- **Length**: 3-50 characters (same rules as username)
- **Format**: No leading/trailing spaces
- **Uniqueness**: Must be unique across all users
- **Purpose**: Public display name

#### Password
- **Length**: 6-128 characters
- **Complexity**: No specific requirements
- **Storage**: Hashed using bcrypt

## Authentication Methods

### 1. REST API Authentication

#### Authorization Header (Recommended)
```http
GET /api/v1/users
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

#### Example Request
```bash
curl -H "Authorization: Bearer ${JWT_TOKEN}" \
     http://localhost:8080/api/v1/users
```

### 2. WebSocket Authentication

#### Option A: Authorization Header
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/client', [], {
    headers: {
        'Authorization': 'Bearer ' + jwtToken
    }
});
```

#### Option B: Query Parameter
```javascript
const ws = new WebSocket(`ws://localhost:8080/ws/client?token=${jwtToken}`);
```

#### Option C: Message-based Authentication
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/client');

ws.onopen = function() {
    ws.send(JSON.stringify({
        type: 'auth',
        token: jwtToken
    }));
};

ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    if (message.type === 'auth_success') {
        console.log('Authenticated successfully');
    }
};
```

## Error Handling

### Authentication Errors

#### Invalid Credentials (401 Unauthorized)
```json
{
  "type": "unauthorized",
  "title": "Unauthorized Access",
  "status": 401,
  "detail": "Invalid username or password",
  "action": "Please check your credentials and try again"
}
```

#### Expired Token (401 Unauthorized)
```json
{
  "type": "unauthorized",
  "title": "Unauthorized Access", 
  "status": 401,
  "detail": "Token has expired",
  "action": "Please login again to obtain a new token"
}
```

#### Username Conflict (409 Conflict)
```json
{
  "type": "conflict",
  "title": "Resource Conflict",
  "status": 409,
  "detail": "Username already exists",
  "action": "Please choose a different username",
  "instance": {
    "username": "existing_user"
  }
}
```

#### Validation Error (422 Unprocessable Entity)
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
    }
  ]
}
```

## Security Considerations

### Token Security

1. **Secret Key Management**
   - Use strong, randomly generated secret keys
   - Store secret keys securely (environment variables, secrets management)
   - Rotate secret keys periodically

2. **Token Expiration**
   - Set appropriate expiration times (24 hours recommended)
   - Implement token refresh mechanism for long-lived applications
   - Consider shorter expiration for sensitive operations

3. **Token Transmission**
   - Always use HTTPS in production
   - Avoid logging tokens
   - Clear tokens from client storage on logout

### Password Security

1. **Bcrypt Configuration**
   - Default cost factor (10) provides good security/performance balance
   - Consider increasing cost factor for higher security requirements
   - Monitor password hashing performance

2. **Password Storage**
   - Never store plaintext passwords
   - Use secure password recovery mechanisms
   - Implement password complexity requirements if needed

### Rate Limiting

Current implementation does not include rate limiting. Consider implementing:

- **Login Attempts**: Limit failed login attempts per IP/user
- **Registration**: Limit account creation per IP
- **Token Generation**: Limit token requests per user

## Configuration

### Environment Variables

```bash
# JWT Configuration
JWT_SECRET=your-super-secret-key-here
JWT_EXPIRATION=86400

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_NAME=controlapp
DB_USER=postgres
DB_PASSWORD=password
```

### Config File Example

```yaml
# config.yaml
auth:
  jwt_secret: "${JWT_SECRET}"
  jwt_expiration: ${JWT_EXPIRATION:86400}

database:
  host: "${DB_HOST:localhost}"
  port: ${DB_PORT:5432}
  name: "${DB_NAME:controlapp}"
  user: "${DB_USER:postgres}"
  password: "${DB_PASSWORD}"
```

## Integration Examples

### Frontend Integration

#### React/JavaScript Example
```javascript
class AuthService {
  constructor() {
    this.token = localStorage.getItem('jwt_token');
  }

  async login(username, password) {
    const response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ username, password }),
    });

    if (response.ok) {
      const data = await response.json();
      this.token = data.token;
      localStorage.setItem('jwt_token', this.token);
      return data.user;
    } else {
      const error = await response.json();
      throw new Error(error.detail);
    }
  }

  async makeAuthenticatedRequest(url, options = {}) {
    return fetch(url, {
      ...options,
      headers: {
        ...options.headers,
        'Authorization': `Bearer ${this.token}`,
      },
    });
  }

  logout() {
    this.token = null;
    localStorage.removeItem('jwt_token');
  }
}
```

### Backend Integration

#### Middleware Example
```go
func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, responses.NewUnauthorizedError(
                "Missing authorization header",
                "Please include a valid JWT token in the Authorization header",
            ))
            c.Abort()
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := authService.JWTManager.ValidateToken(tokenString)
        if err != nil {
            c.JSON(401, responses.NewUnauthorizedError(
                "Invalid or expired token",
                "Please login again to obtain a new token",
            ))
            c.Abort()
            return
        }

        // Add user ID to context
        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

## Testing

### Unit Tests

```go
func TestPasswordHashing(t *testing.T) {
    pm := auth.NewPasswordManager()
    
    password := "testpassword123"
    hash, err := pm.HashPassword(password)
    assert.NoError(t, err)
    assert.NotEmpty(t, hash)
    
    err = pm.VerifyPassword(password, hash)
    assert.NoError(t, err)
    
    err = pm.VerifyPassword("wrongpassword", hash)
    assert.Error(t, err)
}

func TestJWTTokenGeneration(t *testing.T) {
    jm := auth.NewJWTManager("test-secret", time.Hour)
    userID := uuid.New()
    
    token, err := jm.GenerateToken(userID)
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    
    claims, err := jm.ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, userID.String(), claims.UserID)
}
```

### Integration Tests

```bash
# Test registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass","screen_name":"Test User"}'

# Test login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"testpass"}'

# Test authenticated request
curl -H "Authorization: Bearer ${JWT_TOKEN}" \
  http://localhost:8080/api/v1/users
```

## Troubleshooting

### Common Issues

1. **"Invalid token" errors**
   - Check token expiration
   - Verify JWT secret configuration
   - Ensure proper token format

2. **"Username already exists" during registration**
   - Check both login_name and screen_name uniqueness
   - Verify database constraints

3. **Password verification failures**
   - Ensure password is correctly hashed during registration
   - Check bcrypt compatibility

4. **WebSocket authentication issues**
   - Verify token transmission method
   - Check token format and validity
   - Ensure proper message structure for auth messages

### Debug Mode

Enable debug logging to troubleshoot authentication issues:

```bash
LOG_LEVEL=debug go run cmd/server/main.go
```

This will log detailed information about token validation, user lookup, and authentication flow.