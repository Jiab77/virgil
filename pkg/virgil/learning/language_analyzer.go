// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"fmt"
	"os/exec"
)

// LanguageAnalyzer is the interface all language analyzers must implement
type LanguageAnalyzer interface {
	Language() string
	IsAvailable() bool
	AnalyzeCodebase(codebasePath string) ([]CodePattern, error)
}

// NewLanguageAnalyzer creates the appropriate analyzer for a given language
func NewLanguageAnalyzer(language string) (LanguageAnalyzer, error) {
	switch language {
	case "go":
		return &GoAnalyzer{}, nil
	case "python":
		return &PythonAnalyzer{}, nil
	case "javascript":
		return &JavaScriptAnalyzer{}, nil
	case "php":
		return &PHPAnalyzer{}, nil
	case "bash", "sh":
		return &BashAnalyzer{}, nil
	case "ruby":
		return &RubyAnalyzer{}, nil
	case "perl":
		return &PerlAnalyzer{}, nil
	case "rust":
		return &RustAnalyzer{}, nil
	case "c", "cpp", "c++", "cc":
		return &CCppAnalyzer{}, nil
	case "asm", "assembly":
		return &AsmAnalyzer{}, nil
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}
}

// GoAnalyzer handles Go code analysis
type GoAnalyzer struct{}

func (ga *GoAnalyzer) Language() string {
	return "go"
}

func (ga *GoAnalyzer) IsAvailable() bool {
	_, err := exec.Command("go", "version").Output()
	return err == nil
}

// execCommand is a helper to execute system commands
func execCommand(cmd string, args ...string) (string, error) {
	command := exec.Command(cmd, args...)
	output, err := command.Output()
	return string(output), err
}
