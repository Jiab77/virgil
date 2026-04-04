# Virgil

Verification framework for disciplined AI-assisted development.

Named after the Roman poet who guides Dante through Hell and Purgatory in the *Divine Comedy*. **Virgil** serves as guide, protector, teacher, and companion—helping developers navigate the chaos of AI-generated code with confidence.

## Overview

**Virgil** prevents catastrophic security failures in AI-generated code by providing:

- **Assessment before generation** - Understand requirements before AI creates code
- **Comprehensive verification** - Check code against security standards (OWASP, NIST, GDPR)
- **Pattern learning** - Learn from your codebase and apply your standards consistently
- **Transparent decisions** - See exactly why code was approved or rejected

## Project Status

**Phase 1 & 2: Complete** ✅
- Core orchestrator with CLI commands (init, config, create, edit, assess, review)
- Verification pipeline with parallel block execution
- Real OWASP implementation + 7 stub compliance blocks (NIST, GDPR, HIPAA, PCI-DSS, CIS, ISO27001, custom)
- Encryption integration (Tink/ChaCha20-Poly1305) with secure storage
- Web search infrastructure with encrypted caching

**Phase 3: In Progress**
- Code generation from descriptions (LLM integration)
- Assessment phase intelligence
- Augmentation strategy implementation

For detailed status, see [COMPLETION_SUMMARY.md](./docs/COMPLETION_SUMMARY.md)

## Quick Start

```bash
# Initialize a new project
virgil init

# Configure Your Augmentation Strategy
# Local patterns only: learns from your codebase, no external API calls
virgil config --augment learning

# Can use external LLMs (Claude/GPT) with fallback support
virgil config --augment api

# Configure Running Mode
virgil config --mode plan-first   # Establish a project plan first before writing code
# OR
virgil config --mode fast         # Skip plan creation and write code directly (review it later)

# Create new code with assessment gate
virgil create "add user authentication"

# Edit existing code with guidance
virgil edit "add password validation"

# Edit specific file (targeted)
virgil edit "fix hardcoded API keys" pkg/config/

# Assess existing code
virgil assess
virgil assess pkg/auth/handler.go           # Assess specific file
virgil assess internal/                     # Assess directory
virgil assess --augment learning pkg/auth/  # Assess with learned patterns

# View verification results
virgil review

# Get help
virgil --help
virgil config --help
```

## Verification Gates

Virgil enforces mandatory verification before code generation:

1. **Assessment Phase** - Analyzes requirements and existing code
2. **User Approval** - You review assessment results and approve/reject before generation
3. **Code Generation** - Only proceeds after explicit approval
4. **Post-Generation Review** - Generated code is verified against OWASP and custom rules

This prevents "approve blindly" patterns that lead to security failures.

## CLI Commands

- `virgil init` - Initialize a new Virgil project
- `virgil config` - Configure settings (--augment, --mode, --rules)
- `virgil create <description>` - Create new code with assessment gate (primary workflow)
- `virgil edit <description> [path]` - Edit existing code with assessment gate
- `virgil assess [path]` - Assess code against verification rules (utility command)
- `virgil review` - View verification audit trail

## Options

- `--augment api|learning` - Set augmentation strategy (default: api)
- `--mode plan-first|fast` - Set execution mode (default: plan-first)
- `--rules <rules>` - Configure verification rules (comma-separated)
- `--help` - Show help for any command
- `--version` - Show version

## Supported Compliance Standards

**Default (International):**
- `owasp` - OWASP Top 10 security best practices
- `nist` - NIST security guidelines

**Optional (Regional/Industry-Specific):**
- `gdpr` - EU/EEA data protection regulations
- `hipaa` - USA healthcare data protection
- `pci-dss` - Payment Card Industry data security
- `cis` - CIS Controls security framework
- `iso27001` - ISO/IEC 27001 information security management
- `custom` - User-defined rules

## Web Search Integration

Virgil caches web search results for context during verification (Phase 3+):
- SHA256 hashing for query deduplication
- 7-day TTL with encrypted storage
- Graceful degradation if search unavailable
- All queries/results encrypted end-to-end

## Examples

```bash
# Use defaults (OWASP + NIST)
virgil config --rules owasp,nist

# Add GDPR for EU/EEA projects
virgil config --rules owasp,nist,gdpr

# Add HIPAA for healthcare projects
virgil config --rules owasp,nist,hipaa

# Add PCI-DSS for payment processing
virgil config --rules owasp,nist,pci-dss
```

## Documentation

- [Project Plan](./docs/PROJECT_PLAN.md) - Vision, architecture, and implementation roadmap
- [Memory](MEMORY.md) - Quick reference for development decisions
- [Implementation Notes](./docs/IMPLEMENTATION_NOTES.md) - Critical considerations for Phase 1+ development

## Development

Go 1.25+

```bash
go run ./cmd/virgil --help
```

## License

MIT License - See [LICENSE](./LICENSE) file for details.

This project is open-source and free to use, modify, and distribute.

## Credits

* __[Jiab77](https://github.com/Jiab77)__
* __[v0](https://v0.dev)__
