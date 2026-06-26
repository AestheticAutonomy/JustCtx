package agents

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type AgentsProvider struct{}

func init() {
	providers.Register(&AgentsProvider{})
}

func (p *AgentsProvider) Name() string {
	return "agents"
}

func (p *AgentsProvider) SupportedTypes() []schema.Type {
	return []schema.Type{schema.TypeRules}
}

func (p *AgentsProvider) FindFiles(root string, t schema.Type) ([]string, error) {
	if t != schema.TypeRules {
		return nil, providers.ErrNotSupported
	}

	repoRoot := findRepoRoot(root)
	var found []string

	agentsPath := filepath.Clean(filepath.Join(repoRoot, "AGENTS.md"))
	if _, err := os.Stat(agentsPath); err == nil {
		found = append(found, agentsPath)
	}

	return found, nil
}

func (p *AgentsProvider) ParseRules(path string) ([]schema.Section, error) {
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

func (p *AgentsProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
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

	return []providers.OutputFile{{Path: "AGENTS.md", Content: sb.String()}}, nil
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
