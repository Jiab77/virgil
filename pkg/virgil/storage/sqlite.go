// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package storage

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Database wraps SQLite connection and encryption
type Database struct {
	conn      *sql.DB
	cipher    any // Tink AEAD cipher for encryption
	encrypted bool
}

// InitDatabase initializes the SQLite database with schema and encryption
func InitDatabase(dbPath string) (*Database, error) {
	// Create .virgil directory if it doesn't exist
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Initialize encryption: load existing key or generate new one
	var key []byte
	var err error
	
	// Try to load existing key
	key, err = LoadKeyFile()
	if err != nil {
		// Key doesn't exist, generate new one
		key, err = GenerateRandomKey()
		if err != nil {
			return nil, fmt.Errorf("failed to generate encryption key: %w", err)
		}

		// Save the new key
		if err := SaveKeyFile(key); err != nil {
			return nil, fmt.Errorf("failed to save encryption key: %w", err)
		}
	}

	// Create AEAD cipher for encryption
	cipher, err := CreateAEADPrimitive(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create encryption cipher: %w", err)
	}

	// Open database connection
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db := &Database{
		conn:      conn,
		cipher:    cipher,
		encrypted: true,
	}

	// Load and execute schema
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

// initSchema loads and executes the database schema from scripts/init-db.sql
func (db *Database) initSchema() error {
	schemaFile := filepath.Join("scripts", "init-db.sql")
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	_, err = db.conn.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

// Close closes the database connection
func (db *Database) Close() error {
	if db.conn != nil {
		return db.conn.Close()
	}
	return nil
}

// GetConnection returns the underlying database connection
func (db *Database) GetConnection() *sql.DB {
	return db.conn
}

// GetCipher returns the encryption cipher for use in sensitive operations
func (db *Database) GetCipher() any {
	return db.cipher
}

// IsEncrypted returns whether the database has encryption enabled
func (db *Database) IsEncrypted() bool {
	return db.encrypted
}

// CacheWebSearch stores encrypted web search results
func (db *Database) CacheWebSearch(query, queryHash, context, results, createdBy string, expiresAt *time.Time) error {
	// Encrypt sensitive fields before storing
	encryptedQuery, err := db.encryptField(query)
	if err != nil {
		return fmt.Errorf("failed to encrypt query: %w", err)
	}

	encryptedResults, err := db.encryptField(results)
	if err != nil {
		return fmt.Errorf("failed to encrypt results: %w", err)
	}

	stmt := `INSERT INTO web_search_cache (query, query_hash, search_context, results, expires_at, created_by)
			 VALUES (?, ?, ?, ?, ?, ?)`

	_, err = db.conn.Exec(stmt, encryptedQuery, queryHash, context, encryptedResults, expiresAt, createdBy)
	if err != nil {
		return fmt.Errorf("failed to cache web search: %w", err)
	}

	return nil
}

// GetCachedWebSearch retrieves and decrypts web search results
func (db *Database) GetCachedWebSearch(queryHash string) (query string, results string, found bool, err error) {
	var encryptedQuery, encryptedResults string
	var expiresAt *time.Time

	stmt := `SELECT query, results, expires_at FROM web_search_cache 
			 WHERE query_hash = ? AND (expires_at IS NULL OR expires_at > datetime('now'))`

	err = db.conn.QueryRow(stmt, queryHash).Scan(&encryptedQuery, &encryptedResults, &expiresAt)
	if err == sql.ErrNoRows {
		return "", "", false, nil
	}
	if err != nil {
		return "", "", false, fmt.Errorf("failed to query web search cache: %w", err)
	}

	// Decrypt sensitive fields
	query, err = db.decryptField(encryptedQuery)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to decrypt query: %w", err)
	}

	results, err = db.decryptField(encryptedResults)
	if err != nil {
		return "", "", false, fmt.Errorf("failed to decrypt results: %w", err)
	}

	return query, results, true, nil
}

// ListRecentWebSearches retrieves recent searches with pagination
func (db *Database) ListRecentWebSearches(limit int) ([]map[string]interface{}, error) {
	stmt := `SELECT id, query, search_context, timestamp 
			 FROM web_search_cache 
			 ORDER BY timestamp DESC 
			 LIMIT ?`

	rows, err := db.conn.Query(stmt, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent searches: %w", err)
	}
	defer rows.Close()

	var searches []map[string]interface{}
	for rows.Next() {
		var id int
		var encryptedQuery, context string
		var timestamp time.Time

		if err := rows.Scan(&id, &encryptedQuery, &context, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan search result: %w", err)
		}

		// Decrypt query for display
		query, err := db.decryptField(encryptedQuery)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt query: %w", err)
		}

		searches = append(searches, map[string]interface{}{
			"id":        id,
			"query":     query,
			"context":   context,
			"timestamp": timestamp,
		})
	}

	return searches, rows.Err()
}

// PruneExpiredWebSearchCache removes expired cache entries
func (db *Database) PruneExpiredWebSearchCache() (int64, error) {
	stmt := `DELETE FROM web_search_cache WHERE expires_at IS NOT NULL AND expires_at <= datetime('now')`

	result, err := db.conn.Exec(stmt)
	if err != nil {
		return 0, fmt.Errorf("failed to prune web search cache: %w", err)
	}

	return result.RowsAffected()
}

// encryptField encrypts a string field using the Tink AEAD cipher
func (db *Database) encryptField(plaintext string) (string, error) {
	if !db.encrypted {
		return plaintext, nil
	}

	aead, ok := db.cipher.(interface {
		Encrypt(plaintext, additionalData []byte) ([]byte, error)
	})
	if !ok {
		return "", fmt.Errorf("cipher does not support Encrypt operation")
	}

	ciphertext, err := aead.Encrypt([]byte(plaintext), nil)
	if err != nil {
		return "", fmt.Errorf("encryption failed: %w", err)
	}

	// Encode as hex for storage in text field
	return fmt.Sprintf("%x", ciphertext), nil
}

// decryptField decrypts a string field using the Tink AEAD cipher
func (db *Database) decryptField(ciphertext string) (string, error) {
	if !db.encrypted {
		return ciphertext, nil
	}

	aead, ok := db.cipher.(interface {
		Decrypt(ciphertext, additionalData []byte) ([]byte, error)
	})
	if !ok {
		return "", fmt.Errorf("cipher does not support Decrypt operation")
	}

	// Decode from hex
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	plaintext, err := aead.Decrypt(ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

// SaveLearnedPattern stores an extracted pattern in encrypted form
func (db *Database) SaveLearnedPattern(patternID, patternType, language, name, description, filePath, example, metadata string) error {
	// Encrypt example code for security
	var encryptedExample string
	var err error
	if example != "" {
		encryptedExample, err = db.encryptField(example)
		if err != nil {
			return fmt.Errorf("failed to encrypt pattern example: %w", err)
		}
	}

	stmt := `INSERT INTO learned_patterns (pattern_id, pattern_type, language, name, description, file_path, example, metadata)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)
			 ON CONFLICT(pattern_id) DO UPDATE SET 
				 frequency = frequency + 1,
				 updated_at = CURRENT_TIMESTAMP`

	_, err = db.conn.Exec(stmt, patternID, patternType, language, name, description, filePath, encryptedExample, metadata)
	if err != nil {
		return fmt.Errorf("failed to save learned pattern: %w", err)
	}

	return nil
}

// GetLearnedPatterns retrieves all patterns for a language with frequency sorting
func (db *Database) GetLearnedPatterns(language string) ([]map[string]interface{}, error) {
	stmt := `SELECT pattern_id, pattern_type, language, name, description, file_path, example, frequency, metadata
			 FROM learned_patterns
			 WHERE language = ?
			 ORDER BY pattern_type ASC, frequency DESC, updated_at DESC`

	rows, err := db.conn.Query(stmt, language)
	if err != nil {
		return nil, fmt.Errorf("failed to query learned patterns: %w", err)
	}
	defer rows.Close()

	var patterns []map[string]interface{}
	for rows.Next() {
		var patternID, patternType, lang, name, description, filePath, encryptedExample, metadata string
		var frequency int

		if err := rows.Scan(&patternID, &patternType, &lang, &name, &description, &filePath, &encryptedExample, &frequency, &metadata); err != nil {
			return nil, fmt.Errorf("failed to scan pattern: %w", err)
		}

		// Decrypt example if present
		var example string
		if encryptedExample != "" {
			var decErr error
			example, decErr = db.decryptField(encryptedExample)
			if decErr != nil {
				// Log but don't fail - example might be corrupted
				example = "[decryption failed]"
			}
		}

		patterns = append(patterns, map[string]interface{}{
			"pattern_id":   patternID,
			"pattern_type": patternType,
			"language":     lang,
			"name":         name,
			"description":  description,
			"file_path":    filePath,
			"example":      example,
			"frequency":    frequency,
			"metadata":     metadata,
		})
	}

	return patterns, rows.Err()
}

// GetPatternsByType retrieves patterns filtered by type
func (db *Database) GetPatternsByType(language, patternType string) ([]map[string]interface{}, error) {
	stmt := `SELECT pattern_id, pattern_type, language, name, description, file_path, example, frequency, metadata
			 FROM learned_patterns
			 WHERE language = ? AND pattern_type = ?
			 ORDER BY frequency DESC, updated_at DESC`

	rows, err := db.conn.Query(stmt, language, patternType)
	if err != nil {
		return nil, fmt.Errorf("failed to query patterns by type: %w", err)
	}
	defer rows.Close()

	var patterns []map[string]interface{}
	for rows.Next() {
		var patternID, patternType, lang, name, description, filePath, encryptedExample, metadata string
		var frequency int

		if err := rows.Scan(&patternID, &patternType, &lang, &name, &description, &filePath, &encryptedExample, &frequency, &metadata); err != nil {
			return nil, fmt.Errorf("failed to scan pattern: %w", err)
		}

		// Decrypt example if present
		var example string
		if encryptedExample != "" {
			var decErr error
			example, decErr = db.decryptField(encryptedExample)
			if decErr != nil {
				example = "[decryption failed]"
			}
		}

		patterns = append(patterns, map[string]interface{}{
			"pattern_id":   patternID,
			"pattern_type": patternType,
			"language":     lang,
			"name":         name,
			"description":  description,
			"file_path":    filePath,
			"example":      example,
			"frequency":    frequency,
			"metadata":     metadata,
		})
	}

	return patterns, rows.Err()
}
