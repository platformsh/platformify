package main

import (
	"os"

	_ "github.com/golang/mock/mockgen/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/platformsh/platformify/commands"
	"github.com/platformsh/platformify/vendorization"
)

func main() {
	cobra.OnInitialize(func() {
		viper.SetEnvPrefix("PLATFORMIFY")
		viper.AutomaticEnv()
	})

	assets := &vendorization.VendorAssets{
		Binary:       "platform",
		ConfigFlavor: "platform",
		DocsBaseURL:  "https://docs.platform.sh/",
		EnvPrefix:    "PLATFORM",
		ServiceName:  "Platform.sh",
		Use:          "platformify",
	}
	err := commands.Execute(assets)
	if err != nil {
		os.Exit(1)
	}
}
