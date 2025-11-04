package apikey

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	prefix    = "gsk_"
	keyLength = 32
)

// generateKey creates a new API key
func generateKey() (id, key string, err error) {
	// generateKey UUID for the key ID
	id = uuid.New().String()

	// generateKey random bytes
	randomBytes := make([]byte, keyLength)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// Encode to base64 (URL-safe, no padding)
	encoded := base64.RawURLEncoding.EncodeToString(randomBytes)

	// Create the full key with prefix
	key = prefix + encoded

	return id, key, nil
}

func hash(key string) string {
	hash := sha256.Sum256([]byte(key))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func ExtractParts(key string) (prefix, suffix string) {
	if !strings.HasPrefix(key, prefix) {
		return "", ""
	}

	if len(key) < len(prefix)+4 {
		return prefix, ""
	}

	return prefix, key[len(key)-4:]
}

// Validate checks if a key has the correct format
func Validate(key string) bool {
	if !strings.HasPrefix(key, prefix) {
		return false
	}

	// Remove prefix and check length
	withoutPrefix := strings.TrimPrefix(key, prefix)

	// Should be base64 encoded, roughly 43 characters for 32 bytes
	return len(withoutPrefix) >= 40
}
