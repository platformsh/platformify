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

type LaravelPlatformifier struct {
	*UserInput
}

func (p *LaravelPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != "laravel" {
		return fmt.Errorf("cannot platformify non-laravel stack: %s", p.Stack)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	err = fs.WalkDir(templatesFs, "templates/laravel", func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, er := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if er != nil {
			return fmt.Errorf("could not parse template: %w", er)
		}

		filePath = path.Join(cwd, filePath[len("templates/laravel"):])
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
