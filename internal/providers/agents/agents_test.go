package agents

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

func TestAgentsProvider_RenderRules(t *testing.T) {
	p := &AgentsProvider{}
	sections := []schema.Section{
		{Heading: "Core Rules", Content: "Always test."},
		{Heading: "Style", Content: "Be consistent."},
	}
	files, err := p.RenderRules(sections, providers.RenderOpts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("expected 1 output file, got %d", len(files))
	}
	if files[0].Path != "AGENTS.md" {
		t.Errorf("expected path AGENTS.md, got %s", files[0].Path)
	}
	want := "## Core Rules\n\nAlways test.\n\n## Style\n\nBe consistent.\n"
	if files[0].Content != want {
		t.Errorf("content mismatch\nwant: %q\ngot:  %q", want, files[0].Content)
	}
}

func TestAgentsProvider_RenderRules_Empty(t *testing.T) {
	p := &AgentsProvider{}
	files, err := p.RenderRules(nil, providers.RenderOpts{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected empty slice, got %d files", len(files))
	}
}

func TestAgentsProvider_FindFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Mock repo root
	repoRoot := filepath.Join(tmpDir, "myrepo")
	err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Write mock agents file
	agentsPath := filepath.Join(repoRoot, "AGENTS.md")
	err = os.WriteFile(agentsPath, []byte("agents rules"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	provider := &AgentsProvider{}
	files, err := provider.FindFiles(repoRoot, schema.TypeRules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{agentsPath}
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(files), files)
	}

	expectedPath := filepath.Clean(expected[0])
	actualPath := filepath.Clean(files[0])
	if expectedPath != actualPath {
		t.Errorf("expected %s, got %s", expectedPath, actualPath)
	}
}

func TestAgentsProvider_ParseRules(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "AGENTS.md")
	content := []byte("hello world")
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	provider := &AgentsProvider{}
	sections, err := provider.ParseRules(testFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(sections))
	}

	if sections[0].Content != "hello world" {
		t.Errorf("expected content 'hello world', got '%s'", sections[0].Content)
	}
}
