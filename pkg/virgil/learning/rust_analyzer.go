// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// RustAnalyzer analyzes Rust codebases for patterns
type RustAnalyzer struct{}

func (ra *RustAnalyzer) Language() string {
	return "rust"
}

func (ra *RustAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes Rust files and extracts patterns
func (ra *RustAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".rs") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Error handling
		if count := countOccurrences(text, "match "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "match expressions",
				Description: "Using match for pattern matching and error handling",
				Example:     "match result {\n  Ok(v) => { },\n  Err(e) => { }\n}",
				Language:    "rust",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, "Result<"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "Result type",
				Description: "Using Result type for error handling",
				Example:     "fn foo() -> Result<T, E> { }",
				Language:    "rust",
				Frequency:   count,
			})
		}

		// Validation
		if count := countOccurrences(text, "Option<"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "Option type",
				Description: "Using Option type for optional values",
				Example:     "let x: Option<i32> = Some(5)",
				Language:    "rust",
				Frequency:   count,
			})
		}

		// Security
		if count := countOccurrences(text, "unsafe"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeSecurity,
				Name:        "unsafe blocks",
				Description: "Using unsafe code blocks",
				Example:     "unsafe { /* code */ }",
				Language:    "rust",
				Frequency:   count,
			})
		}

		// Logging
		if count := countOccurrences(text, "println!"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "println! macro",
				Description: "Using println! for output",
				Example:     "println!(\"message\");",
				Language:    "rust",
				Frequency:   count,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
