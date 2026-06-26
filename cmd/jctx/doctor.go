package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor [provider]",
	Short: "Validate setup — config, imports, provider files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("doctor command is not implemented")
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
