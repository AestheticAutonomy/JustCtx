package claude

import (
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
	return nil, providers.ErrNotSupported
}

func (p *ClaudeProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	return nil, providers.ErrNotSupported
}

func (p *ClaudeProvider) ParseRules(path string) ([]schema.Section, error) {
	return nil, providers.ErrNotSupported
}
