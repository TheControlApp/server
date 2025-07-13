# ControlMe Legacy Authentication Test

This document explains how to test the legacy authentication system to ensure it works correctly with legacy .NET clients.

## Quick Test

**✅ WORKING AUTHENTICATION TEST:**

```bash
# 1. Encrypt a password
go run cmd/tools/crypto-test/main.go --password "test123"
# Output: AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=

# 2. Test authentication
curl 'http://localhost:8080/TestAuth.aspx?usernm=testuser&pwd=AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=&vrs=012'
```

## Authentication Test Endpoint

The `/TestAuth.aspx` endpoint provides comprehensive debugging information for legacy authentication:

### URL Format
```
http://localhost:8080/TestAuth.aspx?usernm=USERNAME&pwd=ENCRYPTED_PASSWORD&vrs=012
```

### Parameters
- `usernm`: Username (screen_name in database)
- `pwd`: Encrypted password using legacy crypto key
- `vrs`: Version (must be "012" for legacy compatibility)

### Response Example

```
=== ControlMe Legacy Auth Test ===
Timestamp: 2025-07-13 22:02:12

Received Parameters:
  usernm: 'testuser'
  pwd: 'AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=' (encrypted, length: 44)
  vrs: '012'

Server Configuration:
  Crypto Key: 'dev-crypto-key'
  Expected Version: '012'

Step 1: Checking version...
  ✅ Version check passed

Step 2: Decrypting password...
  ✅ Decrypted password: 'test123'

Step 3: Authenticating user...
  ✅ Authentication SUCCESS!
    User ID: 1502755684
    Role: user
    Verified: true

Step 4: What GetCount.aspx would return:
  ✅ User is verified - GetCount would check for pending commands
    Expected response format: [integer] (e.g., '0', '1', '5')

=== Test Complete ===
🎉 If you see ✅ for steps 1-3, your authentication is working correctly!
```

## Creating Test Users

### Method 1: Direct Database Insert
```sql
-- Insert user with encrypted password
INSERT INTO users (id, screen_name, login_name, password, role, verified, created_at, updated_at, login_date) 
VALUES (
    gen_random_uuid(), 
    'testuser', 
    'testuser', 
    'AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=', -- encrypted "test123"
    'user', 
    true, 
    now(), 
    now(), 
    now()
);
```

### Method 2: Using Go Tool (if working)
```bash
# Create table first
go run cmd/tools/user-manager/main.go --create-table

# Create user (may have GORM issues)
go run cmd/tools/user-manager/main.go --username testuser --password test123
```

## Legacy Endpoints Working

### Test Authentication
```bash
curl 'http://localhost:8080/TestAuth.aspx?usernm=testuser&pwd=ENCRYPTED_PASSWORD&vrs=012'
```

### Get Content (Messages)
```bash
curl 'http://localhost:8080/GetContent.aspx?usernm=testuser&pwd=ENCRYPTED_PASSWORD&vrs=012'
```
**Response:** HTML with SenderId, Result, and Varified labels (exactly like ASP.NET)

### Get Count (Pending Commands)
```bash
curl 'http://localhost:8080/GetCount.aspx?usernm=testuser&pwd=ENCRYPTED_PASSWORD&vrs=012'
```
**Response:** HTML with result, next, and vari labels

### App Command
```bash
curl 'http://localhost:8080/AppCommand.aspx?usernm=testuser&pwd=ENCRYPTED_PASSWORD&vrs=012&cmd=Outstanding'
```
**Response:** HTML with result label containing command data

## Expected Output for Legacy Clients

All legacy endpoints return **HTML responses** that match the exact structure of the original ASP.NET pages:

```html
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head runat="server">
    <title></title>
</head>
<body>
    <form id="form1" runat="server">
        <div>
            <asp:Label ID="result" runat="server">DATA_HERE</asp:Label>
        </div>
    </form>
</body>
</html>
```

Legacy .NET clients expect to parse these HTML responses and extract the text content from the labels.

## Troubleshooting

### Authentication Fails
1. **Check password encryption:**
   ```bash
   go run cmd/tools/crypto-test/main.go --password "yourpassword"
   ```

2. **Verify user exists:**
   ```sql
   SELECT screen_name, login_name, verified FROM users WHERE screen_name = 'yourusername';
   ```

3. **Check crypto key in config:**
   - Default: `dev-crypto-key`
   - Config file: `configs/config.yaml` → `legacy.crypto_key`

### Database Issues
- Legacy service expects some tables that may not exist (like "Users" with capital U)
- Main authentication works against `users` table
- Some legacy stored procedures may need table mapping updates

## Summary

✅ **WORKING:**
- Password encryption/decryption
- User authentication
- Legacy endpoint structure
- HTML response format matching ASP.NET

⚠️ **PARTIALLY WORKING:**
- Some legacy stored procedures need table name updates
- GORM auto-migration has issues with complex User model

🎯 **READY FOR LEGACY CLIENTS:**
The authentication system is working correctly and returns the expected HTML format that legacy .NET clients can parse.
