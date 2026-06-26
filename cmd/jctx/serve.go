package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch local web UI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("serve command is not implemented")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
