package claude

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type ClaudeProvider struct{}

func init() {
	providers.Register(&ClaudeProvider{})
}

func (p *ClaudeProvider) Name() string {
	return "claude"
}

func (p *ClaudeProvider) SupportedTypes() []schema.Type {
	return []schema.Type{schema.TypeRules}
}

func (p *ClaudeProvider) FindFiles(root string, t schema.Type) ([]string, error) {
	if t != schema.TypeRules {
		return nil, providers.ErrNotSupported
	}

	var found []string

	// 1. Global config file: ~/.claude/CLAUDE.md
	home, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Clean(filepath.Join(home, ".claude", "CLAUDE.md"))
		if _, err := os.Stat(globalPath); err == nil {
			found = append(found, globalPath)
		}
	}

	// 2. Project hierarchy files
	repoRoot := findRepoRoot(root)
	segments, err := getPathSegments(repoRoot, root)
	if err != nil {
		return nil, err
	}

	for _, dir := range segments {
		path := filepath.Clean(filepath.Join(dir, "CLAUDE.md"))
		if _, err := os.Stat(path); err == nil {
			// Avoid adding the global file twice if repo root is the home dir
			alreadyAdded := false
			for _, f := range found {
				if f == path {
					alreadyAdded = true
					break
				}
			}
			if !alreadyAdded {
				found = append(found, path)
			}
		}
	}

	return found, nil
}

func (p *ClaudeProvider) ParseRules(path string) ([]schema.Section, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return []schema.Section{
		{
			Heading:    "",
			Content:    string(content),
			SourceFile: path,
		},
	}, nil
}

func (p *ClaudeProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	if len(sections) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	for i, s := range sections {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		if s.Heading != "" {
			sb.WriteString("## ")
			sb.WriteString(s.Heading)
			sb.WriteString("\n\n")
		}
		sb.WriteString(strings.TrimRight(s.Content, "\n"))
	}
	sb.WriteString("\n")

	return []providers.OutputFile{{Path: "CLAUDE.md", Content: sb.String()}}, nil
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

func getPathSegments(repoRoot, target string) ([]string, error) {
	absRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, err
	}
	absTarget, err := filepath.Abs(target)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(absTarget, absRoot) {
		return []string{absTarget}, nil
	}

	var segments []string
	curr := absTarget
	for {
		segments = append([]string{curr}, segments...)
		if curr == absRoot {
			break
		}
		parent := filepath.Dir(curr)
		if parent == curr {
			break
		}
		curr = parent
	}
	return segments, nil
}
