# OpenAPI Specification

## Overview
TheControlApp is WebSocket-first with minimal REST endpoints for authentication and file management.

## Base Information
- **Version:** 1.0.0  
- **Base URL:** `http://localhost:8080`
- **Authentication:** JWT Bearer tokens

## Quick Reference

### Authentication Endpoints
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration

### User Endpoints  
- `GET /api/v1/users` - List users (auth required)

### File Endpoints
- `POST /api/v1/files/upload` - Upload file (auth required)
- `GET /api/v1/files/download/{hash}` - Download file (auth required)

### Utility Endpoints
- `GET /health` - Health check

## WebSocket Connection
- `WS /api/ws` - Primary real-time interface
- Authentication via Bearer token in header or query param

## Rate Limits
- Authentication: 5 requests/minute per IP
- File uploads: 10 files/minute per user  
- General API: 100 requests/minute per user

For detailed request/response formats, see the REST API documentation.
    name: MIT
    url: https://opensource.org/licenses/MIT

servers:
  - url: http://localhost:8080
    description: Development server
  - url: https://api.thecontrolapp.com
    description: Production server

paths:
  /api/auth/register:
    post:
      tags:
        - Authentication
      summary: Register new user account
      description: Create a new user account with username, email, and password
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - username
                - display_name
                - email
                - password
              properties:
                username:
                  type: string
                  minLength: 3
                  maxLength: 50
                  pattern: '^[a-zA-Z0-9_-]+$'
                  example: user123
                display_name:
                  type: string
                  minLength: 1
                  maxLength: 100
                  example: John Doe
                email:
                  type: string
                  format: email
                  example: user@example.com
                password:
                  type: string
                  minLength: 8
                  example: securepassword123
      responses:
        '201':
          description: User created successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: User created successfully
                  user:
                    $ref: '#/components/schemas/User'
        '400':
          description: Invalid input data
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                error: validation_failed
                message: Username must be 3-50 characters
        '409':
          description: Username or email already exists
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                error: user_exists
                message: Username already taken

  /api/auth/login:
    post:
      tags:
        - Authentication
      summary: User login
      description: Authenticate user and receive JWT token
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - username
                - password
              properties:
                username:
                  type: string
                  example: user123
                password:
                  type: string
                  example: securepassword123
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                type: object
                properties:
                  token:
                    type: string
                    example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
                  user:
                    $ref: '#/components/schemas/User'
        '401':
          description: Invalid credentials
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                error: invalid_credentials
                message: Invalid username or password

  /api/auth/me:
    get:
      tags:
        - Authentication
      summary: Get current user info
      description: Retrieve current authenticated user details
      security:
        - BearerAuth: []
      responses:
        '200':
          description: User information
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '401':
          description: Authentication required
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /api/files/upload:
    post:
      tags:
        - Files
      summary: Upload file
      description: Upload a file to the server with automatic deduplication
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          multipart/form-data:
            schema:
              type: object
              required:
                - file
              properties:
                file:
                  type: string
                  format: binary
                  description: File to upload (max 50MB)
      responses:
        '200':
          description: File uploaded successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: File uploaded successfully
                  file_hash:
                    type: string
                    example: abc123def456789...
                  file_name:
                    type: string
                    example: document.pdf
                  size_bytes:
                    type: integer
                    example: 1048576
                  mime_type:
                    type: string
                    example: application/pdf
        '400':
          description: Invalid file or size limit exceeded
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '401':
          description: Authentication required
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /api/files/download/{hash}:
    get:
      tags:
        - Files
      summary: Download file by hash
      description: Download a file using its hash identifier
      security:
        - BearerAuth: []
      parameters:
        - name: hash
          in: path
          required: true
          schema:
            type: string
            pattern: '^[a-f0-9]{64}$'
          example: abc123def456789...
          description: SHA-256 hash of the file
      responses:
        '200':
          description: File content
          content:
            application/octet-stream:
              schema:
                type: string
                format: binary
          headers:
            Content-Disposition:
              schema:
                type: string
              example: attachment; filename="document.pdf"
            Content-Type:
              schema:
                type: string
              example: application/pdf
        '404':
          description: File not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                error: file_not_found
                message: File with hash abc123... not found
        '401':
          description: Authentication required
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'

  /api/users/{username}:
    get:
      tags:
        - Users
      summary: Get user by username
      description: Retrieve public user information by username
      security:
        - BearerAuth: []
      parameters:
        - name: username
          in: path
          required: true
          schema:
            type: string
          example: user123
      responses:
        '200':
          description: User information
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/PublicUser'
        '404':
          description: User not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                error: user_not_found
                message: User 'user123' not found

  /api/users/{username}/block:
    post:
      tags:
        - Users
      summary: Block user
      description: Block a user to prevent receiving commands from them
      security:
        - BearerAuth: []
      parameters:
        - name: username
          in: path
          required: true
          schema:
            type: string
          example: user123
      requestBody:
        required: false
        content:
          application/json:
            schema:
              type: object
              properties:
                reason:
                  type: string
                  maxLength: 255
                  example: Inappropriate behavior
      responses:
        '200':
          description: User blocked successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: User blocked successfully
        '404':
          description: User not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
        '409':
          description: User already blocked
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                error: already_blocked
                message: User is already blocked

    delete:
      tags:
        - Users
      summary: Unblock user
      description: Remove block on a user
      security:
        - BearerAuth: []
      parameters:
        - name: username
          in: path
          required: true
          schema:
            type: string
          example: user123
      responses:
        '200':
          description: User unblocked successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: User unblocked successfully
        '404':
          description: Block not found
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
              example:
                error: block_not_found
                message: User is not blocked

  /api/tags:
    get:
      tags:
        - Tags
      summary: List available tags
      description: Get list of all available content tags
      security:
        - BearerAuth: []
      parameters:
        - name: nsfw
          in: query
          required: false
          schema:
            type: boolean
          description: Include NSFW tags in results
      responses:
        '200':
          description: List of tags
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/Tag'

  /api/reports:
    post:
      tags:
        - Reports
      summary: Submit report
      description: Report inappropriate content or behavior
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - reason
              properties:
                reported_username:
                  type: string
                  example: baduser123
                command_id:
                  type: string
                  format: uuid
                  example: 550e8400-e29b-41d4-a716-446655440001
                reason:
                  type: string
                  enum:
                    - spam
                    - harassment
                    - inappropriate_content
                    - fake_profile
                    - other
                  example: harassment
                description:
                  type: string
                  maxLength: 1000
                  example: User sent threatening messages
      responses:
        '201':
          description: Report submitted successfully
          content:
            application/json:
              schema:
                type: object
                properties:
                  message:
                    type: string
                    example: Report submitted successfully
                  report_id:
                    type: string
                    format: uuid
                    example: 123e4567-e89b-12d3-a456-426614174000

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          format: uuid
          example: 123e4567-e89b-12d3-a456-426614174000
        username:
          type: string
          example: user123
        display_name:
          type: string
          example: John Doe
        email:
          type: string
          format: email
          example: user@example.com
        role:
          type: string
          enum: [user, admin, moderator]
          example: user
        preferences:
          type: object
          example:
            theme: "dark"
            notifications: 
              popup: true
              sound: true
        created_at:
          type: string
          format: date-time
          example: 2025-01-24T15:30:00Z
        updated_at:
          type: string
          format: date-time
          example: 2025-01-24T15:30:00Z

    PublicUser:
      type: object
      properties:
        id:
          type: string
          format: uuid
        username:
          type: string
        display_name:
          type: string
        role:
          type: string
        created_at:
          type: string
          format: date-time

    Tag:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
          example: general
        description:
          type: string
          example: General content and announcements
        color:
          type: string
          pattern: '^#[0-9a-fA-F]{6}$'
          example: "#007bff"
        is_nsfw:
          type: boolean
          example: false
        created_at:
          type: string
          format: date-time

    Command:
      type: object
      properties:
        id:
          type: string
          format: uuid
        instructions:
          type: array
          items:
            $ref: '#/components/schemas/Instruction'
        sender:
          type: string
          example: master123
        receiver:
          type: string
          nullable: true
          example: user456
        tags:
          type: array
          items:
            type: string
          example: ["general", "daily"]
        status:
          type: string
          enum: [pending, delivered, completed, failed]
          example: pending
        created_at:
          type: string
          format: date-time
        completed_at:
          type: string
          format: date-time
          nullable: true

    Instruction:
      type: object
      required:
        - type
        - content
      properties:
        type:
          type: string
          enum:
            - popup-msg
            - download-file
            - display-text
            - timer
            - open-url
            - notification
            - form-input
            - display-image
            - play-audio
            - schedule-task
            - vibrate
            - capture-photo
            - request-location
            - device-control
            - announcement
            - custom
          example: popup-msg
        content:
          type: object
          description: Instruction-specific content (varies by type)
          example:
            body: "Complete your daily task"
            button: "Got it!"

    Error:
      type: object
      required:
        - error
        - message
      properties:
        error:
          type: string
          example: validation_failed
        message:
          type: string
          example: Invalid input data
        field:
          type: string
          description: Field name for validation errors
          example: username
        command_id:
          type: string
          format: uuid
          description: Command ID for command-related errors

    WebSocketMessage:
      oneOf:
        - $ref: '#/components/schemas/CommandMessage'
        - $ref: '#/components/schemas/SendCommandMessage'
        - $ref: '#/components/schemas/CommandStatusMessage'
        - $ref: '#/components/schemas/HeartbeatMessage'
        - $ref: '#/components/schemas/ErrorMessage'

    CommandMessage:
      type: object
      required:
        - type
        - data
      properties:
        type:
          type: string
          enum: [command]
        data:
          $ref: '#/components/schemas/Command'

    SendCommandMessage:
      type: object
      required:
        - type
        - data
      properties:
        type:
          type: string
          enum: [send_command]
        data:
          type: object
          required:
            - instructions
          properties:
            instructions:
              type: array
              items:
                $ref: '#/components/schemas/Instruction'
            receiver:
              type: string
              description: Username to send command to (omit for broadcast)
            tags:
              type: array
              items:
                type: string
              description: Tags for content filtering

    CommandStatusMessage:
      type: object
      required:
        - type
        - data
      properties:
        type:
          type: string
          enum: [command_status]
        data:
          type: object
          required:
            - command_id
            - status
          properties:
            command_id:
              type: string
              format: uuid
            status:
              type: string
              enum: [received, completed, failed]
            completed_at:
              type: string
              format: date-time

    HeartbeatMessage:
      type: object
      required:
        - type
        - timestamp
      properties:
        type:
          type: string
          enum: [heartbeat]
        timestamp:
          type: string
          format: date-time

    ErrorMessage:
      type: object
      required:
        - type
        - error
        - message
      properties:
        type:
          type: string
          enum: [error]
        error:
          type: string
        message:
          type: string
        command_id:
          type: string
          format: uuid
          description: Optional command ID if error is command-related

security:
  - BearerAuth: []

tags:
  - name: Authentication
    description: User authentication and session management
  - name: Files
    description: File upload and download operations  
  - name: Users
    description: User management and social features
  - name: Tags
    description: Content categorization and filtering
  - name: Reports
    description: Content and user reporting system

externalDocs:
  description: Find more info about TheControlApp
  url: https://github.com/TheControlApp/server
```

## Code Generation

### Generate Client SDK

```bash
# Install OpenAPI Generator
npm install -g @openapitools/openapi-generator-cli

# Generate JavaScript client
openapi-generator-cli generate \
  -i swagger.yaml \
  -g javascript \
  -o ./clients/javascript \
  --additional-properties=projectName=thecontrolapp-client

# Generate Python client  
openapi-generator-cli generate \
  -i swagger.yaml \
  -g python \
  -o ./clients/python \
  --additional-properties=packageName=thecontrolapp_client

# Generate Go client
openapi-generator-cli generate \
  -i swagger.yaml \
  -g go \
  -o ./clients/go \
  --additional-properties=packageName=thecontrolapp
```

### Integration with Go Server

```go
// cmd/server/main.go
import (
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    _ "github.com/TheControlApp/server/docs" // Generated swagger docs
)

func setupSwagger(r *gin.Engine) {
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
```

### Generate Swagger Documentation

```bash
# Install swag for Go
go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs from Go comments
swag init -g cmd/server/main.go -o docs/

# Generate OpenAPI spec from existing code
swag init --parseDependency --parseInternal -g cmd/server/main.go
```

## API Annotations in Go Code

### Authentication Handler Example

```go
// @Summary User login
// @Description Authenticate user and receive JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func LoginHandler(c *gin.Context) {
    // Implementation
}

type LoginRequest struct {
    Username string `json:"username" binding:"required" example:"user123"`
    Password string `json:"password" binding:"required" example:"password"`
}

type LoginResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
    User  User   `json:"user"`
}
```

### File Upload Handler Example

```go
// @Summary Upload file
// @Description Upload a file with automatic deduplication
// @Tags Files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} FileUploadResponse
// @Failure 400 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/files/upload [post]
func UploadFileHandler(c *gin.Context) {
    // Implementation
}
```

## Documentation Hosting

### Serve Documentation

```go
// Serve swagger UI at /docs
func setupDocs(r *gin.Engine) {
    r.Static("/docs", "./docs/swagger-ui")
    r.GET("/docs/swagger.json", func(c *gin.Context) {
        c.File("./docs/swagger.json")
    })
}
```

### Auto-generated Documentation

The OpenAPI specification above can be automatically generated from Go code annotations and served at development and production endpoints for API exploration and client SDK generation.

This provides comprehensive API documentation that stays synchronized with the actual implementation and enables automatic client generation for multiple programming languages.
