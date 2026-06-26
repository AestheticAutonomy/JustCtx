package differ

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AestheticAutonomy/justctx/internal/generator"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

type DiffOpts struct {
	Root   string
	Target string
	Role   string
	Tags   []string
}

// Diff compares what gen would produce against what's on disk.
func Diff(opts DiffOpts) (*schema.DiffResult, error) {
	results, err := generator.Generate(generator.GenOpts{
		Root:   opts.Root,
		Target: opts.Target,
		Role:   opts.Role,
		Tags:   opts.Tags,
		DryRun: true,
	})
	if err != nil {
		return nil, fmt.Errorf("generating: %w", err)
	}

	res := &schema.DiffResult{
		Envelope: schema.Envelope{
			SchemaVersion: 1,
			Command:       "diff",
			CWD:           opts.Root,
		},
		InSync:  true,
		Changes: []schema.Change{},
	}

	for _, genRes := range results {
		absPath := filepath.Join(opts.Root, genRes.OutputPath)
		actualBytes, err := os.ReadFile(absPath)
		var actual string
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("reading %s: %w", genRes.OutputPath, err)
			}
			// File doesn't exist — treat as empty
			actual = ""
		} else {
			actual = string(actualBytes)
		}

		expected := genRes.Content
		if actual != expected {
			res.InSync = false
			changeType := "modified"
			if actual == "" {
				changeType = "added"
			}
			res.Changes = append(res.Changes, schema.Change{
				Type:   changeType,
				Before: actual,
				After:  expected,
			})
		}
	}

	return res, nil
}

// FormatDiff produces a human-readable diff summary.
func FormatDiff(res *schema.DiffResult) string {
	if res.InSync {
		return "no drift detected\n"
	}
	var sb strings.Builder
	for _, c := range res.Changes {
		sb.WriteString(fmt.Sprintf("--- on disk\n+++ would generate\n"))
		beforeLines := strings.Split(c.Before, "\n")
		afterLines := strings.Split(c.After, "\n")
		// Simple line diff: mark added/removed
		maxLen := len(beforeLines)
		if len(afterLines) > maxLen {
			maxLen = len(afterLines)
		}
		for i := 0; i < maxLen; i++ {
			var before, after string
			if i < len(beforeLines) {
				before = beforeLines[i]
			}
			if i < len(afterLines) {
				after = afterLines[i]
			}
			if before == after {
				sb.WriteString(" " + before + "\n")
			} else {
				if i < len(beforeLines) {
					sb.WriteString("-" + before + "\n")
				}
				if i < len(afterLines) {
					sb.WriteString("+" + after + "\n")
				}
			}
		}
	}
	return sb.String()
}
