package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/AestheticAutonomy/justctx/internal/doctor"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor [provider]",
	Short: "Validate setup — config, imports, provider files",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}

		providerFilter := ""
		if len(args) > 0 {
			providerFilter = args[0]
		}

		res, err := doctor.Run(cwd, providerFilter)
		if err != nil {
			return err
		}

		if jsonFlag {
			data, err := json.MarshalIndent(res, "", "  ")
			if err != nil {
				return err
			}
			fmt.Println(string(data))
		} else {
			for _, c := range res.Checks {
				if c.Pass {
					fmt.Printf("OK   %s\n", c.Name)
				} else {
					fmt.Printf("FAIL %s: %s\n", c.Name, c.Detail)
				}
			}
			fmt.Printf("\n%d checks passed, %d failed\n", res.Passed, res.Failed)
		}

		if !res.AllPass {
			os.Exit(1)
		}
		return nil
	},
}

func runDoctorCmd(cwd, providerFilter string, outputJSON bool, out io.Writer) (bool, error) {
	res, err := doctor.Run(cwd, providerFilter)
	if err != nil {
		return false, err
	}
	if outputJSON {
		data, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			return false, err
		}
		fmt.Fprintln(out, string(data))
	} else {
		for _, c := range res.Checks {
			if c.Pass {
				fmt.Fprintf(out, "OK   %s\n", c.Name)
			} else {
				fmt.Fprintf(out, "FAIL %s: %s\n", c.Name, c.Detail)
			}
		}
		fmt.Fprintf(out, "\n%d checks passed, %d failed\n", res.Passed, res.Failed)
	}
	return res.AllPass, nil
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
