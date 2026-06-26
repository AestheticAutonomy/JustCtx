package antigravity

import (
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
	return nil, providers.ErrNotSupported
}

func (p *AntigravityProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	return nil, providers.ErrNotSupported
}

func (p *AntigravityProvider) ParseRules(path string) ([]schema.Section, error) {
	return nil, providers.ErrNotSupported
}
