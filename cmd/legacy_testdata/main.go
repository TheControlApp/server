package main

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/database"
	"github.com/thecontrolapp/controlme-go/internal/models"
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

	fmt.Println("Creating legacy-compatible test data...")

	// Clear existing users
	db.Exec("DELETE FROM users")

	// Create test users with legacy-encrypted passwords
	testUsers := []struct {
		screenName string
		loginName  string
		password   string
		role       string
	}{
		{"TestDom", "testdom", "testpass1", "dom"},
		{"TestSub", "testsub", "testpass2", "sub"},
		{"AdminUser", "admin", "adminpass", "admin"},
	}

	for _, userData := range testUsers {
		// Encrypt password using legacy crypto (same as C# implementation)
		encryptedPassword, err := authService.LegacyCrypto.Encrypt(userData.password)
		if err != nil {
			log.Fatalf("Failed to encrypt password for %s: %v", userData.screenName, err)
		}

		user := models.User{
			ID:           uuid.New(),
			ScreenName:   userData.screenName,
			LoginName:    userData.loginName,
			Password:     encryptedPassword, // Store encrypted password
			Role:         userData.role,
			RandOpt:      false,
			AnonCmd:      true,
			Verified:     true,
			VerifiedCode: 123,
			ThumbsUp:     0,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			LoginDate:    time.Now(),
		}

		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("Failed to create user %s: %v", userData.screenName, err)
		}

		fmt.Printf("Created user: %s (login: %s) with encrypted password\n", userData.screenName, userData.loginName)
		fmt.Printf("  Plain password: %s\n", userData.password)
		fmt.Printf("  Encrypted password: %s\n", encryptedPassword)

		// Test decryption to verify it works
		decrypted, err := authService.LegacyCrypto.Decrypt(encryptedPassword)
		if err != nil {
			log.Printf("WARNING: Failed to decrypt password for %s: %v", userData.screenName, err)
		} else if decrypted != userData.password {
			log.Printf("WARNING: Decrypted password doesn't match for %s", userData.screenName)
		} else {
			fmt.Printf("  ✓ Encryption/decryption verified\n")
		}
		fmt.Println()
	}

	// Create some test commands
	fmt.Println("Creating test commands...")

	// Get the first two users for testing
	var testDom, testSub models.User
	db.Where("screen_name = ?", "TestDom").First(&testDom)
	db.Where("screen_name = ?", "TestSub").First(&testSub)

	if testDom.ID != uuid.Nil && testSub.ID != uuid.Nil {
		// Create a test command
		command := models.Command{
			ID:        uuid.New(),
			Type:      "message",
			Content:   "Test command from TestDom to TestSub",
			Data:      "This is test command data",
			Status:    "pending",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := db.Create(&command).Error; err != nil {
			log.Printf("Failed to create command: %v", err)
		} else {
			// Create command assignment
			assignment := models.ControlAppCmd{
				ID:        uuid.New(),
				SenderID:  testDom.ID,
				SubID:     testSub.ID,
				CommandID: command.ID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := db.Create(&assignment).Error; err != nil {
				log.Printf("Failed to create command assignment: %v", err)
			} else {
				fmt.Printf("Created test command from %s to %s\n", testDom.ScreenName, testSub.ScreenName)
			}
		}
	}

	fmt.Println("Legacy test data creation complete!")
	fmt.Println("\nYou can now test authentication with:")
	fmt.Println("  curl \"http://localhost:8080/Login.aspx?usernm=TestDom&pwd=testpass1&vrs=012\"")
	fmt.Println("  curl \"http://localhost:8080/Login.aspx?usernm=TestSub&pwd=testpass2&vrs=012\"")
}
