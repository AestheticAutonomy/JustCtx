package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type ScanOpts struct {
	Root     string
	Target   string
	NoGlobal bool
	BottomUp bool
	Depth    int
}

type AssemblyState struct {
	sourceList []schema.Source
	chunks     []schema.Chunk
	nextID     int
}

func Scan(opts ScanOpts) (*schema.ScanResult, error) {
	// 1. Find repo root
	repoRoot := findRepoRoot(opts.Root)

	// 2. Look up provider
	provider, err := providers.Get(opts.Target)
	if err != nil {
		return nil, err
	}

	// 3. Find files from provider
	files, err := provider.FindFiles(opts.Root, schema.TypeRules)
	if err != nil {
		return nil, err
	}

	// 4. Filter and order files
	var finalFiles []string
	home, _ := os.UserHomeDir()
	var globalPath string
	if home != "" {
		globalPath = filepath.Clean(filepath.Join(home, ".claude", "CLAUDE.md"))
		globalGemini := filepath.Clean(filepath.Join(home, ".gemini", "GEMINI.md"))
		for _, file := range files {
			cleaned := filepath.Clean(file)
			if cleaned == globalPath || cleaned == globalGemini {
				if !opts.NoGlobal {
					finalFiles = append(finalFiles, file)
				}
			} else {
				finalFiles = append(finalFiles, file)
			}
		}
	} else {
		finalFiles = files
	}

	// If bottom-up with depth, filter files by directory level
	if opts.BottomUp && opts.Depth > 0 {
		absRoot, _ := filepath.Abs(opts.Root)
		var depthFiltered []string
		for _, file := range finalFiles {
			level := dirLevel(absRoot, filepath.Dir(filepath.Clean(file)))
			if level >= 0 && level <= opts.Depth {
				depthFiltered = append(depthFiltered, file)
			}
		}
		finalFiles = depthFiltered
	}

	// If bottom-up, reverse project files (or all files)
	if opts.BottomUp {
		// Reverse finalFiles
		for i, j := 0, len(finalFiles)-1; i < j; i, j = i+1, j-1 {
			finalFiles[i], finalFiles[j] = finalFiles[j], finalFiles[i]
		}
	}

	// 5. Assemble chunks and resolve imports
	state := &AssemblyState{
		nextID: 1,
	}

	assembledLineStart := 1
	for _, file := range finalFiles {
		loc := getFileLocation(file, repoRoot)
		err := state.parseFile(file, loc, 0, nil, &assembledLineStart)
		if err != nil {
			return nil, err
		}
	}

	// 6. Detect conflicts
	conflicts := DetectConflicts(state.chunks, state.sourceList)

	result := &schema.ScanResult{
		Envelope: schema.Envelope{
			SchemaVersion: 1,
			Command:       "scan",
			CWD:           opts.Root,
		},
		Sources:   state.sourceList,
		Assembled: state.chunks,
		Conflicts: conflicts,
		Findings:  []schema.Finding{},
	}

	return result, nil
}

func (state *AssemblyState) addSource(path string, loc schema.Location) string {
	cleaned := filepath.Clean(path)
	for _, src := range state.sourceList {
		if src.Path == cleaned {
			return src.ID
		}
	}
	id := fmt.Sprintf("s%d", state.nextID)
	state.nextID++

	bytes := 0
	if info, err := os.Stat(cleaned); err == nil {
		bytes = int(info.Size())
	}

	src := schema.Source{
		ID:       id,
		Path:     cleaned,
		Location: loc,
		Type:     schema.TypeRules,
		Bytes:    bytes,
	}
	state.sourceList = append(state.sourceList, src)
	return id
}

func (state *AssemblyState) parseFile(path string, loc schema.Location, depth int, chain []string, assembledLineStart *int) error {
	cleaned := filepath.Clean(path)
	if depth > 5 {
		return fmt.Errorf("max import depth exceeded (limit: 5)")
	}
	for _, p := range chain {
		if p == cleaned {
			cycleStr := strings.Join(append(chain, cleaned), " -> ")
			return fmt.Errorf("import cycle detected: %s", cycleStr)
		}
	}

	contentBytes, err := os.ReadFile(cleaned)
	if err != nil {
		return err
	}

	sourceID := state.addSource(cleaned, loc)

	contentStr := string(contentBytes)
	if contentStr == "" {
		return nil
	}

	lines := strings.Split(contentStr, "\n")
	var buffer []string
	startLine := 1

	for i, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		sourceLineNum := i + 1

		if strings.HasPrefix(trimmed, "@") && len(trimmed) > 1 {
			relPath := trimmed[1:]
			dir := filepath.Dir(cleaned)
			impPath := filepath.Clean(filepath.Join(dir, relPath))

			if _, err := os.Stat(impPath); err == nil {
				// Flush buffer first
				if len(buffer) > 0 {
					chunkContent := strings.Join(buffer, "\n") + "\n"
					linesCount := countLines(chunkContent)

					state.chunks = append(state.chunks, schema.Chunk{
						Content:            chunkContent,
						SourceID:           sourceID,
						AssembledLineStart: *assembledLineStart,
						AssembledLineEnd:   *assembledLineStart + linesCount - 1,
						SourceLineStart:    startLine,
						SourceLineEnd:      startLine + linesCount - 1,
					})
					*assembledLineStart += linesCount
					buffer = nil
				}

				// Recursively parse import
				err = state.parseFile(impPath, schema.LocationImport, depth+1, append(chain, cleaned), assembledLineStart)
				if err != nil {
					return err
				}

				startLine = sourceLineNum + 1
				continue
			}
		}

		buffer = append(buffer, line)
	}

	// Flush remaining buffer
	if len(buffer) > 0 {
		chunkContent := strings.Join(buffer, "\n")
		linesCount := countLines(chunkContent)

		state.chunks = append(state.chunks, schema.Chunk{
			Content:            chunkContent,
			SourceID:           sourceID,
			AssembledLineStart: *assembledLineStart,
			AssembledLineEnd:   *assembledLineStart + linesCount - 1,
			SourceLineStart:    startLine,
			SourceLineEnd:      startLine + linesCount - 1,
		})
		*assembledLineStart += linesCount
	}

	return nil
}

// Helpers

func findRepoRoot(start string) string {
	dir, err := filepath.Abs(start)
	if err != nil {
		return start
	}
	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return start
}

func getFileLocation(path string, repoRoot string) schema.Location {
	home, err := os.UserHomeDir()
	if err == nil {
		globalDir := filepath.Join(home, ".claude")
		globalGeminiDir := filepath.Join(home, ".gemini")
		cleanedPath := filepath.Clean(path)
		if strings.HasPrefix(cleanedPath, filepath.Clean(globalDir)) || strings.HasPrefix(cleanedPath, filepath.Clean(globalGeminiDir)) {
			return schema.LocationUserGlobal
		}
	}

	cleanedPath := filepath.Clean(path)
	cleanedRepo := filepath.Clean(repoRoot)

	if cleanedPath == filepath.Clean(filepath.Join(cleanedRepo, filepath.Base(cleanedPath))) {
		return schema.LocationProjectRoot
	}

	return schema.LocationSubdir
}

func countLines(s string) int {
	if s == "" {
		return 0
	}
	n := strings.Count(s, "\n")
	if !strings.HasSuffix(s, "\n") {
		n++
	}
	return n
}

// dirLevel returns how many directory levels fileDir is above root.
// Returns 0 if fileDir == root, 1 if fileDir is one level up, etc.
// Returns -1 if fileDir is not an ancestor of root.
func dirLevel(root, fileDir string) int {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return -1
	}
	absFile, err := filepath.Abs(fileDir)
	if err != nil {
		return -1
	}

	// fileDir must be root or an ancestor of root
	if absFile == absRoot {
		return 0
	}
	if !strings.HasPrefix(absRoot, absFile+string(filepath.Separator)) {
		return -1
	}

	rel, err := filepath.Rel(absFile, absRoot)
	if err != nil {
		return -1
	}
	return strings.Count(rel, string(filepath.Separator)) + 1
}
