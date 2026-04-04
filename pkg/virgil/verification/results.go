// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package verification

// Issue represents a single security/compliance issue found
type Issue struct {
	IssueID      string // Unique identifier for this issue
	Rule         string // Which rule found this (e.g., "owasp:hardcoded-secrets")
	Severity     string // CRITICAL, HIGH, MEDIUM, LOW
	FilePath     string
	LineNumber   int
	ColumnNumber int
	Message      string // Description of the issue
	Suggestion   string // How to fix it
	CodeSnippet  string // The problematic code
}

// BlockResult represents the result of running one verification block
type BlockResult struct {
	BlockName string  // Name of the block (e.g., "owasp")
	Status    string  // PASS, FAIL, WARNING
	Issues    []Issue
	Summary   string // Human-readable summary
}

// VerificationBlock is the interface all compliance rule blocks must implement
type VerificationBlock interface {
	// Name returns the block name (e.g., "owasp", "nist", "gdpr")
	Name() string

	// Description returns a description of what this block checks
	Description() string

	// Run executes the verification checks on the given target path
	// Returns a BlockResult with any issues found
	Run(targetPath string) (*BlockResult, error)
}

// AggregatedResult combines results from multiple blocks
type AggregatedResult struct {
	AssessmentID string
	TargetPath   string
	RulesEnabled []string
	Status       string // PASS, FAIL, WARNING (based on worst issue severity)
	BlockResults map[string]*BlockResult
	AllIssues    []Issue
	Summary      string
	Context      map[string]interface{} // Additional context like web search results
}

// GetWorstSeverity returns the worst (highest priority) severity from all issues
func (ar *AggregatedResult) GetWorstSeverity() string {
	severityOrder := map[string]int{
		"CRITICAL": 4,
		"HIGH":     3,
		"MEDIUM":   2,
		"LOW":      1,
	}

	worst := ""
	worstScore := 0

	for _, issue := range ar.AllIssues {
		score := severityOrder[issue.Severity]
		if score > worstScore {
			worstScore = score
			worst = issue.Severity
		}
	}

	if worst == "" {
		return "PASS"
	}

	if worstScore >= 3 {
		return "FAIL"
	}
	return "WARNING"
}
