// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// JavaScriptAnalyzer analyzes JavaScript/TypeScript codebases for patterns
type JavaScriptAnalyzer struct{}

func (ja *JavaScriptAnalyzer) Language() string {
	return "javascript"
}

func (ja *JavaScriptAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes JavaScript/TypeScript files and extracts patterns
func (ja *JavaScriptAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		isJSFile := strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".jsx") ||
			strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".tsx")

		if !isJSFile {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Error handling
		if count := countOccurrences(text, "try {"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "try/catch blocks",
				Description: "JavaScript uses try/catch for error handling",
				Example:     "try {\n  // code\n} catch (error) {\n  // handle\n}",
				Language:    "javascript",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, ".catch("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "promise.catch()",
				Description: "Using promise.catch() for async error handling",
				Example:     "promise.catch(error => { /* handle */ })",
				Language:    "javascript",
				Frequency:   count,
			})
		}

		// Validation
		if count := countOccurrences(text, " === null") + countOccurrences(text, " !== null"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "null checks",
				Description: "Using === null for strict null checking",
				Example:     "if (value === null) { }",
				Language:    "javascript",
				Frequency:   count,
			})
		}

		// Logging
		if count := countOccurrences(text, "console.log("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "console.log()",
				Description: "Using console.log() for debugging",
				Example:     "console.log('message')",
				Language:    "javascript",
				Frequency:   count,
			})
		}

		// Security
		if count := countOccurrences(text, "require('crypto')") + countOccurrences(text, "import crypto"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeSecurity,
				Name:        "crypto module",
				Description: "Using crypto module for security",
				Example:     "const crypto = require('crypto')",
				Language:    "javascript",
				Frequency:   count,
			})
		}

		// Async patterns
		if count := countOccurrences(text, "async ") + countOccurrences(text, "await "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeStructure,
				Name:        "async/await",
				Description: "Using async/await for asynchronous operations",
				Example:     "async function() { await somePromise(); }",
				Language:    "javascript",
				Frequency:   count,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
