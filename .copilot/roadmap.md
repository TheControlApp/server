# API Roadmap

## Overview
- Primary endpoints:
  - `/api/ws` (WebSocket communication)
  - `/api/v1/login` (User login)
  - `/api/v1/register` (User registration)
- Secondary endpoints will be defined based on feature requirements and use cases.

## Current Implementation Review
- WebSocket routes are currently exposed at `/ws/client` and `/ws/web`.
- Authentication routes are under `/api/v1/auth/login` and `/api/v1/auth/register`.

## Proposed Plan

1. Align WebSocket endpoint to `/api/ws`
   - Consolidate existing WebSocket handlers under the new path.
   - Define namespace, connection handshake, and message formats.

2. Refactor authentication endpoints
   - Expose login at `/api/v1/login` (instead of `/api/v1/auth/login`).
   - Expose registration at `/api/v1/register` (instead of `/api/v1/auth/register`).
   - Decide on backward compatibility or deprecation of old paths.

3. Define request and response schemas
   - **Login**: JSON payload (email, password), response (JWT token structure), error formats.
   - **Register**: JSON payload (user fields), validation rules, response (user data, tokens).
   - **WebSocket**: handshake parameters (JWT in query or header), message DTOs, error channels.

4. Plan secondary endpoints
   - User profile management (GET/PUT `/api/v1/users/profile`).
   - Command operations (pending, create, complete) under `/api/v1/commands`.
   - Additional resources as needed.

5. Timeline & milestones
   - Phase 1: Route alignment and basic auth endpoints (2 weeks).
   - Phase 2: Detailed schema and validation (1 week).
   - Phase 3: WebSocket enhancements (2 weeks).
   - Phase 4: Secondary endpoints and documentation (1 week).

## Clarifications Received
- **Legacy paths**: Existing `/api/v1/auth/*` paths are for reference only and will eventually be removed.
- **WebSocket auth**: Session-token based authentication. Client logs in via `/api/v1/login` endpoint first, then WebSocket authenticates using the session-token.
- **Command types**: Building custom command types but following similar intentions to legacy codebase patterns.

## Legacy Command Analysis
Based on the legacy codebase review, the original system had these command types:
- **Direct commands**: Text/action commands sent between users
- **System commands**: Special commands like "Outstanding", "Content", "Delete"
- **Feedback commands**: "Thumbs" for rating system

### Current Go Models
The new system has these models:
- `Command`: Type, Content, Data, Status fields
- `ControlAppCmd`: Assignment linking Sender -> Recipient with Command
- `User`: Basic user model with verification, thumbs up system
- `ChatLog`: Direct messaging between users

## Updated Plan

### Phase 1: Core API Alignment (Week 1-2)
1. **Route structure**:
   - Keep existing `/api/v1/auth/login` and `/api/v1/auth/register`
   - Consolidate all WebSocket functionality under single `/api/ws` endpoint
   - Remove separate `/ws/client` and `/ws/web` endpoints

2. **Session-token authentication**:
   - Implement JWT session tokens in login response
   - WebSocket handshake accepts session-token (query param or header)
   - Single WebSocket handles all real-time communication

### Phase 2: WebSocket-First Architecture (Week 3-4)
1. **Primary communication via WebSocket**:
   - Command creation, assignment, completion via WebSocket messages
   - Real-time chat and notifications
   - User status and presence updates
   - Most application logic flows through WebSocket

2. **Minimal REST API for specific needs**:
   - `/api/v1/commands/*` - Only for file uploads/downloads or external integrations
   - `/api/v1/users/profile` - Profile updates (may also support via WebSocket)
   - `/api/v1/health` - System health checks

### Phase 3: WebSocket Message Types & Commands (Week 5-6)
1. **WebSocket message types**:
   - `command.create` - Create new command
   - `command.assign` - Assign command to user
   - `command.complete` - Mark command complete
   - `command.list` - Get pending commands
   - `chat.message` - Send chat message
   - `user.status` - Status updates
   - `system.notification` - System notifications

2. **Command types via WebSocket**:
   - Message commands (basic text)
   - Task commands (actionable items)
   - System commands (outstanding, status)
   - User-to-user command distribution

### Phase 4: Additional Features (Week 7)
1. **User management**: Profile updates, verification system
2. **Rating system**: Thumbs up/down for commands and users
3. **Admin features**: Command logging, user management

## Next Steps
- Begin Phase 1: Consolidate WebSocket endpoints to single `/api/ws`
- Keep existing auth endpoints at `/api/v1/auth/login` and `/api/v1/auth/register`
- Define WebSocket message format and authentication flow
- Plan WebSocket-first architecture with minimal REST API usage
