package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AestheticAutonomy/justctx/internal/generator"
	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
	"github.com/spf13/cobra"
)

var (
	updateDryRun bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Regenerate all targets using declared defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		defaults, err := loadConfigDefaults(cwd)
		if err != nil {
			return fmt.Errorf("reading config: %w", err)
		}
		if defaults == nil {
			return fmt.Errorf("jctx update requires .jctx/config.json with defaults set")
		}

		var targets []string
		if defaults.Target != "" {
			targets = []string{defaults.Target}
		} else {
			for _, p := range providers.All() {
				targets = append(targets, p.Name())
			}
		}

		res := schema.UpdateResult{
			Envelope: schema.Envelope{
				SchemaVersion: 1,
				Command:       "update",
				CWD:           cwd,
			},
		}

		for _, target := range targets {
			opts := generator.GenOpts{
				Root:   cwd,
				Target: target,
				Role:   defaults.Role,
				Tags:   defaults.Tags,
				DryRun: updateDryRun,
			}

			results, err := generator.Generate(opts)
			if err != nil {
				return fmt.Errorf("%s: %w", target, err)
			}

			for _, r := range results {
				if jsonFlag {
					res.TargetsUpdated = append(res.TargetsUpdated, r.OutputPath)
				} else {
					if updateDryRun {
						fmt.Printf("(dry run) %s\n", r.OutputPath)
					} else {
						fmt.Println(r.OutputPath)
					}
				}
			}
		}

		if jsonFlag {
			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		}

		return nil
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "print files that would be written without writing")
	rootCmd.AddCommand(updateCmd)
}
