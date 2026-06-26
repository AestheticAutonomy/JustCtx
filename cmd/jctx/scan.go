package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AestheticAutonomy/justctx/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	scanTarget   string
	scanNoGlobal bool
	scanBottomUp bool
)

type ConfigDefaults struct {
	Target string   `json:"target"`
	Role   string   `json:"role"`
	Tags   []string `json:"tags"`
}

type Config struct {
	SchemaVersion int            `json:"schema_version"`
	Defaults      ConfigDefaults `json:"defaults"`
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Show assembled guidelines for default target from cwd",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting working directory: %v\n", err)
			os.Exit(1)
		}

		target := scanTarget
		if target == "" {
			defaults, err := loadConfigDefaults(cwd)
			if err == nil && defaults != nil {
				target = defaults.Target
			}
		}

		if target == "" {
			fmt.Fprintln(os.Stderr, "specify --target or set a default in .jctx/config.json")
			os.Exit(1)
		}

		res, err := scanner.Scan(scanner.ScanOpts{
			Root:     cwd,
			Target:   target,
			NoGlobal: scanNoGlobal,
			BottomUp: scanBottomUp,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if len(res.Sources) == 0 {
			fmt.Fprintln(os.Stderr, "no provider files found")
			os.Exit(0)
		}

		if jsonFlag {
			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error marshaling JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(data))
			return
		}

		// Plain text pretty printing
		fmt.Printf("Sources (%d):\n", len(res.Sources))
		for _, src := range res.Sources {
			fmt.Printf("  [%s]\t%s\t(%d bytes)\n", src.Location, src.Path, src.Bytes)
		}
		fmt.Println()

		fmt.Println("Assembled context:")
		for _, chunk := range res.Assembled {
			var srcPath string
			for _, src := range res.Sources {
				if src.ID == chunk.SourceID {
					srcPath = src.Path
					break
				}
			}

			lineHeader := fmt.Sprintf("─── %s ", srcPath)
			if len(lineHeader) < 50 {
				lineHeader += strings.Repeat("─", 50-len(lineHeader))
			}
			fmt.Println(lineHeader)
			fmt.Print(chunk.Content)
			if !strings.HasSuffix(chunk.Content, "\n") && chunk.Content != "" {
				fmt.Println()
			}
		}
		fmt.Println()

		if len(res.Conflicts) > 0 {
			fmt.Printf("Conflicts (%d):\n", len(res.Conflicts))
			for _, c := range res.Conflicts {
				if c.Type == "duplicate_heading" {
					var names []string
					for _, sid := range c.SourceIDs {
						for _, src := range res.Sources {
							if src.ID == sid {
								names = append(names, filepath.Base(src.Path))
								break
							}
						}
					}
					fmt.Printf("  [duplicate_heading] \"%s\" in %s\n", c.Heading, strings.Join(names, " + "))
				} else if c.Type == "contradicting_imperative" {
					fmt.Printf("  [contradicting_imperative] %s\n", c.Heading)
				} else {
					fmt.Printf("  [%s] in %s\n", c.Type, strings.Join(c.SourceIDs, " + "))
				}
			}
		}
	},
}

func init() {
	scanCmd.Flags().StringVar(&scanTarget, "target", "", "which tool's guidelines to show")
	scanCmd.Flags().BoolVar(&scanNoGlobal, "no-global", false, "skip ~/.claude/CLAUDE.md")
	scanCmd.Flags().BoolVar(&scanBottomUp, "bottom-up", false, "walk from cwd upward instead of top-down (non-standard)")
	rootCmd.AddCommand(scanCmd)
}

func loadConfigDefaults(startDir string) (*ConfigDefaults, error) {
	dir, err := filepath.Abs(startDir)
	if err != nil {
		return nil, err
	}
	for {
		configPath := filepath.Join(dir, ".jctx", "config.json")
		if _, err := os.Stat(configPath); err == nil {
			data, err := os.ReadFile(configPath)
			if err != nil {
				return nil, err
			}
			var cfg Config
			if err := json.Unmarshal(data, &cfg); err != nil {
				return nil, err
			}
			return &cfg.Defaults, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return nil, nil
}
