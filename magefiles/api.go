//go:build mage

package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// A namespace for API and documentation related commands.
type API mg.Namespace

// Swagger generates swagger documentation.
func (API) Swagger() error {
	mg.Deps(Tools{}.Swag)
	fmt.Println("📚 Generating Swagger documentation...")
	if err := sh.Run("swag", "init", "-g", "cmd/server/main.go", "-o", "docs/swagger"); err != nil {
		return err
	}
	fmt.Println("✅ Swagger documentation generated at docs/swagger/")
	fmt.Println("🌐 Access it at: http://localhost:8080/swagger/index.html")
	return nil
}

// Serve generates swagger docs and runs the server.
func (API) Serve() error {
	mg.Deps(API{}.Swagger)
	fmt.Println("🚀 Starting server with Swagger documentation...")
	fmt.Println("📖 Swagger UI: http://localhost:8080/swagger/index.html")
	fmt.Println("🏥 Health check: http://localhost:8080/health")
	return sh.Run("go", "run", "cmd/server/main.go")
}

// ValidateSwagger validates the swagger documentation.
func (API) ValidateSwagger() error {
	mg.Deps(Tools{}.Swag)
	fmt.Println("🔍 Validating Swagger documentation...")
	return sh.RunV("swag", "fmt", "-g", "cmd/server/main.go")
}

// Clean removes generated API documentation.
func (API) Clean() error {
	fmt.Println("🧹 Cleaning API documentation...")
	if err := os.RemoveAll("docs/swagger"); err != nil {
		return err
	}
	fmt.Println("✅ API documentation cleaned")
	return nil
}
