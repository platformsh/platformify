package platformifier

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

func newGenericPlatformifier(templates fs.FS, file fileCreator) *genericPlatformifier {
	return &genericPlatformifier{
		templates: templates,
		file:      file,
	}
}

// genericPlatformifier contains the configuration for the application to Platformify
type genericPlatformifier struct {
	templates fs.FS
	file      fileCreator
}

// Platformify will generate the .platformifiers.app.yaml and .platformifiers/ directory
func (p *genericPlatformifier) Platformify(_ context.Context, input *UserInput) error {
	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}

	err = fs.WalkDir(p.templates, ".", func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl := template.Must(template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(p.templates, filePath))

		filePath = path.Join(cwd, filePath)
		f, writeErr := p.file.Create(filePath)
		if writeErr != nil {
			return fmt.Errorf("could not write template: %w", writeErr)
		}
		defer f.Close()

		return tpl.Execute(f, input)
	})
	if err != nil {
		return err
	}

	return nil
}
