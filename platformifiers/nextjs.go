package platformifiers

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/platformsh/platformify/internal/models"
)

const nextjsTemplatesPath = "templates/nextjs"

type NextJSPlatformifier struct {
	*UserInput
}

func (p *NextJSPlatformifier) Platformify() error {
	if p.Stack != models.NextJS.String() {
		return fmt.Errorf("cannot platformify non-next.js stack: %s", p.Stack)
	}

	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	err = fs.WalkDir(templatesFs, nextjsTemplatesPath, func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, er := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if er != nil {
			return fmt.Errorf("could not parse template: %w", er)
		}

		filePath = path.Join(cwd, filePath[len(nextjsTemplatesPath):])
		if er := writeTemplate(filePath, tpl, p.UserInput); er != nil {
			return fmt.Errorf("could not write template: %w", er)
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}
