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

	//nolint:lll
	assets := &vendorization.VendorAssets{
		Use:          "upsunify",
		Binary:       "upsun",
		ConfigFlavor: "upsun",
		EnvPrefix:    "UPSUN",
		ServiceName:  "Upsun",
		ProprietaryFiles: []string{
			".environment",
			".upsun/config.yaml",
		},
		Docs: vendorization.Docs{
			AppReference:   "https://docs.deployfriday.net/create-apps/app-reference.html",
			GettingStarted: "https://docs.deployfriday.net/guides/symfony/get-started.html",
			Hooks:          "https://docs.deployfriday.net/create-apps/hooks/hooks-comparison.html",
			PHP:            "https://docs.deployfriday.net/languages/php.html",
			Routes:         "https://docs.deployfriday.net/define-routes.html",
			Services:       "https://docs.deployfriday.net/add-services.html",
			SymfonyCLI:     "https://docs.deployfriday.net/guides/symfony/get-started.html#symfony-cli-tipsl",
			TimeZone:       "https://docs.deployfriday.net/create-apps/timezone.html",
			Variables:      "https://docs.deployfriday.net/development/variables/use-variables.html#use-platformsh-provided-variables",
		},
	}

	err := commands.Execute(assets)
	if err != nil {
		os.Exit(1)
	}
}
