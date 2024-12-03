package vendorization

import (
	"context"
	"fmt"
)

type vendorAssetsKey string

var key vendorAssetsKey = "vendorAssets"

type Docs struct {
	AppReference   string
	GettingStarted string
	Hooks          string
	PHP            string
	Routes         string
	Services       string
	SymfonyCLI     string
	TimeZone       string
	Variables      string
}

type VendorAssets struct {
	Binary       string
	ConfigFlavor string
	DocsBaseURL  string
	EnvPrefix    string
	ServiceName  string
	Use          string
}

func (va *VendorAssets) ProprietaryFiles() []string {
	if va.ConfigFlavor == "upsun" {
		return []string{
			".environment",
			".upsun/config.yaml",
		}
	}

	return []string{
		".environment",
		".platform.app.yaml",
		".platform/services.yaml",
		".platform/routes.yaml",
		".platform/applications.yaml",
	}
}

func (va *VendorAssets) Docs() *Docs {
	return &Docs{
		AppReference:   fmt.Sprintf("%s/create-apps/app-reference/single-runtime-image.html", va.DocsBaseURL),
		GettingStarted: fmt.Sprintf("%s/guides/symfony/get-started.html", va.DocsBaseURL),
		Hooks:          fmt.Sprintf("%s/create-apps/hooks/hooks-comparison.html", va.DocsBaseURL),
		PHP:            fmt.Sprintf("%s/languages/php.html", va.DocsBaseURL),
		Routes:         fmt.Sprintf("%s/define-routes.html", va.DocsBaseURL),
		Services:       fmt.Sprintf("%s/add-services.html", va.DocsBaseURL),
		SymfonyCLI:     fmt.Sprintf("%s/guides/symfony/get-started.html#symfony-cli-tipsl", va.DocsBaseURL),
		TimeZone:       fmt.Sprintf("%s/create-apps/timezone.html", va.DocsBaseURL),
		Variables:      fmt.Sprintf("%s/development/variables/use-variables.html#use-platformsh-provided-variables", va.DocsBaseURL),
	}
}

func defaults() *VendorAssets {
	// Return all values as DEFAULT VALUE key
	return &VendorAssets{
		Binary:       "DEFAULT VALUE BINARY",
		ConfigFlavor: "DEFAULT VALUE CONFIGFLAVOR",
		DocsBaseURL:  "DEFAULT VALUE DOCS BASE URL",
		EnvPrefix:    "DEFAULT VALUE ENVPREFIX",
		ServiceName:  "DEFAULT VALUE SERVICENAME",
		Use:          "DEFAULT VALUE USE",
	}
}

func FromContext(ctx context.Context) (*VendorAssets, bool) {
	assets, ok := ctx.Value(key).(*VendorAssets)
	if !ok {
		return defaults(), false
	}

	return assets, ok
}

func WithVendorAssets(ctx context.Context, assets *VendorAssets) context.Context {
	return context.WithValue(ctx, key, assets)
}
