package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SimpleUser model for testing - minimal fields to avoid GORM issues
type SimpleUser struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key"`
	ScreenName string    `gorm:"size:50;not null"`
	LoginName  string    `gorm:"size:50;not null;unique"`
	Password   string    `gorm:"size:255;not null"`
	Role       string    `gorm:"size:50;default:'user'"`
	Verified   bool      `gorm:"default:true"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (u *SimpleUser) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func main() {
	var (
		username = flag.String("username", "", "Username (required)")
		password = flag.String("password", "", "Password (required)")
		list     = flag.Bool("list", false, "List users")
		create   = flag.Bool("create-table", false, "Create table")
	)
	flag.Parse()

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

	// Create table if requested
	if *create {
		if err := db.AutoMigrate(&SimpleUser{}); err != nil {
			log.Fatalf("Failed to create table: %v", err)
		}
		fmt.Println("✅ Table created")
		return
	}

	// List users if requested
	if *list {
		var users []SimpleUser
		if err := db.Find(&users).Error; err != nil {
			log.Fatalf("Failed to list users: %v", err)
		}
		fmt.Printf("Found %d users:\n", len(users))
		for _, u := range users {
			fmt.Printf("  %s (%s) - Role: %s, Verified: %t\n", u.LoginName, u.ScreenName, u.Role, u.Verified)
		}
		return
	}

	// Create user
	if *username == "" || *password == "" {
		fmt.Println("Usage: go run main.go --username user --password pass")
		fmt.Println("   or: go run main.go --list")
		fmt.Println("   or: go run main.go --create-table")
		return
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	user := SimpleUser{
		ID:         uuid.New(),
		ScreenName: *username,
		LoginName:  *username,
		Password:   string(hash),
		Role:       "user",
		Verified:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := db.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}

	fmt.Printf("✅ Created user: %s\n", *username)
}
