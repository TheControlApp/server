//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Test runs all Go tests in the project
func Test() error {
	cmd := exec.Command("go", "test", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running all tests...")
	return cmd.Run()
}

// Lint runs golangci-lint on the project
func Lint() error {
	cmd := exec.Command("golangci-lint", "run")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running linter...")
	return cmd.Run()
}

// Build builds the main server
func Build() error {
	cmd := exec.Command("go", "build", "-o", "bin/server", "./cmd/server/main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Building server...")
	return cmd.Run()
}
