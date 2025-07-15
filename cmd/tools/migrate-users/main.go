package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User model with correct database column names for the new Users table
type User struct {
	ID           uuid.UUID `gorm:"column:Id;type:uuid;primary_key" json:"id"`
	ScreenName   string    `gorm:"column:Screen Name;size:50;not null" json:"screen_name"`
	LoginName    string    `gorm:"column:Login Name;size:50;not null;unique" json:"login_name"`
	Password     string    `gorm:"column:Password;size:255;not null" json:"-"`
	Role         string    `gorm:"column:Role;size:50" json:"role"`
	RandOpt      bool      `gorm:"column:RandOpt;default:false" json:"rand_opt"`
	AnonCmd      bool      `gorm:"column:AnonCmd;default:false" json:"anon_cmd"`
	Verified     bool      `gorm:"column:Varified;default:false" json:"verified"` // Note: "Varified" matches legacy typo
	VerifiedCode int       `gorm:"column:VarifiedCode;default:0" json:"verified_code"`
	ThumbsUp     int       `gorm:"column:ThumbsUp;default:0" json:"thumbs_up"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
	LoginDate    time.Time `gorm:"column:LoginDate;default:CURRENT_TIMESTAMP" json:"login_date"`
}

func (User) TableName() string {
	return "Users"
}

// OldUser model with users table name
type OldUser struct {
	ID           uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
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

func (OldUser) TableName() string {
	return "users"
}

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.Database.Host,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Port,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("🔄 Dropping existing Users table...")
	db.Exec("DROP TABLE IF EXISTS \"Users\"")

	fmt.Println("🔄 Creating Users table...")

	// Create the new Users table
	if err := db.AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to create Users table: %v", err)
	}

	fmt.Println("✅ Users table created successfully")

	// Copy data from users to Users
	fmt.Println("🔄 Copying data from users table to Users table...")

	var oldUsers []OldUser
	if err := db.Find(&oldUsers).Error; err != nil {
		fmt.Printf("⚠️  Could not find old users table: %v\n", err)
		fmt.Println("✅ Migration complete (no data to copy)")
		return
	}

	for _, oldUser := range oldUsers {
		newUser := User{
			ID:           oldUser.ID,
			ScreenName:   oldUser.ScreenName,
			LoginName:    oldUser.LoginName,
			Password:     oldUser.Password,
			Role:         oldUser.Role,
			RandOpt:      oldUser.RandOpt,
			AnonCmd:      oldUser.AnonCmd,
			Verified:     oldUser.Verified,
			VerifiedCode: oldUser.VerifiedCode,
			ThumbsUp:     oldUser.ThumbsUp,
			CreatedAt:    oldUser.CreatedAt,
			UpdatedAt:    oldUser.UpdatedAt,
			LoginDate:    oldUser.LoginDate,
		}

		if err := db.Create(&newUser).Error; err != nil {
			fmt.Printf("⚠️  Failed to copy user %s: %v\n", oldUser.LoginName, err)
		} else {
			fmt.Printf("✅ Copied user: %s\n", oldUser.LoginName)
		}
	}

	fmt.Printf("🎉 Migration complete! Copied %d users to Users table\n", len(oldUsers))
}
