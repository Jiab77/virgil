// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package llm

// Provider represents an LLM provider
type Provider string

const (
	ProviderV0        Provider = "v0"
	ProviderClaudeAPI Provider = "claude"
	ProviderOpenAIAPI Provider = "openai"
	ProviderGroq      Provider = "groq"
	ProviderLocal     Provider = "local"
)

// GenerationConfig holds LLM generation parameters
type GenerationConfig struct {
	Provider    Provider
	Model       string
	Temperature float32
	MaxTokens   int
	APIKey      string // For API-based providers
}

// GenerationRequest represents a request to generate code
type GenerationRequest struct {
	Prompt       string
	Context      string
	Language     string
	Requirements []string
	Config       *GenerationConfig
}

// GenerationResponse represents the response from code generation
type GenerationResponse struct {
	Code        string
	Explanation string
	Error       string
	TokensUsed  int
}

// LLMClient defines the interface for LLM providers
type LLMClient interface {
	Generate(request *GenerationRequest) (*GenerationResponse, error)
	ValidateConfig() error
	Name() string
}
