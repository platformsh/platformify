package platformifiers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

// GenericPlatformifier contains the configuration for the application to Platformify
type GenericPlatformifier struct {
	*UserInput
}

// Platformify will generate the .platform.app.yaml and .platform/ directory
func (p *GenericPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != "generic" {
		return fmt.Errorf("cannot platformify non-generic stack: %s", p.Stack)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	err = fs.WalkDir(templatesFs, "templates/generic", func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, er := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if er != nil {
			return fmt.Errorf("could not parse template: %w", er)
		}

		filePath = path.Join(cwd, filePath[len("templates/generic"):])
		if er := writeTemplate(ctx, filePath, tpl, p.UserInput); er != nil {
			return fmt.Errorf("could not write template: %w", er)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
