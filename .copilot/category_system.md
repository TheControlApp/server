# Content Category System Design

## Purpose
Allow users to filter incoming commands based on content categories, giving users control over what types of content they receive.

## Category Examples
- `general` - General, safe-for-work content
- `censored` - Adult content that's been censored/blurred
- `adult` - Explicit adult content  
- `feet` - Foot-related content
- `extreme` - Extreme or intense content
- `roleplay` - Roleplay scenarios
- `humiliation` - Humiliation-based content
- `public` - Public setting commands
- `private` - Private/intimate commands

## User Control Flow

### Setting Preferences
1. User sets blocked/allowed categories via WebSocket or REST API
2. Server stores preferences in `user_preferences` table
3. All future commands filtered based on these preferences

### Command Delivery Logic
1. Command created with one or more categories
2. Server checks recipient's preferences before delivery
3. If any command category is in user's blocked list → command not delivered
4. If command has no categories → defaults to "general" category
5. Commands only delivered if user allows the categories

### Database Schema
```sql
-- Commands table with categories
CREATE TABLE commands (
    id UUID PRIMARY KEY,
    sender_id UUID NOT NULL,
    content TEXT NOT NULL,
    command_type VARCHAR(50) NOT NULL,
    categories TEXT[], -- Array of category strings
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '2 weeks'
);

-- User preferences for category filtering
CREATE TABLE user_preferences (
    user_id UUID PRIMARY KEY,
    blocked_categories TEXT[] DEFAULT '{}',
    allowed_categories TEXT[] DEFAULT '{"general"}',
    updated_at TIMESTAMP DEFAULT NOW()
);
```

## WebSocket Messages

### Set User Preferences
```json
{
  "type": "preferences.set",
  "payload": {
    "blockedCategories": ["feet", "extreme"],
    "allowedCategories": ["general", "censored", "roleplay"]
  }
}
```

### Create Categorized Command
```json
{
  "type": "command.create",
  "payload": {
    "content": "Strike a pose",
    "categories": ["general", "roleplay"],
    "targetUser": "user123",
    "commandType": "task"
  }
}
```

## Implementation Notes
- Default category is "general" if none specified
- Users start with only "general" allowed by default
- Server maintains master list of valid categories
- Categories are additive (command blocked if ANY category is blocked)
- All-cast commands still respect individual user preferences
- Admin users may have override capabilities for moderation

## Benefits
- User agency over content they receive
- Content creator awareness of audience preferences  
- Reduced unwanted content delivery
- Scalable filtering system
- Maintains user engagement by respecting boundaries
