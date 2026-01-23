# ControlApp Server Documentation

> **📍 Quick Access**: [Swagger UI](http://localhost:8080/swagger/index.html) | [Health Check](http://localhost:8080/health) | [WebSocket](ws://localhost:8080/ws/client)

## 🚀 Quick Start

```bash
# Docker (Recommended)
make docker-up

# Manual setup
make setup && make dev

# Test connection
curl http://localhost:8080/health
```

**Access Points:**
- 🌐 **API Server**: http://localhost:8080
- 📚 **Swagger UI**: http://localhost:8080/swagger/index.html  
- 🏥 **Health Check**: http://localhost:8080/health
- � **WebSocket**: ws://localhost:8080/ws/client

---

## 📚 API Reference

### REST Endpoints
**Base URL**: `http://localhost:8080/api/v1`

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth/register` | Create user account |
| `POST` | `/auth/login` | Get JWT token |
| `GET` | `/users` | List all users |
| `GET` | `/users/{id}` | Get user by ID |
| `GET` | `/commands/pending` | Get pending commands |
| `POST` | `/commands/complete` | Mark command completed |
| `GET` | `/health` | Server status |

### WebSocket API
```javascript
// Connect
const ws = new WebSocket('ws://localhost:8080/ws/client');

// Send command
ws.send(JSON.stringify({
  instructions: [{
    type: 'std_popup',
    content: { body: 'Hello!', button: 'OK' }
  }],
  tags: 'general',
  receiver_id: null // null = broadcast
}));
```

---

## 🔐 Authentication

JWT tokens required for protected endpoints:
```bash
Authorization: Bearer <token>
```

**Login Example:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "user", "password": "pass"}'
```

---

---

## 📖 Documentation Structure

### Essential Resources
- **[API Reference](API-REFERENCE.md)** - Complete REST + WebSocket API docs
- **[Error Response Guide](ERROR_RESPONSE_REFERENCE.md)** - RFC 7807 error handling  
- **[Swagger UI](swagger/)** - Interactive API testing
- **[Examples](examples/)** - Code samples and integration guides

### Quick Links
- **[Database Schema](database/schema.md)** - Data models reference
- **[Integration Test Tool](../cmd/tools/integration-test/)** - API validation utility

---

## 🧪 Testing Your Setup

```bash
# Quick health check
curl http://localhost:8080/health

# Register test user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "test123", "screen_name": "Test User"}'

# Run integration tests
go run cmd/tools/integration-test/main.go
```

See [Standard Command Types](./standards/command_types.md) for complete specifications.
