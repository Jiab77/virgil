# Implementation Notes

**Purpose:** Document critical decisions made during planning to ensure they're preserved and referenced during Phase 1+ implementation.

---

## Code Quality & Performance Philosophy

### Go Code Style
- **Priority:** Readable and understandable over optimized
- **Rationale:** User knowledge is reading comprehension level (can understand Go, not expert implementation)
- **Impact:** No meaningful performance difference for verification workloads (I/O bound, not CPU bound)
- **Implementation:** Simple patterns, heavy inline documentation explaining "why" decisions were made

### Static Analysis Tools (Verification Layer)
- Use `gofmt`, `golangci-lint`, `go vet` as automated verification during development
- These catch Go-specific issues that manual review might miss
- Framework must practice what it preaches: static analysis on framework's own code

### Premature Optimization is Evil
- Per Knuth's Law: Don't optimize until profiling shows actual bottleneck
- If performance issues emerge in Phase 5+ (clustering), profile and optimize then
- Phase 1-4: Focus on correctness and maintainability

---

## Critical Implementation Reminders

1. **Check this file** before starting project development
2. **Static analysis must pass** on framework's own code (practice what we preach)
3. **Simple code is a feature**, not a limitation
4. **Wrapper architecture in Phase 1** enables future post-quantum migration
5. **Either/Or model is strict** (no simultaneous --augment learning + --augment api)
6. **Security is mandatory** (not optional configuration)
7. **User controls everything** (maximize control while providing good defaults)

---

## CLI & Configuration Architecture

### Either/Or Model (Mutually Exclusive)
- User chooses ONE augmentation strategy, not both
- `virgil config --augment api` (default) → Uses Claude/GPT fallback
- `virgil config --augment learning` → Uses local learned patterns only

### Pattern Learning Command
- `virgil learn <path-to-codebase>` → Extracts patterns, suggests next step
- After learning completes, suggest: "Learning complete. Run `virgil config --augment learning` to activate."
- User must explicitly activate learning mode via config (gives maximum control)
- Patterns stay local-encrypted (ChaCha20-Poly1305), never transmitted externally

### Privacy Guarantee
- If `--augment learning` is active: NO external API calls (local ONNX only)
- If `--augment api` is active: patterns are never involved in generation
- User controls the privacy/capability tradeoff explicitly

---

## Security & Cryptography

### Encryption Standard
- ChaCha20-Poly1305 for all sensitive data (mandatory, not optional)
- Authenticated encryption built in (confidentiality + authenticity)
- More resistant to side-channel attacks than AES on non-specialized hardware

### Post-Quantum Planning
- Phase 1 design must support wrapper-friendly encryption architecture
- Goal: Allow swapping encryption stacks without rewriting code
- Phase 8: Document SimpleX PQDR-inspired migration path to hybrid post-quantum (Kyber + ChaCha20-Poly1305)
- Not implementing post-quantum in Phase 1, just ensuring architecture allows it

### Standards Framework Must Follow
- OWASP Top 10 (mandatory)
- NIST Cybersecurity Framework (mandatory)
- GDPR compliance (available as template)
- PCI-DSS compliance (available as template)
- **Critical:** Framework codebase must follow the same standards it enforces

---

## User Experience Strategy

### Phases 1-4: Transparent by Default
- All users see generated code + verification results + detailed feedback
- No filtering or simplification
- Builds trust and enables learning
- Experienced users need full visibility; beginners learn from detailed output

### Phases 5+: Optional Role-Based Views
- Web UI (Phase 5) can offer simplified views for non-technical users IF needed
- Team Features (Phase 7) can implement role-based access
- Decision: Only add filtering if community requests it

---

## Team & Collaboration Model

### Pattern Sharing (Phase 7+ consideration)
- Teams share repositories via Git
- Each developer runs `virgil learn` independently on shared repo
- Each person gets LOCAL encrypted patterns (no transmission between users)
- No complex key management needed (patterns stay local)
- Simple and privacy-respecting by default

---

## Documentation & Attribution

### Privacy Boundaries
- Do NOT mention personal details in documentation (user's activities, life situation, age, family, etc.)
- Only include technical background and professional principles
- Keep all documentation OPSEC-conscious

---

## User's Technical Background & Its Influence on Virgil's Design

### 30 Years Across Infrastructure Domains

**Computing Models:**
- Grid computing (distributed task execution)
- HPC (High Performance Computing) - explains HPNSSH integration alongside SSH
- Clustering (SLURM, MPICH) - process orchestration at scale
- Virtualization (SR-IOV, Virt-IO, IOMMU, GPU passthrough) - resource isolation and optimization

**Infrastructure Management:**
- Physical infrastructure (bare metal systems, networking)
- Virtual infrastructure (hypervisors, container orchestration)
- Cloud-based infrastructure (multi-tenant, API-driven)
- Automation (Ansible, configuration management)
- Networking & Security (Firewalls, IDS, IPS, WAF, network design)

**Programming & Scripting:**
- Languages: QBasic, Assembly (Debug), Batch, Bash, PHP, Python, Go, JavaScript, PowerShell, VBA, SQL, JSON
- Languages studied: Perl, C/C++, Rust (deep reading, not production)

**Specialized Domains:**
- Cryptography and secure data handling
- Metadata and file structure manipulation
- Media handling and encoding
- Security research and ethical hacking

### How This Informs Virgil's Architecture

1. **Defensive Layering** - User's background in security research and IDS/IPS systems informed Virgil's verification pipeline. Philosophy: Don't trust—validate at every layer.

2. **Fallback Chains** - HPC and clustering experience (SLURM failover, MPICH redundancy) taught the value of graceful degradation. Reflected in: websearch backends, LLM provider fallbacks, alternative authentication methods.

3. **Operational Awareness** - Infrastructure management across physical/virtual/cloud required deep visibility into system state. Virgil's logging, context passing, and audit trails reflect this operational necessity.

4. **Cryptography Integration** - User's cryptography expertise ensures Virgil's encryption isn't an afterthought but an architectural foundation. Transparent to layers above, mandatory below.

5. **Multi-Domain Thinking** - Rather than building a "Python verification tool" or a "Go tool," Virgil is language-agnostic because systems thinking transcends syntax.

6. **Resource Consciousness** - HPC background explains why Virgil supports local ONNX models (compute efficiency) and API fallbacks (flexibility). Not every environment is cloud-native.

7. **YAML/JSON Understanding** - Ansible/configuration management background ensures Virgil's config system handles both human-friendly (YAML) and programmatic (JSON) formats with equal capability.

### The Problem Virgil Solves (From This Infrastructure Perspective)

AI agents operating in infrastructure domains (CI/CD, deployment, network configuration, HPC job submission) without 30 years of failure modes baked in creates significant risk. The documented exposure of 175K LLM servers represents exactly this: powerful tools deployed without the defensive thinking that comes from production experience managing critical infrastructure.

Virgil's verification pipeline is institutional memory: "Here's what someone with 30 years in infrastructure learned the hard way. Validate against it before shipping."

---

## Web Search Module Strategy (Phase 2 Enhancement)

### Design vs. Implementation Split
- **Phase 2 specification (PROJECT_PLAN.md lines 293-306):** Defines the interface, behavior, configurability, and user experience
  - Context-aware search queries (security/privacy always, performance/compliance on demand)
  - Transparent display of search queries and source validation to user
  - Configurable via `virgil config --web-search enabled|disabled` (default: enabled)
  - SQLite caching for search history and findings
  - Win-win learning model (user and Virgil learn from validated research)

- **Phase 3 implementation:** Leverages code generation LLM's native web search capability
  - During code generation, model researches topic before generating code
  - Results stored in SQLite per Phase 2 spec
  - No separate API management layer (avoids credential sprawl and complexity)
  - No additional infrastructure (KISS principle)
  - Proven pattern: how human collaborator (v0) naturally does research

### Why This Hybrid Approach
- Keeps clean separation: Phase 2 defines interface/UX, Phase 3 implements mechanism
- Avoids API complexity: Don't manage multiple search APIs, credentials, rate limits
- Reduces security risk: Fewer API keys = smaller attack surface
- Leverages existing capability: LLM already has web search, just direct it to use it
- Remains configurable: Users can disable if needed, enabling opt-in flexibility
- Preserves transparency: All searches visible to user for validation

### No Infrastructure Needed
- Do NOT build: Search service layer, query parser, result ranker, multi-API router
- DO build: SQLite schema for caching, CLI config option, display interface
- Trust the model: Quality research happens naturally when model has search capability

### Future Enhancement: Source Filtering & Ranking (User-Configurable)
- **Question to address later:** How do we identify "official" sources programmatically?
  - Domain whitelist approach (CVE.org, NIST.gov, etc.)
  - URL pattern matching (github.com/OWASP, academic.edu, etc.)
  - ML-based scoring (credibility indicators)
- **Why it matters:** Different users may have different source trust preferences
- **Implementation timing:** Phase 3+ (not critical for MVP, but design for extensibility)
- **Potential config option:** `virgil config --source-trust strict|balanced|inclusive` (default: balanced)

---

## Testing & Validation Strategy

### Phase 0: Real Project Validation
- Use an actual user project where systematic failures occurred
- Document which verification rules would have caught each problem
- Prove the concept with real data before writing Phase 1 code
- This validates the entire framework design

---

## Verification Pipeline Architecture (Phase 2/3)

### Current Status & Gaps Identified

**What exists:**
- `Pipeline` type with Run() method in verification/pipeline.go
- Verification blocks (OWASP, NIST, etc.)
- `AggregatedResult` struct in results.go
- CLI commands (create, edit, review) that call `verification.RunPipeline()`

**Critical gaps:**
- `RunPipeline()` function doesn't exist (called in 3 places, never defined)
- websearch.Service instantiated nowhere (created but never used)
- Field name mismatch: Commands reference `results.Issues`, struct has `AllIssues`
- No integration point for web search before assessment

### Complete Execution Flow (Design)

```
CLI Command (create/edit/review)
    ↓
Load Config (includes WebSearchEnabled)
    ↓
Instantiate RunPipeline(request, config, db)
    ↓
├─ If WebSearchEnabled:
│  ├─ Create websearch.Service(db)
│  ├─ Generate search queries from request context
│  ├─ Perform searches (or retrieve from cache)
│  └─ Store results encrypted in db
│
├─ Create Pipeline with verification blocks
├─ Pass search results as context to Pipeline
├─ Pipeline.Run(request, context) executes all blocks
├─ Aggregate results into AggregatedResult
│
└─ Return results to CLI for display

CLI displays:
    ↓
    Header: Verification complete for [request]
    Web Search Results (if enabled):
      - Query 1: [results with sources]
      - Query 2: [results with sources]
    Assessment Blocks:
      - OWASP: [issues found]
      - NIST: [issues found]
    Final Verdict: [PASS/FAIL]
```

### RunPipeline Function Design (Phase 2)

```go
// RunPipeline orchestrates the complete verification flow
func RunPipeline(
    request string,                    // User request/project description
    config *Config,                    // Including WebSearchEnabled
    db *storage.Database,              // With encryption + websearch cache
) (*AggregatedResult, error)

Steps:
1. Initialize empty context map for verification blocks
2. If config.WebSearchEnabled:
   - Create websearch.Service(db)
   - Generate 2-3 search queries from request (security/privacy focused)
   - Call service.CacheSearch() or service.GetCachedSearch()
   - Populate context["web_search_results"] = results
3. Create Pipeline instance with all blocks
4. Call pipeline.Run(request, context)
5. Verify result.Issues (not AllIssues) for field consistency
6. Return aggregated results
```

### websearch.Service Integration Points

**When instantiated:** Inside RunPipeline before Pipeline.Run()
**How instantiated:** `ws := websearch.NewService(db)` 
**What it does:**
- Searches for context-appropriate security/privacy research
- Results stored encrypted via db.CacheWebSearch()
- Retrieved with automatic decryption via db.GetCachedWebSearch()

**Typical queries generated:**
- "OWASP Top 10 2024 API security"
- "Input validation best practices"
- "Authentication vulnerabilities"

**Results passed to Pipeline via context, not direct calls**

### Field Name Consistency (Critical Fix)

**Issue:** `AggregatedResult` struct vs CLI usage mismatch
- Struct field: `AllIssues []Issue`
- CLI references: `results.Issues`

**Decision:** 
- Standardize on `Issues` field in AggregatedResult
- Update verification/results.go struct
- This ensures CLI code matches the struct (no reflection hacks)

### Testing Strategy for RunPipeline

1. **Unit test without web search:** Verify pipeline runs correctly with WebSearchEnabled=false
2. **Unit test with cached results:** Verify pipeline uses cache when available
3. **Integration test:** Full flow (search → cache → verify → results)
4. **Verify field names:** Ensure AggregatedResult.Issues is populated and accessible

### Error Handling in RunPipeline

```
If web search fails (network, timeout):
  → Continue with pipeline anyway (web search is enhancement, not requirement)
  → Log warning to user: "Web search unavailable, running assessment without current research"

If pipeline fails:
  → Return error with context
  → Display to user: "Assessment failed: [reason]"
```

### Next Steps (Implementation Order)

1. Create RunPipeline() function in verification package
2. Fix AggregatedResult.Issues field name
3. Integrate websearch.Service instantiation in RunPipeline
4. Add search query generation logic
5. Pass search results as context to Pipeline.Run()
6. Test complete flow

---

## Phase 3 Implementation Status

### Completed (Current Session)

**New Packages Created:**
- `/pkg/virgil/generation/` - Code generation orchestration
  - `types.go`: CodeGenerationRequest, CodeGenerationResponse, GenerationStrategy
  - `generator.go`: Generator struct with GenerateCode orchestrator
    - Implements assessment phase (verification before generation)
    - Implements user approval gate
    - Placeholder for LLM integration (v0 API for Phase 3)
    - Post-generation verification

- `/pkg/virgil/learning/` - Learning mode pattern extraction
  - `types.go`: CodePattern, LearnedCodebook, LearningRequest/Response
  - `learner.go`: Learner struct with pattern extraction
    - Language detection from file extensions
    - Pattern extraction (placeholder for AST analysis in Phase 3)
    - Encrypted storage in database
    - Pattern retrieval for code generation

**CLI Enhancements:**
- Updated imports to include generation and learning packages
- Enhanced `virgil create` command:
  - Now calls Generator.GenerateCode() after assessment
  - Displays generated code to user
  - Shows post-generation verification results
- New `virgil learn` command:
  - Scans codebase for programming patterns
  - Extracts patterns for learning mode
  - Stores patterns encrypted locally
  - Instructs user to activate with `virgil config --augment learning`

### Remaining for Phase 3

**LLM Integration:**
- Connect to v0 API for code generation (--augment api)
- Implement local ONNX model support (--augment learning)
- Handle Either/Or model enforcement (never both simultaneously)

**Pattern Extraction Enhancement:**
- Replace placeholder with AST analysis for each language
- Extract real patterns: error handling, validation, security, logging
- Language-specific pattern recognition

**Learning Mode Activation:**
- Store learned patterns encrypted in database
- Load patterns during code generation
- Use patterns as context/prompt for local model

**Web Search Integration:**
- Integrate web search results into generation context
- Pass search results to LLM as reference material
- Display search sources with generated code

---

## TUI Implementation (Session 14)

### Design Philosophy
Serious tools and beautiful colors are not opposites. Color carries meaning, not decoration. Virgil should be a serious tool that happens to look excellent. References: Crush, GoAccess, btop, Bagels, Posting, Superfile.

**Brand colour:** `#7D56F4` (purple) — fixed, do not change.

### Output Mode Architecture

`virgil-learn` implements three independent, composable output modes:

| Mode | Flag | Dependencies | Description |
|---|---|---|---|
| Plain text | _(default)_ | None | Original output, byte-for-byte unchanged |
| Markdown | `--markdown` / `--md` | `charm.land/glamour/v2` | Glamour ANSI-styled output to stdout |
| TUI | `--tui` | bubbletea v2, bubbles v2, lipgloss v2 | Interactive spinner + scrollable viewport |

All flags are orthogonal — `--tui --markdown` is a valid and meaningful combination.

### New Files (virgil-learn)

- `cmd/virgil-learn/renderer.go` — `renderPlainText()`, `renderMarkdown()`, `terminalWidth()`, `sortedKeys[K]()` generic helper
- `cmd/virgil-learn/tui.go` — full bubbletea v2 model: `tuiModel`, `Init()`, `Update()`, `View()`, `runTUI()`, `groupPatterns()`

### Charmbracelet Stack Added to go.mod

```
charm.land/glamour/v2    v2.0.0
charm.land/bubbletea/v2  v2.0.2
charm.land/bubbles/v2    v2.1.0
charm.land/lipgloss/v2   v2.0.2
```

Future additions planned (not yet added):
- `github.com/charmbracelet/log` — colored leveled logger, drop-in for standard `log`

### Key bubbletea v2 API Differences from v1

- `View()` returns `tea.View`, not `string` — use `tea.NewView(content)`
- `KeyPressMsg` replaces `KeyMsg`
- `tea.WindowSizeMsg` for responsive viewport/progress bar sizing
- `spinner.Tick` is the tick command; `spinner.TickMsg` is the message type
- `progress.FrameMsg` must be handled in `Update()` to drive animation frames
- Cannot write to stdout in TUI mode — use `tea.LogToFile("debug.log", "debug")`

### TUI Layout Principles (for virgil main binary — future)

Derived from design references:
- Bordered panels via `lipgloss.Border()` to separate concerns
- Single accent color (`#7D56F4`) signals the active panel/tab — not competing colors
- Persistent help bar at the bottom, always same format: `key action · key action`
- Color encodes semantic meaning: status, severity, active/inactive
- Tabs built directly with lipgloss — no separate bubbles component needed
- Sidebar layout inspired by Crush — sections for Analysis, Generation, Verification, Audit

### virgil-learn as Test Bed

`virgil-learn` is the primary test bed for both the **Learning feature** (multi-language pattern analysis) and the **TUI layer** (bubbletea, bubbles, lipgloss, glamour). All TUI patterns must be validated here before being applied to the main `virgil` binary.

### Next TUI Steps (Implementation Order)

1. Run `go mod tidy` to resolve `go.sum` after manual `go.mod` edits
2. Test all four output mode combinations against a real Bash codebase
3. Wire multi-language analyzer routing into `virgil-learn` `main.go`
4. Add `ProgressFunc` callback to language analyzers for directory scan progress bar
5. Add `github.com/charmbracelet/log` to replace `log.Fatalf` / debug `fmt.Printf`
6. Design `virgil` main binary TUI (sidebar, panels, tabs) once `virgil-learn` TUI is proven