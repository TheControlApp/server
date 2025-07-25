//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
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

// Coverage runs tests with coverage reporting.
func (Test) Coverage() error {
	fmt.Println("📊 Running tests with coverage...")
	return sh.RunV("go", "test", "-coverprofile=coverage.out", "./...")
}

// CoverageHTML generates HTML coverage report.
func (Test) CoverageHTML() error {
	mg.Deps(Test{}.Coverage)
	fmt.Println("🌐 Generating HTML coverage report...")
	if err := sh.RunV("go", "tool", "cover", "-html=coverage.out", "-o", "coverage.html"); err != nil {
		return err
	}
	fmt.Println("✅ Coverage report generated: coverage.html")
	return nil
}

// Bench runs benchmark tests.
func (Test) Bench() error {
	fmt.Println("⚡ Running benchmark tests...")
	return sh.RunV("go", "test", "-bench=.", "./...")
}

// CreateUser creates a test user using the CLI tool.
func (Test) CreateUser() error {
	mg.Deps(Docker{}.Up)
	fmt.Println("👤 Creating test user...")
	return sh.Run("go", "run", "cmd/tools/create-test-user/main.go")
}

// CreateCommands creates sample commands using the CLI tool.
func (Test) CreateCommands() error {
	mg.Deps(Docker{}.Up)
	fmt.Println("📝 Creating sample commands...")
	return sh.Run("go", "run", "cmd/tools/create-commands/main.go")
}
