package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/platformsh/platformify/vendorization"
)

// Execute executes the ify command and sets flags appropriately.
func Execute(assets *vendorization.VendorAssets) error {
	rootCmd := NewPlatformifyCmd(assets)
	validateCmd := NewValidateCommand(assets)
	rootCmd.AddCommand(validateCmd)

	rootCmd.PersistentFlags().Bool("no-interaction", false, "Disable interactive prompts")
	viper.BindPFlag("no-interaction", rootCmd.PersistentFlags().Lookup("no-interaction"))

	rootCmd.PreRun = func(cmd *cobra.Command, _ []string) {
		ctx := context.WithValue(cmd.Context(), NoInteractionKey, viper.GetBool("no-interaction"))
		cmd.SetContext(ctx)
	}

	return rootCmd.ExecuteContext(vendorization.WithVendorAssets(context.Background(), assets))
}
