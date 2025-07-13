package main

import (
	"encoding/base64"
	"fmt"
)

func main() {
	// The password your .NET client sent
	clientPassword := "BPOcH5nw8HazzzgfkiZxKlvw=="

	fmt.Printf("Analyzing client password: %s\n", clientPassword)
	fmt.Printf("Length: %d\n", len(clientPassword))

	// Try to decode as base64
	decoded, err := base64.StdEncoding.DecodeString(clientPassword)
	if err != nil {
		fmt.Printf("Base64 decode error: %v\n", err)

		// Try with different padding
		fmt.Println("Trying to fix padding...")

		// Remove padding and re-add
		trimmed := clientPassword
		for len(trimmed)%4 != 0 {
			trimmed += "="
		}

		decoded, err = base64.StdEncoding.DecodeString(trimmed)
		if err != nil {
			fmt.Printf("Still can't decode: %v\n", err)
		} else {
			fmt.Printf("Success with padding fix!\n")
		}
	}

	if err == nil {
		fmt.Printf("Decoded bytes: %v\n", decoded)
		fmt.Printf("Decoded string: %s\n", string(decoded))
		fmt.Printf("Decoded hex: %x\n", decoded)
		fmt.Printf("Decoded length: %d bytes\n", len(decoded))
	}
}
