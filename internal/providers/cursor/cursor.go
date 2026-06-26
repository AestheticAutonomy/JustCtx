package cursor

import (
	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type CursorProvider struct{}

func init() {
	providers.Register(&CursorProvider{})
}

func (p *CursorProvider) Name() string {
	return "cursor"
}

func (p *CursorProvider) SupportedTypes() []schema.Type {
	return []schema.Type{schema.TypeRules}
}

func (p *CursorProvider) FindFiles(root string, t schema.Type) ([]string, error) {
	return nil, providers.ErrNotSupported
}

func (p *CursorProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	return nil, providers.ErrNotSupported
}

func (p *CursorProvider) ParseRules(path string) ([]schema.Section, error) {
	return nil, providers.ErrNotSupported
}
