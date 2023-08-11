package vendorization

import "context"

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
	Binary           string
	ConfigFlavor     string
	EnvPrefix        string
	ProprietaryFiles []string
	ServiceName      string

	Docs Docs
}

func defaults() *VendorAssets {
	// Return all values as DEFAULT VALUE key
	return &VendorAssets{
		Binary:           "DEFAULT VALUE BINARY",
		ConfigFlavor:     "DEFAULT VALUE CONFIGFLAVOR",
		EnvPrefix:        "DEFAULT VALUE ENVPREFIX",
		ProprietaryFiles: []string{"DEFAULT VALUE PROPRIETARYFILES"},
		ServiceName:      "DEFAULT VALUE SERVICENAME",
		Docs: Docs{
			AppReference: "DEFAULT VALUE DOCS APPREFERENCE",

			GettingStarted: "DEFAULT VALUE DOCS GETTINGSTARTED",
			Hooks:          "DEFAULT VALUE DOCS HOOKS",
			PHP:            "DEFAULT VALUE DOCS PHP",
			Routes:         "DEFAULT VALUE DOCS ROUTES",
			Services:       "DEFAULT VALUE DOCS SERVICES",
			SymfonyCLI:     "DEFAULT VALUE DOCS SYMFONYCLI",
			TimeZone:       "DEFAULT VALUE DOCS TIMEZONE",
			Variables:      "DEFAULT VALUE DOCS VARIABLES",
		},
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
