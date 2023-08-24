package commands

import (
	"context"

	"github.com/platformsh/platformify/vendorization"
)

// Execute executes the ify command and sets flags appropriately.
func Execute(assets *vendorization.VendorAssets) error {
	return NewPlatformifyCmd(assets).ExecuteContext(vendorization.WithVendorAssets(context.Background(), assets))
}
