package colors

import (
	"fmt"

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

// Colorize the given string with the given color code
func Colorize(code, str string) string {
	return fmt.Sprintf("%s%s%s", code, str, resetCode)
}

func SetColors(brand, accent, secondary, warning, errorColor, defaultColor string) {
	Brand = brand
	Accent = accent
	Secondary = secondary
	Warning = warning
	Error = errorColor
	Default = defaultColor

	BrandCode = ansi.ColorCode(Brand)
	AccentCode = ansi.ColorCode(Accent)
	SecondaryCode = ansi.ColorCode(Secondary)
	WarningCode = ansi.ColorCode(Warning)
	ErrorCode = ansi.ColorCode(Error)
	DefaultCode = ansi.ColorCode(Default)
}
