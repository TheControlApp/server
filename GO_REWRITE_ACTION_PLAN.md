# ControlMe Go Rewrite Action Plan

## Project Overview
Rewriting the ControlMe consensual remote control platform from ASP.NET/C# to Go for better portability, security, and maintainability.

## Current System Analysis

### Existing Architecture
- **Backend**: ASP.NET Web Forms (C#/.NET Framework 4.8)
- **Database**: SQL Server with stored procedures
- **Frontend**: ASP.NET Web Forms with basic HTML/CSS
- **Desktop Client**: .NET executable polling every 30 seconds
- **Authentication**: Cookie-based with custom crypto
- **Deployment**: IIS on Windows Server

### Key Vulnerabilities Identified
- SQL injection risks (string concatenation in queries)
- Weak encryption implementation
- No input validation
- Cookie-based auth without proper security
- Single platform deployment (Windows only)

## Target Architecture

### Technology Stack
- **Backend**: Go 1.21+ with Gin web framework
- **Database**: PostgreSQL 15+ (more portable than SQL Server)
- **Frontend**: Modern web interface (React/Vue or Go templates)
- **Desktop Client**: Go binary (cross-platform)
- **Authentication**: JWT tokens with proper security
- **Encryption**: Go standard library crypto packages
- **Deployment**: Docker containers + single binary

### System Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │    │   Go Backend    │    │   PostgreSQL    │
│   (React/Vue)   │◄──►│   (REST API)    │◄──►│   Database      │
│                 │    │   + WebSocket   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              ▲
                              │ HTTP/WS
                    ┌─────────────────┐
                    │  Desktop Client │
                    │   (Go binary)   │
                    └─────────────────┘
```

### WebSocket Architecture for Real-time Features

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │    │   Go Backend    │    │   PostgreSQL    │
│   (React/Vue)   │◄──►│   (REST API)    │◄──►│   Database      │
│                 │    │   + WebSocket   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       
         │              ┌─────────────────┐              
         │              │  WebSocket Hub  │              
         │              │   (In-Memory)   │              
         │              └─────────────────┘              
         │                       │                       
         │                       ▲                       
         │                       │ WebSocket             
         │             ┌─────────────────┐               
         │             │  Desktop Client │               
         │             │   (Go binary)   │               
         │             └─────────────────┘               
         │                                                
         └─── WebSocket Connection ──────────────────────┘
```

#### WebSocket Message Types
- **Command Messages**: Instant command delivery to desktop clients
- **Status Updates**: Command completion, errors, client online/offline
- **Chat Messages**: Real-time messaging between users
- **Presence Updates**: User online/offline status, typing indicators
- **Notifications**: System alerts, new invites, reports
- **Heartbeat**: Keep-alive messages for connection health

#### WebSocket Hub Implementation
- Connection pool management (users, groups, admins)
- Message routing and broadcasting
- User session tracking
- Rate limiting and spam protection
- Connection cleanup and garbage collection
- Message queuing for offline users

## Phase 1: Project Setup & Database Migration

### 1.1 Initialize Go Project
- [x] Create Go module structure
- [x] Set up development environment
- [x] Configure project dependencies
- [ ] Set up CI/CD pipeline basics

### 1.2 Database Migration
- [x] Install PostgreSQL locally (via Docker)
- [x] Create database schema migration tool (GORM auto-migrate)
- [x] Convert SQL Server tables to PostgreSQL:
  - [x] Users table (modern schema with UUIDs)
  - [x] ControlAppCmd table (modern schema)
  - [x] CommandList table (renamed to Commands, modern schema)
  - [x] ChatLog table
  - [x] Groups table
  - [x] Relationship table
  - [x] Block table
  - [x] Report table
  - [x] Invites table
- [✅] Modern business logic implementation (replaces stored procedures)
- [x] Create database connection pool
- [x] Implement database models with GORM/sqlx
- [🔄] Legacy compatibility layer for client responses (IN PROGRESS)
- [ ] Backup and rollback procedures

## ARCHITECTURE STRATEGY:

### Modern Backend (Core System):
- **Database**: Modern PostgreSQL schema with UUIDs, proper relationships, timestamps
- **Authentication**: JWT tokens, bcrypt password hashing, modern security practices
- **API Design**: RESTful endpoints, proper HTTP status codes, JSON responses
- **Business Logic**: Clean service layer, proper error handling, validation
- **Security**: Rate limiting, CORS, HTTPS, input validation, SQL injection prevention

### Legacy Compatibility Layer:
- **Purpose**: Translate between modern backend and legacy client expectations
- **Scope**: Only client-facing API responses need exact compatibility
- **Implementation**: Wrapper functions that call modern services and format responses
- **Authentication**: Accept legacy encrypted passwords, translate to modern auth internally
- **Response Format**: HTML/text responses matching ASP.NET exactly

### Benefits of This Approach:
1. **Clean Modern Codebase**: Not constrained by legacy limitations
2. **Security**: Modern authentication and security practices from day one
3. **Maintainability**: Proper separation of concerns and clean architecture
4. **Flexibility**: Easy to deprecate legacy compatibility layer later
5. **Performance**: Modern database design and query optimization
6. **Scalability**: Built for modern deployment (Docker, Kubernetes, etc.)

### 1.3 Security Improvements
- [✅] Implement proper password hashing (bcrypt) - FOR NEW USERS/MODERN API
- [✅] Set up JWT authentication - IMPLEMENTED FOR MODERN API
- [🔄] Add input validation and sanitization - IN PROGRESS
- [📋] Implement rate limiting - PLANNED
- [📋] Add CORS configuration - PLANNED  
- [📋] Set up HTTPS/TLS - PLANNED

### 1.4 Legacy Compatibility Layer
- [🔄] Legacy authentication bridge (decrypt legacy passwords, authenticate internally)
- [🔄] Legacy response formatters (HTML responses matching ASP.NET exactly)
- [🔄] Legacy business logic translators (convert modern data to legacy format)
- [📋] User migration utilities (convert plain text passwords to hashed)
- [📋] Legacy client detection and upgrade notifications

## STRATEGIC CLARIFICATION:
**Project Goal: Modern Backend with Legacy Client Compatibility**

This is a complete backend rewrite focusing on modern architecture and standards, with a compatibility layer for legacy clients during transition. The internal backend can be completely overhauled - only client-facing API responses need exact compatibility.

**REVISED APPROACH:**
1. **Modern Backend**: Use modern Go patterns, UUIDs, proper authentication, security
2. **Legacy Compatibility Layer**: Translation layer between modern backend and legacy client expectations
3. **Gradual Migration**: Support both legacy and modern clients simultaneously
4. **Clean Architecture**: Separate concerns between modern API and legacy compatibility

**CURRENT STATUS: LEGACY COMPATIBILITY LAYER IMPLEMENTED ✅**

**Date:** July 11, 2025  
**Phase:** Legacy Compatibility Complete, Modern Backend Ready for Development

### Recently Completed (July 11, 2025)
- ✅ **Legacy Database Models**: Complete mapping to legacy schema with exact field names
- ✅ **Legacy Service Layer**: All stored procedures implemented (USP_Login, USP_GetOutstanding, USP_GetAppContent, etc.)
- ✅ **Legacy Handlers**: Complete rewrite for exact ASP.NET compatibility
- ✅ **Authentication Bridge**: Legacy crypto integration with modern security
- ✅ **Response Format Matching**: Exact HTML and bracketed response formats
- ✅ **Service Architecture**: Clean separation between legacy and modern components
- ✅ **Testing Framework**: Legacy compatibility test scripts created

### Architecture Overview
```
┌─────────────────────────────────────────────────────────────┐
│                    LEGACY COMPATIBILITY LAYER               │
├─────────────────────────────────────────────────────────────┤
│  Legacy Handlers → Legacy Service → Legacy Models          │
│  (ASP.NET format)   (Stored Procs)   (Integer IDs)         │
├─────────────────────────────────────────────────────────────┤
│                    MODERN BACKEND CORE                      │
├─────────────────────────────────────────────────────────────┤
│  Modern Handlers → Modern Services → Modern Models         │
│  (REST/GraphQL)     (Business Logic)  (UUID-based)         │
├─────────────────────────────────────────────────────────────┤
│                    SHARED INFRASTRUCTURE                    │
├─────────────────────────────────────────────────────────────┤
│  Database Layer • Security • Configuration • Logging       │
└─────────────────────────────────────────────────────────────┘
```

## Phase 2: Core API Development

### 2.1 Authentication System
- [ ] User registration endpoint
- [ ] Login endpoint with JWT generation
- [ ] JWT validation middleware
- [ ] Password reset functionality
- [ ] User verification system
- [ ] Session management

### 2.2 Desktop Client API Endpoints
- [ ] `/api/v1/auth/login` - Client authentication
- [ ] `/api/v1/commands/pending` - Get pending commands count (fallback)
- [ ] `/api/v1/commands/fetch` - Fetch command details (fallback)
- [ ] `/api/v1/commands/complete` - Mark command as completed
- [ ] `/api/v1/commands/delete` - Delete outstanding commands
- [ ] `/api/v1/client/heartbeat` - Client status updates
- [ ] `/ws/client` - WebSocket endpoint for real-time command delivery

### 2.2.1 Legacy Client Support Endpoints (COMPLETED ✅)
- [✅] `/AppCommand.aspx` - Legacy command polling (EXACT ASP.NET COMPATIBILITY)
  - ✅ Outstanding command count with next sender info
  - ✅ Content retrieval with sender details
  - ✅ Accept/Reject invite handling
  - ✅ Thumbs up functionality
  - ✅ Delete outstanding commands
  - ✅ Invite and relationship listing
- [✅] `/GetContent.aspx` - Legacy command fetching (EXACT ASP.NET COMPATIBILITY)
  - ✅ Sender ID and content retrieval
  - ✅ Automatic command completion
  - ✅ HTML response format matching
- [✅] `/GetCount.aspx` - Legacy command count (EXACT ASP.NET COMPATIBILITY)
  - ✅ Count, next sender, and verification status
  - ✅ Three-label HTML response format
- [✅] `/ProcessComplete.aspx` - Legacy command completion (EXACT ASP.NET COMPATIBILITY)
- [✅] `/DeleteOut.aspx` - Legacy command deletion (EXACT ASP.NET COMPATIBILITY)
- [✅] `/GetOptions.aspx` - Legacy options (EXACT ASP.NET COMPATIBILITY)
- [✅] Legacy authentication bridge (decrypt legacy passwords, use legacy database)
- [✅] Version detection via `vrs=012` parameter (exact match)
- [✅] Exact query string parameter parsing (`usernm`, `pwd`, `vrs`, `cmd`)
- [✅] Legacy response format matching (HTML/text as expected by .NET client)
- [✅] All stored procedures implemented (USP_Login, USP_GetOutstanding, USP_GetAppContent, etc.)
- [✅] Legacy database schema support with exact field names

### 2.3 Web Interface API Endpoints
- [ ] User management endpoints
- [ ] Command sending endpoints (with WebSocket push)
- [ ] Chat/messaging endpoints
- [ ] Group management endpoints
- [ ] Blocking/reporting endpoints
- [ ] File upload endpoints
- [ ] User settings endpoints
- [ ] `/ws/web` - WebSocket endpoint for web client real-time features

### 2.4 Real-time Features (WebSocket Integration)
- [ ] WebSocket connection handler and hub
- [ ] User session management for WebSocket connections
- [ ] Real-time command delivery system
- [ ] Live chat functionality with typing indicators
- [ ] Real-time command status updates (sent/received/completed)
- [ ] Online user presence and activity status
- [ ] WebSocket authentication and authorization
- [ ] Connection heartbeat and reconnection logic
- [ ] Message broadcasting to groups/individuals
- [ ] Real-time notifications system
- [ ] Legacy client upgrade notifications via WebSocket (for mixed environments)

## Phase 2.5: Legacy Migration Support

### 2.5.1 Legacy Client Bridge (Exact Compatibility Layer)
- [ ] Exact ASP.NET Web Forms URL routing (`/AppCommand.aspx`, `/GetContent.aspx`, etc.)
- [ ] Exact query parameter handling (`usernm`, `pwd`, `vrs`, `cmd`)
- [ ] Legacy crypto decryption for passwords (exact algorithm match)
- [ ] Exact response format matching (XML/plain text as current system)
- [ ] Legacy user session management with exact cookie handling
- [ ] Legacy file upload handling with exact form data processing
- [ ] Exact error message format matching for client compatibility
- [ ] Legacy stored procedure result format translation
- [ ] Exact HTTP status code and header matching

### 2.5.2 Upgrade Notification System
- [ ] Version detection middleware
- [ ] Configurable upgrade message templates
- [ ] Progressive notification escalation
- [ ] Download link generation and tracking
- [ ] Usage analytics for legacy vs new clients
- [ ] Admin dashboard for migration progress

### 2.5.3 Migration Tools
- [ ] User data migration verification
- [ ] Command history migration
- [ ] Legacy client usage reports
- [ ] Automated migration testing
- [ ] Rollback procedures for failed migrations
- [ ] Migration status tracking per user

## Phase 3: Web Frontend Development

### 3.1 Frontend Framework Setup
- [ ] Choose framework (React recommended)
- [ ] Set up build pipeline
- [ ] Configure API client
- [ ] Set up routing
- [ ] Implement responsive design

### 3.2 Core Pages
- [ ] Login/Registration pages
- [ ] Dashboard/Home page
- [ ] Command sending interface
- [ ] Chat/messaging interface
- [ ] User profile management
- [ ] Group management
- [ ] Settings page

### 3.3 Admin Interface
- [ ] User management
- [ ] System monitoring
- [ ] Command logging viewer
- [ ] Report management

## Phase 4: Desktop Client Rewrite

### 4.1 Go Desktop Client
- [ ] Cross-platform GUI framework (Fyne/Wails)
- [ ] System tray integration
- [ ] WebSocket connection for real-time commands
- [ ] Fallback HTTP polling mechanism
- [ ] Command execution engine
- [ ] Auto-update functionality
- [ ] Configuration management
- [ ] Connection resilience and reconnection logic

### 4.2 Command Types Implementation
- [ ] Message box display
- [ ] Website opening
- [ ] File download
- [ ] Wallpaper changing
- [ ] File execution
- [ ] Custom command execution

### 4.3 Security & Safety
- [ ] Command validation
- [ ] Sandboxed execution
- [ ] User consent prompts
- [ ] Emergency stop functionality
- [ ] Audit logging

## Phase 5: Advanced Features

### 5.1 Enhanced Security
- [ ] Two-factor authentication
- [ ] End-to-end encryption for sensitive commands
- [ ] Advanced audit logging
- [ ] Intrusion detection
- [ ] Automated security scanning

### 5.2 Performance Optimizations
- [ ] Database query optimization
- [ ] Caching implementation (Redis)
- [ ] CDN integration for static assets
- [ ] Load balancing support
- [ ] Database connection pooling

### 5.3 Monitoring & Observability
- [ ] Structured logging (logrus/zap)
- [ ] Metrics collection (Prometheus)
- [ ] Health check endpoints
- [ ] Performance monitoring
- [ ] Error tracking

## Phase 6: Testing & Quality Assurance

### 6.1 Testing Strategy
- [ ] Unit tests for all business logic
- [ ] Integration tests for API endpoints
- [ ] End-to-end tests for critical flows
- [ ] Load testing for performance
- [ ] Security testing (OWASP)

### 6.2 Code Quality
- [ ] Code coverage > 80%
- [ ] Static analysis (golangci-lint)
- [ ] Dependency vulnerability scanning
- [ ] Code review processes
- [ ] Documentation

## Phase 7: Deployment & Migration

### 7.1 Infrastructure Setup
- [ ] Docker containerization
- [ ] Kubernetes deployment (optional)
- [ ] Database migration scripts
- [ ] Backup and recovery procedures
- [ ] SSL certificate management

### 7.2 Migration Strategy
- [ ] Parallel deployment setup
- [ ] Data migration tools
- [ ] User migration plan
- [ ] Rollback procedures
- [ ] Performance comparison

### 7.3 Go-Live Preparation
- [ ] Load testing in production environment
- [ ] Security audit
- [ ] User acceptance testing
- [ ] Documentation completion
- [ ] Support procedures

## Phase 8: Post-Launch

### 8.1 Monitoring & Maintenance
- [ ] Performance monitoring
- [ ] Error tracking and resolution
- [ ] Security updates
- [ ] Feature enhancements
- [ ] User feedback integration

### 8.2 Future Enhancements
- [ ] Mobile app development
- [ ] Advanced command scheduling
- [ ] Plugin system for custom commands
- [ ] Multi-language support
- [ ] Advanced analytics

## Technical Specifications

### Dependencies
```go
// Core framework
github.com/gin-gonic/gin
github.com/gorilla/websocket

// Database
github.com/lib/pq
gorm.io/gorm
github.com/golang-migrate/migrate

// Authentication
github.com/golang-jwt/jwt/v5
golang.org/x/crypto/bcrypt

// Validation
github.com/go-playground/validator/v10

// Configuration
github.com/spf13/viper

// Logging
github.com/sirupsen/logrus

// WebSocket and Real-time
github.com/gorilla/websocket
github.com/go-redis/redis/v8 (for WebSocket scaling)

// Message Queue (optional for high load)
github.com/streadway/amqp (RabbitMQ)

// Legacy Support
github.com/gorilla/mux (for legacy URL routing)
golang.org/x/crypto/des (for legacy crypto compatibility)
```

### Database Schema Considerations
- Use UUID for primary keys instead of auto-increment
- Add proper indexes for performance
- Implement soft deletes for audit trails
- Add created_at/updated_at timestamps
- Use JSONB for flexible command payloads

### Security Requirements
- All passwords must be hashed with bcrypt
- JWT tokens with proper expiration
- Input validation on all endpoints
- Rate limiting on authentication endpoints
- HTTPS only in production
- SQL injection prevention with prepared statements

### Performance Targets
- API response time < 200ms (95th percentile)
- WebSocket message delivery < 50ms
- Support 1000+ concurrent WebSocket connections
- Database queries < 50ms average
- Desktop client memory usage < 50MB
- Zero-downtime deployments
- WebSocket reconnection < 5 seconds

### WebSocket Implementation Details

#### Connection Management
- **Connection Pool**: In-memory map of active WebSocket connections
- **User Sessions**: Track multiple connections per user (web + desktop)
- **Connection Types**: Separate handling for web clients vs desktop clients
- **Authentication**: JWT token validation for WebSocket connections
- **Rate Limiting**: Per-connection message rate limiting
- **Graceful Shutdown**: Clean connection closure and cleanup

#### Message Flow for Commands
1. **Command Sending**: Web user sends command via REST API
2. **WebSocket Push**: Server immediately pushes command to target's desktop client
3. **Status Updates**: Desktop client sends status updates via WebSocket
4. **Real-time Feedback**: Web user receives instant status updates
5. **Fallback Mechanism**: REST API polling as backup for WebSocket failures

#### WebSocket Message Format
```json
{
  "type": "command|status|chat|presence|notification",
  "id": "unique-message-id",
  "timestamp": "2025-07-09T12:00:00Z",
  "from": "user-id",
  "to": "user-id|group-id",
  "data": {
    // Message-specific payload
  }
}
```

#### Performance Considerations
- **Connection Pooling**: Efficient memory management for thousands of connections
- **Message Queuing**: Queue messages for offline users
- **Horizontal Scaling**: Redis pub/sub for multi-server WebSocket scaling
- **Compression**: WebSocket message compression for large payloads
- **Heartbeat**: Regular ping/pong to detect dead connections

## Risk Assessment

### High Priority Risks
1. **Data Loss During Migration** - Mitigation: Comprehensive backup and testing
2. **Security Vulnerabilities** - Mitigation: Security audit and penetration testing
3. **User Adoption** - Mitigation: Gradual rollout and user training
4. **Performance Degradation** - Mitigation: Load testing and optimization
5. **Legacy Client Compatibility** - Mitigation: Extensive testing and gradual migration

### Medium Priority Risks
1. **Third-party Dependencies** - Mitigation: Dependency management and alternatives
2. **Cross-platform Compatibility** - Mitigation: Extensive testing on all platforms
3. **Scalability Issues** - Mitigation: Horizontal scaling architecture

## Success Metrics

### Technical Metrics
- 99.9% uptime
- < 200ms API response time
- Zero security vulnerabilities
- 100% test coverage on critical paths
- > 90% legacy client migration rate within 12 weeks

### User Metrics
- User satisfaction > 90%
- Bug reports < 1% of active users
- Feature adoption > 80%
- Support tickets < 5% of active users
- Legacy client complaints < 10% during migration

## Timeline Estimate

- **Phase 1**: 4-6 weeks
- **Phase 2**: 6-8 weeks  
- **Phase 3**: 4-6 weeks
- **Phase 4**: 6-8 weeks
- **Phase 5**: 4-6 weeks
- **Phase 6**: 3-4 weeks
- **Phase 7**: 2-3 weeks
- **Phase 8**: Ongoing

**Total Estimated Time**: 6-9 months for full rewrite and deployment

## Next Steps

1. Review and approve this action plan
2. Set up development environment
3. Begin Phase 1 implementation
4. Establish regular review and update cycles
5. Assign team members to specific phases

---

*This document is a living plan and should be updated as the project progresses and requirements evolve.*

## Legacy Support Strategy

### Overview
Maintain compatibility with existing .NET desktop clients during the transition period while encouraging users to upgrade to the new Go-based client.

### Legacy Client Support Architecture (In-Place Swap)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   New Go Client │    │   Go Backend    │    │   PostgreSQL    │
│   (WebSocket)   │◄──►│   (REST API)    │◄──►│   Database      │
│                 │    │   + WebSocket   │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              ▲
                              │ HTTP (Exact .aspx paths)
                    ┌─────────────────┐
                    │  Legacy Client  │
                    │  (.NET/30sec)   │
                    │  /AppCommand.aspx│
                    │  /GetContent.aspx│
                    └─────────────────┘
```

### Legacy Endpoints (Exact Compatibility for In-Place Swap)
- [ ] `/AppCommand.aspx` - Legacy command polling (30-second interval)
- [ ] `/GetContent.aspx` - Legacy command fetching  
- [ ] `/GetCount.aspx` - Legacy command count
- [ ] `/ProcessComplete.aspx` - Legacy command completion
- [ ] `/DeleteOut.aspx` - Legacy command deletion
- [ ] `/Messages.aspx` - Legacy messaging support
- [ ] `/Upload.aspx` - Legacy file upload
- [ ] `/GetOptions.aspx` - Legacy options retrieval
- [ ] `/AppSendContent.aspx` - Legacy content sending
- [ ] `/BlockReport.aspx` - Legacy blocking/reporting
- [ ] `/NGROK.aspx` - Legacy NGROK functionality
- [ ] `/Default.aspx` - Legacy home page
- [ ] `/Pages/Login.aspx` - Legacy login page
- [ ] `/Pages/Register.aspx` - Legacy registration page
- [ ] `/Pages/ControlPC.aspx` - Legacy control interface

### Legacy Authentication
- [ ] Support existing username/password authentication
- [ ] Convert legacy crypto to modern JWT for new features
- [ ] Maintain cookie-based sessions for legacy clients
- [ ] Gradual migration path for user credentials

### Version Detection and Upgrade Prompts
- [ ] Client version detection in all legacy endpoints
- [ ] Configurable upgrade message system
- [ ] Progressive upgrade notifications (daily -> hourly -> constant)
- [ ] Upgrade download links and instructions
- [ ] Grace period configuration for legacy client support

### Update Notification System
- [ ] **Gentle Reminders**: "New version available! Check it out!"
- [ ] **Persistent Notifications**: "Hey! There's a new version out! Fucking update!"
- [ ] **Urgent Warnings**: "Legacy support ending soon - UPDATE NOW!"
- [ ] **Forced Upgrades**: Block legacy clients after grace period
- [ ] **Custom Messages**: Admin-configurable upgrade messages

### Implementation Details
- [ ] Legacy endpoint wrapper around new API
- [ ] Database migration for existing user data
- [ ] Command format translation between old/new systems
- [ ] Gradual feature deprecation timeline
- [ ] Legacy client usage analytics and monitoring

### Legacy Client Migration Timeline

#### Phase 1: Soft Launch (Weeks 1-4)
- Deploy Go backend with legacy endpoints
- Monitor legacy client usage
- Gentle upgrade notifications: "New version available!"
- Collect migration metrics and user feedback

#### Phase 2: Active Migration (Weeks 5-8)
- Increase notification frequency
- More assertive messages: "Hey! There's a new version out! Fucking update!"
- Provide migration guides and support
- Track adoption rates and address migration blockers

#### Phase 3: Forced Migration (Weeks 9-12)
- Daily persistent notifications
- Feature limitations for legacy clients
- "Legacy support ending soon - UPDATE NOW!"
- Direct user outreach for holdouts

#### Phase 4: Legacy Sunset (Week 13+)
- Block legacy client logins
- Redirect to download page
- Emergency support for critical users only
- Complete migration to new system

### Legacy Client Detection and Messaging

#### Version Detection
```go
// Detect legacy client by exact version parameter matching
if version := r.URL.Query().Get("vrs"); version == "012" {
    // Legacy client detected - handle with exact compatibility
    return handleLegacyClient(w, r, version)
}
```

#### Exact Response Format Examples
```go
// AppCommand.aspx response format (exact match)
fmt.Fprintf(w, "%s\n%s\n%s\n%s", count, nextSender, verified, thumbsUp)

// GetContent.aspx response format (exact match)  
fmt.Fprintf(w, "%s\n%s", senderID, commandContent)

// GetCount.aspx response format (exact match)
fmt.Fprintf(w, "%s", commandCount)
```

#### Upgrade Message Examples
- **Week 1-2**: "📱 New ControlMe client available! Better performance and features await!"
- **Week 3-4**: "⚡ Upgrade to the new client for instant commands (no more 30-second delays!)"
- **Week 5-6**: "🔥 Hey! There's a new version out! Fucking update! WebSocket speed awaits!"
- **Week 7-8**: "⚠️ Legacy support ending soon - UPDATE NOW! Don't get left behind!"
- **Week 9+**: "🚫 Legacy client support has ended. Please download the new client to continue."

### Legacy Feature Limitations
- [ ] Gradual feature restrictions for legacy clients
- [ ] Reduced command priority for legacy users
- [ ] Limited file upload sizes for legacy clients
- [ ] Restricted group features for legacy clients
- [ ] Read-only mode before complete sunset

### In-Place Swap Requirements

#### Exact URL Matching
```go
// Current .NET client expects these exact paths:
GET /AppCommand.aspx?usernm=user&pwd=encrypted&vrs=012&cmd=Outstanding
GET /GetContent.aspx?usernm=user&pwd=encrypted&vrs=012  
GET /GetCount.aspx?usernm=user&pwd=encrypted&vrs=012
POST /ProcessComplete.aspx?usernm=user&pwd=encrypted&vrs=012
```

#### Exact Parameter Handling
- **usernm**: Username parameter (exact spelling)
- **pwd**: Encrypted password using legacy crypto
- **vrs**: Version parameter (currently "012")
- **cmd**: Command type parameter for AppCommand.aspx

#### Exact Response Format Matching
- **AppCommand.aspx**: Plain text response with command count and user info
- **GetContent.aspx**: Plain text with SenderId and command content
- **GetCount.aspx**: Plain text with command count
- **ProcessComplete.aspx**: Plain text success/error response

#### Legacy Crypto Compatibility
- [ ] Decrypt passwords using exact same algorithm as current C# CryptoHelper
- [ ] Maintain exact same encryption/decryption behavior
- [ ] Support legacy password format during transition

#### Database Query Compatibility
- [ ] Exact same stored procedure behavior simulation
- [ ] Same result set format and column order
- [ ] Identical error handling and messages
- [ ] Same database connection string handling
