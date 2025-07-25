# ControlMe Go Backend - Development Roadmap

## 🎯 Project Vision
A modern, secure, and scalable rewrite of the ControlMe platform in Go, featuring a WebSocket-first architecture for real-time command and communication systems with user-controlled content filtering.

## 🏗️ Architecture Overview

### Core Design Principles
- **WebSocket-First**: Primary communication through single `/api/ws` endpoint
- **User Agency**: Content category filtering puts users in control
- **Real-time Priority**: Commands and notifications delivered instantly
- **Secure by Design**: CSAM scanning, virus detection, and robust authentication
- **Offline Resilience**: Persistent command queues for offline users

### Technology Stack
- **Backend**: Go 1.21+ with Gin framework
- **Database**: PostgreSQL with GORM ORM
- **Real-time**: WebSocket with message hub
- **Authentication**: JWT session tokens (24-hour validity)
- **File Security**: CSAM + virus scanning on uploads
- **Deployment**: Docker Compose

## 📋 Development Phases

### 🔧 Phase 1: Core Infrastructure (Weeks 1-2) - ✅ FOUNDATION COMPLETE
**Goal**: Establish WebSocket-first foundation with authentication

**Infrastructure Completed** ✅:
- [x] **Server Running**: Port 8080 with all endpoints functional
- [x] **Database Migration**: All 4 core tables successfully created
- [x] **REST API**: Authentication and user management endpoints working
- [x] **WebSocket Foundation**: Basic endpoints available at `/ws/client` and `/ws/web`
- [x] **Critical Bug Fixed**: GORM migration error resolved

**Current Implementation Status**:
- [x] Session token authentication flow (REST endpoints ready)
- [ ] Consolidate WebSocket endpoints to single `/api/ws` 
- [ ] Implement WebSocket authentication with JWT tokens
- [ ] Build command queue storage and delivery system
- [ ] Add heartbeat mechanism (3-miss disconnection)
- [ ] Set up connection management with error logging

**Key Features Ready for Implementation**:
- Database schema with UUID primary keys established
- JWT authentication infrastructure in place
- WebSocket endpoints available for real-time implementation
- Docker environment stable and tested

**Next Immediate Steps**: 
1. Implement WebSocket protocol at `/api/ws`
2. Add session token validation for WebSocket connections
3. Build command queue service with database persistence
4. Implement connection management and heartbeat system

---

### 📁 Phase 2: File Management System (Weeks 3-4)
**Goal**: Secure file handling with deduplication and scanning

**Deliverables**:
- [ ] Hash-based file storage with deduplication
- [ ] CSAM scanning integration (PhotoDNA/AWS Rekognition)
- [ ] Virus scanning setup (ClamAV or cloud service)
- [ ] File metadata tracking and retention policies
- [ ] RESTful file upload/download APIs

**Key Features**:
- Single copy per file hash, multiple filename mappings
- Comprehensive security scanning pipeline
- 2-week retention policy for command files
- Efficient storage structure: `/storage/files/{hash_prefix}/{full_hash}`

---

### 🎮 Phase 3: Command & Content System (Weeks 5-6)
**Goal**: Full command system with category-based filtering

**Deliverables**:
- [ ] WebSocket command creation and assignment
- [ ] Content category system (general, adult, feet, etc.)
- [ ] User preference management (block/allow categories)
- [ ] Category-based delivery filtering
- [ ] Broadcast functionality (user/category/all-cast)
- [ ] Command completion tracking
- [ ] Queue delivery for reconnecting users

**Key Features**:
- User-controlled content filtering
- Multiple broadcast types with permission controls
- Persistent command history (2-week retention)
- Real-time command status updates

---

### ✨ Phase 4: Polish & Optional Features (Week 7)
**Goal**: Enhanced user experience and administrative tools

**Deliverables**:
- [ ] Optional chat system implementation
- [ ] Advanced moderation tools
- [ ] Admin dashboard for user/category management
- [ ] Performance optimization and caching
- [ ] Comprehensive logging and monitoring

**Key Features**:
- Real-time chat with file attachments (if implemented)
- Administrative oversight tools
- Performance monitoring and optimization
- Advanced moderation capabilities

## 🗂️ Content Category System

### Purpose
Give users complete control over what types of content they receive through a flexible category-based filtering system.

### Category Examples
- `general` - Safe-for-work content
- `censored` - Blurred/censored adult content
- `adult` - Explicit content
- `feet` - Foot-related content
- `extreme` - Intense content
- `roleplay` - Roleplay scenarios
- `humiliation` - Humiliation-based content
- `public` - Public setting commands
- `private` - Private/intimate commands

### User Control
Users can block/allow specific categories, and the server filters all incoming commands based on these preferences before delivery.

## 🔌 API Design

### REST Endpoints (Minimal)
```
POST /api/v1/auth/login      # Authentication
POST /api/v1/auth/register   # User registration
GET/PUT /api/v1/users/profile # Profile management
POST /api/v1/files/upload    # File uploads
GET /api/v1/files/download/{hash}/{filename} # File downloads
GET /api/v1/health          # Health checks
```

### WebSocket Messages
Primary communication through `/api/ws` with message types:
- `command.*` - Command operations (create, complete, queue)
- `chat.*` - Chat functionality (optional)
- `preferences.*` - User preference management
- `user.*` - Status updates
- `system.*` - Notifications and heartbeat

## 🗄️ Database Schema

### Core Tables
- **users** - User accounts and authentication
- **commands** - Command storage with categories and metadata
- **command_assignments** - Command-to-user assignments with delivery tracking
- **user_preferences** - Category filtering preferences
- **file_metadata** - File information and scan results
- **command_files** - File attachments to commands
- **user_groups** - Group membership for broadcast permissions
- **chat_messages** - Optional chat storage

## 🛡️ Security Features

### File Security
- **CSAM Detection**: Automated scanning for illegal content
- **Virus Scanning**: Real-time virus detection on uploads
- **Hash Deduplication**: Prevents storage waste and duplicate scanning
- **Retention Policies**: Automatic cleanup after 2 weeks

### Connection Security
- **Session Tokens**: 24-hour JWT tokens with automatic expiration
- **Heartbeat Monitoring**: Automatic disconnection of inactive connections
- **Error Logging**: Comprehensive logging for security monitoring
- **Content Filtering**: User-controlled category blocking

## 🎯 Success Metrics

### Phase 1 Success ✅ FOUNDATION ACHIEVED
- [x] Database infrastructure with 4 core tables migrated successfully
- [x] Server running with all REST endpoints functional  
- [x] Authentication system ready with JWT token support
- [x] WebSocket endpoints available for real-time implementation
- [ ] Single WebSocket endpoint handling all connections (ready to implement)
- [ ] Sub-second authentication and connection establishment (ready to implement)
- [ ] Zero message loss during normal disconnections (ready to implement)
- [ ] Successful queue delivery for offline users (ready to implement)

### Phase 2 Success
- [ ] 100% file scanning coverage (CSAM + virus)
- [ ] Effective deduplication (storage savings measured)
- [ ] Sub-5-second file upload and processing
- [ ] Zero false positives in content scanning

### Phase 3 Success
- [ ] User preference system with immediate effect
- [ ] Category filtering working across all broadcast types
- [ ] Real-time command delivery and status updates
- [ ] Successful offline queue delivery

### Phase 4 Success
- [ ] Admin tools fully functional
- [ ] Performance optimizations measurably effective
- [ ] Comprehensive monitoring and alerting
- [ ] Optional features stable and tested

## 🚀 Getting Started

### Prerequisites
- Go 1.21+
- PostgreSQL 13+
- Docker & Docker Compose
- External scanning services configured

### Quick Start
1. Clone repository and review architecture documents in `.copilot/`
2. Set up development environment with Docker
3. Configure external scanning services
4. Begin Phase 1 implementation with WebSocket consolidation
5. Follow phase-by-phase development plan

### Key Resources
- Complete Architecture Guide: `.copilot/complete_architecture_guide.md`
- WebSocket Architecture: `.copilot/websocket_architecture.md`
- Category System Design: `.copilot/category_system.md`
- Implementation Plan: `.copilot/implementation_plan.md`

---

**Current Status**: Phase 1 Foundation Complete ✅ - Ready for WebSocket Implementation

**Next Steps**: Begin WebSocket protocol implementation at `/api/ws` with session token authentication and command queue system.

**Timeline**: Foundation established - remaining 6-week development cycle for WebSocket and feature implementation.

**Support**: All architectural decisions documented in `.copilot/` directory for reference.
