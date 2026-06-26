package cursor

import (
	"os"
	"path/filepath"
	"strings"

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
	if t != schema.TypeRules {
		return nil, providers.ErrNotSupported
	}

	repoRoot := findRepoRoot(root)
	var found []string

	// 1. Legacy .cursorrules file at repo root
	legacyPath := filepath.Clean(filepath.Join(repoRoot, ".cursorrules"))
	if _, err := os.Stat(legacyPath); err == nil {
		found = append(found, legacyPath)
	}

	// 2. Project-scoped .cursor/rules/*.mdc
	rulesDir := filepath.Join(repoRoot, ".cursor", "rules")
	files, err := os.ReadDir(rulesDir)
	if err == nil {
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".mdc") {
				found = append(found, filepath.Clean(filepath.Join(rulesDir, file.Name())))
			}
		}
	}

	return found, nil
}

func (p *CursorProvider) ParseRules(path string) ([]schema.Section, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	
	raw := string(content)
	if strings.HasSuffix(path, ".mdc") {
		raw = stripFrontmatter(raw)
	}

	return []schema.Section{
		{
			Heading:    "",
			Content:    raw,
			SourceFile: path,
		},
	}, nil
}

func (p *CursorProvider) RenderRules(sections []schema.Section, opts providers.RenderOpts) ([]providers.OutputFile, error) {
	if len(sections) == 0 {
		return nil, nil
	}

	files := make([]providers.OutputFile, 0, len(sections))
	for _, s := range sections {
		slug := toSlug(s.Heading)
		if slug == "" {
			slug = "rules"
		}

		globs, hasGlobs := s.Dimensions["globs"]
		_, hasRole := s.Dimensions["role"]
		_, hasTag := s.Dimensions["tag"]
		alwaysApply := !hasRole && !hasTag && !hasGlobs

		var fm strings.Builder
		fm.WriteString("---\n")
		fm.WriteString("description: ")
		fm.WriteString(s.Heading)
		fm.WriteString("\n")
		if hasGlobs {
			fm.WriteString("globs: ")
			fm.WriteString(globs)
			fm.WriteString("\n")
		}
		if alwaysApply {
			fm.WriteString("alwaysApply: true\n")
		} else {
			fm.WriteString("alwaysApply: false\n")
		}
		fm.WriteString("---\n\n")
		fm.WriteString(strings.TrimRight(s.Content, "\n"))
		fm.WriteString("\n")

		files = append(files, providers.OutputFile{
			Path:    ".cursor/rules/" + slug + ".mdc",
			Content: fm.String(),
		})
	}
	return files, nil
}

func toSlug(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			b.WriteRune(r)
		} else if r == ' ' || r == '_' {
			b.WriteRune('-')
		}
	}
	// collapse consecutive hyphens
	result := b.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
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

func stripFrontmatter(content string) string {
	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, "---") {
		return content
	}

	lines := strings.Split(content, "\n")
	firstIdx := -1
	secondIdx := -1
	for i, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "---" {
			if firstIdx == -1 {
				firstIdx = i
			} else {
				secondIdx = i
				break
			}
		}
	}

	if firstIdx != -1 && secondIdx != -1 && secondIdx > firstIdx {
		if secondIdx+1 < len(lines) {
			return strings.Join(lines[secondIdx+1:], "\n")
		}
		return ""
	}
	return content
}
