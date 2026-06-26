package agents

import (
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
	return nil, providers.ErrNotSupported
}

func (p *AgentsProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	return nil, providers.ErrNotSupported
}

func (p *AgentsProvider) ParseRules(path string) ([]schema.Section, error) {
	return nil, providers.ErrNotSupported
}
