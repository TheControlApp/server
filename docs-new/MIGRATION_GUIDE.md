# Migration Guide: Old to New Documentation

This guide helps users migrate from the old documentation structure to the new comprehensive documentation system. It highlights breaking changes, deprecated features, and provides migration paths for existing implementations.

## 📋 **Quick Migration Checklist**

- [ ] Update WebSocket authentication method (remove `auth_login`)
- [ ] Verify database schema expectations
- [ ] Update error handling to use RFC 7807 format
- [ ] Test WebSocket client implementations
- [ ] Update API endpoint references
- [ ] Verify command message formats
- [ ] Update documentation bookmarks

## 🚨 **Breaking Changes**

### 1. WebSocket Authentication

#### ❌ **OLD (NOT SUPPORTED)**
```javascript
// This NEVER worked - was documented but not implemented
ws.send(JSON.stringify({
    type: 'auth_login',
    payload: {
        username: 'testuser',
        password: 'password123'
    }
}));
```

#### ✅ **NEW (CORRECT)**
```javascript
// Only JWT token authentication is supported
ws.send(JSON.stringify({
    type: 'auth',
    token: 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'
}));
```

**What Changed:**
- Progressive authentication with username/password was never implemented
- Only JWT token authentication via `auth` message type is supported
- Authentication can also be done via query parameter or Authorization header

**Migration Steps:**
1. Remove any `auth_login` message handling code
2. Ensure you have a valid JWT token from `/api/v1/auth/login`
3. Use the `auth` message type with the JWT token
4. Update error handling for authentication failures

### 2. Database Schema Changes

#### ❌ **OLD SCHEMA (DOCUMENTED BUT INCORRECT)**
```sql
-- These fields/tables don't exist in actual implementation
Users:
  - Email (unique, 255 chars)      -- NOT IN SCHEMA
  - Preferences (JSONB)            -- NOT IN SCHEMA

Files:                             -- TABLE DOESN'T EXIST
  - Hash (64 char primary key)
  - Filename and content type

User_Relationships:                -- TABLE DOESN'T EXIST
  - Relationship type (blocked/friend)
```

#### ✅ **NEW SCHEMA (ACTUAL IMPLEMENTATION)**
```sql
-- Actual database schema
Users:
  - ID (UUID primary key)
  - login_name (unique, not email)
  - screen_name (display name)
  - password_hash
  - created_at, updated_at, deleted_at
  - is_verified, is_active, is_admin

-- Separate tables for relationships
Blocks:
  - BlockerID, BlockedID (user relationships)

Reports:
  - ReporterID, ReportedID, ReasonID
```

**What Changed:**
- No email field - uses `login_name` instead
- No Files table - file handling not implemented
- No JSONB preferences - uses simple boolean fields
- Relationships handled by separate Block and Report tables

**Migration Steps:**
1. Update any code expecting email field to use `login_name`
2. Remove file upload/management features (not implemented)
3. Update user relationship handling to use Block/Report models
4. Use boolean fields instead of JSONB for user preferences

### 3. Error Response Format Standardization

#### ❌ **OLD (INCONSISTENT)**
```json
// Various inconsistent error formats were possible
{
  "error": "Something went wrong"
}
```

#### ✅ **NEW (RFC 7807 COMPLIANT)**
```json
{
  "type": "validation_error",
  "title": "Validation Failed",
  "status": 422,
  "detail": "Username must be at least 3 characters",
  "instance": "/api/v1/auth/register",
  "action": "Please provide a valid username"
}
```

**Migration Steps:**
1. Update error handling to expect RFC 7807 format
2. Use the `type` field for programmatic error handling
3. Display the `detail` field to users
4. Follow `action` guidance when provided

## 📚 **Documentation Structure Changes**

### Old Structure
```
docs/
├── README.md                    # Basic overview
├── API_SWAGGER.md              # Swagger documentation
├── WEBSOCKET_IMPLEMENTATION.md # WebSocket guide
├── ERROR_RESPONSE_REFERENCE.md # Error handling
└── examples/                   # Code examples
```

### New Structure
```
docs-new/
├── README.md                   # Getting started guide
├── api/
│   ├── rest-api.md            # Complete REST API reference
│   ├── websocket-api.md       # Complete WebSocket API reference
│   └── authentication.md     # Authentication guide
├── architecture/
│   └── overview.md            # System architecture
├── error-handling.md          # RFC 7807 error guide
├── getting-started/
│   └── quickstart.md         # Quick setup guide
└── examples/
    └── websocket-client.js   # Updated working examples
```

**Migration Steps:**
1. Update bookmarks to new documentation URLs
2. Use new comprehensive API references
3. Follow new authentication guide
4. Use updated code examples

## 🔧 **Updated Examples**

### WebSocket Client Update

#### ❌ **OLD CLIENT (PARTIAL FEATURES)**
```javascript
// Old client had incorrect authentication and missing features
const client = new ControlMeClient('ws://localhost:8080', token);
client.authenticateWithLogin(username, password); // Didn't work
```

#### ✅ **NEW CLIENT (FULL FEATURES)**
```javascript
// Updated client with correct authentication and full feature set
const client = new ControlMeClient('ws://localhost:8080', jwtToken);
client.connect(); // Authentication handled automatically

// Or anonymous connection
const anonClient = new ControlMeClient('ws://localhost:8080');
anonClient.connect();

// Message-based authentication
client.authenticateWithToken(jwtToken);
```

### API Endpoint Updates

#### ❌ **OLD (UNDOCUMENTED ENDPOINTS)**
```javascript
// Some endpoints were missing from old docs
POST /api/v1/users              // Not documented
PUT /api/v1/users/:id           // Not documented  
DELETE /api/v1/users/:id        // Not documented
```

#### ✅ **NEW (COMPLETE DOCUMENTATION)**
```javascript
// All endpoints now documented with examples
POST /api/v1/users              // Create user
GET /api/v1/users               // List users
GET /api/v1/users/:id           // Get user by ID
PUT /api/v1/users/:id           // Update user
DELETE /api/v1/users/:id        // Delete user
```

## 🎯 **Feature Status Clarification**

### ✅ **Fully Implemented & Documented**
- JWT authentication for REST and WebSocket
- User registration and login
- Command creation and management
- Real-time WebSocket communication
- Anonymous and authenticated connections
- RFC 7807 compliant error responses
- User blocking and reporting system

### ⚠️ **Partially Implemented**
- Command instruction execution (basic types only)
- File upload/download (infrastructure exists, limited functionality)
- Administrative features (basic user management)

### ❌ **Not Implemented (Despite Old Documentation)**
- Progressive authentication with username/password
- Complex file management system
- Advanced user relationship features
- Email notifications
- Advanced admin panel

## 🛠 **Code Migration Examples**

### Authentication Migration

#### Old Code (Broken)
```javascript
// This never worked
ws.onopen = () => {
    ws.send(JSON.stringify({
        type: 'auth_login',
        payload: { username: 'user', password: 'pass' }
    }));
};
```

#### New Code (Working)
```javascript
// Method 1: Query parameter authentication
const ws = new WebSocket(`ws://localhost:8080/ws/client?token=${jwtToken}`);

// Method 2: Message-based authentication
ws.onopen = () => {
    ws.send(JSON.stringify({
        type: 'auth',
        token: jwtToken
    }));
};

// Method 3: Header-based authentication (if supported by client)
const ws = new WebSocket('ws://localhost:8080/ws/client', [], {
    headers: { 'Authorization': `Bearer ${jwtToken}` }
});
```

### Error Handling Migration

#### Old Code (Inconsistent)
```javascript
// Old error handling was inconsistent
.catch(error => {
    console.log(error.message); // Might not exist
    // No standard error structure
});
```

#### New Code (RFC 7807)
```javascript
// New standardized error handling
.catch(error => {
    console.error(`${error.type}: ${error.detail}`);
    
    if (error.type === 'validation_error') {
        // Handle validation errors
        error.errors?.forEach(fieldError => {
            showFieldError(fieldError.field, fieldError.message);
        });
    }
    
    if (error.action) {
        showActionGuidance(error.action);
    }
});
```

## 📖 **New Documentation Features**

### What's New
- **Complete API Coverage**: Every endpoint documented with examples
- **Accurate Implementation Details**: Documentation matches actual code
- **RFC 7807 Error Handling**: Standardized error responses
- **Working Code Examples**: All examples tested and verified
- **Architecture Documentation**: Comprehensive system overview
- **Migration Guides**: Clear upgrade paths
- **Validation Results**: API accuracy verification included

### Improved Sections
- **REST API Reference**: Complete with request/response examples
- **WebSocket API Reference**: Corrected authentication methods
- **Authentication Guide**: Comprehensive JWT handling
- **Error Handling Guide**: RFC 7807 compliant examples
- **Architecture Overview**: System design and data flow
- **Getting Started**: Quick setup and testing guide

## 🔍 **Verification Steps**

After migrating, verify your implementation:

1. **Test Authentication**
   ```bash
   # Verify JWT token generation
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"login_name": "testuser", "password": "password"}'
   ```

2. **Test WebSocket Connection**
   ```javascript
   // Verify WebSocket authentication works
   const ws = new WebSocket('ws://localhost:8080/ws/client?token=YOUR_JWT_TOKEN');
   ws.onopen = () => console.log('Connected successfully');
   ws.onmessage = (event) => console.log('Message:', JSON.parse(event.data));
   ```

3. **Test Error Handling**
   ```bash
   # Verify RFC 7807 error format
   curl -X POST http://localhost:8080/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"login_name": "ab"}' # Trigger validation error
   ```

4. **Verify API Responses**
   ```bash
   # Check user endpoints work correctly
   curl -X GET http://localhost:8080/api/v1/users \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

## 📞 **Getting Help**

If you encounter issues during migration:

1. **Check the validation results**: `docs-new/API_VALIDATION_RESULTS.md`
2. **Review working examples**: `docs-new/examples/`
3. **Test with Swagger**: http://localhost:8080/swagger/index.html
4. **Compare with old docs validation**: `.copilot/old_docs_validation.md`

## 📅 **Migration Timeline**

### Immediate (Required)
- Update WebSocket authentication method
- Fix error handling to use RFC 7807 format
- Update documentation bookmarks

### Short Term (Recommended)
- Update code examples and client implementations
- Test all API integrations
- Verify database schema expectations

### Long Term (Optional)
- Adopt new architecture patterns
- Implement improved error handling
- Use new documentation structure for development

---

**Note**: The old documentation is preserved in the `docs/` directory for reference, but should not be used for new development. The new documentation in `docs-new/` is the authoritative source and matches the actual implementation.