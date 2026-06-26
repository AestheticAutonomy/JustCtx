package antigravity

import (
	"os"
	"path/filepath"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type AntigravityProvider struct{}

func init() {
	providers.Register(&AntigravityProvider{})
}

func (p *AntigravityProvider) Name() string {
	return "antigravity"
}

func (p *AntigravityProvider) SupportedTypes() []schema.Type {
	return []schema.Type{schema.TypeRules}
}

func (p *AntigravityProvider) FindFiles(root string, t schema.Type) ([]string, error) {
	if t != schema.TypeRules {
		return nil, providers.ErrNotSupported
	}

	var found []string

	// 1. Global config file: ~/.gemini/GEMINI.md
	home, err := os.UserHomeDir()
	if err == nil {
		globalPath := filepath.Clean(filepath.Join(home, ".gemini", "GEMINI.md"))
		if _, err := os.Stat(globalPath); err == nil {
			found = append(found, globalPath)
		}
	}

	// 2. Project workspace files: GEMINI.md and AGENTS.md at repo root
	repoRoot := findRepoRoot(root)

	geminiPath := filepath.Clean(filepath.Join(repoRoot, "GEMINI.md"))
	if _, err := os.Stat(geminiPath); err == nil {
		alreadyAdded := false
		for _, f := range found {
			if f == geminiPath {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			found = append(found, geminiPath)
		}
	}

	agentsPath := filepath.Clean(filepath.Join(repoRoot, "AGENTS.md"))
	if _, err := os.Stat(agentsPath); err == nil {
		alreadyAdded := false
		for _, f := range found {
			if f == agentsPath {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			found = append(found, agentsPath)
		}
	}

	return found, nil
}

func (p *AntigravityProvider) ParseRules(path string) ([]schema.Section, error) {
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

func (p *AntigravityProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	return nil, providers.ErrNotSupported
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
