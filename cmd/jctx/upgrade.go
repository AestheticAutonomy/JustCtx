package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	upgradeCheck bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Update jctx to the latest version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("upgrade command is not implemented")
	},
}

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeCheck, "check", false, "check for new version only, no install")
	rootCmd.AddCommand(upgradeCmd)
}
