package main

import (
	"fmt"
	"log"

	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Test password: test123
	plainPassword := "test123"

	// Create legacy crypto
	legacyCrypto := auth.NewLegacyCrypto(cfg.Legacy.CryptoKey)

	// Encrypt password
	encrypted, err := legacyCrypto.Encrypt(plainPassword)
	if err != nil {
		log.Fatalf("Failed to encrypt: %v", err)
	}

	fmt.Printf("Original password: %s\n", plainPassword)
	fmt.Printf("Encrypted password: %s\n", encrypted)
	fmt.Printf("Crypto key: %s\n", cfg.Legacy.CryptoKey)

	// Test decryption
	decrypted, err := legacyCrypto.Decrypt(encrypted)
	if err != nil {
		log.Fatalf("Failed to decrypt: %v", err)
	}

	fmt.Printf("Decrypted password: %s\n", decrypted)
	fmt.Printf("Round-trip successful: %t\n", plainPassword == decrypted)
}
