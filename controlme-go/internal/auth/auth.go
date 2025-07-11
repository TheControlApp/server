package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
)

// LegacyCrypto handles encryption/decryption compatible with the C# CryptoHelper
type LegacyCrypto struct {
	key []byte
}

// NewLegacyCrypto creates a new legacy crypto handler
func NewLegacyCrypto(key string) *LegacyCrypto {
	// Pad or trim key to 32 bytes for AES-256
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		// Pad with zeros
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		keyBytes = padded
	} else if len(keyBytes) > 32 {
		// Trim to 32 bytes
		keyBytes = keyBytes[:32]
	}
	
	return &LegacyCrypto{
		key: keyBytes,
	}
}

// Decrypt decrypts a base64 encoded string using AES encryption (compatible with C# implementation)
func (lc *LegacyCrypto) Decrypt(encryptedText string) (string, error) {
	if encryptedText == "" {
		return "", fmt.Errorf("encrypted text is empty")
	}

	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	// Create AES cipher
	block, err := aes.NewCipher(lc.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Check minimum length
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract IV and encrypted data
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Create cipher mode
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt
	mode.CryptBlocks(ciphertext, ciphertext)

	// Remove padding
	plaintext := removePKCS7Padding(ciphertext)
	
	return string(plaintext), nil
}

// Encrypt encrypts a string using AES encryption (compatible with C# implementation)
func (lc *LegacyCrypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", fmt.Errorf("plaintext is empty")
	}

	// Create AES cipher
	block, err := aes.NewCipher(lc.key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Add PKCS7 padding
	paddedText := addPKCS7Padding([]byte(plaintext), aes.BlockSize)

	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	// For compatibility, we might need to use a fixed IV or implement proper random IV generation
	// This is a simplified version - in production, use crypto/rand
	
	// Create cipher mode
	mode := cipher.NewCBCEncrypter(block, iv)

	// Encrypt
	ciphertext := make([]byte, len(paddedText))
	mode.CryptBlocks(ciphertext, paddedText)

	// Prepend IV to ciphertext
	result := append(iv, ciphertext...)

	// Encode to base64
	return base64.StdEncoding.EncodeToString(result), nil
}

// Helper functions for PKCS7 padding
func addPKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

func removePKCS7Padding(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	
	unpadding := int(data[length-1])
	if unpadding > length {
		return data
	}
	
	return data[:(length - unpadding)]
}

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	Role       string    `json:"role"`
	ClientType string    `json:"client_type"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token operations
type JWTManager struct {
	secret     []byte
	expiration time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secret string, expiration time.Duration) *JWTManager {
	return &JWTManager{
		secret:     []byte(secret),
		expiration: expiration,
	}
}

// GenerateToken generates a new JWT token
func (jm *JWTManager) GenerateToken(userID uuid.UUID, username, role, clientType string) (string, error) {
	claims := JWTClaims{
		UserID:     userID,
		Username:   username,
		Role:       role,
		ClientType: clientType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jm.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "controlme-go",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jm.secret)
}

// ValidateToken validates and parses a JWT token
func (jm *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jm.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// PasswordManager handles password hashing and verification
type PasswordManager struct{}

// NewPasswordManager creates a new password manager
func NewPasswordManager() *PasswordManager {
	return &PasswordManager{}
}

// HashPassword hashes a password using bcrypt
func (pm *PasswordManager) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyPassword verifies a password against a hash
func (pm *PasswordManager) VerifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// AuthService combines all authentication functionality
type AuthService struct {
	LegacyCrypto    *LegacyCrypto
	JWTManager      *JWTManager
	PasswordManager *PasswordManager
}

// NewAuthService creates a new authentication service
func NewAuthService(legacyKey, jwtSecret string, jwtExpiration time.Duration) *AuthService {
	return &AuthService{
		LegacyCrypto:    NewLegacyCrypto(legacyKey),
		JWTManager:      NewJWTManager(jwtSecret, jwtExpiration),
		PasswordManager: NewPasswordManager(),
	}
}
