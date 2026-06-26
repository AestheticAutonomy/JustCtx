package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
	"github.com/spf13/cobra"
)

var (
	convertFrom     string
	convertTo       string
	convertType     string
	convertAllTypes bool
	convertDryRun   bool
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert guidelines directly from one provider to another",
	RunE: func(cmd *cobra.Command, args []string) error {
		if convertFrom == "" {
			return fmt.Errorf("--from required")
		}

		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		p, err := providers.Get(convertFrom)
		if err != nil {
			return err
		}

		files, err := p.FindFiles(cwd, schema.TypeRules)
		if err != nil {
			return fmt.Errorf("finding files: %w", err)
		}

		if len(files) == 0 {
			fmt.Printf("no %s files found\n", convertFrom)
			return nil
		}

		var written []string

		for _, f := range files {
			sections, err := p.ParseRules(f)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", f, err)
			}
			if len(sections) == 0 {
				continue
			}

			// Build .jctx/rules/<slug>.md
			base := filepath.Base(f)
			ext := filepath.Ext(base)
			slug := strings.TrimSuffix(base, ext)

			outPath := filepath.Join(cwd, ".jctx", "rules", slug+".md")
			content := buildJctxSource(convertFrom, sections)

			if convertDryRun {
				fmt.Printf("(dry run) .jctx/rules/%s.md\n", slug)
				written = append(written, outPath)
				continue
			}

			if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
				return err
			}
			if err := os.WriteFile(outPath, []byte(content), 0644); err != nil {
				return err
			}
			fmt.Printf(".jctx/rules/%s.md\n", slug)
			written = append(written, outPath)
		}

		if jsonFlag {
			res := schema.ConvertResult{
				Envelope: schema.Envelope{
					SchemaVersion: 1,
					Command:       "convert",
					CWD:           cwd,
				},
				From: convertFrom,
				To:   convertTo,
				Type: schema.TypeRules,
			}
			_ = written
			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		}

		return nil
	},
}

// buildJctxSource renders parsed sections as a .jctx/rules/*.md file.
func buildJctxSource(target string, sections []schema.Section) string {
	var sb strings.Builder
	sb.WriteString("---\ntargets: [")
	sb.WriteString(target)
	sb.WriteString("]\n---\n")

	for _, s := range sections {
		sb.WriteString("\n")
		if s.Heading != "" {
			sb.WriteString("@@@ ")
			sb.WriteString(s.Heading)
			sb.WriteString("\n")
		} else {
			sb.WriteString("@@@ Content\n")
		}
		content := strings.TrimRight(s.Content, "\n")
		sb.WriteString(content)
		sb.WriteString("\n")
	}
	return sb.String()
}

func init() {
	convertCmd.Flags().StringVar(&convertFrom, "from", "", "source provider")
	convertCmd.Flags().StringVar(&convertTo, "to", "", "target provider")
	convertCmd.Flags().StringVar(&convertType, "type", "rules", "which type to convert")
	convertCmd.Flags().BoolVar(&convertAllTypes, "all-types", false, "convert all supported types")
	convertCmd.Flags().BoolVar(&convertDryRun, "dry-run", false, "dry run conversion")
	rootCmd.AddCommand(convertCmd)
}
