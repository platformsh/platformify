package main

import (
	"os"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/platformsh/platformify/commands"
)

func main() {
	cobra.OnInitialize(initConfig)

	commands.PlatformifyCmd.AddCommand(commands.ValidateCmd)
	err := commands.Execute("upsun")
	if err != nil {
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("UPSUNIFY")
	viper.AutomaticEnv() // read in environment variables that match
}
