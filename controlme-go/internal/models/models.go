package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
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
}

// Command represents a command in the system
type Command struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Content   string    `gorm:"type:text" json:"content"`
	Data      string    `gorm:"type:text" json:"data"`
	Status    string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ControlAppCmd represents the command assignment table
type ControlAppCmd struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SenderID   uuid.UUID  `gorm:"type:uuid;not null" json:"sender_id"`
	SubID      uuid.UUID  `gorm:"type:uuid;not null" json:"sub_id"`
	CommandID  uuid.UUID  `gorm:"type:uuid;not null" json:"command_id"`
	GroupRefID *uuid.UUID `gorm:"type:uuid" json:"group_ref_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`

	// Relationships
	Sender  User    `gorm:"foreignKey:SenderID" json:"sender"`
	Sub     User    `gorm:"foreignKey:SubID" json:"sub"`
	Command Command `gorm:"foreignKey:CommandID" json:"command"`
}

// ChatLog represents chat messages
type ChatLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID uuid.UUID `gorm:"type:uuid;not null" json:"receiver_id"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Sender   User `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver"`
}

// Group represents user groups
type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Owner User `gorm:"foreignKey:OwnerID" json:"owner"`
}

// GroupMember represents group membership
type GroupMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Role      string    `gorm:"size:50;default:'member'" json:"role"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Group Group `gorm:"foreignKey:GroupID" json:"group"`
	User  User  `gorm:"foreignKey:UserID" json:"user"`
}

// Relationship represents dom/sub relationships
type Relationship struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DomID     uuid.UUID `gorm:"type:uuid;not null" json:"dom_id"`
	SubID     uuid.UUID `gorm:"type:uuid;not null" json:"sub_id"`
	Status    string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Dom User `gorm:"foreignKey:DomID" json:"dom"`
	Sub User `gorm:"foreignKey:SubID" json:"sub"`
}

// Block represents blocked users
type Block struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BlockerID uuid.UUID `gorm:"type:uuid;not null" json:"blocker_id"`
	BlockeeID uuid.UUID `gorm:"type:uuid;not null" json:"blockee_id"`
	Reason    string    `gorm:"type:text" json:"reason"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Blocker User `gorm:"foreignKey:BlockerID" json:"blocker"`
	Blockee User `gorm:"foreignKey:BlockeeID" json:"blockee"`
}

// Report represents user reports
type Report struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ReporterID uuid.UUID `gorm:"type:uuid;not null" json:"reporter_id"`
	ReportedID uuid.UUID `gorm:"type:uuid;not null" json:"reported_id"`
	Reason     string    `gorm:"type:text;not null" json:"reason"`
	Status     string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	Reporter User `gorm:"foreignKey:ReporterID" json:"reporter"`
	Reported User `gorm:"foreignKey:ReportedID" json:"reported"`
}

// Invite represents user invitations
type Invite struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID uuid.UUID `gorm:"type:uuid;not null" json:"receiver_id"`
	Type       string    `gorm:"size:50;not null" json:"type"`
	Status     string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relationships
	Sender   User `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver"`
}

// BeforeCreate hook to set UUID for models
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (c *Command) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (cac *ControlAppCmd) BeforeCreate(tx *gorm.DB) error {
	if cac.ID == uuid.Nil {
		cac.ID = uuid.New()
	}
	return nil
}
