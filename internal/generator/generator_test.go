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

func TestGenerate_RemoteLayerOverride(t *testing.T) {
	// Same basename in remote + base: base wins (last-write wins in fileMap)
	root := buildJctxTree(t, map[string]string{
		".jctx/.remote/pkg1/rules/shared.md": "---\ntargets: [claude]\n---\n\n@@@ Remote Rule\nRemote content.\n",
		".jctx/rules/shared.md":              "---\ntargets: [claude]\n---\n\n@@@ Base Rule\nBase content.\n",
	})

	results, err := Generate(GenOpts{Root: root, Target: "claude", DryRun: true})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if strings.Contains(results[0].Content, "Remote content.") {
		t.Error("remote content should be overridden by base (same basename)")
	}
	if !strings.Contains(results[0].Content, "Base content.") {
		t.Errorf("base content should win over remote:\n%s", results[0].Content)
	}
}

func TestGenerate_RemoteOnlyFile(t *testing.T) {
	// File only in .remote/ with unique basename should be included
	root := buildJctxTree(t, map[string]string{
		".jctx/.remote/pkg1/rules/remote-only.md": "---\ntargets: [claude]\n---\n\n@@@ Remote Only\nOnly in remote.\n",
		".jctx/rules/base-rules.md":                "---\ntargets: [claude]\n---\n\n@@@ Base Rules\nBase content.\n",
	})

	results, err := Generate(GenOpts{Root: root, Target: "claude", DryRun: true})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Content, "Only in remote.") {
		t.Errorf("remote-only file should be included:\n%s", results[0].Content)
	}
	if !strings.Contains(results[0].Content, "Base content.") {
		t.Errorf("base content should also be included:\n%s", results[0].Content)
	}
}

func TestGenerate_TagFiltering(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/rules.md": `---
targets: [claude]
---

@@@ Core Rules
Always write tests.

@@@ Postgres Rules [tag:pg]
Use pgx driver.
`,
	})

	// No tags: postgres section excluded
	results, err := Generate(GenOpts{Root: root, Target: "claude", DryRun: true})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}
	if strings.Contains(results[0].Content, "pgx driver") {
		t.Error("tag-filtered section should not appear when no tags")
	}
	if !strings.Contains(results[0].Content, "Always write tests.") {
		t.Error("untagged section should appear")
	}

	// Tag=pg: postgres section included
	results2, err := Generate(GenOpts{Root: root, Target: "claude", Tags: []string{"pg"}, DryRun: true})
	if err != nil {
		t.Fatalf("Generate with tag: %v", err)
	}
	if !strings.Contains(results2[0].Content, "pgx driver") {
		t.Error("tag-matched section should appear when tag=pg")
	}

	// Tag=mysql: postgres section excluded
	results3, err := Generate(GenOpts{Root: root, Target: "claude", Tags: []string{"mysql"}, DryRun: true})
	if err != nil {
		t.Fatalf("Generate with mysql tag: %v", err)
	}
	if strings.Contains(results3[0].Content, "pgx driver") {
		t.Error("tag=[mysql] should not include pg section")
	}
}

func TestGenerate_WriteToFile(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nWrite clean code.\n",
	})

	_, err := Generate(GenOpts{Root: root, Target: "claude", DryRun: false})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	outPath := filepath.Join(root, "CLAUDE.md")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("CLAUDE.md not written: %v", err)
	}
	if !strings.Contains(string(data), "Write clean code.") {
		t.Errorf("output file content incorrect:\n%s", string(data))
	}
}

func TestGenerate_OverwritesExisting(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Rules\nFirst generation.\n",
	})

	// First gen
	_, err := Generate(GenOpts{Root: root, Target: "claude", DryRun: false})
	if err != nil {
		t.Fatalf("first Generate: %v", err)
	}

	// Update the source
	err = os.WriteFile(filepath.Join(root, ".jctx", "rules", "rules.md"),
		[]byte("---\ntargets: [claude]\n---\n\n@@@ Rules\nSecond generation.\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Second gen
	_, err = Generate(GenOpts{Root: root, Target: "claude", DryRun: false})
	if err != nil {
		t.Fatalf("second Generate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, "CLAUDE.md"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "First generation.") {
		t.Error("expected first gen content to be overwritten")
	}
	if !strings.Contains(string(data), "Second generation.") {
		t.Errorf("expected second gen content:\n%s", string(data))
	}
}

func TestGenerate_ManifestWritten(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ Core Rules\nContent here.\n",
	})

	_, err := Generate(GenOpts{Root: root, Target: "claude", DryRun: false})
	if err != nil {
		t.Fatalf("Generate: %v", err)
	}

	manifestPath := filepath.Join(root, ".jctx", ".manifest", "CLAUDE.md.json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Fatal("manifest not written after gen")
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"Core Rules"`) {
		t.Errorf("manifest should contain section name:\n%s", string(data))
	}
	if !strings.Contains(string(data), `"claude"`) {
		t.Errorf("manifest should contain target:\n%s", string(data))
	}
}

func TestGenerate_ManifestUpdatedOnRegen(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/rules.md": "---\ntargets: [claude]\n---\n\n@@@ First Section\nOriginal.\n",
	})

	_, err := Generate(GenOpts{Root: root, Target: "claude", DryRun: false})
	if err != nil {
		t.Fatalf("first Generate: %v", err)
	}

	// Replace source
	os.WriteFile(filepath.Join(root, ".jctx", "rules", "rules.md"),
		[]byte("---\ntargets: [claude]\n---\n\n@@@ Updated Section\nNew content.\n"), 0644)

	_, err = Generate(GenOpts{Root: root, Target: "claude", DryRun: false})
	if err != nil {
		t.Fatalf("second Generate: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(root, ".jctx", ".manifest", "CLAUDE.md.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), `"First Section"`) {
		t.Error("old section should not appear in updated manifest")
	}
	if !strings.Contains(string(data), `"Updated Section"`) {
		t.Errorf("updated section should appear in manifest:\n%s", string(data))
	}
}

func TestGenerate_AllProviders(t *testing.T) {
	// Each provider should generate its own output file from the same .jctx/ source
	cases := []struct {
		target   string
		outFile  string
		fileBody string
	}{
		{"claude", "CLAUDE.md", "---\ntargets: [claude]\n---\n\n@@@ Rules\nClaude rules.\n"},
		{"antigravity", "GEMINI.md", "---\ntargets: [antigravity]\n---\n\n@@@ Rules\nGemini rules.\n"},
		{"agents", "AGENTS.md", "---\ntargets: [agents]\n---\n\n@@@ Rules\nAgents rules.\n"},
	}

	for _, tc := range cases {
		t.Run(tc.target, func(t *testing.T) {
			root := buildJctxTree(t, map[string]string{
				".jctx/rules/rules.md": tc.fileBody,
			})

			_, err := Generate(GenOpts{Root: root, Target: tc.target, DryRun: false})
			if err != nil {
				t.Fatalf("Generate for %s: %v", tc.target, err)
			}

			outPath := filepath.Join(root, tc.outFile)
			if _, err := os.Stat(outPath); os.IsNotExist(err) {
				t.Errorf("%s not written for target %s", tc.outFile, tc.target)
			}
		})
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
