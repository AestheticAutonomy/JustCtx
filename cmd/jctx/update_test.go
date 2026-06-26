package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildUpdateProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	os.MkdirAll(filepath.Join(root, ".jctx", "rules"), 0755)
	os.WriteFile(filepath.Join(root, ".jctx", "rules", "rules.md"),
		[]byte("---\ntargets: [claude]\n---\n\n@@@ Rules\nContent.\n"), 0644)
	os.WriteFile(filepath.Join(root, ".jctx", "config.json"),
		[]byte(`{"schema_version":1,"defaults":{"target":"claude","role":"","tags":[]}}`), 0644)
	return root
}

func TestRunUpdate_RegeneratesFromConfig(t *testing.T) {
	root := buildUpdateProject(t)

	var buf bytes.Buffer
	if err := runUpdate(root, false, false, &buf); err != nil {
		t.Fatalf("runUpdate: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "CLAUDE.md")); os.IsNotExist(err) {
		t.Fatal("CLAUDE.md not written by update")
	}
	if !strings.Contains(buf.String(), "CLAUDE.md") {
		t.Errorf("expected output path in stdout:\n%s", buf.String())
	}
}

func TestRunUpdate_NoConfig(t *testing.T) {
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	os.MkdirAll(filepath.Join(root, ".jctx"), 0755)

	var buf bytes.Buffer
	err := runUpdate(root, false, false, &buf)
	if err == nil {
		t.Fatal("expected error when no config.json")
	}
	if !strings.Contains(err.Error(), "config.json") {
		t.Errorf("expected config.json mention in error: %v", err)
	}
}
