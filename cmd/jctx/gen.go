package main

import (
	"fmt"

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
	Short: "Generate guidelines for a tool from .justctx source",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gen command is not implemented")
	},
}

func init() {
	genCmd.Flags().StringVar(&genRole, "role", "", "filter to this role")
	genCmd.Flags().StringSliceVar(&genTags, "tag", []string{}, "filter to this tag (repeatable)")
	genCmd.Flags().BoolVar(&genAll, "all", false, "generate all targets declared in source")
	genCmd.Flags().BoolVar(&genDryRun, "dry-run", false, "print output, no writes")
	genCmd.Flags().BoolVar(&genAnnotate, "annotate", false, "include generation metadata as comments in output")
	rootCmd.AddCommand(genCmd)
}
