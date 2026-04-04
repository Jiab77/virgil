// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"os"
	"path/filepath"
	"strings"
)

// AsmAnalyzer analyzes Assembly codebases for patterns
type AsmAnalyzer struct{}

func (aa *AsmAnalyzer) Language() string {
	return "asm"
}

func (aa *AsmAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes Assembly files and extracts patterns
func (aa *AsmAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		isAsmFile := strings.HasSuffix(path, ".asm") || strings.HasSuffix(path, ".s") || strings.HasSuffix(path, ".S")
		if !isAsmFile {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		text := string(content)

		// Control flow
		if count := countOccurrences(text, "jmp "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeStructure,
				Name:        "jump instructions",
				Description: "Using jump instructions for control flow",
				Example:     "jmp label",
				Language:    "asm",
				Frequency:   count,
			})
		}

		if count := countOccurrences(text, "je ") + countOccurrences(text, "jne "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeStructure,
				Name:        "conditional jumps",
				Description: "Using conditional jump instructions",
				Example:     "je label\njne label",
				Language:    "asm",
				Frequency:   count,
			})
		}

		// Memory operations
		if count := countOccurrences(text, "mov "); count > 0 {
			patterns = append(patterns, CodePattern{
				Type:        PatternTypeStructure,
				Name:        "move instructions",
				Description: "Using mov for data transfer",
				Example:     "mov rax, rbx",
				Language:    "asm",
				Frequency:   count / 10,
			})
		}

		return nil
	})

	patterns = deduplicatePatterns(patterns)
	return patterns, nil
}
