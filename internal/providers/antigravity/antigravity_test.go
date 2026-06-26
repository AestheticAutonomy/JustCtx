package antigravity

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

func TestAntigravityProvider_FindFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Mock repo root
	repoRoot := filepath.Join(tmpDir, "myrepo")
	err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Write mock workspace files
	geminiPath := filepath.Join(repoRoot, "GEMINI.md")
	err = os.WriteFile(geminiPath, []byte("gemini rules"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	agentsPath := filepath.Join(repoRoot, "AGENTS.md")
	err = os.WriteFile(agentsPath, []byte("agents rules"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Set HOME and USERPROFILE to isolate global files
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	// Write global file
	globalGemini := filepath.Clean(filepath.Join(tempHome, ".gemini", "GEMINI.md"))
	err = os.MkdirAll(filepath.Dir(globalGemini), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(globalGemini, []byte("global gemini"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	provider := &AntigravityProvider{}
	files, err := provider.FindFiles(repoRoot, schema.TypeRules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{globalGemini, geminiPath, agentsPath}
	if len(files) != len(expected) {
		t.Fatalf("expected %d files, got %d: %v", len(expected), len(files), files)
	}

	for i, f := range files {
		expectedPath := filepath.Clean(expected[i])
		actualPath := filepath.Clean(f)
		if expectedPath != actualPath {
			t.Errorf("at index %d: expected %s, got %s", i, expectedPath, actualPath)
		}
	}
}

func TestAntigravityProvider_ParseRules(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "GEMINI.md")
	content := []byte("hello world")
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	provider := &AntigravityProvider{}
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
