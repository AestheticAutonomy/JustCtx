package main

import (
	"fmt"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("diff command is not implemented")
	},
}

func init() {
	diffCmd.Flags().StringVar(&diffRole, "role", "", "scope to this role")
	diffCmd.Flags().StringSliceVar(&diffTags, "tag", []string{}, "scope to this tag (repeatable)")
	diffCmd.Flags().StringVar(&diffTarget, "target", "", "scope to one tool")
	rootCmd.AddCommand(diffCmd)
}
