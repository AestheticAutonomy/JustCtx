package manifest

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListManifests(t *testing.T) {
	root := t.TempDir()

	m1 := &Manifest{SchemaVersion: 1, Target: "claude", OutputPath: "CLAUDE.md"}
	m2 := &Manifest{SchemaVersion: 1, Target: "antigravity", OutputPath: "GEMINI.md"}
	if err := Write(root, "CLAUDE.md", m1); err != nil {
		t.Fatal(err)
	}
	if err := Write(root, "GEMINI.md", m2); err != nil {
		t.Fatal(err)
	}

	manifests, err := ListManifests(root)
	if err != nil {
		t.Fatalf("ListManifests: %v", err)
	}
	if len(manifests) != 2 {
		t.Fatalf("expected 2 manifests, got %d", len(manifests))
	}
}

func TestListManifests_Empty(t *testing.T) {
	root := t.TempDir()
	manifests, err := ListManifests(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(manifests) != 0 {
		t.Errorf("expected empty, got %d", len(manifests))
	}
}

func TestDelete(t *testing.T) {
	root := t.TempDir()

	m := &Manifest{SchemaVersion: 1, Target: "claude", OutputPath: "CLAUDE.md"}
	if err := Write(root, "CLAUDE.md", m); err != nil {
		t.Fatal(err)
	}

	// Verify file exists
	sidecar := filepath.Join(root, ".jctx", ".manifest", "CLAUDE.md.json")
	if _, err := os.Stat(sidecar); err != nil {
		t.Fatalf("expected sidecar to exist: %v", err)
	}

	if err := Delete(root, "CLAUDE.md"); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Should be gone
	if _, err := os.Stat(sidecar); !os.IsNotExist(err) {
		t.Error("expected sidecar to be deleted")
	}

	// Second delete should be no-op
	if err := Delete(root, "CLAUDE.md"); err != nil {
		t.Errorf("second Delete should be no-op: %v", err)
	}
}
