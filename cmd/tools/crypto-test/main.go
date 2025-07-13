package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
)

func main() {
	var (
		password = flag.String("password", "", "Password to encrypt")
		decrypt  = flag.String("decrypt", "", "Encrypted password to decrypt")
		help     = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help || (*password == "" && *decrypt == "") {
		fmt.Println("🔐 ControlMe Legacy Crypto Tool")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/tools/crypto-test/main.go [options]")
		fmt.Println()
		fmt.Println("Options:")
		fmt.Println("  --password string    Password to encrypt")
		fmt.Println("  --decrypt string     Encrypted password to decrypt")
		fmt.Println("  --help               Show this help")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  # Encrypt a password")
		fmt.Println("  go run cmd/tools/crypto-test/main.go --password secret123")
		fmt.Println()
		fmt.Println("  # Decrypt a password")
		fmt.Println("  go run cmd/tools/crypto-test/main.go --decrypt 'encrypted_base64_string'")
		return
	}

	// Load config to get crypto key
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load config: %v", err)
	}

	// Create legacy crypto instance
	legacyCrypto := auth.NewLegacyCrypto(cfg.Legacy.CryptoKey)

	if *password != "" {
		// Encrypt password
		encrypted, err := legacyCrypto.Encrypt(*password)
		if err != nil {
			log.Fatalf("❌ Failed to encrypt: %v", err)
		}

		fmt.Printf("🔐 Password Encryption Result:\n")
		fmt.Printf("   Original: %s\n", *password)
		fmt.Printf("   Encrypted: %s\n", encrypted)
		fmt.Printf("   Crypto Key: %s\n", cfg.Legacy.CryptoKey)
		fmt.Printf("\n")
		fmt.Printf("💡 Use this encrypted value in your test:\n")
		fmt.Printf("   curl 'http://localhost:8080/TestAuth.aspx?usernm=testuser&pwd=%s&vrs=012'\n", encrypted)
	}

	if *decrypt != "" {
		// Decrypt password
		decrypted, err := legacyCrypto.Decrypt(*decrypt)
		if err != nil {
			log.Fatalf("❌ Failed to decrypt: %v", err)
		}

		fmt.Printf("🔓 Password Decryption Result:\n")
		fmt.Printf("   Encrypted: %s\n", *decrypt)
		fmt.Printf("   Decrypted: %s\n", decrypted)
		fmt.Printf("   Crypto Key: %s\n", cfg.Legacy.CryptoKey)
	}
}
