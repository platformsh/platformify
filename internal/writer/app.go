package writer

import (
	"fmt"
	"github.com/platformsh/platformify/platformifiers"
	"os"
	"text/template"
)

const filename = ".platform.app.yaml"
const tmplName = "platform.app.goyaml"
const tmplPath = "./templates/common/"

var blockNames = []string{"appComments", "mounts"}

// App is the main template writer which writes to .dot_platform.app.yaml.
type App struct {
	name string
}

func (app App) Override(tmpl *template.Template) (*template.Template, error) {
	// We don't need to override anything yet.
	return tmpl, nil
}

func (app App) Write(pfier platformifiers.Platformifier) error {
	//cwd, err := os.Getwd()
	//if err != nil {
	//	return fmt.Errorf("could not get current working directory: %w", err)
	//}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create '%s' at path '%s'", filename)
	}
	defer file.Close()

	// Parse the template.
	tmpl, err := Parse(tmplPath, tmplName)
	if err != nil {
		return fmt.Errorf("parse error: %s", err)
	}

	// Check for overrides.
	tmpl, err = app.Override(tmpl)
	if err != nil {
		return fmt.Errorf("override error: %s", err)
	}

	// Write the file.
	if err = tmpl.Execute(file, pfier.GetPshConfig()); err != nil {
		return fmt.Errorf("could not write file: %s", filename)
	}

	return file.Close()
}
