// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// PHPAnalyzer analyzes PHP codebases for patterns
type PHPAnalyzer struct{}

func (pa *PHPAnalyzer) Language() string {
	return "php"
}

func (pa *PHPAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes PHP files and extracts patterns
func (pa *PHPAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".php") {
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
				Description: "PHP uses try/catch for error handling",
				Example:     "try {\n} catch (Exception $e) {\n}",
				Language:    "php",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, "throw new"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "throw exceptions",
				Description: "Throwing exceptions for error handling",
				Example:     "throw new Exception('error message')",
				Language:    "php",
				Frequency:   count,
			})
		}

		// Validation
		if count := countOccurrences(text, "is_null("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "is_null() checks",
				Description: "Using is_null() for null validation",
				Example:     "if (is_null($var)) { }",
				Language:    "php",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, "isset("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "isset() checks",
				Description: "Using isset() to check if variable is set",
				Example:     "if (isset($var)) { }",
				Language:    "php",
				Frequency:   count,
			})
		}

		// Logging
		if count := countOccurrences(text, "error_log("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "error_log()",
				Description: "Using error_log() for logging",
				Example:     "error_log('message')",
				Language:    "php",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, "echo "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "echo output",
				Description: "Using echo for output",
				Example:     "echo 'message'",
				Language:    "php",
				Frequency:   count,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
