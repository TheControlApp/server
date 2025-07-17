package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ScreenName   string    `gorm:"size:50;not null" json:"screen_name"`
	LoginName    string    `gorm:"size:50;not null;unique" json:"login_name"`
	Email        string    `gorm:"size:300;not null;unique" json:"email"`
	Password     string    `gorm:"size:255;not null" json:"-"`
	Role         string    `gorm:"size:50" json:"role"`
	RandomOptIn  bool      `gorm:"default:false" json:"random_opt_in"`
	AnonCmd      bool      `gorm:"default:false" json:"anon_cmd"`
	Verified     bool      `gorm:"default:false" json:"verified"`
	VerifiedCode int       `gorm:"default:0" json:"verified_code"`
	ThumbsUp     int       `gorm:"default:0" json:"thumbs_up"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	LoginDate    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"login_date"`
}

// BeforeCreate sets the ID before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// Command represents a command in the system
type Command struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Content   string    `gorm:"type:text" json:"content"`
	Data      string    `gorm:"type:text" json:"data"`
	Status    string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate sets the ID before creating a command
func (c *Command) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// ControlAppCmd represents the command assignment table
type ControlAppCmd struct {
	ID         uuid.UUID  `gorm:"type:uuid;primary_key" json:"id"`
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

// BeforeCreate sets the ID before creating a control app command
func (cac *ControlAppCmd) BeforeCreate(tx *gorm.DB) error {
	if cac.ID == uuid.Nil {
		cac.ID = uuid.New()
	}
	return nil
}

// ChatLog represents chat messages
type ChatLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	SenderID   uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`
	ReceiverID uuid.UUID `gorm:"type:uuid;not null" json:"receiver_id"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Sender   User `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver User `gorm:"foreignKey:ReceiverID" json:"receiver"`
}

// BeforeCreate sets the ID before creating a chat log
func (cl *ChatLog) BeforeCreate(tx *gorm.DB) error {
	if cl.ID == uuid.Nil {
		cl.ID = uuid.New()
	}
	return nil
}

// Group represents user groups
type Group struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	OwnerID     uuid.UUID `gorm:"type:uuid;not null" json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	Owner User `gorm:"foreignKey:OwnerID" json:"owner"`
}

// BeforeCreate sets the ID before creating a group
func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// GroupMember represents group membership
type GroupMember struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	GroupID   uuid.UUID `gorm:"type:uuid;not null" json:"group_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	Role      string    `gorm:"size:50;default:'member'" json:"role"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Group Group `gorm:"foreignKey:GroupID" json:"group"`
	User  User  `gorm:"foreignKey:UserID" json:"user"`
}

// BeforeCreate sets the ID before creating a group member
func (gm *GroupMember) BeforeCreate(tx *gorm.DB) error {
	if gm.ID == uuid.Nil {
		gm.ID = uuid.New()
	}
	return nil
}

// Relationship represents relationships between users
type Relationship struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	RelatedID uuid.UUID `gorm:"type:uuid;not null" json:"related_id"`
	Type      string    `gorm:"size:50;not null" json:"type"`
	Status    string    `gorm:"size:20;default:'active'" json:"status"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User    User `gorm:"foreignKey:UserID" json:"user"`
	Related User `gorm:"foreignKey:RelatedID" json:"related"`
}

// BeforeCreate sets the ID before creating a relationship
func (r *Relationship) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Block represents blocked users
type Block struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	BlockedID uuid.UUID `gorm:"type:uuid;not null" json:"blocked_id"`
	Reason    string    `gorm:"type:text" json:"reason"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User    User `gorm:"foreignKey:UserID" json:"user"`
	Blocked User `gorm:"foreignKey:BlockedID" json:"blocked"`
}

// BeforeCreate sets the ID before creating a block
func (b *Block) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// Report represents user reports
type Report struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	ReporterID uuid.UUID `gorm:"type:uuid;not null" json:"reporter_id"`
	ReportedID uuid.UUID `gorm:"type:uuid;not null" json:"reported_id"`
	Reason     string    `gorm:"type:text;not null" json:"reason"`
	Status     string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt  time.Time `json:"created_at"`

	// Relationships
	Reporter User `gorm:"foreignKey:ReporterID" json:"reporter"`
	Reported User `gorm:"foreignKey:ReportedID" json:"reported"`
}

// BeforeCreate sets the ID before creating a report
func (r *Report) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

// Invite represents user invitations
type Invite struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
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

// BeforeCreate sets the ID before creating an invite
func (i *Invite) BeforeCreate(tx *gorm.DB) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return nil
}
