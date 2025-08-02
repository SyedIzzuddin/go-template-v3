package tokens

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	// TokenLength defines the byte length of verification tokens
	TokenLength = 32
	// DefaultTokenExpiry defines default token expiration time
	DefaultTokenExpiry = 24 * time.Hour
)

// GenerateVerificationToken generates a cryptographically secure random token
func GenerateVerificationToken() (string, error) {
	bytes := make([]byte, TokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate verification token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateTokenWithExpiry generates a token and calculates its expiry time
func GenerateTokenWithExpiry() (string, time.Time, error) {
	token, err := GenerateVerificationToken()
	if err != nil {
		return "", time.Time{}, err
	}
	
	expiry := time.Now().Add(DefaultTokenExpiry)
	return token, expiry, nil
}

// IsTokenExpired checks if a token has expired
func IsTokenExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return true
	}
	return time.Now().After(*expiresAt)
}

// ValidateToken performs basic token validation
func ValidateToken(token string) error {
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}
	
	if len(token) != TokenLength*2 { // hex encoding doubles the length
		return fmt.Errorf("invalid token length")
	}
	
	// Verify it's valid hex
	if _, err := hex.DecodeString(token); err != nil {
		return fmt.Errorf("invalid token format: %w", err)
	}
	
	return nil
}

// GeneratePasswordResetToken generates a cryptographically secure password reset token
// This is an alias to GenerateVerificationToken for clarity and future extensibility
func GeneratePasswordResetToken() (string, error) {
	return GenerateVerificationToken()
}

// GeneratePasswordResetTokenWithExpiry generates a password reset token and calculates its expiry time
func GeneratePasswordResetTokenWithExpiry() (string, time.Time, error) {
	token, err := GeneratePasswordResetToken()
	if err != nil {
		return "", time.Time{}, err
	}
	
	expiry := time.Now().Add(DefaultTokenExpiry)
	return token, expiry, nil
}