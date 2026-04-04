// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"

	"github.com/jiab77/virgil/pkg/virgil/cli"
)

func main() {
	rootCmd := cli.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
