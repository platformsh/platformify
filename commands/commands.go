package commands

import (
	"context"

	"github.com/platformsh/platformify/vendorization"
)

// Execute executes the ify command and sets flags appropriately.
func Execute(assets *vendorization.VendorAssets) error {
	rootCmd := NewPlatformifyCmd(assets)
	validateCmd := NewValidateCommand(assets)
	rootCmd.AddCommand(validateCmd)
	return rootCmd.ExecuteContext(vendorization.WithVendorAssets(context.Background(), assets))
}
