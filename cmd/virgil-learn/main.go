package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jiab77/virgil/pkg/virgil/learning"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("Usage: virgil-learn <codebase-path>")
		fmt.Println("\nExample:")
		fmt.Println("  virgil-learn /path/to/bash/scripts")
		os.Exit(1)
	}

	codebasePath := args[0]

	// Verify path exists
	info, err := os.Stat(codebasePath)
	if err != nil {
		log.Fatalf("Error: codebase path not found: %v", err)
	}

	// Display what we're analyzing
	if info.IsDir() {
		fmt.Printf("\n[virgil-learn] Analyzing Bash scripts in: %s\n", codebasePath)
	} else {
		fmt.Printf("\n[virgil-learn] Analyzing Bash script: %s\n", codebasePath)
	}
	fmt.Println(strings.Repeat("=", 80))

	// Create analyzer - it handles both files and directories
	analyzer := &learning.BashAnalyzer{}
	patterns, err := analyzer.AnalyzeCodebase(codebasePath)
	if err != nil {
		log.Fatalf("Error analyzing: %v", err)
	}

	if len(patterns) == 0 {
		fmt.Println("No patterns detected.")
		os.Exit(0)
	}

	// Group patterns by file
	patternsByFile := make(map[string][]learning.CodePattern)
	for _, pattern := range patterns {
		patternsByFile[pattern.FilePath] = append(patternsByFile[pattern.FilePath], pattern)
	}

	// Sort filenames for consistent output
	var sortedFiles []string
	for filePath := range patternsByFile {
		sortedFiles = append(sortedFiles, filePath)
	}
	sort.Strings(sortedFiles)

	// Display results grouped by file
	fmt.Printf("\n[RESULTS] Detected %d total patterns across %d files\n\n", len(patterns), len(patternsByFile))

	// Aggregate all pattern types for summary
	allPatternTypes := make(map[string]int)

	for _, filePath := range sortedFiles {
		filePatterns := patternsByFile[filePath]
		fmt.Printf("%s:\n", filePath)

		// Group by pattern type within this file
		patternsByType := make(map[string][]learning.CodePattern)
		for _, pattern := range filePatterns {
			typeStr := string(pattern.Type)
			patternsByType[typeStr] = append(patternsByType[typeStr], pattern)
			allPatternTypes[typeStr]++ // Aggregate for summary
		}

		// Sort types for consistent output
		var sortedTypes []string
		for typeStr := range patternsByType {
			sortedTypes = append(sortedTypes, typeStr)
		}
		sort.Strings(sortedTypes)

		// Display patterns for this file
		for _, typeStr := range sortedTypes {
			typePatterns := patternsByType[typeStr]
			for _, p := range typePatterns {
				fmt.Printf("  [%s] %s: %d occurrence(s)\n", typeStr, p.Name, p.Frequency)
				if p.Present && len(p.LineNumbers) > 0 {
					fmt.Printf("    Lines: %v\n", p.LineNumbers)
				}
			}
		}
		fmt.Println()
	}

	// Summary
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\n[SUMMARY]")
	fmt.Printf("Total patterns detected: %d\n", len(patterns))

	// Count by type
	fmt.Println("\nPattern Breakdown:")
	var sortedTypeNames []string
	for typeStr := range allPatternTypes {
		sortedTypeNames = append(sortedTypeNames, typeStr)
	}
	sort.Strings(sortedTypeNames)

	for _, typeStr := range sortedTypeNames {
		count := allPatternTypes[typeStr]
		fmt.Printf("  %s: %d pattern(s)\n", typeStr, count)
	}

	// Check for critical Phase 2 patterns
	fmt.Println("\n[PHASE 2 PATTERNS - Systems Engineering Validation]")
	phase2Patterns := map[string]bool{
		"configuration_center":    false,
		"defensive_prevalidation": false,
		"operation_validation":    false,
	}

	for _, typeStr := range sortedTypeNames {
		if _, exists := phase2Patterns[typeStr]; exists {
			phase2Patterns[typeStr] = true
		}
	}

	for patternName, found := range phase2Patterns {
		status := "✗ NOT DETECTED"
		if found {
			status = "✓ DETECTED"
		}
		fmt.Printf("  %s: %s\n", patternName, status)
	}

	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
