// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/keyset"
)

// EncryptionMethod represents the encryption strategy
type EncryptionMethod string

const (
	EncryptionRandom     EncryptionMethod = "random"
	EncryptionPassphrase EncryptionMethod = "passphrase"
)

// EncryptionConfig holds encryption settings
type EncryptionConfig struct {
	Method EncryptionMethod
	KeyHex string // Hex-encoded key for random method
}

// GenerateRandomKey creates a 32-byte random key
func GenerateRandomKey() ([]byte, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}
	return key, nil
}

// SaveKeyFile saves hex-encoded key to .virgil/encryption.key
func SaveKeyFile(key []byte) error {
	keyDir := filepath.Join(".", ".virgil")
	if err := os.MkdirAll(keyDir, 0700); err != nil {
		return fmt.Errorf("failed to create .virgil directory: %w", err)
	}

	keyFile := filepath.Join(keyDir, "encryption.key")
	keyHex := hex.EncodeToString(key)
	
	if err := os.WriteFile(keyFile, []byte(keyHex), 0600); err != nil {
		return fmt.Errorf("failed to write key file: %w", err)
	}
	
	return nil
}

// LoadKeyFile loads hex-encoded key from .virgil/encryption.key
func LoadKeyFile() ([]byte, error) {
	keyFile := filepath.Join(".", ".virgil", "encryption.key")
	
	data, err := os.ReadFile(keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}
	
	key, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode key: %w", err)
	}
	
	return key, nil
}

// CreateAEADPrimitive creates a Tink AEAD cipher from a key
func CreateAEADPrimitive(key []byte) (any, error) {
	// Create a keyset with ChaCha20-Poly1305
	template := aead.ChaCha20Poly1305KeyTemplate()
	handle, err := keyset.NewHandle(template)
	if err != nil {
		return nil, fmt.Errorf("failed to create keyset handle: %w", err)
	}

	primitive, err := aead.New(handle)
	if err != nil {
		return nil, fmt.Errorf("failed to create AEAD primitive: %w", err)
	}

	return primitive, nil
}
