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

	//nolint:lll
	assets := &vendorization.VendorAssets{
		Binary:       "platform",
		ConfigFlavor: "platform",
		EnvPrefix:    "PLATFORM",
		ProprietaryFiles: []string{
			".environment",
			".platform.app.yaml",
			".platform/services.yaml",
			".platform/routes.yaml",
			".platform/applications.yaml",
		},
		ServiceName: "Platform.sh",
		Docs: vendorization.Docs{
			AppReference:   "https://docs.platform.sh/create-apps/app-reference.html",
			GettingStarted: "https://docs.platform.sh/guides/symfony/get-started.html",
			Hooks:          "https://docs.platform.sh/create-apps/hooks/hooks-comparison.html",
			PHP:            "https://docs.platform.sh/languages/php.html",
			Routes:         "https://docs.platform.sh/define-routes.html",
			Services:       "https://docs.platform.sh/add-services.html",
			SymfonyCLI:     "https://docs.platform.sh/guides/symfony/get-started.html#symfony-cli-tips",
			TimeZone:       "https://docs.platform.sh/create-apps/timezone.html",
			Variables:      "https://docs.platform.sh/development/variables/use-variables.html#use-platformsh-provided-variables",
		},
	}
	err := commands.Execute(assets)
	if err != nil {
		os.Exit(1)
	}
}
