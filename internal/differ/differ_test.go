package differ

import (
	"os"
	"path/filepath"
	"testing"

	_ "github.com/AestheticAutonomy/justctx/internal/providers/agents"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/antigravity"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/claude"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/cursor"
)

func buildTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		abs := filepath.Join(root, rel)
		if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(abs, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestDiff_Clean(t *testing.T) {
	content := "## Core Rules\n\nAlways write tests.\n"
	root := buildTree(t, map[string]string{
		".jctx/rules/coding.md": "---\ntargets: [claude]\n---\n\n@@@ Core Rules\nAlways write tests.\n",
		"CLAUDE.md":             content,
	})

	res, err := Diff(DiffOpts{Root: root, Target: "claude"})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if !res.InSync {
		t.Errorf("expected in_sync=true, got changes: %+v", res.Changes)
	}
}

func TestDiff_Modified(t *testing.T) {
	root := buildTree(t, map[string]string{
		".jctx/rules/coding.md": "---\ntargets: [claude]\n---\n\n@@@ Core Rules\nAlways write tests.\n",
		"CLAUDE.md":             "some different content\n",
	})

	res, err := Diff(DiffOpts{Root: root, Target: "claude"})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if res.InSync {
		t.Error("expected in_sync=false for modified content")
	}
	if len(res.Changes) != 1 {
		t.Errorf("expected 1 change, got %d", len(res.Changes))
	}
}

func TestDiff_MissingFile(t *testing.T) {
	root := buildTree(t, map[string]string{
		".jctx/rules/coding.md": "---\ntargets: [claude]\n---\n\n@@@ Core Rules\nAlways write tests.\n",
		// No CLAUDE.md on disk
	})

	res, err := Diff(DiffOpts{Root: root, Target: "claude"})
	if err != nil {
		t.Fatalf("Diff: %v", err)
	}
	if res.InSync {
		t.Error("expected in_sync=false when file doesn't exist")
	}
	if len(res.Changes) != 1 || res.Changes[0].Type != "added" {
		t.Errorf("expected change type 'added', got: %+v", res.Changes)
	}
}
