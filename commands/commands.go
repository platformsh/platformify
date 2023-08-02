package commands

import (
	"context"
)

// Execute executes the Platformify command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the PlatformifyCmd.
func Execute() error {
	return PlatformifyCmd.ExecuteContext(context.Background())
}
