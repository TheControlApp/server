# WebSocket Quick Reference

## 🚀 **TL;DR**

```bash
# Connect
wscli.exe ws://localhost:8080/ws/client

# Send command
{"instructions":[{"type":"std_popup","content":{"body":"Hello"}}],"tags":"general"}

# Authenticate (optional)
{"type":"auth","token":"your_jwt_token"}
```

## 📨 **Message Types**

### System Messages (Direct Response)
| Send | Receive |
|------|---------|
| `{"type":"ping"}` | `{"type":"pong","timestamp":123}` |
| `{"type":"auth","token":"..."}` | `{"type":"auth_success",...}` |

### Commands (Broadcasted to All)
```json
{
  "instructions": [{"type":"INSTRUCTION_TYPE","content":{...}}],
  "tags": "general"
}
```

## 🎯 **Common Instructions**

| Type | Content Example |
|------|-----------------|
| `std_popup` | `{"body":"Message","title":"Title"}` |
| `std_timer` | `{"duration":300,"title":"5 min timer"}` |
| `std_download` | `{"file_hash":"abc","file_name":"file.pdf"}` |
| `std_notification` | `{"title":"Alert","message":"Text"}` |

## 🔐 **Authentication**

```bash
# Get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"pass"}'

# Use in WebSocket  
{"type":"auth","token":"returned_token"}
```

## ✅ **Working Examples**

**Simple popup:**
```json
{"instructions":[{"type":"std_popup","content":{"body":"Hello World"}}],"tags":"general"}
```

**Timer command:**
```json
{"instructions":[{"type":"std_timer","content":{"duration":60,"title":"1 minute"}}],"tags":"general"}
```

**Multi-instruction:**
```json
{
  "instructions": [
    {"type":"std_popup","content":{"body":"Starting..."}},
    {"type":"std_timer","content":{"duration":30,"title":"Wait"}}
  ],
  "tags":"general"
}
```

**Legacy format (still works):**
```json
{"type":"std_popup","content":{"body":"Hello"}}
```