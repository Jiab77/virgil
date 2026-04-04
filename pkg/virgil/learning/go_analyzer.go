// Copyright 2026 jiab77
// SPDX-License-Identifier: MIT

package learning

import (
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// AnalyzeCodebase analyzes Go code and extracts patterns
func (ga *GoAnalyzer) AnalyzeCodebase(codebasePath string) ([]CodePattern, error) {
	patterns := make([]CodePattern, 0)

	// Walk through directory and parse Go files
	err := filepath.Walk(codebasePath, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}

		if strings.Contains(path, "vendor") || strings.Contains(path, ".test.go") {
			return nil
		}

		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
		if err != nil {
			log.Printf("[virgil] Warning: could not parse %s: %v", path, err)
			return nil
		}

		ga.analyzeFile(file, path, &patterns)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return deduplicatePatterns(patterns), nil
}

// analyzeFile extracts patterns from a single Go file
func (ga *GoAnalyzer) analyzeFile(file *ast.File, filepath string, patterns *[]CodePattern) {
	// Extract package structure pattern
	ga.extractPackageStructure(file, filepath, patterns)

	// Extract error handling patterns
	ga.extractErrorHandlingPatterns(file, patterns)

	// Extract validation patterns
	ga.extractValidationPatterns(file, patterns)

	// Extract logging patterns
	ga.extractLoggingPatterns(file, patterns)

	// Extract security patterns
	ga.extractSecurityPatterns(file, patterns)

	// Extract naming conventions
	ga.extractNamingConventions(file, patterns)
}

// extractPackageStructure identifies package organization
func (ga *GoAnalyzer) extractPackageStructure(file *ast.File, filepath string, patterns *[]CodePattern) {
	pkg := file.Name.Name

	*patterns = append(*patterns, CodePattern{
		Type:        PatternTypeStructure,
		Name:        "Package: " + pkg,
		Description: "Package-based organization using '" + pkg + "'",
		Language:    "go",
		Frequency:   1,
	})
}

// extractErrorHandlingPatterns identifies how errors are handled
func (ga *GoAnalyzer) extractErrorHandlingPatterns(file *ast.File, patterns *[]CodePattern) {
	patternMap := make(map[string]int)

	// Look for error returns and error checks
	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		// Check if function returns error
		if funcDecl.Type.Results != nil {
			for _, field := range funcDecl.Type.Results.List {
				if ident, ok := field.Type.(*ast.Ident); ok && ident.Name == "error" {
					patternMap["returns error"]++
				}
			}
		}

		// Walk function body for error checks
		if funcDecl.Body != nil {
			ast.Walk(&errorCheckVisitor{patterns: patternMap}, funcDecl.Body)
		}
	}

	// Add detected patterns
	for pattern, freq := range patternMap {
		*patterns = append(*patterns, CodePattern{
			Type:        PatternTypeErrorHandling,
			Name:        "Error handling: " + pattern,
			Description: "Error handling pattern detected: " + pattern,
			Language:    "go",
			Frequency:   freq,
		})
	}
}

// extractValidationPatterns identifies input validation approaches
func (ga *GoAnalyzer) extractValidationPatterns(file *ast.File, patterns *[]CodePattern) {
	// Look for common validation patterns
	patternMap := make(map[string]int)

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Body != nil {
			ast.Walk(&validationVisitor{patterns: patternMap}, funcDecl.Body)
		}
	}

	for pattern, freq := range patternMap {
		*patterns = append(*patterns, CodePattern{
			Type:        PatternTypeValidation,
			Name:        "Validation: " + pattern,
			Description: "Input validation pattern detected: " + pattern,
			Language:    "go",
			Frequency:   freq,
		})
	}
}

// extractLoggingPatterns identifies logging strategies
func (ga *GoAnalyzer) extractLoggingPatterns(file *ast.File, patterns *[]CodePattern) {
	patternMap := make(map[string]int)

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Body != nil {
			ast.Walk(&loggingVisitor{patterns: patternMap}, funcDecl.Body)
		}
	}

	for pattern, freq := range patternMap {
		*patterns = append(*patterns, CodePattern{
			Type:        PatternTypeLogging,
			Name:        "Logging: " + pattern,
			Description: "Logging pattern detected: " + pattern,
			Language:    "go",
			Frequency:   freq,
		})
	}
}

// extractSecurityPatterns identifies security-related practices
func (ga *GoAnalyzer) extractSecurityPatterns(file *ast.File, patterns *[]CodePattern) {
	patternMap := make(map[string]int)

	for _, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		if funcDecl.Body != nil {
			ast.Walk(&securityVisitor{patterns: patternMap}, funcDecl.Body)
		}
	}

	for pattern, freq := range patternMap {
		*patterns = append(*patterns, CodePattern{
			Type:        PatternTypeSecurity,
			Name:        "Security: " + pattern,
			Description: "Security pattern detected: " + pattern,
			Language:    "go",
			Frequency:   freq,
		})
	}
}

// extractNamingConventions identifies naming patterns
func (ga *GoAnalyzer) extractNamingConventions(file *ast.File, patterns *[]CodePattern) {
	patternMap := make(map[string]int)

	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if strings.HasPrefix(d.Name.Name, "Test") {
				patternMap["Test prefix for tests"]++
			}
			if strings.HasPrefix(d.Name.Name, "Benchmark") {
				patternMap["Benchmark prefix for benchmarks"]++
			}
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				// Analyze type naming
				patternMap["Type definitions"]++
			}
			if d.Tok == token.CONST {
				// Analyze constant naming
				patternMap["Constant definitions"]++
			}
		}
	}

	for pattern, freq := range patternMap {
		*patterns = append(*patterns, CodePattern{
			Type:        PatternTypeNaming,
			Name:        "Naming: " + pattern,
			Description: "Naming convention detected: " + pattern,
			Language:    "go",
			Frequency:   freq,
		})
	}
}

// AST Visitors for pattern detection

// errorCheckVisitor detects error checking patterns
type errorCheckVisitor struct {
	patterns map[string]int
}

func (v *errorCheckVisitor) Visit(node ast.Node) ast.Visitor {
	if ifStmt, ok := node.(*ast.IfStmt); ok {
		// Check for "if err != nil" pattern
		if binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
			if ident, ok := binExpr.X.(*ast.Ident); ok && ident.Name == "err" {
				v.patterns["if err != nil"]++
			}
		}
	}
	return v
}

// validationVisitor detects validation patterns
type validationVisitor struct {
	patterns map[string]int
}

func (v *validationVisitor) Visit(node ast.Node) ast.Visitor {
	if ifStmt, ok := node.(*ast.IfStmt); ok {
		// Check for common validation patterns
		if binExpr, ok := ifStmt.Cond.(*ast.BinaryExpr); ok {
			if _, ok := binExpr.X.(*ast.Ident); ok {
				// Detect pattern like "if len(x) == 0" or "if x == nil"
				v.patterns["conditional validation"]++
			}
		}
	}
	return v
}

// loggingVisitor detects logging patterns
type loggingVisitor struct {
	patterns map[string]int
}

func (v *loggingVisitor) Visit(node ast.Node) ast.Visitor {
	if callExpr, ok := node.(*ast.CallExpr); ok {
		// Check for log.Printf, log.Println, etc.
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selExpr.X.(*ast.Ident); ok && ident.Name == "log" {
				v.patterns[selExpr.Sel.Name]++
			}
		}
	}
	return v
}

// securityVisitor detects security patterns
type securityVisitor struct {
	patterns map[string]int
}

func (v *securityVisitor) Visit(node ast.Node) ast.Visitor {
	if callExpr, ok := node.(*ast.CallExpr); ok {
		// Check for crypto usage
		if selExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			if ident, ok := selExpr.X.(*ast.Ident); ok {
				if strings.Contains(ident.Name, "crypto") {
					v.patterns["cryptographic operation"]++
				}
				if strings.Contains(ident.Name, "hash") {
					v.patterns["hashing operation"]++
				}
			}
		}
	}
	return v
}
