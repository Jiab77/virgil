// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/jiab77/virgil/pkg/virgil/config"
	"github.com/jiab77/virgil/pkg/virgil/generation"
	"github.com/jiab77/virgil/pkg/virgil/learning"
	"github.com/jiab77/virgil/pkg/virgil/storage"
	"github.com/jiab77/virgil/pkg/virgil/verification"
)

// NewRootCommand creates the root Virgil command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "virgil",
		Short: "Virgil - Code Verification Framework",
		Long: `Virgil is a verification framework for disciplined AI-assisted development.

It prevents catastrophic security failures by enforcing verification gates
before code is generated and deployed.`,
		Version: "0.1.0",
	}

	rootCmd.AddCommand(
		newInitCommand(),
		newConfigCommand(),
		newCreateCommand(),
		newEditCommand(),
		newAssessCommand(),
		newReviewCommand(),
		newLearnCommand(),
	)

	return rootCmd
}

// newInitCommand creates the 'init' command
func newInitCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initialize a new Virgil project",
		Long: `Initialize a new Virgil project in the current directory.

Creates .virgil/ directory with configuration and encrypted database.
Prompts for encryption strategy: random key (default) or passphrase-based.`,
		Example: `  virgil init`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing Virgil project...")

			// Create .virgil directory
			virgilDir := ".virgil"
			if err := os.MkdirAll(virgilDir, 0700); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating .virgil directory: %v\n", err)
				os.Exit(1)
			}

			// Prompt for encryption strategy
			fmt.Println("\nEncryption Strategy:")
			fmt.Println("1. Random key (default, stored in .virgil/key)")
			fmt.Println("2. Passphrase (user-provided, derived with Argon2)")
			fmt.Print("\nChoose encryption strategy [1/2]: ")

			reader := bufio.NewReader(os.Stdin)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			if choice == "" {
				choice = "1"
			}

			// Initialize database with chosen encryption
			dbPath := filepath.Join(virgilDir, "virgil.db")
			keyPath := filepath.Join(virgilDir, "key")

			var key []byte
			var err error

			if choice == "1" {
				// Random key
				key, err = storage.GenerateRandomKey()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error generating encryption key: %v\n", err)
					os.Exit(1)
				}
				// Store key in file
				if err := os.WriteFile(keyPath, key, 0600); err != nil {
					fmt.Fprintf(os.Stderr, "Error storing encryption key: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("✓ Generated random encryption key (stored in .virgil/key)")
			} else if choice == "2" {
				// Passphrase-based
				fmt.Print("Enter passphrase for database encryption: ")
				passphrase, _ := reader.ReadString('\n')
				passphrase = strings.TrimSpace(passphrase)

				key, err = storage.DeriveKeyFromPassphrase(passphrase)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error deriving encryption key: %v\n", err)
					os.Exit(1)
				}
				fmt.Println("✓ Derived encryption key from passphrase (no key file)")
			} else {
				fmt.Println("Invalid choice. Using random key.")
				key, _ = storage.GenerateRandomKey()
				os.WriteFile(keyPath, key, 0600)
			}

			// Initialize database
			if err := storage.InitializeDatabase(dbPath, key); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				os.Exit(1)
			}

			// Create default config
			cfg := config.NewConfig()
			fmt.Println("✓ Created default configuration (owasp, nist)")
			fmt.Println("✓ Initialized encrypted database")
			fmt.Printf("\nProject initialized successfully!\n")
			fmt.Printf("Next: virgil create \"your feature description\"  # Start the workflow\n")
			fmt.Printf("      virgil config --rules gdpr               # Add more rules if needed\n")
		},
	}
}

// newConfigCommand creates the 'config' command
func newConfigCommand() *cobra.Command {
	var augment string
	var mode string
	var rules string
	var webSearch string
	var modelV0 string
	var modelClaude string
	var modelOpenAI string
	var modelGroq string

	cmd := &cobra.Command{
		Use:   "config [flags]",
		Short: "Configure Virgil settings",
		Long: `Configure Virgil settings for this project.

Configuration options:
  --augment api|learning    Set augmentation strategy
                            - api (default): Use external APIs (Claude/GPT)
                            - learning: Use learned patterns from your codebase
  --mode plan-first|fast    Set execution mode
                            - plan-first (default): Assessment before generation
                            - fast: Skip assessment step
  --rules <rules>           Set verification rules (comma-separated)
                            Default: owasp, nist (international standards)
                            Optional: gdpr (EU/EEA), hipaa (USA healthcare),
                                     pci-dss (payment cards), cis (controls),
                                     iso27001 (information security),
                                     custom (user-defined)
  --web-search enabled|disabled  Enable/disable web search module
                            - enabled (default): Ground assessment in current research
                            - disabled: Use static rules only
  --model-v0 <model>        Set v0 model (default: v0-1.5-md, or v0-1.5-lg)
  --model-claude <model>    Set Claude model (default: claude-3-5-sonnet-20241022)
  --model-openai <model>    Set OpenAI model (default: gpt-4-turbo)
  --model-groq <model>      Set Groq model (default: mixtral-8x7b-32768)

When no flags are provided, displays current configuration.`,
		Example: `  virgil config                           # Show current configuration
  virgil config --augment learning         # Use learned patterns
  virgil config --mode fast                # Skip assessment gate
  virgil config --rules owasp,nist,gdpr    # Add GDPR (EU/EEA users)
  virgil config --rules owasp,nist,hipaa   # Add HIPAA (USA healthcare)
  virgil config --rules owasp,nist,pci-dss # Add PCI-DSS (payment processing)
  virgil config --rules owasp,nist,custom  # Add custom rules
  virgil config --web-search disabled      # Disable web search
  virgil config --model-v0 v0-1.5-lg       # Use v0 large model for complex projects
  virgil config --model-claude claude-opus-4 # Switch to Claude Opus`,
		Run: func(cmd *cobra.Command, args []string) {
			// Load current configuration
			cfg, err := config.LoadConfig(".")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// If no flags provided, show current configuration
			if augment == "" && mode == "" && rules == "" && webSearch == "" && modelV0 == "" && modelClaude == "" && modelOpenAI == "" && modelGroq == "" {
				fmt.Println("=== Current Configuration ===\n")
				fmt.Printf("Augment Strategy: %s\n", cfg.AugmentStrategy)
				fmt.Printf("Execution Mode: %s\n", cfg.Mode)
				fmt.Printf("Rules Enabled: %s\n", strings.Join(cfg.RulesEnabled, ", "))
				webSearchStatus := "enabled"
				if !cfg.WebSearchEnabled {
					webSearchStatus = "disabled"
				}
				fmt.Printf("Web Search: %s\n\n", webSearchStatus)

				// Show model configuration
				fmt.Println("=== LLM Model Configuration ===")
				fmt.Printf("v0 Model: %s\n", cfg.ModelV0)
				fmt.Printf("Claude Model: %s\n", cfg.ModelClaude)
				fmt.Printf("OpenAI Model: %s\n", cfg.ModelOpenAI)
				fmt.Printf("Groq Model: %s\n", cfg.ModelGroq)
				
				// Show config format
				format := config.GetConfigFormat(".")
				if format != "none" {
					fmt.Printf("\nConfig Format: %s (.virgil/config.%s)\n", format, format)
				}
				return
			}

			// Update configuration with provided flags
			if augment != "" {
				if augment != "api" && augment != "learning" {
					fmt.Fprintf(os.Stderr, "Invalid augment strategy: %s (must be 'api' or 'learning')\n", augment)
					os.Exit(1)
				}
				cfg.AugmentStrategy = augment
			}

			if mode != "" {
				if mode != "plan-first" && mode != "fast" {
					fmt.Fprintf(os.Stderr, "Invalid mode: %s (must be 'plan-first' or 'fast')\n", mode)
					os.Exit(1)
				}
				cfg.Mode = mode
			}

			if rules != "" {
				cfg.RulesEnabled = strings.Split(rules, ",")
				// Trim whitespace from each rule
				for i, rule := range cfg.RulesEnabled {
					cfg.RulesEnabled[i] = strings.TrimSpace(rule)
				}
			}

			if webSearch != "" {
				if webSearch != "enabled" && webSearch != "disabled" {
					fmt.Fprintf(os.Stderr, "Invalid web search option: %s (must be 'enabled' or 'disabled')\n", webSearch)
					os.Exit(1)
				}
				cfg.WebSearchEnabled = webSearch == "enabled"
			}

			// Update model configuration if provided
			if modelV0 != "" {
				cfg.ModelV0 = modelV0
			}
			if modelClaude != "" {
				cfg.ModelClaude = modelClaude
			}
			if modelOpenAI != "" {
				cfg.ModelOpenAI = modelOpenAI
			}
			if modelGroq != "" {
				cfg.ModelGroq = modelGroq
			}

			// Save configuration to file
			if err := config.SaveConfig(cfg, "."); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving configuration: %v\n", err)
				os.Exit(1)
			}

			// Display updated configuration
			fmt.Println("✓ Configuration updated:\n")
			fmt.Printf("Augment Strategy: %s\n", cfg.AugmentStrategy)
			fmt.Printf("Execution Mode: %s\n", cfg.Mode)
			fmt.Printf("Rules Enabled: %s\n", strings.Join(cfg.RulesEnabled, ", "))
			webSearchStatus := "enabled"
			if !cfg.WebSearchEnabled {
				webSearchStatus = "disabled"
			}
			fmt.Printf("Web Search: %s\n\n", webSearchStatus)

			// Display updated model configuration
			fmt.Println("=== LLM Model Configuration ===")
			fmt.Printf("v0 Model: %s\n", cfg.ModelV0)
			fmt.Printf("Claude Model: %s\n", cfg.ModelClaude)
			fmt.Printf("OpenAI Model: %s\n", cfg.ModelOpenAI)
			fmt.Printf("Groq Model: %s\n", cfg.ModelGroq)
		},
	}

	cmd.Flags().StringVar(&augment, "augment", "", "augmentation strategy (api|learning)")
	cmd.Flags().StringVar(&mode, "mode", "", "execution mode (plan-first|fast)")
	cmd.Flags().StringVar(&rules, "rules", "", "verification rules (default: owasp,nist; optional: gdpr,hipaa,pci-dss,cis,iso27001,custom)")
	cmd.Flags().StringVar(&webSearch, "web-search", "", "web search module (enabled|disabled)")
	cmd.Flags().StringVar(&modelV0, "model-v0", "", "v0 model selection")
	cmd.Flags().StringVar(&modelClaude, "model-claude", "", "Claude model selection")
	cmd.Flags().StringVar(&modelOpenAI, "model-openai", "", "OpenAI model selection")
	cmd.Flags().StringVar(&modelGroq, "model-groq", "", "Groq model selection")

	return cmd
}

// newCreateCommand creates the 'create' command
func newCreateCommand() *cobra.Command {
	var augment string

	cmd := &cobra.Command{
		Use:   "create <description>",
		Short: "Create new code with assessment gate (primary workflow)",
		Long: `Create new code based on your description.

This is the primary Virgil workflow. It guides you through a disciplined
process to ensure generated code meets security requirements.

Workflow:
  1. ASSESS - Framework analyzes current code and requirements
  2. USER DECISION - You review the approach and approve/reject
  3. GENERATE - Code is created (Phase 3)
  4. VERIFY - Generated code is verified (Phase 3)

Flags:
  --augment api|learning    Use specific augmentation strategy
                            - api: External APIs (Claude/GPT)
                            - learning: Learned patterns from your codebase`,
		Example: `  virgil create "add user authentication"
  virgil create "add user authentication" --augment learning
  virgil create "implement password hashing with bcrypt"`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := args[0]
			augmentStrategy := ""
			if augment != "" {
				augmentStrategy = augment
			}

			fmt.Println("=== Code Creation Workflow ===\n")
			fmt.Printf("Request: %s\n", description)
			if augmentStrategy != "" {
				fmt.Printf("Strategy: %s\n\n", augmentStrategy)
			}

			// Step 1: Assess
			fmt.Println("Step 1: ASSESS - Evaluating codebase and requirements...")
			
			// Load configuration (with flag overrides)
			cfg, err := config.LoadConfig(".")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// Apply augment strategy override if provided via flag
			if augmentStrategy != "" {
				cfg.AugmentStrategy = augmentStrategy
			}

			// Load database connection
			dbPath := filepath.Join(".virgil", "virgil.db")
			db, err := storage.InitDatabase(dbPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			// Run verification pipeline with web search (if enabled)
			results, err := verification.RunPipeline(description, ".", cfg, db)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running assessment: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("  Found %d issues in assessment\n", len(results.AllIssues))
			if len(results.AllIssues) > 0 {
				fmt.Println("  Issues to address:")
				for _, issue := range results.AllIssues {
					fmt.Printf("    - [%s] %s (in %s:%d)\n", issue.Severity, issue.Message, issue.FilePath, issue.LineNumber)
				}
			}

			// Step 2: Ask for approval
			fmt.Println("\nStep 2: USER DECISION - Review approach")
			fmt.Println("Proposed approach:")
			fmt.Printf("  1. Analyze: %s\n", description)
			fmt.Printf("  2. Plan: Create implementation using %s\n", cfg.AugmentStrategy)
			fmt.Printf("  3. Verify: Ensure compliance with %s rules\n", strings.Join(cfg.RulesEnabled, ", "))
			fmt.Print("\nApprove this approach? [y/n]: ")

			reader := bufio.NewReader(os.Stdin)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			if choice != "y" && choice != "Y" {
				fmt.Println("Aborted.")
				return
			}

			// Step 3: Generate code with Phase 3 implementation
			fmt.Println("\nStep 3: GENERATE - Creating code...")

			// Initialize generator
			gen := generation.NewGenerator(db, cfg)

			// Prepare generation request
			genRequest := &generation.CodeGenerationRequest{
				Description: description,
				ProjectPath: ".",
				Language:    "go", // Default language, could be parameterized
				Context:     "",
			}

			// Generate code
			genResponse, err := gen.GenerateCode(genRequest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
				os.Exit(1)
			}

			// Display generated code
			fmt.Println("\n=== Generated Code ===\n")
			fmt.Println(genResponse.Code)

			// Step 4: User approval of generated code
			fmt.Println("\nStep 4: USER DECISION - Review generated code")
			fmt.Print("Accept generated code? [y/n]: ")

			choice2, _ := reader.ReadString('\n')
			choice2 = strings.TrimSpace(choice2)

			if choice2 != "y" && choice2 != "Y" {
				fmt.Println("Code rejected. Please refine your requirements and try again.")
				return
			}

			// Step 5: Post-generation verification
			fmt.Println("\nStep 5: VERIFY - Final verification")
			if genResponse.VerificationResult != nil {
				if len(genResponse.VerificationResult.FailedChecks) == 0 {
					fmt.Println("✓ All verification checks passed")
				} else {
					fmt.Printf("⚠ %d verification issues found. Review before committing.\n", len(genResponse.VerificationResult.FailedChecks))
				}
			}

			fmt.Println("\n✓ Workflow complete. Generated code is ready for review.")
		},
	}

	cmd.Flags().StringVar(&augment, "augment", "", "augmentation strategy (api|learning)")

	return cmd
}

// newEditCommand creates the 'edit' command
func newEditCommand() *cobra.Command {
	var augment string

	cmd := &cobra.Command{
		Use:   "edit <description> [path]",
		Short: "Edit existing code with assessment gate",
		Long: `Edit existing code based on your description.

This workflow guides you through modifying existing code with the same
verification gates as code creation. Use this to fix issues, add features,
or improve existing implementations.

If no path is provided, analyzes the entire project and suggests changes.
If a path is provided, focuses the edit on that specific file or directory.

Workflow:
  1. ASSESS - Analyze current code and proposed changes
  2. USER DECISION - Review the approach and approve/reject
  3. GENERATE - Modified code is created (Phase 3)
  4. VERIFY - Updated code is verified (Phase 3)

Flags:
  --augment api|learning    Use specific augmentation strategy
                            - api: External APIs (Claude/GPT)
                            - learning: Learned patterns from your codebase`,
		Example: `  virgil edit "add password validation"
  virgil edit "add password validation" pkg/auth/handler.go
  virgil edit "fix hardcoded API keys" app/
  virgil edit "add parameterized queries" pkg/db/query.go --augment learning`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			description := args[0]
			path := "."
			if len(args) > 1 {
				path = args[1]
			}
			augmentStrategy := ""
			if augment != "" {
				augmentStrategy = augment
			}

			fmt.Println("=== Code Edit Workflow ===\n")
			fmt.Printf("Change: %s\n", description)
			fmt.Printf("Target: %s\n", path)
			if augmentStrategy != "" {
				fmt.Printf("Strategy: %s\n\n", augmentStrategy)
			}

			// Step 1: Assess
			fmt.Println("Step 1: ASSESS - Analyzing code and proposed changes...")
			
			// Load configuration (with flag overrides)
			cfg, err := config.LoadConfig(".")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// Apply augment strategy override if provided via flag
			if augmentStrategy != "" {
				cfg.AugmentStrategy = augmentStrategy
			}

			// Load database connection
			dbPath := filepath.Join(".virgil", "virgil.db")
			db, err := storage.InitDatabase(dbPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			// Run verification pipeline with web search (if enabled)
			results, err := verification.RunPipeline(description, path, cfg, db)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running assessment: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("  Found %d issues in target code\n", len(results.AllIssues))
			if len(results.AllIssues) > 0 {
				fmt.Println("  Issues to address:")
				for _, issue := range results.AllIssues {
					fmt.Printf("    - [%s] %s (in %s:%d)\n", issue.Severity, issue.Message, issue.FilePath, issue.LineNumber)
				}
			}

			// Step 2: Ask for approval
			fmt.Println("\nStep 2: USER DECISION - Review approach")
			fmt.Println("Proposed approach:")
			fmt.Printf("  1. Change: %s\n", description)
			fmt.Printf("  2. Target: %s\n", path)
			fmt.Printf("  3. Plan: Edit using %s\n", cfg.AugmentStrategy)
			fmt.Printf("  4. Verify: Ensure compliance with %s rules\n", strings.Join(cfg.RulesEnabled, ", "))
			fmt.Print("\nApprove this approach? [y/n]: ")

			reader := bufio.NewReader(os.Stdin)
			choice, _ := reader.ReadString('\n')
			choice = strings.TrimSpace(choice)

			if choice != "y" && choice != "Y" {
				fmt.Println("Aborted.")
				return
			}

			// Store decision
			dbPath := filepath.Join(".virgil", "virgil.db")
			if err := storage.StoreAssessment(dbPath, results); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not store assessment: %v\n", err)
			}

			// Step 3 & 4 are in Phase 3
			fmt.Println("\nStep 3: GENERATE - Code modification (Phase 3)")
			fmt.Println("Step 4: VERIFY - Updated code verification (Phase 3)")
			fmt.Println("\n✓ Approach approved and stored")
			fmt.Println("→ Code generation coming in Phase 3")
		},
	}

	cmd.Flags().StringVar(&augment, "augment", "", "augmentation strategy (api|learning)")

	return cmd
}
func newAssessCommand() *cobra.Command {
	var augment string

	cmd := &cobra.Command{
		Use:   "assess [path]",
		Short: "Assess code against verification rules (utility)",
		Long: `Assess existing code against verification rules.

Use this utility command to review code independently of the create workflow.
Good for debugging, auditing, or checking code before committing.

If a path is provided, assesses that specific file or directory.
If no path is provided, assesses the entire project.

Rules checked:
  - OWASP Top 10 (security best practices - international, default)
  - NIST security guidelines (international, default)
  - GDPR compliance (optional - EU/EEA specific)
  - HIPAA compliance (optional - USA healthcare)
  - PCI-DSS compliance (optional - payment processing)
  - CIS Controls (optional)
  - ISO 27001 (optional - information security)
  - Custom rules (optional - user-defined)

Note: For the primary code creation workflow, use 'virgil create' instead.`,
		Example: `  virgil assess                                # Assess full project
  virgil assess pkg/auth/handler.go           # Assess specific file
  virgil assess internal/                     # Assess directory
  virgil assess --augment learning pkg/auth/  # Assess with learned patterns`,
		Run: func(cmd *cobra.Command, args []string) {
			// Get target path
			target := "."
			if len(args) > 0 {
				target = args[0]
			}

			// Load configuration
			cfg, err := config.LoadConfig(".")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
				os.Exit(1)
			}

			// Apply augment strategy override if provided via flag
			if augment != "" {
				cfg.AugmentStrategy = augment
			}

			// Load database connection
			dbPath := filepath.Join(".virgil", "virgil.db")
			db, err := storage.InitDatabase(dbPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			// Run verification pipeline
			fmt.Println("Assessing code...")
			results, err := verification.RunPipeline(target, target, cfg, db)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error running assessment: %v\n", err)
				os.Exit(1)
			}

			// Display results
			fmt.Println("\n=== Assessment Results ===\n")
			fmt.Printf("Target: %s\n", target)
			fmt.Printf("Total Issues: %d\n", len(results.AllIssues))

			if len(results.AllIssues) == 0 {
				fmt.Println("\n✓ No issues found!")
			} else {
				fmt.Println("\nIssues:")
				for i, issue := range results.AllIssues {
					fmt.Printf("\n%d. [%s] %s\n", i+1, issue.Severity, issue.Message)
					fmt.Printf("   Location: %s:%d\n", issue.FilePath, issue.LineNumber)
					fmt.Printf("   Rule: %s\n", issue.Rule)
					if issue.Suggestion != "" {
						fmt.Printf("   Suggestion: %s\n", issue.Suggestion)
					}
				}
			}

			// Store results to database
			if err := storage.StoreAssessment(dbPath, results); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Could not store assessment to database: %v\n", err)
			} else {
				fmt.Println("\n✓ Assessment stored to audit trail")
			}

			fmt.Printf("\nRun 'virgil review' to see assessment history\n")
		},
	}

	cmd.Flags().StringVar(&augment, "augment", "", "augmentation strategy (api|learning)")

	return cmd
}

// newReviewCommand creates the 'review' command
func newReviewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "review",
		Short: "View verification results",
		Long: `View results from recent assessments and code generation.

Displays:
  - Verification results from last assessment
  - Audit trail of actions taken
  - Summary of issues found
  - Previous assessments and their outcomes`,
		Example: `  virgil review`,
		Run: func(cmd *cobra.Command, args []string) {
			dbPath := filepath.Join(".virgil", "virgil.db")

			// Check if database exists
			if _, err := os.Stat(dbPath); os.IsNotExist(err) {
				fmt.Println("No assessments yet. Run 'virgil create' to start the workflow or 'virgil assess' to review existing code.")
				return
			}

			// Retrieve assessments from database
			assessments, err := storage.GetAssessments(dbPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error retrieving assessments: %v\n", err)
				os.Exit(1)
			}

			if len(assessments) == 0 {
				fmt.Println("No assessments found in audit trail.")
				return
			}

			fmt.Println("=== Verification Audit Trail ===\n")
			fmt.Printf("Total Assessments: %d\n\n", len(assessments))

			for i, assessment := range assessments {
				fmt.Printf("Assessment #%d\n", i+1)
				fmt.Printf("  Timestamp: %s\n", assessment.Timestamp)
				fmt.Printf("  Issues Found: %d\n", assessment.IssueCount)
				fmt.Printf("  Severity: %s\n", assessment.MaxSeverity)
				fmt.Println()
			}

			fmt.Println("Run 'virgil assess [path]' to perform a new assessment")
		},
	}
	return cmd
}

// newLearnCommand creates the 'learn' command for learning mode
func newLearnCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "learn <path>",
		Short: "Learn code patterns from existing codebase (learning mode)",
		Long: `Extract programming patterns from your codebase for use in learning mode.

Learning mode enables Virgil to generate code that matches YOUR engineering style,
without requiring external API calls. Learned patterns are stored encrypted locally.

Workflow:
  1. SCAN - Analyzes your codebase for programming patterns
  2. EXTRACT - Identifies structure, error handling, validation, security patterns
  3. STORE - Saves patterns encrypted in local database
  4. USE - Generated code will follow these patterns in learning mode

After learning, activate with: virgil config --augment learning`,
		Example: `  virgil learn .                     # Learn from current directory
  virgil learn ./src                 # Learn from specific directory
  virgil learn /path/to/project      # Learn from absolute path`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			codepath := args[0]

			fmt.Println("=== Code Pattern Learning ===\n")
			fmt.Printf("Analyzing: %s\n\n", codepath)

			// Load database connection
			dbPath := filepath.Join(".virgil", "virgil.db")
			db, err := storage.InitDatabase(dbPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing database: %v\n", err)
				os.Exit(1)
			}
			defer db.Close()

			// Initialize learner
			learner := learning.NewLearner(db)

			// Create learning request
			learnRequest := &learning.LearningRequest{
				CodebasePath: codepath,
				Languages:    make([]string, 0), // Auto-detect
			}

			// Learn patterns
			learnResponse, err := learner.Learn(learnRequest)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during learning: %v\n", err)
				os.Exit(1)
			}

			// Display results
			if !learnResponse.Success {
				fmt.Fprintf(os.Stderr, "Learning failed: %s\n", learnResponse.Error)
				os.Exit(1)
			}

			fmt.Printf("✓ %s\n\n", learnResponse.Message)

			// Display learned patterns
			if len(learnResponse.Patterns) > 0 {
				fmt.Println("Learned Patterns:")
				for i, pattern := range learnResponse.Patterns {
					fmt.Printf("  %d. [%s] %s (%s)\n", i+1, pattern.Type, pattern.Name, pattern.Language)
				}
			}

			fmt.Println("\n✓ Patterns saved to encrypted database")
			fmt.Println("→ Enable with: virgil config --augment learning")
			fmt.Println("→ Use with: virgil create \"your feature description\"")
		},
	}
}
