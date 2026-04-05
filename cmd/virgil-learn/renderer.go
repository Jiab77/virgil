package main

import (
	"fmt"
	"sort"
	"strings"

	"charm.land/glamour/v2"

	"github.com/jiab77/virgil/pkg/virgil/learning"
)

// renderPlainText returns the current plain-text output as a string.
// This is the default renderer — byte-for-byte identical to the original
// fmt.Printf output that was previously inlined in main().
func renderPlainText(
	codebasePath string,
	patterns []learning.CodePattern,
	patternsByFile map[string][]learning.CodePattern,
	sortedFiles []string,
) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\n[RESULTS] Detected %d total patterns across %d files\n\n",
		len(patterns), len(patternsByFile)))

	allPatternTypes := make(map[string]int)

	for _, filePath := range sortedFiles {
		filePatterns := patternsByFile[filePath]
		sb.WriteString(fmt.Sprintf("%s:\n", filePath))

		patternsByType := make(map[string][]learning.CodePattern)
		for _, pattern := range filePatterns {
			typeStr := string(pattern.Type)
			patternsByType[typeStr] = append(patternsByType[typeStr], pattern)
			allPatternTypes[typeStr]++
		}

		sortedTypes := sortedKeys(patternsByType)

		for _, typeStr := range sortedTypes {
			typePatterns := patternsByType[typeStr]
			for _, p := range typePatterns {
				sb.WriteString(fmt.Sprintf("  [%s] %s: %d occurrence(s)\n", typeStr, p.Name, p.Frequency))
				if len(p.LineNumbers) > 0 {
					sb.WriteString(fmt.Sprintf("    Lines: %v\n", p.LineNumbers))
				}
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(strings.Repeat("=", 80) + "\n")
	sb.WriteString("\n[SUMMARY]\n")
	sb.WriteString(fmt.Sprintf("Total patterns detected: %d\n", len(patterns)))

	sb.WriteString("\nPattern Breakdown:\n")
	sortedTypeNames := sortedKeys(allPatternTypes)
	for _, typeStr := range sortedTypeNames {
		count := allPatternTypes[typeStr]
		sb.WriteString(fmt.Sprintf("  %s: %d pattern(s)\n", typeStr, count))
	}

	sb.WriteString("\n[PHASE 2 PATTERNS - Systems Engineering Validation]\n")

	// Build map: pattern type -> list of files that contain it
	phase2Names := []string{"configuration_center", "defensive_prevalidation", "operation_validation"}
	phase2Files := map[string][]string{
		"configuration_center":    {},
		"defensive_prevalidation": {},
		"operation_validation":    {},
	}
	for _, filePath := range sortedFiles {
		for _, p := range patternsByFile[filePath] {
			typeStr := string(p.Type)
			if _, isPhase2 := phase2Files[typeStr]; isPhase2 {
				phase2Files[typeStr] = append(phase2Files[typeStr], filePath)
			}
		}
	}
	for _, patternName := range phase2Names {
		files := phase2Files[patternName]
		if len(files) > 0 {
			sb.WriteString(fmt.Sprintf("  v DETECTED   %s (%d file(s))\n", patternName, len(files)))
			for _, f := range files {
				sb.WriteString(fmt.Sprintf("      - %s\n", f))
			}
		} else {
			sb.WriteString(fmt.Sprintf("  x NOT DETECTED  %s\n", patternName))
		}
	}

	sb.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	return sb.String()
}

// renderMarkdown converts analysis results to markdown and passes them through
// glamour for ANSI-styled terminal output. Falls back to plain text on error.
func renderMarkdown(
	codebasePath string,
	patterns []learning.CodePattern,
	patternsByFile map[string][]learning.CodePattern,
	sortedFiles []string,
) string {
	var sb strings.Builder

	sb.WriteString("# virgil-learn Analysis\n\n")
	sb.WriteString(fmt.Sprintf("**Path:** `%s`  \n", codebasePath))
	sb.WriteString(fmt.Sprintf("**Total patterns:** %d across %d file(s)\n\n",
		len(patterns), len(patternsByFile)))
	sb.WriteString("---\n\n")

	allPatternTypes := make(map[string]int)

	for _, filePath := range sortedFiles {
		filePatterns := patternsByFile[filePath]
		sb.WriteString(fmt.Sprintf("## `%s`\n\n", filePath))

		patternsByType := make(map[string][]learning.CodePattern)
		for _, pattern := range filePatterns {
			typeStr := string(pattern.Type)
			patternsByType[typeStr] = append(patternsByType[typeStr], pattern)
			allPatternTypes[typeStr]++
		}

		sortedTypes := sortedKeys(patternsByType)

		for _, typeStr := range sortedTypes {
			typePatterns := patternsByType[typeStr]
			sb.WriteString(fmt.Sprintf("### %s\n\n", typeStr))
			for _, p := range typePatterns {
				sb.WriteString(fmt.Sprintf("- **%s** — %d occurrence(s)", p.Name, p.Frequency))
				if len(p.LineNumbers) > 0 {
					lineStrs := make([]string, len(p.LineNumbers))
					for i, ln := range p.LineNumbers {
						lineStrs[i] = fmt.Sprintf("%d", ln)
					}
					sb.WriteString(fmt.Sprintf("  \n  Lines: `%s`", strings.Join(lineStrs, ", ")))
				}
				sb.WriteString("\n")
			}
			sb.WriteString("\n")
		}
		sb.WriteString("---\n\n")
	}

	// Summary table
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Pattern | Occurrences |\n")
	sb.WriteString("|---|---|\n")
	sortedTypeNames := sortedKeys(allPatternTypes)
	for _, typeStr := range sortedTypeNames {
		sb.WriteString(fmt.Sprintf("| `%s` | %d |\n", typeStr, allPatternTypes[typeStr]))
	}
	sb.WriteString("\n")

	// Phase 2 validation — per-file breakdown
	sb.WriteString("## Phase 2 Patterns — Systems Engineering Validation\n\n")
	phase2Names := []string{"configuration_center", "defensive_prevalidation", "operation_validation"}
	phase2Files := map[string][]string{
		"configuration_center":    {},
		"defensive_prevalidation": {},
		"operation_validation":    {},
	}
	for _, fp := range sortedFiles {
		for _, p := range patternsByFile[fp] {
			typeStr := string(p.Type)
			if _, isPhase2 := phase2Files[typeStr]; isPhase2 {
				phase2Files[typeStr] = append(phase2Files[typeStr], fp)
			}
		}
	}
	for _, patternName := range phase2Names {
		files := phase2Files[patternName]
		if len(files) > 0 {
			sb.WriteString(fmt.Sprintf("- `%s`: **DETECTED** (%d file(s))\n", patternName, len(files)))
			for _, f := range files {
				sb.WriteString(fmt.Sprintf("  - `%s`\n", f))
			}
		} else {
			sb.WriteString(fmt.Sprintf("- `%s`: **NOT DETECTED**\n", patternName))
		}
	}
	sb.WriteString("\n")

	// Render through glamour
	out, err := glamour.Render(sb.String(), "dark")
	if err != nil {
		// Fall back to plain text if glamour fails to render
		return renderPlainText(codebasePath, patterns, patternsByFile, sortedFiles)
	}

	return out
}

// sortedKeys returns the sorted keys of a map[string]T.
func sortedKeys[T any](m map[string]T) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
