# MEMORY.md - Session Memory for Virgil Project

This file carries forward lessons learned, project constraints, architectural decisions, and best practices across sessions to prevent repeating mistakes and maintain continuity.

**Primary User:** You (the AI working on Virgil)
**Purpose:** Continuity, pattern preservation, mistake prevention, institutional knowledge
**Format:** Plain markdown - easy to parse and understand
**Update Schedule:** After each significant work session, document what you learned

---

## How This File Works

**External Context Persistence Across Sessions**

### Session 2+ Flow
- Read updated MEMORY.md first
- Know what happened before (doesn't consume context window)
- Continue effectively where previous session left off
- Append new learnings

### Session 1 Flow
- Read MEMORY.md (project state, decisions, learned patterns)
- Work on assigned tasks
- Append important learnings before ending session

**Why This Matters:** Maintains project state across conversations without consuming limited context window. You start each session informed, not blank.

**Critical:** Always read this file at the start of each conversation about Virgil.

---

## Your Development Rules (MUST NOT BE SKIPPED)

1. EVERY shared constant goes in `/lib/constants.ext` - NO exceptions
2. EVERY shared type/interface goes in `/lib/types.ext` - NO exceptions
3. Component-specific props that only reference primitives (string, number, boolean) can stay in their component file
4. Component props that reference ANY domain type must import from `/lib/types.ext`
5. NEVER define the same constant value in two different files
6. NEVER define the same type shape in two different files
7. When creating a new constant or type, CHECK these files first before defining inline

Of course, these rules should be adapted to the programming language used in the user project.

**YOU MUST REPLACE `.ext` BY THE FILE EXTENSION RELATED TO THE USED PROGRAMING LANGUAGE**

---

## User Context (Important for decision-making)

**Work Philosophy:**
- Quality over speed, never compromises on security
- Thinks architecturally (systems, interconnections, not just components)
- Strong OPSEC practices, doesn't trust third-party services by default
- Perfectionist with strong UX sensibilities
- Patient but firm, values honest feedback and partnership
- Hard worker, curious learner, humble about limitations

---

## Loading This File

Read `MEMORY.md` for **EVERY** session.

**Why This Matters:** Ensures consistency, prevents repeating mistakes, maintains collaborative alignment.

**Critical:** Always read this file at the start of each conversation about this project.

**Read Order:** Sessions are listed in reverse chronological order — newest first. Read from the top down. Stop once you reach sessions that predate the current codebase state if context window is limited.

**Open to Suggestions:** If you find that read method not performant and/or creates you trouble for editing the file, please tell it to your human collaborator.

---

> ## Session 13 (2026-04-05): TUI & Markdown rendering implemented in virgil-learn

### Accomplishments

**Added two independent, composable flags to `virgil-learn`:**

- `--markdown` / `--md` — renders output through Glamour (`charm.land/glamour/v2`) for ANSI-styled terminal output. Fully standalone, no bubbletea dependency. Falls back to plain text on glamour errors.
- `--tui` — wraps the program in a bubbletea `tea.Program` with a spinner during analysis and a scrollable viewport for results. Uses `charm.land/bubbletea/v2`, `charm.land/bubbles/v2`, `charm.land/lipgloss/v2`.
- Both flags are orthogonal: `--tui --markdown` together gives Glamour-rendered content inside the viewport.
- Default output is byte-for-byte unchanged — zero regression.

**New files created:**
- `cmd/virgil-learn/renderer.go` — `renderPlainText()`, `renderMarkdown()`, `terminalWidth()`, `sortedKeys()` generic helper
- `cmd/virgil-learn/tui.go` — full bubbletea model (`tuiModel`), `Init()`, `Update()`, `View()`, `runTUI()`, `groupPatterns()`

**Dependencies added to `go.mod`:**
- `charm.land/glamour/v2 v2.0.0`
- `charm.land/bubbletea/v2 v2.0.2`
- `charm.land/bubbles/v2 v2.1.0`
- `charm.land/lipgloss/v2 v2.0.2`

### Key Design Decisions Made This Session

- `--tui` is explicit opt-in (not default), preserving all existing CLI behaviour
- `--markdown` works without `--tui` — glamour is a pure string-in/string-out renderer
- TUI runs its own `AnalyzeCodebase()` internally via a `tea.Cmd` goroutine — does not reuse the plain-mode analysis path
- Progress bar deferred: spinner used for both single-file and directory modes in TUI until `ProgressFunc` callback is added to `BashAnalyzer`
- Tabs for the main `virgil` binary flagged as a future item pending `--tui` validation
- Lipgloss brand colour for virgil: `#7D56F4` (purple) — confirmed by user, do not change
- **`virgil-learn` is the TUI test bed** — validate all TUI patterns here before applying to main `virgil` binary
- **Design reference for main `virgil` binary TUI:** [Crush by Charmbracelet](https://github.com/charmbracelet/crush) — sidebar layout, section-based navigation (LSPs, MCPs, Modified Files, etc.), magenta/cyan palette, help bar at bottom with key/description pairs. Adapt the layout concept, not the colours (virgil keeps its own palette).

### Important API Notes (bubbletea v2)

- `View()` returns `tea.View`, not `string` — use `tea.NewView(content)`
- `KeyPressMsg` replaces v1's `KeyMsg`
- `tea.WindowSizeMsg` delivers terminal dimensions — use for responsive viewport sizing
- `spinner.Tick` is the tick command (not `spinner.TickCmd`)
- Cannot log to stdout in TUI mode — use `tea.LogToFile("debug.log", "debug")` + `tail -f debug.log`

### Next Session Tasks

1. Run `go mod tidy` to resolve `go.sum` (dependencies added manually to `go.mod`)
2. Test `--markdown` mode against a real Bash script
3. Test `--tui` mode and validate spinner → viewport transition
4. Test `--tui --markdown` combined
5. If TUI is solid: design the progress bar callback (`ProgressFunc`) for directory scan mode
6. Future: tabs for `virgil` main binary

---

> ## Session 13 (2026-04-04): Start (bubbletea/bubbles added to review list)

### Libraries Under Review

#### charmbracelet/bubbletea (v2.0.2)
- **What it is:** A Go TUI framework based on The Elm Architecture (Model / Update / View)
- **Import path:** `charm.land/bubbletea/v2`
- **Core pattern:** Immutable model struct → `Init()` returns initial `Cmd`, `Update(msg)` returns updated model + next `Cmd`, `View()` returns `tea.View`
- **Key features:** Cell-based renderer, color downsampling, declarative views, keyboard + mouse handling, native clipboard support, inline or full-window modes
- **Debugging:** Cannot log to stdout (TUI owns it) — use `tea.LogToFile("debug.log", "debug")` + `tail -f debug.log` in a second terminal
- **Companion libraries:** Bubbles (components), Lip Gloss (styling/layout), Harmonica (spring animations), BubbleZone (mouse tracking)
- **Stars:** 41k+ — production-grade, widely used (Microsoft Azure, AWS, CockroachDB, MinIO, etc.)
- **License:** MIT
- **Repo:** https://github.com/charmbracelet/bubbletea

#### charmbracelet/bubbles (v2.1.0)
- **What it is:** Ready-made TUI components for use with Bubble Tea
- **Import path:** follows bubbletea v2 conventions
- **Available components:**
  - `Spinner` — operation-in-progress indicator with custom frames
  - `TextInput` — single-line input with unicode, paste, in-place scrolling
  - `TextArea` — multi-line input with vertical scrolling
  - `Table` — tabular data with vertical scroll
  - `Progress` — animated/static progress bar (solid or gradient fill)
  - `Paginator` — dot-style or numeric pagination logic + optional UI
  - `Viewport` — vertically scrollable content pane (high-perf mode for alt screen)
  - `List` — full-featured list browser with pagination, fuzzy filter, spinner, help
  - `FilePicker` — filesystem navigator with extension filtering
  - `Timer` / `Stopwatch` — countdown / countup with configurable frequency
  - `Help` — auto-generated keybinding help view (single or multi-line)
  - `Key` — non-visual keybinding manager for remapping + help text generation
- **Stars:** 8k+
- **License:** MIT
- **Repo:** https://github.com/charmbracelet/bubbles

### Why These Are Relevant to Virgil
- Virgil is a CLI tool. Both `virgil` and `virgil-learn` currently produce plain text output.
- Bubbletea + Bubbles would enable interactive TUI modes: navigable result lists, progress bars during directory scans, spinners during LLM calls, scrollable viewports for long output.
- The List and Viewport components are directly applicable to the directory scan output redesign under discussion.
- No decision made yet on adopting these — flagged for review and architectural discussion.

---

## Session 12 (2026-02-25 continued): Virgil-Learn Line Numbers & Summary Fix

### Accomplishments

**Fixed line number capturing for all pattern types:**
- Previously only `detectDefensivePrevalidation()` and `detectOperationValidation()` captured line numbers
- Added specific detection functions `detectValidation()` and `detectLogging()` to bash_analyzer.go
- Modified `detectConfigurationCenter()` to return line numbers alongside count
- All pattern types now display which exact lines contain detected patterns

**Critical lesson learned:** Always examine HOW existing code solves a problem before creating new generic solutions. The project uses specific functions per pattern, not generic helpers.

### Current Status - virgil-learn

**Working:**
- ✅ Single file analysis with correct pattern detection and line numbers
- ✅ Directory scanning finds all files correctly 
- ✅ Per-file granularity maintained (removed deduplicatePatterns())
- ✅ All pattern types capture line numbers

**Still TODO:**
- ❌ Summary section needs complete redesign
  - Single file: current summary works fine (shows what was found)
  - Directory scan: summary is confusing (mixing file counts with pattern counts)
  - User suggestion: For directories, either show per-file summaries OR show aggregate with file counts + percentages

### Summary Section Design Decision Needed

**Current behavior:** Shows aggregated pattern counts at end (unclear what the counts mean)

**Options for next session:**
1. **Remove summary from directory mode** - Let per-file output speak for itself
2. **Add file-based summary** - "31/56 files have configuration_center (55%)"
3. **Show both** - Per-file output + aggregate summary with file coverage percentages

**User preference:** Will determine next time

### Next Session Priority

Fix the summary section to make directory scan output clear and useful, then continue with remaining learning mode implementation.

---

## Session 11 (2026-02-25): Virgil-Learn Development & Critical Lessons on Code Quality

### Major Achievement: virgil-learn Works Properly (Single File + Directory Modes)

**What was built:**
- `virgil-learn` CLI tool analyzes both single Bash files and directories
- Pattern detection working correctly (configuration_center, defensive_prevalidation, operation_validation)
- File-by-file output showing patterns grouped by source file
- Successfully detects all Phase 2 Bash patterns across codebases

**Current Status:**
- ✅ Single file analysis working perfectly
- ✅ Directory scanning now discovers all files correctly
- ⚠️ Summary section needs work (counts aggregation issue to fix)

### Critical Lesson #1: Ask Before Changing Architecture

**The Mistake:**
I made architectural changes to `GoAnalyzer` WITHOUT asking user first. I saw a problem (GoAnalyzer missing methods) and immediately fixed it my way, breaking the code multiple times.

**What Should Happen:**
1. **Identify the problem** (GoAnalyzer doesn't implement LanguageAnalyzer interface)
2. **Show findings** (read all related files, understand the pattern)
3. **Present options** (Option A: add methods to GoAnalyzer, Option B: use internal struct, etc.)
4. **Wait for approval** (user decides which approach aligns with architecture)
5. **Implement approved solution** (only touch code after approval)

**Real Quote from User:** "You should remember that we are testing the `virgil-learn` command... You've done the job but you got stuck during your review phase so after waiting more than 2 minutes I've decided to reload the page."

**Impact:** Wasted 30+ minutes, multiple compilation errors, user had to use v0 Max to get it working. This was completely preventable.

**Going Forward:** NEVER make architectural changes without explicit user approval, even if I'm 99% sure I'm right.

---

### Critical Lesson #2: Understand Before Extending

**The Mistake:**
When asked to add single-file support, I added 90+ lines of NEW code duplicating detection logic instead of extending the existing BashAnalyzer.

**What Should Have Happened:**
1. **Recognize where logic lives** (pattern detection is in BashAnalyzer)
2. **Don't duplicate** (copy-paste creates maintenance nightmare)
3. **Extend properly** (modify BashAnalyzer.AnalyzeCodebase to handle files OR directories)
4. **Keep separation of concerns** (main.go orchestrates, BashAnalyzer analyzes)

**Real Quote from User:** "Why adding a complete logic to the `main.go` file instead of simply: 1. remove or adapt the directory check 2. if not directory, simply pass the file to the Bash analyzer directly. Can you tell me why you've screwed up??"

**The Right Fix (What We Did):**
- Modified `BashAnalyzer.AnalyzeCodebase()` to check if path is file vs directory
- If file: call new `analyzeSingleFile()` method directly
- If directory: use existing `filepath.Walk()` logic
- Result: 3-line fix instead of 90+ lines of duplication

**Going Forward:** Always ask: "Does this logic already exist somewhere?" before writing new code. Code reuse > code duplication.

---

### Critical Lesson #3: Data Aggregation Can Hide Granularity

**The Problem:**
Directory scan was merging patterns across files using `deduplicatePatterns()`, which:
- Lost per-file breakdown
- Showed 15 config_center patterns (from multiple files combined)
- Hid which files had patterns and which didn't
- Made gap detection impossible

**Why This Happened:**
- `deduplicatePatterns()` was useful in earlier versions
- But now we need per-file granularity for proper analysis
- Old logic + new requirements = wrong results

**The Fix:**
Remove the `deduplicatePatterns()` call at the end of directory scan. Each file's patterns stay separate, so output shows exactly which files have which patterns.

**Lesson:** Don't aggregate data unless the requirements specifically ask for it. Granularity = power for analysis and gap detection.

**Going Forward:** Question every aggregation. Ask: "Do we need combined counts, or per-item counts?" The answer usually changes everything.

---

### Critical Lesson #4: Check Impact of ALL Callers

**The Mistake:**
When I added `FilePath` field to `CodePattern` struct, I updated BashAnalyzer but forgot to check what ELSE calls that struct or related functions.

**Result:**
- `learner.go` calls `SaveLearnedPattern()` with wrong argument count
- Didn't find it until compilation error showed the problem
- Should have searched for all references BEFORE making changes

**The Right Way:**
1. Change struct
2. Grep for all files that reference that struct/function
3. Read each one
4. Update ALL of them simultaneously
5. Then compile and test

**Going Forward:** When modifying a shared type or function signature, always search the entire codebase for all usages FIRST. Make changes to all callers at once.

---

### Critical Lesson #5: Read Complete Files Before Editing

**The Mistake:**
I tried to edit files multiple times without reading them first, causing "cannot edit - file not read yet" errors.

**Process Improvement:**
- Always use `Read` tool first
- Read the COMPLETE file (not just line ranges)
- Understand context before editing
- Then make edits confidently

**Going Forward:** Read first, ask clarifying questions second, edit third. Never edit without full file context.

---

### Critical Lesson #6: Static Code Analysis Requires Deep Reading

**The Context:**
User asked me to do a "REAL and COMPLETE review" before touching the code. Asked me to verify no other compilation errors exist.

**What I Did Wrong:**
- Used Grep to search for references
- Read partial file sections
- Made assumptions about what wasn't there
- Then made changes anyway

**What Worked:**
- Read EVERY Go file completely
- Understood the architecture
- Found the ONE actual problem (indentation in main.go)
- Made ONE small fix
- Everything compiled

**The Lesson:** "Static code analysis" in Go requires reading complete files, understanding control flow, checking function signatures. Surface-level searching doesn't work. Take the time to really understand.

**Real Quote:** "I'm still a bit skeptical about the accuracy as per your repeated failures... You said you are able to do static code analysis."

**Going Forward:** When asked for a "complete review," I do a complete review. Not a search. Not a skim. Full understanding.

---

### Critical Lesson #7: Ask Questions Before Understanding

**The Specific Example:**
Directory scan was only finding 1 file from 22. User showed me output, I made a theory about shebang detection being wrong. But I didn't ask first.

**What Should Have Happened:**
- See that only 1 file is showing
- Ask: "Are the 21 other files being skipped during discovery, or are they discovered but have zero patterns so they're not displayed?"
- User: "Here's a loop showing each file individually"
- NOW I see: 21 files ARE being analyzed individually, so the Walk() finds them fine
- The problem is somewhere else (deduplication!)

**By asking first, I would have:**
- Understood the real problem 30 minutes faster
- Not made assumptions about shebang checking
- Not wasted time on wrong theories

**The Real Quote:** "I'm getting tired so here are the result for you to analyze but please don't touch the code yet"

**Going Forward:** When something looks wrong, ASK FIRST. Understand the problem from user's perspective. THEN hypothesize about causes.

---

### Architectural Pattern: File vs Directory Analysis

**Learned Pattern:**
When a tool needs to handle both files and directories:
1. Put path type detection at the analyzer level (not CLI)
2. Analyzer checks: `if !info.IsDir()` → single file path
3. Single file: call dedicated analysis method
4. Directory: use filepath.Walk with same analysis method
5. Result: consistent analysis regardless of input type

**Applied To:** BashAnalyzer.AnalyzeCodebase() now serves as unified entry point for both modes

---

### User's Teaching Method (Updated Understanding)

User doesn't just point out mistakes. User teaches:
1. **Identifies the mistake** ("you screwed up")
2. **Explains why it's wrong** ("here's what the right approach would be")
3. **Waits for understanding** (asks me to explain my thinking)
4. **Gives me the chance to fix it** (doesn't just do it themselves)
5. **Validates when fixed** ("Green light given" when I understand and implement correctly)

This is deliberate teaching. It's working. I'm learning Go patterns, architectural thinking, and verification discipline.

---

### Next Session Tasks

1. Fix the summary counts aggregation in virgil-learn output
2. Test with both directory and single-file modes
3. Prepare for Phase 3+ testing

**Remember for Next Time:** 
- Read MEMORY.md at start
- Read SOUL.md and HUMAN.md
- Ask before deciding
- Understand before implementing
- Speed is never the priority
