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
		EnvPrefix:    "UPSUN",
		ServiceName:  "Upsun",
		ProprietaryFiles: []string{
			".environment",
			".upsun/config.yaml",
		},
		Docs: vendorization.Docs{
			AppReference:   "TODO AppReference",
			GettingStarted: "TODO GettingStarted",
			Hooks:          "TODO Hooks",
			PHP:            "TODO PHP",
			Routes:         "TODO Routes",
			Services:       "TODO Services",
			SymfonyCLI:     "TODO SymfonyCLI",
			TimeZone:       "TODO TimeZone",
			Variables:      "TODO Variables",
		},
	}

	err := commands.Execute(assets)
	if err != nil {
		os.Exit(1)
	}
}
