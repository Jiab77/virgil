// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// RubyAnalyzer analyzes Ruby codebases for patterns
type RubyAnalyzer struct{}

func (ra *RubyAnalyzer) Language() string {
	return "ruby"
}

func (ra *RubyAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes Ruby files and extracts patterns
func (ra *RubyAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".rb") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Error handling
		if count := countOccurrences(text, "begin"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "begin/rescue blocks",
				Description: "Ruby uses begin/rescue for error handling",
				Example:     "begin\n  # code\nrescue => e\nend",
				Language:    "ruby",
				Frequency:   count,
			})
		}

		// Validation
		if count := countOccurrences(text, ".nil?"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        ".nil? checks",
				Description: "Using .nil? to check for nil",
				Example:     "if value.nil?",
				Language:    "ruby",
				Frequency:   count,
			})
		}

		// Logging
		if count := countOccurrences(text, "puts "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "puts output",
				Description: "Using puts for output",
				Example:     "puts \"message\"",
				Language:    "ruby",
				Frequency:   count,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
