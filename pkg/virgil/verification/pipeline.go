// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package verification

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/jiab77/virgil/pkg/virgil/config"
)

// Pipeline orchestrates verification blocks and aggregates results
type Pipeline struct {
	blocks map[string]VerificationBlock
}

// NewPipeline creates a new verification pipeline
func NewPipeline() *Pipeline {
	return &Pipeline{
		blocks: make(map[string]VerificationBlock),
	}
}

// RegisterBlock registers a verification block
func (p *Pipeline) RegisterBlock(block VerificationBlock) {
	p.blocks[block.Name()] = block
}

// Run executes all enabled blocks and returns aggregated results
func (p *Pipeline) Run(targetPath string, enabledRules []string) (*AggregatedResult, error) {
	// Validate target path exists
	if _, err := os.Stat(targetPath); err != nil {
		return nil, fmt.Errorf("target path invalid: %w", err)
	}

	result := &AggregatedResult{
		AssessmentID: fmt.Sprintf("assess_%d", time.Now().Unix()),
		TargetPath:   targetPath,
		RulesEnabled: enabledRules,
		BlockResults: make(map[string]*BlockResult),
		AllIssues:    []Issue{},
	}

	// Filter blocks to only enabled ones
	var blocksToRun []VerificationBlock
	for _, ruleName := range enabledRules {
		if block, exists := p.blocks[ruleName]; exists {
			blocksToRun = append(blocksToRun, block)
		}
	}

	if len(blocksToRun) == 0 {
		result.Status = "PASS"
		result.Summary = "No verification blocks enabled"
		return result, nil
	}

	// Run blocks in parallel using goroutines
	blockResultsChan := make(chan map[string]interface{}, len(blocksToRun))
	var wg sync.WaitGroup

	for _, block := range blocksToRun {
		wg.Add(1)
		go func(b VerificationBlock) {
			defer wg.Done()
			blockResult, err := b.Run(targetPath)
			blockResultsChan <- map[string]interface{}{
				"block": b.Name(),
				"result": blockResult,
				"error": err,
			}
		}(block)
	}

	// Wait for all blocks to complete
	wg.Wait()
	close(blockResultsChan)

	// Collect results
	for res := range blockResultsChan {
		blockName := res["block"].(string)
		if err := res["error"]; err != nil {
			fmt.Printf("Block %s failed: %v\n", blockName, err)
			continue
		}
		blockResult := res["result"].(*BlockResult)
		result.BlockResults[blockName] = blockResult
		result.AllIssues = append(result.AllIssues, blockResult.Issues...)
	}

	// Determine overall status
	result.Status = result.GetWorstSeverity()

	// Generate summary
	criticalCount := 0
	highCount := 0
	for _, issue := range result.AllIssues {
		if issue.Severity == "CRITICAL" {
			criticalCount++
		} else if issue.Severity == "HIGH" {
			highCount++
		}
	}

	if result.Status == "PASS" {
		result.Summary = "All checks passed"
	} else {
		result.Summary = fmt.Sprintf("Found %d critical, %d high severity issues", criticalCount, highCount)
	}

	return result, nil
}

// LoadDefaultBlocks loads all available blocks (real and stubs)
func LoadDefaultBlocks() *Pipeline {
	pipeline := NewPipeline()

	// Register all blocks (real and stubs will be in separate files)
	// This will be populated as we create the block files
	// For now, this is a placeholder

	return pipeline
}
