package main

import (
	"os"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/platformsh/platformify/commands"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("PLATFORMIFY")
	viper.AutomaticEnv() // read in environment variables that match
}

func main() {
	cobra.OnInitialize(initConfig)
	commands.PlatformifyCmd.AddCommand(commands.ValidateCmd)
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}
