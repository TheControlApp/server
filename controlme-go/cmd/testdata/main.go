package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/database"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create auth service
	authService := auth.NewAuthService(cfg.Legacy.CryptoKey, cfg.Auth.JWTSecret, time.Duration(cfg.Auth.JWTExpiration)*time.Hour)

	fmt.Println("Adding test data...")

	// Hash passwords
	password1, _ := bcrypt.GenerateFromPassword([]byte("testpass1"), bcrypt.DefaultCost)
	password2, _ := bcrypt.GenerateFromPassword([]byte("testpass2"), bcrypt.DefaultCost)

	// Create test users
	user1 := models.User{
		ID:         uuid.New(),
		ScreenName: "TestDom",
		LoginName:  "testdom",
		Password:   string(password1),
		Role:       "dom",
		Verified:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		LoginDate:  time.Now(),
	}

	user2 := models.User{
		ID:         uuid.New(),
		ScreenName: "TestSub",
		LoginName:  "testsub",
		Password:   string(password2),
		Role:       "sub",
		Verified:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		LoginDate:  time.Now(),
	}

	// Insert users
	if err := db.Create(&user1).Error; err != nil {
		log.Printf("Failed to create user1: %v", err)
	} else {
		fmt.Printf("Created user: %s (ID: %s)\n", user1.LoginName, user1.ID)
	}

	if err := db.Create(&user2).Error; err != nil {
		log.Printf("Failed to create user2: %v", err)
	} else {
		fmt.Printf("Created user: %s (ID: %s)\n", user2.LoginName, user2.ID)
	}

	// Create a relationship between the users
	relationship := models.Relationship{
		ID:        uuid.New(),
		DomID:     user1.ID,
		SubID:     user2.ID,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := db.Create(&relationship).Error; err != nil {
		log.Printf("Failed to create relationship: %v", err)
	} else {
		fmt.Printf("Created relationship: %s -> %s\n", user1.LoginName, user2.LoginName)
	}

	// Test encryption/decryption
	fmt.Println("\nTesting legacy crypto...")
	testData := "testpass1"
	encrypted, err := authService.LegacyCrypto.Encrypt(testData)
	if err != nil {
		log.Printf("Encryption failed: %v", err)
	} else {
		fmt.Printf("Encrypted: %s\n", encrypted)
		
		decrypted, err := authService.LegacyCrypto.Decrypt(encrypted)
		if err != nil {
			log.Printf("Decryption failed: %v", err)
		} else {
			fmt.Printf("Decrypted: %s\n", decrypted)
			if decrypted == testData {
				fmt.Println("✓ Crypto test passed!")
			} else {
				fmt.Println("✗ Crypto test failed!")
			}
		}
	}

	fmt.Println("\nTest data setup complete!")
}
