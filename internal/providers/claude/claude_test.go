package claude

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

func TestClaudeProvider_FindFiles(t *testing.T) {
	// Create temp dir
	tmpDir := t.TempDir()

	// Create mock repo root
	repoRoot := filepath.Join(tmpDir, "myrepo")
	err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Create subdirectories
	subDir1 := filepath.Join(repoRoot, "services")
	subDir2 := filepath.Join(subDir1, "auth")
	err = os.MkdirAll(subDir2, 0755)
	if err != nil {
		t.Fatal(err)
	}

	// Write CLAUDE.md files
	rootClaude := filepath.Join(repoRoot, "CLAUDE.md")
	err = os.WriteFile(rootClaude, []byte("root rules"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	authClaude := filepath.Join(subDir2, "CLAUDE.md")
	err = os.WriteFile(authClaude, []byte("auth rules"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Set HOME and USERPROFILE to a temp dir to isolate from host config
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	t.Setenv("USERPROFILE", tempHome)

	// Write a mock global file
	globalClaude := filepath.Clean(filepath.Join(tempHome, ".claude", "CLAUDE.md"))
	err = os.MkdirAll(filepath.Dir(globalClaude), 0755)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(globalClaude, []byte("global rules"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	provider := &ClaudeProvider{}
	files, err := provider.FindFiles(subDir2, schema.TypeRules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []string{globalClaude, rootClaude, authClaude}
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

func TestClaudeProvider_ParseRules(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "CLAUDE.md")
	content := []byte("hello world")
	err := os.WriteFile(testFile, content, 0644)
	if err != nil {
		t.Fatal(err)
	}

	provider := &ClaudeProvider{}
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
	if sections[0].SourceFile != testFile {
		t.Errorf("expected source file %s, got %s", testFile, sections[0].SourceFile)
	}
}
