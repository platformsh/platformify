package platformifier

import (
	"bytes"
	"context"
	"io/fs"
	"strings"
	"text/template"

	"github.com/platformsh/platformify/vendorization"

	"github.com/Masterminds/sprig/v3"
)

func newGenericPlatformifier(templates, fileSystem fs.FS) *genericPlatformifier {
	return &genericPlatformifier{
		templates:  templates,
		fileSystem: fileSystem,
	}
}

// genericPlatformifier contains the configuration for the application to Platformify
type genericPlatformifier struct {
	templates  fs.FS
	fileSystem fs.FS
}

// Platformify will generate the needed configuration files in the current directory.
func (p *genericPlatformifier) Platformify(ctx context.Context, input *UserInput) (map[string][]byte, error) {
	assets, _ := vendorization.FromContext(ctx)
	files := make(map[string][]byte)
	err := fs.WalkDir(p.templates, ".", func(name string, d fs.DirEntry, _ error) error {
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

		files[name] = contents.Bytes()
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
