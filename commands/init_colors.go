package commands

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"

	"github.com/platformsh/platformify/internal/colors"
)

func init() {
	survey.InputQuestionTemplate = fmt.Sprintf(`
{{- color "%[1]s"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
	{{- color "%[2]s"}}[{{.Answer}}]{{color "reset"}}{{"\n"}}
{{ else }}
	{{- if .Default }}{{color "%[3]s"}}[{{.Default}}] {{color "reset"}}{{ end }}
{{- end }}`, colors.Accent, colors.Brand, colors.Secondary, colors.Error, colors.Default)

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
{{ color "%[5]s"}}Use arrows to move up and down, type to filter{{color "reset"}}{{"\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
{{- end }}`, colors.Accent, colors.Brand, colors.Secondary, colors.Error, colors.Default)

	survey.ConfirmQuestionTemplate = fmt.Sprintf(`
{{- color "%[1]s"}}{{ .Message }} {{color "reset"}}
{{- if .Answer}}
  {{- color "%[2]s"}}[{{.Answer}}]{{color "reset"}}{{"\n"}}
{{- else }}
  {{- color "%[3]s"}}{{if .Default}}Press "y" for yes, "n" for no. {{else}}(y/N){{end}}{{color "reset"}}{{"\n"}}
{{- end }}`, colors.Accent, colors.Brand, colors.Secondary, colors.Error, colors.Default)

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
{{ color "%[5]s"}}Use arrows to move, space to select, type to filter{{color "reset"}}{{"\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- template "option" $.IterateOption $ix $option}}
  {{- end}}
{{- end }}`, colors.Accent, colors.Brand, colors.Secondary, colors.Error, colors.Default)
}
