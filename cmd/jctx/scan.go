package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	scanTarget   string
	scanNoGlobal bool
	scanBottomUp bool
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Show assembled guidelines for default target from cwd",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("scan command is not implemented")
	},
}

func init() {
	scanCmd.Flags().StringVar(&scanTarget, "target", "", "which tool's guidelines to show")
	scanCmd.Flags().BoolVar(&scanNoGlobal, "no-global", false, "skip ~/.claude/CLAUDE.md")
	scanCmd.Flags().BoolVar(&scanBottomUp, "bottom-up", false, "walk from cwd upward instead of top-down (non-standard)")
	rootCmd.AddCommand(scanCmd)
}
