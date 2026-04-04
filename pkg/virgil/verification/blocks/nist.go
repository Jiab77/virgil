// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package blocks

import (
	"github.com/jiab77/virgil/pkg/virgil/verification"
)

// NISTBlock implements NIST security guidelines checks
type NISTBlock struct{}

func NewNISTBlock() verification.VerificationBlock {
	return &NISTBlock{}
}

func (b *NISTBlock) Name() string {
	return "nist"
}

func (b *NISTBlock) Description() string {
	return "NIST security guidelines"
}

func (b *NISTBlock) Run(targetPath string) (*verification.BlockResult, error) {
	result := &verification.BlockResult{
		BlockName: b.Name(),
		Issues:    []verification.Issue{},
		Status:    "PASS",
		Summary:   "[Phase 2 Stub] NIST verification not yet implemented",
	}
	return result, nil
}
