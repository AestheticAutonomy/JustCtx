package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/AestheticAutonomy/justctx/internal/manifest"
	"github.com/AestheticAutonomy/justctx/pkg/schema"
	"github.com/spf13/cobra"
)

var (
	cleanTarget string
	cleanDryRun bool
	cleanAll    bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove generated artifacts",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		if cleanTarget == "" && !cleanAll {
			return fmt.Errorf("--target or --all required")
		}

		manifests, err := manifest.ListManifests(cwd)
		if err != nil {
			return err
		}

		var toClean []*manifest.Manifest
		for _, m := range manifests {
			if cleanAll || m.Target == cleanTarget {
				toClean = append(toClean, m)
			}
		}

		if len(toClean) == 0 {
			if cleanTarget != "" {
				fmt.Printf("nothing to clean for %s\n", cleanTarget)
			} else {
				fmt.Println("nothing to clean")
			}
			return nil
		}

		res := schema.CleanResult{
			Envelope: schema.Envelope{
				SchemaVersion: 1,
				Command:       "clean",
				CWD:           cwd,
			},
		}

		for _, m := range toClean {
			absPath := filepath.Join(cwd, m.OutputPath)

			if cleanDryRun {
				fmt.Printf("(dry run) %s\n", m.OutputPath)
				res.RemovedFiles = append(res.RemovedFiles, m.OutputPath)
				continue
			}

			if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("removing %s: %w", m.OutputPath, err)
			}
			if err := manifest.Delete(cwd, m.OutputPath); err != nil {
				return fmt.Errorf("removing manifest for %s: %w", m.OutputPath, err)
			}
			fmt.Println(m.OutputPath)
			res.RemovedFiles = append(res.RemovedFiles, m.OutputPath)
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

func runClean(cwd, target string, all, dryRun, outputJSON bool, out io.Writer) error {
	if target == "" && !all {
		return fmt.Errorf("--target or --all required")
	}

	manifests, err := manifest.ListManifests(cwd)
	if err != nil {
		return err
	}

	var toClean []*manifest.Manifest
	for _, m := range manifests {
		if all || m.Target == target {
			toClean = append(toClean, m)
		}
	}

	if len(toClean) == 0 {
		if target != "" {
			fmt.Fprintf(out, "nothing to clean for %s\n", target)
		} else {
			fmt.Fprintln(out, "nothing to clean")
		}
		return nil
	}

	res := schema.CleanResult{
		Envelope: schema.Envelope{SchemaVersion: 1, Command: "clean", CWD: cwd},
	}

	for _, m := range toClean {
		absPath := filepath.Join(cwd, m.OutputPath)
		if dryRun {
			fmt.Fprintf(out, "(dry run) %s\n", m.OutputPath)
			res.RemovedFiles = append(res.RemovedFiles, m.OutputPath)
			continue
		}
		if err := os.Remove(absPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("removing %s: %w", m.OutputPath, err)
		}
		if err := manifest.Delete(cwd, m.OutputPath); err != nil {
			return fmt.Errorf("removing manifest for %s: %w", m.OutputPath, err)
		}
		fmt.Fprintln(out, m.OutputPath)
		res.RemovedFiles = append(res.RemovedFiles, m.OutputPath)
	}

	if outputJSON {
		data, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return err
		}
		fmt.Fprintln(out, string(data))
	}
	return nil
}

func init() {
	cleanCmd.Flags().StringVar(&cleanTarget, "target", "", "scope to one provider")
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "list files that would be removed without removing")
	cleanCmd.Flags().BoolVar(&cleanAll, "all", false, "clean all providers")
	rootCmd.AddCommand(cleanCmd)
}
