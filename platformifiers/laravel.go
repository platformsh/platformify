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
	"text/template"
)

//go:embed templates/laravel/**/*
var laravelTemplates embed.FS

type LaravelPlatformifier struct {
	platformify.UserInput
}

func (p *LaravelPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != "laravel" {
		return fmt.Errorf("cannot platformify non-laravel stack: %s", p.Stack)
	}
	laravel := template.New("laravel").Funcs(sprig.FuncMap())
	_, err := laravel.ParseFS(
		laravelTemplates,
		"templates/laravel/.platform.app.yaml",
		"templates/laravel/.platform/*.yaml",
	)
	if err != nil {
		return fmt.Errorf("could not parse laravel templates: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	for _, t := range laravel.Templates() {
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

func (p *LaravelPlatformifier) writeTemplate(ctx context.Context, tPath string, t *template.Template) error {
	if err := os.MkdirAll(path.Dir(tPath), os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("could not create directory %s: %w", path.Dir(tPath), err)
	}

	f, err := os.Create(tPath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %w", tPath, err)
	}

	// Antonis says we can ignore this.
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	if err := t.Execute(f, p); err != nil {
		return fmt.Errorf("could not execute template %s: %w", t.Name(), err)
	}

	return nil
}
