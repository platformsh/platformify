package commands

import (
	"context"
	"fmt"
	"github.com/platformsh/platformify/internal/colors"
	"github.com/platformsh/platformify/internal/validator"
	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate your Platform.sh application config.",
	Long: "Check the application yaml configuration files are valid for deploying an application to Platform.sh.\n\n" +
		"This will check your git repository, and validate .platform.app.yaml, services.yaml and routes.yaml files.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var validator = validator.ConfigValidator{}
		var ctx = context.Background()
		ctx = colors.ToContext(
			ctx,
			cmd.OutOrStderr(),
			cmd.ErrOrStderr(),
		)
		err := validator.Validate(ctx)
		if err != nil {
			fmt.Fprintln(cmd.ErrOrStderr(), colors.Colorize(colors.ErrorCode, err.Error()))
			return fmt.Errorf("validation failed: #{err}")
		}

		return nil
	},
}
