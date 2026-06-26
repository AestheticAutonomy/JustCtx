package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func buildGenProject(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	os.MkdirAll(filepath.Join(root, ".git"), 0755)
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		os.MkdirAll(filepath.Dir(abs), 0755)
		os.WriteFile(abs, []byte(content), 0644)
	}
	return root
}

func TestRunGen_WritesCLAUDEmd(t *testing.T) {
	root := buildGenProject(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nWrite clean code.\n",
	})

	var buf bytes.Buffer
	err := runGen(root, []string{"claude"}, "", nil, false, false, &buf)
	if err != nil {
		t.Fatalf("runGen: %v", err)
	}

	outPath := filepath.Join(root, "CLAUDE.md")
	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Fatal("CLAUDE.md not written")
	}
	if !strings.Contains(buf.String(), "CLAUDE.md") {
		t.Errorf("expected output path in stdout:\n%s", buf.String())
	}
}

func TestRunGen_DryRun(t *testing.T) {
	root := buildGenProject(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nContent.\n",
	})

	var buf bytes.Buffer
	err := runGen(root, []string{"claude"}, "", nil, true, false, &buf)
	if err != nil {
		t.Fatalf("runGen dry-run: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "CLAUDE.md")); !os.IsNotExist(err) {
		t.Error("dry run should not write CLAUDE.md")
	}
	if !strings.Contains(buf.String(), "(dry run)") {
		t.Errorf("expected '(dry run)' in output:\n%s", buf.String())
	}
}

func TestRunGen_JSON(t *testing.T) {
	root := buildGenProject(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nContent.\n",
	})

	var buf bytes.Buffer
	err := runGen(root, []string{"claude"}, "", nil, true, true, &buf)
	if err != nil {
		t.Fatalf("runGen JSON: %v", err)
	}

	// Should be valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\nraw: %s", err, buf.String())
	}
	if _, ok := result["output_path"]; !ok {
		t.Errorf("expected output_path in JSON:\n%s", buf.String())
	}
}

func TestRunGen_AllProviders(t *testing.T) {
	root := buildGenProject(t, map[string]string{
		".jctx/rules/claude.md":      "---\ntargets: [claude]\n---\n\n@@@ Rules\nClaude rules.\n",
		".jctx/rules/agents.md":      "---\ntargets: [agents]\n---\n\n@@@ Rules\nAgents rules.\n",
		".jctx/rules/antigravity.md": "---\ntargets: [antigravity]\n---\n\n@@@ Rules\nGemini rules.\n",
	})

	var buf bytes.Buffer
	err := runGen(root, allTargets(), "", nil, false, false, &buf)
	if err != nil {
		t.Fatalf("runGen all: %v", err)
	}

	for _, f := range []string{"CLAUDE.md", "AGENTS.md", "GEMINI.md"} {
		if _, err := os.Stat(filepath.Join(root, f)); os.IsNotExist(err) {
			t.Errorf("%s not written", f)
		}
	}
}
