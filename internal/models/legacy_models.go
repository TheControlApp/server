package models

import (
	"time"

	"github.com/google/uuid"
)

// Legacy models for database compatibility
// These models map to the exact legacy database schema

// LegacyUser represents the legacy Users table with exact field names
type LegacyUser struct {
	ID           int       `gorm:"column:Id;primaryKey;autoIncrement" json:"id"`
	ScreenName   string    `gorm:"column:Screen Name;size:50;not null" json:"screen_name"`
	LoginName    string    `gorm:"column:Login Name;size:50;not null" json:"login_name"`
	Password     string    `gorm:"column:Password;size:50;not null" json:"password"`
	Role         string    `gorm:"column:Role;size:50" json:"role"`
	RandOpt      bool      `gorm:"column:RandOpt;default:false" json:"rand_opt"`
	AnonCmd      bool      `gorm:"column:AnonCmd;default:false" json:"anon_cmd"`
	Verified     bool      `gorm:"column:Varified;default:false" json:"verified"` // Note: "Varified" is correct typo in legacy
	VerifiedCode int       `gorm:"column:VarifiedCode;default:0" json:"verified_code"`
	LoginDate    time.Time `gorm:"column:LoginDate;default:CURRENT_TIMESTAMP" json:"login_date"`
	ThumbsUp     int       `gorm:"column:ThumbsUp;default:0" json:"thumbs_up"`
}

// TableName returns the legacy table name
func (LegacyUser) TableName() string {
	return "Users"
}

// LegacyControlAppCmd represents the legacy ControlAppCmd table
type LegacyControlAppCmd struct {
	ID         int  `gorm:"column:Id;primaryKey;autoIncrement" json:"id"`
	SenderID   int  `gorm:"column:SenderId;not null" json:"sender_id"`
	SubID      int  `gorm:"column:SubId;not null" json:"sub_id"`
	CommandID  int  `gorm:"column:CmdId;not null" json:"command_id"`
	GroupRefID *int `gorm:"column:GroupRefId" json:"group_ref_id"`
}

// TableName returns the legacy table name
func (LegacyControlAppCmd) TableName() string {
	return "ControlAppCmd"
}

// LegacyCommand represents the legacy CommandList table
type LegacyCommand struct {
	ID       int       `gorm:"column:CmdId;primaryKey;autoIncrement" json:"id"`
	Content  string    `gorm:"column:Content;type:nvarchar(max);not null" json:"content"`
	SendDate time.Time `gorm:"column:SendDate;default:CURRENT_TIMESTAMP" json:"send_date"`
}

// TableName returns the legacy table name
func (LegacyCommand) TableName() string {
	return "CommandList"
}

// LegacyBlock represents the legacy Block table
type LegacyBlock struct {
	ID        int `gorm:"column:Id;primaryKey;autoIncrement" json:"id"`
	BlockerID int `gorm:"column:BlockerId;not null" json:"blocker_id"`
	BlockeeID int `gorm:"column:BlockeeId;not null" json:"blockee_id"`
}

// TableName returns the legacy table name
func (LegacyBlock) TableName() string {
	return "Block"
}

// LegacyInvite represents the legacy Invites table
type LegacyInvite struct {
	ID    int `gorm:"column:Id;primaryKey;autoIncrement" json:"id"`
	SubID int `gorm:"column:SubId;not null" json:"sub_id"`
	DomID int `gorm:"column:DomId;not null" json:"dom_id"`
}

// TableName returns the legacy table name
func (LegacyInvite) TableName() string {
	return "Invites"
}

// LegacyRelationship represents the legacy Relationship table
type LegacyRelationship struct {
	ID    int `gorm:"column:Id;primaryKey;autoIncrement" json:"id"`
	DomID int `gorm:"column:DomId;not null" json:"dom_id"`
	SubID int `gorm:"column:SubId;not null" json:"sub_id"`
}

// TableName returns the legacy table name
func (LegacyRelationship) TableName() string {
	return "Relationship"
}

// UserIDMapping maps between legacy integer IDs and modern UUIDs
type UserIDMapping struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LegacyID int       `gorm:"not null;unique" json:"legacy_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
}

// CommandIDMapping maps between legacy integer IDs and modern UUIDs
type CommandIDMapping struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LegacyID  int       `gorm:"not null;unique" json:"legacy_id"`
	CommandID uuid.UUID `gorm:"type:uuid;not null" json:"command_id"`
}

// ControlAppCmdIDMapping maps between legacy integer IDs and modern UUIDs
type ControlAppCmdIDMapping struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LegacyID        int       `gorm:"not null;unique" json:"legacy_id"`
	ControlAppCmdID uuid.UUID `gorm:"type:uuid;not null" json:"control_app_cmd_id"`
}

// LegacyCompatUser is a hybrid model for legacy compatibility
type LegacyCompatUser struct {
	// Modern fields
	ID           uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ScreenName   string    `gorm:"size:50;not null" json:"screen_name"`
	LoginName    string    `gorm:"size:50;not null;unique" json:"login_name"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:50" json:"role"`
	RandOpt      bool      `gorm:"default:false" json:"rand_opt"`
	AnonCmd      bool      `gorm:"default:false" json:"anon_cmd"`
	Verified     bool      `gorm:"default:false" json:"verified"`
	VerifiedCode int       `gorm:"default:0" json:"verified_code"`
	ThumbsUp     int       `gorm:"default:0" json:"thumbs_up"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LoginDate    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"login_date"`

	// Legacy compatibility fields
	LegacyID         int    `gorm:"unique" json:"legacy_id"`
	LegacyScreenName string `gorm:"column:legacy_screen_name;size:50" json:"legacy_screen_name"`
	LegacyLoginName  string `gorm:"column:legacy_login_name;size:50" json:"legacy_login_name"`
}
