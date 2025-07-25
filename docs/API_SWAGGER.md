# ControlMe Go API - Swagger Documentation

This project now includes comprehensive Swagger/OpenAPI documentation for all API endpoints.

## 🚀 Quick Start

### 1. Start the Server
```bash
# Local development
go run ./cmd/server

# Or using Docker
docker-compose up
```

### 2. Access Swagger UI
Once the server is running, you can access the interactive Swagger documentation at:

**🌐 [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

## 📋 Available Endpoints

### Health Check
- **GET** `/health` - Check server status

### Authentication
- **POST** `/api/v1/auth/login` - User login
- **POST** `/api/v1/auth/register` - User registration

### Users
- **GET** `/api/v1/users` - Get all users
- **GET** `/api/v1/users/{id}` - Get user by ID
- **POST** `/api/v1/users` - Create new user

### Commands
- **GET** `/api/v1/commands/pending` - Get pending commands for a user
- **POST** `/api/v1/commands/complete` - Mark command as completed

### WebSocket
- **GET** `/ws/client` - Client WebSocket connection
- **GET** `/ws/web` - Web WebSocket connection

## 🔧 API Documentation Features

### Response Models
The API now uses structured response models instead of generic maps:

- `AuthResponse` - Login/authentication responses
- `UserResponse` - Single user responses
- `UsersResponse` - Multiple users responses
- `CommandsResponse` - Commands list responses
- `MessageResponse` - Simple success messages
- `ErrorResponse` - Error responses
- `HealthResponse` - Health check responses

### Request Models
- `LoginRequest` - User login credentials
- `RegisterRequest` - User registration data
- `CreateUserRequest` - User creation data

## 📖 Using the Swagger UI

1. **Interactive Testing**: Click on any endpoint to expand it and see detailed information
2. **Try It Out**: Use the "Try it out" button to test endpoints directly from the browser
3. **Model Documentation**: Click on model names to see the structure of request/response objects
4. **Authentication**: Some endpoints may require authentication (JWT tokens)

## 🛠️ Development

### Regenerating Swagger Documentation

When you make changes to the API handlers or add new endpoints, regenerate the documentation:

```bash
# Install swag tool (one-time setup)
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g cmd/server/main.go -o docs/swagger
```

### Adding New Endpoints

1. Add Swagger annotations to your handler functions:
```go
// GetUsers godoc
// @Summary      Get all users
// @Description  Retrieves a list of all users
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      200  {object}  responses.UsersResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /users [get]
func (h *UserHandlers) GetUsers(c *gin.Context) {
    // implementation
}
```

2. Regenerate the documentation using `swag init`

### Response Model Standards

- Use structured response models from `internal/api/responses`
- Always include proper HTTP status codes
- Provide meaningful error messages
- Include example values in response models

## 🔍 API Testing

You can test the API using:

1. **Swagger UI** (recommended for exploration)
2. **curl** commands
3. **Postman** or similar tools
4. **HTTP clients** in your applications

### Example curl Commands

```bash
# Health check
curl http://localhost:8080/health

# User registration
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123",
    "screen_name": "Test User",
    "email": "test@example.com",
    "random_opt_in": false
  }'

# User login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "password123"
  }'
```

## 📚 Additional Resources

- [Swagger/OpenAPI Specification](https://swagger.io/specification/)
- [swaggo Documentation](https://github.com/swaggo/swag)
- [Gin Framework](https://gin-gonic.com/)

---

**Happy API Development! 🚀**
