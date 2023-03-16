package writer

import (
	"github.com/platformsh/platformify/platformifiers"
	"text/template"
)

const tmplName = "platform.app.goyaml"
const tmplPath = "./templates/common/"

var blockNames = []string{"appComments", "mounts"}

// App is the main template writer which writes to .platform.app.yaml.
type App struct {
	name string
}

func (app App) Override(tmpl template.Template) (template.Template, error) {
	// We don't need to override anything yet.
	return tmpl, nil
}

func (app App) Write(pfier platformifiers.Platformifier) error {
	//TODO implement me
	panic("implement me")
}
