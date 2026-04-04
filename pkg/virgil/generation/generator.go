// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package generation

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jiab77/virgil/pkg/virgil/config"
	"github.com/jiab77/virgil/pkg/virgil/learning"
	"github.com/jiab77/virgil/pkg/virgil/llm"
	"github.com/jiab77/virgil/pkg/virgil/storage"
	"github.com/jiab77/virgil/pkg/virgil/verification"
)

// Generator handles code generation with verification gates
type Generator struct {
	db   *storage.Database
	cfg  *config.Config
	mode AugmentationMode
}

// NewGenerator creates a new code generator
func NewGenerator(db *storage.Database, cfg *config.Config) *Generator {
	mode := AugmentationModeLearning
	if cfg.AugmentStrategy == "api" {
		mode = AugmentationModeAPI
	}

	return &Generator{
		db:   db,
		cfg:  cfg,
		mode: mode,
	}
}

// GenerateCode orchestrates the code generation workflow:
// 1. Assessment phase (verification before generation)
// 2. User approval gate
// 3. Code generation
// 4. Post-generation verification
func (g *Generator) GenerateCode(request *CodeGenerationRequest) (*CodeGenerationResponse, error) {
	response := &CodeGenerationResponse{
		Language: request.Language,
	}

	// Step 1: Assessment phase - verify requirements before generation
	log.Printf("[virgil] Assessment phase: analyzing requirements...")

	// Run verification pipeline against description
	assessmentResult, err := verification.RunPipeline(
		request.Description,
		request.ProjectPath,
		g.cfg,
		g.db,
	)
	if err != nil {
		response.Error = fmt.Sprintf("Assessment failed: %v", err)
		return response, err
	}

	response.VerificationResult = assessmentResult

	// Display assessment results
	g.displayAssessmentResults(assessmentResult)

	// Step 2: User approval gate (in actual implementation, this prompts user)
	// For now, we log that approval is needed
	log.Printf("[virgil] Waiting for user approval to proceed with code generation...")

	// Step 3: Code generation (placeholder for Phase 3 implementation)
	// This will call either:
	// - Local ONNX model (learning mode)
	// - External API (Claude/GPT via Vercel AI Gateway)
	response.Code = g.generateCodeFromAssessment(request, assessmentResult)
	response.Approved = false // Will be set to true after user approval

	return response, nil
}

// generateCodeFromAssessment generates code based on assessment results
// Uses either API mode (Claude/OpenAI/Groq) or learning mode (local patterns)
func (g *Generator) generateCodeFromAssessment(request *CodeGenerationRequest, assessment *verification.AggregatedResult) string {
	switch g.mode {
	case AugmentationModeAPI:
		return g.generateWithAPI(request, assessment)
	case AugmentationModeLearning:
		return g.generateWithLearning(request, assessment)
	default:
		return g.generateStubCode(request, assessment)
	}
}

// generateWithAPI generates code using external LLM API
func (g *Generator) generateWithAPI(request *CodeGenerationRequest, assessment *verification.AggregatedResult) string {
	log.Printf("[virgil] Code generation mode: API (external LLM)")

	// Determine which LLM provider to use
	// Default to Claude if available, fall back to others
	var client llm.LLMClient
	var provider llm.Provider

	// Try providers in order of preference
	// v0 is first (specialized for web code generation)
	// Then fallback to Claude, OpenAI, Groq
	for _, p := range []llm.Provider{llm.ProviderV0, llm.ProviderClaudeAPI, llm.ProviderOpenAIAPI, llm.ProviderGroq} {
		apiClient := llm.NewAPIClient(p, g.getModelForProvider(p))
		if err := apiClient.ValidateConfig(); err == nil {
			client = apiClient
			provider = p
			break
		}
	}

	if client == nil {
		log.Printf("[virgil] No LLM API configured. Falling back to stub code.")
		return g.generateStubCode(request, assessment)
	}

	log.Printf("[virgil] Using %s for code generation", client.Name())

	// Build generation request
	requirements := []string{
		fmt.Sprintf("Language: %s", request.Language),
		"Follow the assessment findings below",
		"Implement proper error handling",
		"Include input validation",
		"Add logging statements",
	}

	genRequest := &llm.GenerationRequest{
		Prompt:       request.Description,
		Context:      g.buildAssessmentContext(assessment),
		Language:     request.Language,
		Requirements: requirements,
		Config: &llm.GenerationConfig{
			Provider:    provider,
			Model:       g.getModelForProvider(provider),
			Temperature: 0.7,
			MaxTokens:   4096,
		},
	}

	// Generate code
	response, err := client.Generate(genRequest)
	if err != nil {
		log.Printf("[virgil] Code generation failed: %v. Using stub code.", err)
		return g.generateStubCode(request, assessment)
	}

	log.Printf("[virgil] Generated code using %s (%d tokens)", client.Name(), response.TokensUsed)
	return response.Code
}

// generateWithLearning generates code using learned patterns from codebase
func (g *Generator) generateWithLearning(request *CodeGenerationRequest, assessment *verification.AggregatedResult) string {
	log.Printf("[virgil] Code generation mode: LEARNING (local patterns from your codebase)")

	// Initialize learner to retrieve learned patterns
	learner := learning.NewLearner(g.db)

	// Get learned patterns for this language
	patterns, err := learner.GetLearnedPatterns(request.Language)
	if err != nil || len(patterns) == 0 {
		log.Printf("[virgil] No learned patterns found for %s. Generating stub code.", request.Language)
		return g.generateStubCode(request, assessment)
	}

	log.Printf("[virgil] Found %d learned patterns for %s", len(patterns), request.Language)

	// Build prompt context from learned patterns
	patternContext := g.buildPatternContext(patterns)

	// Generate system prompt that embeds learned patterns
	systemPrompt := fmt.Sprintf(`You are an expert code generator trained on the user's codebase patterns.

%s

Requirements:
- Follow the patterns shown above
- Match the user's coding style and conventions
- Include appropriate error handling
- Add validation as shown in patterns
- Use logging patterns from the codebase
- Follow security practices observed in patterns
- Match the naming conventions used

Verification Checks:
%s

Generate production-ready code that feels native to this codebase.`, 
		patternContext,
		g.buildAssessmentContext(assessment),
	)

	// Build user prompt
	userPrompt := fmt.Sprintf(`Generate %s code for:

%s

Follow all learned patterns and verification requirements.`, 
		request.Language, 
		request.Description,
	)

	// Since learning mode uses local patterns, still use API for actual generation
	// but with learned patterns as context (API mode with learning-augmented prompts)
	// This ensures consistency while leveraging user's patterns

	// Try to generate with API first, using learned patterns as context
	for _, p := range []llm.Provider{llm.ProviderV0, llm.ProviderClaudeAPI, llm.ProviderOpenAIAPI, llm.ProviderGroq} {
		apiClient := llm.NewAPIClient(p, g.getModelForProvider(p))
		if err := apiClient.ValidateConfig(); err == nil {
			log.Printf("[virgil] Using %s for learning-augmented code generation", p)

			genRequest := &llm.GenerationRequest{
				Prompt:   userPrompt,
				Context:  systemPrompt,
				Language: request.Language,
				Requirements: []string{
					"Use learned patterns",
					"Match codebase style",
					"Follow verification rules",
				},
				Config: &llm.GenerationConfig{
					Provider:    p,
					Model:       g.getModelForProvider(p),
					Temperature: 0.6, // Lower temp for pattern adherence
					MaxTokens:   4096,
				},
			}

			response, err := apiClient.Generate(genRequest)
			if err == nil {
				log.Printf("[virgil] Generated code with learned patterns (%d tokens)", response.TokensUsed)
				return response.Code
			}
			log.Printf("[virgil] Generation with %s failed, trying next provider", p)
		}
	}

	// Fallback: generate stub with pattern context
	log.Printf("[virgil] All API providers failed. Generating stub with learned pattern context.")
	return g.generateStubCode(request, assessment)
}

// buildPatternContext formats learned patterns as context for LLM
func (g *Generator) buildPatternContext(patterns []learning.CodePattern) string {
	context := "## Learned Code Patterns From Your Codebase\n\n"

	// Group patterns by type
	patternsByType := make(map[learning.PatternType][]learning.CodePattern)
	for _, p := range patterns {
		patternsByType[p.Type] = append(patternsByType[p.Type], p)
	}

	// Format each pattern type
	for _, pType := range []learning.PatternType{
		learning.PatternTypeErrorHandling,
		learning.PatternTypeValidation,
		learning.PatternTypeLogging,
		learning.PatternTypeSecurity,
		learning.PatternTypeNaming,
		learning.PatternTypeStructure,
	} {
		pats := patternsByType[pType]
		if len(pats) == 0 {
			continue
		}

		context += fmt.Sprintf("### %s Patterns (Frequency score: ", strings.Title(string(pType)))
		for i, p := range pats {
			context += fmt.Sprintf("- **%s** (seen %d times): %s\n", p.Name, p.Frequency, p.Description)
			if p.Example != "" {
				context += fmt.Sprintf("  Example: `%s`\n", p.Example)
			}
			if i >= 2 { // Limit to 3 top patterns per type
				if len(pats) > 3 {
					context += fmt.Sprintf("  ... and %d more %s patterns\n", len(pats)-3, pType)
				}
				break
			}
		}
		context += "\n"
	}

	return context
}

// generateStubCode generates stub code when LLM is unavailable
func (g *Generator) generateStubCode(request *CodeGenerationRequest, assessment *verification.AggregatedResult) string {
	return fmt.Sprintf(`// Generated code for: %s
// Language: %s
// Augmentation Mode: %s
// Generated: Automatic code generation
// 
// IMPORTANT: This is stub code. Please implement based on requirements.

package main

import (
	"fmt"
	"log"
)

// Assessment Results:
// - Passed Checks: %d
// - Failed Checks: %d
// - Warnings: %d

func main() {
	log.Println("Stub implementation - replace with actual code")
	fmt.Println("TODO: Implement %s")
}

// TODO: Implement error handling
// TODO: Implement input validation
// TODO: Implement logging

// Assessment Findings:
// %s
`,
		request.Description,
		request.Language,
		g.mode,
		len(assessment.PassedChecks),
		len(assessment.FailedChecks),
		len(assessment.WarningChecks),
		request.Description,
		g.buildAssessmentContext(assessment),
	)
}

// buildAssessmentContext formats assessment results as context for LLM
func (g *Generator) buildAssessmentContext(assessment *verification.AggregatedResult) string {
	context := "Assessment Results:\n"

	if len(assessment.PassedChecks) > 0 {
		context += fmt.Sprintf("\nPassed (%d):\n", len(assessment.PassedChecks))
		for _, check := range assessment.PassedChecks {
			context += fmt.Sprintf("  - %s: %s\n", check.Block, check.Message)
		}
	}

	if len(assessment.FailedChecks) > 0 {
		context += fmt.Sprintf("\nFailed (%d):\n", len(assessment.FailedChecks))
		for _, check := range assessment.FailedChecks {
			context += fmt.Sprintf("  - %s: %s\n", check.Block, check.Message)
		}
	}

	if len(assessment.WarningChecks) > 0 {
		context += fmt.Sprintf("\nWarnings (%d):\n", len(assessment.WarningChecks))
		for _, check := range assessment.WarningChecks {
			context += fmt.Sprintf("  - %s: %s\n", check.Block, check.Message)
		}
	}

	return context
}

// getModelForProvider returns the model for a given provider
// Reads from config if set, otherwise uses hardcoded defaults
func (g *Generator) getModelForProvider(provider llm.Provider) string {
	switch provider {
	case llm.ProviderV0:
		if g.cfg.ModelV0 != "" {
			return g.cfg.ModelV0
		}
		return "v0-1.5-md"
	case llm.ProviderClaudeAPI:
		if g.cfg.ModelClaude != "" {
			return g.cfg.ModelClaude
		}
		return "claude-3-5-sonnet-20241022"
	case llm.ProviderOpenAIAPI:
		if g.cfg.ModelOpenAI != "" {
			return g.cfg.ModelOpenAI
		}
		return "gpt-4-turbo"
	case llm.ProviderGroq:
		if g.cfg.ModelGroq != "" {
			return g.cfg.ModelGroq
		}
		return "mixtral-8x7b-32768"
	default:
		return "unknown"
	}
}

// displayAssessmentResults displays verification results to user
func (g *Generator) displayAssessmentResults(result *verification.AggregatedResult) {
	fmt.Println("\n=== Assessment Results ===\n")

	if len(result.PassedChecks) > 0 {
		fmt.Printf("✓ Passed Checks (%d):\n", len(result.PassedChecks))
		for _, check := range result.PassedChecks {
			fmt.Printf("  - %s: %s\n", check.Block, check.Message)
		}
		fmt.Println()
	}

	if len(result.FailedChecks) > 0 {
		fmt.Printf("✗ Failed Checks (%d):\n", len(result.FailedChecks))
		for _, check := range result.FailedChecks {
			fmt.Printf("  - %s: %s\n", check.Block, check.Message)
		}
		fmt.Println()
	}

	if len(result.WarningChecks) > 0 {
		fmt.Printf("⚠ Warnings (%d):\n", len(result.WarningChecks))
		for _, check := range result.WarningChecks {
			fmt.Printf("  - %s: %s\n", check.Block, check.Message)
		}
		fmt.Println()
	}

	// Display web search context if available
	if searchResults, ok := result.Context["web_search_results"]; ok {
		fmt.Printf("Research Context: %d search results available\n", len(searchResults.([]interface{})))
	}

	fmt.Println("Assessment complete. Ready for code generation after approval.")
}
