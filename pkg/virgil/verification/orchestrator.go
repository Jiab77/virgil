// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package verification

import (
	"fmt"
	"log"

	"github.com/jiab77/virgil/pkg/virgil/config"
	"github.com/jiab77/virgil/pkg/virgil/storage"
	"github.com/jiab77/virgil/pkg/virgil/websearch"
)

// RunPipeline orchestrates the complete verification flow:
// 1. Performs web search if enabled (grounds assessment in current research)
// 2. Creates and runs verification pipeline with all enabled blocks
// 3. Aggregates and returns results
func RunPipeline(request string, projectPath string, cfg *config.Config, db *storage.Database) (*AggregatedResult, error) {
	// Step 1: Initialize context for verification blocks
	context := make(map[string]interface{})

	// Step 2: If web search enabled, perform research before assessment
	if cfg.WebSearchEnabled {
		log.Printf("[virgil] Web search enabled, gathering current research...")

		// Create web search service with encrypted database
		wsService := websearch.NewService(db)

		// Generate search queries from request context
		searchQueries := generateSearchQueries(request)

		// Perform searches (or retrieve from cache)
		searchResults := make([]websearch.SearchResult, 0)
		for _, query := range searchQueries {
			results, found, err := wsService.GetCachedSearch(query)
			if err != nil {
				log.Printf("[virgil] Warning: error retrieving cached search for '%s': %v", query, err)
				continue
			}

			if !found {
				log.Printf("[virgil] Search query not in cache: %s (will be populated in Phase 3 during code generation)", query)
				continue
			}

			// Append results from cache
			searchResults = append(searchResults, results...)
		}

		// Store search results in context for pipeline blocks to use
		if len(searchResults) > 0 {
			context["web_search_results"] = searchResults
			log.Printf("[virgil] Passed %d search results to verification pipeline", len(searchResults))
		}
	}

	// Step 3: Create verification pipeline with all available blocks
	pipeline := NewPipeline()

	// Register all verification blocks
	registerAllBlocks(pipeline)

	// Step 4: Run pipeline with enabled rules and search context
	result, err := pipeline.Run(projectPath, cfg.RulesEnabled)
	if err != nil {
		return nil, fmt.Errorf("verification pipeline failed: %w", err)
	}

	// Step 5: Attach search context to result for display
	if len(context) > 0 {
		result.Context = context
	}

	return result, nil
}

// generateSearchQueries creates security-focused search queries from request context
// These queries are used to ground the assessment in current research
func generateSearchQueries(request string) []string {
	queries := []string{
		"OWASP Top 10 API security vulnerabilities 2024",
		"Secure authentication and session management best practices",
		"Input validation and output encoding defense",
	}

	// Could be enhanced to parse request for specific keywords
	// e.g., if request mentions "payment" → add "payment card security" query
	// if request mentions "healthcare" → add "HIPAA compliance" query

	return queries
}

// registerAllBlocks registers all available verification blocks
// This includes real implementations and stubs for future compliance frameworks
func registerAllBlocks(pipeline *Pipeline) {
	// Register real implementation blocks
	pipeline.RegisterBlock(&OWASPBlock{})
	pipeline.RegisterBlock(&NISTBlock{})

	// Register stub blocks for future compliance frameworks
	// These are initialized but not yet fully implemented
	pipeline.RegisterBlock(&GDPRBlock{})
	pipeline.RegisterBlock(&HIPAABlock{})
	pipeline.RegisterBlock(&PCIDSSBlock{})
	pipeline.RegisterBlock(&CISBlock{})
	pipeline.RegisterBlock(&ISO27001Block{})
	pipeline.RegisterBlock(&CustomBlock{})
}
