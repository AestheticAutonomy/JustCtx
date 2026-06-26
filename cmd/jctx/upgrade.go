package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/AestheticAutonomy/justctx/internal/upgrade"
	"github.com/spf13/cobra"
)

var (
	upgradeCheck bool
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Update jctx to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		latest, isNewer, err := upgrade.CheckLatest(http.DefaultClient, version)
		if err != nil {
			return fmt.Errorf("checking for updates: %w", err)
		}

		if jsonFlag {
			data, _ := json.MarshalIndent(map[string]interface{}{
				"current":    version,
				"latest":     latest,
				"is_newer":   isNewer,
				"up_to_date": !isNewer,
			}, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		if !isNewer {
			fmt.Printf("already up to date (jctx %s)\n", version)
			return nil
		}

		fmt.Printf("jctx %s is available (current: %s)\n", latest, version)

		if upgradeCheck {
			return nil
		}

		// Full download + replace not yet implemented
		fmt.Fprintln(os.Stderr, "binary replacement not yet implemented — download manually from GitHub releases")
		return nil
	},
}

func init() {
	upgradeCmd.Flags().BoolVar(&upgradeCheck, "check", false, "check for new version only, no install")
	rootCmd.AddCommand(upgradeCmd)
}
