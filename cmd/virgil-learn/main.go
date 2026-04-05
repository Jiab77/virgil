package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/jiab77/virgil/pkg/virgil/learning"
)

func main() {
	markdownFlag := flag.Bool("markdown", false, "Render output as markdown via Glamour")
	mdFlag       := flag.Bool("md", false, "Alias for --markdown")
	tuiFlag      := flag.Bool("tui", false, "Enable interactive TUI mode (spinner + scrollable viewport)")
	flag.Parse()
	args := flag.Args()

	useMarkdown := *markdownFlag || *mdFlag
	useTUI      := *tuiFlag

	if len(args) < 1 {
		fmt.Println("Usage: virgil-learn [--tui] [--markdown|--md] <codebase-path>")
		fmt.Println("\nExamples:")
		fmt.Println("  virgil-learn /path/to/bash/scripts")
		fmt.Println("  virgil-learn --markdown /path/to/bash/scripts")
		fmt.Println("  virgil-learn --tui /path/to/bash/scripts")
		fmt.Println("  virgil-learn --tui --markdown /path/to/bash/scripts")
		os.Exit(1)
	}

	codebasePath := args[0]

	// Verify path exists
	info, err := os.Stat(codebasePath)
	if err != nil {
		log.Fatalf("Error: codebase path not found: %v", err)
	}

	// TUI mode: hand off to bubbletea immediately; it runs its own analysis internally.
	if useTUI {
		if err := runTUI(codebasePath, useMarkdown); err != nil {
			log.Fatalf("TUI error: %v", err)
		}
		return
	}

	// Plain / markdown mode: print header, run analysis, render output.
	if info.IsDir() {
		fmt.Printf("\n[virgil-learn] Analyzing Bash scripts in: %s\n", codebasePath)
	} else {
		fmt.Printf("\n[virgil-learn] Analyzing Bash script: %s\n", codebasePath)
	}
	fmt.Println(strings.Repeat("=", 80))

	// Create analyzer - it handles both files and directories
	analyzer := &learning.BashAnalyzer{}
	patterns, err := analyzer.AnalyzeCodebase(codebasePath)
	if err != nil {
		log.Fatalf("Error analyzing: %v", err)
	}

	if len(patterns) == 0 {
		fmt.Println("No patterns detected.")
		os.Exit(0)
	}

	// Group patterns by file
	patternsByFile := make(map[string][]learning.CodePattern)
	for _, pattern := range patterns {
		patternsByFile[pattern.FilePath] = append(patternsByFile[pattern.FilePath], pattern)
	}

	// Sort filenames for consistent output
	var sortedFiles []string
	for filePath := range patternsByFile {
		sortedFiles = append(sortedFiles, filePath)
	}
	sort.Strings(sortedFiles)

	// Route to the appropriate renderer
	var output string
	if useMarkdown {
		output = renderMarkdown(codebasePath, patterns, patternsByFile, sortedFiles)
	} else {
		output = renderPlainText(codebasePath, patterns, patternsByFile, sortedFiles)
	}

	fmt.Print(output)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
