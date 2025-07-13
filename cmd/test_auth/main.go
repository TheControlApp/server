package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run test_auth.go <password>")
		fmt.Println("Example: go run test_auth.go testpass1")
		os.Exit(1)
	}

	password := os.Args[1]

	// Load configuration to get the encryption key
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create crypto instance with the same key as the server
	crypto := auth.NewLegacyCrypto(cfg.Legacy.CryptoKey)

	// Encrypt the password
	encrypted, err := crypto.Encrypt(password)
	if err != nil {
		log.Fatalf("Failed to encrypt password: %v", err)
	}

	fmt.Printf("Original password: %s\n", password)
	fmt.Printf("Encrypted password: %s\n", encrypted)
	fmt.Println()

	// Test decryption to make sure it works
	decrypted, err := crypto.Decrypt(encrypted)
	if err != nil {
		log.Fatalf("Failed to decrypt password: %v", err)
	}

	fmt.Printf("Decrypted password: %s\n", decrypted)
	if decrypted != password {
		log.Fatalf("Decryption failed - passwords don't match!")
	}
	fmt.Println("✓ Encryption/Decryption test passed!")
	fmt.Println()

	// Now test with the actual endpoints
	fmt.Println("=== Testing Legacy Endpoints with Authentication ===")

	// URL encode the encrypted password
	encodedPassword := url.QueryEscape(encrypted)

	// Test GetCount.aspx
	fmt.Println("1. Testing GetCount.aspx...")
	testURL := fmt.Sprintf("http://localhost:8080/GetCount.aspx?usernm=testdom&pwd=%s&vrs=012", encodedPassword)
	resp, err := http.Get(testURL)
	if err != nil {
		log.Printf("Error testing GetCount: %v", err)
	} else {
		defer resp.Body.Close()
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		fmt.Printf("Response: %s\n", string(body[:n]))
	}
	fmt.Println()

	// Test GetOptions.aspx
	fmt.Println("2. Testing GetOptions.aspx...")
	testURL = fmt.Sprintf("http://localhost:8080/GetOptions.aspx?usernm=testdom&pwd=%s&vrs=012", encodedPassword)
	resp, err = http.Get(testURL)
	if err != nil {
		log.Printf("Error testing GetOptions: %v", err)
	} else {
		defer resp.Body.Close()
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		fmt.Printf("Response: %s\n", string(body[:n]))
	}
	fmt.Println()

	// Test GetContent.aspx
	fmt.Println("3. Testing GetContent.aspx...")
	testURL = fmt.Sprintf("http://localhost:8080/GetContent.aspx?usernm=testdom&pwd=%s&vrs=012", encodedPassword)
	resp, err = http.Get(testURL)
	if err != nil {
		log.Printf("Error testing GetContent: %v", err)
	} else {
		defer resp.Body.Close()
		body := make([]byte, 1024)
		n, _ := resp.Body.Read(body)
		fmt.Printf("Response: %s\n", string(body[:n]))
	}
	fmt.Println()

	fmt.Println("=== Authentication Tests Complete ===")
}
