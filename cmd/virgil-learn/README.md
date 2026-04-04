# virgil-learn CLI Tool

A command-line tool to analyze Bash codebases and extract systems engineering patterns.

## Building

```bash
cd cmd/virgil-learn
go build -o virgil-learn
```

Or from project root:
```bash
go build -o bin/virgil-learn ./cmd/virgil-learn
```

## Usage

```bash
./virgil-learn /path/to/bash/scripts
```

## Example

If you have your 9 production Bash scripts in `/home/user/scripts/bash`:

```bash
./virgil-learn /home/user/scripts/bash
```

## Output

The tool will display:
1. **Total patterns detected** - How many patterns were found
2. **Pattern breakdown by type** - Each pattern type and frequency
3. **Line numbers** - Where each pattern occurs (first 5 shown)
4. **Phase 2 validation** - Which critical systems engineering patterns were detected

## What to Look For

The tool checks for three critical Phase 2 patterns:
- `configuration_center` - Centralized configuration at top of scripts
- `defensive_prevalidation` - Validation checks before resource use
- `operation_validation` - Exit code checking after operations

If all three show `✓ DETECTED`, the analyzer is working correctly!

## Important

Make sure you update the import path in `main.go` to match your actual module name in `go.mod`.

Current import:
```go
"github.com/your-org/virgil/pkg/virgil/learning"
```

Replace with your actual module path.
