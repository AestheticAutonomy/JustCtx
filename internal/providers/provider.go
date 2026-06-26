package providers

import (
	"errors"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type RenderOpts struct {
	Role     string
	Tags     []string
	Annotate bool
}

type OutputFile struct {
	Path    string
	Content string
}

type Provider interface {
	Name() string
	SupportedTypes() []schema.Type
	FindFiles(root string, t schema.Type) ([]string, error)
	RenderRules(sections []schema.Section, opts RenderOpts) ([]OutputFile, error)
	ParseRules(path string) ([]schema.Section, error)
}

var ErrNotSupported = errors.New("not supported")
