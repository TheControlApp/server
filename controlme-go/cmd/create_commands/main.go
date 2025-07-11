package main

import (
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/database"
	"github.com/thecontrolapp/controlme-go/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.Initialize(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create services
	cmdService := services.NewCommandService(db)

	// Find testdom and testsub users  
	var users []struct {
		ID         uuid.UUID `gorm:"column:id"`
		LoginName  string    `gorm:"column:login_name"`
		ScreenName string    `gorm:"column:screen_name"`
		Role       string    `gorm:"column:role"`
	}

	err = db.Table("users").Select("id, login_name, screen_name, role").Find(&users).Error
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}

	fmt.Println("=== Current Users ===")
	var domID, subID uuid.UUID
	for _, user := range users {
		fmt.Printf("ID: %s, Login: %s, Screen: %s, Role: %s\n", 
			user.ID, user.LoginName, user.ScreenName, user.Role)
		
		if user.LoginName == "testdom" {
			domID = user.ID
		}
		if user.LoginName == "testsub" {
			subID = user.ID
		}
	}

	if domID == uuid.Nil || subID == uuid.Nil {
		fmt.Println("Warning: testdom or testsub user not found!")
		os.Exit(1)
	}

	fmt.Printf("\nDom ID: %s\nSub ID: %s\n\n", domID, subID)

	// Create some test commands
	fmt.Println("=== Creating Test Commands ===")

	commands := []struct {
		Type    string
		Content string
		Data    string
	}{
		{"message", "Hello there, pet!", ""},
		{"task", "Go make me a sandwich", "kitchen"},
		{"pose", "Kneel and wait for 5 minutes", "duration:5min"},
		{"rule", "No speaking unless spoken to", "silence"},
	}

	for i, cmd := range commands {
		// Create command
		command, err := cmdService.CreateCommand(cmd.Type, cmd.Content, cmd.Data)
		if err != nil {
			log.Printf("Failed to create command %d: %v", i+1, err)
			continue
		}

		// Assign to sub
		assignment, err := cmdService.AssignCommandToUser(domID, subID, command.ID, nil)
		if err != nil {
			log.Printf("Failed to assign command %d: %v", i+1, err)
			continue
		}

		fmt.Printf("✓ Created command: %s -> %s (Assignment ID: %s)\n", 
			cmd.Content, assignment.ID, command.ID)
	}

	// Show pending commands count
	count, err := cmdService.GetPendingCommandCount(subID)
	if err != nil {
		log.Fatalf("Failed to get command count: %v", err)
	}

	fmt.Printf("\n=== Result ===\n")
	fmt.Printf("Pending commands for testsub: %d\n", count)
	fmt.Println("Now test with: go run cmd/test_auth/main.go testpass1")
}
