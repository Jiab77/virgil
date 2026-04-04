// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package blocks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jiab77/virgil/pkg/virgil/verification"
)

// OWASPBlock implements OWASP Top 10 security checks
type OWASPBlock struct{}

func NewOWASPBlock() verification.VerificationBlock {
	return &OWASPBlock{}
}

func (b *OWASPBlock) Name() string {
	return "owasp"
}

func (b *OWASPBlock) Description() string {
	return "OWASP Top 10 security best practices"
}

func (b *OWASPBlock) Run(targetPath string) (*verification.BlockResult, error) {
	result := &verification.BlockResult{
		BlockName: b.Name(),
		Issues:    []verification.Issue{},
	}

	// Recursively check all files in target path
	err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files for now
		if info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		fileContent := string(content)

		// Run all OWASP checks on this file
		b.checkHardcodedSecrets(path, fileContent, result)
		b.checkSQLInjection(path, fileContent, result)
		b.checkPathTraversal(path, fileContent, result)
		b.checkXSSPatterns(path, fileContent, result)
		b.checkCommandInjection(path, fileContent, result)

		return nil
	})

	if err != nil {
		return result, err
	}

	if len(result.Issues) == 0 {
		result.Status = "PASS"
		result.Summary = "No OWASP issues detected"
	} else {
		result.Status = "FAIL"
		result.Summary = fmt.Sprintf("Found %d OWASP security issues", len(result.Issues))
	}

	return result, nil
}

// checkHardcodedSecrets looks for hardcoded API keys, passwords, tokens
func (b *OWASPBlock) checkHardcodedSecrets(filePath string, content string, result *verification.BlockResult) {
	patterns := []struct {
		regex       string
		name        string
		suggestion  string
	}{
		{
			regex:      `(?i)(password|passwd|pwd)\s*[:=]\s*["']([^"']{5,})["']`,
			name:       "Hardcoded Password",
			suggestion: "Use environment variables or secure secret management (e.g., vault, KMS)",
		},
		{
			regex:      `(?i)(api[_-]?key|apikey|secret|token)\s*[:=]\s*["']([^"']{10,})["']`,
			name:       "Hardcoded API Key/Secret",
			suggestion: "Move secrets to environment variables or configuration management service",
		},
		{
			regex:      `(?i)aws[_-]?secret[_-]?access[_-]?key\s*[:=]\s*["']([^"']{10,})["']`,
			name:       "Hardcoded AWS Secret",
			suggestion: "Use IAM roles or AWS credentials provider instead of hardcoding",
		},
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.regex)
		matches := re.FindAllStringIndex(content, -1)

		for _, match := range matches {
			lineNum := strings.Count(content[:match[0]], "\n") + 1
			issue := verification.Issue{
				IssueID:     fmt.Sprintf("owasp:hardcoded-secret:%s:%d", filePath, lineNum),
				Rule:        "owasp:hardcoded-secret",
				Severity:    "CRITICAL",
				FilePath:    filePath,
				LineNumber:  lineNum,
				Message:     pattern.name,
				Suggestion:  pattern.suggestion,
				CodeSnippet: extractLineContent(content, lineNum),
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// checkSQLInjection looks for SQL injection vulnerabilities
func (b *OWASPBlock) checkSQLInjection(filePath string, content string, result *verification.BlockResult) {
	patterns := []struct {
		regex      string
		name       string
		suggestion string
	}{
		{
			regex:      `(?i)(query|sql|SELECT|INSERT|UPDATE|DELETE).*\+.*variable`,
			name:       "Potential SQL Injection (string concatenation)",
			suggestion: "Use parameterized queries or prepared statements instead of string concatenation",
		},
		{
			regex:      `(?i)fmt\.Sprintf.*SELECT.*%s`,
			name:       "SQL Injection (fmt.Sprintf with user input)",
			suggestion: "Use database/sql with parameterized queries instead of string formatting",
		},
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.regex)
		matches := re.FindAllStringIndex(content, -1)

		for _, match := range matches {
			lineNum := strings.Count(content[:match[0]], "\n") + 1
			issue := verification.Issue{
				IssueID:     fmt.Sprintf("owasp:sql-injection:%s:%d", filePath, lineNum),
				Rule:        "owasp:sql-injection",
				Severity:    "HIGH",
				FilePath:    filePath,
				LineNumber:  lineNum,
				Message:     pattern.name,
				Suggestion:  pattern.suggestion,
				CodeSnippet: extractLineContent(content, lineNum),
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// checkPathTraversal looks for path traversal vulnerabilities
func (b *OWASPBlock) checkPathTraversal(filePath string, content string, result *verification.BlockResult) {
	patterns := []struct {
		regex      string
		name       string
		suggestion string
	}{
		{
			regex:      `(?i)filepath\.Join.*userInput`,
			name:       "Potential Path Traversal (user-controlled path)",
			suggestion: "Validate and sanitize user input; use filepath.Clean() or restrict to safe directories",
		},
		{
			regex:      `(?i)ioutil\.ReadFile\(userInput\)`,
			name:       "Path Traversal (direct file read from user input)",
			suggestion: "Validate file paths before reading; reject paths with '..' or absolute paths",
		},
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.regex)
		matches := re.FindAllStringIndex(content, -1)

		for _, match := range matches {
			lineNum := strings.Count(content[:match[0]], "\n") + 1
			issue := verification.Issue{
				IssueID:     fmt.Sprintf("owasp:path-traversal:%s:%d", filePath, lineNum),
				Rule:        "owasp:path-traversal",
				Severity:    "HIGH",
				FilePath:    filePath,
				LineNumber:  lineNum,
				Message:     pattern.name,
				Suggestion:  pattern.suggestion,
				CodeSnippet: extractLineContent(content, lineNum),
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// checkXSSPatterns looks for potential XSS vulnerabilities (in web contexts)
func (b *OWASPBlock) checkXSSPatterns(filePath string, content string, result *verification.BlockResult) {
	patterns := []struct {
		regex      string
		name       string
		suggestion string
	}{
		{
			regex:      `(?i)w\.Write.*userInput`,
			name:       "Potential XSS (unescaped user input in response)",
			suggestion: "Use HTML templating with auto-escaping (html/template) or sanitize user input",
		},
		{
			regex:      `(?i)html\.UnescapeString\(userInput\)`,
			name:       "Potential XSS (unescaping user input)",
			suggestion: "Escape user input before rendering in HTML context; only unescape trusted data",
		},
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.regex)
		matches := re.FindAllStringIndex(content, -1)

		for _, match := range matches {
			lineNum := strings.Count(content[:match[0]], "\n") + 1
			issue := verification.Issue{
				IssueID:     fmt.Sprintf("owasp:xss:%s:%d", filePath, lineNum),
				Rule:        "owasp:xss",
				Severity:    "HIGH",
				FilePath:    filePath,
				LineNumber:  lineNum,
				Message:     pattern.name,
				Suggestion:  pattern.suggestion,
				CodeSnippet: extractLineContent(content, lineNum),
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// checkCommandInjection looks for command injection vulnerabilities
func (b *OWASPBlock) checkCommandInjection(filePath string, content string, result *verification.BlockResult) {
	patterns := []struct {
		regex      string
		name       string
		suggestion string
	}{
		{
			regex:      `(?i)exec\.Command\(userInput\)`,
			name:       "Command Injection (user-controlled command)",
			suggestion: "Use exec.Command() with separate arguments array, never pass shell strings with user input",
		},
		{
			regex:      `(?i)os\.Popen.*userInput`,
			name:       "Command Injection (user-controlled shell command)",
			suggestion: "Avoid shell execution with user input; use exec.Command() with Args array instead",
		},
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern.regex)
		matches := re.FindAllStringIndex(content, -1)

		for _, match := range matches {
			lineNum := strings.Count(content[:match[0]], "\n") + 1
			issue := verification.Issue{
				IssueID:     fmt.Sprintf("owasp:command-injection:%s:%d", filePath, lineNum),
				Rule:        "owasp:command-injection",
				Severity:    "CRITICAL",
				FilePath:    filePath,
				LineNumber:  lineNum,
				Message:     pattern.name,
				Suggestion:  pattern.suggestion,
				CodeSnippet: extractLineContent(content, lineNum),
			}
			result.Issues = append(result.Issues, issue)
		}
	}
}

// extractLineContent extracts the content of a specific line from text
func extractLineContent(content string, lineNum int) string {
	lines := strings.Split(content, "\n")
	if lineNum > 0 && lineNum <= len(lines) {
		return strings.TrimSpace(lines[lineNum-1])
	}
	return ""
}
