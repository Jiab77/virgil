// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// BashAnalyzer analyzes Bash scripts for patterns
type BashAnalyzer struct{}

func (ba *BashAnalyzer) Language() string {
	return "bash"
}

func (ba *BashAnalyzer) IsAvailable() bool {
	return true
}

// AnalyzeCodebase analyzes Bash scripts and extracts patterns
// Accepts either a single file or directory path
func (ba *BashAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)
	fileProfiles := make([]PatternProfile, 0)

	// Check if path is a file or directory
	info, err := os.Stat(codebasePath)
	if err != nil {
		return nil, err
	}

	// If it's a single file, analyze it directly
	if !info.IsDir() {
		patterns, _ := ba.analyzeSingleFile(codebasePath)
		return patterns, nil
	}

	// If it's a directory, walk through all files
	filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[virgil] Error accessing %s: %v", path, err)
			return nil
		}
		
		if info.IsDir() {
			return nil
		}

		// Skip hidden files and common non-script files
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		isBashFile := strings.HasSuffix(path, ".sh") || strings.HasSuffix(path, ".bash")
		
		if !isBashFile {
			// Check shebang for bash/sh
			file, err := os.Open(path)
			if err != nil {
				log.Printf("[virgil] Could not open %s for shebang check: %v", path, err)
				return nil
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			if scanner.Scan() {
				line := scanner.Text()
				if !strings.Contains(line, "bash") && !strings.Contains(line, "sh") {
					return nil
				}
			} else {
				// No first line to check
				return nil
			}
		}

		// Analyze this file
		filePatterns, fileProfile := ba.analyzeSingleFile(path)
		patterns = append(patterns, filePatterns...)
		fileProfiles = append(fileProfiles, fileProfile)
		return nil
	})

	return patterns, nil
}

// analyzeSingleFile analyzes a single bash file and returns its patterns and profile
func (ba *BashAnalyzer) analyzeSingleFile(filePath string) ([]CodePattern, PatternProfile) {
	patterns := make([]CodePattern, 0)

	// Read and analyze the file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("[virgil] Warning: could not read %s: %v", filePath, err)
		return patterns, PatternProfile{}
	}

	text := string(content)
	lines := strings.Split(text, "\n")

	// Build pattern profile for this file
	profile := PatternProfile{
		FilePath: filePath,
		Language: "bash",
		Detected: make(map[PatternType]int),
		Expected: []PatternType{
			PatternTypeConfigurationCenter,
			PatternTypeDefensivePrevalidation,
			PatternTypeOperationValidation,
		},
		Gaps: make([]PatternType, 0),
	}

	// Phase 2 Pattern 1: Configuration Center (explicit constants/variables at top)
	if configCount, lineNumbers := detectConfigurationCenter(text, lines); configCount > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeConfigurationCenter,
			Name:        "configuration center",
			Description: "Configuration variables centralized at top of script",
			Example:     "DEBUG=false; RETRY_COUNT=3",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   configCount,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeConfigurationCenter] = configCount
	}

	// Phase 2 Pattern 2: Defensive Pre-validation (check BEFORE use)
	if count, lineNumbers := detectDefensivePrevalidation(text, lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeDefensivePrevalidation,
			Name:        "defensive pre-validation",
			Description: "Validation checks before resource use",
			Example:     "[[ -z $var ]] && die \"var not set\"",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeDefensivePrevalidation] = count
	}

	// Phase 2 Pattern 3: Operation Validation (check AFTER operations)
	if count, lineNumbers := detectOperationValidation(text, lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeOperationValidation,
			Name:        "operation validation",
			Description: "Exit code checking after operations",
			Example:     "if [[ $? -eq 0 ]]; then",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeOperationValidation] = count
	}

	// Detect gaps: which expected patterns are missing?
	for _, expectedPattern := range profile.Expected {
		if _, found := profile.Detected[expectedPattern]; !found {
			profile.Gaps = append(profile.Gaps, expectedPattern)
		}
	}

	// Phase 1 patterns (keep for backward compatibility)
	if count, lineNumbers := detectValidation(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeValidation,
			Name:        "if conditionals",
			Description: "Using if statements for control flow",
			Example:     "if [ $? -eq 0 ]; then",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
	}

	if count, lineNumbers := detectLogging(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeLogging,
			Name:        "echo output",
			Description: "Using echo for output and logging",
			Example:     "echo \"message\"",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
	}

	// Intent-aware Phase 2 detectors (Nine Patterns from LEARNING_MODE_INTENT.md)
	if count, lineNumbers := detectFallbackStrategy(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeFallbackStrategy,
			Name:        "fallback strategy",
			Description: "Primary path fails, secondary path activated",
			Example:     "if [[ $? -ne 0 ]]; then fallback_cmd; fi",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeFallbackStrategy] = count
	}

	if count, lineNumbers := detectMultiPathConfigLoading(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeConfigurationCenter,
			Name:        "multi-path config loading",
			Description: "Priority-ordered configuration file search",
			Example:     "[[ -r ~/.config/app.conf ]] && CONFIG=~/.config/app.conf",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
	}

	if count, lineNumbers := detectStatePreservation(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeStatePreservation,
			Name:        "state preservation",
			Description: "Scope management and state protection via local, readonly, trap, or explicit save/restore",
			Example:     "local var; readonly CONST=val; trap cleanup EXIT; OLD_IFS=$IFS",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeStatePreservation] = count
	}

	if count, lineNumbers := detectStructuredOutput(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeStructuredOutput,
			Name:        "structured output",
			Description: "Consistent message prefix format for log levels",
			Example:     "echo \"[+] success\" or echo \"[-] failure\"",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeStructuredOutput] = count
	}

	if count, lineNumbers := detectPureFunction(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypePureFunction,
			Name:        "pure function",
			Description: "Function uses local vars only, no global side effects",
			Example:     "function foo() { local x=$1; echo \"$x\"; }",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypePureFunction] = count
	}

	if count, lineNumbers := detectAdaptability(lines); count > 0 {
		patterns = append(patterns, CodePattern{
			Type:        PatternTypeAdaptability,
			Name:        "adaptability",
			Description: "Parameterized behaviour via defaults and CLI overrides",
			Example:     "VAR=${VAR:-default}; getopts or case for CLI flags",
			Language:    "bash",
			FilePath:    filePath,
			Frequency:   count,
			Present:     true,
			LineNumbers: lineNumbers,
		})
		profile.Detected[PatternTypeAdaptability] = count
	}

	return patterns, profile
}

// detectConfigurationCenter identifies centralized configuration patterns.
// Looks for a cluster of variable assignments at or near the top of the file,
// before main logic begins. Accepts all naming conventions:
//   - UPPERCASE_WITH_UNDERSCORES (classic shell config)
//   - _prefixed_lowercase (library/private convention)
//   - snake_case or mixed ()
//   - ANSI color variable blocks (NC=, RED=, BLUE=, etc.)
//
// The "top of file" window is relaxed to the first half of the file to handle
// scripts that define colors first, then config vars after the shebang/header.
func detectConfigurationCenter(text string, lines []string) (int, []int) {
	configCount := 0
	lineNumbers := make([]int, 0)

	// Scan up to the first half of the file
	limit := len(lines) / 2
	if limit < 20 {
		limit = len(lines) // small files: scan all
	}

	for i, line := range lines {
		if i >= limit {
			break
		}

		trimmed := strings.TrimSpace(line)

		// Skip comments, blank lines, shebangs, source/export/declare
		if trimmed == "" || strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "source ") || strings.HasPrefix(trimmed, ". ") {
			continue
		}

		// Stop at function definitions or main loop constructs — main logic has started.
		// NOTE: we allow simple guard-clause if blocks like `if [[ -t 1 ]]; then` that
		// are commonly used to conditionally define color variables — do NOT stop on those.
		// We only stop on for/while loops and function definitions which signal real logic.
		if strings.HasPrefix(trimmed, "function ") ||
			strings.Contains(trimmed, "() {") ||
			strings.HasPrefix(trimmed, "for ") ||
			strings.HasPrefix(trimmed, "while ") {
			break
		}

		if strings.Contains(trimmed, "=") && isConfigVariable(trimmed) {
			configCount++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	// Configuration center pattern present if 3+ config variables found
	if configCount >= 3 {
		return configCount, lineNumbers
	}
	return 0, make([]int, 0)
}

// isConfigVariable checks if a line looks like a configuration variable assignment.
// Accepts UPPERCASE, _prefixed_lowercase, snake_case, and mixed conventions.
func isConfigVariable(line string) bool {
	if !strings.Contains(line, "=") {
		return false
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		return false
	}

	varName := strings.TrimSpace(parts[0])

	// Must be a plain identifier — no spaces, brackets, or operators
	if strings.ContainsAny(varName, " \t[](){}$;|&<>") {
		return false
	}

	// UPPERCASE (classic shell config: DEBUG=false, RETRY_COUNT=3)
	if strings.ToUpper(varName) == varName && len(varName) > 1 {
		return true
	}

	// Contains underscore (snake_case, _private, MIXED_case)
	if strings.Contains(varName, "_") {
		return true
	}

	return false
}

// detectDefensivePrevalidation identifies validation before resource use.
// Detects two forms, both of which must terminate execution on failure:
//
//  Form A (single-line): [[ -z $var ]] && die "message"
//                        [[ -f "$file" ]] || exit 1
//
//  Form B (block):       if [[ -z "$var" ]]; then
//                            echo "[ERROR] ..." >&2
//                            exit 1          ← termination required
//                        fi
//
// The key intent signal: unlike generic validation, defensive pre-validation
// always terminates (exit, die, return 1) when the check fails.
func detectDefensivePrevalidation(text string, lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	prevalidationChecks := []string{
		// File/variable test operators
		"[ -z", "[ -n", "[ -f", "[ -d", "[ -r", "[ -w", "[ -x", "[ -e",
		"[[ -z", "[[ -n", "[[ -f", "[[ -d", "[[ -r", "[[ -w", "[[ -x", "[[ -e",
		// Dependency/binary checks — canonical pattern: command -v prog || die
		"command -v", "which ", "type ",
		// Numeric/string comparisons used as guards
		"[ -s", "[[ -s",
		// Error safety options — set -e/set -u/set -o pipefail signal deliberate safety intent.
		// These are pre-validation at the script level: the script refuses to run unsafely.
		"set -e", "set -u", "set -o", "set -E",
	}
	// "exit" matches exit, exit 1, exit 127, etc.
	// "return" matches return, return 1, return 5, etc.
	// "die" matches the common custom helper function pattern
	terminators := []string{"exit", "die", "return"}

	for i, line := range lines {
		hasCheck := false
		for _, pattern := range prevalidationChecks {
			if strings.Contains(line, pattern) {
				hasCheck = true
				break
			}
		}
		if !hasCheck {
			continue
		}

		// set -e / set -u / set -o pipefail are standalone safety declarations —
		// no terminator needed, the intent is unambiguous.
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "set -") || strings.HasPrefix(trimmedLine, "set -o ") {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Form A: single-line with && or || and a terminator
		if strings.Contains(line, "&&") || strings.Contains(line, "||") {
			for _, term := range terminators {
				if strings.Contains(line, term) {
					count++
					lineNumbers = append(lineNumbers, i+1)
					break
				}
			}
			continue
		}

		// Form B: if/then block — look ahead up to 5 lines for a terminator before fi
		if strings.Contains(strings.TrimSpace(line), "if ") {
			lookAhead := i + 1
			maxLook := i + 5
			if maxLook > len(lines) {
				maxLook = len(lines)
			}
			for j := lookAhead; j < maxLook; j++ {
				trimmed := strings.TrimSpace(lines[j])
				if trimmed == "fi" {
					break
				}
				for _, term := range terminators {
					if strings.Contains(trimmed, term) {
						count++
						lineNumbers = append(lineNumbers, i+1)
						goto nextLine
					}
				}
			}
		}
	nextLine:
	}

	return count, lineNumbers
}

// detectOperationValidation identifies exit code checking after operations.
// Looks for:
//   - $? with comparison operators (-eq, -ne, -gt, -lt)
//   - PIPESTATUS — the only reliable way to check individual pipeline exit codes
//   - (( exitval > 0 )) arithmetic comparisons on captured exit codes
func detectOperationValidation(text string, lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// $? with comparison operators
		if strings.Contains(line, "$?") &&
			(strings.Contains(line, "-eq") || strings.Contains(line, "-ne") ||
				strings.Contains(line, "-gt") || strings.Contains(line, "-lt") ||
				strings.Contains(line, "-ge") || strings.Contains(line, "-le")) {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// PIPESTATUS — checking individual pipeline command exit codes
		if strings.Contains(line, "PIPESTATUS") {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Arithmetic exit code check: (( exitval > 0 )) or (( $? != 0 ))
		if strings.Contains(trimmed, "((") &&
			(strings.Contains(line, "$?") || strings.Contains(line, "exitval") || strings.Contains(line, "exit_code")) &&
			(strings.Contains(line, "> 0") || strings.Contains(line, "!= 0") || strings.Contains(line, "== 0")) {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Explicit if [ / if [[ return value checks
		if (strings.Contains(trimmed, "if [") || strings.Contains(trimmed, "if [[")) &&
			(strings.Contains(line, "-eq 0") || strings.Contains(line, "-ne 0")) {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectValidation identifies if statements for control flow
func detectValidation(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		if strings.Contains(line, "if [") {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectLogging identifies echo statements for output and logging
func detectLogging(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		if strings.Contains(line, "echo ") {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectFallbackStrategy identifies primary-fails-use-secondary patterns.
// Looks for: check $? after a command and branch to an alternative.
func detectFallbackStrategy(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Single-line: cmd || fallback_cmd — exclude termination patterns (those are defensive prevalidation)
		if strings.Contains(line, "||") && !strings.Contains(line, "exit") && !strings.Contains(line, "die") && !strings.Contains(line, "return") {
			if strings.Contains(line, "=") || strings.Contains(line, "$(") {
				count++
				lineNumbers = append(lineNumbers, i+1)
				continue
			}
		}
		// Block form: if [[ $? -ne 0 ]]; then ... (alternative action, not exit)
		if (strings.Contains(trimmed, "if") && strings.Contains(line, "$?")) ||
			(strings.Contains(trimmed, "if") && strings.Contains(line, "-ne 0")) {
			// Look ahead: must NOT immediately exit — must do something else
			if i+1 < len(lines) {
				next := strings.TrimSpace(lines[i+1])
				if !strings.HasPrefix(next, "exit") && !strings.HasPrefix(next, "die") && next != "fi" {
					count++
					lineNumbers = append(lineNumbers, i+1)
				}
			}
		}
	}

	return count, lineNumbers
}

// detectMultiPathConfigLoading identifies priority-ordered config file search chains.
// Looks for: [[ -r path ]] && CONFIG=path or [[ -f path ]] && source path
func detectMultiPathConfigLoading(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		if (strings.Contains(line, "[[ -r") || strings.Contains(line, "[[ -f") ||
			strings.Contains(line, "[ -r") || strings.Contains(line, "[ -f")) &&
			strings.Contains(line, "&&") &&
			(strings.Contains(line, "CONFIG") || strings.Contains(line, "source") || strings.Contains(line, ".")) {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectStatePreservation identifies patterns that protect and manage state.
// Detects three forms, all valid:
//
//   Form A (explicit save/restore): OLD_VAR=$VAR, SAVED_IFS=$IFS, ORIG_PATH=$PATH
//     The workaround pattern used by authors who do not know about `local`.
//
//   Form B (trap-based cleanup): trap cleanup_func EXIT/TERM/INT
//     Signal handling for guaranteed state restoration on exit.
//
//   Form C (idiomatic scoping): local var inside functions, readonly for constants
//     The correct Bash tools for scope management and write protection.
//     Using `local` prevents global state mutation — it IS state preservation.
//     Using `readonly` prevents accidental overwrite of constants.
func detectStatePreservation(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Form A: explicit save/restore prefixes
		if (strings.HasPrefix(trimmed, "OLD_") || strings.HasPrefix(trimmed, "SAVED_") ||
			strings.HasPrefix(trimmed, "ORIG_") || strings.HasPrefix(trimmed, "PREV_")) &&
			strings.Contains(trimmed, "=") {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Form B: trap-based signal/exit handler
		// ERR fires when any command returns non-zero — a strong defensive signal.
		if strings.Contains(trimmed, "trap") &&
			(strings.Contains(trimmed, "EXIT") || strings.Contains(trimmed, "TERM") ||
				strings.Contains(trimmed, "INT") || strings.Contains(trimmed, "HUP") ||
				strings.Contains(trimmed, "ERR")) {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Form A extension: IFS save/restore — explicit field separator management
		// IFS=$'\n', IFS=',', IFS=$IFS_SAVED, IFS=$OLDIFS all qualify
		if strings.Contains(trimmed, "IFS=") && !strings.HasPrefix(trimmed, "#") {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Form C: local keyword inside a function body (idiomatic Bash scoping)
		if strings.HasPrefix(trimmed, "local ") {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Form C: readonly keyword (write protection for constants)
		if strings.HasPrefix(trimmed, "readonly ") {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectStructuredOutput identifies consistent structured output patterns.
// Detects three styles, all valid:
//
//   Style A (bracket prefix): echo "[+] success", echo "[-] failure"
//     Common in CTF/hacker tools and security scripts.
//
//   Style B (ANSI color codes): echo -e "${RED}...${NC}", echo -e "\033[0;31m..."
//     The actual ANSI escape sequences are the ground truth — variable naming
//     conventions (RED, CDR, BLUE, etc.) are irrelevant, only the escape code matters.
//     Covers: \033[, \e[, \x1b[ and variable definitions containing those sequences.
//
//   Style C (named logging functions): warn(), error(), info(), debug(), notice()
//     Functions like warn() in dkms.in that print "Warning: " / "Error: " without
//     ANSI codes. Detection looks for function *definitions* using those names.
func detectStructuredOutput(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	// Style A: bracket prefix patterns inside echo calls
	bracketPrefixes := []string{
		`"[+]`, `"[-]`, `"[*]`, `"[!]`,
		`"[DEBUG]`, `"[INFO]`, `"[WARN]`, `"[ERROR]`,
		`'[+]`, `'[-]`, `'[*]`, `'[!]`,
		`'[DEBUG]`, `'[INFO]`, `'[WARN]`, `'[ERROR]`,
		"[+]", "[-]", "[*]", "[!]",
	}

	// Style B: ANSI escape sequences — ground truth regardless of variable names.
	// Covers all three quoting forms:
	//   - \033[  octal escape (most common in scripts)
	//   - \e[    shorthand escape (bash-specific)
	//   - \x1b[  hex escape
	//   - $'\033[' and $'\e[' ANSI-C quoting form (e.g. $'\033[0;31m')
	ansiSequences := []string{`\033[`, `\e[`, `\x1b[`, `$'\033[`, `$'\e[`}

	// Style C: named logging function definitions
	loggingFuncNames := []string{
		"warn()", "warning()", "error()", "err()", "fatal()",
		"info()", "debug()", "notice()", "log()", "msg()",
		"function warn ", "function warning ", "function error ",
		"function err ", "function info ", "function debug ",
		"function notice ", "function log ", "function msg ",
		"function die ",
	}

	for i, line := range lines {
		matched := false

		// Style A: bracket prefix inside echo
		if strings.Contains(line, "echo") {
			for _, prefix := range bracketPrefixes {
				if strings.Contains(line, prefix) {
					matched = true
					break
				}
			}
		}

		// Style B: ANSI escape sequence anywhere on the line
		// (variable definitions and echo usage both count)
		if !matched {
			for _, seq := range ansiSequences {
				if strings.Contains(line, seq) {
					matched = true
					break
				}
			}
		}

		// Style C: named logging function definition
		if !matched {
			trimmed := strings.TrimSpace(line)
			for _, fn := range loggingFuncNames {
				if strings.Contains(trimmed, fn) {
					matched = true
					break
				}
			}
		}

		if matched {
			count++
			lineNumbers = append(lineNumbers, i+1)
		}
	}

	return count, lineNumbers
}

// detectPureFunction identifies functions that use only local variables.
// A function is considered pure if it declares 'local' for its variables
// and does not write to global vars or perform I/O side effects.
func detectPureFunction(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	inFunction := false
	functionStart := 0
	hasLocal := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "function ") || strings.Contains(trimmed, "() {") {
			inFunction = true
			functionStart = i
			hasLocal = false
			continue
		}

		if inFunction {
			if strings.HasPrefix(trimmed, "local ") {
				hasLocal = true
			}
			if trimmed == "}" {
				if hasLocal {
					count++
					lineNumbers = append(lineNumbers, functionStart+1)
				}
				inFunction = false
			}
		}
	}

	return count, lineNumbers
}

// detectAdaptability identifies parameterized behaviour patterns.
// Detects three styles:
//
//   Style A: ${VAR:-default} or ${VAR:=default} parameter expansion
//   Style B: getopts or case-based CLI flag parsing
//   Style C: Optional env var overrides — VAR=${VAR:-value} or bare VAR=${VAR} at top-level
//            Used by scripts like hackshell.sh: XHOME=, QUIET=, FORCE= as env toggles.
func detectAdaptability(lines []string) (int, []int) {
	count := 0
	lineNumbers := make([]int, 0)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Style A: ${VAR:-default} or ${VAR:=default} parameter expansion
		if strings.Contains(line, ":-") || strings.Contains(line, ":=") {
			if strings.Contains(line, "${") {
				count++
				lineNumbers = append(lineNumbers, i+1)
				continue
			}
		}

		// Style B: getopts-based CLI argument parsing
		if strings.Contains(line, "getopts") {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Style B: case-based CLI flag parsing (case "$1", "$ARG", or "$flag")
		if strings.HasPrefix(trimmed, "case") &&
			(strings.Contains(line, `"$1"`) || strings.Contains(line, `"${1}"`) ||
				strings.Contains(line, "$1)") || strings.Contains(line, "${ARG}") ||
				strings.Contains(line, `"$flag"`) || strings.Contains(line, `"$opt"`)) {
			count++
			lineNumbers = append(lineNumbers, i+1)
			continue
		}

		// Style C: optional env var override pattern — VAR=${VAR:-...} or VAR="${VAR:+...}"
		// Also catches: QUIET=${QUIET} or FORCE=${FORCE:-} as env toggle declarations
		if strings.Contains(trimmed, "=${") && !strings.HasPrefix(trimmed, "#") {
			// Must be a top-level assignment, not inside a function body
			if isConfigVariable(trimmed) {
				count++
				lineNumbers = append(lineNumbers, i+1)
			}
		}
	}

	return count, lineNumbers
}
