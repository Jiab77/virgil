// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// PerlAnalyzer analyzes Perl codebases for patterns
type PerlAnalyzer struct{}

func (pa *PerlAnalyzer) Language() string {
	return "perl"
}

func (pa *PerlAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes Perl files and extracts patterns
func (pa *PerlAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		isPerlFile := strings.HasSuffix(path, ".pl") || strings.HasSuffix(path, ".pm")
		if !isPerlFile {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Error handling
		if count := countOccurrences(text, "eval"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "eval error handling",
				Description: "Using eval for error handling",
				Example:     "eval { ... }; if ($@) { }",
				Language:    "perl",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, "die "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "die statements",
				Description: "Using die to exit with error",
				Example:     "die 'error message'",
				Language:    "perl",
				Frequency:   count,
			})
		}

		// Validation
		if count := countOccurrences(text, "defined("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "defined() checks",
				Description: "Using defined() to check if variable is defined",
				Example:     "if (defined($var)) { }",
				Language:    "perl",
				Frequency:   count,
			})
		}

		// Logging
		if count := countOccurrences(text, "print "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "print output",
				Description: "Using print for output",
				Example:     "print \"message\\n\"",
				Language:    "perl",
				Frequency:   count,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
