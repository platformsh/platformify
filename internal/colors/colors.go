package colors

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/mgutz/ansi"
)

var (
	Brand     = "210"
	Accent    = "33"
	Secondary = "247"
	Warning   = "yellow"
	Error     = "red"
	Default   = "default"

	BrandCode     = ansi.ColorCode(Brand)
	AccentCode    = ansi.ColorCode(Accent)
	SecondaryCode = ansi.ColorCode(Secondary)
	WarningCode   = ansi.ColorCode(Warning)
	ErrorCode     = ansi.ColorCode(Error)
	DefaultCode   = ansi.ColorCode(Default)
	resetCode     = ansi.ColorCode("reset")
)

func init() {
	// If NO_COLOR environment variable is set and is not 0, reset colors to ""
	if os.Getenv("NO_COLOR") != "" && os.Getenv("NO_COLOR") != "0" {
		BrandCode = ""
		AccentCode = ""
		SecondaryCode = ""
		WarningCode = ""
		ErrorCode = ""
		DefaultCode = ""
		resetCode = ""
		core.DisableColor = true
	}
}

// Colorize the given string with the given color code
func Colorize(code, str string) string {
	return fmt.Sprintf("%s%s%s", code, str, resetCode)
}
