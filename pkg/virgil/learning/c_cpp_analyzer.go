// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// CCppAnalyzer analyzes C and C++ codebases for patterns
type CCppAnalyzer struct{}

func (ca *CCppAnalyzer) Language() string {
	return "c++"
}

func (ca *CCppAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes C/C++ files and extracts patterns
func (ca *CCppAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		isCFile := strings.HasSuffix(path, ".c") || strings.HasSuffix(path, ".h")
		isCppFile := strings.HasSuffix(path, ".cpp") || strings.HasSuffix(path, ".cc") ||
			strings.HasSuffix(path, ".cxx") || strings.HasSuffix(path, ".hpp")

		if !isCFile && !isCppFile {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Error handling
		if count := countOccurrences(text, "if ("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "if error checks",
				Description: "Using if statements for error checking",
				Example:     "if (result == NULL) { }",
				Language:    "c++",
				Frequency:   count / 10,
			})
		}

		if isCppFile && countOccurrences(text, "try {") > 0 {
			count := countOccurrences(text, "try {")
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeErrorHandling,
				Name:        "try/catch (C++)",
				Description: "C++ exception handling with try/catch",
				Example:     "try {\n} catch (exception& e) {\n}",
				Language:    "c++",
				Frequency:   count,
			})
		}

		// Validation
		if count := countOccurrences(text, "NULL"); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeValidation,
				Name:        "NULL checks",
				Description: "Checking for NULL pointers",
				Example:     "if (ptr != NULL) { }",
				Language:    "c++",
				Frequency:   count,
			})
		}

		// Logging
		if count := countOccurrences(text, "printf("); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "printf output",
				Description: "Using printf for output",
				Example:     "printf(\"message\\n\")",
				Language:    "c++",
				Frequency:   count,
			})
		}

		if isCppFile && countOccurrences(text, "std::cout") > 0 {
			count := countOccurrences(text, "std::cout")
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeLogging,
				Name:        "std::cout output",
				Description: "Using std::cout for C++ output",
				Example:     "std::cout << \"message\" << std::endl;",
				Language:    "c++",
				Frequency:   count,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
