// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package generation

import (
	"github.com/jiab77/virgil/pkg/virgil/verification"
)

// CodeGenerationRequest represents a request to generate code
type CodeGenerationRequest struct {
	Description string // User's feature description
	ProjectPath string // Path to project
	Language    string // Target language (go, python, typescript, etc.)
	Context     string // Additional context about the project
}

// CodeGenerationResponse represents the generated code
type CodeGenerationResponse struct {
	Code              string                                 // Generated source code
	Language          string                                 // Language used
	Approved          bool                                   // User approval status
	VerificationResult *verification.AggregatedResult        // Verification results
	Error             string                                 // Error message if any
}

// AugmentationMode determines how code generation should work
type AugmentationMode string

const (
	AugmentationModeAPI      AugmentationMode = "api"      // Use external APIs (Claude/GPT)
	AugmentationModeLearning AugmentationMode = "learning" // Use learned patterns only
)

// GenerationStrategy determines the code generation approach
type GenerationStrategy struct {
	Mode           AugmentationMode // api or learning
	WebSearchEnabled bool            // Whether to use web search context
	Rules          []string         // Verification rules to apply
}
