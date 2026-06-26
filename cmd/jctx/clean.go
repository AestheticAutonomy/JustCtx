package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cleanTarget string
	cleanDryRun bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove generated artifacts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("clean command is not implemented")
	},
}

func init() {
	cleanCmd.Flags().StringVar(&cleanTarget, "target", "", "scope to one provider")
	cleanCmd.Flags().BoolVar(&cleanDryRun, "dry-run", false, "dry run clean")
	rootCmd.AddCommand(cleanCmd)
}
