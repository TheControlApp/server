package auth

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"fmt"
	"strings"
)

// LegacyDESCrypto handles DES encryption/decryption compatible with the C# client
type LegacyDESCrypto struct {
	publicKey []byte
	secretKey []byte
}

// NewLegacyDESCrypto creates a new legacy DES crypto handler with the exact keys from C# client
func NewLegacyDESCrypto() *LegacyDESCrypto {
	return &LegacyDESCrypto{
		publicKey: []byte("santhosh"), // Exactly as in C# client
		secretKey: []byte("engineer"), // Exactly as in C# client
	}
}

// Decrypt decrypts a string that was encrypted by the C# client
func (ldc *LegacyDESCrypto) Decrypt(encryptedText string) (string, error) {
	if encryptedText == "" {
		return "", fmt.Errorf("encrypted text is empty")
	}

	// First, reverse the character replacements made by C# client
	decoded := encryptedText
	decoded = strings.ReplaceAll(decoded, "xxx", "\\")
	decoded = strings.ReplaceAll(decoded, "yyy", "&")
	decoded = strings.ReplaceAll(decoded, "zzz", "/")
	decoded = strings.ReplaceAll(decoded, "aaa", "]")
	decoded = strings.ReplaceAll(decoded, "ppp", "G0")
	decoded = strings.ReplaceAll(decoded, "lll", "0x")

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(decoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create DES cipher with the same keys as C# client
	block, err := des.NewCipher(ldc.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to create DES cipher: %w", err)
	}

	// Check minimum length (must be multiple of DES block size)
	if len(ciphertext) < des.BlockSize || len(ciphertext)%des.BlockSize != 0 {
		return "", fmt.Errorf("invalid ciphertext length")
	}

	// Use the secret key as IV (this matches C# DESCryptoServiceProvider behavior)
	iv := make([]byte, des.BlockSize)
	copy(iv, ldc.secretKey)

	// Create CBC decrypter
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS5 padding
	plaintext = removePKCS5Padding(plaintext)

	return string(plaintext), nil
}

// Encrypt encrypts a string using the same method as the C# client
func (ldc *LegacyDESCrypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", fmt.Errorf("plaintext is empty")
	}

	// Create DES cipher
	block, err := des.NewCipher(ldc.publicKey)
	if err != nil {
		return "", fmt.Errorf("failed to create DES cipher: %w", err)
	}

	// Add PKCS5 padding
	paddedText := addPKCS5Padding([]byte(plaintext), des.BlockSize)

	// Use the secret key as IV (matching C# behavior)
	iv := make([]byte, des.BlockSize)
	copy(iv, ldc.secretKey)

	// Create CBC encrypter
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt
	ciphertext := make([]byte, len(paddedText))
	mode.CryptBlocks(ciphertext, paddedText)

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	// Apply the same character replacements as C# client
	result := encoded
	result = strings.ReplaceAll(result, "\\", "xxx")
	result = strings.ReplaceAll(result, "&", "yyy")
	result = strings.ReplaceAll(result, "/", "zzz")
	result = strings.ReplaceAll(result, "]", "aaa")
	result = strings.ReplaceAll(result, "G0", "ppp")
	result = strings.ReplaceAll(result, "0x", "lll")

	return result, nil
}

// Helper functions for PKCS5 padding (DES uses PKCS5, not PKCS7)
func addPKCS5Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

func removePKCS5Padding(data []byte) []byte {
	if len(data) == 0 {
		return data
	}
	padding := int(data[len(data)-1])
	if padding > len(data) || padding == 0 {
		return data
	}
	return data[:len(data)-padding]
}
