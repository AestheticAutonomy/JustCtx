package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/AestheticAutonomy/justctx/internal/providers/agents"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/antigravity"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/claude"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/cursor"
)

// buildJctxTree creates a minimal .jctx/rules/ tree in a temp dir.
func buildJctxTree(t *testing.T, files map[string]string) string {
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
	// Need a .git dir so provider helpers can find repo root
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestGenerate_BasicOutput(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/coding.md": `---
targets: [claude]
---

@@@ Core Rules
Always write tests.
`,
	})

	results, err := Generate(GenOpts{
		Root:   root,
		Target: "claude",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].OutputPath != "CLAUDE.md" {
		t.Errorf("unexpected output path: %s", results[0].OutputPath)
	}
	if !strings.Contains(results[0].Content, "Core Rules") {
		t.Errorf("expected heading in content:\n%s", results[0].Content)
	}
	if !strings.Contains(results[0].Content, "Always write tests.") {
		t.Errorf("expected content body:\n%s", results[0].Content)
	}
}

func TestGenerate_DryRun_NoFiles(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/coding.md": `---
targets: [claude]
---

@@@ Core Rules
Always write tests.
`,
	})

	results, err := Generate(GenOpts{
		Root:   root,
		Target: "claude",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}

	// Verify no file was written
	outPath := filepath.Join(root, "CLAUDE.md")
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Error("expected no output file in dry-run mode")
	}
}

func TestGenerate_DimensionFiltering(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/coding.md": `---
targets: [claude]
---

@@@ Core Rules
Always write tests.

@@@ DB Review [role:dbreviewer]
Check for N+1s.
`,
	})

	// Without role: DB Review section excluded
	results, err := Generate(GenOpts{
		Root:   root,
		Target: "claude",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if strings.Contains(results[0].Content, "N+1") {
		t.Error("role-filtered section should not appear without role flag")
	}
	if !strings.Contains(results[0].Content, "Always write tests.") {
		t.Error("unfiltered section should be present")
	}

	// With role: DB Review included
	results2, err := Generate(GenOpts{
		Root:   root,
		Target: "claude",
		Role:   "dbreviewer",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("Generate with role: %v", err)
	}
	if len(results2) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results2))
	}
	if !strings.Contains(results2[0].Content, "N+1") {
		t.Error("role-matched section should appear when role is active")
	}
}

func TestGenerate_LocalOverridesBase(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		// Base file
		".jctx/rules/coding.md": `---
targets: [claude]
---

@@@ Core Rules
Base content.
`,
		// Local override with same filename
		".jctx/.local/rules/coding.md": `---
targets: [claude]
---

@@@ Core Rules
Local override content.
`,
	})

	results, err := Generate(GenOpts{
		Root:   root,
		Target: "claude",
		DryRun: true,
	})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if strings.Contains(results[0].Content, "Base content.") {
		t.Error("base content should be overridden by local")
	}
	if !strings.Contains(results[0].Content, "Local override content.") {
		t.Errorf("expected local override content, got:\n%s", results[0].Content)
	}
}
