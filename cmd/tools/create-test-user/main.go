package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
	"github.com/thecontrolapp/controlme-go/internal/database"
	"github.com/thecontrolapp/controlme-go/internal/models"
	"github.com/thecontrolapp/controlme-go/internal/services"
)

func main() {
	// Initialize logger
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatal("Failed to load configuration: ", err)
	}

	// Initialize database
	db, err := database.Connect(cfg.Database)
	if err != nil {
		logrus.Fatal("Failed to connect to database: ", err)
	}

	// Initialize services
	userService := services.NewUserService(db)
	authService := auth.NewAuthService(cfg.Auth.JWTSecret)

	// Interactive user creation
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("=== ControlMe Test User Creation ===")
	fmt.Println()

	// Get username
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	if username == "" {
		logrus.Fatal("Username cannot be empty")
	}

	// Get password
	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if password == "" {
		logrus.Fatal("Password cannot be empty")
	}

	// Get email
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	if email == "" {
		logrus.Fatal("Email cannot be empty")
	}

	// Create user
	user := &models.User{
		Username:  username,
		Email:     email,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Hash password
	hashedPassword, err := authService.HashPassword(password)
	if err != nil {
		logrus.Fatal("Failed to hash password: ", err)
	}
	user.PasswordHash = hashedPassword

	// Save user
	if err := userService.CreateUser(user); err != nil {
		logrus.Fatal("Failed to create user: ", err)
	}

	fmt.Printf("\n✅ Test user created successfully!\n")
	fmt.Printf("Username: %s\n", user.Username)
	fmt.Printf("Email: %s\n", user.Email)
	fmt.Printf("User ID: %d\n", user.ID)
	fmt.Printf("Active: %t\n", user.IsActive)
	fmt.Printf("Created: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
	fmt.Println("You can now use this user to test the API endpoints.")
}
