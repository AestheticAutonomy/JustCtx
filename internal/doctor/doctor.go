package doctor

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AestheticAutonomy/justctx/internal/manifest"
	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type Check struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Detail string `json:"detail,omitempty"`
}

type Result struct {
	Checks  []Check `json:"checks"`
	Passed  int     `json:"passed"`
	Failed  int     `json:"failed"`
	AllPass bool    `json:"all_pass"`
}

// Run runs all doctor checks for the given root. providerFilter limits provider checks to one
// provider if non-empty.
func Run(root, providerFilter string) (*Result, error) {
	var checks []Check

	// 1. .jctx/ exists
	jctxDir := filepath.Join(root, ".jctx")
	if _, err := os.Stat(jctxDir); err == nil {
		checks = append(checks, Check{Name: ".jctx/ exists", Pass: true})
	} else {
		checks = append(checks, Check{Name: ".jctx/ exists", Pass: false, Detail: "directory not found"})
	}

	// 2. config.json valid JSON (if present)
	configPath := filepath.Join(jctxDir, "config.json")
	if data, err := os.ReadFile(configPath); err == nil {
		var v interface{}
		if jsonErr := json.Unmarshal(data, &v); jsonErr != nil {
			checks = append(checks, Check{Name: "config.json valid", Pass: false, Detail: jsonErr.Error()})
		} else {
			checks = append(checks, Check{Name: "config.json valid", Pass: true})
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		checks = append(checks, Check{Name: "config.json valid", Pass: false, Detail: err.Error()})
	}

	// 3. Provider FindFiles checks
	allProviders := providers.All()
	for _, p := range allProviders {
		if providerFilter != "" && p.Name() != providerFilter {
			continue
		}
		name := fmt.Sprintf("provider %s: FindFiles", p.Name())
		_, err := p.FindFiles(root, schema.TypeRules)
		if err != nil {
			checks = append(checks, Check{Name: name, Pass: false, Detail: err.Error()})
		} else {
			checks = append(checks, Check{Name: name, Pass: true})
		}
	}

	// 4. Manifest output files exist
	manifests, err := manifest.ListManifests(root)
	if err != nil {
		return nil, fmt.Errorf("listing manifests: %w", err)
	}
	for _, m := range manifests {
		absPath := filepath.Join(root, m.OutputPath)
		name := fmt.Sprintf("manifest target exists: %s", m.OutputPath)
		if _, err := os.Stat(absPath); err == nil {
			checks = append(checks, Check{Name: name, Pass: true})
		} else {
			checks = append(checks, Check{Name: name, Pass: false, Detail: "output file missing"})
		}
	}

	res := &Result{Checks: checks}
	for _, c := range checks {
		if c.Pass {
			res.Passed++
		} else {
			res.Failed++
		}
	}
	res.AllPass = res.Failed == 0
	return res, nil
}
