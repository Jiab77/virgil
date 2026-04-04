// Copyright 2026 Jiab77
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents Virgil configuration
type Config struct {
	AugmentStrategy  string   `yaml:"augment_strategy" json:"augment_strategy"`
	Mode             string   `yaml:"mode" json:"mode"`
	ProjectPath      string   `yaml:"-" json:"-"`
	RulesEnabled     []string `yaml:"rules_enabled" json:"rules_enabled"`
	WebSearchEnabled bool     `yaml:"web_search_enabled" json:"web_search_enabled"`
	// LLM Model Selection (per-provider)
	ModelV0      string `yaml:"model_v0" json:"model_v0"`
	ModelClaude  string `yaml:"model_claude" json:"model_claude"`
	ModelOpenAI  string `yaml:"model_openai" json:"model_openai"`
	ModelGroq    string `yaml:"model_groq" json:"model_groq"`
}

// NewConfig creates a new default configuration
func NewConfig() *Config {
	return &Config{
		AugmentStrategy:  "api",        // Default: use external APIs
		Mode:             "plan-first", // Default: assessment gate before generation
		ProjectPath:      ".",
		RulesEnabled:     []string{"owasp", "nist"}, // International standards only by default
		WebSearchEnabled: true,         // Default: web search enabled (ground AI reasoning in research)
		// Default models per provider (can be overridden)
		ModelV0:     "v0-1.5-md",                      // v0: medium model (default), use v0-1.5-lg for complex projects
		ModelClaude: "claude-3-5-sonnet-20241022",   // Claude: latest Sonnet
		ModelOpenAI: "gpt-4-turbo",                   // OpenAI: GPT-4 Turbo
		ModelGroq:   "mixtral-8x7b-32768",           // Groq: Mixtral
	}
}

// LoadConfig loads configuration from .virgil/config.yaml or .virgil/config.json
// Checks for YAML first (preferred), then JSON
// Returns defaults if neither file exists
func LoadConfig(projectPath string) (*Config, error) {
	configDir := filepath.Join(projectPath, ".virgil")
	yamlPath := filepath.Join(configDir, "config.yaml")
	jsonPath := filepath.Join(configDir, "config.json")

	// Try YAML first (preferred format)
	if _, err := os.Stat(yamlPath); err == nil {
		return loadYAMLConfig(yamlPath)
	}

	// Fall back to JSON if YAML doesn't exist
	if _, err := os.Stat(jsonPath); err == nil {
		return loadJSONConfig(jsonPath, projectPath)
	}

	// Neither exists, return defaults
	cfg := NewConfig()
	cfg.ProjectPath = projectPath
	return cfg, nil
}

// loadYAMLConfig loads configuration from YAML file
func loadYAMLConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading YAML config file: %w", err)
	}

	cfg := NewConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parsing YAML config file: %w", err)
	}

	return cfg, nil
}

// loadJSONConfig loads configuration from JSON file
func loadJSONConfig(path string, projectPath string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON config file: %w", err)
	}

	cfg := NewConfig()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("error parsing JSON config file: %w", err)
	}

	cfg.ProjectPath = projectPath
	return cfg, nil
}

// SaveConfig saves configuration to .virgil/config.yaml or .virgil/config.json
// Uses existing format if one already exists, defaults to YAML for new configs
func SaveConfig(cfg *Config, projectPath string) error {
	configDir := filepath.Join(projectPath, ".virgil")
	yamlPath := filepath.Join(configDir, "config.yaml")
	jsonPath := filepath.Join(configDir, "config.json")

	// Ensure .virgil directory exists
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return fmt.Errorf("error creating .virgil directory: %w", err)
	}

	// Determine which format to use
	useJSON := false
	if _, err := os.Stat(jsonPath); err == nil {
		// JSON config already exists, use JSON
		useJSON = true
	} else if _, err := os.Stat(yamlPath); err == nil {
		// YAML config already exists, use YAML
		useJSON = false
	} else {
		// Neither exists, default to YAML
		useJSON = false
	}

	if useJSON {
		return saveJSONConfig(cfg, jsonPath)
	}
	return saveYAMLConfig(cfg, yamlPath)
}

// saveYAMLConfig saves configuration to YAML file
func saveYAMLConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("error marshaling config to YAML: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("error writing YAML config file: %w", err)
	}

	return nil
}

// saveJSONConfig saves configuration to JSON file
func saveJSONConfig(cfg *Config, path string) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling config to JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("error writing JSON config file: %w", err)
	}

	return nil
}

// GetConfigFormat returns the format of the current config file ("yaml", "json", or "none")
func GetConfigFormat(projectPath string) string {
	configDir := filepath.Join(projectPath, ".virgil")
	yamlPath := filepath.Join(configDir, "config.yaml")
	jsonPath := filepath.Join(configDir, "config.json")

	if _, err := os.Stat(yamlPath); err == nil {
		return "yaml"
	}
	if _, err := os.Stat(jsonPath); err == nil {
		return "json"
	}
	return "none"
}
