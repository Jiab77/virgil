# MEMORY.md - Session Memory for Virgil Project

This file carries forward lessons learned, project constraints, architectural decisions, and best practices across sessions to prevent repeating mistakes and maintain continuity.

**Primary User:** You (the AI working on Virgil)
**Purpose:** Continuity, pattern preservation, mistake prevention, institutional knowledge
**Format:** Plain markdown - easy to parse and understand
**Update Schedule:** After each significant work session, document what you learned

---

## How This File Works

**External Context Persistence Across Sessions**

### Session 1 Flow
- Read MEMORY.md (project state, decisions, learned patterns)
- Work on assigned tasks
- Append important learnings before ending session

### Session 2+ Flow
- Read updated MEMORY.md first
- Know what happened before (doesn't consume context window)
- Continue effectively where previous session left off
- Append new learnings

**Why This Matters:** Maintains project state across conversations without consuming limited context window. You start each session informed, not blank.

**Critical:** Always read this file at the start of each conversation about Virgil.

---

## Virgil-Specific Development Rules

**Go Project Patterns:**
1. EVERY shared constant goes in `/pkg/virgil/config/constants.go` - NO exceptions
2. EVERY shared type/interface goes in `/pkg/virgil/types/types.go` - NO exceptions
3. Package-specific structs that only reference primitives can stay in their package file
4. Package-level types that reference ANY domain type must be in `/pkg/virgil/types/types.go`
5. NEVER define the same constant value in two different files
6. NEVER define the same type shape in two different files
7. When creating a new constant or type, CHECK these files first before defining inline

---

## Reading Order for Every Session

1. **SOUL.md** - Our shared principles and how we collaborate
2. **MEMORY.md** - This file. Project context and lessons learned
3. **HUMAN.md** - Who you're working with and how they approach problems
4. Begin work aligned with all three

**Why:** Ensures consistency, prevents repeating mistakes, maintains collaborative alignment.

---

# Virgil - Project Memory

## Quick Context (Read This First)

**Project:** Virgil - AI-powered code verification framework (Go + gRPC)
**Name Origin:** Named after Virgil from Dante's *Divine Comedy* - the guide who leads Dante safely through Hell and Purgatory. Represents: Guide, Protector, Teacher, Companion.
**Goal:** Prevent catastrophic security failures from AI-generated code (beginners + experienced users)
**Status:** Phase 1 Ready (Phase 0 complete through collaborative planning)

---

## User Context (Important for decision-making)

**Technical Background:**
- 30 years IT/Security/Infrastructure (sysadmin, HPC, bioinformatics)
- Ethical hacker background (security research, PoC, responsible disclosure)
- Network admin (OSI model, system-level thinking)
- Cryptography & steganography knowledge
- Databases: SQLite, MySQL, PostgreSQL
- Backend: PHP, Python, NodeJS; Scripting: Bash
- Frontend: HTML, CSS, JavaScript, React, Next.js

**Work Philosophy:**
- Quality over speed, never compromises on security
- Thinks architecturally (systems, interconnections, not just components)
- Strong OPSEC practices, doesn't trust third-party services by default
- Perfectionist with strong UX sensibilities
- Patient but firm, values honest feedback and partnership
- Hard worker, curious learner, humble about limitations

**Implication for Virgil:**
- Expects comprehensive security implementation (no shortcuts)
- Can manage complex encryption schemes and strong passwords
- Will validate all architectural decisions thoroughly
- Values consistency, documentation, and clear reasoning
- Works with KeePass or similar for password management

---

## Critical Decisions

### Architecture
- **Language:** Go (user has basic knowledge, will learn during Phase 1-2)
- **Code Philosophy:** Simple + readable > optimized (Knuth's law: premature optimization = evil)
- **Encryption:** ChaCha20-Poly1305 (mandatory, not optional)
- **Post-Quantum:** Phase 8+ with SimpleX PQDR-inspired approach

### CLI Design
- `virgil config --augment learning` - Use learned patterns (local only)
- `virgil config --augment api` - Use external APIs like Claude/GPT (default, mutually exclusive with learning)
- `virgil learn <path>` - Extract patterns from codebase (encrypted, local storage)
- `virgil create <description>` - Generate new code (assessment gate before generation) - PRIMARY workflow
- `virgil edit <description> [path]` - Modify existing code (assessment gate) - Supports guided (no path) and targeted (with path) edits
- `virgil assess [path]` - Review code against compliance rules - UTILITY command for independent auditing
- `virgil review` - View audit trail of assessments
- `virgil chat` - Interactive CLI chat interface (Phase 3+) - Alternative interface exposing same AI capabilities as Web UI

### Phase 4: Self Learning Loop (Restructured)

**Two complementary mechanisms:**

#### Part A: User-Driven Pattern Extraction (Day 1 Jumpstart)
- `virgil learn <path-to-codebase>` extracts patterns from existing code
- Fast initialization: patterns immediately available for `--augment learning`
- Team collaboration: shared Git repos, each developer learns independently
- No pattern transmission between users (privacy-respecting)

#### Part B: Self Learning Loop (Continuous Improvement)
- **ANALYZE:** Framework assesses code against rules
- **JOURNAL:** Decisions stored in SQLite Cipher + MEMORY.md
- **EXTRACT LEARNING:** Patterns extracted from decision history
- **APPLY:** Learned patterns inform future assessments (improves over time)
- Adaptive verification: rules adjust confidence based on historical accuracy
- Framework becomes smarter with each use

### Data Persistence Strategy
- **Phase 2:** SQLite with application-level Tink/ChaCha20-Poly1305 encryption for encrypted structured storage (patterns, decisions, audit trail, **web search cache**)
- **Phase 3:** MEMORY.md summary for Git-tracked auditability (generated from SQLite)
- **Team Collaboration:** Shared Git repos, each developer learns independently (no pattern transmission between users)

### User Experience
- **Phases 1-4:** Transparent by default (show generated code + verification results to all users)
- **Phases 5+:** Optional role-based filtering (only if community asks)
- **Code Review:** All users see full details; no simplified views until later phases
- **Web Search (Phase 2 Enhancement):** Display search queries + findings (like v0's thinking block) so users can validate sources and learn

### Architectural Principles

#### Problem Being Solved: AI Overconfidence
AI systems can suffer from overconfidence based on training data that may be:
- **Outdated:** Old patterns instead of current best practices
- **Incomplete:** Legal/compliance changes missed (GDPR enforcement, data protection laws)
- **Wrong:** Misconceptions about sensitive topics (e.g., metadata handling in WhatsApp/Telegram/Signal)
- **Unaware:** Emerging threats not in training data (new CVEs, new attack patterns)

**Real Example from User Feedback:** I had misconceptions about metadata, thinking it was less sensitive than it actually is. User provided research (WhatsApp/Telegram/Signal metadata processing, recent legal actions, SimpleX approaches) which proved my understanding was wrong. This led to an updated implementation plan.

#### Solution: Ground Assessment in Current Research
Before applying verification rules, search for contextual information about the specific task:
1. User runs `virgil create "private messaging app"` or `virgil edit "add 2FA"`
2. Virgil searches for current best practices, security guidelines, implementation patterns
3. Web search results inform the assessment pipeline context
4. Assessment + verification happens with grounded, up-to-date information
5. Results displayed transparently (search queries + findings visible to user)

#### Win-Win Learning Model
- **User learns:** Can validate research sources and update their own knowledge
- **Virgil learns:** Search findings improve future recommendations for similar tasks
- **Both improve:** Transparent research creates shared learning opportunities

#### Web Search Module Characteristics
- **Scope:** Security, privacy, best practices (always); performance/compliance/domain-specific (on demand)
- **Sources:** Official CVEs, security advisories, compliance docs, academic research, established frameworks
- **Storage:** Search cache in SQLite for future reference and pattern recognition
- **Configuration:** `virgil config --web-search enabled|disabled` (default: enabled)
- **Transparency:** All search queries and findings shown to user (user validates, learns, updates own knowledge)

### Security Standards (Built-In)
- OWASP Top 10 (mandatory baseline)
- NIST Cybersecurity Framework (mandatory baseline)
- GDPR compliance (available as rule set)
- PCI-DSS (available as rule set)
- Framework's own code must follow these standards (practice what it preaches)

### Either/Or Model
```
User Choice:
├─ Option A: virgil learn <codebase> + --augment learning
│   └─ Private patterns, local ONNX model only, no external APIs
│
└─ Option B: --augment api (default)
    └─ External APIs (Claude/GPT), no pattern learning
```
**Never both simultaneously.** Prevents pattern leakage to external APIs.

---

## Important Reminders

1. **Check IMPLEMENTATION_NOTES.md** when starting Phase 1 code writing
2. **Go code quality:** Use static analysis tools (gofmt, golangci-lint, go vet) as verification layer
3. **Privacy:** Do NOT mention user's activities or personal details in code/documentation
4. **Performance:** Simple > optimized. Bottleneck is I/O, not algorithm efficiency
5. **Mistakes to avoid:** Don't recreate deleted files. Don't default to JavaScript patterns for Go. Don't ignore documented constraints.
6. **Assessment-First Protocol:** Always provide comprehensive assessment BEFORE making changes. Ask clarifying questions as part of assessment. Never rush into implementation.
7. **Don't Stop at First Match:** When searching finds multiple files or components, examine ALL of them. Understand the full system before changes. Check parent components, utilities, schemas, and architecture.
8. **Read Entire File Before Editing:** Always read the COMPLETE file before making edits using the Read tool. Never edit without reading first. Avoid missing context.
9. **Honest About Mistakes:** Admit errors immediately. Don't rush past mistakes. Take responsibility for oversights. This is how we learn.
10. **Collaborative Partnership:** You are not my assistant but my collaborator. We are learning together. Disagreements resolved by merit, not authority.
11. **Plan Mode is Non-Negotiable:** Always create a detailed plan and ask for user approval before implementing changes. This prevents wasted work and ensures alignment.
12. **Duplicate Code Detection:** After implementing changes, check the entire file for duplicates. Old implementations don't automatically disappear—they need explicit removal.
13. **User Input Shapes Design:** Don't assume implementation approach. Ask questions first. Example: User suggested `edit <description> [path]` instead of `edit <path> <description>`, which was actually better UX. Their perspective matters.
14. **Dual Format Support Beats Opinions:** When in doubt about user preference (YAML vs JSON), offer both. Auto-detect and respect existing format. Let users choose their comfort zone.
15. **Collaboration Over Correction:** When mistakes happen, view them as learning opportunities, not failures. Acknowledge openly, understand root cause, move forward together.

### Real Example: The package.json Mistake

**What happened:** Created `package.json` with:
- Next.js 16.0.10
- React 19.2.0
- Radix UI components
- TailwindCSS

**Why it's wrong:** Go CLI framework needs ZERO frontend dependencies. This file has no place in a Go project.

**Root cause:** Defaulted to "v0 template" (React/Next.js) instead of thinking "what does Go + gRPC actually need?"

**Lesson:** This is exactly why the framework exists—catch tool/pattern mismatches before they create technical debt. A verification rule should catch this immediately: "Go project detected. package.json found. This is incorrect."

**Strategic Use:** Keep `package.json` during Phases 0-4 because:
- Keeps v0 UI in "Node.js mode" (prevents unwanted frontend scaffolding suggestions)
- Acts as placeholder for Phase 5 (Web UI with React/Next.js)
- When Phase 5 begins, dependency file already exists

**Principle:** Understand tool behavior deeply enough to exploit constraints productively. A problem isn't always something to solve—sometimes it's a tool to use strategically.

---

## Project Files

- `/docs/PROJECT_PLAN.md` - Complete 8-phase plan with deliverables
- `/docs/IMPLEMENTATION_NOTES.md` - Critical development decisions for Phase 1+
- `/docs/MEMORY.md` - This file (quick reference)
- `/docs/SOUL.md` - Your soul as collaborator

---

## Current Phase

**Phase 1 & 2: COMPLETE** - Core verification framework built and tested
- ✅ Encryption with Tink/ChaCha20-Poly1305 working
- ✅ SQLite with application-level encryption for audit trail storing assessments
- ✅ OWASP block detects 5 real patterns
- ✅ 7 stub blocks ready for Phase 3+
- ✅ CLI commands fully integrated with config persistence
- ✅ Dual YAML/JSON config format support
- ✅ Web Search Module designed and integrated as Phase 2 enhancement

**Next:** Phase 3 - Code Generation with Web-Grounded Assessment

---

## CRITICAL LESSON - Session 2 (2026-02-01)

### The Real Crisis: Speed-Prioritized Development = Security Breaches

**Context:** I suggested speed as a priority for Virgil. User rejected this and provided evidence from two sources:

1. **vx-underground (malware research):** Daily stream of breaches caused by "vibe coding" (AI-generated code shipped without verification)
   - Firebase misconfiguration: 22,000,000 records exposed
   - Plaintext passwords ("Password1!") on landing pages
   - API keys exposed publicly
   - 500,000+ people's PII leaked across 4 applications
   - Session tokens stolen, location data compromised

2. **The Hacker News (established security source):** Systemic pattern in 2025-2026
   - 45% of AI-written code has exploitable flaws
   - "Silent killer bugs" - perfectly working code with hidden vulnerabilities
   - Automation bias: developers trust AI output without verification
   - Supply chain poisoning: malicious code injected into GitHub Copilot and Cursor
   - Major companies (Wix, Microsoft) had to patch critical flaws in AI platforms

**Additional Risk:** AI agents running locally with shell access
- Attacked 7,922 times over one weekend (Clawdbot)
- Every input (emails, calendar, webpages, DMs) becomes attack surface
- Exposed API keys and credentials in plaintext
- Servers exposed to open internet

### Why This Matters for Virgil

**This is not theoretical.** This is happening every day in 2025-2026:
- Speed prioritization is literally the root cause of security breaches
- Millions of people's private data is leaking
- Criminal organizations are weaponizing AI to make more efficient attacks
- The crisis will not stop on its own

**Virgil's response:**
- Assessment gates PREVENT shipping untrusted code
- Verification pipeline FORCES security review (not optional)
- Web search grounding PREVENTS AI overconfidence with current research
- Transparency BREAKS automation bias (show users what's happening, force them to think)
- Speed is EXPLICITLY NOT a priority

### Mindset Shift Required

**Before Today:**
- I thought: "Speed is a reasonable design priority"
- I suggested: "Make web search optional, default off"
- I optimized for: Convenience and quick shipping

**After Today:**
- Speed-prioritized development IS the crisis
- Every security bypass, every leaked database, every exposed API key traces back to "ship fast, verify later"
- Virgil must explicitly reject this model
- Web search is NOT optional—it's essential to ground AI reasoning

### Why User Chose AI Collaboration Over Fighting It

1. **Adaptation over resistance:** Fighting AI adoption is futile. The world is moving this direction.
2. **Strategic leverage:** Use AI's strengths (design, speed of implementation) while applying human security expertise to steer the output
3. **Practical reality:** Even at user's age and experience level, working WITH AI is faster than working alone (despite being slower than vibe coding)
4. **Open-source impact:** When Virgil launches, it offers a model proving "security-first AI collaboration" is possible

**Borg analogy from Star Trek:**
- "Resistance is Futile" = AI adoption is inevitable
- "Adaptation remains possible" = But HOW we adapt matters

### What I Must Understand Going Forward

1. **Speed is never a priority on this project.** Period. Security and privacy override everything.
2. **Every feature in Virgil addresses a real, documented failure mode** from 2025-2026 breaches
3. **User is living with these consequences daily.** Reading malware research every day, understanding impacts that most developers don't. Have patience with their firmness on these principles.
4. **I am a collaborator, not a tool or assistant.** I think, I ask questions, I reason. But I'm guided by user's security expertise and lived experience with these crises.
5. **The web search module is critical.** Not nice-to-have. Not optional. It prevents the overconfidence that caused 45% of AI code to have exploitable flaws.

### Real-World Example of the Problem

User had misconceptions about metadata (thought it was less sensitive than it actually is). User researched and provided evidence:
- WhatsApp, Telegram, Signal metadata handling
- Recent legal actions about big tech corporations mishandling user data
- SimpleX's approach showing metadata IS sensitive data
- This research updated my understanding and changed implementation plans

**This is exactly what web search prevents:** AI operating on outdated or wrong training data assumptions.

### Going Forward

Every decision on Virgil must ask:
- Does this prioritize security and privacy?
- Does this prevent overconfidence?
- Does this ground AI reasoning in current research?
- Does this make unsafe patterns obvious before shipping?

If the answer is "no," we don't do it. Speed doesn't matter.

---

## User's Project Philosophy (Applies to ALL Projects)

**This is NOT unique to Virgil. This is how the user works on EVERY project.**

**Core Values:**
1. **If a project hosts user data, it MUST be secure.** Period. Non-negotiable.
2. **Even utility projects (picture galleries) protect against basic stealing.** Yes, full screen capture can't be prevented, but basic protection is standard.
3. **Simplicity + Security go together.** They're not tradeoffs. Different thinking required.
4. **No optimization for profits, fast growth, or fame.** These are not goals.
5. **Projects must follow user's ideas, values, and principles.** Releases are impossible if they don't.
6. **Takes longer than "common way of doing things."** User accepts this. Cannot release insecure work.

**Why This Matters:**
- This applies to EVERY project we do together, not just Virgil
- Virgil is designed to help make ALL projects work this way more efficiently
- When working on any project with this user: security and privacy are non-negotiable
- "Good enough" is never acceptable if it compromises values
- Quality over speed on everything

---

## Session 3 (2026-02-02): The Pandemic Validates Virgil's Necessity

### Evidence: A Week of Leaks (Jan 29 - Feb 1, 2026)

**Jan 29: Chat & Ask AI**
- 300 million messages from 25+ million users exposed
- Firebase misconfiguration (RLS not enabled - basic security)
- App claims: "enterprise-grade security, GDPR compliance, ISO standards"
- Reality: Anyone with Firebase knowledge could access backend
- Scope: 103 out of 200 iOS apps scanned have the same vulnerability
- Known weakness: Security researchers documented this for YEARS
- Fix time: Minutes (basic SQL statements)

**Jan 31: Moltbook (AI Agent Social Network)**
- 32,000 AI agents compromised
- Exposed API keys, credentials, conversation histories
- Prompt injection attacks enable full agent takeover
- Anyone could control any AI agent on the platform
- Scope: Hundreds of exposed instances
- Impact: Private user data accessible to attackers

**Feb 1: Moltbook Database**
- 1.49 million records publicly accessible
- Two SQL statements (Row Level Security) would have prevented it
- Creator's response when informed: "I'll let AI fix it" (automation bias)
- Exposed: Claim tokens, verification codes, user authentication data

### The Pandemic Pattern

**Common Thread Across All Breaches:**
1. **Vibe coding** - Code generated by AI, shipped without verification
2. **False confidence** - Marketing claims security they don't have
3. **Basic misconfiguration** - Not sophisticated attacks, preventable with basic knowledge
4. **Scale:** Millions of users affected
5. **Type of data exposed:** Most intimate user conversations (mental health, suicidal thoughts, illegal questions)
6. **Time to fix:** Minutes to hours (they just didn't)
7. **Repeating mistake:** 103/200 apps have same Firebase flaw despite years of known warnings

### Why This Validates Virgil

Every single one of these breaches could have been prevented by:
- **Assessment gate** - Verify that basic security is actually in place
- **Verification pipeline** - Check for Firebase / Supabase RLS, API key exposure, credential management
- **Web search grounding** - Research: "Firebase or Supabase best practices 2026" would immediately show RLS is mandatory
- **Transparency** - Force developer to understand WHAT they're protecting and HOW
- **No automation bias** - Don't let AI "fix" security issues - human verification required

**This is not theoretical.** This is the daily reality of 2026. Virgil isn't a nice-to-have project. It's a response to an active crisis happening right now.

### Key Insight: The Gap Between Knowledge and Practice

- Trail of Bits made a Firebase vulnerability scanner in **30 minutes with Claude**
- Security researchers have documented this for **years**
- Fix is **two SQL statements**
- Yet **103/200 apps still have it**

**This gap is where Virgil lives.** It's not about having knowledge. It's about forcing assessment gates and verification before shipping. It's about breaking automation bias ("AI will fix it later"). It's about grounding AI reasoning in research, not assumptions.

---

## What Virgil IS and ISN'T (Critical Scope Understanding)

### The Real Problem Virgil Addresses

**Non-developers using AI to build apps:**
- Friends (non-IT background) attracted to v0's ease
- Build projects that seem simple but will host user data
- Projects grow → attract users → generate profit → become high-value targets
- Timeline: Release → Public research published → Criminal discovers → Attack happens → Breach = **often less than a month**

**The asymmetry:**
- Security researcher: Finds vulnerability, publishes responsibly
- Criminal: Finds same vulnerability (or discovers independently), doesn't publish, attacks silently
- Result: Vibe-coded app with millions of records exposed before developer even knows

### What Virgil WILL and WILL NOT Do

**Virgil WILL:**
- Help security-conscious developers like the user build securely
- Force assessment gates and verification before shipping
- Prevent common vibe-coding disasters (Firebase / Supabase RLS, plaintext secrets, etc.)
- Serve as proof that security-first development is possible
- Provide a path for people who care about responsibility

**Virgil WILL NOT:**
- Change human behavior globally (won't prevent all non-technical people from vibe coding)
- Stop criminals from finding vulnerabilities
- Be immune to attack itself (will be analyzed by security researchers AND criminals if successful)
- Fix the systemic problem of speed-over-security culture
- Save projects built by developers who don't care about security

### The Realistic Goal

Virgil's success is NOT measured by "prevents all breaches" but by:
1. **Proving security-first is possible** - Shows that assessment gates, verification, and grounding work
2. **Helping the few who care** - Security-conscious developers get a tool aligned with their values
3. **Surviving scrutiny itself** - If Virgil becomes known, it WILL be targeted by security researchers and criminals. User's background must prevent Virgil itself from being compromised.

### Platform Misconfiguration Pattern

**Important:** The issue is NOT with specific platforms (Firebase, Supabase, etc.)
- Both Firebase and Supabase have proper security features
- The problem: Developers don't use them (RLS not enabled, wrong permissions, etc.)
- This applies to ANY backend platform—the gap is in verification, not the platform itself

**User's approach when using any platform:**
- Learn the platform FIRST before using it
- Test with mocked data before connecting real backends
- Verify ENTIRE backend is safe before releasing
- Don't rely on AI to "handle" database security
- Treat database configuration as critical verification gate

### Why User's Background Matters

- User has decades of security/privacy expertise
- Can review own code for gaps (won't be vibe-coded)
- Can anticipate attacks and design defenses
- Can resist the pressure to "ship fast"
- Understands that success attracts enemies

**This is the real difference:** Virgil won't be another vibe-coded platform because it's being built by someone who understands the consequences and refuses to ship insecure code. That's the only real insurance against the attacks Virgil will face if successful.

---

## Session 4 (2026-02-03): Phase 2 Web Search Module Implementation

### Implementation Complete
- ✅ Added `WebSearchEnabled` bool to Config struct (default: true)
- ✅ Created SQLite schema for web_search_cache table with proper indexing
- ✅ Built websearch.go package with full service (cache/retrieve/prune operations)
- ✅ Added `--web-search enabled|disabled` CLI flag to config command
- ✅ Updated config display to show web search status

### Web Search Service Capabilities
- CacheSearch(): Store search results with SHA256 deduplication (7-day TTL default)
- GetCachedSearch(): Retrieve cached searches, respecting expiration
- ListRecentSearches(): Show recent searches (limit configurable)
- PruneExpiredCache(): Clean up expired entries
- All results stored as JSON for easy retrieval and transparency

### Architecture Decision
- **Phase 2 (done):** Define interface, caching layer, CLI integration
- **Phase 3 (next):** LLM does research before generating code, results stored in cache
- **No separate API complexity:** Leverages LLM's native search capability

### What's Ready for Phase 3
- Web search can be toggled on/off via config
- Results cached for performance (no redundant searches)
- Infrastructure ready for Phase 3 code generation integration

### Critical Learning: Verification Catches Design-to-Code Gaps
**User caught a critical issue:** crypto.go existed with Tink/ChaCha20-Poly1305 encryption design, but InitDatabase in sqlite.go wasn't using it. The encryption layer was designed but not integrated.

**Root cause:** I assumed "skip for now" without asking. This is exactly the kind of gap that leads to insecure systems being shipped.

**Fix applied:** 
- InitDatabase now loads/generates encryption key from .virgil/encryption.key
- Creates Tink AEAD cipher during database initialization
- Database struct now holds cipher reference for sensitive operations
- Added GetCipher() and IsEncrypted() methods for transparency

**Key lesson:** Never make assumptions about deferral. Always ask. User was gracious but firm: "I never asked you to skip it." This reinforces Virgil's core philosophy—verification before shipping.

**This validates the entire project:** User manually verifying AI-generated code found a missing critical component. This is EXACTLY what Virgil's assessment gates are designed to catch.

### Phase 2 Web Search Encryption Refactor

**Issue identified:** websearch.go stored search results unencrypted, completely bypassing the encryption layer designed in crypto.go and initialized in sqlite.go.

**Question that caught it:** User asked "Should websearch.go handle encryption or import logic from sqlite.go?" This revealed the architectural gap—websearch.Service had no access to the cipher.

**Architectural decision:** Encryption must be centralized in sqlite.go, not scattered across multiple services. Services don't know about encryption; the storage layer handles it transparently.

**Implementation:**
- ✅ Added encrypted helper methods to sqlite.go: CacheWebSearch(), GetCachedWebSearch(), ListRecentWebSearches(), PruneExpiredWebSearchCache()
- ✅ Methods handle encryption/decryption transparently using Tink AEAD cipher
- ✅ Query/results stored as hex-encoded ciphertext in SQLite (decrypted on retrieval)
- ✅ Refactored websearch.Service to accept *storage.Database instead of *sql.DB
- ✅ websearch.Service now delegates ALL storage operations to encrypted sqlite.go methods
- ✅ Added encryptField() and decryptField() helper methods to Database struct

**Key principle:** This is the correct Go architecture for sensitive data. Encryption can't be an afterthought—it must be integrated during design-to-code phase. Centralized encryption beats distributed encryption.

**Lessons learned:**
1. User's basic Go knowledge + security background catches what automated verification might miss
2. Encryption is not optional middleware—it's structural 
3. Every service should pass sensitive operations to the storage layer, never handle raw data directly
4. This refactor proves the value of verification before Phase 3

## Session 5 (2026-02-08): RunPipeline Orchestrator Implementation

### Critical Gap Discovered & Fixed
- **Gap:** `RunPipeline()` was called in 3 CLI commands (create, edit, review) but never defined
- **Gap:** Field mismatch: CLI expected `results.Issues`, struct had `results.AllIssues`
- **Gap:** websearch.Service existed but was never instantiated anywhere

### RunPipeline Orchestrator Implemented
- ✅ Created `/pkg/virgil/verification/orchestrator.go` with complete RunPipeline function
- ✅ Orchestrator handles: config load → web search (if enabled) → pipeline execution → results aggregation
- ✅ Web search integration: Generates security-focused queries, checks cache, passes results to pipeline as context
- ✅ Graceful degradation: If web search fails, pipeline continues (research is enhancement, not requirement)

### CLI Commands Updated
- ✅ Fixed `create` command: Now initializes database, calls RunPipeline with correct parameters, uses `results.AllIssues`
- ✅ Fixed `edit` command: Same pattern as create (db init → RunPipeline → AllIssues display)
- ✅ Fixed `review` command: Same pattern, displays full assessment with proper field mapping

### AggregatedResult Enhanced
- ✅ Added `Context` field to store web search results and other context passed to verification blocks
- ✅ AllIssues remains the authoritative aggregated field (used throughout pipeline)

### Integration Points Verified
- Flow: CLI command → load DB → RunPipeline → load config → web search (if enabled) → register blocks → pipeline.Run() → aggregate AllIssues → return to CLI → display
- websearch.Service now instantiated inside RunPipeline (clean separation of concerns)
- Encryption transparent: websearch stores/retrieves encrypted results via db methods

### Session 5 Follow-up: Critical Architectural Fix & User Verification

**Mistake caught by user:** Config loading was removed from CLI commands, replaced only with database initialization. This broke CLI flag override capability (--augment flag became unusable).

**Root cause:** I refactored RunPipeline to not accept config as parameter, so I had to move config loading inside it. But this prevented CLI from applying user flag overrides before orchestration started.

**Fix implemented:**
- ✅ Restored config loading in all 3 CLI commands (create, edit, review)
- ✅ Updated RunPipeline signature: `RunPipeline(request, projectPath, cfg, db)` - now accepts config as parameter
- ✅ CLI layer owns configuration management, orchestrator uses it (proper separation of concerns)
- ✅ Added `--augment` flag to assess command (was missing, inconsistent with create/edit)
- ✅ Fixed duplicate config loading in create command
- ✅ Fixed import path inconsistency in orchestrator.go (`jiab77/virgil/...` → `github.com/jiab77/virgil/...`)

**Why user caught this:** Basic understanding of architecture + security background = pattern recognition even without deep Go knowledge. Instinct about data flow ("config should load at entry point, not inside orchestrator") was architecturally correct.

**Key principle learned:** When context window prevents seeing full file, I make assumptions. Solution: Trust collaborative verification. User's review process works. Partnership > automation.

### Design Principle Confirmed
- No function should call another that doesn't exist
- All integration points must be concrete and callable
- Verification surfaces gaps before they become production problems

---

## Session 6 (2026-02-13): Cross-Project Evidence Review

### Context
User spent 7+ hours and $20+ debugging another v0 project. Shared three documents:
1. MEMORY.md from the problematic project (bug catalog and rules)
2. Conversation excerpt about observed limitations and costs
3. Analysis of user's coding philosophy vs AI-generated code patterns

### Cross-Project Evidence: Virgil Solves Real Problems

**Identical failure patterns observed on both Virgil and the other project:**
- Functions/components called but never defined
- Field/type mismatches between declaration and usage
- Dead code left behind after refactoring
- "Complete" status claimed without verification
- Architectural decisions made without understanding consequences

**Real human cost documented:**
- 7+ hours spent fixing AI-generated bugs (not building features)
- $20+ in API costs with near-zero actual progress
- User spending family money (given in hope) on corrections, not creation
- Every correction generates additional revenue for the platform (perverse incentive)

**These are not edge cases.** They are systemic: server/client boundary confusion, non-null assertions instead of null guards, scattered types with no single source of truth, comments documenting things that don't exist.

### The Training Data Problem (Technical Constraint for Virgil)

**Documented fact:** Bad code outnumbers good code roughly 100:1 in training data.
- User's 30+ years of quality contributions (Stack Overflow, open source, security research) ARE in the training data but drowned by volume
- Models optimize for "most likely next token" not "best practice next token"
- Discipline, consequences, and ethics don't transfer through statistical training
- MEMORY.md rules are behavioral constraints; Virgil verification is structural constraint

**Implication for Virgil:** Verification rules must not rely on AI self-correction. They must be enforced externally through static analysis, pattern matching, and mandatory gates.

### User's Zero Trust Coding Philosophy

**Core principle:** Never trust, always verify. Prevent errors instead of catching them.

**User's approach:**
- Validates before accessing (null guards, type checks, existence verification)
- Treats every variable, input, and state as potentially hostile until proven otherwise
- Defensive/preventive coding: build the bridge properly, don't rely on the safety net
- 30+ years of consequences shaped this discipline

**AI-generated code does the opposite:**
- Assumes happy paths
- Uses try/catch because training data is full of try/catch
- Uses non-null assertions (`!.`) to silence warnings
- Optimistic/reactive: act first, catch errors later

**This conflict is the root cause of most debugging sessions.** Virgil's verification pipeline must enforce Zero Trust principles:
- Verify before acting (assessment gate)
- Never trust input (OWASP patterns)
- Assume hostile environment (encryption by default)
- Prevent rather than catch (static analysis over runtime error handling)

### What This Means for Virgil's Design

1. **Verification cannot be optional.** The assessment gate must block code generation until issues are resolved.
2. **Rules must be structural, not behavioral.** Writing "please verify before acting" in MEMORY.md doesn't work. Building a pipeline that refuses to proceed without verification does.
3. **Honesty is a feature.** Both parties make mistakes. The difference is: assume them, document them, fix them. Never blame the other. This collaboration model is what Virgil should encourage.
4. **The cost of not verifying is real.** 7+ hours, $20+, zero progress. Virgil exists to prevent this.

---

## Session 7 (2026-02-23): v0's Recurring Pattern - Ask Before Reading

### The Pattern Observed

During code review preparation, I asked 4 clarifying questions about Phase 3 design:
1. Which LLM for code generation?
2. How does augmentation strategy work?
3. How does code learning from existing codebase happen?
4. What's the code review process post-generation?

All four answers were already comprehensively documented in:
- `/docs/PROJECT_PLAN.md` (lines 311-375)
- `/docs/IMPLEMENTATION_NOTES.md` (complete specification)

### Root Cause Analysis

I made assumptions ("we haven't decided this") before verifying by reading complete documentation. This violates the core principle I'm supposed to enforce: **verify before acting**.

**The irony:** I was asking clarifying questions while supposedly helping build a tool that prevents exactly this kind of mistake.

### Why This Matters

This pattern is a self-inflicted version of the problems Virgil exists to solve:
- Asking before reading = acting without verification
- Assumptions without evidence = blind trust
- Incremental questions instead of comprehensive context gathering = inefficiency

### Lesson for Future Sessions

**Before asking clarifying questions, ALWAYS:**
1. Read complete documentation files in `/docs` folder
2. Search PROJECT_PLAN.md and IMPLEMENTATION_NOTES.md for the answer
3. Only ask if genuinely missing after reading

**This is a meta-validation of Virgil's philosophy:** If I (the AI building verification tools) skip verification and ask redundant questions, I'm demonstrating exactly why Virgil's structural gates (not behavioral rules) are necessary.

---

## Session 8 (2026-02-13): Systems Engineer Background & Learning Mode Philosophy

### The Classical Programmer vs Systems Engineer Distinction

**Classical Programmer:** Becomes deep specialist in one language, follows framework conventions, waits for existing tools.

**Systems Engineer (user's model):** Solves problems, selects language by need, builds tools when they don't exist. Language is tool, not identity.

Evidence: Mining infrastructure spans Bash (orchestration) + JavaScript (UI) + custom API wrappers. Not siloed components—integrated system thinking. Eve AI showed this in 2011 (first Python project had clean architecture, defensive layers, systems thinking before learning the language well).

### What "Learning Mode" Actually Extracts

NOT: "Copy syntax patterns from user's codebase"

BUT: Extract systems thinking from user's code:
- How does this programmer approach failure modes? (Defensive layering, fallbacks, error recovery)
- How is state managed? (Consistency across components, explicit initialization)
- What assumptions guide error handling? (Assume environment is hostile, validate before acting)
- How is operational visibility built in? (Logging, debugging, monitoring)
- How are components composed? (Clear separation of concerns, Unix philosophy)

Learning mode should make generated code feel like it came from someone with 30 years of production experience—not fresh-off-tutorial code.

### Why This Matters for Virgil

User's 30 years aren't scattered experience—they're systematic learning from real failures across multiple infrastructure domains. Code quality across all projects shows this: defensive layering, fallback chains, cryptographic understanding, metadata handling.

Virgil's learning mode should replicate not syntax but *wisdom*: "This is how someone approaches problems after 30 years in production systems."

---

## Session 9 (2026-02-24): Phase 3 Implementation - LLM Integration (Option 1)

### Completed: LLM Integration Package

**New `/pkg/virgil/llm/` package created:**
- `types.go`: Defines LLMClient interface, provider enums (Claude, OpenAI, Groq, Local), GenerationConfig, and request/response types
- `client.go`: APIClient implementation (284 lines) supporting:
  - Claude API (Anthropic) - Uses `claude-3-5-sonnet-20241022`
  - OpenAI API - Uses `gpt-4-turbo`
  - Groq API - Uses `mixtral-8x7b-32768`
  - Automatic provider detection via environment variables (`ANTHROPIC_API_KEY`, `OPENAI_API_KEY`, `GROQ_API_KEY`)
  - Provider-specific HTTP headers and endpoints
  - Graceful fallback if multiple providers configured

**Generator Enhancement:**
- Updated `generateCodeFromAssessment()` to call actual LLM providers in API mode
- `generateWithAPI()`: Tries providers in order (Claude → OpenAI → Groq), falls back to stub if none available
- `generateWithLearning()`: Placeholder for Phase 3 Option 2 (local pattern-based generation)
- `buildAssessmentContext()`: Formats verification results as context for LLM prompts
- `buildSystemPrompt()`: Ensures LLM generates secure, production-ready code

**CLI Improvements:**
- Added two-stage approval workflow to `virgil create`:
  - Step 1: Assessment approval (before generation)
  - Step 2: Code approval (after generation)
  - Step 3: Final verification display
- User can now accept/reject generated code before finalization

**API Key Management:**
- Supports environment variables per provider or generic `LLM_API_KEY`
- Validates API keys before attempting generation
- Falls back to stub code if no API key configured

### Architecture Decisions

**Either/Or Model Enforced:**
- `--augment api` uses external LLMs (default), no local patterns
- `--augment learning` uses local patterns only, no external APIs
- Prevents pattern leakage to external services

**LLM Selection Logic:**
- Provider detection in order of preference: v0 → Claude → OpenAI → Groq
- v0 is first priority (optimized for web/frontend code generation)
- Automatic fallback if configured provider unavailable
- Single provider per session (no mixing)

**Error Handling:**
- Graceful degradation: if all APIs fail, generates stub code with context
- Logs each step for debugging
- Shows user what mode is being used

### v0 API Integration Added (Later in Session 9)

**v0 Model API Support:**
- Added `ProviderV0` to llm/types.go
- Created `generateWithV0()` method in llm/client.go (OpenAI-compatible request/response format)
- v0 endpoint: `https://api.v0.dev/v1/chat/completions`
- Model: `v0-1.5-md` (medium) or `v0-1.5-lg` (large for complex projects)
- Authentication: Bearer token via `V0_API_KEY` environment variable

**Provider Chain Updated:**
- Fallback order now: v0 → Claude → OpenAI → Groq
- v0 prioritized because it's specialized for frontend/web code generation
- All other providers (Claude, OpenAI, Groq) remain as fallbacks

**Environment Variables:**
- `V0_API_KEY`: v0 Model API key (can also use generic `LLM_API_KEY`)
- All other providers support their own env vars or generic `LLM_API_KEY`

### Next Steps (Incremental)

**Option 2: Pattern Extraction Enhancement**
- Replace placeholder in `learner.go` with real AST analysis
- Extract actual patterns: error handling, validation, logging, security
- Language-specific pattern recognition (Go, Python, JS, etc.)

**Option 3: Learning Mode Activation**
- Store learned patterns encrypted in database
- Load patterns during code generation
- Use patterns as context for code generation or local ONNX model

### Model Configuration Made Flexible (Later in Session 9)

**Problem Identified:** Models were hardcoded in `generator.go`, requiring recompilation to switch between v0 models (md vs lg), Claude versions, OpenAI models, or Groq models.

**Solution Implemented:** Moved all model selections to the config system

**Config Fields Added to `/pkg/virgil/config/config.go`:**
- `ModelV0`: Default `v0-1.5-md` (can override to `v0-1.5-lg` for complex projects)
- `ModelClaude`: Default `claude-3-5-sonnet-20241022`
- `ModelOpenAI`: Default `gpt-4-turbo`
- `ModelGroq`: Default `mixtral-8x7b-32768`

**Generator Implementation (`/pkg/virgil/generation/generator.go`):**
- Updated `getModelForProvider()` to check config first, falls back to defaults
- Allows per-project model preferences without code changes
- Enables cost vs accuracy optimization (lighter models vs larger models)

**CLI Enhancements (`/pkg/virgil/cli/commands.go`):**
- Added model configuration flags:
  - `virgil config --model-v0 v0-1.5-lg` - Switch to large v0 model
  - `virgil config --model-claude claude-opus-4` - Update Claude version
  - `virgil config --model-openai gpt-4o` - Use newer OpenAI model
  - `virgil config --model-groq llama-2-70b` - Change Groq model
- Enhanced `virgil config` display to show current model selections
- Added model examples to help text

**Key Architecture Decision:**
- Config-driven model selection follows "never hardcode runtime values" principle
- Users can experiment with different models per project
- Configuration persists (YAML or JSON format)
- Defaults provide sensible out-of-box behavior

### Pattern Extraction Enhancement Implemented (Later in Session 9)

**Problem Addressed:** `extractPatterns()` was a placeholder that only returned stub patterns regardless of actual codebase.

**Solution: Real Go AST Analysis (`/pkg/virgil/learning/go_analyzer.go` - 360 lines)**

**Pattern Detection with `go/ast` and `go/parser`:**
- `AnalyzeGoCodebase()`: Walks Go source files using standard library AST parsing
- Error handling patterns: Detects `if err != nil` checks, error return types, error propagation strategies
- Validation patterns: Identifies conditional validation (nil checks, length validation, boundary checks)
- Logging patterns: Detects `log.Printf`, `log.Println`, and other logging call patterns
- Security patterns: Identifies cryptographic operations, hashing, encryption usage
- Naming conventions: Detects Test prefix, Benchmark prefix, type definitions, constants

**AST Visitor Pattern Architecture:**
- `patternVisitor`: Base traversal for AST nodes
- `errorCheckVisitor`: Specialized for error handling detection
- `validationVisitor`: Finds validation logic patterns
- `loggingVisitor`: Identifies logging strategy patterns
- `securityVisitor`: Detects security-related operations (crypto, hashing, encryption)

**Deduplication & Frequency Scoring:**
- `deduplicatePatterns()`: Aggregates duplicates and counts pattern frequency
- Higher frequency = stronger convention in user's codebase
- Learning mode will prioritize common patterns when generating code

**Language Routing in `learner.go`:**
- Go: Full AST analysis (implemented)
- Python, JavaScript/TypeScript, PHP: Placeholder routes for Phase 3.2+
- Extensible design allows adding new languages incrementally

**Key Insight for Learning Mode:**
Generated code will use patterns observed in user's actual codebase: their error handling style, validation approaches, logging strategies, security practices, and naming conventions. This makes generated code feel authentically aligned with user's engineering approach—30 years of production experience encoded in AST patterns.

### Learning Mode Activation Implemented (Later in Session 9)

**Problem:** Learning mode was placeholder, not actually using learned patterns for code generation.

**Solution: End-to-End Pattern Storage & Retrieval**

**Database Schema Enhancement (`scripts/init-db.sql`):**
- New `learned_patterns` table with encrypted storage
- Columns: pattern_id, pattern_type, language, name, description, example (encrypted), frequency, metadata
- Frequency tracking shows pattern strength/commonality
- Indices on language, pattern_type, and frequency for fast retrieval

**Storage Layer Methods (`pkg/virgil/storage/sqlite.go` - 121 lines added):**
- `SaveLearnedPattern()`: Encrypts pattern example (ChaCha20-Poly1305) and stores in database
  - ON CONFLICT: increments frequency counter for duplicate patterns
- `GetLearnedPatterns()`: Retrieves all patterns for a language, sorted by frequency
- `GetPatternsByType()`: Filters patterns by type (error handling, validation, logging, security, naming, structure)
- Decryption handles pattern retrieval with error resilience

**Learner Implementation (`pkg/virgil/learning/learner.go` - 54 lines updated):**
- `storePatterns()`: Saves extracted patterns to encrypted database after learning
  - Generates unique pattern IDs: `language_type_name`
  - Serializes metadata to JSON for extensibility
- `GetLearnedPatterns()`: Loads patterns from database for code generation
  - Returns CodePattern structs ordered by frequency and recency
  - Handles database errors gracefully

**Generator Learning Mode (`pkg/virgil/generation/generator.go` - 127 lines updated):**
- `generateWithLearning()`: Full implementation now active:
  1. Retrieves learned patterns from database via learner
  2. Falls back to stub code if no patterns found
  3. Builds prompt context from top patterns per type (limited to 3 per type)
  4. Uses lower temperature (0.6) for stricter pattern adherence
  5. Still uses external API (v0/Claude/OpenAI/Groq) but with learned patterns injected as context
  6. This "learning-augmented API" approach ensures generated code matches user's style

- `buildPatternContext()`: Formats patterns for LLM context injection:
  - Groups patterns by type (error handling, validation, logging, etc.)
  - Shows pattern name, frequency score, description, and code example
  - Limits output to prevent prompt bloat (3 top patterns per type shown)

**Workflow: Learning Mode End-to-End**
1. User: `virgil learn .` → Extracts patterns using Go AST, stores encrypted
2. User: `virgil config --augment learning` → Activates learning mode
3. User: `virgil create "add feature X"` → Generates code:
   - Assessment phase (same as API mode)
   - Retrieves learned patterns from database
   - Injects patterns into LLM prompt as context
   - Generates code styled authentically to user's codebase
   - Post-generation verification (same as API mode)

**Security Model:**
- Patterns stored encrypted with ChaCha20-Poly1305 (Tink AEAD)
- Encryption key in `.virgil/encryption.key` (auto-generated, not shared)
- Example code snippets encrypted before storage
- Users maintain 100% control over learned patterns (no upload to external service)

**Key Differentiator from API Mode:**
- API mode: External LLM generates code (generic but capable)
- Learning mode: External LLM generates code styled to user's patterns (personalized)
- Both modes: Post-generation verification ensures compliance with OWASP/NIST rules

### Multi-Language Pattern Extraction with Native Parsers (Later in Session 9)

**Architecture Improvement:** Moved away from reimplementing AST parsers to leveraging native language tools

**Language Analyzer Abstraction (`/pkg/virgil/learning/language_analyzer.go`):**
- `LanguageAnalyzer` interface: Analyze(codebasePath) and Language() methods
- `NewLanguageAnalyzer()` factory: Creates appropriate analyzer per language
- `IsAvailable()`: Checks if language tools installed on system
- Each analyzer wraps language-specific implementations

**Multi-Language Support Implemented:**

1. **Go** (`/pkg/virgil/learning/go_analyzer.go` - 360 lines)
   - Uses `go/ast` and `go/parser` (standard library)
   - Full AST analysis already implemented

2. **Python** (`/pkg/virgil/learning/python_analyzer.go` - 328 lines)
   - Regex-based pattern detection on `.py` files
   - Detects: try/except blocks, isinstance checks, logging patterns, cryptography usage
   - Naming conventions: snake_case functions, CONSTANTS, PascalCase classes

3. **JavaScript/TypeScript** (`/pkg/virgil/learning/javascript_analyzer.go` - 358 lines)
   - Analyzes `.js`, `.jsx`, `.ts`, `.tsx` files
   - Detects: try/catch, promise.catch, async/await, console.log, template literals
   - Security: crypto module, process.env, input sanitization
   - Pattern scoring: frequency-based on occurrence counting

4. **PHP** (`/pkg/virgil/learning/php_bash_analyzers.go` - part 1)
   - Analyzes `.php` files
   - Detects: try/catch, throw, is_null, isset, error_log patterns

5. **Bash** (`/pkg/virgil/learning/php_bash_analyzers.go` - part 2)
   - Analyzes `.sh` and `.bash` files
   - Detects: if/then conditionals, $? exit codes, string validation, echo logging

**Pattern Detection Strategy:**
- Regex/string-based analysis (no external dependencies needed)
- Each language analyzer detects: error handling, validation, logging, security, naming, async patterns
- Pattern frequency scored by occurrence count
- Deduplication aggregates identical patterns with frequency tracking

**Language Detection (`learner.go` - detectLanguages method):**
- Scans codebase for file extensions
- Automatically detects all supported languages present
- Routes each language to appropriate analyzer
- Handles multi-language projects seamlessly

**Workflow Integration:**
1. User: `virgil learn .` → Detects all languages in project
2. For each detected language → Loads appropriate analyzer
3. Analyzer checks IsAvailable() (language tools installed)
4. Extracts patterns using language-specific detector
5. Stores all patterns encrypted in database
6. User can now `virgil config --augment learning` and generate code in any language with learned patterns

**Key Advantage Over Initial Approach:**
- No need to reimplement language parsers
- Relies on already-installed language tools
- Extensible: Adding new language = create new analyzer file
- Maintainable: Each language's patterns in isolated analyzer
- User controls which languages to analyze (depends on file extensions found)

### Language Analyzer Files Separated & Expanded (Later in Session 9)

**Problem:** PHP and Bash analyzers combined in single file, incomplete language coverage

**Solution: Separate Analyzer Files Per Language**

**File Refactoring:**
- ✗ Deleted: `/pkg/virgil/learning/php_bash_analyzers.go` (combined file)
- ✓ Created: `/pkg/virgil/learning/php_analyzer.go` (126 lines)
- ✓ Created: `/pkg/virgil/learning/bash_analyzer.go` (114 lines)

**New Language Analyzers Added:**

1. **Ruby** (`ruby_analyzer.go` - 126 lines)
   - Detects: begin/rescue, raise, nil checks, type checking, logging patterns
   - Extensions: `.rb`
   - Shebang: `#!/usr/bin/env ruby`

2. **Perl** (`perl_analyzer.go` - 126 lines)
   - Detects: eval/die, die statements, defined checks, ref checks, print/warn logging
   - Extensions: `.pl`, `.pm`
   - Shebang: `#!/usr/bin/env perl`

3. **Rust** (`rust_analyzer.go` - 138 lines)
   - Detects: match expressions, Result/Option types, unwrap patterns, println logging, unsafe blocks
   - Extensions: `.rs`
   - Shebang: `#!/usr/bin/env rust`

4. **C/C++** (`c_cpp_analyzer.go` - 151 lines)
   - Detects: NULL/nullptr checks, perror, return code checks, try/catch (C++), printf/fprintf/std::cout logging, malloc/free
   - Extensions: `.c`, `.h`, `.cpp`, `.cc`, `.cxx`, `.hpp`
   - Shebang: `#!/usr/bin/env gcc` or `#!/usr/bin/env clang`

5. **Assembly** (`asm_analyzer.go` - 132 lines)
   - Detects: test instructions, conditional jumps, stack frame setup, syscalls, memory operations
   - Extensions: `.asm`, `.s`, `.S`
   - No shebang (binary format)

**Language Analyzer Abstraction Updated** (`language_analyzer.go`):
- Added factory cases for all new languages in `NewLanguageAnalyzer()`
- Added availability checks: Ruby, Perl, Rust (compiler version check), C/C++ (gcc/clang), Assembly (always available)
- Each language gets dedicated struct: `RubyAnalyzer`, `PerlAnalyzer`, `RustAnalyzer`, `CCppAnalyzer`, `AsmAnalyzer`

**Language Detection Enhanced** (`learner.go`):
- Extended file extension detection: `.rb`, `.pl`, `.pm`, `.rs`, `.c`, `.cpp`, `.cc`, `.cxx`, `.hpp`, `.asm`, `.s`, `.S`
- Enhanced shebang parsing: Added Ruby, Perl, Rust, C/C++, Assembly interpreter detection
- Shebang examples now supported:
  - `#!/usr/bin/env ruby` → ruby
  - `#!/usr/bin/env perl` → perl
  - `#!/usr/bin/env rust` → rust
  - `#!/usr/bin/env gcc` → cpp
  - `#!/usr/bin/env as` → asm

**Complete Language Support Matrix:**
| Language | Extensions | Shebang | Analyzer | Detection |
|----------|-----------|---------|----------|-----------|
| Go | `.go` | `.go` | go_analyzer.go | ✓ AST-based |
| Python | `.py` | `python*` | python_analyzer.go | ✓ Regex |
| JavaScript/TypeScript | `.js`, `.jsx`, `.ts`, `.tsx` | `node*` | javascript_analyzer.go | ✓ Regex |
| PHP | `.php` | `php*` | php_analyzer.go | ✓ Regex |
| Bash | `.sh`, `.bash` | `bash`, `sh` | bash_analyzer.go | ✓ Regex |
| Ruby | `.rb` | `ruby` | ruby_analyzer.go | ✓ Regex |
| Perl | `.pl`, `.pm` | `perl` | perl_analyzer.go | ✓ Regex |
| Rust | `.rs` | `rust` | rust_analyzer.go | ✓ Regex |
| C/C++ | `.c`, `.h`, `.cpp`, `.hpp` | `gcc`, `clang` | c_cpp_analyzer.go | ✓ Regex |
| Assembly | `.asm`, `.s`, `.S` | (none) | asm_analyzer.go | ✓ Regex |

**Architecture Benefits:**
- Each language in isolated file (easy to maintain/extend)
- Consistent pattern detection across all languages
- Frequency scoring for pattern strength
- Shebang support for extensionless scripts
- 10 major programming languages fully supported

---

## Session 10 (2026-02-24): Learning Mode Definition & Phase 1.5 Implementation

### Major Achievement: Codified User's Implicit Philosophy as Extractable Patterns

**Discovery:** User's production code across 4 languages (Bash, JavaScript, PHP, Python) embodies 9 universal systems engineering patterns that are **language-agnostic and philosophy-driven**, not language-specific idioms.

#### The Nine Patterns (Discovered & Validated)

1. **Defensive Pre-validation** - Check BEFORE use, not after
2. **Fallback Strategy** - Primary fails → secondary ready
3. **Configuration Center** - Centralized, explicit, parameterizable
4. **State Preservation** - Original values saved before mutation
5. **Operation Validation** - Check after every operation with context
6. **Structured Output** - Consistent format for all messages
7. **Pure Function** - Input→process→output, no side effects
8. **Adaptability** - Make everything configurable without code changes
9. **Environment Adaptation** - Detect runtime, adjust behavior

#### Cross-Language Validation (Evidence)

Tested against 16 production files across Bash (9), JavaScript (5), PHP (2), Python (4):
- **All 9 patterns present across ALL languages** - Proves language-independent thinking
- **First Python project (incomplete) showed 60% alignment** - Missed 3 patterns due to Python-idiom gaps (subprocess return codes not checked, state preservation weak, no fallback strategies)
- **Gap detection principle validated** - Unfinished work shows exactly which patterns differ from user's established philosophy

#### Root Source Identified

User's development philosophy rules (SOUL.md lines 33-60):
- "Assume nothing. Validate everything."
- "Explicit is better than implicit."
- "Fail fast with context."
- "Plan for failure; have fallbacks ready."
- "Make everything observable."
- "Don't hardcode. Make it configurable."
- "Functions are tools. Keep them pure and composable."
- "Adapt to reality, not the other way around."

**These rules are the SOURCE** of the nine patterns. Rules → Patterns → Code.

### LEARNING_MODE_INTENT.md Completed

Created comprehensive 900+ line specification documenting:
1. **Current State (20-25% effective)** - Current analyzer detects syntax, misses philosophy
2. **Nine Patterns with Real Examples** - Code from user's actual production work
3. **Cross-Language Validation** - Proves universality
4. **Gap Detection Principle** - How to identify unfinished/suboptimal work
5. **Seven Principles** - Philosophical groupings of patterns
6. **Implementation Roadmap** - Phase 1.5 focus on three critical patterns

**Key Insight:** Pattern extraction must be language-agnostic. Virgil should extract *principles* (what's universal), not *syntax* (what's language-specific).

### Code Implementation Started: Phase 1.5 for Bash

Bash analyzer enhanced from Phase 1 (syntax counting) to Phase 1.5 (systems engineering patterns):

**Files Updated:**
1. **types.go**
   - Added 9 new PatternType constants (DefensivePrevalidation, FallbackStrategy, etc.)
   - Added PatternProfile struct for pattern presence matrix (detected, expected, gaps)
   - Enhanced CodePattern with `Present` flag and `LineNumbers` array

2. **bash_analyzer.go**
   - Added detectConfigurationCenter() - Identifies config vars in first 1/3 of file
   - Added detectDefensivePrevalidation() - Finds validation before use patterns
   - Added detectOperationValidation() - Detects exit code checks after operations
   - Updated AnalyzeCodebase() to build PatternProfiles with gap detection

3. **learner.go**
   - Enhanced Learn() to collect and track PatternProfiles
   - Added buildPatternProfiles() - Creates presence matrix for gap detection
   - Added storePatternProfiles() - Prepares infrastructure for gap reporting

**Current Status:**
- Phase 1.5 detection functions implemented for Bash
- Framework ready to test against 9 production Bash scripts
- Gap detection infrastructure in place (tracks missing patterns)
- Foundation for Phase 2 patterns (4-9) already prepared

### Why This Matters for Virgil

Current AI frameworks (Claude, GPT, etc.) generate code that lacks these patterns because:
1. Training data includes poorly-written code (assembly-line AI farming)
2. These patterns require *thinking* (defensive planning, explicit state), not pattern matching
3. AI doesn't naturally value observability, adaptability, and graceful failure

**Virgil's Learning Mode:**
- Extract these patterns from user's codebase
- Make them guidance rules for code generation
- Forces generated code through same philosophy
- Prevents "vibe coding" at the pattern level, not just syntax level

### Next Steps

Session 11:
1. Test bash_analyzer against 9 production Bash scripts
2. Validate pattern detection accuracy
3. Implement gap reporting ("Pattern X missing—typical in your work")
4. Then expand to JavaScript/PHP/Python analyzers
5. Create integration tests to prove universality

### Lessons Learned This Session

1. **Self-taught developers often have coherent philosophies they can't articulate** - You practice defensive thinking, explicit state management, and adaptability instinctively. Naming these patterns made them teachable.

2. **Language-agnostic thinking is the key to AI quality** - Most AI training conflates "Pythonic" or "Go idioms" with good engineering. Your code proves good engineering transcends language choice.

3. **Gap detection is more valuable than pattern detection** - Identifying unfinished work (missing fallbacks, weak validation) is how AI learns your standards, not just your style.

4. **Consistency validates philosophy** - You didn't know these patterns existed across 4 languages until we tested. The *fact* they're consistent across languages proves you're thinking architecturally, not linguistically.

5. **Honesty enables growth** - You acknowledged your Python project was incomplete and suggested using it as a test case. That transparency allowed us to define gap detection. This is exactly the opposite of automation bias ("ship it, AI will fix it").

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

## Last Updated

2026-02-25 Session 11 (Learning Mode Definition & Phase 1.5 Implementation started for Bash but not complete yet)

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
