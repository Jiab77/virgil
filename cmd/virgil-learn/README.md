# virgil-learn CLI Tool

A command-line tool to analyze codebases and extract systems engineering patterns.
Currently supports Bash. Additional language analyzers are planned (Go, Python, JavaScript, Rust, C/C++, Ruby, PHP, Perl, ASM).

## Building

From the project root:

```bash
go build -o bin/virgil-learn ./cmd/virgil-learn
```

Or from this directory:

```bash
go build -o virgil-learn .
```

## Usage

```
virgil-learn [flags] <path>
```

`<path>` can be a single file or a directory. When a directory is provided, all supported files are analyzed recursively.

### Flags

| Flag | Alias | Description |
|---|---|---|
| `--markdown` | `--md` | Render output through Glamour for ANSI-styled terminal output |
| `--tui` | — | Enable interactive TUI mode (spinner during analysis + scrollable viewport for results) |

All flags are independent and composable. The default output (no flags) is plain text, identical to previous behaviour.

## Examples

```bash
# Plain text output (default — unchanged from previous behaviour)
./virgil-learn /path/to/bash/scripts

# Glamour markdown rendering (no TUI required)
./virgil-learn --markdown /path/to/bash/scripts
./virgil-learn --md /path/to/bash/scripts

# Interactive TUI with spinner and scrollable viewport
./virgil-learn --tui /path/to/bash/scripts

# Interactive TUI with Glamour rendering inside the viewport
./virgil-learn --tui --markdown /path/to/bash/scripts
```

## Output

All modes display the same data:

1. **Results per file** - Each detected pattern, its frequency, and line numbers
2. **Summary** - Total pattern count, breakdown by type, and Phase 2 systems engineering validation

### Phase 2 Systems Engineering Patterns

The tool specifically validates three critical patterns:

- `configuration_center` - Centralized configuration at the top of scripts
- `defensive_prevalidation` - Validation checks before resource use
- `operation_validation` - Exit code checking after operations

If all three show `DETECTED`, the codebase follows the expected systems engineering conventions.

## Output Modes

### Default (plain text)
Standard formatted output written directly to stdout. Safe to pipe to other tools.

```bash
./virgil-learn /path/to/scripts | grep DETECTED
```

### `--markdown` / `--md`
Passes the output through [Glamour](https://github.com/charmbracelet/glamour) for ANSI-styled rendering in the terminal. Style is auto-detected from the terminal background. Degrades gracefully to plain text if Glamour encounters an error.

### `--tui`
Runs a full [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI:
- Spinner displayed while analysis runs
- Results rendered in a scrollable [Viewport](https://github.com/charmbracelet/bubbles)
- Navigate with `j`/`k` or arrow keys
- Quit with `q` or `ctrl+c`

Combine with `--markdown` to render Glamour-styled content inside the viewport.

## Dependencies

| Package | Version | Purpose |
|---|---|---|
| `charm.land/glamour/v2` | v2.0.0 | Markdown rendering |
| `charm.land/bubbletea/v2` | v2.0.2 | TUI framework |
| `charm.land/bubbles/v2` | v2.1.0 | Spinner + Viewport components |
| `charm.land/lipgloss/v2` | v2.0.2 | TUI styling |
