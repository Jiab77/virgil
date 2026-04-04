// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package blocks

import (
	"github.com/jiab77/virgil/pkg/virgil/verification"
)

// CustomBlock implements user-defined rule checks
type CustomBlock struct{}

func NewCustomBlock() verification.VerificationBlock {
	return &CustomBlock{}
}

func (b *CustomBlock) Name() string {
	return "custom"
}

func (b *CustomBlock) Description() string {
	return "User-defined custom rules"
}

func (b *CustomBlock) Run(targetPath string) (*verification.BlockResult, error) {
	result := &verification.BlockResult{
		BlockName: b.Name(),
		Issues:    []verification.Issue{},
		Status:    "PASS",
		Summary:   "[Phase 2 Stub] Custom rules not yet configured",
	}
	return result, nil
}
