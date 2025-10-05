# .copilot Knowledge Base - Agent Bootstrap Guide

## 🎯 **Quick Context for New Agents**

**Project**: ControlMe Go Server - Real-time command delivery platform  
**Status**: Production-ready with WebSocket-first architecture  
**Last Updated**: September 28, 2025  
**Previous Agent**: GitHub Copilot

## ⚡ **Immediate Context**

### Current State
- ✅ **Server Running**: `http://localhost:8082` with SQLite database
- ✅ **WebSocket Endpoint**: `ws://localhost:8082/ws/client` 
- ✅ **Authentication**: Progressive auth (anonymous → authenticated)
- ✅ **Documentation**: Complete API reference and implementation guides
- ✅ **Test Tools**: Working WebSocket test clients

### User Preferences
- Values **thorough documentation** and **complete implementations**
- Prefers **working code** with proper testing over quick fixes
- Expects **security considerations** in all implementations
- Appreciates **clean, well-organized** code structure
- Wants **flexible, extensible** architectures

## 📚 **Essential Reading Order**

**For Immediate Understanding**:
1. `project-state/current-status.md` - What's working now
2. `session-logs/2024-01-conversation-summary.md` - Our journey
3. `architecture/system-overview.md` - Technical architecture

**For Development Work**:
4. `development/environment-setup.md` - How to run and test
5. `project-state/working-features.md` - Verified functionality
6. `development/future-roadmap.md` - What's planned next

## 🔑 **Key Implementation Details**

### WebSocket Authentication Architecture
The user specifically requested a **progressive authentication flow**:
- **Anonymous connections**: Receive broadcasts only
- **Authenticated connections**: Full command access  
- **Progressive upgrade**: Connect anonymous → send auth message → upgrade session
- **Multiple auth methods**: Header, query parameter, WebSocket message

### Current Working Configuration
```yaml
# config.test.yaml (CURRENTLY ACTIVE)
environment: development
server:
  port: 8082
database:
  type: sqlite
  path: data/controlme.db
auth:
  jwt_secret: "test-jwt-secret-key-for-development"
```

### Message Format Standard
```json
{
  "type": "message_type",
  "payload": { "data": "here" },
  "timestamp": "2024-01-01T00:00:00.000Z"
}
```

## 🚀 **How to Continue Development**

### Starting the Server
```powershell
cd D:\Workspace\github.com\TheControlApp\server
$env:CONFIG_FILE="config.test.yaml"
.\tmp\server.exe
```

### Testing WebSocket Connection
```powershell
# Use the test client
.\bin\test-websocket-auth.exe

# Or manually with wscat
wscat -c ws://localhost:8082/ws/client
```

### Key Files to Know
- `internal/api/handlers/websocket_handlers.go` - WebSocket connection handling
- `internal/websocket/hub.go` - Connection management
- `internal/api/handlers/auth_handlers.go` - REST authentication
- `config.test.yaml` - Current active configuration

## 🛠️ **Common Tasks**

### If Server Won't Start
1. Check if port 8082 is in use: `netstat -ano | findstr :8082`
2. Verify config file: `$env:CONFIG_FILE="config.test.yaml"`
3. Check database file exists: `data/controlme.db`
4. Rebuild if needed: `go build -o tmp/server.exe cmd/server/main.go`

### If WebSocket Issues
1. Check server logs for WebSocket connections
2. Test with browser console: `new WebSocket('ws://localhost:8082/ws/client')`
3. Use test client: `.\bin\test-websocket-auth.exe`
4. Verify message format matches standard

### Adding New Features
1. Follow existing patterns in `internal/` directory
2. Add proper error handling and logging
3. Update documentation in `docs/` 
4. Create or update test clients as needed
5. Follow the user's preference for thorough implementation

## 📋 **Development Standards**

### Code Quality
- Use proper Go conventions and formatting
- Include comprehensive error handling
- Add structured logging for debugging
- Follow existing architectural patterns
- Write tests for new functionality

### Documentation
- Update relevant documentation files
- Include code examples and usage patterns
- Document API changes in Swagger annotations
- Update `.copilot/` notes for future agents
- Maintain consistency with existing docs

### Testing
- Test all WebSocket message types
- Verify authentication flows work
- Check database operations function correctly
- Test error conditions and edge cases
- Use provided test clients for verification

## 🚨 **Critical Things to Remember**

1. **Progressive Authentication**: This was the user's specific requirement - don't break it
2. **WebSocket-First**: Primary communication method, REST is for auth only
3. **SQLite for Development**: Don't change to PostgreSQL without user request
4. **Port 8082**: Currently using this port due to conflicts with 8081
5. **Test Configuration**: Using `config.test.yaml` not `config.yaml`
6. **Message Format**: Standardized JSON format must be maintained
7. **Documentation**: User values comprehensive documentation - update it

## 📞 **When You Need Help**

### Check These First
- Server logs (console output)
- Configuration files (`config.test.yaml`)
- Test client output
- Swagger UI at `http://localhost:8082/swagger/index.html`
- Database file existence and permissions

### Common Solutions
- **Port conflicts**: Change port in config or kill conflicting process
- **Database issues**: Check file permissions and path
- **Authentication errors**: Verify JWT secret and token format
- **WebSocket problems**: Check message format and connection handling
- **Build errors**: Run `go mod tidy` and rebuild

## 🎉 **Success Indicators**

You'll know everything is working when:
- ✅ Server starts without errors on port 8082
- ✅ Swagger UI loads at `/swagger/index.html`
- ✅ WebSocket connections work (test with client)
- ✅ Anonymous and authenticated sessions both function
- ✅ Progressive authentication upgrades work
- ✅ REST API endpoints return proper responses
- ✅ Database operations complete successfully

## 💡 **Pro Tips for This Project**

1. **Always test WebSocket functionality** - it's the core feature
2. **Use the provided test clients** - they're comprehensive and reliable  
3. **Check both anonymous and authenticated flows** - user requires both
4. **Update documentation as you go** - user values this highly
5. **Follow established patterns** - consistency is important
6. **Consider security implications** - user mentioned security frequently
7. **Test error conditions** - robust error handling is expected

This knowledge base contains everything needed to continue development effectively. The user built a solid foundation with comprehensive documentation, working code, and clear architectural decisions. Any future agent should be able to pick up where we left off and continue building upon this strong foundation.

## 📁 **File Navigation**

**Quick Access to Key Information**:
- 📊 Current Status: `project-state/current-status.md`
- 📝 Session History: `session-logs/2024-01-conversation-summary.md`  
- 🏗️ Architecture: `architecture/system-overview.md`
- 🚀 Setup Guide: `development/environment-setup.md`
- 🎯 Future Plans: `development/future-roadmap.md`
- ✅ Working Features: `project-state/working-features.md`

**Main Project Documentation**:
- 📖 Complete API Reference: `../../docs/COMPLETE_API_REFERENCE.md`
- 🔌 WebSocket Implementation: `../../docs/WEBSOCKET_IMPLEMENTATION.md`
- 📋 Swagger Documentation: `http://localhost:8082/swagger/index.html`