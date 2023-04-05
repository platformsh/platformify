package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

var (
	// Colors. See Documentation: https://github.com/mgutz/ansi#style-format
	color1      = "33"  // DodgerBlue1
	color2      = "210" // LightCoral
	color3      = "247" // Grey62
	color4      = "red" // Red3
	color5      = "default"
	colorSchema = []*string{&color1, &color2, &color3, &color4, &color5}
)

func setColors(colors ...string) {
	if len(colors) > len(colorSchema) {
		colors = colors[:len(colorSchema)]
	}
	for i, color := range colors {
		*colorSchema[i] = color
	}

	survey.InputQuestionTemplate = fmt.Sprintf(`
{{- color "%[1]s"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
	{{- color "%[2]s"}}[{{.Answer}}]{{color "reset"}}{{"\n"}}
{{ else }}
	{{- if .Default }}{{color "%[3]s"}}[{{.Default}}] {{color "reset"}}{{ end }}
{{- end }}`, color1, color2, color3, color4, color5)

	//nolint:lll
	survey.SelectQuestionTemplate = fmt.Sprintf(`
{{- define "option"}}
    {{- if eq .SelectedIndex .CurrentIndex }}{{color "%[2]s" }}{{ .Config.Icons.SelectFocus.Text }} {{else}}{{color "%[5]s"}}  {{end}}
    {{- .CurrentOpt.Value}}
    {{- color "reset"}}
{{end}}
{{- color "%[1]s"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "%[2]s"}} [{{.Answer}}]{{color "reset"}}{{"\n"}}{{"\n"}}
{{- else}}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
  {{- color "%[3]s"}}Use arrows to move up and down, type to filter{{color "reset"}}{{"\n"}}
{{- end}}`, color1, color2, color3, color4, color5)

	survey.ConfirmQuestionTemplate = fmt.Sprintf(`
{{- color "%[1]s"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "%[2]s"}}[{{.Answer}}]{{color "reset"}}{{"\n"}}
{{- else }}
  {{- color "%[3]s"}}{{if .Default}}Press "y" for yes, "n" for no. {{else}}(y/N){{end}}{{color "reset"}}{{"\n"}}
{{- end}}`, color1, color2, color3, color4, color5)

	//nolint:lll
	survey.MultiSelectQuestionTemplate = fmt.Sprintf(`
{{- define "option"}}
    {{- if eq .SelectedIndex .CurrentIndex }}{{color "%[5]s" }}{{ .Config.Icons.SelectFocus.Text }}{{color "reset"}}{{else}} {{end}}
    {{- if index .Checked .CurrentOpt.Index }}{{color "%[2]s" }} {{ .Config.Icons.MarkedOption.Text }} {{else}}{{color "%[5]s" }} {{ .Config.Icons.UnmarkedOption.Text }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{- if index .Checked .CurrentOpt.Index }}{{color "%[2]s" }}{{else}}{{color "%[5]s" }}{{end}}{{- .CurrentOpt.Value}}{{color "reset"}}
{{end}}
{{- color "%[1]s"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "%[2]s"}} [{{.Answer}}]{{color "reset"}}{{"\n"}}{{"\n"}}
{{- else }}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
  {{- color "%[3]s"}}Use arrows to move, space to select, type to filter{{color "reset"}}{{"\n"}}
{{- end}}`, color1, color2, color3, color4, color5)
}
