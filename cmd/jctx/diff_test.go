package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func buildDiffProject(t *testing.T, files map[string]string) string {
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

func TestRunDiff_InSync(t *testing.T) {
	root := buildDiffProject(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nWrite tests.\n",
	})

	// First gen to produce the output file
	runGen(root, []string{"claude"}, "", nil, false, false, &bytes.Buffer{})

	var buf bytes.Buffer
	res, err := runDiff(root, "claude", "", nil, false, &buf)
	if err != nil {
		t.Fatalf("runDiff: %v", err)
	}
	if !res.InSync {
		t.Errorf("expected InSync=true:\n%s", buf.String())
	}
}

func TestRunDiff_Modified(t *testing.T) {
	root := buildDiffProject(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nOriginal.\n",
	})

	// Gen to produce output
	runGen(root, []string{"claude"}, "", nil, false, false, &bytes.Buffer{})

	// Modify the source after gen
	os.WriteFile(filepath.Join(root, ".jctx", "rules", "rules.md"),
		[]byte("---\ntargets: [claude]\n---\n\n@@@ Rules\nModified.\n"), 0644)

	var buf bytes.Buffer
	res, err := runDiff(root, "claude", "", nil, false, &buf)
	if err != nil {
		t.Fatalf("runDiff: %v", err)
	}
	if res.InSync {
		t.Error("expected InSync=false after source change")
	}
}

func TestRunDiff_JSON(t *testing.T) {
	root := buildDiffProject(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nContent.\n",
	})

	runGen(root, []string{"claude"}, "", nil, false, false, &bytes.Buffer{})

	var buf bytes.Buffer
	_, err := runDiff(root, "claude", "", nil, true, &buf)
	if err != nil {
		t.Fatalf("runDiff JSON: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\nraw: %s", err, buf.String())
	}
	if _, ok := result["in_sync"]; !ok {
		t.Errorf("expected in_sync in JSON:\n%s", buf.String())
	}
}
