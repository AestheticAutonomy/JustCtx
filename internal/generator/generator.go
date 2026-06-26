package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/AestheticAutonomy/justctx/internal/manifest"
	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type GenOpts struct {
	Root   string
	Target string
	Role   string
	Tags   []string
	DryRun bool
}

type excludedSection struct {
	heading string
	reason  string
}

type fileFrontmatter struct {
	Targets []string `yaml:"targets"`
}

// Generate generates output files for a single target. Returns one GenResult per output file.
func Generate(opts GenOpts) ([]schema.GenResult, error) {
	p, err := providers.Get(opts.Target)
	if err != nil {
		return nil, err
	}

	sections, err := loadSections(opts.Root, opts.Target)
	if err != nil {
		return nil, err
	}

	filtered, excluded := filterSections(sections, opts)

	renderOpts := providers.RenderOpts{
		Role: opts.Role,
		Tags: opts.Tags,
	}

	outputFiles, err := p.RenderRules(filtered, renderOpts)
	if err != nil {
		return nil, err
	}

	var included []schema.SectionIncluded
	for _, s := range filtered {
		included = append(included, schema.SectionIncluded{
			Heading: s.Heading,
			Source:  s.SourceFile,
		})
	}
	var excludedOut []schema.SectionExcluded
	for _, s := range excluded {
		excludedOut = append(excludedOut, schema.SectionExcluded{
			Heading: s.heading,
			Reason:  s.reason,
		})
	}

	var results []schema.GenResult

	for _, f := range outputFiles {
		absPath := filepath.Join(opts.Root, f.Path)
		mPath := filepath.Join(opts.Root, ".jctx", ".manifest", filepath.Base(f.Path)+".json")

		res := schema.GenResult{
			Envelope: schema.Envelope{
				SchemaVersion: 1,
				Command:       "gen",
				CWD:           opts.Root,
			},
			OutputPath:       f.Path,
			Role:             opts.Role,
			Tags:             opts.Tags,
			SectionsIncluded: included,
			SectionsExcluded: excludedOut,
			ManifestPath:     mPath,
			Content:          f.Content,
		}

		if !opts.DryRun {
			if err := os.MkdirAll(filepath.Dir(absPath), 0755); err != nil {
				return nil, fmt.Errorf("creating output dir for %s: %w", f.Path, err)
			}
			if err := os.WriteFile(absPath, []byte(f.Content), 0644); err != nil {
				return nil, fmt.Errorf("writing %s: %w", f.Path, err)
			}
			m := buildManifest(opts.Target, f.Path, filtered)
			if err := manifest.Write(opts.Root, f.Path, m); err != nil {
				return nil, fmt.Errorf("writing manifest for %s: %w", f.Path, err)
			}
		}

		results = append(results, res)
	}

	return results, nil
}

// loadSections reads all .jctx/rules/*.md files in merge order, parsing file frontmatter
// and @@@ sections. Files with targets that don't include the requested target are skipped.
func loadSections(root, target string) ([]schema.Section, error) {
	// Build merged file map: basename → absolute path. Later dirs win.
	fileMap := map[string]string{}
	var fileOrder []string

	addDir := func(dir string) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
				continue
			}
			base := e.Name()
			if _, exists := fileMap[base]; !exists {
				fileOrder = append(fileOrder, base)
			}
			fileMap[base] = filepath.Join(dir, e.Name())
		}
	}

	// 1. Remote packages
	remoteDir := filepath.Join(root, ".jctx", ".remote")
	if entries, err := os.ReadDir(remoteDir); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				addDir(filepath.Join(remoteDir, e.Name(), "rules"))
			}
		}
	}

	// 2. Base rules
	addDir(filepath.Join(root, ".jctx", "rules"))

	// 3. Local overrides
	addDir(filepath.Join(root, ".jctx", ".local", "rules"))

	sort.Strings(fileOrder)

	var sections []schema.Section
	for _, base := range fileOrder {
		path := fileMap[base]
		fileSections, err := parseJctxFile(path, target)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", path, err)
		}
		sections = append(sections, fileSections...)
	}
	return sections, nil
}

// parseJctxFile reads a .jctx/rules/*.md file, extracts YAML frontmatter and @@@ sections.
// Returns nil if the file's targets don't include the requested target.
func parseJctxFile(path, target string) ([]schema.Section, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(data)

	// Extract YAML frontmatter
	var fm fileFrontmatter
	body := content
	if strings.HasPrefix(strings.TrimSpace(content), "---") {
		parts := strings.SplitN(content, "---", 3)
		if len(parts) >= 3 {
			if err := yaml.Unmarshal([]byte(parts[1]), &fm); err != nil {
				return nil, fmt.Errorf("parsing frontmatter in %s: %w", path, err)
			}
			body = parts[2]
		}
	}

	// File-level target filter: if targets declared and our target not in the list, skip
	if len(fm.Targets) > 0 && !containsStr(fm.Targets, target) {
		return nil, nil
	}

	return parseSections(body, path), nil
}

var dimRegexp = regexp.MustCompile(`\[([^\]]+)\]`)

// parseSections splits content on @@@ markers and parses dimension annotations.
func parseSections(content, sourcePath string) []schema.Section {
	lines := strings.Split(content, "\n")

	type rawSection struct {
		heading    string
		dims       map[string][]string // key → []values (multiple [tag:x] allowed)
		bodyLines  []string
		lineStart  int
	}

	var rawSections []rawSection
	var current *rawSection
	lineNum := 1

	for _, line := range lines {
		trimmed := strings.TrimRight(line, "\r")
		if strings.HasPrefix(trimmed, "@@@") {
			if current != nil {
				rawSections = append(rawSections, *current)
			}
			rest := strings.TrimSpace(trimmed[3:])
			heading, dims := parseSectionHeader(rest)
			current = &rawSection{
				heading:   heading,
				dims:      dims,
				lineStart: lineNum + 1,
			}
		} else {
			if current != nil {
				current.bodyLines = append(current.bodyLines, trimmed)
			}
		}
		lineNum++
	}
	if current != nil {
		rawSections = append(rawSections, *current)
	}

	// If no @@@ markers, treat the whole file as one unnamed section
	if len(rawSections) == 0 {
		trimmed := strings.TrimSpace(content)
		if trimmed != "" {
			return []schema.Section{{
				Content:    trimmed,
				SourceFile: sourcePath,
				LineStart:  1,
				LineEnd:    strings.Count(content, "\n") + 1,
			}}
		}
		return nil
	}

	var sections []schema.Section
	for _, r := range rawSections {
		body := strings.Join(r.bodyLines, "\n")
		body = strings.TrimSpace(body)
		if body == "" && r.heading == "" {
			continue
		}

		dims := make(map[string]string)
		for k, vals := range r.dims {
			dims[k] = strings.Join(vals, ",")
		}

		sections = append(sections, schema.Section{
			Heading:    r.heading,
			Content:    body,
			Dimensions: dims,
			SourceFile: sourcePath,
			LineStart:  r.lineStart,
		})
	}
	return sections
}

// parseSectionHeader extracts heading and dimension map from a @@@ line's content.
func parseSectionHeader(rest string) (string, map[string][]string) {
	dims := map[string][]string{}
	matches := dimRegexp.FindAllStringIndex(rest, -1)

	headingEnd := len(rest)
	if len(matches) > 0 {
		headingEnd = matches[0][0]
	}
	heading := strings.TrimSpace(rest[:headingEnd])

	for _, m := range dimRegexp.FindAllString(rest, -1) {
		inner := m[1 : len(m)-1] // strip [ ]
		// Handle & within group (AND): for now treat as separate constraints
		parts := strings.Split(inner, "&")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			idx := strings.Index(part, ":")
			if idx < 0 {
				continue
			}
			k := strings.TrimSpace(part[:idx])
			v := strings.TrimSpace(part[idx+1:])
			dims[k] = append(dims[k], v)
		}
	}
	return heading, dims
}

// filterSections partitions sections into included and excluded based on opts.
func filterSections(sections []schema.Section, opts GenOpts) ([]schema.Section, []excludedSection) {
	var included []schema.Section
	var excluded []excludedSection

	for _, s := range sections {
		reason := matchesDimensions(s, opts)
		if reason == "" {
			included = append(included, s)
		} else {
			excluded = append(excluded, excludedSection{heading: s.Heading, reason: reason})
		}
	}
	return included, excluded
}

// matchesDimensions returns "" if the section should be included, or a reason string if excluded.
func matchesDimensions(s schema.Section, opts GenOpts) string {
	dims := s.Dimensions

	// target dimension
	if t, ok := dims["target"]; ok {
		targets := splitComma(t)
		if !containsStr(targets, opts.Target) {
			return "target:" + t
		}
	}

	// role dimension
	if r, ok := dims["role"]; ok {
		roles := splitComma(r)
		if opts.Role == "" || !containsStr(roles, opts.Role) {
			return "role:" + r
		}
	}

	// tag dimension
	if tg, ok := dims["tag"]; ok {
		tags := splitComma(tg)
		if len(opts.Tags) == 0 {
			return "tag:" + tg
		}
		matched := false
		for _, t := range tags {
			if containsStr(opts.Tags, t) {
				matched = true
				break
			}
		}
		if !matched {
			return "tag:" + tg
		}
	}

	return ""
}

func buildManifest(target, outputPath string, sections []schema.Section) *manifest.Manifest {
	m := &manifest.Manifest{
		SchemaVersion: 1,
		GeneratedAt:   time.Now().UTC().Format(time.RFC3339),
		Target:        target,
		OutputPath:    outputPath,
	}
	for _, s := range sections {
		m.Chunks = append(m.Chunks, manifest.Chunk{
			SourceFile: s.SourceFile,
			Section:    s.Heading,
		})
	}
	return m
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if strings.TrimSpace(v) == s {
			return true
		}
	}
	return false
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	for i, p := range parts {
		parts[i] = strings.TrimSpace(p)
	}
	return parts
}
