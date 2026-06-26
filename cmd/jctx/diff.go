package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/AestheticAutonomy/justctx/internal/differ"
	"github.com/spf13/cobra"
)

var (
	diffRole   string
	diffTags   []string
	diffTarget string
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Diff generated guidelines against what source would produce",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		if diffTarget == "" {
			return fmt.Errorf("--target required")
		}

		res, err := differ.Diff(differ.DiffOpts{
			Root:   cwd,
			Target: diffTarget,
			Role:   diffRole,
			Tags:   diffTags,
		})
		if err != nil {
			return err
		}

		if jsonFlag {
			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			fmt.Print(differ.FormatDiff(res))
		}

		if !res.InSync {
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	diffCmd.Flags().StringVar(&diffRole, "role", "", "scope to this role")
	diffCmd.Flags().StringSliceVar(&diffTags, "tag", []string{}, "scope to this tag (repeatable)")
	diffCmd.Flags().StringVar(&diffTarget, "target", "", "scope to one tool")
	rootCmd.AddCommand(diffCmd)
}
