package main

import (
	"fmt"
	"os"

	_ "github.com/AestheticAutonomy/justctx/internal/providers/agents"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/antigravity"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/claude"
	_ "github.com/AestheticAutonomy/justctx/internal/providers/cursor"
	"github.com/spf13/cobra"
)

var (
	jsonFlag bool
	version  = "0.1.0"
)

var rootCmd = &cobra.Command{
	Use:     "jctx",
	Short:   "justctx is a build system for AI guidelines",
	Long:    `One source of truth for every AI instruction file in your stack.`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "output results in JSON format")
}
