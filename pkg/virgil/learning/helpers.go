package learning

import (
	"fmt"
	"strings"
)

// countOccurrences counts how many times a pattern appears in text
func countOccurrences(text, pattern string) int {
	return strings.Count(text, pattern)
}

// deduplicatePatterns aggregates duplicate patterns and counts their frequency
func deduplicatePatterns(patterns []CodePattern) []CodePattern {
	patternMap := make(map[string]*CodePattern)

	for _, pattern := range patterns {
		// Create unique key from pattern type and name
		key := fmt.Sprintf("%s_%s", pattern.Type, pattern.Name)

		if existing, exists := patternMap[key]; exists {
			// Pattern already exists, increment frequency
			existing.Frequency++
		} else {
			// New pattern, add to map
			p := pattern
			patternMap[key] = &p
		}
	}

	// Convert map back to slice
	result := make([]CodePattern, 0, len(patternMap))
	for _, pattern := range patternMap {
		result = append(result, *pattern)
	}

	return result
}

// hasUppercaseWords checks if a string contains uppercase words (for naming convention detection)
func hasUppercaseWords(text string) bool {
	for _, char := range text {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}
