package commands

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute executes the Platformify command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the PlatformifyCmd.
func Execute() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return PlatformifyCmd.ExecuteContext(ctx)
}

func init() {
	cobra.OnInitialize(initConfig)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	PlatformifyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("PLATFORMIFY")
	viper.AutomaticEnv() // read in environment variables that match
}
