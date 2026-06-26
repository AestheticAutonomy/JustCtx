package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/AestheticAutonomy/justctx/internal/generator"
	"github.com/AestheticAutonomy/justctx/internal/providers"
	"github.com/spf13/cobra"
)

var (
	genRole     string
	genTags     []string
	genAll      bool
	genDryRun   bool
	genAnnotate bool
)

var genCmd = &cobra.Command{
	Use:   "gen [target]",
	Short: "Generate guidelines for a tool from .jctx source",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		var targets []string
		if genAll {
			for _, p := range providers.All() {
				targets = append(targets, p.Name())
			}
		} else {
			if len(args) == 0 {
				return fmt.Errorf("target required (or use --all)")
			}
			targets = []string{args[0]}
		}

		for _, target := range targets {
			opts := generator.GenOpts{
				Root:   cwd,
				Target: target,
				Role:   genRole,
				Tags:   genTags,
				DryRun: genDryRun,
			}

			results, err := generator.Generate(opts)
			if err != nil {
				return fmt.Errorf("%s: %w", target, err)
			}

			if jsonFlag {
				for _, r := range results {
					data, err := json.MarshalIndent(r, "", "  ")
					if err != nil {
						return err
					}
					fmt.Println(string(data))
				}
				continue
			}

			for _, r := range results {
				if genDryRun {
					fmt.Printf("(dry run) %s\n", r.OutputPath)
				} else {
					fmt.Println(r.OutputPath)
				}
			}
		}

		return nil
	},
}

func runGen(cwd string, targets []string, role string, tags []string, dryRun bool, outputJSON bool, out io.Writer) error {
	for _, target := range targets {
		opts := generator.GenOpts{
			Root:   cwd,
			Target: target,
			Role:   role,
			Tags:   tags,
			DryRun: dryRun,
		}
		results, err := generator.Generate(opts)
		if err != nil {
			return fmt.Errorf("%s: %w", target, err)
		}
		if outputJSON {
			for _, r := range results {
				data, err := json.MarshalIndent(r, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(out, string(data))
			}
			continue
		}
		for _, r := range results {
			if dryRun {
				fmt.Fprintf(out, "(dry run) %s\n", r.OutputPath)
			} else {
				fmt.Fprintln(out, r.OutputPath)
			}
		}
	}
	return nil
}

func allTargets() []string {
	var targets []string
	for _, p := range providers.All() {
		targets = append(targets, p.Name())
	}
	return targets
}

func init() {
	genCmd.Flags().StringVar(&genRole, "role", "", "filter to this role")
	genCmd.Flags().StringSliceVar(&genTags, "tag", []string{}, "filter to this tag (repeatable)")
	genCmd.Flags().BoolVar(&genAll, "all", false, "generate all targets declared in source")
	genCmd.Flags().BoolVar(&genDryRun, "dry-run", false, "print output, no writes")
	genCmd.Flags().BoolVar(&genAnnotate, "annotate", false, "include generation metadata as comments in output")
	rootCmd.AddCommand(genCmd)
}
