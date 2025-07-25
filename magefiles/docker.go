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
	if err := sh.Run("docker-compose", "up", "-d", "database"); err != nil {
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

// Build builds the production Docker image.
func (Docker) Build() error {
	fmt.Println("🐳 Building production Docker image...")
	return sh.Run("docker", "build", "-f", "docker/Dockerfile.prod", "-t", "controlme-server", ".")
}

// BuildDev builds the development Docker image.
func (Docker) BuildDev() error {
	fmt.Println("🐳 Building development Docker image...")
	return sh.Run("docker", "build", "-f", "docker/Dockerfile.dev", "-t", "controlme-server:dev", ".")
}

// Restart restarts Docker services.
func (Docker) Restart() error {
	fmt.Println("🔄 Restarting Docker services...")
	mg.Deps(Docker{}.Down)
	return Docker{}.Up()
}

// Status shows the status of Docker services.
func (Docker) Status() error {
	fmt.Println("📊 Docker service status:")
	return sh.RunV("docker-compose", "ps")
}
