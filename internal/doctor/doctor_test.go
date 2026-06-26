package doctor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AestheticAutonomy/justctx/internal/manifest"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/agents"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/antigravity"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/claude"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/cursor"
)

func TestDoctor_ValidProject(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jctx", "rules"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0755); err != nil {
		t.Fatal(err)
	}

	res, err := Run(root, "")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if !res.AllPass {
		for _, c := range res.Checks {
			if !c.Pass {
				t.Errorf("check failed: %s — %s", c.Name, c.Detail)
			}
		}
	}
}

func TestDoctor_MissingJctx(t *testing.T) {
	root := t.TempDir()
	// No .jctx/ directory

	res, err := Run(root, "")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	found := false
	for _, c := range res.Checks {
		if c.Name == ".jctx/ exists" && !c.Pass {
			found = true
		}
	}
	if !found {
		t.Error("expected .jctx/ check to fail when directory missing")
	}
	if res.AllPass {
		t.Error("expected AllPass=false")
	}
}

func TestDoctor_StaleManifest(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".jctx", "rules"), 0755); err != nil {
		t.Fatal(err)
	}

	// Write manifest pointing to a file that doesn't exist
	m := &manifest.Manifest{
		SchemaVersion: 1,
		Target:        "claude",
		OutputPath:    "CLAUDE.md",
	}
	if err := manifest.Write(root, "CLAUDE.md", m); err != nil {
		t.Fatal(err)
	}

	res, err := Run(root, "")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}

	found := false
	for _, c := range res.Checks {
		if c.Name == "manifest target exists: CLAUDE.md" && !c.Pass {
			found = true
		}
	}
	if !found {
		t.Error("expected stale manifest check to fail")
	}
}
