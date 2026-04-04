// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// PythonAnalyzer analyzes Python codebases for patterns
type PythonAnalyzer struct{}

func (pa *PythonAnalyzer) Language() string {
	return "python"
}

func (pa *PythonAnalyzer) IsAvailable() bool {
	// Python is available if .py files exist or no system check needed
	return true
}

// AnalyzeCodebase analyzes Python files and extracts patterns
func (pa *PythonAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".py") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Error handling
		if count := countOccurrences(text, "try:"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "try/except blocks",
				Description: "Python uses try/except for error handling",
				Example:     "try:\n    # code\nexcept Exception as e:\n    # handle",
				Language:    "python",
				Frequency:   count,
			})
		}

		// Validation
		if count := countOccurrences(text, "isinstance("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "isinstance checks",
				Description: "Type checking using isinstance()",
				Example:     "if isinstance(var, str):",
				Language:    "python",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, " is None") + countOccurrences(text, " is not None"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "None checks",
				Description: "Using 'is None' for null checks",
				Example:     "if value is None:",
				Language:    "python",
				Frequency:   count,
			})
		}

		// Logging
		if count := countOccurrences(text, "print("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "print() logging",
				Description: "Using print() for output",
				Example:     "print(f'Value: {value}')",
				Language:    "python",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, "logging."); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "logging module",
				Description: "Using Python logging module",
				Example:     "logging.info('message')",
				Language:    "python",
				Frequency:   count,
			})
		}

		// Security
		if count := countOccurrences(text, "from cryptography") + countOccurrences(text, "import hashlib"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeSecurity,
				Name:        "cryptography usage",
				Description: "Using cryptography or hashlib for security",
				Example:     "from cryptography.fernet import Fernet",
				Language:    "python",
				Frequency:   count,
			})
		}

		// Naming
		if strings.Count(text, "_") > 10 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeNaming,
				Name:        "snake_case naming",
				Description: "Using snake_case for function and variable names",
				Example:     "def my_function_name():",
				Language:    "python",
				Frequency:   1,
			})
		}

		if count := countOccurrences(text, "class "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeNaming,
				Name:        "PascalCase class names",
				Description: "Using PascalCase for class names",
				Example:     "class MyClass:",
				Language:    "python",
				Frequency:   count,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
