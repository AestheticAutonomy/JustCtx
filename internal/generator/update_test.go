package generator

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// Tests for update-related gen behavior (called by cmd-update)

func TestGenerate_WithConfigDefaults(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/coding.md": "---\ntargets: [claude]\n---\n\n@@@ Core Rules\nAlways write tests.\n",
		".jctx/config.json":     `{"schema_version":1,"defaults":{"target":"claude","role":"","tags":[]}}`,
	})

	// Read the config and use it (mirrors what update command does)
	configPath := filepath.Join(root, ".jctx", "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatal(err)
	}
	var cfg struct {
		Defaults struct {
			Target string   `json:"target"`
			Role   string   `json:"role"`
			Tags   []string `json:"tags"`
		} `json:"defaults"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}

	results, err := Generate(GenOpts{
		Root:   root,
		Target: cfg.Defaults.Target,
		Role:   cfg.Defaults.Role,
		Tags:   cfg.Defaults.Tags,
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
}

func TestGenerate_UpdateDryRun(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/coding.md": "---\ntargets: [claude]\n---\n\n@@@ Core Rules\nContent here.\n",
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
	// Dry run: no file should be written
	for _, r := range results {
		if _, err := os.Stat(filepath.Join(root, r.OutputPath)); !os.IsNotExist(err) {
			t.Errorf("dry run should not write %s", r.OutputPath)
		}
	}
}
