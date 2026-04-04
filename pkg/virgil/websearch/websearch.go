// Copyright 2026 Jiab77
// SPDX-License-Identifier: MIT

package websearch

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jiab77/virgil/pkg/virgil/storage"
)

// SearchResult represents a single web search result
type SearchResult struct {
	URL     string `json:"url"`
	Title   string `json:"title"`
	Snippet string `json:"snippet"`
	Source  string `json:"source"` // Domain or publication name
}

// SearchQuery represents a web search query and its results
type SearchQuery struct {
	Query     string         `json:"query"`
	Context   string         `json:"context"` // "create", "edit", "verify"
	Results   []SearchResult `json:"results"`
	Timestamp time.Time      `json:"timestamp"`
	ExpiresAt *time.Time     `json:"expires_at,omitempty"`
	CreatedBy string         `json:"created_by,omitempty"`
}

// Service handles web search operations with encrypted storage
type Service struct {
	db *storage.Database
}

// NewService creates a new web search service with encrypted database access
func NewService(db *storage.Database) *Service {
	return &Service{
		db: db,
	}
}

// queryHash generates a SHA256 hash of a query for deduplication
func queryHash(query string) string {
	hash := sha256.Sum256([]byte(query))
	return hex.EncodeToString(hash[:])
}

// CacheSearch stores encrypted search results in the database
func (s *Service) CacheSearch(query string, context string, results []SearchResult, createdBy string) error {
	if s.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	qHash := queryHash(query)
	resultsJSON, err := json.Marshal(results)
	if err != nil {
		return fmt.Errorf("error marshaling search results: %w", err)
	}

	// Cache for 7 days by default
	expiresAt := time.Now().AddDate(0, 0, 7)

	// Use encrypted storage method from Database struct
	err = s.db.CacheWebSearch(query, qHash, context, string(resultsJSON), createdBy, &expiresAt)
	if err != nil {
		return fmt.Errorf("error caching search results: %w", err)
	}

	return nil
}

// GetCachedSearch retrieves decrypted cached search results
func (s *Service) GetCachedSearch(query string) (*SearchQuery, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	qHash := queryHash(query)

	// Use encrypted retrieval method from Database struct
	decryptedQuery, resultsJSON, found, err := s.db.GetCachedWebSearch(qHash)
	if err != nil {
		return nil, fmt.Errorf("error retrieving cached search: %w", err)
	}

	if !found {
		return nil, nil // Cache miss, not an error
	}

	// Unmarshal results
	var results []SearchResult
	if err := json.Unmarshal([]byte(resultsJSON), &results); err != nil {
		return nil, fmt.Errorf("error unmarshaling cached results: %w", err)
	}

	sq := &SearchQuery{
		Query:     decryptedQuery,
		Context:   "", // Will be populated from ListRecentSearches if needed
		Results:   results,
		Timestamp: time.Now(),
		CreatedBy: "",
	}

	return sq, nil
}

// ListRecentSearches returns recent decrypted cached searches
func (s *Service) ListRecentSearches(limit int) ([]SearchQuery, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database connection not initialized")
	}

	if limit <= 0 {
		limit = 10
	}

	// Use encrypted retrieval method from Database struct
	searches, err := s.db.ListRecentWebSearches(limit)
	if err != nil {
		return nil, fmt.Errorf("error querying recent searches: %w", err)
	}

	var results []SearchQuery
	for _, search := range searches {
		sq := SearchQuery{
			Query:     search["query"].(string),
			Context:   search["context"].(string),
			Timestamp: search["timestamp"].(time.Time),
		}
		results = append(results, sq)
	}

	return results, nil
}

// PruneExpiredCache removes expired cached searches
func (s *Service) PruneExpiredCache() error {
	if s.db == nil {
		return fmt.Errorf("database connection not initialized")
	}

	// Use encrypted pruning method from Database struct
	rowsAffected, err := s.db.PruneExpiredWebSearchCache()
	if err != nil {
		return fmt.Errorf("error pruning expired cache: %w", err)
	}

	if rowsAffected > 0 {
		fmt.Printf("[websearch] Pruned %d expired cache entries\n", rowsAffected)
	}

	return nil
}
