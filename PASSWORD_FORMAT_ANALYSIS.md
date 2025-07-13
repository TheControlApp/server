# ControlMe Password Format Analysis & Testing

## Password Storage vs Authentication Format

There are **two different password systems** in ControlMe:

### 1. 🔒 Legacy Authentication (for .NET clients)
- **Storage**: Passwords stored **encrypted** in database using AES-256-CBC
- **Authentication**: Client sends the **same encrypted password** from database
- **Algorithm**: AES-256-CBC + PKCS7 padding + Base64 encoding
- **Key**: `dev-crypto-key` (configurable)

### 2. 🆕 Modern Authentication (for new clients)  
- **Storage**: Passwords stored **hashed** using bcrypt
- **Authentication**: Client sends plain password, server hashes and compares
- **Algorithm**: bcrypt with salt

## Legacy Authentication Flow

```
1. User creates account → Password encrypted with AES → Stored in DB
2. Client encrypts password with same AES key → Sends encrypted password
3. Server compares encrypted client password with stored encrypted password
4. If match → Authentication success
```

## Testing Password Encryption

### Test the Crypto Tool
```bash
# Encrypt a password
cd /workspace/server
go run cmd/tools/crypto-test/main.go --password "mypassword123"

# Output example:
# 🔐 Password Encryption Result:
#    Original: mypassword123
#    Encrypted: AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=
#    Crypto Key: dev-crypto-key
```

### Test Decryption
```bash
# Decrypt a password
go run cmd/tools/crypto-test/main.go --decrypt "AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw="

# Output example:
# 🔓 Password Decryption Result:
#    Encrypted: AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=
#    Decrypted: mypassword123
#    Crypto Key: dev-crypto-key
```

## What Your .NET Client Should Do

### 1. Password Encryption (Client Side)
Your .NET client needs to encrypt passwords using **the same AES key**:

```csharp
// C# example (you'll need to match this exactly)
public static string EncryptPassword(string password)
{
    string key = "dev-crypto-key"; // Must match server config
    // Use AES-256-CBC with PKCS7 padding
    // Return base64 encoded result
    // IV should be 16 zero bytes for compatibility
}
```

### 2. Authentication Request
```
GET /TestAuth.aspx?usernm=USERNAME&pwd=ENCRYPTED_PASSWORD&vrs=012
```

Where:
- `usernm` = username (screen_name in database)
- `pwd` = **encrypted password** (NOT plain password)
- `vrs` = "012" (version check)

## Expected Authentication Response Formats

### Test Endpoint (`/TestAuth.aspx`)
**Format:** Plain text with debugging info
```
=== ControlMe Legacy Auth Test ===
Timestamp: 2025-07-13 22:10:59

Step 1: Checking version...
  ✅ Version check passed

Step 2: Decrypting password...
  ✅ Decrypted password: 'mypassword123'

Step 3: Authenticating user...
  ✅ Authentication SUCCESS!
    User ID: 1502755684
    Role: user
    Verified: true
```

### Production Endpoints
**Format:** HTML matching ASP.NET structure

#### GetCount.aspx
```html
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head runat="server">
    <title></title>
</head>
<body>
    <form id="form1" runat="server">
        <div>
            <asp:Label ID="result" runat="server">0</asp:Label>
            <asp:Label ID="next" runat="server"></asp:Label>
            <asp:Label ID="vari" runat="server">true</asp:Label>
        </div>
    </form>
</body>
</html>
```

#### GetContent.aspx
```html
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head runat="server">
    <title></title>
</head>
<body>
    <form id="form1" runat="server">
        <div>
            <asp:Label ID="SenderId" runat="server"></asp:Label>
            <asp:Label ID="Result" runat="server"></asp:Label>
            <asp:Label ID="Varified" runat="server"></asp:Label>
        </div>
    </form>
</body>
</html>
```

#### AppCommand.aspx
```html
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head runat="server">
    <title></title>
</head>
<body>
    <form id="form1" runat="server">
        <div>
            <asp:Label ID="result" runat="server">[]</asp:Label>
        </div>
    </form>
</body>
</html>
```

## Creating Test Users with Encrypted Passwords

### Method 1: Manual Database Insert (Recommended)
```sql
-- Step 1: Encrypt password using Go tool
-- go run cmd/tools/crypto-test/main.go --password "mypassword123"
-- Result: AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=

-- Step 2: Insert user with encrypted password
INSERT INTO users (id, screen_name, login_name, password, role, verified, created_at, updated_at, login_date) 
VALUES (
    gen_random_uuid(), 
    'myuser',           -- This is what client sends as 'usernm'
    'myuser', 
    'AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=', -- Encrypted password
    'user', 
    true, 
    now(), 
    now(), 
    now()
);
```

### Method 2: User Manager Tool (if fixed)
```bash
# This creates bcrypt hashed passwords (for modern auth only)
go run cmd/tools/user-manager/main.go --username myuser --password mypassword123
```

## Troubleshooting "Failed to Decrypt Password"

### Check 1: Password Format
```bash
# Test what your client is sending
curl 'http://localhost:8080/TestAuth.aspx?usernm=yourusername&pwd=YOUR_ENCRYPTED_PASSWORD&vrs=012'

# Look for these in the response:
# ✅ Is valid base64: true
# ✅ Decrypted password: 'yourplainpassword'
```

### Check 2: Crypto Key Match
Your .NET client MUST use the same crypto key as the server:
- **Server key**: Check `configs/config.yaml` → `legacy.crypto_key`
- **Default key**: `dev-crypto-key`

### Check 3: AES Algorithm Match
- **Algorithm**: AES-256-CBC
- **Padding**: PKCS7
- **IV**: 16 zero bytes (for compatibility)
- **Encoding**: Base64

### Common Issues:
1. **"illegal base64"** → Check URL encoding of '=' characters
2. **"cipher: message authentication failed"** → Wrong crypto key
3. **"ciphertext too short"** → Encryption algorithm mismatch
4. **"user not found"** → Wrong username or user doesn't exist

## Working Test Example

```bash
# 1. Create user with encrypted password
docker exec controlme-postgres psql -U postgres -d controlme -c "
INSERT INTO users (id, screen_name, login_name, password, role, verified, created_at, updated_at, login_date) 
VALUES (gen_random_uuid(), 'testuser', 'testuser', 'AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw=', 'user', true, now(), now(), now())
ON CONFLICT (login_name) DO UPDATE SET password = EXCLUDED.password;"

# 2. Test authentication
curl 'http://localhost:8080/TestAuth.aspx?usernm=testuser&pwd=AAAAAAAAAAAAAAAAAAAAACm1LLz77SmR8t3v12QoWhw%3D&vrs=012'

# 3. Should show: ✅ Authentication SUCCESS!
```

## Summary

The key insight is that **legacy authentication doesn't use password hashing** - it uses **password encryption**. Your .NET client needs to encrypt passwords with the same AES key and send the encrypted result, which the server then decrypts and compares with the stored encrypted password.
