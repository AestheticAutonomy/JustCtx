package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildCleanProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	// Write source and gen to produce CLAUDE.md + manifest
	os.MkdirAll(filepath.Join(root, ".jctx", "rules"), 0755)
	os.WriteFile(filepath.Join(root, ".jctx", "rules", "rules.md"),
		[]byte("---\ntargets: [claude]\n---\n\n@@@ Rules\nContent.\n"), 0644)
	runGen(root, []string{"claude"}, "", nil, false, false, &bytes.Buffer{})
	return root
}

func TestRunClean_RemovesFileAndManifest(t *testing.T) {
	root := buildCleanProject(t)

	outPath := filepath.Join(root, "CLAUDE.md")
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatal("CLAUDE.md should exist before clean")
	}

	var buf bytes.Buffer
	if err := runClean(root, "claude", false, false, false, &buf); err != nil {
		t.Fatalf("runClean: %v", err)
	}

	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Error("CLAUDE.md should be removed after clean")
	}

	manifestPath := filepath.Join(root, ".jctx", ".manifest", "CLAUDE.md.json")
	if _, err := os.Stat(manifestPath); !os.IsNotExist(err) {
		t.Error("manifest should be removed after clean")
	}
}

func TestRunClean_NoManifest(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	os.MkdirAll(filepath.Join(root, ".jctx"), 0755)

	var buf bytes.Buffer
	if err := runClean(root, "claude", false, false, false, &buf); err != nil {
		t.Fatalf("runClean with no manifest: %v", err)
	}
	if !strings.Contains(buf.String(), "nothing to clean") {
		t.Errorf("expected 'nothing to clean' message:\n%s", buf.String())
	}
}

func TestRunClean_DryRun(t *testing.T) {
	root := buildCleanProject(t)

	outPath := filepath.Join(root, "CLAUDE.md")
	var buf bytes.Buffer
	if err := runClean(root, "claude", false, true, false, &buf); err != nil {
		t.Fatalf("runClean dry-run: %v", err)
	}

	// File should still exist
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("dry run should not remove CLAUDE.md")
	}
	if !strings.Contains(buf.String(), "(dry run)") {
		t.Errorf("expected '(dry run)' in output:\n%s", buf.String())
	}
}
