package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	initImport    bool
	initOverwrite bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Scaffold a blank .jctx/ in the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init command is not implemented")
	},
}

func init() {
	initCmd.Flags().BoolVar(&initImport, "import", false, "scan existing guidelines, convert to .jctx source")
	initCmd.Flags().BoolVar(&initOverwrite, "overwrite", false, "overwrite existing .jctx/ when used with --import")
	rootCmd.AddCommand(initCmd)
}
