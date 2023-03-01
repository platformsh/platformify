package platformifiers

import (
	"context"
	"embed"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"github.com/platformsh/platformify"
	"golang.org/x/exp/slices"
	"os"
	"path"
	"strings"
)

//go:embed templates/generic/**/*
var genericTemplates embed.FS

var platformFiles = []string{
	"services.yaml",
	"routes.yaml",
}

// GenericPlatformifier contains the configuration for the application to Platformify
type GenericPlatformifier struct {
	platformify.UserInput
}

// Platformify will generate the .platform.app.yaml and .platform/ directory
func (p *GenericPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != "generic" {
		return fmt.Errorf("cannot platformify non-generic stack: %s", p.Stack)
	}
	generic := template.New("generic").Funcs(sprig.FuncMap())
	_, err := generic.ParseFS(
		genericTemplates,
		"templates/generic/.platform.app.yaml",
		"templates/generic/.platform/*.yaml",
	)
	if err != nil {
		return fmt.Errorf("could not parse generic templates: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	for _, t := range generic.Templates() {
		tPath := path.Join(cwd, t.Name())
		if slices.Contains(platformFiles, t.Name()) {
			tPath = path.Join(cwd, ".platform", t.Name())
		}
		if err := p.writeTemplate(ctx, tPath, t); err != nil {
			return fmt.Errorf("could not write template: %w", err)
		}
	}

	return nil
}

func (p *GenericPlatformifier) writeTemplate(ctx context.Context, tPath string, t *template.Template) error {
	if err := os.MkdirAll(path.Dir(tPath), os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory %s: %w", path.Dir(tPath), err)
	}

	f, err := os.Create(tPath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", tPath, err)
	}
	defer f.Close()
	if err := t.Execute(f, p); err != nil {
		return fmt.Errorf("could not execute template %s: %w", t.Name(), err)
	}

	return nil
}

func (p *GenericPlatformifier) Relationships() map[string]string {
	relationships := make(map[string]string)
	for _, service := range p.Services {
		endpoint := strings.Split(service.Type, ":")[0]
		relationships[service.Name] = fmt.Sprintf("%s:%s", service.Name, endpoint)
	}
	return relationships
}
