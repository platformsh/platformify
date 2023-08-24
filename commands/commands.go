package commands

import (
	"context"

	"github.com/platformsh/platformify/vendorization"
)

// Execute executes the ify command and sets flags appropriately.
func Execute(assets *vendorization.VendorAssets) error {
	cmd := NewPlatformifyCmd(assets)
	validateCmd := NewValidateCommand(assets)
	cmd.AddCommand(validateCmd)
	return cmd.ExecuteContext(vendorization.WithVendorAssets(context.Background(), assets))
}
