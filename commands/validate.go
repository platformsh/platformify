package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/validator"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your Platform.sh application config.",
	Long: "Check the application yaml configuration files are valid for deploying an application to Platform.sh.\n\n" +
		"This will check your git repository, and validate .platform.app.yaml, services.yaml and routes.yaml files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(
				cmd.ErrOrStderr(),
				colors.Colorize(
					colors.ErrorCode,
					fmt.Sprintf("Platform.sh config validation failed: %s", err.Error()),
				),
			)
			return err
		}

		if err = validator.ValidateConfig(cwd, cmd.Context().Value(FlavorKey).(string)); err != nil {
			fmt.Fprintln(
				cmd.ErrOrStderr(),
				colors.Colorize(
					colors.ErrorCode,
					fmt.Sprintf("Platform.sh config validation failed: %s", err.Error()),
				),
			)
			return err
		}

		fmt.Fprintln(cmd.ErrOrStderr(), colors.Colorize(colors.BrandCode, "Your Platform.sh application config is valid!"))
		return nil
	},
}
