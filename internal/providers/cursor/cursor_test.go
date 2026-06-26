package cursor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

func TestCursorProvider_FindFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Mock repo root
	repoRoot := filepath.Join(tmpDir, "myrepo")
	err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create .cursor/rules directory
	rulesDir := filepath.Join(repoRoot, ".cursor", "rules")
	err = os.MkdirAll(rulesDir, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Write mock legacy file
	legacyPath := filepath.Join(repoRoot, ".cursorrules")
	err = os.WriteFile(legacyPath, []byte("legacy rules"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Write mock rules files
	mdc1 := filepath.Join(rulesDir, "rule1.mdc")
	err = os.WriteFile(mdc1, []byte("rule 1"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mdc2 := filepath.Join(rulesDir, "rule2.mdc")
	err = os.WriteFile(mdc2, []byte("rule 2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Write ignored files
	err = os.WriteFile(filepath.Join(rulesDir, "ignore.md"), []byte("ignored"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	provider := &CursorProvider{}
	files, err := provider.FindFiles(repoRoot, schema.TypeRules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{legacyPath, mdc1, mdc2}
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(files), files)
	}

	for i, f := range files {
		expectedPath := filepath.Clean(expected[i])
		actualPath := filepath.Clean(f)
		if expectedPath != actualPath {
			t.Errorf("at index %d: expected %s, got %s", i, expectedPath, actualPath)
		}
	}
}

func TestCursorProvider_ParseRules(t *testing.T) {
	provider := &CursorProvider{}

	// Test .cursorrules (no frontmatter)
	tmpDir := t.TempDir()
	legacyFile := filepath.Join(tmpDir, ".cursorrules")
	err := os.WriteFile(legacyFile, []byte("some rules\nline 2"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	sec1, err := provider.ParseRules(legacyFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(sec1) != 1 || sec1[0].Content != "some rules\nline 2" {
		t.Errorf("unexpected legacy content: %v", sec1)
	}

	// Test .mdc with frontmatter
	mdcFile := filepath.Join(tmpDir, "rule.mdc")
	mdcContent := `---
description: Test rule
globs: *.go
---
Actual rule content here
and line 2`

	err = os.WriteFile(mdcFile, []byte(mdcContent), 0644)
	if err != nil {
		t.Fatal(err)
	}

	sec2, err := provider.ParseRules(mdcFile)
	if err != nil {
		t.Fatal(err)
	}
	if len(sec2) != 1 {
		t.Fatalf("expected 1 section, got %d", len(sec2))
	}

	expectedContent := "Actual rule content here\nand line 2"
	if strings.TrimSpace(sec2[0].Content) != expectedContent {
		t.Errorf("expected stripped content '%s', got '%s'", expectedContent, sec2[0].Content)
	}
}
