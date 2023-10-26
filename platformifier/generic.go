package platformifier

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"text/template"

	"github.com/platformsh/platformify/vendorization"

	"github.com/Masterminds/sprig/v3"
)

func newGenericPlatformifier(templates fs.FS, fileSystem FS) *genericPlatformifier {
	return &genericPlatformifier{
		templates:  templates,
		fileSystem: fileSystem,
	}
}

// genericPlatformifier contains the configuration for the application to Platformify
type genericPlatformifier struct {
	templates  fs.FS
	fileSystem FS
}

// Platformify will generate the needed configuration files in the current directory.
func (p *genericPlatformifier) Platformify(ctx context.Context, input *UserInput) error {
	assets, _ := vendorization.FromContext(ctx)
	err := fs.WalkDir(p.templates, ".", func(name string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl := template.Must(template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(p.templates, name))
		contents := &bytes.Buffer{}
		if err := tpl.Execute(contents, templateData{input, assets}); err != nil {
			return err
		}

		// Skip empty files
		if strings.TrimSpace(contents.String()) == "" {
			return nil
		}

		f, writeErr := p.fileSystem.Create(name)
		if writeErr != nil {
			return fmt.Errorf("could not write template: %w", writeErr)
		}
		defer f.Close()

		if _, err := io.Copy(f, contents); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
