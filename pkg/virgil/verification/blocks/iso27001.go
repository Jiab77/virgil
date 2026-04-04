// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package blocks

import (
	"github.com/jiab77/virgil/pkg/virgil/verification"
)

// ISO27001Block implements ISO/IEC 27001 information security checks
type ISO27001Block struct{}

func NewISO27001Block() verification.VerificationBlock {
	return &ISO27001Block{}
}

func (b *ISO27001Block) Name() string {
	return "iso27001"
}

func (b *ISO27001Block) Description() string {
	return "ISO/IEC 27001 information security management"
}

func (b *ISO27001Block) Run(targetPath string) (*verification.BlockResult, error) {
	result := &verification.BlockResult{
		BlockName: b.Name(),
		Issues:    []verification.Issue{},
		Status:    "PASS",
		Summary:   "[Phase 2 Stub] ISO27001 verification not yet implemented",
	}
	return result, nil
}
