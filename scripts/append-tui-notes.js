import { appendFileSync } from 'fs'

const section = `
---

## TUI Implementation (Session 14)

### Design Philosophy
Serious tools and beautiful colors are not opposites. Color carries meaning, not decoration. Virgil should be a serious tool that happens to look excellent. References: Crush, GoAccess, btop, Bagels, Posting, Superfile.

**Brand colour:** \`#7D56F4\` (purple) — fixed, do not change.

### Output Mode Architecture

\`virgil-learn\` implements three independent, composable output modes:

| Mode | Flag | Dependencies | Description |
|---|---|---|---|
| Plain text | _(default)_ | None | Original output, byte-for-byte unchanged |
| Markdown | \`--markdown\` / \`--md\` | \`charm.land/glamour/v2\` | Glamour ANSI-styled output to stdout |
| TUI | \`--tui\` | bubbletea v2, bubbles v2, lipgloss v2 | Interactive spinner + scrollable viewport |

All flags are orthogonal — \`--tui --markdown\` is a valid and meaningful combination.

### New Files (virgil-learn)

- \`cmd/virgil-learn/renderer.go\` — \`renderPlainText()\`, \`renderMarkdown()\`, \`terminalWidth()\`, \`sortedKeys[K]()\` generic helper
- \`cmd/virgil-learn/tui.go\` — full bubbletea v2 model: \`tuiModel\`, \`Init()\`, \`Update()\`, \`View()\`, \`runTUI()\`, \`groupPatterns()\`

### Charmbracelet Stack Added to go.mod

\`\`\`
charm.land/glamour/v2    v2.0.0
charm.land/bubbletea/v2  v2.0.2
charm.land/bubbles/v2    v2.1.0
charm.land/lipgloss/v2   v2.0.2
\`\`\`

Future additions planned (not yet added):
- \`github.com/charmbracelet/log\` — colored leveled logger, drop-in for standard \`log\`

### Key bubbletea v2 API Differences from v1

- \`View()\` returns \`tea.View\`, not \`string\` — use \`tea.NewView(content)\`
- \`KeyPressMsg\` replaces \`KeyMsg\`
- \`tea.WindowSizeMsg\` for responsive viewport/progress bar sizing
- \`spinner.Tick\` is the tick command; \`spinner.TickMsg\` is the message type
- \`progress.FrameMsg\` must be handled in \`Update()\` to drive animation frames
- Cannot write to stdout in TUI mode — use \`tea.LogToFile("debug.log", "debug")\`

### TUI Layout Principles (for virgil main binary — future)

Derived from design references:
- Bordered panels via \`lipgloss.Border()\` to separate concerns
- Single accent color (\`#7D56F4\`) signals the active panel/tab — not competing colors
- Persistent help bar at the bottom, always same format: \`key action · key action\`
- Color encodes semantic meaning: status, severity, active/inactive
- Tabs built directly with lipgloss — no separate bubbles component needed
- Sidebar layout inspired by Crush — sections for Analysis, Generation, Verification, Audit

### virgil-learn as Test Bed

\`virgil-learn\` is the primary test bed for both the **Learning feature** (multi-language pattern analysis) and the **TUI layer** (bubbletea, bubbles, lipgloss, glamour). All TUI patterns must be validated here before being applied to the main \`virgil\` binary.

### Next TUI Steps (Implementation Order)

1. Run \`go mod tidy\` to resolve \`go.sum\` after manual \`go.mod\` edits
2. Test all four output mode combinations against a real Bash codebase
3. Wire multi-language analyzer routing into \`virgil-learn\` \`main.go\` (language auto-detection from file extension)
4. Add \`ProgressFunc\` callback to language analyzers for directory scan progress bar
5. Add \`github.com/charmbracelet/log\` to replace \`log.Fatalf\` / debug \`fmt.Printf\`
6. Design \`virgil\` main binary TUI (sidebar, panels, tabs) once \`virgil-learn\` TUI is proven
`

appendFileSync('/vercel/share/v0-project/docs/IMPLEMENTATION_NOTES.md', section, 'utf8')
console.log('[v0] TUI section appended to IMPLEMENTATION_NOTES.md')
