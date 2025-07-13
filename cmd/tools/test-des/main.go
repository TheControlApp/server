package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/thecontrolapp/controlme-go/internal/auth"
)

func main() {
	// Test the DES decryption with the exact values from your .NET client
	desCrypto := auth.NewLegacyDESCrypto()

	// Your client's encrypted value
	clientEncrypted := "BPOcH5nw8HazzzgfkiZxKlvw=="

	fmt.Printf("Testing DES decryption...\n")
	fmt.Printf("Input from .NET client: %s\n", clientEncrypted)

	// Try to decrypt
	decrypted, err := desCrypto.Decrypt(clientEncrypted)
	if err != nil {
		log.Printf("❌ Decryption failed: %v", err)

		// Let's analyze what we have
		fmt.Printf("\nAnalyzing the encrypted string:\n")
		fmt.Printf("Length: %d\n", len(clientEncrypted))
		fmt.Printf("Contains 'zzz': %t (should be '/' in original base64)\n", strings.Contains(clientEncrypted, "zzz"))

		// Show character replacements
		restored := restoreClientEncoding(clientEncrypted)
		fmt.Printf("After character restoration: %s\n", restored)

	} else {
		fmt.Printf("✅ Decryption successful!\n")
		fmt.Printf("Decrypted password: %s\n", decrypted)

		// Verify it matches the expected password
		if decrypted == "secret123" {
			fmt.Printf("🎉 Perfect! Matches expected password 'secret123'\n")
		} else {
			fmt.Printf("⚠️  Doesn't match expected 'secret123', got '%s'\n", decrypted)
		}
	}

	// Test encryption round-trip
	fmt.Printf("\nTesting encryption round-trip...\n")
	testPassword := "secret123"
	encrypted, err := desCrypto.Encrypt(testPassword)
	if err != nil {
		log.Printf("❌ Encryption failed: %v", err)
	} else {
		fmt.Printf("Original: %s\n", testPassword)
		fmt.Printf("Encrypted: %s\n", encrypted)

		// Decrypt it back
		roundTrip, err := desCrypto.Decrypt(encrypted)
		if err != nil {
			log.Printf("❌ Round-trip decryption failed: %v", err)
		} else {
			fmt.Printf("Round-trip: %s\n", roundTrip)
			fmt.Printf("Round-trip success: %t\n", roundTrip == testPassword)
		}
	}
}

func restoreClientEncoding(encoded string) string {
	result := encoded
	result = strings.ReplaceAll(result, "xxx", "\\")
	result = strings.ReplaceAll(result, "yyy", "&")
	result = strings.ReplaceAll(result, "zzz", "/")
	result = strings.ReplaceAll(result, "aaa", "]")
	result = strings.ReplaceAll(result, "ppp", "G0")
	result = strings.ReplaceAll(result, "lll", "0x")
	return result
}
