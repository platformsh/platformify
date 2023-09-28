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
		viper.SetEnvPrefix("UPSUNIFY")
		viper.AutomaticEnv() // read in environment variables that match
	})

	assets := &vendorization.VendorAssets{
		Binary:       "upsun",
		ConfigFlavor: "upsun",
		DocsBaseURL:  "https://docs.upsun.com",
		EnvPrefix:    "PLATFORM",
		ServiceName:  "Upsun",
		Use:          "upsunify",
	}

	err := commands.Execute(assets)
	if err != nil {
		os.Exit(1)
	}
}
