package scanner

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

// Helper mock provider
type mockProvider struct {
	name  string
	files []string
}

func (p *mockProvider) Name() string { return p.name }
func (p *mockProvider) SupportedTypes() []schema.Type { return []schema.Type{schema.TypeRules} }
func (p *mockProvider) FindFiles(root string, t schema.Type) ([]string, error) {
	return p.files, nil
}
func (p *mockProvider) ParseRules(path string) ([]schema.Section, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return []schema.Section{{Content: string(content), SourceFile: path}}, nil
}
func (p *mockProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	return nil, providers.ErrNotSupported
}

func TestScan_TopDownAndBottomUp(t *testing.T) {
	tmpDir := t.TempDir()

	repoRoot := filepath.Join(tmpDir, "myrepo")
	err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	file1 := filepath.Join(repoRoot, "rule1.md")
	file2 := filepath.Join(repoRoot, "rule2.md")

	err = os.WriteFile(file1, []byte("rule 1 content\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(file2, []byte("rule 2 content\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mp := &mockProvider{name: "mock1", files: []string{file1, file2}}
	providers.Register(mp)

	// Test Top-Down
	res, err := Scan(ScanOpts{
		Root:     repoRoot,
		Target:   "mock1",
		NoGlobal: true,
		BottomUp: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Assembled) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(res.Assembled))
	}
	if res.Assembled[0].Content != "rule 1 content\n" {
		t.Errorf("unexpected chunk 0 content: '%s'", res.Assembled[0].Content)
	}

	// Test Bottom-Up
	resBU, err := Scan(ScanOpts{
		Root:     repoRoot,
		Target:   "mock1",
		NoGlobal: true,
		BottomUp: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(resBU.Assembled) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(resBU.Assembled))
	}
	if resBU.Assembled[0].Content != "rule 2 content\n" {
		t.Errorf("expected rule 2 content first for bottom-up, got '%s'", resBU.Assembled[0].Content)
	}
}

func TestScan_ImportResolutionAndCycle(t *testing.T) {
	tmpDir := t.TempDir()

	repoRoot := filepath.Join(tmpDir, "myrepo")
	err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	// 1. Valid import chain
	fileA := filepath.Join(repoRoot, "a.md")
	fileB := filepath.Join(repoRoot, "b.md")

	err = os.WriteFile(fileA, []byte("start a\n@b.md\nend a\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(fileB, []byte("content b\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mp := &mockProvider{name: "mock2", files: []string{fileA}}
	providers.Register(mp)

	res, err := Scan(ScanOpts{
		Root:     repoRoot,
		Target:   "mock2",
		NoGlobal: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// We expect 3 chunks: "start a\n", "content b\n", "end a\n"
	if len(res.Assembled) != 3 {
		t.Fatalf("expected 3 chunks, got %d: %v", len(res.Assembled), res.Assembled)
	}
	if res.Assembled[1].Content != "content b\n" {
		t.Errorf("expected chunk 1 to be import b, got '%s'", res.Assembled[1].Content)
	}

	// 2. Cycle detection
	fileCycleA := filepath.Join(repoRoot, "cycleA.md")
	fileCycleB := filepath.Join(repoRoot, "cycleB.md")

	err = os.WriteFile(fileCycleA, []byte("cycle a\n@cycleB.md\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(fileCycleB, []byte("cycle b\n@cycleA.md\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mpCycle := &mockProvider{name: "mockCycle", files: []string{fileCycleA}}
	providers.Register(mpCycle)

	_, err = Scan(ScanOpts{
		Root:     repoRoot,
		Target:   "mockCycle",
		NoGlobal: true,
	})
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestScan_Conflicts(t *testing.T) {
	tmpDir := t.TempDir()

	repoRoot := filepath.Join(tmpDir, "myrepo")
	err := os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)
	if err != nil {
		t.Fatal(err)
	}

	file1 := filepath.Join(repoRoot, "rule1.md")
	file2 := filepath.Join(repoRoot, "rule2.md")

	err = os.WriteFile(file1, []byte("## Core Principles\nAlways use spaces instead of tabs\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(file2, []byte("## Core Principles\nNever use spaces instead of tabs\n"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	mp := &mockProvider{name: "mockConflict", files: []string{file1, file2}}
	providers.Register(mp)

	res, err := Scan(ScanOpts{
		Root:     repoRoot,
		Target:   "mockConflict",
		NoGlobal: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// We expect two types of conflicts: duplicate_heading and contradicting_imperative
	if len(res.Conflicts) == 0 {
		t.Fatal("expected conflicts to be detected")
	}

	hasDuplicateHeading := false
	hasContradiction := false
	for _, c := range res.Conflicts {
		if c.Type == "duplicate_heading" {
			hasDuplicateHeading = true
		}
		if c.Type == "contradicting_imperative" {
			hasContradiction = true
		}
	}

	if !hasDuplicateHeading {
		t.Error("expected duplicate_heading conflict")
	}
	if !hasContradiction {
		t.Error("expected contradicting_imperative conflict")
	}
}

func TestScan_BottomUpDepth(t *testing.T) {
	// Build a 5-level deep tree:
	// repo/.git
	// repo/L0.md         (level 4 above cwd)
	// repo/a/L1.md       (level 3)
	// repo/a/b/L2.md     (level 2)
	// repo/a/b/c/L3.md   (level 1)
	// repo/a/b/c/d/L4.md (level 0 = cwd)
	tmpDir := t.TempDir()
	repoRoot := filepath.Join(tmpDir, "repo")
	os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)

	dirs := []string{
		repoRoot,
		filepath.Join(repoRoot, "a"),
		filepath.Join(repoRoot, "a", "b"),
		filepath.Join(repoRoot, "a", "b", "c"),
		filepath.Join(repoRoot, "a", "b", "c", "d"),
	}
	var allFiles []string
	for i, dir := range dirs {
		os.MkdirAll(dir, 0755)
		f := filepath.Join(dir, "L"+string(rune('0'+i))+".md")
		os.WriteFile(f, []byte("level "+string(rune('0'+i))+"\n"), 0644)
		allFiles = append(allFiles, f)
	}

	cwd := dirs[4] // repo/a/b/c/d

	mp := &mockProvider{name: "mockDepth", files: allFiles}
	providers.Register(mp)

	// Without depth: all 5 files returned
	res, err := Scan(ScanOpts{
		Root:     cwd,
		Target:   "mockDepth",
		NoGlobal: true,
		BottomUp: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sources) != 5 {
		t.Fatalf("no depth: expected 5 sources, got %d", len(res.Sources))
	}

	// Depth 2: cwd + 2 levels up = 3 directories (d, c, b)
	res, err = Scan(ScanOpts{
		Root:     cwd,
		Target:   "mockDepth",
		NoGlobal: true,
		BottomUp: true,
		Depth:    2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sources) != 3 {
		var paths []string
		for _, s := range res.Sources {
			paths = append(paths, s.Path)
		}
		t.Fatalf("depth 2: expected 3 sources, got %d: %v", len(res.Sources), paths)
	}

	// Depth 0 (unlimited) with bottom-up: all files
	res, err = Scan(ScanOpts{
		Root:     cwd,
		Target:   "mockDepth",
		NoGlobal: true,
		BottomUp: true,
		Depth:    0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sources) != 5 {
		t.Fatalf("depth 0: expected 5 sources, got %d", len(res.Sources))
	}
}

func TestScan_DepthWithoutBottomUp(t *testing.T) {
	// Depth without bottom-up should have no filtering effect at the scanner level.
	// The warning is printed by the CLI layer (cmd/jctx/scan.go), not the scanner.
	tmpDir := t.TempDir()
	repoRoot := filepath.Join(tmpDir, "repo")
	os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)

	f1 := filepath.Join(repoRoot, "rules.md")
	os.WriteFile(f1, []byte("content\n"), 0644)

	mp := &mockProvider{name: "mockDepthNoBottomUp", files: []string{f1}}
	providers.Register(mp)

	res, err := Scan(ScanOpts{
		Root:     repoRoot,
		Target:   "mockDepthNoBottomUp",
		NoGlobal: true,
		BottomUp: false,
		Depth:    2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(res.Sources))
	}
}
