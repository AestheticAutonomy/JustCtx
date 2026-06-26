package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

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
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}

		err = runInit(cwd, cmd.OutOrStdout(), jsonFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func runInit(cwd string, out io.Writer, outputJSON bool) error {
	jctxDir := filepath.Join(cwd, ".jctx")
	if _, err := os.Stat(jctxDir); err == nil {
		fmt.Fprintln(out, ".jctx/ already exists")
		return nil
	}

	subdirs := []string{
		"rules",
		"hooks",
		"mcp",
		"commands",
		"skills",
		"ignores",
	}

	var created []string

	if err := os.Mkdir(jctxDir, 0755); err != nil {
		return fmt.Errorf("creating .jctx directory: %w", err)
	}
	created = append(created, ".jctx/")

	for _, sub := range subdirs {
		path := filepath.Join(jctxDir, sub)
		if err := os.Mkdir(path, 0755); err != nil {
			return fmt.Errorf("creating .jctx/%s directory: %w", sub, err)
		}
		created = append(created, fmt.Sprintf(".jctx/%s/", sub))
	}

	starterFile := filepath.Join(jctxDir, "rules", "main.md")
	starterContent := `---
target: [claude]
---
@@@ General Rules

Add your coding guidelines here.
`
	if err := os.WriteFile(starterFile, []byte(starterContent), 0644); err != nil {
		return fmt.Errorf("creating .jctx/rules/main.md: %w", err)
	}
	created = append(created, ".jctx/rules/main.md")

	if outputJSON {
		data, err := json.Marshal(created)
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Fprintln(out, string(data))
	} else {
		for _, item := range created {
			fmt.Fprintln(out, item)
		}
	}

	return nil
}

func init() {
	initCmd.Flags().BoolVar(&initImport, "import", false, "scan existing guidelines, convert to .jctx source")
	initCmd.Flags().BoolVar(&initOverwrite, "overwrite", false, "overwrite existing .jctx/ when used with --import")
	rootCmd.AddCommand(initCmd)
}
