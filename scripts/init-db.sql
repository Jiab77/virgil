-- Copyright 2026 jiab77
-- SPDX-License-Identifier: MIT

-- Virgil Audit Trail Schema
-- Stores all verification assessments and decisions

CREATE TABLE IF NOT EXISTS assessments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id TEXT UNIQUE NOT NULL,
    target_path TEXT NOT NULL,
    rules_enabled TEXT NOT NULL, -- JSON array of rules checked
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT NOT NULL, -- PASS, FAIL, WARNING
    summary TEXT -- Short summary of findings
);

CREATE TABLE IF NOT EXISTS issues (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id TEXT NOT NULL,
    issue_id TEXT UNIQUE NOT NULL,
    rule TEXT NOT NULL, -- Which rule found this (e.g., "owasp:hardcoded-secrets")
    severity TEXT NOT NULL, -- CRITICAL, HIGH, MEDIUM, LOW
    file_path TEXT NOT NULL,
    line_number INTEGER,
    column_number INTEGER,
    message TEXT NOT NULL,
    suggestion TEXT, -- How to fix it
    code_snippet TEXT, -- The problematic code
    FOREIGN KEY (assessment_id) REFERENCES assessments(assessment_id)
);

CREATE TABLE IF NOT EXISTS decisions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    assessment_id TEXT NOT NULL,
    decision TEXT NOT NULL, -- APPROVED, REJECTED, PENDING
    reason TEXT, -- User's reason for decision
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (assessment_id) REFERENCES assessments(assessment_id)
);

CREATE INDEX IF NOT EXISTS idx_assessment_timestamp ON assessments(timestamp);
CREATE INDEX IF NOT EXISTS idx_assessment_status ON assessments(status);
CREATE INDEX IF NOT EXISTS idx_issue_assessment ON issues(assessment_id);
CREATE INDEX IF NOT EXISTS idx_issue_severity ON issues(severity);

CREATE TABLE IF NOT EXISTS web_search_cache (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    query TEXT NOT NULL,
    query_hash TEXT UNIQUE NOT NULL, -- Hash for deduplication
    search_context TEXT NOT NULL, -- Context of request (e.g., "create", "edit", "verify")
    results TEXT NOT NULL, -- JSON array of search results with URL, title, snippet, source
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME, -- TTL for cache invalidation
    created_by TEXT -- Which assessment/feature requested the search
);

CREATE INDEX IF NOT EXISTS idx_web_search_query_hash ON web_search_cache(query_hash);
CREATE INDEX IF NOT EXISTS idx_web_search_timestamp ON web_search_cache(timestamp);
CREATE INDEX IF NOT EXISTS idx_web_search_context ON web_search_cache(search_context);

CREATE TABLE IF NOT EXISTS learned_patterns (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pattern_id TEXT UNIQUE NOT NULL,
    pattern_type TEXT NOT NULL, -- 'structure', 'error', 'security', 'validation', 'logging', 'naming'
    language TEXT NOT NULL, -- 'go', 'python', 'javascript', etc.
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    file_path TEXT, -- Path to the source file where pattern was detected
    example TEXT, -- Encrypted example code
    frequency INTEGER DEFAULT 1, -- How often this pattern appears
    metadata TEXT, -- JSON: additional pattern metadata
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_learned_patterns_language ON learned_patterns(language);
CREATE INDEX IF NOT EXISTS idx_learned_patterns_type ON learned_patterns(pattern_type);
CREATE INDEX IF NOT EXISTS idx_learned_patterns_frequency ON learned_patterns(frequency);
