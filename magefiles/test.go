//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/thecontrolapp/controlme-go/internal/services"
)

// A namespace for test related commands.
type Test mg.Namespace

// All runs all tests.
func (Test) All() error {
	fmt.Println("🧪 Running all tests...")
	return sh.RunV("go", "test", "-v", "./...")
}

// Services runs tests for the services package.
func (Test) Services() error {
	fmt.Println("🧪 Running service tests...")
	return sh.RunV("go", "test", "-v", "./internal/services/...")
}

// API runs tests for the api package.
func (Test) API() error {
	fmt.Println("🧪 Running api tests...")
	return sh.RunV("go", "test", "-v", "./internal/api/...")
}

// CreateUser is an example of calling code directly from a mage task.
func (Test) CreateUser() error {
	mg.Deps(Docker.Up)
	fmt.Println("✨ Creating a test user...")

	// This is a simplified example. In a real scenario, you would
	// properly initialize your services and dependencies.
	// This also assumes your database connection logic can be
	// called from here.
	userService, err := getTestUserService()
	if err != nil {
		return fmt.Errorf("failed to get user service: %w", err)
	}

	user, err := userService.CreateUser(services.CreateUserRequest{
		LoginName:  "mage-test-user",
		ScreenName: "Mage Test User",
		Password:   "password",
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("✅ User created with ID: %s\n", user.ID)
	return nil
}

// getTestUserService is a helper to initialize services for testing.
// This would need to be adapted to your project's structure.
func getTestUserService() (*services.UserService, error) {
	// This is where you would put your service initialization logic.
	// It's separated to show how you might structure it.
	// You'll need to replace this with your actual database connection
	// and service setup.
	return nil, fmt.Errorf("database connection not implemented in magefile")
}
