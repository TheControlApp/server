//go:build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// A namespace for installing tools.
type Tools mg.Namespace

// Install installs all development tools.
func (Tools) Install() {
	mg.Deps(Tools.GolangciLint, Tools.Air, Tools.Swag, Tools.Mockgen)
}

// GolangciLint installs the golangci-lint tool.
func (Tools) GolangciLint() error {
	fmt.Println("🔧 Installing golangci-lint...")
	return sh.RunV("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@latest")
}

// Air installs the air tool for hot reloading.
func (Tools) Air() error {
	fmt.Println("🔧 Installing air...")
	return sh.RunV("go", "install", "github.com/air-verse/air@latest")
}

// Swag installs the swag tool for swagger documentation.
func (Tools) Swag() error {
	fmt.Println("🔧 Installing swag...")
	return sh.RunV("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
}

// Mockgen installs the mockgen tool for generating mocks.
func (Tools) Mockgen() error {
	fmt.Println("🔧 Installing mockgen...")
	return sh.RunV("go", "install", "github.com/golang/mock/mockgen@latest")
}

// Gofumpt installs gofumpt for stricter formatting.
func (Tools) Gofumpt() error {
	fmt.Println("🔧 Installing gofumpt...")
	return sh.RunV("go", "install", "mvdan.cc/gofumpt@latest")
}
