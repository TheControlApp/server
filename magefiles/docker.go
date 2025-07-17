//go:build mage

package main

import (
	"fmt"
	"time"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// A namespace for Docker related commands.
type Docker mg.Namespace

// Up starts docker services and waits for them to be healthy.
func (Docker) Up() error {
	fmt.Println("🐳 Starting Docker services...")
	if err := sh.Run("docker-compose", "up", "-d", "postgres"); err != nil {
		return err
	}

	fmt.Println("⏳ Waiting for services to be ready...")
	// Simple sleep, can be replaced with a more robust health check
	time.Sleep(5 * time.Second)

	fmt.Println("📊 Database services started successfully!")
	return nil
}

// Down stops docker services.
func (Docker) Down() error {
	fmt.Println("🛑 Stopping Docker services...")
	return sh.Run("docker-compose", "down")
}

// Clean stops services and removes docker volumes.
func (Docker) Clean() error {
	fmt.Println("🧹 Cleaning Docker volumes...")
	return sh.Run("docker-compose", "down", "-v", "--remove-orphans")
}

// Logs shows docker service logs.
func (Docker) Logs() error {
	fmt.Println("📋 Docker service logs:")
	return sh.RunV("docker-compose", "logs", "-f")
}
