package manifest

import (
	"testing"
	"time"
)

func TestWriteRead_Roundtrip(t *testing.T) {
	root := t.TempDir()

	m := &Manifest{
		SchemaVersion: 1,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Target:        "claude",
		OutputPath:    "CLAUDE.md",
		Chunks: []Chunk{
			{
				SourceFile:         ".jctx/rules/coding.md",
				Section:            "Core Rules",
				AssembledLineStart: 1,
				AssembledLineEnd:   10,
				SourceLineStart:    1,
				SourceLineEnd:      10,
			},
		},
	}

	if err := Write(root, "CLAUDE.md", m); err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := Read(root, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if got == nil {
		t.Fatal("expected manifest, got nil")
	}
	if got.SchemaVersion != m.SchemaVersion {
		t.Errorf("SchemaVersion: want %d, got %d", m.SchemaVersion, got.SchemaVersion)
	}
	if got.Target != m.Target {
		t.Errorf("Target: want %s, got %s", m.Target, got.Target)
	}
	if got.OutputPath != m.OutputPath {
		t.Errorf("OutputPath: want %s, got %s", m.OutputPath, got.OutputPath)
	}
	if len(got.Chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(got.Chunks))
	}
	if got.Chunks[0].Section != "Core Rules" {
		t.Errorf("chunk section: want Core Rules, got %s", got.Chunks[0].Section)
	}
}

func TestRead_Missing(t *testing.T) {
	root := t.TempDir()
	m, err := Read(root, "CLAUDE.md")
	if err != nil {
		t.Fatalf("expected nil error on missing file, got: %v", err)
	}
	if m != nil {
		t.Errorf("expected nil manifest, got: %+v", m)
	}
}

func TestManifest_UpdateOnRegen(t *testing.T) {
	root := t.TempDir()

	m1 := &Manifest{
		SchemaVersion: 1,
		Target:        "claude",
		OutputPath:    "CLAUDE.md",
		Chunks:        []Chunk{{SourceFile: "rules.md", Section: "Original Section"}},
	}
	if err := Write(root, "CLAUDE.md", m1); err != nil {
		t.Fatal(err)
	}

	// Overwrite with updated manifest
	m2 := &Manifest{
		SchemaVersion: 1,
		Target:        "claude",
		OutputPath:    "CLAUDE.md",
		Chunks:        []Chunk{{SourceFile: "rules.md", Section: "Updated Section"}},
	}
	if err := Write(root, "CLAUDE.md", m2); err != nil {
		t.Fatal(err)
	}

	got, err := Read(root, "CLAUDE.md")
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if len(got.Chunks) != 1 || got.Chunks[0].Section != "Updated Section" {
		t.Errorf("expected updated section, got: %+v", got.Chunks)
	}
}

func TestListAll(t *testing.T) {
	root := t.TempDir()

	m1 := &Manifest{SchemaVersion: 1, Target: "claude", OutputPath: "CLAUDE.md"}
	m2 := &Manifest{SchemaVersion: 1, Target: "antigravity", OutputPath: "GEMINI.md"}

	if err := Write(root, "CLAUDE.md", m1); err != nil {
		t.Fatal(err)
	}
	if err := Write(root, "GEMINI.md", m2); err != nil {
		t.Fatal(err)
	}

	paths, err := ListAll(root)
	if err != nil {
		t.Fatalf("ListAll: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 paths, got %d: %v", len(paths), paths)
	}

	seen := map[string]bool{}
	for _, p := range paths {
		seen[p] = true
	}
	if !seen["CLAUDE.md"] {
		t.Error("expected CLAUDE.md in ListAll results")
	}
	if !seen["GEMINI.md"] {
		t.Error("expected GEMINI.md in ListAll results")
	}
}
