package commands

import (
	"context"
)

// Execute executes the Platformify command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the PlatformifyCmd.
func Execute(flavor string) error {
	return PlatformifyCmd.ExecuteContext(context.WithValue(context.Background(), FlavorKey, flavor))
}
