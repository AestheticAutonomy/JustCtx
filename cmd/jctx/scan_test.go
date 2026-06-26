package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

func buildScanProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	// Isolate global dir
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return root
}

func TestRunScan_ClaudeTarget(t *testing.T) {
	root := buildScanProject(t)
	os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("Always write tests.\n"), 0644)

	var buf bytes.Buffer
	err := runScan(root, "claude", true, false, 0, false, &buf)
	if err != nil {
		t.Fatalf("runScan: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Sources") {
		t.Errorf("expected Sources in output:\n%s", output)
	}
	if !strings.Contains(output, "Always write tests.") {
		t.Errorf("expected CLAUDE.md content in output:\n%s", output)
	}
}

func TestRunScan_JSON(t *testing.T) {
	root := buildScanProject(t)
	os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("Test rule.\n"), 0644)

	var buf bytes.Buffer
	err := runScan(root, "claude", true, false, 0, true, &buf)
	if err != nil {
		t.Fatalf("runScan JSON: %v", err)
	}

	var res schema.ScanResult
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		t.Fatalf("unmarshal JSON: %v\nraw: %s", err, buf.String())
	}
	if len(res.Sources) == 0 {
		t.Error("expected at least one source in JSON output")
	}
}

func TestRunScan_UnknownTarget(t *testing.T) {
	root := buildScanProject(t)
	var buf bytes.Buffer
	err := runScan(root, "nonexistent-provider", true, false, 0, false, &buf)
	if err == nil {
		t.Fatal("expected error for unknown target")
	}
}

func TestRunScan_BottomUpDepth(t *testing.T) {
	root := buildScanProject(t)

	// Create a simple hierarchy
	sub := filepath.Join(root, "sub")
	os.MkdirAll(sub, 0755)
	os.WriteFile(filepath.Join(root, "CLAUDE.md"), []byte("root rule\n"), 0644)
	os.WriteFile(filepath.Join(sub, "CLAUDE.md"), []byte("sub rule\n"), 0644)

	var buf bytes.Buffer
	// bottom-up depth=1 from sub: should only include sub level (depth 0) and root (depth 1)
	err := runScan(sub, "claude", true, true, 1, false, &buf)
	if err != nil {
		t.Fatalf("runScan bottom-up depth: %v", err)
	}
	output := buf.String()
	if !strings.Contains(output, "Sources") {
		t.Errorf("expected Sources in output:\n%s", output)
	}
}
