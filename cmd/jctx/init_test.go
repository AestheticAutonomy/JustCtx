package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInit_NewScaffold(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	err := runInit(tmpDir, &buf, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify standard directories were created
	subdirs := []string{"rules", "hooks", "mcp", "commands", "skills", "ignores"}
	for _, sub := range subdirs {
		path := filepath.Join(tmpDir, ".jctx", sub)
		fi, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected subdirectory %s to exist, error: %v", sub, err)
		} else if !fi.IsDir() {
			t.Errorf("expected %s to be a directory", sub)
		}
	}

	// Verify starter rule file
	ruleFile := filepath.Join(tmpDir, ".jctx", "rules", "main.md")
	contentBytes, err := os.ReadFile(ruleFile)
	if err != nil {
		t.Fatalf("expected starter rule file to exist, error: %v", err)
	}

	expectedContent := `---
target: [claude]
---
@@@ General Rules

Add your coding guidelines here.
`
	if string(contentBytes) != expectedContent {
		t.Errorf("starter rule file content mismatch.\nExpected:\n%s\nGot:\n%s", expectedContent, string(contentBytes))
	}

	// Verify stdout output format
	expectedOutput := strings.Join([]string{
		".jctx/",
		".jctx/rules/",
		".jctx/hooks/",
		".jctx/mcp/",
		".jctx/commands/",
		".jctx/skills/",
		".jctx/ignores/",
		".jctx/rules/main.md",
	}, "\n") + "\n"

	// Replace Windows line endings in buffer for cross-platform comparison
	actualOutput := strings.ReplaceAll(buf.String(), "\r\n", "\n")
	if actualOutput != expectedOutput {
		t.Errorf("stdout output mismatch.\nExpected:\n%q\nGot:\n%q", expectedOutput, actualOutput)
	}
}

func TestRunInit_Exists(t *testing.T) {
	tmpDir := t.TempDir()

	// Pre-create .jctx
	jctxDir := filepath.Join(tmpDir, ".jctx")
	if err := os.Mkdir(jctxDir, 0755); err != nil {
		t.Fatalf("failed to create pre-existing .jctx: %v", err)
	}

	var buf bytes.Buffer
	err := runInit(tmpDir, &buf, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify it printed the correct message
	expectedMessage := ".jctx/ already exists\n"
	actualOutput := strings.ReplaceAll(buf.String(), "\r\n", "\n")
	if actualOutput != expectedMessage {
		t.Errorf("expected output %q, got %q", expectedMessage, actualOutput)
	}

	// Verify no subfolders or files were created
	subdirs := []string{"rules", "hooks", "mcp", "commands", "skills", "ignores"}
	for _, sub := range subdirs {
		path := filepath.Join(jctxDir, sub)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Errorf("subdirectory %s should not have been created", sub)
		}
	}
}

func TestRunInit_JSON(t *testing.T) {
	tmpDir := t.TempDir()

	var buf bytes.Buffer
	err := runInit(tmpDir, &buf, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var paths []string
	if err := json.Unmarshal(buf.Bytes(), &paths); err != nil {
		t.Fatalf("failed to unmarshal JSON output: %v, raw output: %s", err, buf.String())
	}

	expectedPaths := []string{
		".jctx/",
		".jctx/rules/",
		".jctx/hooks/",
		".jctx/mcp/",
		".jctx/commands/",
		".jctx/skills/",
		".jctx/ignores/",
		".jctx/rules/main.md",
	}

	if len(paths) != len(expectedPaths) {
		t.Fatalf("expected %d paths in JSON, got %d", len(expectedPaths), len(paths))
	}

	for i, p := range paths {
		if p != expectedPaths[i] {
			t.Errorf("at index %d: expected path %q, got %q", i, expectedPaths[i], p)
		}
	}
}
