// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

// PatternType represents different code patterns to learn
type PatternType string

const (
	// Phase 1: Syntax Patterns (Original)
	PatternTypeStructure      PatternType = "structure"      // File organization, package layout
	PatternTypeErrorHandling  PatternType = "error"          // Error handling patterns
	PatternTypeSecurity       PatternType = "security"       // Security-related patterns
	PatternTypeValidation     PatternType = "validation"     // Input validation patterns
	PatternTypeLogging        PatternType = "logging"        // Logging and observability
	PatternTypeNaming         PatternType = "naming"         // Naming conventions

	// Phase 2: Systems Engineering Patterns (New)
	PatternTypeDefensivePrevalidation PatternType = "defensive_prevalidation" // Check BEFORE use
	PatternTypeFallbackStrategy       PatternType = "fallback_strategy"       // Primary fails → secondary ready
	PatternTypeConfigurationCenter    PatternType = "configuration_center"    // Centralized, explicit config
	PatternTypeStatePreservation      PatternType = "state_preservation"      // Save originals before mutation
	PatternTypeOperationValidation    PatternType = "operation_validation"    // Check AFTER operation with context
	PatternTypeStructuredOutput       PatternType = "structured_output"       // Consistent format for messages
	PatternTypePureFunction           PatternType = "pure_function"           // Input→process→output, no side effects
	PatternTypeAdaptability           PatternType = "adaptability"            // Parameterize, don't hardcode
	PatternTypeEnvironmentAdaptation  PatternType = "environment_adaptation"  // Detect runtime, adjust behavior
)

// CodePattern represents an extracted code pattern
type CodePattern struct {
	Type        PatternType `json:"type"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Example     string      `json:"example"`
	Language    string      `json:"language"`
	FilePath    string      `json:"file_path"`    // Path to the file where pattern was found
	Frequency   int         `json:"frequency"` // How often this pattern appears in codebase
	Present     bool        `json:"present"`   // Whether pattern was detected in file
	LineNumbers []int       `json:"line_numbers"` // Where pattern occurs
}

// PatternProfile represents the pattern presence matrix for a file/module
type PatternProfile struct {
	FilePath string                     `json:"file_path"`
	Language string                     `json:"language"`
	Detected map[PatternType]int        `json:"detected"`   // PatternType → count of occurrences
	Expected []PatternType             `json:"expected"`   // Patterns baseline expects to find
	Gaps     []PatternType             `json:"gaps"`       // Missing patterns from baseline
}

// LearnedCodebook represents the collected patterns from a codebase
type LearnedCodebook struct {
	Language          string                 `json:"language"`
	Patterns          []CodePattern          `json:"patterns"`
	Summary           string                 `json:"summary"`
	PatternProfiles   []PatternProfile       `json:"pattern_profiles"`  // Per-file pattern presence
	EstablishedTiers  map[string][]PatternType `json:"established_tiers"` // Which patterns typically present at each tier
}

// LearningRequest represents a request to learn from codebase
type LearningRequest struct {
	CodebasePath string // Path to codebase to learn from
	Languages    []string // Languages to analyze (leave empty for auto-detect)
}

// LearningResponse represents the result of learning
type LearningResponse struct {
	Success   bool
	Message   string
	Patterns  []CodePattern
	Error     string
}
