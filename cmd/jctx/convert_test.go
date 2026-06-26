package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunConvertCmd_ConvertsCLAUDEmd(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("Always write tests.\n"), 0644)

	var buf bytes.Buffer
	if err := runConvertCmd(root, "claude", false, false, &buf); err != nil {
		t.Fatalf("runConvertCmd: %v", err)
	}

	// Should have written .jctx/rules/CLAUDE.md
	outPath := filepath.Join(root, ".jctx", "rules", "CLAUDE.md")
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatal(".jctx/rules/CLAUDE.md not written")
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Always write tests.") {
		t.Errorf("converted content missing original:\n%s", string(data))
	}
	if !strings.Contains(string(data), "targets:") {
		t.Errorf("converted file should have targets frontmatter:\n%s", string(data))
	}
}

func TestRunConvertCmd_DryRun(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("Rules here.\n"), 0644)

	var buf bytes.Buffer
	if err := runConvertCmd(root, "claude", true, false, &buf); err != nil {
		t.Fatalf("runConvertCmd dry-run: %v", err)
	}

	// No file should be written
	outPath := filepath.Join(root, ".jctx", "rules", "CLAUDE.md")
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Error("dry run should not write .jctx/rules/CLAUDE.md")
	}
	if !strings.Contains(buf.String(), "(dry run)") {
		t.Errorf("expected '(dry run)' in output:\n%s", buf.String())
	}
}
