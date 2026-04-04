// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// BashAnalyzer analyzes Bash scripts for patterns
type BashAnalyzer struct{}

func (ba *BashAnalyzer) Language() string {
	return "bash"
}

func (ba *BashAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes Bash scripts and extracts patterns
// Accepts either a single file or directory path
func (ba *BashAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)
	fileProfiles := make([]PatternProfile, 0)

	// Check if path is a file or directory
	info, err := os.Stat(codebasePath)
	if err != nil {
		return nil, err
	}

	// If it's a single file, analyze it directly
	if !info.IsDir() {
		patterns, _ := ba.analyzeSingleFile(codebasePath)
		return patterns, nil
	}

	// If it's a directory, walk through all files
	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[virgil] Error accessing %s: %v", path, err)
			return nil
		}
		
		if info.IsDir() {
			return nil
		}

		// Skip hidden files and common non-script files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		isBashFile := strings.HasSuffix(path, ".sh") || strings.HasSuffix(path, ".bash")
		
		if !isBashFile {
			// Check shebang for bash/sh
			file, err := os.Open(path)
			if err != nil {
				log.Printf("[virgil] Could not open %s for shebang check: %v", path, err)
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			if scanner.Scan() {
				line := scanner.Text()
				if !strings.Contains(line, "bash") && !strings.Contains(line, "sh") {
					return nil
				}
			} else {
				// No first line to check
				return nil
			}
		}

		// Analyze this file
		filePatterns, fileProfile := ba.analyzeSingleFile(path)
		patterns = append(patterns, filePatterns...)
		fileProfiles = append(fileProfiles, fileProfile)
		return nil
	})

	return patterns, nil
}

// analyzeSingleFile analyzes a single bash file and returns its patterns and profile
func (ba *BashAnalyzer) analyzeSingleFile(filePath string) ([]CodePattern, PatternProfile) {
	patterns := make([]CodePattern, 0)

	// Read and analyze the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("[virgil] Warning: could not read %s: %v", filePath, err)
		return patterns, PatternProfile{}
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Build pattern profile for this file
	profile := PatternProfile{
		FilePath: filePath,
		Language: "bash",
		Detected: make(map[PatternType]int),
		Expected: []PatternType{
			PatternTypeConfigurationCenter,
			PatternTypeDefensivePrevalidation,
			PatternTypeOperationValidation,
		},
		Gaps: make([]PatternType, 0),
	}

	// Phase 2 Pattern 1: Configuration Center (explicit constants/variables at top)
	if configCount, lineNumbers := detectConfigurationCenter(text, lines); configCount > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeConfigurationCenter,
			Name:        "configuration center",
			Description: "Configuration variables centralized at top of script",
			Example:     "DEBUG=false; RETRY_COUNT=3",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   configCount,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeConfigurationCenter] = configCount
	}

	// Phase 2 Pattern 2: Defensive Pre-validation (check BEFORE use)
	if count, lineNumbers := detectDefensivePrevalidation(text, lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeDefensivePrevalidation,
			Name:        "defensive pre-validation",
			Description: "Validation checks before resource use",
			Example:     "[[ -z $var ]] && die \"var not set\"",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeDefensivePrevalidation] = count
	}

	// Phase 2 Pattern 3: Operation Validation (check AFTER operations)
	if count, lineNumbers := detectOperationValidation(text, lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeOperationValidation,
			Name:        "operation validation",
			Description: "Exit code checking after operations",
			Example:     "if [[ $? -eq 0 ]]; then",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeOperationValidation] = count
	}

	// Detect gaps: which expected patterns are missing?
	for _, expectedPattern := range profile.Expected {
		if _, found := profile.Detected[expectedPattern]; !found {
			profile.Gaps = append(profile.Gaps, expectedPattern)
		}
	}

	// Phase 1 patterns (keep for backward compatibility)
	if count, lineNumbers := detectValidation(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeValidation,
			Name:        "if conditionals",
			Description: "Using if statements for control flow",
			Example:     "if [ $? -eq 0 ]; then",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			LineNumbers: lineNumbers,
		})
	}

	if count, lineNumbers := detectLogging(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeLogging,
			Name:        "echo output",
			Description: "Using echo for output and logging",
			Example:     "echo \"message\"",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			LineNumbers: lineNumbers,
		})
	}

	return patterns, profile
}

// detectConfigurationCenter identifies centralized configuration patterns
// Looks for explicit constants/variables at top of file before main logic
// Returns count and line numbers where config variables are found
func detectConfigurationCenter(text string, lines []string) (int, []int) {
	configCount := 0
	lineNumbers := make([]int, 0)
	inConfigSection := true

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Stop looking once we hit function definitions or main logic
		if strings.HasPrefix(trimmed, "function ") || strings.HasPrefix(trimmed, "main(") {
			inConfigSection = false
		}

		// Count variable assignments that look like config (CAPS_WITH_UNDERSCORE or snake_case)
		if inConfigSection && i < len(lines)/3 { // Config typically in first 1/3 of file
			if strings.Contains(trimmed, "=") && !strings.HasPrefix(trimmed, "#") && !strings.HasPrefix(trimmed, "//") {
				if isConfigVariable(trimmed) {
					configCount++
					lineNumbers = append(lineNumbers, i+1)
				}
			}
		}
	}

	// Configuration center pattern present if 3+ config variables at top
	if configCount >= 3 {
		return configCount, lineNumbers
	}
	return 0, make([]int, 0)
}

// isConfigVariable checks if a line looks like a configuration variable
func isConfigVariable(line string) bool {
	// Look for patterns like VAR=value or VAR="${DEFAULT:-value}"
	if !strings.Contains(line, "=") {
		return false
	}

	parts := strings.Split(line, "=")
	if len(parts) < 2 {
		return false
	}

	varName := strings.TrimSpace(parts[0])

	// Configuration variables are typically UPPERCASE or contain underscores
	if strings.ToUpper(varName) == varName && len(varName) > 1 {
		return true
	}

	// Or snake_case with some capitals
	if strings.Contains(varName, "_") {
		return true
	}

	return false
}

// detectDefensivePrevalidation identifies validation before resource use
// Looks for checks like [[ -z $var ]], [ -f "$file" ], etc.
func detectDefensivePrevalidation(text string, lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)
	prevalidationPatterns := []string{
		"[ -z",
		"[ -n",
		"[ -f",
		"[ -d",
		"[ -r",
		"[ -w",
		"[ -x",
		"[ -e",
		"[[ -z",
		"[[ -n",
		"[[ -f",
		"[[ -d",
		"[[ -r",
		"[[ -w",
		"[[ -x",
		"[[ -e",
	}

	for i, line := range lines {
		for _, pattern := range prevalidationPatterns {
			if strings.Contains(line, pattern) && strings.Contains(line, "&&") && strings.Contains(line, "die") {
				// Pattern: [[ -X $var ]] && die "message"
				count++
				lineNumbers = append(lineNumbers, i+1)
				break
			}
		}
	}

	return count, lineNumbers
}

// detectOperationValidation identifies exit code checking after operations
// Looks for patterns like: if [[ $? -eq 0 ]], checking return values, etc.
func detectOperationValidation(text string, lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Look for exit code checks
		if strings.Contains(line, "$?") && (strings.Contains(line, "-eq") || strings.Contains(line, "-ne")) {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}

		// Look for explicit return value checks
		if (strings.Contains(trimmed, "if [") || strings.Contains(trimmed, "if [[")) &&
			(strings.Contains(line, "-eq 0") || strings.Contains(line, "-ne 0")) {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectValidation identifies if statements for control flow
func detectValidation(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		if strings.Contains(line, "if [") {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectLogging identifies echo statements for output and logging
func detectLogging(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		if strings.Contains(line, "echo ") {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}
