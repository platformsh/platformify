package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/validator"
	"github.com/platformsh/platformify/vendorization"
)

func NewValidateCommand(assets *vendorization.VendorAssets) *cobra.Command {
	return &cobra.Command{
		Use: "validate",
		Short: fmt.Sprintf(
			"Validate your %s config.",
			assets.ServiceName,
		),
		Long: fmt.Sprintf(
			"Check the application yaml configuration files are valid for deploying an application to %s.\n\n"+
				"This will check your git repository and validate your files.",
			assets.ServiceName,
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Fprintln(
					cmd.ErrOrStderr(),
					colors.Colorize(
						colors.ErrorCode,
						fmt.Sprintf("%s config validation failed: %s", assets.ServiceName, err.Error()),
					),
				)
				return err
			}

			if err = validator.ValidateConfig(os.DirFS(cwd), assets.ConfigFlavor); err != nil {
				fmt.Fprintf(
					cmd.ErrOrStderr(),
					colors.Colorize(
						colors.ErrorCode,
						"%s config validation failed: %s\n",
					),
					assets.ServiceName,
					err.Error(),
				)
				return err
			}

			fmt.Fprintf(
				cmd.ErrOrStderr(),
				colors.Colorize(colors.BrandCode, "Your %s application config is valid!\n"),
				assets.ServiceName,
			)
			return nil
		},
	}
}
