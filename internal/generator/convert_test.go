package generator

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildJctxSource mirrors the logic in cmd/jctx/convert.go (tested here at the generator level)
func buildJctxSource(target string, heading, content string) string {
	var sb strings.Builder
	sb.WriteString("---\ntargets: [")
	sb.WriteString(target)
	sb.WriteString("]\n---\n\n@@@ ")
	sb.WriteString(heading)
	sb.WriteString("\n")
	sb.WriteString(strings.TrimRight(content, "\n"))
	sb.WriteString("\n")
	return sb.String()
}

func TestConvert_GenAfterConvert(t *testing.T) {
	// Simulate converting CLAUDE.md → .jctx/rules/CLAUDE.md → gen reproduces content
	original := "Always write tests."

	root := buildJctxTree(t, map[string]string{
		// Simulated convert output
		".jctx/rules/CLAUDE.md": buildJctxSource("claude", "Core Rules", original),
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
	if !strings.Contains(results[0].Content, original) {
		t.Errorf("expected converted content to be reproducible:\n%s", results[0].Content)
	}
}

func TestConvert_DryRun_NoFiles(t *testing.T) {
	root := buildJctxTree(t, map[string]string{
		".jctx/rules/rules.md": buildJctxSource("claude", "Rules", "content"),
	})

	// Dry run: verify gen produces no files on disk
	_, err := Generate(GenOpts{
		Root:   root,
		Target: "claude",
		DryRun: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Confirm no output file was written
	outPath := filepath.Join(root, "CLAUDE.md")
	if _, err := os.Stat(outPath); !os.IsNotExist(err) {
		t.Error("dry run should not write CLAUDE.md")
	}
}
