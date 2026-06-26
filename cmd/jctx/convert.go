package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	convertFrom     string
	convertTo       string
	convertType     string
	convertAllTypes bool
	convertDryRun   bool
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert guidelines directly from one provider to another",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("convert command is not implemented")
	},
}

func init() {
	convertCmd.Flags().StringVar(&convertFrom, "from", "", "source provider")
	convertCmd.Flags().StringVar(&convertTo, "to", "", "target provider")
	convertCmd.Flags().StringVar(&convertType, "type", "rules", "which type to convert")
	convertCmd.Flags().BoolVar(&convertAllTypes, "all-types", false, "convert all supported types")
	convertCmd.Flags().BoolVar(&convertDryRun, "dry-run", false, "dry run conversion")
	rootCmd.AddCommand(convertCmd)
}
