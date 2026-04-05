package learning

import (
	"fmt"
	"sort"
	"strings"
)

// countOccurrences counts how many times a pattern appears in text
func countOccurrences(text, pattern string) int {
	return strings.Count(text, pattern)
}

// deduplicatePatterns aggregates duplicate patterns and counts their frequency
func deduplicatePatterns(patterns []CodePattern) []CodePattern {
	patternMap := make(map[string]*CodePattern)

	for _, pattern := range patterns {
		// Create unique key from pattern type and name
		key := fmt.Sprintf("%s_%s", pattern.Type, pattern.Name)

		if existing, exists := patternMap[key]; exists {
			// Pattern already exists, increment frequency
			existing.Frequency++
		} else {
			// New pattern, add to map
			p := pattern
			patternMap[key] = &p
		}
	}

	// Convert map back to slice
	result := make([]CodePattern, 0, len(patternMap))
	for _, pattern := range patternMap {
		result = append(result, *pattern)
	}

	return result
}

// hasUppercaseWords checks if a string contains uppercase words (for naming convention detection)
func hasUppercaseWords(text string) bool {
	for _, char := range text {
		if char >= 'A' && char <= 'Z' {
			return true
		}
	}
	return false
}

// GenerateReport produces a language-agnostic markdown report from a learned
// codebase. It operates purely on LearnedCodebook and patternsByFile — both are
// language-agnostic types — so any language analyzer can call it for free.
//
// The report has four sections designed to be consumed by both humans and LLMs:
//
//  1. Codebase Profile Summary  — synthesized paragraph describing the author's style
//  2. Pattern Density           — how consistently each pattern appears (%, not raw count)
//  3. Co-occurrence Matrix      — which patterns appear together (teaches the model to generate them together)
//  4. Gaps                      — Nine Patterns absent from the codebase (honest signal to the generation model)
func GenerateReport(
	codebasePath string,
	language string,
	patternsByFile map[string][]CodePattern,
) string {
	if len(patternsByFile) == 0 {
		return ""
	}

	totalFiles := len(patternsByFile)

	// -- Build per-pattern file sets and co-occurrence data ------------------
	// patternFiles: patternType -> set of files that contain it
	patternFiles := make(map[PatternType]map[string]struct{})
	// coOccurrence: patternA -> patternB -> count of files both appear in
	coOccurrence := make(map[PatternType]map[PatternType]int)

	for filePath, filePatterns := range patternsByFile {
		// Collect unique pattern types for this file
		fileTypes := make(map[PatternType]struct{})
		for _, p := range filePatterns {
			fileTypes[p.Type] = struct{}{}
		}
		// Record file in patternFiles
		for pt := range fileTypes {
			if patternFiles[pt] == nil {
				patternFiles[pt] = make(map[string]struct{})
			}
			patternFiles[pt][filePath] = struct{}{}
		}
		// Record co-occurrences
		typeList := make([]PatternType, 0, len(fileTypes))
		for pt := range fileTypes {
			typeList = append(typeList, pt)
		}
		for i, a := range typeList {
			for _, b := range typeList[i+1:] {
				if coOccurrence[a] == nil {
					coOccurrence[a] = make(map[PatternType]int)
				}
				if coOccurrence[b] == nil {
					coOccurrence[b] = make(map[PatternType]int)
				}
				coOccurrence[a][b]++
				coOccurrence[b][a]++
			}
		}
	}

	// -- Nine Patterns baseline (language-agnostic intent layer) -------------
	ninePatterns := []PatternType{
		PatternTypeDefensivePrevalidation,
		PatternTypeFallbackStrategy,
		PatternTypeConfigurationCenter,
		PatternTypeStatePreservation,
		PatternTypeOperationValidation,
		PatternTypeStructuredOutput,
		PatternTypePureFunction,
		PatternTypeAdaptability,
		PatternTypeEnvironmentAdaptation,
	}

	// -- Section 1: Codebase Profile Summary ---------------------------------
	var sb strings.Builder

	sb.WriteString("## Codebase Profile\n\n")
	sb.WriteString(fmt.Sprintf("**Path:** `%s`  \n", codebasePath))
	sb.WriteString(fmt.Sprintf("**Language:** %s  \n", language))
	sb.WriteString(fmt.Sprintf("**Files analysed:** %d\n\n", totalFiles))

	// Synthesize a summary sentence from what is strongly present (>= 50% of files)
	strongPatterns := []string{}
	for _, pt := range ninePatterns {
		files := patternFiles[pt]
		if len(files)*100/totalFiles >= 50 {
			strongPatterns = append(strongPatterns, string(pt))
		}
	}
	if len(strongPatterns) > 0 {
		sb.WriteString(fmt.Sprintf(
			"> This codebase consistently demonstrates: **%s**.\n",
			strings.Join(strongPatterns, "**, **"),
		))
	} else {
		sb.WriteString("> No dominant patterns detected above 50%% file coverage.\n")
	}
	sb.WriteString("\n")

	// -- Section 2: Pattern Density ------------------------------------------
	sb.WriteString("## Pattern Density\n\n")
	sb.WriteString("| Pattern | Files | Coverage | Signal |\n")
	sb.WriteString("|---|---|---|---|\n")

	// Sort patterns by coverage descending
	type densityRow struct {
		pt       PatternType
		fileCount int
		pct      int
	}
	rows := make([]densityRow, 0, len(patternFiles))
	for _, pt := range ninePatterns {
		fc := len(patternFiles[pt])
		pct := 0
		if totalFiles > 0 {
			pct = fc * 100 / totalFiles
		}
		rows = append(rows, densityRow{pt, fc, pct})
	}
	sort.Slice(rows, func(i, j int) bool {
		return rows[i].pct > rows[j].pct
	})
	for _, r := range rows {
		signal := "weak"
		if r.pct >= 75 {
			signal = "strong"
		} else if r.pct >= 50 {
			signal = "moderate"
		} else if r.pct >= 25 {
			signal = "present"
		}
		sb.WriteString(fmt.Sprintf("| `%s` | %d/%d | %d%% | %s |\n",
			r.pt, r.fileCount, totalFiles, r.pct, signal))
	}
	sb.WriteString("\n")

	// -- Section 3: Co-occurrence Matrix -------------------------------------
	sb.WriteString("## Pattern Co-occurrence\n\n")
	sb.WriteString("Patterns that consistently appear together — the model should generate these as a set.\n\n")

	// Only show pairs that co-occur in >= 25% of files
	threshold := totalFiles / 4
	if threshold < 1 {
		threshold = 1
	}
	type coRow struct {
		a, b  PatternType
		count int
	}
	coRows := []coRow{}
	seen := make(map[string]struct{})
	for _, a := range ninePatterns {
		if coOccurrence[a] == nil {
			continue
		}
		for _, b := range ninePatterns {
			if a == b {
				continue
			}
			key := string(a) + "+" + string(b)
			rev := string(b) + "+" + string(a)
			if _, exists := seen[key]; exists {
				continue
			}
			if _, exists := seen[rev]; exists {
				continue
			}
			c := coOccurrence[a][b]
			if c >= threshold {
				coRows = append(coRows, coRow{a, b, c})
				seen[key] = struct{}{}
			}
		}
	}
	sort.Slice(coRows, func(i, j int) bool {
		return coRows[i].count > coRows[j].count
	})
	if len(coRows) > 0 {
		sb.WriteString("| Pattern A | Pattern B | Co-occurs in |\n")
		sb.WriteString("|---|---|---|\n")
		for _, r := range coRows {
			sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %d file(s) |\n", r.a, r.b, r.count))
		}
	} else {
		sb.WriteString("_No strong co-occurrence pairs detected._\n")
	}
	sb.WriteString("\n")

	// -- Section 4: Gaps -----------------------------------------------------
	sb.WriteString("## Gaps\n\n")
	sb.WriteString("Patterns from the Nine not detected in this codebase.\n")
	sb.WriteString("The generation model will **not** produce these unless explicitly asked.\n\n")

	gaps := []PatternType{}
	for _, pt := range ninePatterns {
		if len(patternFiles[pt]) == 0 {
			gaps = append(gaps, pt)
		}
	}
	if len(gaps) > 0 {
		for _, pt := range gaps {
			sb.WriteString(fmt.Sprintf("- `%s`\n", pt))
		}
	} else {
		sb.WriteString("_All Nine Patterns detected. No gaps._\n")
	}
	sb.WriteString("\n")

	return sb.String()
}
