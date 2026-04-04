// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// APIClient represents an HTTP-based LLM client
type APIClient struct {
	provider   Provider
	model      string
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

// NewAPIClient creates a new API-based LLM client
func NewAPIClient(provider Provider, model string) *APIClient {
	apiKey := os.Getenv("LLM_API_KEY")
	if apiKey == "" {
		// Try provider-specific env vars
		switch provider {
		case ProviderV0:
			apiKey = os.Getenv("V0_API_KEY")
		case ProviderClaudeAPI:
			apiKey = os.Getenv("ANTHROPIC_API_KEY")
		case ProviderOpenAIAPI:
			apiKey = os.Getenv("OPENAI_API_KEY")
		case ProviderGroq:
			apiKey = os.Getenv("GROQ_API_KEY")
		}
	}

	client := &APIClient{
		provider:   provider,
		model:      model,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 0}, // No timeout for long generations
	}

	client.setEndpoint()
	return client
}

// setEndpoint sets the correct API endpoint based on provider
func (c *APIClient) setEndpoint() {
	switch c.provider {
	case ProviderV0:
		c.endpoint = "https://api.v0.dev/v1/chat/completions"
	case ProviderClaudeAPI:
		c.endpoint = "https://api.anthropic.com/v1/messages"
	case ProviderOpenAIAPI:
		c.endpoint = "https://api.openai.com/v1/chat/completions"
	case ProviderGroq:
		c.endpoint = "https://api.groq.com/openai/v1/chat/completions"
	}
}

// ValidateConfig validates the client configuration
func (c *APIClient) ValidateConfig() error {
	if c.apiKey == "" {
		return fmt.Errorf("API key not configured for %s provider", c.provider)
	}
	if c.endpoint == "" {
		return fmt.Errorf("unknown provider: %s", c.provider)
	}
	return nil
}

// Name returns the client name
func (c *APIClient) Name() string {
	return string(c.provider)
}

// Generate sends a generation request to the LLM API
func (c *APIClient) Generate(request *GenerationRequest) (*GenerationResponse, error) {
	if err := c.ValidateConfig(); err != nil {
		return nil, err
	}

	// Build the system and user prompts
	systemPrompt := c.buildSystemPrompt(request)
	userPrompt := c.buildUserPrompt(request)

	log.Printf("[virgil] Calling %s API for code generation...", c.provider)

	switch c.provider {
	case ProviderV0:
		return c.generateWithV0(systemPrompt, userPrompt)
	case ProviderClaudeAPI:
		return c.generateWithClaude(systemPrompt, userPrompt)
	case ProviderOpenAIAPI:
		return c.generateWithOpenAI(systemPrompt, userPrompt)
	case ProviderGroq:
		return c.generateWithGroq(systemPrompt, userPrompt)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", c.provider)
	}
}

// generateWithClaude generates code using Claude API
func (c *APIClient) generateWithClaude(system, user string) (*GenerationResponse, error) {
	payload := map[string]interface{}{
		"model":     c.model,
		"max_tokens": 4096,
		"system":    system,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": user,
			},
		},
	}

	respBody, statusCode, err := c.doRequest(payload)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("Claude API error (status %d): %s", statusCode, respBody)
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal([]byte(respBody), &result); err != nil {
		return nil, fmt.Errorf("error parsing Claude response: %w", err)
	}

	if len(result.Content) == 0 {
		return nil, fmt.Errorf("no content in Claude response")
	}

	response := &GenerationResponse{
		Code:       result.Content[0].Text,
		TokensUsed: result.Usage.InputTokens + result.Usage.OutputTokens,
	}

	return response, nil
}

// generateWithV0 generates code using v0 Model API (OpenAI-compatible)
func (c *APIClient) generateWithV0(system, user string) (*GenerationResponse, error) {
	payload := map[string]interface{}{
		"model":       c.model,
		"temperature": 0.7,
		"max_tokens":  4096,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": system,
			},
			{
				"role":    "user",
				"content": user,
			},
		},
	}

	respBody, statusCode, err := c.doRequest(payload)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("v0 API error (status %d): %s", statusCode, respBody)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal([]byte(respBody), &result); err != nil {
		return nil, fmt.Errorf("error parsing v0 response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in v0 response")
	}

	response := &GenerationResponse{
		Code:       result.Choices[0].Message.Content,
		TokensUsed: result.Usage.PromptTokens + result.Usage.CompletionTokens,
	}

	return response, nil
}

// generateWithOpenAI generates code using OpenAI API
func (c *APIClient) generateWithOpenAI(system, user string) (*GenerationResponse, error) {
	payload := map[string]interface{}{
		"model":       c.model,
		"temperature": 0.7,
		"max_tokens":  4096,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": system,
			},
			{
				"role":    "user",
				"content": user,
			},
		},
	}

	respBody, statusCode, err := c.doRequest(payload)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", statusCode, respBody)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal([]byte(respBody), &result); err != nil {
		return nil, fmt.Errorf("error parsing OpenAI response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in OpenAI response")
	}

	response := &GenerationResponse{
		Code:       result.Choices[0].Message.Content,
		TokensUsed: result.Usage.PromptTokens + result.Usage.CompletionTokens,
	}

	return response, nil
}

// generateWithGroq generates code using Groq API
func (c *APIClient) generateWithGroq(system, user string) (*GenerationResponse, error) {
	// Groq uses OpenAI-compatible API
	return c.generateWithOpenAI(system, user)
}

// doRequest performs the HTTP request
func (c *APIClient) doRequest(payload interface{}) (string, int, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", 0, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Set provider-specific headers
	switch c.provider {
	case ProviderClaudeAPI:
		req.Header.Set("x-api-key", c.apiKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	case ProviderV0, ProviderOpenAIAPI, ProviderGroq:
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), resp.StatusCode, nil
}

// buildSystemPrompt builds the system prompt for code generation
func (c *APIClient) buildSystemPrompt(request *GenerationRequest) string {
	return fmt.Sprintf(`You are an expert software engineer creating production-ready code.

Your code MUST:
1. Follow security best practices (OWASP Top 10)
2. Include comprehensive error handling
3. Have proper input validation
4. Include logging and monitoring hooks
5. Be well-documented with comments
6. Follow %s conventions and patterns
7. Handle edge cases gracefully

Generate clean, maintainable, secure code. No placeholders or TODO comments.`, request.Language)
}

// buildUserPrompt builds the user prompt for code generation
func (c *APIClient) buildUserPrompt(request *GenerationRequest) string {
	prompt := fmt.Sprintf(`Generate %s code for the following requirements:

%s

Context:
%s

Requirements:
`, request.Language, request.Prompt, request.Context)

	for i, req := range request.Requirements {
		prompt += fmt.Sprintf("%d. %s\n", i+1, req)
	}

	return prompt
}
