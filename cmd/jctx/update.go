package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	updateDryRun bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Regenerate all targets using declared defaults",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("update command is not implemented")
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "dry run update")
	rootCmd.AddCommand(updateCmd)
}
