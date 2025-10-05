# Known Issues and Limitations

## 🚨 **Current Known Issues**

### Configuration and Environment
1. **Port Configuration**
   - **Issue**: Default port 8081 conflicts with other services
   - **Current Workaround**: Using port 8082 in test configuration
   - **Status**: Working solution in place
   - **Solution**: Update default configuration or implement dynamic port detection

2. **JWT Secret in Development**
   - **Issue**: Using test JWT secret "test-jwt-secret-key-for-development"
   - **Security Risk**: Not suitable for production
   - **Status**: Acceptable for development, must change for production
   - **Solution**: Use environment variable for production deployment

3. **Database Configuration Priority**
   - **Issue**: Main config.yaml tries PostgreSQL first, then falls back to SQLite
   - **Impact**: Unnecessary connection attempts and log warnings
   - **Status**: Working but produces noise in logs
   - **Solution**: Set explicit database type in configuration

### Code and Architecture
4. **Single Instance Limitation**
   - **Issue**: Current architecture doesn't support horizontal scaling
   - **Impact**: Limited to single server deployment
   - **Status**: By design for current requirements
   - **Future Solution**: Redis for session storage and pub/sub

5. **In-Memory Session State**
   - **Issue**: WebSocket sessions not persisted across server restarts
   - **Impact**: All connections lost on restart
   - **Status**: By design, acceptable for current use case
   - **Future Solution**: Persistent session storage if needed

6. **Limited Error Context**
   - **Issue**: Some error messages could be more descriptive
   - **Impact**: Slightly harder debugging experience
   - **Status**: Functional but could be improved
   - **Solution**: Enhanced error messages with more context

### Testing and Development
7. **Test Coverage Gaps**
   - **Issue**: No comprehensive automated test suite
   - **Impact**: Manual testing required for changes
   - **Status**: Have working test clients but no unit tests
   - **Solution**: Implement unit and integration test suites

8. **Development Hot Reload**
   - **Issue**: No automatic restart on code changes
   - **Impact**: Manual rebuild and restart required
   - **Status**: Workable but not optimal
   - **Solution**: Implement hot reload with Air or similar tool

## ⚠️ **Technical Debt Items**

### Code Quality
1. **Magic Numbers**: Some timeouts and limits are hardcoded
2. **Error Handling**: Could be more granular and user-friendly
3. **Logging Consistency**: Mix of different logging approaches
4. **Configuration Validation**: Limited validation of config file values

### Documentation
1. **Code Comments**: Some complex functions need better documentation
2. **API Examples**: More comprehensive examples in Swagger docs
3. **Error Response Documentation**: Better documentation of error conditions

### Performance
1. **Database Queries**: Some queries could be optimized with indexes
2. **Memory Usage**: No comprehensive memory profiling done
3. **Connection Cleanup**: Could be more aggressive about cleaning stale connections

## 🔧 **Workarounds in Place**

### Port 8082 Usage
```yaml
# Current working configuration in config.test.yaml
server:
  port: 8082  # Changed from 8081 due to conflicts
```

### SQLite for Development
```yaml
# Explicit SQLite configuration
database:
  type: sqlite
  path: data/controlme.db
```

### Environment Variable Override
```powershell
# Use test configuration
$env:CONFIG_FILE="config.test.yaml"
```

## 🎯 **Non-Issues (Intentional Design)**

### WebSocket-First Architecture
- **Not an Issue**: REST API is minimal by design
- **Reason**: WebSocket is the primary communication method
- **Status**: Working as intended

### Anonymous Session Limitations
- **Not an Issue**: Anonymous users have limited functionality
- **Reason**: Security and business logic requirement
- **Status**: Working as designed

### SQLite in Development
- **Not an Issue**: Using SQLite instead of PostgreSQL in development
- **Reason**: Simpler setup and no external dependencies
- **Status**: Appropriate for development environment

## 🚫 **False Alarms to Ignore**

### Build Warnings
1. **Unused Imports**: May appear during development, cleaned up in final builds
2. **Go Module Checksums**: Occasional checksum warnings are normal
3. **GORM Debug Logs**: Verbose SQL logging in development mode is expected

### Runtime Messages
1. **WebSocket Upgrade Messages**: Connection upgrade logs are informational
2. **JWT Validation Messages**: Token validation logs in debug mode are normal
3. **CORS Preflight**: OPTIONS requests and CORS headers are expected

## 🔍 **Monitoring Points**

### Performance Metrics to Watch
- **Connection Count**: Monitor active WebSocket connections
- **Memory Usage**: Watch for memory leaks in long-running sessions
- **Database Connection Pool**: Ensure proper connection recycling
- **Message Processing Time**: Monitor WebSocket message handling latency

### Error Conditions to Monitor
- **Authentication Failures**: High rate could indicate attack
- **WebSocket Connection Failures**: Could indicate network issues
- **Database Connection Errors**: May indicate resource exhaustion
- **JSON Parsing Errors**: Could indicate malformed client messages

## 🛠️ **Quick Fixes Available**

### Immediate Improvements (< 1 hour)
1. **Better Error Messages**: Add more context to error responses
2. **Configuration Validation**: Add startup validation for config values
3. **Health Check Enhancement**: Add more detailed health check information
4. **Logging Standardization**: Standardize log message format

### Short-term Improvements (< 1 day)
1. **Production JWT Secret**: Environment variable configuration
2. **Rate Limiting**: Basic rate limiting implementation
3. **Enhanced Monitoring**: Add basic metrics endpoints
4. **Input Validation**: More comprehensive input validation

## 📋 **Issue Tracking**

### Priority Levels
- **P0 (Critical)**: Security vulnerabilities, data loss risks
- **P1 (High)**: Functionality breaking issues, poor user experience
- **P2 (Medium)**: Performance issues, minor functionality gaps
- **P3 (Low)**: Cosmetic issues, nice-to-have improvements

### Current Issue Distribution
- **P0**: None currently identified
- **P1**: JWT secret for production (item #2)
- **P2**: Port conflicts (item #1), error messaging (item #6)
- **P3**: Hot reload (item #8), code comments (documentation debt)

## 🎯 **Resolution Strategy**

### For New Development
1. **Always test** both anonymous and authenticated WebSocket flows
2. **Verify** REST API endpoints after any authentication changes
3. **Check** database operations work correctly
4. **Update** documentation when making architectural changes
5. **Consider** security implications of any new features

### For Bug Reports
1. **Reproduce** the issue with current configuration
2. **Check** if it's a known issue from this list
3. **Verify** it's not an intentional design decision
4. **Test** with both SQLite and PostgreSQL if database-related
5. **Document** any new issues discovered

This list should be updated as new issues are discovered or existing ones are resolved.