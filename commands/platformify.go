package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/platformifiers"
)

// PlatformifyCmd represents the base Platformify command when called without any subcommands
var PlatformifyCmd = &cobra.Command{
	Use:   "platformify",
	Short: "Platfomrify your application, and deploy it to the Platform.sh",
	Long: `Platformify your application, creating all the needed files
for it to be deployed to Platform.sh.

This will create the needed YAML files for both your application and your
services, choosing from a variety of stacks or simple runtimes.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var data []byte
		if len(args) == 0 || args[0] == "-" {
			fin := cmd.InOrStdin()
			var err error
			if data, err = io.ReadAll(fin); err != nil {
				return fmt.Errorf("could not read from stdin: %w", err)
			}
		} else {
			fin, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("could not open file %s: %w", args[0], err)
			}
			defer fin.Close()
			if data, err = io.ReadAll(fin); err != nil {
				return fmt.Errorf("could not read from file %s: %w", args[0], err)
			}
		}
		input := &platformifiers.UserInput{}
		if err := json.Unmarshal(data, input); err != nil {
			return fmt.Errorf("could not unmarshal json: %w", err)
		}
		var pfier platformifiers.Platformifier
		pfier, err := platformifiers.NewPlatformifier(input)
		if err != nil {
			return fmt.Errorf("creating platformifier failed: %s", err)
		}

		if err := pfier.Platformify(cmd.Context()); err != nil {
			return fmt.Errorf("could not platformify project: %w", err)
		}

		return nil
	},
}
