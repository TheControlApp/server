//go:build mage

// Main magefile for the ControlMe server project.
// This file provides main targets while organized namespace targets are in magefiles/ directory.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	serverBinary = "server"
)

func init() {
	if runtime.GOOS == "windows" {
		serverBinary = "server.exe"
	}
}

// Default target to run when none is specified
var Default = Build

// Help shows available targets and namespace organization
func Help() {
	fmt.Println("ControlMe Server - Mage Build System")
	fmt.Println("=====================================")
	fmt.Println("")
	fmt.Println("📋 Available targets:")
	sh.RunV(os.Args[0], "-l")
	fmt.Println("")
	fmt.Println("🏗️  Namespaces:")
	fmt.Println("  docker:*  - Docker container management")
	fmt.Println("  tools:*   - Development tool installation")
	fmt.Println("  test:*    - Testing utilities")
	fmt.Println("")
	fmt.Println("Example: mage docker:up")
}

// Build builds the application binary
func Build() error {
	fmt.Println("🔨 Building server...")
	if err := os.MkdirAll("bin", os.ModePerm); err != nil {
		return err
	}
	if err := sh.Run("go", "build", "-o", filepath.Join("bin", serverBinary), "cmd/server/main.go"); err != nil {
		return err
	}
	fmt.Printf("✅ Build complete: bin/%s\n", serverBinary)
	return nil
}

// Run runs the main application
func Run() error {
	fmt.Println("🚀 Running server...")
	return sh.Run("go", "run", "cmd/server/main.go")
}

// Dev starts the development environment with hot reloading
func Dev() error {
	mg.Deps(Docker{}.Up)
	fmt.Println("🚀 Starting development environment with hot-reload...")
	return sh.RunV("air")
}

// Fmt formats the code
func Fmt() error {
	fmt.Println("🎨 Formatting code...")
	return sh.RunV("go", "fmt", "./...")
}

// Vet runs go vet
func Vet() error {
	fmt.Println("🔍 Running go vet...")
	return sh.RunV("go", "vet", "./...")
}

// Lint runs the linter (uses Tools namespace)
func Lint() error {
	mg.Deps(Tools{}.GolangciLint)
	fmt.Println("🔍 Running linter...")
	return sh.RunV("golangci-lint", "run")
}

// Clean removes build artifacts
func Clean() {
	fmt.Println("🧹 Cleaning build artifacts...")
	os.RemoveAll("bin")
	os.RemoveAll("tmp")
	os.Remove("build-errors.log")
	fmt.Println("✅ Clean complete")
}

// Setup runs the complete development setup
func Setup() error {
	fmt.Println("🚀 Setting up development environment...")

	fmt.Println("📁 Creating directories...")
	for _, dir := range []string{"bin", "logs", "tmp", "docs/swagger", "docs/examples"} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	fmt.Println("📦 Installing Go dependencies...")
	if err := sh.Run("go", "mod", "tidy"); err != nil {
		return err
	}

	fmt.Println("🔧 Installing development tools...")
	mg.Deps(Tools{}.Install)

	fmt.Println("📚 Generating API documentation...")
	mg.Deps(Swagger)

	fmt.Println("🐳 Starting Docker services...")
	mg.Deps(Docker{}.Up)

	fmt.Println("✅ Setup complete!")
	return nil
}

// Swagger generates swagger documentation (uses API namespace)
func Swagger() error {
	return API{}.Swagger()
}

// Serve generates swagger docs and runs the server (uses API namespace)
func Serve() error {
	return API{}.Serve()
}
