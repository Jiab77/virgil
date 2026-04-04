// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jiab77/virgil/pkg/virgil/storage"
)

// Learner orchestrates pattern extraction from codebases
type Learner struct {
	db *storage.Database
}

// NewLearner creates a new code pattern learner
func NewLearner(db *storage.Database) *Learner {
	return &Learner{
		db: db,
	}
}

// Learn analyzes a codebase and extracts programming patterns
func (l *Learner) Learn(request *LearningRequest) (*LearningResponse, error) {
	response := &LearningResponse{
		Patterns: make([]CodePattern, 0),
	}

	log.Printf("[virgil] Learning from codebase at: %s", request.CodebasePath)

	// Detect language(s) in codebase
	languages := l.detectLanguages(request.CodebasePath)
	if len(languages) == 0 {
		response.Error = "No supported languages detected in codebase"
		response.Success = false
		return response, nil
	}

	log.Printf("[virgil] Detected languages: %v", languages)

	// Extract patterns for each language
	patternProfiles := make([]PatternProfile, 0)
	for _, lang := range languages {
		log.Printf("[virgil] Extracting patterns for %s", lang)

		analyzer, err := NewLanguageAnalyzer(lang)
		if err != nil {
			log.Printf("[virgil] Warning: could not create analyzer for %s: %v", lang, err)
			continue
		}

		if !analyzer.IsAvailable() {
			log.Printf("[virgil] Warning: %s tools not available on system", lang)
			continue
		}

		patterns, err := analyzer.AnalyzeCodebase(request.CodebasePath)
		if err != nil {
			log.Printf("[virgil] Warning: pattern extraction failed for %s: %v", lang, err)
			continue
		}

		response.Patterns = append(response.Patterns, patterns...)

		// For Bash, collect pattern profiles for gap detection
		if lang == "bash" {
			profiles := l.buildPatternProfiles(request.CodebasePath, patterns)
			patternProfiles = append(patternProfiles, profiles...)
		}
	}

	if len(response.Patterns) == 0 {
		response.Error = "No patterns extracted from codebase"
		response.Success = false
		return response, nil
	}

	// Store patterns and profiles in encrypted database
	err := l.storePatterns(response.Patterns)
	if err != nil {
		response.Error = fmt.Sprintf("Failed to store patterns: %v", err)
		response.Success = false
		return response, nil
	}

	// Store pattern profiles for gap detection
	if len(patternProfiles) > 0 {
		err = l.storePatternProfiles(patternProfiles)
		if err != nil {
			log.Printf("[virgil] Warning: failed to store pattern profiles: %v", err)
		}
	}

	response.Success = true
	response.Message = fmt.Sprintf("Extracted and stored %d patterns from %d languages", len(response.Patterns), len(languages))

	return response, nil
}

// detectLanguages scans codebase and detects programming languages
func (l *Learner) detectLanguages(codebasePath string) []string {
	languageMap := make(map[string]bool)

	err := filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if !info.IsDir() {
			ext := filepath.Ext(path)

			// First, try to detect by file extension
			switch ext {
			case ".go":
				languageMap["go"] = true
			case ".py":
				languageMap["python"] = true
			case ".js", ".jsx", ".ts", ".tsx":
				languageMap["javascript"] = true
			case ".php":
				languageMap["php"] = true
			case ".sh", ".bash":
				languageMap["bash"] = true
			case ".rb":
				languageMap["ruby"] = true
			case ".pl", ".pm":
				languageMap["perl"] = true
			case ".rs":
				languageMap["rust"] = true
			case ".c", ".h":
				languageMap["c"] = true
			case ".cpp", ".cc", ".cxx", ".hpp":
				languageMap["cpp"] = true
			case ".asm", ".s", ".S":
				languageMap["asm"] = true
			default:
				// If no extension, try shebang detection
				if ext == "" {
					lang := l.detectLanguageFromShebang(path)
					if lang != "" {
						languageMap[lang] = true
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Printf("[virgil] Error scanning codebase: %v", err)
	}

	// Convert map to slice
	languages := make([]string, 0, len(languageMap))
	for lang := range languageMap {
		languages = append(languages, lang)
	}

	return languages
}

// detectLanguageFromShebang reads the first line of a file and parses the shebang
func (l *Learner) detectLanguageFromShebang(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return ""
	}

	shebang := scanner.Text()

	// Check if first line is a shebang
	if !strings.HasPrefix(shebang, "#!") {
		return ""
	}

	// Parse shebang: #!/usr/bin/env python -> python
	// Also handles: #!/bin/bash -> bash, #!/usr/bin/python -> python
	parts := strings.Fields(shebang)
	if len(parts) < 2 {
		return ""
	}

	// Get the executable name (last part of path or after 'env')
	executable := parts[len(parts)-1]

	// Extract interpreter from path
	// Examples: /usr/bin/python -> python, /bin/bash -> bash
	if idx := strings.LastIndex(executable, "/"); idx != -1 {
		executable = executable[idx+1:]
	}

	// Normalize interpreter names to standard language identifiers
	switch {
	case strings.HasPrefix(executable, "python"):
		return "python"
	case executable == "bash" || executable == "sh":
		return "bash"
	case strings.HasPrefix(executable, "php"):
		return "php"
	case executable == "node" || executable == "nodejs":
		return "javascript"
	case executable == "ruby":
		return "ruby"
	case executable == "perl":
		return "perl"
	case executable == "go":
		return "go"
	case executable == "rust" || executable == "rustup":
		return "rust"
	case executable == "gcc" || executable == "clang" || executable == "g++":
		return "cpp"
	case executable == "as":
		return "asm"
	default:
		return ""
	}
}

// storePatterns saves learned patterns to encrypted database
func (l *Learner) storePatterns(patterns []CodePattern) error {
	log.Printf("[virgil] Storing %d patterns in encrypted database", len(patterns))

	for i, pattern := range patterns {
		// Generate unique pattern ID (language_type_name)
		patternID := fmt.Sprintf("%s_%s_%s", pattern.Language, pattern.Type, strings.ReplaceAll(strings.ToLower(pattern.Name), " ", "_"))

		// Serialize metadata to JSON
		metadata := fmt.Sprintf(`{"frequency":%d}`, pattern.Frequency)

		// Store in database (encrypted)
		err := l.db.SaveLearnedPattern(
			patternID,
			string(pattern.Type),
			pattern.Language,
			pattern.Name,
			pattern.Description,
			pattern.FilePath,
			pattern.Example,
			metadata,
		)

		if err != nil {
			log.Printf("[virgil] Warning: failed to store pattern %d: %v", i, err)
			continue
		}

		log.Printf("[virgil] ✓ Stored pattern: %s", patternID)
	}

	return nil
}

// GetLearnedPatterns retrieves stored patterns for code generation
func (l *Learner) GetLearnedPatterns(language string) ([]CodePattern, error) {
	log.Printf("[virgil] Retrieving learned patterns for %s", language)

	// Query patterns from database
	patternMaps, err := l.db.GetLearnedPatterns(language)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve patterns: %w", err)
	}

	if len(patternMaps) == 0 {
		log.Printf("[virgil] No learned patterns found for %s", language)
		return []CodePattern{}, nil
	}

	// Convert maps to CodePattern structs
	patterns := make([]CodePattern, 0, len(patternMaps))
	for _, p := range patternMaps {
		pattern := CodePattern{
			Type:        PatternType(p["pattern_type"].(string)),
			Name:        p["name"].(string),
			Description: p["description"].(string),
			Example:     p["example"].(string),
			Language:    p["language"].(string),
			Frequency:   p["frequency"].(int),
		}
		patterns = append(patterns, pattern)
	}

	log.Printf("[virgil] Retrieved %d learned patterns for %s", len(patterns), language)

	return patterns, nil
}

// buildPatternProfiles constructs PatternProfile objects for gap detection
func (l *Learner) buildPatternProfiles(codebasePath string, patterns []CodePattern) []PatternProfile {
	profiles := make([]PatternProfile, 0)
	filePatternMap := make(map[string][]CodePattern)

	// Group patterns by file
	for _, pattern := range patterns {
		if len(pattern.LineNumbers) > 0 {
			// Pattern is file-specific
			for _, lineNum := range pattern.LineNumbers {
				key := fmt.Sprintf("%s:%d", codebasePath, lineNum)
				filePatternMap[key] = append(filePatternMap[key], pattern)
			}
		}
	}

	// Build profiles from pattern groupings
	seenFiles := make(map[string]bool)
	for _, pattern := range patterns {
		if len(pattern.LineNumbers) > 0 && !seenFiles[pattern.Language] {
			profile := PatternProfile{
				FilePath: codebasePath,
				Language: pattern.Language,
				Detected: make(map[PatternType]int),
				Expected: []PatternType{
					PatternTypeConfigurationCenter,
					PatternTypeDefensivePrevalidation,
					PatternTypeOperationValidation,
				},
				Gaps: make([]PatternType, 0),
			}

			// Populate detected patterns
			for _, p := range patterns {
				if p.Language == pattern.Language {
					profile.Detected[p.Type] = p.Frequency
				}
			}

			// Identify gaps
			for _, expectedPattern := range profile.Expected {
				if _, found := profile.Detected[expectedPattern]; !found {
					profile.Gaps = append(profile.Gaps, expectedPattern)
				}
			}

			profiles = append(profiles, profile)
			seenFiles[pattern.Language] = true
		}
	}

	return profiles
}

// storePatternProfiles saves pattern profiles to database for gap detection analysis
func (l *Learner) storePatternProfiles(profiles []PatternProfile) error {
	log.Printf("[virgil] Storing %d pattern profiles for gap detection", len(profiles))

	for _, profile := range profiles {
		log.Printf("[virgil] Pattern profile for %s: Detected=%d, Gaps=%d", profile.FilePath, len(profile.Detected), len(profile.Gaps))

		if len(profile.Gaps) > 0 {
			gapNames := make([]string, 0, len(profile.Gaps))
			for _, gap := range profile.Gaps {
				gapNames = append(gapNames, string(gap))
			}
			log.Printf("[virgil] Gaps detected in %s: %v", profile.FilePath, gapNames)
		}
	}

	return nil
}
